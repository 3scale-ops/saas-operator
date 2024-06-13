package controllers

import (
	"context"
	"fmt"

	"github.com/3scale-ops/basereconciler/util"
	marin3rv1alpha1 "github.com/3scale-ops/marin3r/apis/marin3r/v1alpha1"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	testutil "github.com/3scale-ops/saas-operator/test/util"
	grafanav1beta1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Apicast controller", func() {
	var namespace string
	var apicast *saasv1alpha1.Apicast

	BeforeEach(func() {
		// Create a namespace for each block
		namespace = "test-ns-" + nameGenerator.Generate()

		// Add any setup steps that needs to be executed before each test
		testNamespace := &corev1.Namespace{
			TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Namespace"},
			ObjectMeta: metav1.ObjectMeta{Name: namespace},
		}

		err := k8sClient.Create(context.Background(), testNamespace)
		Expect(err).ToNot(HaveOccurred())

		n := &corev1.Namespace{}
		Eventually(func() error {
			return k8sClient.Get(context.Background(), types.NamespacedName{Name: namespace}, n)
		}, timeout, poll).ShouldNot(HaveOccurred())

	})

	When("deploying a defaulted Apicast instance", func() {

		BeforeEach(func() {
			By("creating an Apicast simple resource", func() {
				apicast = &saasv1alpha1.Apicast{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "instance",
						Namespace: namespace,
					},
					Spec: saasv1alpha1.ApicastSpec{
						Staging: saasv1alpha1.ApicastEnvironmentSpec{
							Config: saasv1alpha1.ApicastConfig{
								ConfigurationCache:       30,
								ThreescalePortalEndpoint: "http://example/config",
							},
							Endpoint: &saasv1alpha1.Endpoint{
								DNS: []string{"apicast-staging.example.com"},
							},
						},
						Production: saasv1alpha1.ApicastEnvironmentSpec{
							Config: saasv1alpha1.ApicastConfig{
								ConfigurationCache:       300,
								ThreescalePortalEndpoint: "http://example/config",
							},
							Endpoint: &saasv1alpha1.Endpoint{
								DNS: []string{"apicast-production.example.com"},
							},
						},
					},
				}
				err := k8sClient.Create(context.Background(), apicast)
				Expect(err).ToNot(HaveOccurred())
				Eventually(func() error {
					return k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, apicast)
				}, timeout, poll).ShouldNot(HaveOccurred())
			})
		})

		It("creates the required Apicast resources", func() {

			dep := &appsv1.Deployment{}
			By("deploying an apicast-production workload",
				(&testutil.ExpectedWorkload{
					Name:          "apicast-production",
					Namespace:     namespace,
					Replicas:      2,
					ContainerName: "apicast",
					PDB:           true,
					HPA:           true,
					PodMonitor:    true,
				}).Assert(k8sClient, dep, timeout, poll),
			)
			for _, env := range dep.Spec.Template.Spec.Containers[0].Env {
				switch env.Name {
				case "THREESCALE_DEPLOYMENT_ENV":
					Expect(env.Value).To(Equal("production"))
				case "APICAST_CONFIGURATION_LOADER":
					Expect(env.Value).To(Equal("lazy"))
				case "APICAST_LOG_LEVEL":
					Expect(env.Value).To(Equal("warn"))
				case "APICAST_RESPONSE_CODES":
					Expect(env.Value).To(Equal("true"))
				}
			}

			By("deploying an apicast-staging workload",
				(&testutil.ExpectedWorkload{
					Name:          "apicast-staging",
					Namespace:     namespace,
					Replicas:      2,
					ContainerName: "apicast",
					PDB:           true,
					HPA:           true,
					PodMonitor:    true,
				}).Assert(k8sClient, dep, timeout, poll),
			)
			for _, env := range dep.Spec.Template.Spec.Containers[0].Env {
				switch env.Name {
				case "APICAST_CONFIGURATION_LOADER":
					Expect(env.Value).To(Equal("lazy"))
				case "APICAST_LOG_LEVEL":
					Expect(env.Value).To(Equal("warn"))
				case "THREESCALE_DEPLOYMENT_ENV":
					Expect(env.Value).To(Equal("staging"))
				case "APICAST_RESPONSE_CODES":
					Expect(env.Value).To(Equal("true"))
				}
			}

			svc := &corev1.Service{}
			By("deploying the apicast-production service",
				(&testutil.ExpectedResource{Name: "apicast-production", Namespace: namespace}).
					Assert(k8sClient, svc, timeout, poll))

			Expect(svc.Spec.Selector["deployment"]).To(Equal("apicast-production"))
			Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("apicast-production"))
			Expect(
				svc.Annotations["external-dns.alpha.kubernetes.io/hostname"],
			).To(Equal("apicast-production.example.com"))
			Expect(
				svc.Annotations["service.beta.kubernetes.io/aws-load-balancer-proxy-protocol"],
			).To(Equal("*"))
			Expect(
				svc.Annotations["service.beta.kubernetes.io/aws-load-balancer-connection-draining-enabled"],
			).To(Equal("true"))
			Expect(
				svc.Annotations["service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled"],
			).To(Equal("true"))

			By("deploying the apicast-production-management service",
				(&testutil.ExpectedResource{Name: "apicast-production-management", Namespace: namespace}).
					Assert(k8sClient, svc, timeout, poll),
			)
			Expect(svc.Spec.Selector["deployment"]).To(Equal("apicast-production"))
			Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("apicast-production"))
			Expect(svc.Spec.Ports[0].Name).To(Equal("management"))
			Expect(svc.Annotations).To(HaveLen(0))

			By("deploying the apicast-staging service",
				(&testutil.ExpectedResource{
					Name: "apicast-staging", Namespace: namespace}).
					Assert(k8sClient, svc, timeout, poll))

			Expect(svc.Spec.Selector["deployment"]).To(Equal("apicast-staging"))
			Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("apicast-staging"))
			Expect(
				svc.Annotations["external-dns.alpha.kubernetes.io/hostname"],
			).To(Equal("apicast-staging.example.com"))
			Expect(
				svc.Annotations["service.beta.kubernetes.io/aws-load-balancer-proxy-protocol"],
			).To(Equal("*"))
			Expect(
				svc.Annotations["service.beta.kubernetes.io/aws-load-balancer-connection-draining-enabled"],
			).To(Equal("true"))
			Expect(
				svc.Annotations["service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled"],
			).To(Equal("true"))

			By("deploying the apicast-staging-management service",
				(&testutil.ExpectedResource{Name: "apicast-staging-management", Namespace: namespace}).
					Assert(k8sClient, svc, timeout, poll))

			Expect(svc.Spec.Selector["deployment"]).To(Equal("apicast-staging"))
			Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("apicast-staging"))
			Expect(svc.Spec.Ports[0].Name).To(Equal("management"))
			Expect(svc.Annotations).To(HaveLen(0))

			By("deploying an apicast grafana dashboard",
				(&testutil.ExpectedResource{Name: "apicast", Namespace: namespace}).
					Assert(k8sClient, &grafanav1beta1.GrafanaDashboard{}, timeout, poll))

			By("deploying an apicast-services grafana dashboard",
				(&testutil.ExpectedResource{Name: "apicast-services", Namespace: namespace}).
					Assert(k8sClient, &grafanav1beta1.GrafanaDashboard{}, timeout, poll))

		})

		It("doesn't create the non-default resources", func() {

			dep := &appsv1.Deployment{}
			By("ensuring an apicast-production-canary workload is not created",
				(&testutil.ExpectedResource{Name: "apicast-production-canary", Namespace: namespace, Missing: true}).
					Assert(k8sClient, dep, timeout, poll))
			By("ensuring an apicast-staging-canary workload is not created",
				(&testutil.ExpectedResource{Name: "apicast-staging-canary", Namespace: namespace, Missing: true}).
					Assert(k8sClient, dep, timeout, poll))

			ec := &marin3rv1alpha1.EnvoyConfig{}
			By("ensuring an apicast-production envoyconfig is not created",
				(&testutil.ExpectedResource{Name: "apicast-staging-production", Namespace: namespace, Missing: true}).
					Assert(k8sClient, ec, timeout, poll))
			By("ensuring an apicast-staging envoyconfig is not created",
				(&testutil.ExpectedResource{Name: "apicast-staging-canary", Namespace: namespace, Missing: true}).
					Assert(k8sClient, ec, timeout, poll))

		})

		When("updating an Apicast resource with customizations", func() {

			// Resource Versions
			rvs := make(map[string]string)

			BeforeEach(func() {
				Eventually(func() error {

					apicast := &saasv1alpha1.Apicast{}
					if err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						apicast,
					); err != nil {
						return err
					}

					rvs["deployment/apicast-production"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "apicast-production", namespace, timeout, poll)

					rvs["deployment/apicast-staging"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "apicast-staging", namespace, timeout, poll)

					patch := client.MergeFrom(apicast.DeepCopy())
					apicast.Spec.Production = saasv1alpha1.ApicastEnvironmentSpec{
						Config: saasv1alpha1.ApicastConfig{
							ConfigurationCache:       42,
							ThreescalePortalEndpoint: "http://updated-example/config",
						},
						Endpoint: &saasv1alpha1.Endpoint{
							DNS: []string{"updated-apicast-production.example.com"},
						},
						HPA: &saasv1alpha1.HorizontalPodAutoscalerSpec{
							MinReplicas: util.Pointer[int32](3),
						},
						LivenessProbe:  &saasv1alpha1.ProbeSpec{},
						ReadinessProbe: &saasv1alpha1.ProbeSpec{},
						Marin3r: &saasv1alpha1.Marin3rSidecarSpec{
							NodeID: util.Pointer("apicast-production"),
							EnvoyDynamicConfig: saasv1alpha1.MapOfEnvoyDynamicConfig{
								"http": {
									GeneratorVersion: util.Pointer("v1"),
									ListenerHttp: &saasv1alpha1.ListenerHttp{
										Port:            8080,
										RouteConfigName: "route",
									},
								}},
						},
					}
					apicast.Spec.Staging = saasv1alpha1.ApicastEnvironmentSpec{
						Config: saasv1alpha1.ApicastConfig{
							ConfigurationCache:       42,
							ThreescalePortalEndpoint: "http://updated-example/config",
						},
						Endpoint: &saasv1alpha1.Endpoint{
							DNS: []string{"updated-apicast-staging.example.com"},
						},
						HPA: &saasv1alpha1.HorizontalPodAutoscalerSpec{
							MinReplicas: util.Pointer[int32](3),
						},
						LivenessProbe:  &saasv1alpha1.ProbeSpec{},
						ReadinessProbe: &saasv1alpha1.ProbeSpec{},
						Marin3r: &saasv1alpha1.Marin3rSidecarSpec{
							NodeID: util.Pointer("apicast-production"),
							EnvoyDynamicConfig: saasv1alpha1.MapOfEnvoyDynamicConfig{
								"http": {

									GeneratorVersion: util.Pointer("v1"),
									ListenerHttp: &saasv1alpha1.ListenerHttp{
										Port:            8080,
										RouteConfigName: "route",
									},
								},
							},
						},
					}
					apicast.Spec.GrafanaDashboard = &saasv1alpha1.GrafanaDashboardSpec{}

					return k8sClient.Patch(context.Background(), apicast, patch)

				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("updates the Apicast resources", func() {

				By("ensuring the Apicast grafana dashboard is gone",
					(&testutil.ExpectedResource{
						Name:      "apicast-production",
						Namespace: namespace,
						Missing:   true,
					}).Assert(k8sClient, &grafanav1beta1.GrafanaDashboard{}, timeout, poll),
				)

				dep := &appsv1.Deployment{}

				By("updating the Apicast Production workload",
					(&testutil.ExpectedWorkload{
						Name:          "apicast-production",
						Namespace:     namespace,
						Replicas:      3,
						ContainerName: "apicast",
						HPA:           true,
						PDB:           true,
						PodMonitor:    true,
						EnvoyConfig:   true,
						LastVersion:   rvs["deployment/apicast-production"],
					}).Assert(k8sClient, dep, timeout, poll))

				for _, env := range dep.Spec.Template.Spec.Containers[0].Env {
					switch env.Name {
					case "THREESCALE_PORTAL_ENDPOINT":
						Expect(env.Value).To(Equal("http://updated-example/config"))
					case "THREESCALE_DEPLOYMENT_ENV":
						Expect(env.Value).To(Equal("production"))
					case "APICAST_RESPONSE_CODES":
						Expect(env.Value).To(Equal("true"))
					case "APICAST_CONFIGURATION_CACHE":
						Expect(env.Value).To(Equal("42"))
					}
				}

				By("updating the Apicast Staging workload",
					(&testutil.ExpectedWorkload{
						Name:          "apicast-staging",
						Namespace:     namespace,
						Replicas:      3,
						ContainerName: "apicast",
						HPA:           true,
						PDB:           true,
						PodMonitor:    true,
						EnvoyConfig:   true,
						LastVersion:   rvs["deployment/apicast-staging"],
					}).Assert(k8sClient, dep, timeout, poll))

				for _, env := range dep.Spec.Template.Spec.Containers[0].Env {
					switch env.Name {
					case "THREESCALE_PORTAL_ENDPOINT":
						Expect(env.Value).To(Equal("http://updated-example/config"))
					case "APICAST_LOG_LEVEL":
						Expect(env.Value).To(Equal("warn"))
					case "THREESCALE_DEPLOYMENT_ENV":
						Expect(env.Value).To(Equal("staging"))
					case "APICAST_RESPONSE_CODES":
						Expect(env.Value).To(Equal("true"))
					case "APICAST_CONFIGURATION_CACHE":
						Expect(env.Value).To(Equal("42"))
					}
				}

				svc := &corev1.Service{}
				By("updating the apicast-production service annotation",
					(&testutil.ExpectedResource{
						Name: "apicast-production", Namespace: namespace,
					}).Assert(k8sClient, svc, timeout, poll),
				)
				Expect(svc.Annotations["external-dns.alpha.kubernetes.io/hostname"]).To(
					Equal("updated-apicast-production.example.com"),
				)

				By("updating the apicast-staging service annotation",
					(&testutil.ExpectedResource{
						Name: "apicast-staging", Namespace: namespace,
					}).Assert(k8sClient, svc, timeout, poll),
				)
				Expect(svc.Annotations["external-dns.alpha.kubernetes.io/hostname"]).To(
					Equal("updated-apicast-staging.example.com"),
				)

			})

		})

		When("updating an Apicast resource with canary", func() {

			// Resource Versions
			rvs := make(map[string]string)

			BeforeEach(func() {
				Eventually(func() error {
					apicast := &saasv1alpha1.Apicast{}
					if err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						apicast,
					); err != nil {
						return err
					}

					rvs["svc/apicast-production"] = testutil.GetResourceVersion(
						k8sClient, &corev1.Service{}, "apicast-production", namespace, timeout, poll)
					rvs["deployment/apicast-production"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "apicast-production", namespace, timeout, poll)
					rvs["svc/apicast-staging"] = testutil.GetResourceVersion(
						k8sClient, &corev1.Service{}, "apicast-staging", namespace, timeout, poll)
					rvs["deployment/apicast-staging"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "apicast-staging", namespace, timeout, poll)

					patch := client.MergeFrom(apicast.DeepCopy())
					apicast.Spec.Production.Canary = &saasv1alpha1.Canary{
						ImageName: util.Pointer("newImage"),
						ImageTag:  util.Pointer("newTag"),
						Replicas:  util.Pointer[int32](1),
					}
					apicast.Spec.Staging.Canary = &saasv1alpha1.Canary{
						ImageName: util.Pointer("newImage"),
						ImageTag:  util.Pointer("newTag"),
						Replicas:  util.Pointer[int32](1),
					}

					if err := k8sClient.Patch(context.Background(), apicast, patch); err != nil {
						return err
					}

					if testutil.GetResourceVersion(k8sClient, &appsv1.Deployment{}, "apicast-production-canary", namespace, timeout, poll) == "" {
						return fmt.Errorf("not ready")
					}
					if testutil.GetResourceVersion(k8sClient, &appsv1.Deployment{}, "apicast-staging-canary", namespace, timeout, poll) == "" {
						return fmt.Errorf("not ready")
					}

					return nil

				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("creates the required canary resources", func() {

				dep := &appsv1.Deployment{}
				By("deploying an Apicast-production-canary workload",
					(&testutil.ExpectedWorkload{
						Name:          "apicast-production-canary",
						Namespace:     namespace,
						Replicas:      1,
						ContainerName: "apicast",
						PodMonitor:    true,
					}).Assert(k8sClient, dep, timeout, poll))

				for _, env := range dep.Spec.Template.Spec.Containers[0].Env {
					switch env.Name {
					case "THREESCALE_DEPLOYMENT_ENV":
						Expect(env.Value).To(Equal("production"))
					case "APICAST_LOG_LEVEL":
						Expect(env.Value).To(Equal("warn"))
					case "APICAST_RESPONSE_CODES":
						Expect(env.Value).To(Equal("true"))
					}
				}

				By("deploying an Apicast-staging-canary workload",
					(&testutil.ExpectedWorkload{
						Name:          "apicast-staging-canary",
						Namespace:     namespace,
						Replicas:      1,
						ContainerName: "apicast",
						PodMonitor:    true,
					}).Assert(k8sClient, dep, timeout, poll))

				for _, env := range dep.Spec.Template.Spec.Containers[0].Env {
					switch env.Name {
					case "THREESCALE_DEPLOYMENT_ENV":
						Expect(env.Value).To(Equal("staging"))
					case "APICAST_LOG_LEVEL":
						Expect(env.Value).To(Equal("warn"))
					case "APICAST_RESPONSE_CODES":
						Expect(env.Value).To(Equal("true"))
					}
				}

				svc := &corev1.Service{}
				By("keeping the apicast-production service spec",
					(&testutil.ExpectedResource{
						Name: "apicast-production", Namespace: namespace}).Assert(k8sClient, svc, timeout, poll))

				Expect(svc.Spec.Selector["deployment"]).To(Equal("apicast-production"))
				Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("apicast-production"))
				Expect(
					svc.Annotations["external-dns.alpha.kubernetes.io/hostname"],
				).To(Equal("apicast-production.example.com"))
				Expect(
					svc.Annotations["service.beta.kubernetes.io/aws-load-balancer-proxy-protocol"],
				).To(Equal("*"))
				Expect(
					svc.Annotations["service.beta.kubernetes.io/aws-load-balancer-connection-draining-enabled"],
				).To(Equal("true"))
				Expect(
					svc.Annotations["service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled"],
				).To(Equal("true"))

				By("keeping the apicast-production-management service spec",
					(&testutil.ExpectedResource{
						Name: "apicast-production-management", Namespace: namespace}).Assert(k8sClient, svc, timeout, poll))

				Expect(svc.Spec.Selector["deployment"]).To(Equal("apicast-production"))
				Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("apicast-production"))
				Expect(svc.Spec.Ports[0].Name).To(Equal("management"))
				Expect(svc.Annotations).To(HaveLen(0))

				By("keeping the apicast-staging service spec",
					(&testutil.ExpectedResource{
						Name: "apicast-staging", Namespace: namespace}).Assert(k8sClient, svc, timeout, poll))

				Expect(svc.Spec.Selector["deployment"]).To(Equal("apicast-staging"))
				Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("apicast-staging"))
				Expect(svc.Annotations["external-dns.alpha.kubernetes.io/hostname"]).To(Equal("apicast-staging.example.com"))
				Expect(svc.Annotations["service.beta.kubernetes.io/aws-load-balancer-proxy-protocol"]).To(Equal("*"))
				Expect(svc.Annotations["service.beta.kubernetes.io/aws-load-balancer-connection-draining-enabled"]).To(Equal("true"))
				Expect(svc.Annotations["service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled"]).To(Equal("true"))

				By("keeping the apicast-staging-management service spec",
					(&testutil.ExpectedResource{
						Name: "apicast-staging-management", Namespace: namespace}).Assert(k8sClient, svc, timeout, poll))

				Expect(svc.Spec.Selector["deployment"]).To(Equal("apicast-staging"))
				Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("apicast-staging"))
				Expect(svc.Spec.Ports[0].Name).To(Equal("management"))
				Expect(svc.Annotations).To(HaveLen(0))

			})

			When("enabling canary traffic", func() {

				BeforeEach(func() {
					Eventually(func() error {
						apicast := &saasv1alpha1.Apicast{}
						if err := k8sClient.Get(
							context.Background(),
							types.NamespacedName{Name: "instance", Namespace: namespace},
							apicast,
						); err != nil {
							return err
						}
						rvs["svc/apicast-production"] = testutil.GetResourceVersion(
							k8sClient, &corev1.Service{}, "apicast-production", namespace, timeout, poll)
						rvs["svc/apicast-staging"] = testutil.GetResourceVersion(
							k8sClient, &corev1.Service{}, "apicast-staging", namespace, timeout, poll)

						patch := client.MergeFrom(apicast.DeepCopy())
						apicast.Spec.Production.Canary = &saasv1alpha1.Canary{
							SendTraffic: *util.Pointer(true),
						}
						apicast.Spec.Staging.Canary = &saasv1alpha1.Canary{
							SendTraffic: *util.Pointer(true),
						}
						return k8sClient.Patch(context.Background(), apicast, patch)
					}, timeout, poll).ShouldNot(HaveOccurred())
				})

				It("updates the apicast service", func() {

					svc := &corev1.Service{}
					By("removing the apicast-production service deployment label selector",
						(&testutil.ExpectedResource{
							Name: "apicast-production", Namespace: namespace,
							LastVersion: rvs["svc/apicast-production"],
						}).Assert(k8sClient, svc, timeout, poll))

					Expect(svc.Spec.Selector).NotTo(HaveKey("deployment"))
					Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("apicast-production"))

					By("removing the apicast-staging service deployment label selector",
						(&testutil.ExpectedResource{
							Name: "apicast-staging", Namespace: namespace,
							LastVersion: rvs["svc/apicast-staging"],
						}).Assert(k8sClient, svc, timeout, poll))

					Expect(svc.Spec.Selector).NotTo(HaveKey("deployment"))
					Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("apicast-staging"))

				})

			})

			When("disabling the canary", func() {

				BeforeEach(func() {

					Eventually(func() error {
						apicast := &saasv1alpha1.Apicast{}
						if err := k8sClient.Get(
							context.Background(),
							types.NamespacedName{Name: "instance", Namespace: namespace},
							apicast,
						); err != nil {
							return err
						}
						patch := client.MergeFrom(apicast.DeepCopy())
						apicast.Spec.Production.Canary = nil
						apicast.Spec.Staging.Canary = nil
						return k8sClient.Patch(context.Background(), apicast, patch)
					}, timeout, poll).ShouldNot(HaveOccurred())
				})

				It("deletes the canary resources", func() {

					dep := &appsv1.Deployment{}
					By("removing the apicast-production-canary Deployment",
						(&testutil.ExpectedResource{
							Name: "apicast-production-canary", Namespace: namespace, Missing: true}).Assert(k8sClient, dep, timeout, poll))
					By("removing the apicast-staging-canary Deployment",
						(&testutil.ExpectedResource{
							Name: "apicast-staging-canary", Namespace: namespace, Missing: true}).Assert(k8sClient, dep, timeout, poll))

					pm := &monitoringv1.PodMonitor{}
					By("removing the apicast-production-canary PodMonitor",
						(&testutil.ExpectedResource{
							Name: "apicast-production-canary", Namespace: namespace, Missing: true}).Assert(k8sClient, pm, timeout, poll))
					By("removing the apicast-staging-canary PodMonitor",
						(&testutil.ExpectedResource{
							Name: "apicast-staging-canary", Namespace: namespace, Missing: true}).Assert(k8sClient, pm, timeout, poll))
				})
			})
		})

		When("removing the PDB and HPA from an Apicast instance", func() {

			// Resource Versions
			rvs := make(map[string]string)

			BeforeEach(func() {
				Eventually(func() error {

					apicast := &saasv1alpha1.Apicast{}
					if err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						apicast,
					); err != nil {
						return err
					}

					rvs["deployment/apicast-production"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "apicast-production", namespace, timeout, poll)
					rvs["deployment/apicast-staging"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "apicast-staging", namespace, timeout, poll)

					patch := client.MergeFrom(apicast.DeepCopy())

					apicast.Spec.Production.Replicas = util.Pointer[int32](0)
					apicast.Spec.Production.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{}
					apicast.Spec.Production.PDB = &saasv1alpha1.PodDisruptionBudgetSpec{}

					apicast.Spec.Staging.Replicas = util.Pointer[int32](0)
					apicast.Spec.Staging.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{}
					apicast.Spec.Staging.PDB = &saasv1alpha1.PodDisruptionBudgetSpec{}

					return k8sClient.Patch(context.Background(), apicast, patch)

				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("removes the Apicast disabled resources", func() {

				dep := &appsv1.Deployment{}
				By("updating the Apicast Production workload",
					(&testutil.ExpectedWorkload{
						Name:        "apicast-production",
						Namespace:   namespace,
						Replicas:    0,
						HPA:         false,
						PDB:         false,
						PodMonitor:  true,
						LastVersion: rvs["deployment/apicast-production"],
					}).Assert(k8sClient, dep, timeout, poll),
				)
				By("updating the Apicast Staging workload",
					(&testutil.ExpectedWorkload{
						Name:        "apicast-staging",
						Namespace:   namespace,
						Replicas:    0,
						HPA:         false,
						PDB:         false,
						PodMonitor:  true,
						LastVersion: rvs["deployment/apicast-staging"],
					}).Assert(k8sClient, dep, timeout, poll))

			})

		})

	})

})
