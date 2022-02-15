package controllers

import (
	"context"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	grafanav1alpha1 "github.com/3scale/saas-operator/pkg/apis/grafana/v1alpha1"
	secretsmanagerv1alpha1 "github.com/3scale/saas-operator/pkg/apis/secrets-manager/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Backend controller", func() {
	var namespace string
	var backend *saasv1alpha1.Backend

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

	When("deploying a defaulted Backend instance", func() {

		BeforeEach(func() {
			By("creating an Backend simple resource")
			backend = &saasv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance",
					Namespace: namespace,
				},
				Spec: saasv1alpha1.BackendSpec{
					Listener: saasv1alpha1.ListenerSpec{
						Endpoint: saasv1alpha1.Endpoint{
							DNS: []string{"backend-listener.example.com"},
						},
					},
					Config: saasv1alpha1.BackendConfig{
						RedisStorageDSN: "storageDSN",
						RedisQueuesDSN:  "queuesDSN",
						SystemEventsHookURL: saasv1alpha1.SecretReference{
							FromVault: &saasv1alpha1.VaultSecretReference{
								Path: "some-path",
								Key:  "some-key",
							},
						},
						SystemEventsHookPassword: saasv1alpha1.SecretReference{
							FromVault: &saasv1alpha1.VaultSecretReference{
								Path: "some-path",
								Key:  "some-key",
							},
						},
						InternalAPIUser: saasv1alpha1.SecretReference{
							FromVault: &saasv1alpha1.VaultSecretReference{
								Path: "some-path",
								Key:  "some-key",
							},
						},
						InternalAPIPassword: saasv1alpha1.SecretReference{
							FromVault: &saasv1alpha1.VaultSecretReference{
								Path: "some-path",
								Key:  "some-key",
							},
						},
					},
				},
			}
			err := k8sClient.Create(context.Background(), backend)
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, backend)
			}, timeout, poll).ShouldNot(HaveOccurred())

			Eventually(func() bool {
				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, backend)
				Expect(err).ToNot(HaveOccurred())
				return len(backend.GetFinalizers()) > 0
			}, timeout, poll).Should(BeTrue())
		})

		It("creates the required backend workload resources", func() {

			dep := &appsv1.Deployment{}
			By("deploying the backend-listener workload",
				checkWorkloadResources(dep,
					expectedWorkload{
						Name:          "backend-listener",
						Namespace:     namespace,
						Replicas:      2,
						ContainerName: "backend-listener",
						ContainterArgs: []string{
							"bin/3scale_backend", "start",
							"-e", "production",
							"-p", "3000",
							"-x", "/dev/stdout",
						},
						PDB:        true,
						HPA:        true,
						PodMonitor: true,
					},
				),
			)

			svc := &corev1.Service{}
			By("deploying the backend-listener service",
				checkResource(svc, expectedResource{
					Name: "backend-listener", Namespace: namespace,
				}),
			)
			Expect(svc.Spec.Selector["deployment"]).To(Equal("backend-listener"))

			By("deploying the backend-listener-internal service",
				checkResource(svc, expectedResource{
					Name: "backend-listener-internal", Namespace: namespace,
				}),
			)
			Expect(svc.Spec.Selector["deployment"]).To(Equal("backend-listener"))

			By("deploying the backend-worker workload",
				checkWorkloadResources(dep,
					expectedWorkload{
						Name:          "backend-worker",
						Namespace:     namespace,
						Replicas:      2,
						ContainerName: "backend-worker",
						ContainterArgs: []string{
							"bin/3scale_backend_worker", "run",
						},
						PDB:        true,
						HPA:        true,
						PodMonitor: true,
					},
				),
			)

			By("deploying the backend-cron workload",
				checkWorkloadResources(dep,
					expectedWorkload{
						Name:          "backend-cron",
						Namespace:     namespace,
						Replicas:      1,
						ContainerName: "backend-cron",
						ContainterArgs: []string{
							"backend-cron",
						},
					},
				),
			)

		})

		It("creates the required backend shared resources", func() {

			gd := &grafanav1alpha1.GrafanaDashboard{}
			By("deploying the backend grafana dashboard",
				checkResource(gd, expectedResource{
					Name: "backend", Namespace: namespace,
				}),
			)

			for _, sdn := range []string{
				"backend-internal-api",
				"backend-system-events-hook",
			} {
				sd := &secretsmanagerv1alpha1.SecretDefinition{}
				By("deploying the backend secret definitions",
					checkResource(sd, expectedResource{Name: sdn, Namespace: namespace}),
				)
			}

		})

		When("updating a backend resource with some customizations", func() {

			rvs := make(map[string]string)

			BeforeEach(func() {
				Eventually(func() error {

					backend := &saasv1alpha1.Backend{}
					if err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						backend,
					); err != nil {
						return err
					}

					rvs["svc/backend-listener"] = getResourceVersion(
						&corev1.Service{}, "backend-listener", namespace,
					)
					rvs["deployment/backend-listener"] = getResourceVersion(
						&appsv1.Deployment{}, "backend-listener", namespace,
					)
					rvs["hpa/backend-worker"] = getResourceVersion(
						&autoscalingv2beta2.HorizontalPodAutoscaler{}, "backend-worker", namespace,
					)
					rvs["deployment/backend-cron"] = getResourceVersion(
						&appsv1.Deployment{}, "backend-cron", namespace,
					)

					patch := client.MergeFrom(backend.DeepCopy())
					backend.Spec.Image = &saasv1alpha1.ImageSpec{
						Name: pointer.StringPtr("newImage"),
						Tag:  pointer.StringPtr("newTag"),
					}
					backend.Spec.Listener.Replicas = pointer.Int32(3)
					backend.Spec.Listener.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{}
					backend.Spec.Listener.Config = &saasv1alpha1.ListenerConfig{
						RedisAsync: pointer.BoolPtr(true),
					}
					backend.Spec.Worker = &saasv1alpha1.WorkerSpec{
						HPA: &saasv1alpha1.HorizontalPodAutoscalerSpec{
							MinReplicas: pointer.Int32(3),
						},
					}
					backend.Spec.Cron = &saasv1alpha1.CronSpec{
						Replicas: pointer.Int32(3),
					}
					backend.Spec.Listener.LoadBalancer = &saasv1alpha1.NLBLoadBalancerSpec{
						CrossZoneLoadBalancingEnabled: pointer.BoolPtr(false),
					}

					return k8sClient.Patch(context.Background(), backend, patch)

				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("updates the backend-listener resources", func() {

				dep := &appsv1.Deployment{}
				By("updating the backend-listener workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:           "backend-listener",
							Namespace:      namespace,
							Replicas:       3,
							ContainerName:  "backend-listener",
							ContainerImage: "newImage:newTag",
							ContainterArgs: []string{
								"bin/3scale_backend", "-s", "falcon", "start",
								"-e", "production",
								"-p", "3000",
								"-x", "/dev/stdout",
							},
							PDB:         true,
							PodMonitor:  true,
							LastVersion: rvs["deployment/backend-listener"],
						},
					),
				)

				svc := &corev1.Service{}
				By("updating backend-listener service",
					checkResource(svc, expectedResource{
						Name: "backend-listener", Namespace: namespace,
						LastVersion: rvs["svc/backend-listener"],
					}),
				)
				Expect(svc.Spec.Selector["deployment"]).To(Equal("backend-listener"))
				Expect(svc.GetAnnotations()["service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled"]).To(Equal("false"))

				hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{}
				By("updating the backend-worker workload",
					checkResource(hpa,
						expectedResource{
							Name:        "backend-worker",
							Namespace:   namespace,
							LastVersion: rvs["hpa/backend-worker"],
						},
					),
				)
				Expect(hpa.Spec.MinReplicas).To(Equal(pointer.Int32(3)))

				By("updating the backend-cron workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:           "backend-cron",
							Namespace:      namespace,
							Replicas:       3,
							ContainerName:  "backend-cron",
							ContainerImage: "newImage:newTag",
							LastVersion:    rvs["deployment/backend-cron"],
						},
					),
				)

			})
		})

		When("updating a backend resource with canary", func() {

			// Resource Versions
			rvs := make(map[string]string)

			BeforeEach(func() {
				Eventually(func() error {

					backend := &saasv1alpha1.Backend{}
					if err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						backend,
					); err != nil {
						return err
					}

					rvs["svc/backend-listener"] = getResourceVersion(
						&corev1.Service{}, "backend-listener", namespace,
					)
					rvs["deployment/backend-listener"] = getResourceVersion(
						&appsv1.Deployment{}, "backend-listener", namespace,
					)
					rvs["deployment/backend-worker"] = getResourceVersion(
						&appsv1.Deployment{}, "backend-worker", namespace,
					)

					patch := client.MergeFrom(backend.DeepCopy())
					backend.Spec.Listener.Canary = &saasv1alpha1.Canary{
						ImageName: pointer.StringPtr("newImage"),
						ImageTag:  pointer.StringPtr("newTag"),
					}
					backend.Spec.Worker = &saasv1alpha1.WorkerSpec{
						Canary: &saasv1alpha1.Canary{
							ImageName: pointer.StringPtr("newImage"),
							ImageTag:  pointer.StringPtr("newTag"),
							Patches: []string{
								`[{"op": "add", "path": "/config/rackEnv", "value": "test"}]`,
								`[{"op": "replace", "path": "/config/redisStorageDSN", "value": "testDSN"}]`,
							},
						},
					}
					return k8sClient.Patch(context.Background(), backend, patch)

				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("creates the required cannary resources", func() {

				dep := &appsv1.Deployment{}
				By("deploying the backend-listener-canary workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:           "backend-listener-canary",
							Namespace:      namespace,
							Replicas:       2,
							ContainerName:  "backend-listener",
							ContainerImage: "newImage:newTag",
							ContainterArgs: []string{
								"bin/3scale_backend", "start",
								"-e", "production",
								"-p", "3000",
								"-x", "/dev/stdout",
							},
							PodMonitor:  true,
							LastVersion: rvs["deployment/backend-listener"],
						},
					),
				)

				svc := &corev1.Service{}
				By("keeps the backend-listener service deployment label selector",
					checkResource(svc, expectedResource{
						Name: "backend-listener", Namespace: namespace,
					}),
				)
				Expect(svc.Spec.Selector["deployment"]).To(Equal("backend-listener"))
				Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("backend-listener"))

				By("keeps the backend-listener-internal service deployment label selector",
					checkResource(svc, expectedResource{
						Name: "backend-listener-internal", Namespace: namespace,
					}),
				)
				Expect(svc.Spec.Selector["deployment"]).To(Equal("backend-listener"))
				Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("backend-listener"))

				By("deploying the backend-worker-canary workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:           "backend-worker-canary",
							Namespace:      namespace,
							Replicas:       2,
							ContainerName:  "backend-worker",
							ContainerImage: "newImage:newTag",
							ContainterArgs: []string{
								"bin/3scale_backend_worker", "run",
							},
							PodMonitor: true,
						},
					),
				)
				Expect(dep.Spec.Template.Spec.Containers[0].Env[0].Name).To(Equal("RACK_ENV"))
				Expect(dep.Spec.Template.Spec.Containers[0].Env[0].Value).To(Equal("test"))
				Expect(dep.Spec.Template.Spec.Containers[0].Env[1].Name).To(Equal("CONFIG_REDIS_PROXY"))
				Expect(dep.Spec.Template.Spec.Containers[0].Env[1].Value).To(Equal("testDSN"))

			})

			Context("and enabling canary traffic", func() {

				BeforeEach(func() {
					Eventually(func() error {

						backend := &saasv1alpha1.Backend{}
						if err := k8sClient.Get(
							context.Background(),
							types.NamespacedName{Name: "instance", Namespace: namespace},
							backend,
						); err != nil {
							return err
						}

						rvs["deployment/backend-listener-canary"] = getResourceVersion(
							&appsv1.Deployment{}, "backend-listener-canary", namespace,
						)
						rvs["svc/backend-listener"] = getResourceVersion(
							&corev1.Service{}, "backend-listener", namespace,
						)
						rvs["deployment/backend-worker-canary"] = getResourceVersion(
							&appsv1.Deployment{}, "backend-worker-canary", namespace,
						)

						patch := client.MergeFrom(backend.DeepCopy())
						backend.Spec.Listener.Replicas = pointer.Int32(3)
						backend.Spec.Listener.Canary = &saasv1alpha1.Canary{
							SendTraffic: *pointer.Bool(true),
							Replicas:    pointer.Int32(3),
						}
						backend.Spec.Worker = &saasv1alpha1.WorkerSpec{
							Replicas: pointer.Int32(3),
						}
						backend.Spec.Worker.Canary = &saasv1alpha1.Canary{
							Replicas: pointer.Int32(3),
						}
						return k8sClient.Patch(context.Background(), backend, patch)

					}, timeout, poll).ShouldNot(HaveOccurred())
				})

				It("updates the backend resources", func() {

					dep := &appsv1.Deployment{}
					By("scaling up the backend-listener-canary workload",
						checkWorkloadResources(dep,
							expectedWorkload{
								Name:        "backend-listener-canary",
								Namespace:   namespace,
								Replicas:    3,
								PodMonitor:  true,
								LastVersion: rvs["deployment/backend-listener-canary"],
							},
						),
					)
					Expect(dep.Spec.Replicas).To(Equal(pointer.Int32(3)))

					svc := &corev1.Service{}
					By("removing the backend-listener service deployment label selector",
						checkResource(svc, expectedResource{
							Name: "backend-listener", Namespace: namespace,
							LastVersion: rvs["svc/backend-listener"],
						}),
					)
					Expect(svc.Spec.Selector).ToNot(HaveKey("deployment"))
					Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("backend-listener"))

					By("scaling up the backend-worker-canary workload",
						checkWorkloadResources(&appsv1.Deployment{},
							expectedWorkload{
								Name:        "backend-worker-canary",
								Namespace:   namespace,
								Replicas:    3,
								PodMonitor:  true,
								LastVersion: rvs["deployment/backend-worker-canary"],
							},
						),
					)

				})

			})

		})

		When("updating a backend resource with twemproxyconfig", func() {

			// Resource Versions
			rvs := make(map[string]string)

			BeforeEach(func() {
				Eventually(func() error {

					backend := &saasv1alpha1.Backend{}
					if err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						backend,
					); err != nil {
						return err
					}

					rvs["deployment/backend-listener"] = getResourceVersion(
						&appsv1.Deployment{}, "backend-listener", namespace,
					)
					rvs["deployment/backend-worker"] = getResourceVersion(
						&appsv1.Deployment{}, "backend-worker", namespace,
					)
					rvs["deployment/backend-cron"] = getResourceVersion(
						&appsv1.Deployment{}, "backend-cron", namespace,
					)

					patch := client.MergeFrom(backend.DeepCopy())
					backend.Spec.Listener.Replicas = pointer.Int32(2)
					backend.Spec.Listener.Canary = &saasv1alpha1.Canary{
						Replicas: pointer.Int32(2),
						Patches: []string{
							`[{"op":"add","path":"/twemproxy","value":{"twemproxyConfigRef":"backend-canary-twemproxyconfig","options":{"logLevel":3}}}]`,
						},
					}
					backend.Spec.Worker = &saasv1alpha1.WorkerSpec{
						Replicas: pointer.Int32(2),
					}
					backend.Spec.Worker.Canary = &saasv1alpha1.Canary{
						Replicas: pointer.Int32(2),
						Patches: []string{
							`[{"op":"add","path":"/twemproxy/options","value":{"logLevel":4}}]`,
						},
					}
					backend.Spec.Twemproxy = &saasv1alpha1.TwemproxySpec{
						TwemproxyConfigRef: "backend-twemproxyconfig",
						Options: &saasv1alpha1.TwemproxyOptions{
							LogLevel: pointer.Int32Ptr(2),
						},
					}

					return k8sClient.Patch(context.Background(), backend, patch)

				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("updates the backend-listener resources", func() {

				dep := &appsv1.Deployment{}
				By("adding a twemproxy sidecar to the backend-listener workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:        "backend-listener",
							Namespace:   namespace,
							Replicas:    2,
							PDB:         true,
							HPA:         true,
							PodMonitor:  true,
							LastVersion: rvs["deployment/backend-listener"],
						},
					),
				)

				Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("twemproxy-config"))
				Expect(dep.Spec.Template.Spec.Volumes[0].VolumeSource.ConfigMap.LocalObjectReference.Name).To(Equal("backend-twemproxyconfig"))
				Expect(dep.Spec.Template.Spec.Containers[1].Name).To(Equal("twemproxy"))
				Expect(dep.Spec.Template.Spec.Containers[1].VolumeMounts[0].Name).To(Equal("twemproxy-config"))

				By("adding a twemproxy sidecar to the backend-listener-canary workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:       "backend-listener-canary",
							Replicas:   2,
							Namespace:  namespace,
							PodMonitor: true,
						},
					),
				)

				Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("twemproxy-config"))
				Expect(dep.Spec.Template.Spec.Volumes[0].VolumeSource.ConfigMap.LocalObjectReference.Name).To(Equal("backend-canary-twemproxyconfig"))
				Expect(dep.Spec.Template.Spec.Containers[1].Name).To(Equal("twemproxy"))
				Expect(dep.Spec.Template.Spec.Containers[1].VolumeMounts[0].Name).To(Equal("twemproxy-config"))
				Expect(dep.Spec.Template.Spec.Containers[1].Env[3].Name).To(Equal("TWEMPROXY_LOG_LEVEL"))
				Expect(dep.Spec.Template.Spec.Containers[1].Env[3].Value).To(Equal("3"))

				By("adding a twemproxy sidecar to the backend-worker workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:        "backend-worker",
							Namespace:   namespace,
							Replicas:    2,
							PDB:         true,
							HPA:         true,
							PodMonitor:  true,
							LastVersion: rvs["deployment/backend-worker"],
						},
					),
				)

				Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("twemproxy-config"))
				Expect(dep.Spec.Template.Spec.Volumes[0].VolumeSource.ConfigMap.LocalObjectReference.Name).To(Equal("backend-twemproxyconfig"))
				Expect(dep.Spec.Template.Spec.Containers[1].Name).To(Equal("twemproxy"))
				Expect(dep.Spec.Template.Spec.Containers[1].VolumeMounts[0].Name).To(Equal("twemproxy-config"))

				By("adding a twemproxy sidecar to the backend-worker-canary workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:       "backend-worker-canary",
							Replicas:   2,
							Namespace:  namespace,
							PodMonitor: true,
						},
					),
				)

				Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("twemproxy-config"))
				Expect(dep.Spec.Template.Spec.Volumes[0].VolumeSource.ConfigMap.LocalObjectReference.Name).To(Equal("backend-twemproxyconfig"))
				Expect(dep.Spec.Template.Spec.Containers[1].Name).To(Equal("twemproxy"))
				Expect(dep.Spec.Template.Spec.Containers[1].VolumeMounts[0].Name).To(Equal("twemproxy-config"))
				Expect(dep.Spec.Template.Spec.Containers[1].Env[3].Name).To(Equal("TWEMPROXY_LOG_LEVEL"))
				Expect(dep.Spec.Template.Spec.Containers[1].Env[3].Value).To(Equal("4"))

				By("not updating the backend-cron workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:      "backend-cron",
							Namespace: namespace,
							Replicas:  1,
						},
					),
				)

				Expect(dep.GetResourceVersion()).To(Equal(rvs["deployment/backend-cron"]))
				Expect(dep.Spec.Template.Spec.Containers).To(HaveLen(1))
				Expect(dep.Spec.Template.Spec.Containers[0].Name).To(Equal("backend-cron"))

			})
		})

	})
})
