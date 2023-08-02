package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/resource_builders/twemproxy"
	testutil "github.com/3scale/saas-operator/test/util"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
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
				Spec:       saasv1alpha1.RedisShardSpec{MasterIndex: pointer.Int32(0), SlaveCount: pointer.Int32(2)},
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

		sentinel = saasv1alpha1.Sentinel{
			ObjectMeta: metav1.ObjectMeta{Name: "sentinel", Namespace: ns},
			Spec: saasv1alpha1.SentinelSpec{
				Config: &saasv1alpha1.SentinelConfig{
					MonitoredShards: map[string][]string{
						shards[0].GetName(): {
							"redis://" + shards[0].Status.ShardNodes.GetHostPortByPodIndex(0),
							"redis://" + shards[0].Status.ShardNodes.GetHostPortByPodIndex(1),
							"redis://" + shards[0].Status.ShardNodes.GetHostPortByPodIndex(2),
						},
						shards[1].GetName(): {
							"redis://" + shards[1].Status.ShardNodes.GetHostPortByPodIndex(0),
							"redis://" + shards[1].Status.ShardNodes.GetHostPortByPodIndex(1),
							"redis://" + shards[1].Status.ShardNodes.GetHostPortByPodIndex(2),
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

	})

	When("TwemproxyConfig resource is created targeting redis masters", func() {

		BeforeEach(func() {

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

			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{Name: twemproxyconfig.GetName(), Namespace: ns}, &twemproxyconfig)
			}, timeout, poll).ShouldNot(HaveOccurred())
		})

		It("deploys a ConfigMap with twemproxy configuration that points to redis masters", func() {

			cm := &corev1.ConfigMap{}
			By("getting the generated ConfigMap",
				(&testutil.ExpectedResource{Name: "tmc-instance", Namespace: ns}).
					Assert(k8sClient, cm, timeout, poll))

			config := map[string]twemproxy.ServerPoolConfig{}
			data, err := yaml.YAMLToJSON([]byte(cm.Data["nutcracker.yml"]))
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal(data, &config)
			Expect(err).ToNot(HaveOccurred())

			Expect(config["test-pool"].Servers).To(Equal(
				[]twemproxy.Server{
					{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard00"},
					{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard01"},
					{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard02"},
					{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard03"},
					{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard04"},
				},
			))
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

			It("updates the twemproxy configuration with new master", func() {
				Eventually(func() error {

					cm := &corev1.ConfigMap{}
					err := k8sClient.Get(context.Background(), types.NamespacedName{Name: twemproxyconfig.GetName(), Namespace: ns}, cm)
					if err != nil {
						return err
					}

					config := map[string]twemproxy.ServerPoolConfig{}
					data, err := yaml.YAMLToJSON([]byte(cm.Data["nutcracker.yml"]))
					if err != nil {
						return err
					}
					err = json.Unmarshal(data, &config)
					if err != nil {
						return err
					}

					if diff := cmp.Diff(config["test-pool"].Servers,
						[]twemproxy.Server{
							{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(1), Priority: 1, Name: "l-shard00"},
							{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(1), Priority: 1, Name: "l-shard01"},
							{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(1), Priority: 1, Name: "l-shard02"},
							{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard03"},
							{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard04"},
						},
					); diff != "" {
						return fmt.Errorf("got unexpected pool servers %s", diff)
					}

					return nil
				}, timeout, poll).ShouldNot(HaveOccurred())
			})
		})

		AfterEach(func() {
			// Delete twemproxyconfig
			err := k8sClient.Delete(context.Background(), &twemproxyconfig, client.PropagationPolicy(metav1.DeletePropagationForeground))
			Expect(err).ToNot(HaveOccurred())
		})

	})

	When("TwemproxyConfig resource is created targeting redis rw-slaves", func() {
		cm := &corev1.ConfigMap{}

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

			By("creating the TwemproxyConfig resource pointing to rw-slaves", func() {
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

				Eventually(func() error {
					return k8sClient.Get(context.Background(), types.NamespacedName{Name: twemproxyconfig.GetName(), Namespace: ns}, &twemproxyconfig)
				}, timeout, poll).ShouldNot(HaveOccurred())

				By("getting the generated ConfigMap",
					(&testutil.ExpectedResource{Name: "tmc-instance", Namespace: ns}).
						Assert(k8sClient, cm, timeout, poll))
			})

		})

		It("deploys a ConfigMap with twemproxy configuration that points to redis rw-slaves", func() {

			config := map[string]twemproxy.ServerPoolConfig{}
			data, err := yaml.YAMLToJSON([]byte(cm.Data["nutcracker.yml"]))
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal(data, &config)
			Expect(err).ToNot(HaveOccurred())

			Expect(config["test-pool"].Servers).To(Equal(
				[]twemproxy.Server{
					{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard00"},
					{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard01"},
					{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard02"},
					{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard03"},
					{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard04"},
				},
			))
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

					Eventually(func() []twemproxy.Server {

						By("getting the ConfigMap",
							(&testutil.ExpectedResource{Name: "tmc-instance", Namespace: ns}).
								Assert(k8sClient, cm, timeout, poll))

						config := map[string]twemproxy.ServerPoolConfig{}
						data, err := yaml.YAMLToJSON([]byte(cm.Data["nutcracker.yml"]))
						Expect(err).ToNot(HaveOccurred())
						err = json.Unmarshal(data, &config)
						Expect(err).ToNot(HaveOccurred())
						return config["test-pool"].Servers

					}, timeout, poll).Should(Equal(
						[]twemproxy.Server{
							{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard00"},
							{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard01"},
							{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(0), Priority: 1, Name: "l-shard02"},
							{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard03"},
							{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard04"},
						},
					))

				})

				By("checking the config for rs0 points back to rw-slave once it's recovered", func() {

					Eventually(func() []twemproxy.Server {

						By("getting the ConfigMap",
							(&testutil.ExpectedResource{Name: "tmc-instance", Namespace: ns}).
								Assert(k8sClient, cm, timeout, poll))

						config := map[string]twemproxy.ServerPoolConfig{}
						data, err := yaml.YAMLToJSON([]byte(cm.Data["nutcracker.yml"]))
						Expect(err).ToNot(HaveOccurred())
						err = json.Unmarshal(data, &config)
						Expect(err).ToNot(HaveOccurred())
						return config["test-pool"].Servers

					}, timeout, poll).Should(Equal(
						[]twemproxy.Server{
							{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard00"},
							{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard01"},
							{Address: shards[0].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard02"},
							{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard03"},
							{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard04"},
						},
					))

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

			It("reconfigures shard rs0 to point to the first slave in alphabetical order (by address)", func() {

				By("checking the twemproxy config for shard rs0", func() {

					// determine which should be the chosen rw-slave
					addresses := []string{
						shards[0].Status.ShardNodes.GetHostPortByPodIndex(1),
						shards[0].Status.ShardNodes.GetHostPortByPodIndex(2),
					}
					sort.Strings(addresses)
					expectedRWSlave := addresses[0]

					Eventually(func() []twemproxy.Server {

						By("getting the ConfigMap",
							(&testutil.ExpectedResource{Name: "tmc-instance", Namespace: ns}).
								Assert(k8sClient, cm, timeout, poll))

						config := map[string]twemproxy.ServerPoolConfig{}
						data, err := yaml.YAMLToJSON([]byte(cm.Data["nutcracker.yml"]))
						Expect(err).ToNot(HaveOccurred())
						err = json.Unmarshal(data, &config)
						Expect(err).ToNot(HaveOccurred())
						return config["test-pool"].Servers

					}, timeout, poll).Should(Equal(
						[]twemproxy.Server{
							{Address: expectedRWSlave, Priority: 1, Name: "l-shard00"},
							{Address: expectedRWSlave, Priority: 1, Name: "l-shard01"},
							{Address: expectedRWSlave, Priority: 1, Name: "l-shard02"},
							{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard03"},
							{Address: shards[1].Status.ShardNodes.GetHostPortByPodIndex(2), Priority: 1, Name: "l-shard04"},
						},
					))

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
