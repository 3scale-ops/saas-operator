package controllers

import (
	"context"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	externalsecretsv1alpha1 "github.com/3scale/saas-operator/pkg/apis/externalsecrets/v1alpha1"
	grafanav1alpha1 "github.com/3scale/saas-operator/pkg/apis/grafana/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("System controller", func() {
	var namespace string
	var system *saasv1alpha1.System

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

	When("deploying a defaulted system instance", func() {

		BeforeEach(func() {
			system = &saasv1alpha1.System{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance",
					Namespace: namespace,
				},
				Spec: saasv1alpha1.SystemSpec{
					Config: saasv1alpha1.SystemConfig{
						DatabaseDSN: saasv1alpha1.SecretReference{
							FromVault: &saasv1alpha1.VaultSecretReference{
								Path: "some-path",
								Key:  "some-key",
							},
						},
						EventsSharedSecret: saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
						Recaptcha: saasv1alpha1.SystemRecaptchaSpec{
							PublicKey:  saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
							PrivateKey: saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
						},
						SecretKeyBase: saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
						AccessCode:    saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
						Segment: saasv1alpha1.SegmentSpec{
							DeletionWorkspace: "value",
							DeletionToken:     saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
							WriteKey:          saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
						},
						Github: saasv1alpha1.GithubSpec{
							ClientID:     saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
							ClientSecret: saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
						},
						RedHatCustomerPortal: saasv1alpha1.RedHatCustomerPortalSpec{
							ClientID:     saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
							ClientSecret: saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
						},
						Bugsnag: &saasv1alpha1.BugsnagSpec{
							APIKey: saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
						},
						DatabaseSecret:   saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
						MemcachedServers: "value",
						Redis: saasv1alpha1.RedisSpec{
							QueuesDSN: "value",
						},
						SMTP: saasv1alpha1.SMTPSpec{
							Address:           "value",
							User:              saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
							Password:          saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
							Port:              1000,
							AuthProtocol:      "value",
							OpenSSLVerifyMode: "value",
							STARTTLSAuto:      false,
						},
						MappingServiceAccessToken: saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
						ZyncAuthToken:             saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
						Backend: saasv1alpha1.SystemBackendSpec{
							ExternalEndpoint:    "value",
							InternalEndpoint:    "value",
							InternalAPIUser:     saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
							InternalAPIPassword: saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
							RedisDSN:            "value",
						},
						Assets: saasv1alpha1.AssetsSpec{
							Host:      pointer.StringPtr("test.cloudfront.net"),
							Bucket:    "bucket",
							Region:    "us-east-1",
							AccessKey: saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
							SecretKey: saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
						},
					},
				},
			}
			err := k8sClient.Create(context.Background(), system)
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, system)
			}, timeout, poll).ShouldNot(HaveOccurred())

			Eventually(func() bool {
				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, system)
				Expect(err).ToNot(HaveOccurred())
				return len(system.GetFinalizers()) > 0
			}, timeout, poll).Should(BeTrue())

		})

		It("creates the required system-app resources", func() {

			dep := &appsv1.Deployment{}
			By("deploying a system-app workload",
				checkWorkloadResources(dep,
					expectedWorkload{
						Name:          "system-app",
						Namespace:     namespace,
						Replicas:      2,
						ContainerName: "system-app",
						ContainterArgs: []string{
							"env", "PORT=3000", "container-entrypoint", "bundle", "exec",
							"unicorn", "-c", "config/unicorn.rb",
						},
						PDB:        true,
						HPA:        true,
						PodMonitor: true,
					},
				),
			)
			Expect(dep.Spec.Template.Spec.Volumes[0].Secret.SecretName).To(Equal("system-config"))

			svc := &corev1.Service{}
			By("deploying the system-app service",
				checkResource(svc, expectedResource{
					Name: "system-app", Namespace: namespace,
				}),
			)
			Expect(svc.Spec.Selector["deployment"]).To(Equal("system-app"))

		})

		It("creates the required system-sidekiq resources", func() {

			dep := &appsv1.Deployment{}
			By("deploying a system-sidekiq-default workload",
				checkWorkloadResources(dep,
					expectedWorkload{
						Name:          "system-sidekiq-default",
						Namespace:     namespace,
						Replicas:      2,
						ContainerName: "system-sidekiq",
						ContainterArgs: []string{"sidekiq",
							"--queue", "critical", "--queue", "backend_sync",
							"--queue", "events", "--queue", "zync,40",
							"--queue", "priority,25", "--queue", "default,15",
							"--queue", "web_hooks,10", "--queue", "deletion,5",
						},
						PDB:        true,
						HPA:        true,
						PodMonitor: true,
					},
				),
			)
			Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
			Expect(dep.Spec.Template.Spec.Volumes[1].Secret.SecretName).To(Equal("system-config"))

			By("deploying a system-sidekiq-billing workload",
				checkWorkloadResources(dep,
					expectedWorkload{
						Name:           "system-sidekiq-billing",
						Namespace:      namespace,
						Replicas:       2,
						ContainerName:  "system-sidekiq",
						ContainterArgs: []string{"sidekiq", "--queue", "billing"},
						PDB:            true,
						HPA:            true,
						PodMonitor:     true,
					},
				),
			)
			Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
			Expect(dep.Spec.Template.Spec.Volumes[1].Secret.SecretName).To(Equal("system-config"))

			By("deploying a system-sidekiq-low workload",
				checkWorkloadResources(dep,
					expectedWorkload{
						Name:           "system-sidekiq-low",
						Namespace:      namespace,
						Replicas:       2,
						ContainerName:  "system-sidekiq",
						ContainterArgs: []string{"sidekiq", "--queue", "mailers", "--queue", "low"},
						PDB:            true,
						HPA:            true,
						PodMonitor:     true,
					},
				),
			)
			Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
			Expect(dep.Spec.Template.Spec.Volumes[1].Secret.SecretName).To(Equal("system-config"))

		})

		It("creates the system-sphinx resources", func() {

			sts := &appsv1.StatefulSet{}
			By("deploying the system-sphinx statefulset",
				checkResource(sts, expectedResource{
					Name: "system-sphinx", Namespace: namespace,
				}),
			)

			svc := &corev1.Service{}
			By("deploying the system-sphinx service",
				checkResource(svc, expectedResource{
					Name: "system-sphinx", Namespace: namespace,
				}),
			)
			Expect(svc.Spec.Selector["deployment"]).To(Equal("system-sphinx"))

		})

		It("creates the required system shared resources", func() {

			gd := &grafanav1alpha1.GrafanaDashboard{}
			By("deploying the system grafana dashboard",
				checkResource(gd, expectedResource{
					Name: "system", Namespace: namespace,
				}),
			)

			for _, esn := range []string{
				"system-database",
				"system-recaptcha",
				"system-events-hook",
				"system-smtp",
				"system-master-apicast",
				"system-zync",
				"system-backend",
				"system-multitenant-assets-s3",
				"system-app",
			} {
				es := &externalsecretsv1alpha1.ExternalSecret{}

				By("deploying the system external secret",
					checkResource(es, expectedResource{
						Name: esn, Namespace: namespace,
					}),
				)
			}

		})

		It("doesn't creates the non-default resources", func() {

			sts := &appsv1.StatefulSet{}
			By("ensuring the system-console statefulset",
				checkResource(sts, expectedResource{
					Name:      "system-console",
					Namespace: namespace, Missing: true,
				}),
			)

			dep := &appsv1.Deployment{}
			By("ensuring the system-app-canary deployment",
				checkResource(dep, expectedResource{
					Name:      "system-app-canary",
					Namespace: namespace, Missing: true,
				}),
			)

			By("ensuring the system-sidekiq-default-canary deployment",
				checkResource(dep, expectedResource{
					Name:      "system-sidekiq-default-canary",
					Namespace: namespace, Missing: true,
				}),
			)

			By("ensuring the system-sidekiq-billing-canary deployment",
				checkResource(dep, expectedResource{
					Name:      "system-sidekiq-billing-canary",
					Namespace: namespace, Missing: true,
				}),
			)

			By("ensuring the system-sidekiq-low-canary deployment",
				checkResource(dep, expectedResource{
					Name:      "system-sidekiq-low-canary",
					Namespace: namespace, Missing: true,
				}),
			)

		})

		When("updating a System resource with console", func() {

			BeforeEach(func() {
				Eventually(func() error {
					err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						system,
					)
					Expect(err).ToNot(HaveOccurred())
					patch := client.MergeFrom(system.DeepCopy())
					system.Spec.Image = &saasv1alpha1.ImageSpec{
						Name: pointer.StringPtr("newImage"),
						Tag:  pointer.StringPtr("newTag"),
					}
					system.Spec.Config.Rails = &saasv1alpha1.SystemRailsSpec{
						Console: pointer.Bool(true),
					}
					return k8sClient.Patch(context.Background(), system, patch)
				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("creates the required console resources", func() {

				sts := &appsv1.StatefulSet{}
				By("deploying the system-console StatefulSet",
					checkResource(sts, expectedResource{
						Name: "system-console", Namespace: namespace,
					}),
				)
				Expect(sts.Spec.Template.Spec.Containers[0].Image).Should((Equal("newImage:newTag")))
				Expect(sts.Spec.Template.Spec.Volumes[0].Secret.SecretName).Should((Equal("system-config")))

				pdb := &policyv1beta1.PodDisruptionBudget{}
				By("ensuring the system-console PDB",
					checkResource(pdb, expectedResource{
						Name: "system-console", Namespace: namespace, Missing: true,
					}),
				)

				hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{}
				By("ensuring the system-console HPA",
					checkResource(hpa, expectedResource{
						Name: "system-console", Namespace: namespace, Missing: true,
					}),
				)

			})

		})

		When("updating a System resource with canary", func() {

			// Resource Versions
			rvs := make(map[string]string)

			BeforeEach(func() {
				Eventually(func() error {
					system := &saasv1alpha1.System{}
					if err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						system,
					); err != nil {
						return err
					}

					rvs["svc/system-app"] = getResourceVersion(
						&corev1.Service{}, "system-app", namespace,
					)
					rvs["deployment/system-app"] = getResourceVersion(
						&appsv1.Deployment{}, "system-app", namespace,
					)

					patch := client.MergeFrom(system.DeepCopy())
					system.Spec.App = &saasv1alpha1.SystemAppSpec{
						Canary: &saasv1alpha1.Canary{
							ImageName: pointer.StringPtr("newImage"),
							ImageTag:  pointer.StringPtr("newTag"),
							Replicas:  pointer.Int32Ptr(2)},
					}
					system.Spec.SidekiqDefault = &saasv1alpha1.SystemSidekiqSpec{
						Canary: &saasv1alpha1.Canary{
							ImageName: pointer.StringPtr("newImage"),
							ImageTag:  pointer.StringPtr("newTag"),
							Replicas:  pointer.Int32Ptr(2)},
					}
					system.Spec.SidekiqBilling = &saasv1alpha1.SystemSidekiqSpec{
						Canary: &saasv1alpha1.Canary{
							ImageName: pointer.StringPtr("newImage"),
							ImageTag:  pointer.StringPtr("newTag"),
							Replicas:  pointer.Int32Ptr(2)},
					}
					system.Spec.SidekiqLow = &saasv1alpha1.SystemSidekiqSpec{
						Canary: &saasv1alpha1.Canary{
							ImageName: pointer.StringPtr("newImage"),
							ImageTag:  pointer.StringPtr("newTag"),
							Replicas:  pointer.Int32Ptr(2)},
					}
					return k8sClient.Patch(context.Background(), system, patch)
				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("creates the required canary resources", func() {

				dep := &appsv1.Deployment{}
				By("deploying a system-app-canary workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:          "system-app-canary",
							Namespace:     namespace,
							Replicas:      2,
							ContainerName: "system-app",
							ContainterArgs: []string{
								"env", "PORT=3000", "container-entrypoint", "bundle", "exec",
								"unicorn", "-c", "config/unicorn.rb",
							},
							PodMonitor: true,
						},
					),
				)
				Expect(dep.Spec.Template.Spec.Volumes[0].Secret.SecretName).To(Equal("system-config"))

				svc := &corev1.Service{}
				By("keeps the system-app service deployment label selector",
					checkResource(svc, expectedResource{
						Name: "system-app", Namespace: namespace,
					}),
				)
				Expect(svc.Spec.Selector["deployment"]).To(Equal("system-app"))
				Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("system-app"))

				By("deploying a system-sidekiq-default-canary workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:          "system-sidekiq-default-canary",
							Namespace:     namespace,
							Replicas:      2,
							ContainerName: "system-sidekiq",
							ContainterArgs: []string{"sidekiq",
								"--queue", "critical", "--queue", "backend_sync",
								"--queue", "events", "--queue", "zync,40",
								"--queue", "priority,25", "--queue", "default,15",
								"--queue", "web_hooks,10", "--queue", "deletion,5",
							},
							PodMonitor: true,
						},
					),
				)
				Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
				Expect(dep.Spec.Template.Spec.Volumes[1].Secret.SecretName).To(Equal("system-config"))

				By("deploying a system-sidekiq-billing-canary workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:           "system-sidekiq-billing-canary",
							Namespace:      namespace,
							Replicas:       2,
							ContainerName:  "system-sidekiq",
							ContainterArgs: []string{"sidekiq", "--queue", "billing"},
							PodMonitor:     true,
						},
					),
				)
				Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
				Expect(dep.Spec.Template.Spec.Volumes[1].Secret.SecretName).To(Equal("system-config"))

				By("deploying a system-sidekiq-low-canary workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:           "system-sidekiq-low-canary",
							Namespace:      namespace,
							Replicas:       2,
							ContainerName:  "system-sidekiq",
							ContainterArgs: []string{"sidekiq", "--queue", "mailers", "--queue", "low"},
							PodMonitor:     true,
						},
					),
				)
				Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
				Expect(dep.Spec.Template.Spec.Volumes[1].Secret.SecretName).To(Equal("system-config"))

			})

			When("enabling canary traffic", func() {

				BeforeEach(func() {
					Eventually(func() error {
						system := &saasv1alpha1.System{}
						if err := k8sClient.Get(
							context.Background(),
							types.NamespacedName{Name: "instance", Namespace: namespace},
							system,
						); err != nil {
							return err
						}
						rvs["svc/system-app"] = getResourceVersion(
							&corev1.Service{}, "system-app", namespace,
						)
						patch := client.MergeFrom(system.DeepCopy())
						system.Spec.App = &saasv1alpha1.SystemAppSpec{
							Canary: &saasv1alpha1.Canary{
								SendTraffic: *pointer.Bool(true),
							},
						}
						return k8sClient.Patch(context.Background(), system, patch)
					}, timeout, poll).ShouldNot(HaveOccurred())
				})

				It("updates the system-app service", func() {

					svc := &corev1.Service{}
					By("removing the system-app service deployment label selector",
						checkResource(svc, expectedResource{
							Name: "system-app", Namespace: namespace,
							LastVersion: rvs["svc/system-app"],
						}),
					)
					Expect(svc.Spec.Selector).NotTo(HaveKey("deployment"))
					Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("system-app"))

				})

			})
		})

		When("updating a system resource with twemproxyconfig", func() {

			// Resource Versions
			rvs := make(map[string]string)

			BeforeEach(func() {
				Eventually(func() error {

					system := &saasv1alpha1.System{}
					if err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						system,
					); err != nil {
						return err
					}

					rvs["deployment/system-app"] = getResourceVersion(
						&appsv1.Deployment{}, "system-app", namespace,
					)
					rvs["deployment/system-sidekiq-billing"] = getResourceVersion(
						&appsv1.Deployment{}, "system-sidekiq-billing", namespace,
					)
					rvs["deployment/system-sidekiq-default"] = getResourceVersion(
						&appsv1.Deployment{}, "system-sidekiq-default", namespace,
					)
					rvs["deployment/system-sidekiq-low"] = getResourceVersion(
						&appsv1.Deployment{}, "system-sidekiq-low", namespace,
					)

					patch := client.MergeFrom(system.DeepCopy())

					system.Spec.Twemproxy = &saasv1alpha1.TwemproxySpec{
						TwemproxyConfigRef: "system-twemproxyconfig",
						Options: &saasv1alpha1.TwemproxyOptions{
							LogLevel: pointer.Int32Ptr(2),
						},
					}

					system.Spec.App = &saasv1alpha1.SystemAppSpec{
						Canary: &saasv1alpha1.Canary{
							Replicas: pointer.Int32(2),
							Patches: []string{
								`[{"op":"add","path":"/twemproxy","value":{"twemproxyConfigRef":"system-canary-twemproxyconfig","options":{"logLevel":2}}}]`,
							},
						},
					}

					system.Spec.SidekiqBilling = &saasv1alpha1.SystemSidekiqSpec{
						Canary: &saasv1alpha1.Canary{
							Replicas: pointer.Int32(3),
							Patches: []string{
								`[{"op":"add","path":"/twemproxy","value":{"twemproxyConfigRef":"system-canary-twemproxyconfig","options":{"logLevel":3}}}]`,
							},
						},
					}

					system.Spec.SidekiqDefault = &saasv1alpha1.SystemSidekiqSpec{
						Replicas: pointer.Int32(2),
						Canary: &saasv1alpha1.Canary{
							Replicas: pointer.Int32(4),
							Patches: []string{
								`[{"op":"add","path":"/twemproxy/options","value":{"logLevel":4}}]`,
							},
						},
					}

					system.Spec.SidekiqLow = &saasv1alpha1.SystemSidekiqSpec{
						Replicas: pointer.Int32(2),
						Canary: &saasv1alpha1.Canary{
							Replicas: pointer.Int32(5),
							Patches: []string{
								`[{"op":"add","path":"/twemproxy/options","value":{"logLevel":5}}]`,
							},
						},
					}

					return k8sClient.Patch(context.Background(), system, patch)

				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("updates the system-app resources", func() {

				dep := &appsv1.Deployment{}
				By("adding a twemproxy sidecar to the system-app workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:        "system-app",
							Namespace:   namespace,
							Replicas:    2,
							PDB:         true,
							HPA:         true,
							PodMonitor:  true,
							LastVersion: rvs["deployment/system-app"],
						},
					),
				)
				Expect(dep.Spec.Template.Spec.Volumes).To(HaveLen(2))
				Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-config"))
				Expect(dep.Spec.Template.Spec.Volumes[0].VolumeSource.Secret.SecretName).To(Equal("system-config"))
				Expect(dep.Spec.Template.Spec.Volumes[1].Name).To(Equal("twemproxy-config"))
				Expect(dep.Spec.Template.Spec.Volumes[1].VolumeSource.ConfigMap.LocalObjectReference.Name).To(Equal("system-twemproxyconfig"))
				Expect(dep.Spec.Template.Spec.Containers).To(HaveLen(2))
				Expect(dep.Spec.Template.Spec.Containers[1].Name).To(Equal("twemproxy"))
				Expect(dep.Spec.Template.Spec.Containers[1].VolumeMounts[0].Name).To(Equal("twemproxy-config"))

				By("adding a twemproxy sidecar to the system-app-canary workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:       "system-app-canary",
							Replicas:   2,
							Namespace:  namespace,
							PodMonitor: true,
						},
					),
				)
				Expect(dep.Spec.Template.Spec.Volumes).To(HaveLen(2))
				Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-config"))
				Expect(dep.Spec.Template.Spec.Volumes[0].VolumeSource.Secret.SecretName).To(Equal("system-config"))
				Expect(dep.Spec.Template.Spec.Volumes[1].Name).To(Equal("twemproxy-config"))
				Expect(dep.Spec.Template.Spec.Volumes[1].VolumeSource.ConfigMap.LocalObjectReference.Name).To(Equal("system-canary-twemproxyconfig"))
				Expect(dep.Spec.Template.Spec.Containers).To(HaveLen(2))
				Expect(dep.Spec.Template.Spec.Containers[1].Name).To(Equal("twemproxy"))
				Expect(dep.Spec.Template.Spec.Containers[1].VolumeMounts[0].Name).To(Equal("twemproxy-config"))
				Expect(dep.Spec.Template.Spec.Containers[1].Env[3].Name).To(Equal("TWEMPROXY_LOG_LEVEL"))
				Expect(dep.Spec.Template.Spec.Containers[1].Env[3].Value).To(Equal("2"))

				By("adding a twemproxy sidecar to the system-sidekiq-billing workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:        "system-sidekiq-billing",
							Namespace:   namespace,
							Replicas:    2,
							PDB:         true,
							HPA:         true,
							PodMonitor:  true,
							LastVersion: rvs["deployment/system-sidekiq-billing"],
						},
					),
				)
				Expect(dep.Spec.Template.Spec.Volumes).To(HaveLen(3))
				Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
				Expect(dep.Spec.Template.Spec.Volumes[0].VolumeSource.EmptyDir).ShouldNot(BeNil())
				Expect(dep.Spec.Template.Spec.Volumes[1].Name).To(Equal("system-config"))
				Expect(dep.Spec.Template.Spec.Volumes[1].VolumeSource.Secret.SecretName).To(Equal("system-config"))
				Expect(dep.Spec.Template.Spec.Volumes[2].Name).To(Equal("twemproxy-config"))
				Expect(dep.Spec.Template.Spec.Volumes[2].VolumeSource.ConfigMap.LocalObjectReference.Name).To(Equal("system-twemproxyconfig"))
				Expect(dep.Spec.Template.Spec.Containers).To(HaveLen(2))
				Expect(dep.Spec.Template.Spec.Containers[1].Name).To(Equal("twemproxy"))
				Expect(dep.Spec.Template.Spec.Containers[1].VolumeMounts[0].Name).To(Equal("twemproxy-config"))
				Expect(dep.Spec.Template.Spec.Containers[1].Env[3].Name).To(Equal("TWEMPROXY_LOG_LEVEL"))
				Expect(dep.Spec.Template.Spec.Containers[1].Env[3].Value).To(Equal("2"))

				By("adding a twemproxy sidecar to the system-sidekiq-billing-canary workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:       "system-sidekiq-billing-canary",
							Replicas:   3,
							Namespace:  namespace,
							PodMonitor: true,
						},
					),
				)
				Expect(dep.Spec.Template.Spec.Volumes).To(HaveLen(3))
				Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
				Expect(dep.Spec.Template.Spec.Volumes[0].VolumeSource.EmptyDir).ShouldNot(BeNil())
				Expect(dep.Spec.Template.Spec.Volumes[1].Name).To(Equal("system-config"))
				Expect(dep.Spec.Template.Spec.Volumes[1].VolumeSource.Secret.SecretName).To(Equal("system-config"))
				Expect(dep.Spec.Template.Spec.Volumes[2].Name).To(Equal("twemproxy-config"))
				Expect(dep.Spec.Template.Spec.Volumes[2].VolumeSource.ConfigMap.LocalObjectReference.Name).To(Equal("system-canary-twemproxyconfig"))
				Expect(dep.Spec.Template.Spec.Containers).To(HaveLen(2))
				Expect(dep.Spec.Template.Spec.Containers[1].Name).To(Equal("twemproxy"))
				Expect(dep.Spec.Template.Spec.Containers[1].VolumeMounts[0].Name).To(Equal("twemproxy-config"))
				Expect(dep.Spec.Template.Spec.Containers[1].Env[3].Name).To(Equal("TWEMPROXY_LOG_LEVEL"))
				Expect(dep.Spec.Template.Spec.Containers[1].Env[3].Value).To(Equal("3"))

				By("adding a twemproxy sidecar to the system-sidekiq-low workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:        "system-sidekiq-low",
							Namespace:   namespace,
							Replicas:    2,
							PDB:         true,
							HPA:         true,
							PodMonitor:  true,
							LastVersion: rvs["deployment/system-sidekiq-low"],
						},
					),
				)
				Expect(dep.Spec.Template.Spec.Volumes).To(HaveLen(3))
				Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
				Expect(dep.Spec.Template.Spec.Volumes[0].VolumeSource.EmptyDir).ShouldNot(BeNil())
				Expect(dep.Spec.Template.Spec.Volumes[1].Name).To(Equal("system-config"))
				Expect(dep.Spec.Template.Spec.Volumes[1].VolumeSource.Secret.SecretName).To(Equal("system-config"))
				Expect(dep.Spec.Template.Spec.Volumes[2].Name).To(Equal("twemproxy-config"))
				Expect(dep.Spec.Template.Spec.Volumes[2].VolumeSource.ConfigMap.LocalObjectReference.Name).To(Equal("system-twemproxyconfig"))
				Expect(dep.Spec.Template.Spec.Containers).To(HaveLen(2))
				Expect(dep.Spec.Template.Spec.Containers[1].Name).To(Equal("twemproxy"))
				Expect(dep.Spec.Template.Spec.Containers[1].VolumeMounts[0].Name).To(Equal("twemproxy-config"))
				Expect(dep.Spec.Template.Spec.Containers[1].Env[3].Name).To(Equal("TWEMPROXY_LOG_LEVEL"))
				Expect(dep.Spec.Template.Spec.Containers[1].Env[3].Value).To(Equal("2"))

				By("adding a twemproxy sidecar to the system-sidekiq-low-canary workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:       "system-sidekiq-low-canary",
							Replicas:   5,
							Namespace:  namespace,
							PodMonitor: true,
						},
					),
				)
				Expect(dep.Spec.Template.Spec.Volumes).To(HaveLen(3))
				Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
				Expect(dep.Spec.Template.Spec.Volumes[0].VolumeSource.EmptyDir).ShouldNot(BeNil())
				Expect(dep.Spec.Template.Spec.Volumes[1].Name).To(Equal("system-config"))
				Expect(dep.Spec.Template.Spec.Volumes[1].VolumeSource.Secret.SecretName).To(Equal("system-config"))
				Expect(dep.Spec.Template.Spec.Volumes[2].Name).To(Equal("twemproxy-config"))
				Expect(dep.Spec.Template.Spec.Volumes[2].VolumeSource.ConfigMap.LocalObjectReference.Name).To(Equal("system-twemproxyconfig"))
				Expect(dep.Spec.Template.Spec.Containers).To(HaveLen(2))
				Expect(dep.Spec.Template.Spec.Containers[1].Name).To(Equal("twemproxy"))
				Expect(dep.Spec.Template.Spec.Containers[1].VolumeMounts[0].Name).To(Equal("twemproxy-config"))
				Expect(dep.Spec.Template.Spec.Containers[1].Env[3].Name).To(Equal("TWEMPROXY_LOG_LEVEL"))
				Expect(dep.Spec.Template.Spec.Containers[1].Env[3].Value).To(Equal("5"))

			})
		})

	})
})
