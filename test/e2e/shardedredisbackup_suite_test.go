package e2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/3scale-ops/basereconciler/util"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	testutil "github.com/3scale-ops/saas-operator/test/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	awsCredentials = "aws-credentials"
	sshPrivateKey  = "redis-backup-ssh-private-key"
	bucketName     = "backups"
	backupsPath    = "redis"
	minioNamespace = "minio"
)

var _ = Describe("shardedredisbackup e2e suite", func() {
	var ns string
	var shards []saasv1alpha1.RedisShard
	var sentinel saasv1alpha1.Sentinel
	var backup saasv1alpha1.ShardedRedisBackup

	BeforeEach(func() {
		// Create a namespace for each block
		ns = "test-ns-" + nameGenerator.Generate()

		// Add any setup steps that needs to be executed before each test
		testNamespace := &corev1.Namespace{
			TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Namespace"},
			ObjectMeta: metav1.ObjectMeta{Name: ns},
		}

		err := k8sClient.Create(context.Background(), testNamespace)
		Expect(err).ToNot(HaveOccurred())

		n := &corev1.Namespace{}
		Eventually(func() error {
			return k8sClient.Get(context.Background(), types.NamespacedName{Name: ns}, n)
		}, timeout, poll).ShouldNot(HaveOccurred())

		// create redis shards
		shards = []saasv1alpha1.RedisShard{
			{
				ObjectMeta: metav1.ObjectMeta{Name: "rs0", Namespace: ns},
				Spec: saasv1alpha1.RedisShardSpec{
					MasterIndex: util.Pointer[int32](0),
					SlaveCount:  util.Pointer[int32](2),
					Command:     util.Pointer("/entrypoint.sh"),
					Image: &saasv1alpha1.ImageSpec{
						Name: util.Pointer("redis-with-ssh"),
						Tag:  util.Pointer("6.2.13-alpine"),
					},
				},
			},
		}

		for i, shard := range shards {
			err = k8sClient.Create(context.Background(), &shard)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() error {
				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: shard.GetName(), Namespace: ns}, &shard)
				if err != nil {
					return err
				}
				if shard.Status.ShardNodes != nil && shard.Status.ShardNodes.Master != nil {
					// store the resource for later use
					shards[i] = shard
					GinkgoWriter.Printf("[debug] Shard %s topology: %+v\n", shard.GetName(), *shard.Status.ShardNodes)
					return nil
				} else {
					return fmt.Errorf("RedisShard %s not ready", shard.ObjectMeta.Name)
				}

			}, timeout, poll).ShouldNot(HaveOccurred())
		}

		// create sentinel
		sentinel = saasv1alpha1.Sentinel{
			ObjectMeta: metav1.ObjectMeta{Name: "sentinel", Namespace: ns},
			Spec: saasv1alpha1.SentinelSpec{
				Replicas: util.Pointer(int32(1)),
				Config: &saasv1alpha1.SentinelConfig{
					MonitoredShards: map[string][]string{
						shards[0].GetName(): {
							"redis://" + shards[0].Status.ShardNodes.GetHostPortByPodIndex(0),
							"redis://" + shards[0].Status.ShardNodes.GetHostPortByPodIndex(1),
							"redis://" + shards[0].Status.ShardNodes.GetHostPortByPodIndex(2),
						},
					},
				},
			},
		}

		err = k8sClient.Create(context.Background(), &sentinel)
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() error {

			err := k8sClient.Get(context.Background(), types.NamespacedName{Name: sentinel.GetName(), Namespace: ns}, &sentinel)
			Expect(err).ToNot(HaveOccurred())

			if len(sentinel.Status.MonitoredShards) != len(shards) {
				return fmt.Errorf("sentinel not ready")
			}
			return nil
		}, timeout, poll).ShouldNot(HaveOccurred())

		// Load a test dataset into redis
		rclient, stopCh, err := testutil.RedisClient(cfg,
			types.NamespacedName{
				Name:      "redis-shard-rs0-0",
				Namespace: ns,
			})
		Expect(err).ToNot(HaveOccurred())
		defer close(stopCh)

		dir, _ := os.Getwd()
		// got the dataset from https://redis.com/blog/datasets-for-test-databases/
		err = testutil.LoadRedisDataset(context.Background(), rclient, filepath.Join(dir, "../assets/redis-datasets/supernovas.csv"))
		Expect(err).ToNot(HaveOccurred())

		// copy over required credentials from default namespace
		for _, creds := range []string{awsCredentials, sshPrivateKey} {
			secret := &corev1.Secret{}
			err := k8sClient.Get(context.Background(), types.NamespacedName{Name: creds, Namespace: "default"}, secret)
			Expect(err).ToNot(HaveOccurred())

			secret.ObjectMeta = metav1.ObjectMeta{Name: creds, Namespace: ns}
			err = k8sClient.Create(context.Background(), secret)
			Expect(err).ToNot(HaveOccurred())
		}

		// Create a shardedredisbackup resource
		backup = saasv1alpha1.ShardedRedisBackup{
			ObjectMeta: metav1.ObjectMeta{Name: "backup", Namespace: ns},
			Spec: saasv1alpha1.ShardedRedisBackupSpec{
				SentinelRef: sentinel.GetName(),
				Schedule:    "* * * * *",
				DBFile:      "/data/dump.rdb",
				SSHOptions: saasv1alpha1.SSHOptions{
					User: "docker",
					PrivateKeySecretRef: corev1.LocalObjectReference{
						Name: "redis-backup-ssh-private-key",
					},
					Port: util.Pointer(uint32(2222)),
					Sudo: util.Pointer(true),
				},
				S3Options: saasv1alpha1.S3Options{
					Bucket: bucketName,
					Path:   backupsPath,
					Region: "us-east-1",
					CredentialsSecretRef: corev1.LocalObjectReference{
						Name: "aws-credentials",
					},
					ServiceEndpoint: util.Pointer(fmt.Sprintf("http://minio.%s.svc.cluster.local:9000", minioNamespace)),
				},
				PollInterval: &metav1.Duration{Duration: 1 * time.Second},
			},
		}

		err = k8sClient.Create(context.Background(), &backup)
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() error {

			err := k8sClient.Get(context.Background(), types.NamespacedName{Name: backup.GetName(), Namespace: ns}, &backup)
			Expect(err).ToNot(HaveOccurred())

			if len(backup.Status.Backups) == 0 {
				msg := "[debug] waiting for backup to be scheduled"
				GinkgoWriter.Println(msg)
				return fmt.Errorf(msg)
			}
			return nil
		}, timeout, poll).ShouldNot(HaveOccurred())

	})

	AfterEach(func() {

		// Delete backup
		err := k8sClient.Delete(context.Background(), &backup, client.PropagationPolicy(metav1.DeletePropagationForeground))
		Expect(err).ToNot(HaveOccurred())

		// Delete sentinel
		err = k8sClient.Delete(context.Background(), &sentinel, client.PropagationPolicy(metav1.DeletePropagationForeground))
		Expect(err).ToNot(HaveOccurred())

		// Delete redis shards
		for _, shard := range shards {
			err := k8sClient.Delete(context.Background(), &shard, client.PropagationPolicy(metav1.DeletePropagationForeground))
			Expect(err).ToNot(HaveOccurred())
		}

		// Delete the namespace
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}
		err = k8sClient.Delete(context.Background(), ns, client.PropagationPolicy(metav1.DeletePropagationForeground))
		Expect(err).ToNot(HaveOccurred())
	})

	It("runs a backup that completes successfully", func() {
		var backupResult saasv1alpha1.BackupStatus

		Eventually(func() error {

			err := k8sClient.Get(context.Background(), types.NamespacedName{Name: backup.GetName(), Namespace: ns}, &backup)
			Expect(err).ToNot(HaveOccurred())
			backupResult = backup.Status.Backups[len(backup.Status.Backups)-1]

			switch backupResult.State {
			case saasv1alpha1.BackupPendingState:
				GinkgoWriter.Printf("[debug %s] backup has not yet started\n", time.Now())
				return fmt.Errorf("")
			case saasv1alpha1.BackupRunningState:
				GinkgoWriter.Printf("[debug %s] backup is running\n", time.Now())
				return fmt.Errorf("")
			case saasv1alpha1.BackupCompletedState:
				GinkgoWriter.Printf("[debug %s] backup completed successfully\n", time.Now())
				return nil
			default:
				GinkgoWriter.Printf("[debug %s] backup failed: '%s'\n", time.Now(), backupResult.Message)
				return fmt.Errorf(backupResult.Message)
			}

		}, timeout, poll).ShouldNot(HaveOccurred())

		By("checking that the backup is actually where the status says it is and has the reported size")
		ctx := context.Background()
		list := &corev1.PodList{}
		err := k8sClient.List(context.Background(), list,
			client.InNamespace(minioNamespace),
			client.MatchingLabels{"app": "minio"})
		Expect(err).ToNot(HaveOccurred())
		Expect(len(list.Items)).To(Equal(1))

		s3client, stopCh, err := testutil.MinioClient(ctx, cfg, client.ObjectKeyFromObject(&list.Items[0]), "admin", "admin123")
		Expect(err).ToNot(HaveOccurred())
		defer close(stopCh)

		result, err := s3client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(strings.TrimPrefix(*backupResult.BackupFile, fmt.Sprintf("s3://%s/", bucketName))),
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(*result.ContentLength).To(Equal(*backupResult.BackupSize))
	})
})
