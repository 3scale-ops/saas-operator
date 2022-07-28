package e2e

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/redis"
	redisclient "github.com/3scale/saas-operator/pkg/redis/crud/client"
	testutil "github.com/3scale/saas-operator/test/util"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("sentinel e2e suite", func() {
	var ns string
	var shards []saasv1alpha1.RedisShard
	var sentinel saasv1alpha1.Sentinel

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

		shards = []saasv1alpha1.RedisShard{
			{
				ObjectMeta: metav1.ObjectMeta{Name: "rs1", Namespace: ns},
				Spec:       saasv1alpha1.RedisShardSpec{MasterIndex: pointer.Int32(0), SlaveCount: pointer.Int32(2)},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Name: "rs2", Namespace: ns},
				Spec:       saasv1alpha1.RedisShardSpec{MasterIndex: pointer.Int32(2), SlaveCount: pointer.Int32(2)},
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
					return nil
				} else {
					return fmt.Errorf("RedisShard %s not ready", shard.ObjectMeta.Name)
				}

			}, timeout, poll).ShouldNot(HaveOccurred())
		}
	})

	When("Sentinel resource is created and ready", func() {

		BeforeEach(func() {
			sentinel = saasv1alpha1.Sentinel{
				ObjectMeta: metav1.ObjectMeta{Name: "sentinel", Namespace: ns},
				Spec: saasv1alpha1.SentinelSpec{
					Config: &saasv1alpha1.SentinelConfig{
						MonitoredShards: map[string][]string{
							shards[0].GetName(): {
								*shards[0].Status.ShardNodes.Master,
								shards[0].Status.ShardNodes.Slaves[0],
								shards[0].Status.ShardNodes.Slaves[1],
							},
							shards[1].GetName(): {
								*shards[1].Status.ShardNodes.Master,
								shards[1].Status.ShardNodes.Slaves[0],
								shards[1].Status.ShardNodes.Slaves[1],
							},
						},
					},
				},
			}

			err := k8sClient.Create(context.Background(), &sentinel)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() error {

				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: sentinel.GetName(), Namespace: ns}, &sentinel)
				Expect(err).ToNot(HaveOccurred())

				if len(sentinel.Status.MonitoredShards) != len(shards) {
					return fmt.Errorf("sentinel not ready")
				}
				return nil
			}, timeout, poll).ShouldNot(HaveOccurred())
		})

		It("deploys sentinel Pods that monitor each of the redis shards", func() {

			By("issuing a 'sentinel masters' command to ensure all shards are monitored by sentinel", func() {

				sentinelPod := fmt.Sprintf("redis-sentinel-%d", rand.Intn(int(saasv1alpha1.SentinelDefaultReplicas))-1)
				localPort, stopCh, err := testutil.PortForward(cfg, types.NamespacedName{Name: sentinelPod, Namespace: ns}, saasv1alpha1.SentinelPort)
				Expect(err).ToNot(HaveOccurred())
				defer close(stopCh)

				ss, err := redis.NewSentinelServerFromConnectionString(sentinelPod, fmt.Sprintf("redis://localhost:%d", localPort))
				Expect(err).ToNot(HaveOccurred())

				ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
				defer cancel()

				masters, err := ss.CRUD.SentinelMasters(ctx)
				Expect(err).ToNot(HaveOccurred())

				for _, shard := range shards {
					found := false
					for _, master := range masters {
						if strings.Contains(*shard.Status.ShardNodes.Master, master.IP) {
							found = true
							break
						}
					}
					if found == false {
						Fail(fmt.Sprintf("master for shard %s not found in 'sentinel master' response", shard.GetName()))
					}
				}
			})
		})

		It("updates the resource status with the current status of each redis shard", func() {

			for i, shard := range shards {

				if diff := cmp.Diff(sentinel.Status.MonitoredShards[i],
					saasv1alpha1.MonitoredShard{
						Name:   shard.GetName(),
						Master: "",
						Servers: map[string]saasv1alpha1.RedisServerDetails{
							strings.TrimPrefix(*shard.Status.ShardNodes.Master, "redis://"): {
								Role:   redisclient.Master,
								Config: map[string]string{"save": "900 1 300 10"},
							},
							strings.TrimPrefix(shard.Status.ShardNodes.Slaves[0], "redis://"): {
								Role:   redisclient.Slave,
								Config: map[string]string{"save": "900 1 300 10", "slave-read-only": "yes"},
							},
							strings.TrimPrefix(shard.Status.ShardNodes.Slaves[1], "redis://"): {
								Role:   redisclient.Slave,
								Config: map[string]string{"save": "900 1 300 10", "slave-read-only": "yes"},
							},
						},
					},
				); diff != "" {
					Fail(fmt.Sprintf("got unexpected sentinel status %s", diff))
				}
			}
		})

	})

	AfterEach(func() {

		// Delete redis shards
		for _, shard := range shards {
			err := k8sClient.Delete(context.Background(), &shard, client.PropagationPolicy(metav1.DeletePropagationForeground))
			Expect(err).ToNot(HaveOccurred())
		}

		// Delete the namespace
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}
		err := k8sClient.Delete(context.Background(), ns, client.PropagationPolicy(metav1.DeletePropagationForeground))
		Expect(err).ToNot(HaveOccurred())
	})

})
