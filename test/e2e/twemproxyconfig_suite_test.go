package e2e

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/3scale-ops/basereconciler/util"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/twemproxy"
	testutil "github.com/3scale-ops/saas-operator/test/util"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

var _ = Describe("twemproxyconfig e2e suite", func() {
	var ns string
	var shards []saasv1alpha1.RedisShard
	var sentinel saasv1alpha1.Sentinel
	var twemproxyconfig saasv1alpha1.TwemproxyConfig

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

		Eventually(func() error {
			return k8sClient.Get(context.Background(), client.ObjectKeyFromObject(testNamespace), testNamespace)
		}, timeout, poll).ShouldNot(HaveOccurred())

		shards = []saasv1alpha1.RedisShard{
			{
				ObjectMeta: metav1.ObjectMeta{Name: "rs0", Namespace: ns},
				Spec:       saasv1alpha1.RedisShardSpec{MasterIndex: util.Pointer[int32](0), SlaveCount: util.Pointer[int32](2)},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Name: "rs1", Namespace: ns},
				Spec:       saasv1alpha1.RedisShardSpec{MasterIndex: util.Pointer[int32](0), SlaveCount: util.Pointer[int32](2)},
			},
		}

		for i, shard := range shards {
			err = k8sClient.Create(context.Background(), &shard)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() error {
				err := k8sClient.Get(context.Background(), client.ObjectKeyFromObject(&shard), &shard)
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

		sentinel = saasv1alpha1.Sentinel{
			ObjectMeta: metav1.ObjectMeta{Name: "sentinel", Namespace: ns},
			Spec: saasv1alpha1.SentinelSpec{
				Config: &saasv1alpha1.SentinelConfig{
					ClusterTopology: map[string]map[string]string{
						shards[0].GetName(): {
							"redis-shard-rs0-0": "redis://" + shards[0].Status.ShardNodes.GetHostPortByPodIndex(0), // master
							"redis-shard-rs0-1": "redis://" + shards[0].Status.ShardNodes.GetHostPortByPodIndex(1),
							"redis-shard-rs0-2": "redis://" + shards[0].Status.ShardNodes.GetHostPortByPodIndex(2),
						},
						shards[1].GetName(): {
							"redis-shard-rs1-0": "redis://" + shards[1].Status.ShardNodes.GetHostPortByPodIndex(0), // master
							"redis-shard-rs1-1": "redis://" + shards[1].Status.ShardNodes.GetHostPortByPodIndex(1),
							"redis-shard-rs1-2": "redis://" + shards[1].Status.ShardNodes.GetHostPortByPodIndex(2),
						},
					},
				},
			},
		}

		err = k8sClient.Create(context.Background(), &sentinel)
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() error {

			err := k8sClient.Get(context.Background(), client.ObjectKeyFromObject(&sentinel), &sentinel)
			Expect(err).ToNot(HaveOccurred())

			if len(sentinel.Status.MonitoredShards) != len(shards) {
				return fmt.Errorf("sentinel not ready")
			}
			return nil
		}, timeout, poll).ShouldNot(HaveOccurred())

	})

	When("TwemproxyConfig resource is created targeting redis masters", func() {

		BeforeEach(func() {

			By("creating the TwemproxyConfig resource targeting masters", func() {
				twemproxyconfig = saasv1alpha1.TwemproxyConfig{
					ObjectMeta: metav1.ObjectMeta{Name: "tmc-instance", Namespace: ns},
					Spec: saasv1alpha1.TwemproxyConfigSpec{
						ServerPools: []saasv1alpha1.TwemproxyServerPool{{
							Name: "test-pool",
							Topology: []saasv1alpha1.ShardedRedisTopology{
								{ShardName: "l-shard00", PhysicalShard: shards[0].GetName()},
								{ShardName: "l-shard01", PhysicalShard: shards[0].GetName()},
								{ShardName: "l-shard02", PhysicalShard: shards[0].GetName()},
								{ShardName: "l-shard03", PhysicalShard: shards[1].GetName()},
								{ShardName: "l-shard04", PhysicalShard: shards[1].GetName()},
							},
							BindAddress: "0.0.0.0:22121",
							Timeout:     5000,
							TCPBacklog:  512,
							PreConnect:  false,
						}},
						GrafanaDashboard: &saasv1alpha1.GrafanaDashboardSpec{},
					},
				}

				err := k8sClient.Create(context.Background(), &twemproxyconfig)
				Expect(err).ToNot(HaveOccurred())
			})

			By("wating until the TwemproxyConfig resource is ready", func() {
				Eventually(func() error {
					if err := k8sClient.Get(context.Background(), client.ObjectKeyFromObject(&twemproxyconfig), &twemproxyconfig); err != nil {
						return err
					}
					if twemproxyconfig.Status.SelectedTargets == nil {
						return fmt.Errorf("status.targets is empty")
					}
					return nil
				}, timeout, poll).ShouldNot(HaveOccurred())
			})
		})

		It("deploys a ConfigMap with twemproxy configuration that points to redis masters", func() {
			Eventually(assertTwemproxyConfigStatus(&twemproxyconfig, &sentinel,
				&saasv1alpha1.TwemproxyConfigStatus{
					SelectedTargets: map[string]saasv1alpha1.TargetServer{
						shards[0].GetName(): {
							ServerAlias:   util.Pointer(shards[0].Status.ShardNodes.GetAliasByPodIndex(0)),
							ServerAddress: shards[0].Status.ShardNodes.GetHostPortByPodIndex(0),
						},
						shards[1].GetName(): {
							ServerAlias:   util.Pointer(shards[1].Status.ShardNodes.GetAliasByPodIndex(0)),
							ServerAddress: shards[1].Status.ShardNodes.GetHostPortByPodIndex(0),
						},
					},
				}), timeout, poll).Should(Not(HaveOccurred()))

			Eventually(assertTwemproxyConfigServerPool(&twemproxyconfig,
				[]twemproxy.Server{
					{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard00"},
					{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard01"},
					{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard02"},
					{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard03"},
					{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard04"},
				}), timeout, poll).Should(Not(HaveOccurred()))

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

			It("updates the twemproxy configuration with the new master", func() {

				Eventually(assertTwemproxyConfigStatus(&twemproxyconfig, &sentinel,
					&saasv1alpha1.TwemproxyConfigStatus{
						SelectedTargets: map[string]saasv1alpha1.TargetServer{
							shards[0].GetName(): {
								ServerAlias:   util.Pointer(shards[0].Status.ShardNodes.GetAliasByPodIndex(1)),
								ServerAddress: shards[0].Status.ShardNodes.GetHostPortByPodIndex(1),
							},
							shards[1].GetName(): {
								ServerAlias:   util.Pointer(shards[1].Status.ShardNodes.GetAliasByPodIndex(0)),
								ServerAddress: shards[1].Status.ShardNodes.GetHostPortByPodIndex(0),
							},
						},
					}), timeout, poll).Should(Not(HaveOccurred()))

				Eventually(assertTwemproxyConfigServerPool(&twemproxyconfig,
					[]twemproxy.Server{
						{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(1), Priority: 1, Name: "l-shard00"},
						{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(1), Priority: 1, Name: "l-shard01"},
						{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(1), Priority: 1, Name: "l-shard02"},
						{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard03"},
						{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard04"},
					}), timeout, poll).Should(Not(HaveOccurred()))
			})
		})

		AfterEach(func() {
			// Delete twemproxyconfig
			err := k8sClient.Delete(context.Background(), &twemproxyconfig, client.PropagationPolicy(metav1.DeletePropagationForeground))
			Expect(err).ToNot(HaveOccurred())
		})

	})

	When("TwemproxyConfig resource is created targeting redis rw-slaves", func() {

		BeforeEach(func() {

			By("configuring rw-slaves in each shard", func() {

				for _, slave := range []string{"redis-shard-rs0-2", "redis-shard-rs1-2"} {
					rclient, stopCh, err := testutil.RedisClient(cfg,
						types.NamespacedName{
							Name:      slave,
							Namespace: ns,
						})
					Expect(err).ToNot(HaveOccurred())
					defer close(stopCh)
					defer rclient.CloseClient()

					ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
					defer cancel()

					err = rclient.RedisConfigSet(ctx, "slave-read-only", "no")
					Expect(err).ToNot(HaveOccurred())
				}

			})

			By("creating the TwemproxyConfig resource targeting rw-slaves", func() {
				twemproxyconfig = saasv1alpha1.TwemproxyConfig{
					ObjectMeta: metav1.ObjectMeta{Name: "tmc-instance", Namespace: ns},
					Spec: saasv1alpha1.TwemproxyConfigSpec{
						ServerPools: []saasv1alpha1.TwemproxyServerPool{{
							Name:   "test-pool",
							Target: func() *saasv1alpha1.TargetRedisServers { t := saasv1alpha1.SlavesRW; return &t }(),
							Topology: []saasv1alpha1.ShardedRedisTopology{
								{ShardName: "l-shard00", PhysicalShard: shards[0].GetName()},
								{ShardName: "l-shard01", PhysicalShard: shards[0].GetName()},
								{ShardName: "l-shard02", PhysicalShard: shards[0].GetName()},
								{ShardName: "l-shard03", PhysicalShard: shards[1].GetName()},
								{ShardName: "l-shard04", PhysicalShard: shards[1].GetName()},
							},
							BindAddress: "0.0.0.0:22121",
							Timeout:     5000,
							TCPBacklog:  512,
							PreConnect:  false,
						}},
						GrafanaDashboard: &saasv1alpha1.GrafanaDashboardSpec{},
					},
				}

				err := k8sClient.Create(context.Background(), &twemproxyconfig)
				Expect(err).ToNot(HaveOccurred())
			})

			By("wating until the TwemproxyConfig resource is ready", func() {
				Eventually(func() error {
					if err := k8sClient.Get(context.Background(), client.ObjectKeyFromObject(&twemproxyconfig), &twemproxyconfig); err != nil {
						return err
					}
					if twemproxyconfig.Status.SelectedTargets == nil {
						return fmt.Errorf("status.targets is empty")
					}
					return nil
				}, timeout, poll).ShouldNot(HaveOccurred())
			})
		})

		It("deploys a ConfigMap with twemproxy configuration that points to redis rw-slaves", func() {

			Eventually(assertTwemproxyConfigStatus(&twemproxyconfig, &sentinel,
				&saasv1alpha1.TwemproxyConfigStatus{
					SelectedTargets: map[string]saasv1alpha1.TargetServer{
						shards[0].GetName(): {
							ServerAlias:   util.Pointer(shards[0].Status.ShardNodes.GetAliasByPodIndex(2)),
							ServerAddress: shards[0].Status.ShardNodes.GetHostPortByPodIndex(2),
						},
						shards[1].GetName(): {
							ServerAlias:   util.Pointer(shards[1].Status.ShardNodes.GetAliasByPodIndex(2)),
							ServerAddress: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2),
						},
					},
				}), timeout, poll).Should(Not(HaveOccurred()))

			Eventually(assertTwemproxyConfigServerPool(&twemproxyconfig,
				[]twemproxy.Server{
					{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard00"},
					{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard01"},
					{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard02"},
					{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard03"},
					{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard04"},
				}), timeout, poll).Should(Not(HaveOccurred()))
		})

		When("there are no rw-slaves available in a shard it failovers to the master", func() {

			BeforeEach(func() {
				By("simulating a failure in 'redis-shard-rs0-2' rw-slave", func() {

					go func() {
						defer GinkgoRecover()

						rclient, stopCh, err := testutil.RedisClient(cfg,
							types.NamespacedName{
								Name:      "redis-shard-rs0-2",
								Namespace: ns,
							})
						Expect(err).ToNot(HaveOccurred())
						defer close(stopCh)
						defer rclient.CloseClient()

						rclient.RedisDebugSleep(context.TODO(), 10*time.Second)
					}()
				})
			})

			It("reconfigures shard rs0 to point to the master", func() {

				By("checking the config for rs0 points to master", func() {

					Eventually(assertTwemproxyConfigStatus(&twemproxyconfig, &sentinel,
						&saasv1alpha1.TwemproxyConfigStatus{
							SelectedTargets: map[string]saasv1alpha1.TargetServer{
								shards[0].GetName(): {
									ServerAlias:   util.Pointer(shards[0].Status.ShardNodes.GetAliasByPodIndex(0)),
									ServerAddress: shards[0].Status.ShardNodes.GetHostPortByPodIndex(0),
								},
								shards[1].GetName(): {
									ServerAlias:   util.Pointer(shards[1].Status.ShardNodes.GetAliasByPodIndex(2)),
									ServerAddress: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2),
								},
							},
						}), timeout, poll).Should(Not(HaveOccurred()))

					Eventually(assertTwemproxyConfigServerPool(&twemproxyconfig,
						[]twemproxy.Server{
							{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard00"},
							{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard01"},
							{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard02"},
							{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard03"},
							{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard04"},
						}), timeout, poll).Should(Not(HaveOccurred()))
				})

				By("checking the config for rs0 points back to rw-slave once it's recovered", func() {

					Eventually(assertTwemproxyConfigStatus(&twemproxyconfig, &sentinel,
						&saasv1alpha1.TwemproxyConfigStatus{
							SelectedTargets: map[string]saasv1alpha1.TargetServer{
								shards[0].GetName(): {
									ServerAlias:   util.Pointer(shards[0].Status.ShardNodes.GetAliasByPodIndex(2)),
									ServerAddress: shards[0].Status.ShardNodes.GetHostPortByPodIndex(2),
								},
								shards[1].GetName(): {
									ServerAlias:   util.Pointer(shards[1].Status.ShardNodes.GetAliasByPodIndex(2)),
									ServerAddress: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2),
								},
							},
						}), timeout, poll).Should(Not(HaveOccurred()))

					Eventually(assertTwemproxyConfigServerPool(&twemproxyconfig,
						[]twemproxy.Server{
							{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard00"},
							{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard01"},
							{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard02"},
							{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard03"},
							{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard04"},
						}), timeout, poll).Should(Not(HaveOccurred()))

				})

			})
		})

		When("when there are several rw-slaves, chooses the first in alphabetical order (by address)", func() {

			BeforeEach(func() {
				By("configuring 'redis-shard-rs0-1' also as a rw-slave", func() {

					rclient, stopCh, err := testutil.RedisClient(cfg,
						types.NamespacedName{
							Name:      "redis-shard-rs0-1",
							Namespace: ns,
						})
					Expect(err).ToNot(HaveOccurred())
					defer close(stopCh)
					defer rclient.CloseClient()

					ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
					defer cancel()

					err = rclient.RedisConfigSet(ctx, "slave-read-only", "no")
					Expect(err).ToNot(HaveOccurred())
				})
			})

			It("reconfigures shard rs0 to point to the first slave in alphabetical order (by hostport)", func() {

				By("checking the twemproxy config for shard rs0", func() {

					// determine which should be the chosen rw-slave
					addresses := []string{
						shards[0].Status.ShardNodes.GetHostPortByPodIndex(1),
						shards[0].Status.ShardNodes.GetHostPortByPodIndex(2),
					}
					sort.Strings(addresses)
					idx := shards[0].Status.ShardNodes.GetIndexByHostPort(addresses[0])

					Eventually(assertTwemproxyConfigStatus(&twemproxyconfig, &sentinel,
						&saasv1alpha1.TwemproxyConfigStatus{
							SelectedTargets: map[string]saasv1alpha1.TargetServer{
								shards[0].GetName(): {
									ServerAlias:   util.Pointer(shards[0].Status.ShardNodes.GetAliasByPodIndex(idx)),
									ServerAddress: shards[0].Status.ShardNodes.GetHostPortByPodIndex(idx),
								},
								shards[1].GetName(): {
									ServerAlias:   util.Pointer(shards[1].Status.ShardNodes.GetAliasByPodIndex(2)),
									ServerAddress: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2),
								},
							},
						}), timeout, poll).Should(Not(HaveOccurred()))

					Eventually(assertTwemproxyConfigServerPool(&twemproxyconfig,
						[]twemproxy.Server{
							{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(idx), Priority: 1, Name: "l-shard00"},
							{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(idx), Priority: 1, Name: "l-shard01"},
							{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(idx), Priority: 1, Name: "l-shard02"},
							{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard03"},
							{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard04"},
						}), timeout, poll).Should(Not(HaveOccurred()))

				})

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

func assertTwemproxyConfigStatus(tmc *saasv1alpha1.TwemproxyConfig, sentinel *saasv1alpha1.Sentinel,
	want *saasv1alpha1.TwemproxyConfigStatus) func() error {

	return func() error {
		if err := k8sClient.Get(context.Background(), client.ObjectKeyFromObject(tmc), tmc); err != nil {
			return err
		}

		if err := k8sClient.Get(context.Background(), client.ObjectKeyFromObject(sentinel), sentinel); err != nil {
			return err
		}

		monitoredShards, _ := yaml.Marshal(sentinel.Status.MonitoredShards)
		GinkgoWriter.Printf("[debug] cluster status:\n\n %s\n", monitoredShards)
		selectedTargets, _ := yaml.Marshal(tmc.Status.SelectedTargets)
		GinkgoWriter.Printf("[debug] selected targets:\n\n %s\n", selectedTargets)

		if diff := cmp.Diff(*want, tmc.Status); diff != "" {
			return fmt.Errorf("got unexpected status %s", diff)
		}

		return nil
	}
}

func assertTwemproxyConfigServerPool(tmc *saasv1alpha1.TwemproxyConfig, want []twemproxy.Server) func() error {

	return func() error {
		cm := &corev1.ConfigMap{}
		if err := k8sClient.Get(context.Background(), client.ObjectKeyFromObject(tmc), cm); err != nil {
			return err
		}

		config := map[string]twemproxy.ServerPoolConfig{}
		if err := yaml.Unmarshal([]byte(cm.Data["nutcracker.yml"]), &config); err != nil {
			return err
		}

		if diff := cmp.Diff(want, config["test-pool"].Servers, cmpopts.IgnoreUnexported(twemproxy.Server{})); diff != "" {
			return fmt.Errorf("got unexpected pool servers %s", diff)
		}

		return nil
	}
}
