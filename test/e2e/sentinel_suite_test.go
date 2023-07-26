package e2e

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	redisclient "github.com/3scale/saas-operator/pkg/redis_v2/client"
	testutil "github.com/3scale/saas-operator/test/util"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
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
				ObjectMeta: metav1.ObjectMeta{Name: "rs0", Namespace: ns},
				Spec:       saasv1alpha1.RedisShardSpec{MasterIndex: pointer.Int32(0), SlaveCount: pointer.Int32(2)},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Name: "rs1", Namespace: ns},
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

				sclient, stopCh, err := testutil.SentinelClient(cfg,
					types.NamespacedName{
						Name:      fmt.Sprintf("redis-sentinel-%d", rand.Intn(int(saasv1alpha1.SentinelDefaultReplicas))),
						Namespace: ns,
					})
				Expect(err).ToNot(HaveOccurred())
				defer close(stopCh)

				ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
				defer cancel()

				masters, err := sclient.SentinelMasters(ctx)
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

			Eventually(func() error {

				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: sentinel.GetName(), Namespace: ns}, &sentinel)
				if err != nil {
					return err
				}

				for i, shard := range shards {

					if diff := cmp.Diff(sentinel.Status.MonitoredShards[i],
						saasv1alpha1.MonitoredShard{
							Name: shard.GetName(),
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
						return fmt.Errorf("got unexpected sentinel status %s", diff)
					}
				}

				return nil
			}, timeout, poll).ShouldNot(HaveOccurred())
		})

		When("a redis master is unavailable", func() {

			BeforeEach(func() {

				By("configuring shard0's slave 2 with priority '0' to ensure always the same slave gets promoted to master", func() {

					rclient, stopCh, err := testutil.RedisClient(cfg,
						types.NamespacedName{
							Name:      "redis-shard-rs0-2",
							Namespace: ns,
						})
					Expect(err).ToNot(HaveOccurred())
					defer close(stopCh)

					ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
					defer cancel()

					err = rclient.RedisConfigSet(ctx, "slave-priority", "0")
					Expect(err).ToNot(HaveOccurred())
				})

				By("making the shard0's master unavailable to force a failover", func() {

					go func() {
						defer GinkgoRecover()

						rclient, stopCh, err := testutil.RedisClient(cfg,
							types.NamespacedName{
								Name:      "redis-shard-rs0-0",
								Namespace: ns,
							})
						Expect(err).ToNot(HaveOccurred())
						defer close(stopCh)

						ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
						defer cancel()

						// a master is considered down after 5s, so we sleep the current master
						// for 10 seconds to simulate a failure and trigger a master failover
						rclient.RedisDebugSleep(ctx, 10*time.Second)
					}()
				})

			})

			It("triggers a failover", func() {

				rclient, stopCh, err := testutil.RedisClient(cfg,
					types.NamespacedName{
						Name:      "redis-shard-rs0-1",
						Namespace: ns,
					})
				Expect(err).ToNot(HaveOccurred())
				defer close(stopCh)

				Eventually(func() error {
					role, _, err := rclient.RedisRole(context.TODO())
					if err != nil {
						return err
					}
					if role != redisclient.Master {
						return fmt.Errorf("expected 'master' but got %s", role)
					}
					return nil
				}, timeout, poll).ShouldNot(HaveOccurred())

			})

			It("updates the status appropriately", func() {
				Eventually(func() error {

					err := k8sClient.Get(context.Background(), types.NamespacedName{Name: sentinel.GetName(), Namespace: ns}, &sentinel)
					if err != nil {
						return err
					}

					if diff := cmp.Diff(sentinel.Status.MonitoredShards[0],
						saasv1alpha1.MonitoredShard{
							Name: shards[0].GetName(),
							Servers: map[string]saasv1alpha1.RedisServerDetails{
								// old master is now a slave
								strings.TrimPrefix(*shards[0].Status.ShardNodes.Master, "redis://"): {
									Role:   redisclient.Slave,
									Config: map[string]string{"save": "900 1 300 10", "slave-read-only": "yes"},
								},
								// first slave is now the master
								strings.TrimPrefix(shards[0].Status.ShardNodes.Slaves[0], "redis://"): {
									Role:   redisclient.Master,
									Config: map[string]string{"save": "900 1 300 10"},
								},
								strings.TrimPrefix(shards[0].Status.ShardNodes.Slaves[1], "redis://"): {
									Role:   redisclient.Slave,
									Config: map[string]string{"save": "900 1 300 10", "slave-read-only": "yes"},
								},
							},
						},
					); diff != "" {
						return fmt.Errorf("got unexpected sentinel status %s", diff)
					}

					return nil
				}, timeout, poll).ShouldNot(HaveOccurred())
			})
		})

	})

	AfterEach(func() {

		// Delete sentinel
		err := k8sClient.Delete(context.Background(), &sentinel, client.PropagationPolicy(metav1.DeletePropagationForeground))
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

})
