package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/3scale-ops/basereconciler/util"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	testutil "github.com/3scale-ops/saas-operator/test/util"
	externalsecretsv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	grafanav1beta1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	pipelinev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
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
								Path: "some-path-db",
								Key:  "some-key-db",
							},
						},
						EventsSharedSecret: saasv1alpha1.SecretReference{Override: util.Pointer("override")},
						Recaptcha: saasv1alpha1.SystemRecaptchaSpec{
							PublicKey:  saasv1alpha1.SecretReference{Override: util.Pointer("override")},
							PrivateKey: saasv1alpha1.SecretReference{Override: util.Pointer("override")},
						},
						SecretKeyBase: saasv1alpha1.SecretReference{Override: util.Pointer("override")},
						AccessCode:    &saasv1alpha1.SecretReference{Override: util.Pointer("override")},
						Segment: saasv1alpha1.SegmentSpec{
							DeletionWorkspace: "value",
							DeletionToken:     saasv1alpha1.SecretReference{Override: util.Pointer("override")},
							WriteKey:          saasv1alpha1.SecretReference{Override: util.Pointer("override")},
						},
						Github: saasv1alpha1.GithubSpec{
							ClientID:     saasv1alpha1.SecretReference{Override: util.Pointer("override")},
							ClientSecret: saasv1alpha1.SecretReference{Override: util.Pointer("override")},
						},
						RedHatCustomerPortal: saasv1alpha1.RedHatCustomerPortalSpec{
							ClientID:     saasv1alpha1.SecretReference{Override: util.Pointer("override")},
							ClientSecret: saasv1alpha1.SecretReference{Override: util.Pointer("override")},
							Realm:        util.Pointer("sso.example.net"),
						},
						Bugsnag: &saasv1alpha1.BugsnagSpec{
							ReleaseStage: util.Pointer("staging"),
							APIKey:       saasv1alpha1.SecretReference{Override: util.Pointer("override")},
						},
						// DatabaseSecret:   saasv1alpha1.SecretReference{Override: util.Pointer("override")},
						DatabaseSecret: saasv1alpha1.SecretReference{FromVault: &saasv1alpha1.VaultSecretReference{
							Path: "path",
							Key:  "key",
						}},
						MemcachedServers: "value",
						Redis: saasv1alpha1.RedisSpec{
							QueuesDSN: "value",
						},
						SMTP: saasv1alpha1.SMTPSpec{
							Address:           "value",
							User:              saasv1alpha1.SecretReference{Override: util.Pointer("override")},
							Password:          saasv1alpha1.SecretReference{Override: util.Pointer("override")},
							Port:              1000,
							AuthProtocol:      "value",
							OpenSSLVerifyMode: "value",
							STARTTLS:          util.Pointer(false),
						},
						MappingServiceAccessToken: saasv1alpha1.SecretReference{Override: util.Pointer("override")},
						Zync: saasv1alpha1.SystemZyncSpec{
							AuthToken: saasv1alpha1.SecretReference{Override: util.Pointer("override")},
							Endpoint:  "value",
						},
						Backend: saasv1alpha1.SystemBackendSpec{
							ExternalEndpoint:    "value",
							InternalEndpoint:    "value",
							InternalAPIUser:     saasv1alpha1.SecretReference{Override: util.Pointer("override")},
							InternalAPIPassword: saasv1alpha1.SecretReference{Override: util.Pointer("override")},
							RedisDSN:            "value",
						},
						Assets: saasv1alpha1.AssetsSpec{
							Host:      util.Pointer("test.cloudfront.net"),
							Bucket:    "bucket",
							Region:    "us-east-1",
							AccessKey: saasv1alpha1.SecretReference{Override: util.Pointer("override")},
							SecretKey: saasv1alpha1.SecretReference{Override: util.Pointer("override")},
						},
					},
					App: &saasv1alpha1.SystemAppSpec{
						DeploymentStrategy: &saasv1alpha1.DeploymentStrategySpec{
							Type: appsv1.RollingUpdateDeploymentStrategyType,
							RollingUpdate: &appsv1.RollingUpdateDeployment{
								MaxSurge:       util.Pointer(intstr.FromString("20%")),
								MaxUnavailable: util.Pointer(intstr.FromInt(0)),
							},
						}},
					SidekiqDefault: &saasv1alpha1.SystemSidekiqSpec{
						DeploymentStrategy: &saasv1alpha1.DeploymentStrategySpec{
							Type: appsv1.RollingUpdateDeploymentStrategyType,
							RollingUpdate: &appsv1.RollingUpdateDeployment{
								MaxSurge:       util.Pointer(intstr.FromString("15%")),
								MaxUnavailable: util.Pointer(intstr.FromString("5%")),
							},
						},
						HPA: &saasv1alpha1.HorizontalPodAutoscalerSpec{
							Behavior: &autoscalingv2.HorizontalPodAutoscalerBehavior{
								ScaleUp: &autoscalingv2.HPAScalingRules{
									Policies: []autoscalingv2.HPAScalingPolicy{
										{
											Type:          autoscalingv2.PodsScalingPolicy,
											Value:         4,
											PeriodSeconds: 60,
										},
										{
											Type:          autoscalingv2.PercentScalingPolicy,
											Value:         10,
											PeriodSeconds: 60,
										},
									},
								},
							},
						},
					},
				},
			}
			err := k8sClient.Create(context.Background(), system)
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, system)
			}, timeout, poll).ShouldNot(HaveOccurred())

		})

		It("creates the required system-app resources", func() {

			dep := &appsv1.Deployment{}
			By("deploying a system-app workload",
				(&testutil.ExpectedWorkload{
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
				}).Assert(k8sClient, dep, timeout, poll))

			for _, env := range dep.Spec.Template.Spec.Containers[0].Env {
				switch env.Name {
				case "THINKING_SPHINX_ADDRESS":
					Expect(env.Value).To(Equal("system-searchd"))
				case "THINKING_SPHINX_PORT":
					Expect(env.Value).To(Equal("9306"))
				}
			}

			Expect(dep.Spec.Template.Spec.Volumes[0].Secret.SecretName).To(Equal("system-config"))
			Expect(dep.Spec.Strategy.Type).To(Equal(appsv1.RollingUpdateDeploymentStrategyType))
			Expect(dep.Spec.Strategy.RollingUpdate.MaxSurge).To(Equal(util.Pointer(intstr.FromString("20%"))))
			Expect(dep.Spec.Strategy.RollingUpdate.MaxUnavailable).To(Equal(util.Pointer(intstr.FromInt(0))))

			svc := &corev1.Service{}
			By("deploying the system-app service",
				(&testutil.ExpectedResource{Name: "system-app-http-svc", Namespace: namespace}).
					Assert(k8sClient, svc, timeout, poll))

			Expect(svc.Spec.Selector["deployment"]).To(Equal("system-app"))
			Expect(dep.Spec.Template.Spec.TerminationGracePeriodSeconds).To(Equal(util.Pointer[int64](60)))

		})

		It("creates the required system-sidekiq resources", func() {

			dep := &appsv1.Deployment{}
			By("deploying a system-sidekiq-default workload",
				(&testutil.ExpectedWorkload{
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
				}).Assert(k8sClient, dep, timeout, poll))

			Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
			Expect(dep.Spec.Template.Spec.Volumes[1].Secret.SecretName).To(Equal("system-config"))
			Expect(dep.Spec.Template.Spec.TerminationGracePeriodSeconds).To(Equal(util.Pointer[int64](60)))
			Expect(dep.Spec.Strategy.Type).To(Equal(appsv1.RollingUpdateDeploymentStrategyType))
			Expect(dep.Spec.Strategy.RollingUpdate.MaxSurge).To(Equal(util.Pointer(intstr.FromString("15%"))))
			Expect(dep.Spec.Strategy.RollingUpdate.MaxUnavailable).To(Equal(util.Pointer(intstr.FromString("5%"))))

			hpa := &autoscalingv2.HorizontalPodAutoscaler{}
			By("updates system-sidekiq-default hpa behaviour",
				(&testutil.ExpectedResource{Name: "system-sidekiq-default", Namespace: namespace}).
					Assert(k8sClient, hpa, timeout, poll))
			Expect(hpa.Spec.Behavior.ScaleUp.Policies).To(Not(BeEmpty()))
			Expect(hpa.Spec.Behavior.ScaleUp.Policies[0].Type).To(Equal(autoscalingv2.PodsScalingPolicy))
			Expect(hpa.Spec.Behavior.ScaleUp.Policies[0].Value).To(Equal(int32(4)))

			By("deploying a system-sidekiq-billing workload",
				(&testutil.ExpectedWorkload{
					Name:           "system-sidekiq-billing",
					Namespace:      namespace,
					Replicas:       2,
					ContainerName:  "system-sidekiq",
					ContainterArgs: []string{"sidekiq", "--queue", "billing"},
					PDB:            true,
					HPA:            true,
					PodMonitor:     true,
				}).Assert(k8sClient, dep, timeout, poll))

			Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
			Expect(dep.Spec.Template.Spec.Volumes[1].Secret.SecretName).To(Equal("system-config"))
			Expect(dep.Spec.Template.Spec.TerminationGracePeriodSeconds).To(Equal(util.Pointer[int64](60)))
			Expect(dep.Spec.Strategy.Type).To(Equal(appsv1.RollingUpdateDeploymentStrategyType))
			Expect(dep.Spec.Strategy.RollingUpdate.MaxSurge).To(Equal(util.Pointer(intstr.FromInt(1))))
			Expect(dep.Spec.Strategy.RollingUpdate.MaxUnavailable).To(Equal(util.Pointer(intstr.FromInt(0))))

			By("deploying a system-sidekiq-low workload",
				(&testutil.ExpectedWorkload{
					Name:           "system-sidekiq-low",
					Namespace:      namespace,
					Replicas:       2,
					ContainerName:  "system-sidekiq",
					ContainterArgs: []string{"sidekiq", "--queue", "mailers", "--queue", "low", "--queue", "bulk_indexing"},
					PDB:            true,
					HPA:            true,
					PodMonitor:     true,
				}).Assert(k8sClient, dep, timeout, poll))

			Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
			Expect(dep.Spec.Template.Spec.Volumes[1].Secret.SecretName).To(Equal("system-config"))
			Expect(dep.Spec.Template.Spec.TerminationGracePeriodSeconds).To(Equal(util.Pointer[int64](60)))

		})

		It("creates the system-searchd resources", func() {

			sts := &appsv1.StatefulSet{}
			By("deploying the system-searchd statefulset",
				(&testutil.ExpectedResource{Name: "system-searchd", Namespace: namespace}).
					Assert(k8sClient, sts, timeout, poll))

			Expect(sts.Spec.Template.Spec.Containers[0].Args).To(BeEmpty())
			Expect(sts.Spec.Template.Spec.TerminationGracePeriodSeconds).To(Equal(util.Pointer[int64](60)))
			Expect(sts.Spec.Template.Spec.Containers[0].Env).To(BeEmpty())

			svc := &corev1.Service{}
			By("deploying the system-searchd service",
				(&testutil.ExpectedResource{Name: "system-searchd", Namespace: namespace}).
					Assert(k8sClient, svc, timeout, poll))
			Expect(svc.Spec.Selector["deployment"]).To(Equal("system-searchd"))

		})

		It("creates the required system shared resources", func() {

			gd := &grafanav1beta1.GrafanaDashboard{}
			By("deploying the system grafana dashboard",
				(&testutil.ExpectedResource{Name: "system", Namespace: namespace}).
					Assert(k8sClient, gd, timeout, poll))

			for _, esn := range []string{
				"system-database",
				// "system-recaptcha",
				// "system-events-hook",
				// "system-smtp",
				// "system-master-apicast",
				// "system-zync",
				// "system-backend",
				// "system-multitenant-assets-s3",
				"system-app",
			} {
				es := &externalsecretsv1beta1.ExternalSecret{}

				By("deploying the system external secret",
					(&testutil.ExpectedResource{Name: esn, Namespace: namespace}).
						Assert(k8sClient, es, timeout, poll))
			}

			es := &externalsecretsv1beta1.ExternalSecret{}
			By("deploying the system-database external secret with specific configuration",
				(&testutil.ExpectedResource{Name: "system-database", Namespace: namespace}).
					Assert(k8sClient, es, timeout, poll))

			Expect(es.Spec.RefreshInterval.ToUnstructured()).To(Equal("1m0s"))
			Expect(es.Spec.SecretStoreRef.Name).To(Equal("vault-mgmt"))
			Expect(es.Spec.SecretStoreRef.Kind).To(Equal("ClusterSecretStore"))

			for _, data := range es.Spec.Data {
				switch data.SecretKey {
				case "DATABASE_URL":
					Expect(data.RemoteRef.Property).To(Equal("some-key-db"))
					Expect(data.RemoteRef.Key).To(Equal("some-path-db"))
				}
			}
		})

		It("creates the required tekton resources", func() {

			for _, tr := range []string{
				"system-db-migrate",
				"system-searchd-reindex",
				"system-backend-sync",
			} {
				task := &pipelinev1beta1.Task{}
				pipeline := &pipelinev1beta1.Pipeline{}

				By("deploying the system task",
					(&testutil.ExpectedResource{Name: tr, Namespace: namespace}).
						Assert(k8sClient, task, timeout, poll))

				By("deploying the system pipeline",
					(&testutil.ExpectedResource{Name: tr, Namespace: namespace}).
						Assert(k8sClient, pipeline, timeout, poll))
			}

		})

		It("doesn't creates the non-default resources", func() {

			sts := &appsv1.StatefulSet{}
			By("ensuring the system-console statefulset",
				(&testutil.ExpectedResource{Name: "system-console", Namespace: namespace, Missing: true}).
					Assert(k8sClient, sts, timeout, poll))

			dep := &appsv1.Deployment{}
			By("ensuring the system-app-canary deployment",
				(&testutil.ExpectedResource{Name: "system-app-canary", Namespace: namespace, Missing: true}).
					Assert(k8sClient, dep, timeout, poll))

			By("ensuring the system-sidekiq-default-canary deployment",
				(&testutil.ExpectedResource{Name: "system-sidekiq-default-canary", Namespace: namespace, Missing: true}).
					Assert(k8sClient, dep, timeout, poll))

			By("ensuring the system-sidekiq-billing-canary deployment",
				(&testutil.ExpectedResource{Name: "system-sidekiq-billing-canary", Namespace: namespace, Missing: true}).
					Assert(k8sClient, dep, timeout, poll))

			By("ensuring the system-sidekiq-low-canary deployment",
				(&testutil.ExpectedResource{Name: "system-sidekiq-low-canary", Namespace: namespace, Missing: true}).
					Assert(k8sClient, dep, timeout, poll))

		})

		When("updating a System resource without searchd", func() {

			BeforeEach(func() {
				Eventually(func() error {
					err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						system,
					)
					Expect(err).ToNot(HaveOccurred())

					patch := client.MergeFrom(system.DeepCopy())
					system.Spec.Searchd = &saasv1alpha1.SystemSearchdSpec{
						Enabled: util.Pointer(false),
					}
					return k8sClient.Patch(context.Background(), system, patch)
				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("removes the system-searchd resources", func() {

				sts := &appsv1.StatefulSet{}
				By("removing the system-searchd statefulset",
					(&testutil.ExpectedResource{Name: "system-searchd", Namespace: namespace, Missing: true}).
						Assert(k8sClient, sts, timeout, poll))

			})
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
					system.Spec.Config.Rails = &saasv1alpha1.SystemRailsSpec{
						Console: util.Pointer(true),
					}
					system.Spec.Console = &saasv1alpha1.SystemRailsConsoleSpec{
						Image: &saasv1alpha1.ImageSpec{
							Name: util.Pointer("newImage"),
							Tag:  util.Pointer("newTag"),
						},
					}
					return k8sClient.Patch(context.Background(), system, patch)
				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("creates the required console resources", func() {

				sts := &appsv1.StatefulSet{}
				By("deploying the system-console StatefulSet",
					(&testutil.ExpectedResource{Name: "system-console", Namespace: namespace}).
						Assert(k8sClient, sts, timeout, poll))

				Expect(sts.Spec.Template.Spec.Containers[0].Image).Should((Equal("newImage:newTag")))
				Expect(sts.Spec.Template.Spec.Volumes[0].Secret.SecretName).Should((Equal("system-config")))

				pdb := &policyv1.PodDisruptionBudget{}
				By("ensuring the system-console PDB",
					(&testutil.ExpectedResource{Name: "system-console", Namespace: namespace, Missing: true}).
						Assert(k8sClient, pdb, timeout, poll))

				hpa := &autoscalingv2.HorizontalPodAutoscaler{}
				By("ensuring the system-console HPA",
					(&testutil.ExpectedResource{Name: "system-console", Namespace: namespace, Missing: true}).
						Assert(k8sClient, hpa, timeout, poll))

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

					rvs["svc/system-app-http-svc"] = testutil.GetResourceVersion(
						k8sClient, &corev1.Service{}, "system-app-http-svc", namespace, timeout, poll)
					rvs["deployment/system-app"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "system-app", namespace, timeout, poll)

					patch := client.MergeFrom(system.DeepCopy())
					system.Spec.App = &saasv1alpha1.SystemAppSpec{
						Canary: &saasv1alpha1.Canary{
							ImageName: util.Pointer("newImage"),
							ImageTag:  util.Pointer("newTag"),
							Replicas:  util.Pointer[int32](2)},
					}
					system.Spec.SidekiqDefault = &saasv1alpha1.SystemSidekiqSpec{
						Canary: &saasv1alpha1.Canary{
							ImageName: util.Pointer("newImage"),
							ImageTag:  util.Pointer("newTag"),
							Replicas:  util.Pointer[int32](2)},
					}
					system.Spec.SidekiqBilling = &saasv1alpha1.SystemSidekiqSpec{
						Canary: &saasv1alpha1.Canary{
							ImageName: util.Pointer("newImage"),
							ImageTag:  util.Pointer("newTag"),
							Replicas:  util.Pointer[int32](2)},
					}
					system.Spec.SidekiqLow = &saasv1alpha1.SystemSidekiqSpec{
						Canary: &saasv1alpha1.Canary{
							ImageName: util.Pointer("newImage"),
							ImageTag:  util.Pointer("newTag"),
							Replicas:  util.Pointer[int32](2)},
					}
					// return k8sClient.Patch(context.Background(), system, patch)

					if err := k8sClient.Patch(context.Background(), system, patch); err != nil {
						return err
					}

					// waiting for on of the Deployments is enough ...
					if testutil.GetResourceVersion(k8sClient, &appsv1.Deployment{}, "system-app-canary", namespace, timeout, poll) == "" {
						return fmt.Errorf("not ready")
					}

					return nil
				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("creates the required canary resources", func() {

				dep := &appsv1.Deployment{}
				By("deploying a system-app-canary workload",
					(&testutil.ExpectedWorkload{
						Name:          "system-app-canary",
						Namespace:     namespace,
						Replicas:      2,
						ContainerName: "system-app",
						ContainterArgs: []string{
							"env", "PORT=3000", "container-entrypoint", "bundle", "exec",
							"unicorn", "-c", "config/unicorn.rb",
						},
						PodMonitor: true,
					}).Assert(k8sClient, dep, timeout, poll))

				Expect(dep.Spec.Template.Spec.Volumes[0].Secret.SecretName).To(Equal("system-config"))
				Expect(dep.Spec.Template.Spec.TerminationGracePeriodSeconds).To(Equal(util.Pointer[int64](60)))

				svc := &corev1.Service{}
				By("keeps the system-app service deployment label selector",
					(&testutil.ExpectedResource{Name: "system-app-http-svc", Namespace: namespace}).
						Assert(k8sClient, svc, timeout, poll))

				Expect(svc.Spec.Selector["deployment"]).To(Equal("system-app"))
				Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("system-app"))

				By("deploying a system-sidekiq-default-canary workload",
					(&testutil.ExpectedWorkload{
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
					}).Assert(k8sClient, dep, timeout, poll))

				Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
				Expect(dep.Spec.Template.Spec.Volumes[1].Secret.SecretName).To(Equal("system-config"))
				Expect(dep.Spec.Template.Spec.TerminationGracePeriodSeconds).To(Equal(util.Pointer[int64](60)))

				By("deploying a system-sidekiq-billing-canary workload",
					(&testutil.ExpectedWorkload{
						Name:           "system-sidekiq-billing-canary",
						Namespace:      namespace,
						Replicas:       2,
						ContainerName:  "system-sidekiq",
						ContainterArgs: []string{"sidekiq", "--queue", "billing"},
						PodMonitor:     true,
					}).Assert(k8sClient, dep, timeout, poll))

				Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
				Expect(dep.Spec.Template.Spec.Volumes[1].Secret.SecretName).To(Equal("system-config"))
				Expect(dep.Spec.Template.Spec.TerminationGracePeriodSeconds).To(Equal(util.Pointer[int64](60)))

				By("deploying a system-sidekiq-low-canary workload",
					(&testutil.ExpectedWorkload{
						Name:           "system-sidekiq-low-canary",
						Namespace:      namespace,
						Replicas:       2,
						ContainerName:  "system-sidekiq",
						ContainterArgs: []string{"sidekiq", "--queue", "mailers", "--queue", "low", "--queue", "bulk_indexing"},
						PodMonitor:     true,
					}).Assert(k8sClient, dep, timeout, poll))

				Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
				Expect(dep.Spec.Template.Spec.Volumes[1].Secret.SecretName).To(Equal("system-config"))
				Expect(dep.Spec.Template.Spec.TerminationGracePeriodSeconds).To(Equal(util.Pointer[int64](60)))

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
						rvs["svc/system-app-http-svc"] = testutil.GetResourceVersion(
							k8sClient, &corev1.Service{}, "system-app-http-svc", namespace, timeout, poll)

						patch := client.MergeFrom(system.DeepCopy())
						system.Spec.App = &saasv1alpha1.SystemAppSpec{
							Canary: &saasv1alpha1.Canary{
								SendTraffic: *util.Pointer(true),
							},
						}
						return k8sClient.Patch(context.Background(), system, patch)
					}, timeout, poll).ShouldNot(HaveOccurred())
				})

				It("updates the system-app service", func() {

					svc := &corev1.Service{}
					By("removing the system-app service deployment label selector",
						(&testutil.ExpectedResource{
							Name:        "system-app-http-svc",
							Namespace:   namespace,
							LastVersion: rvs["svc/system-app-http-svc"],
						}).Assert(k8sClient, svc, timeout, poll))

					Expect(svc.Spec.Selector).NotTo(HaveKey("deployment"))
					Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("system-app"))

				})

			})

			When("disabling the canary", func() {

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
						patch := client.MergeFrom(system.DeepCopy())
						system.Spec.App.Canary = nil
						system.Spec.SidekiqDefault.Canary = nil
						system.Spec.SidekiqBilling.Canary = nil
						system.Spec.SidekiqLow.Canary = nil
						return k8sClient.Patch(context.Background(), system, patch)
					}, timeout, poll).ShouldNot(HaveOccurred())
				})

				It("deletes the canary resources", func() {

					dep := &appsv1.Deployment{}
					By("removing the system-app-canary Deployment",
						(&testutil.ExpectedResource{
							Name: "system-app-canary", Namespace: namespace, Missing: true}).Assert(k8sClient, dep, timeout, poll))
					By("removing the system-sidekiq-default-canary Deployment",
						(&testutil.ExpectedResource{
							Name: "system-sidekiq-default-canary", Namespace: namespace, Missing: true}).Assert(k8sClient, dep, timeout, poll))
					By("removing the system-sidekiq-billing-canary Deployment",
						(&testutil.ExpectedResource{
							Name: "system-sidekiq-billing-canary", Namespace: namespace, Missing: true}).Assert(k8sClient, dep, timeout, poll))
					By("removing the system-sidekiq-low-canary Deployment",
						(&testutil.ExpectedResource{
							Name: "system-sidekiq-low-canary", Namespace: namespace, Missing: true}).Assert(k8sClient, dep, timeout, poll))

					pm := &monitoringv1.PodMonitor{}
					By("removing the system-app-canary PodMonitor",
						(&testutil.ExpectedResource{
							Name: "system-app-canary", Namespace: namespace, Missing: true}).Assert(k8sClient, pm, timeout, poll))
					By("removing the system-sidekiq-default-canary PodMonitor",
						(&testutil.ExpectedResource{
							Name: "system-sidekiq-default-canary", Namespace: namespace, Missing: true}).Assert(k8sClient, pm, timeout, poll))
					By("removing the system-sidekiq-billing-canary PodMonitor",
						(&testutil.ExpectedResource{
							Name: "system-sidekiq-billing-canary", Namespace: namespace, Missing: true}).Assert(k8sClient, pm, timeout, poll))
					By("removing the system-sidekiq-low-canary PodMonitor",
						(&testutil.ExpectedResource{
							Name: "system-sidekiq-low-canary", Namespace: namespace, Missing: true}).Assert(k8sClient, pm, timeout, poll))
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

					rvs["deployment/system-app"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "system-app", namespace, timeout, poll)
					rvs["deployment/system-sidekiq-billing"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "system-sidekiq-billing", namespace, timeout, poll)
					rvs["deployment/system-sidekiq-default"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "system-sidekiq-default", namespace, timeout, poll)
					rvs["deployment/system-sidekiq-low"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "system-sidekiq-low", namespace, timeout, poll)

					patch := client.MergeFrom(system.DeepCopy())

					system.Spec.Config.Rails = &saasv1alpha1.SystemRailsSpec{
						Console: util.Pointer(true),
					}

					system.Spec.Twemproxy = &saasv1alpha1.TwemproxySpec{
						TwemproxyConfigRef: "system-twemproxyconfig",
						Options: &saasv1alpha1.TwemproxyOptions{
							LogLevel: util.Pointer[int32](2),
						},
					}

					system.Spec.App = &saasv1alpha1.SystemAppSpec{
						Canary: &saasv1alpha1.Canary{
							Replicas: util.Pointer[int32](2),
							Patches: []string{
								`[{"op":"add","path":"/twemproxy","value":{"twemproxyConfigRef":"system-canary-twemproxyconfig","options":{"logLevel":2}}}]`,
							},
						},
					}

					system.Spec.SidekiqBilling = &saasv1alpha1.SystemSidekiqSpec{
						Canary: &saasv1alpha1.Canary{
							Replicas: util.Pointer[int32](3),
							Patches: []string{
								`[{"op":"add","path":"/twemproxy","value":{"twemproxyConfigRef":"system-canary-twemproxyconfig","options":{"logLevel":3}}}]`,
							},
						},
					}

					system.Spec.SidekiqDefault = &saasv1alpha1.SystemSidekiqSpec{
						Replicas: util.Pointer[int32](2),
						Canary: &saasv1alpha1.Canary{
							Replicas: util.Pointer[int32](4),
							Patches: []string{
								`[{"op":"add","path":"/twemproxy/options","value":{"logLevel":4}}]`,
							},
						},
					}

					system.Spec.SidekiqLow = &saasv1alpha1.SystemSidekiqSpec{
						Replicas: util.Pointer[int32](2),
						Canary: &saasv1alpha1.Canary{
							Replicas: util.Pointer[int32](5),
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
					(&testutil.ExpectedWorkload{
						Name:        "system-app",
						Namespace:   namespace,
						Replicas:    2,
						PDB:         true,
						HPA:         true,
						PodMonitor:  true,
						LastVersion: rvs["deployment/system-app"],
					}).Assert(k8sClient, dep, timeout, poll))

				Expect(dep.Spec.Template.Spec.Volumes).To(HaveLen(2))
				Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-config"))
				Expect(dep.Spec.Template.Spec.Volumes[0].VolumeSource.Secret.SecretName).To(Equal("system-config"))
				Expect(dep.Spec.Template.Spec.Volumes[1].Name).To(Equal("twemproxy-config"))
				Expect(dep.Spec.Template.Spec.Volumes[1].VolumeSource.ConfigMap.LocalObjectReference.Name).To(Equal("system-twemproxyconfig"))
				Expect(dep.Spec.Template.Spec.Containers).To(HaveLen(2))
				Expect(dep.Spec.Template.Spec.Containers[1].Name).To(Equal("twemproxy"))
				Expect(dep.Spec.Template.Spec.Containers[1].VolumeMounts[0].Name).To(Equal("twemproxy-config"))

				By("adding a twemproxy sidecar to the system-app-canary workload",
					(&testutil.ExpectedWorkload{
						Name:       "system-app-canary",
						Replicas:   2,
						Namespace:  namespace,
						PodMonitor: true,
					}).Assert(k8sClient, dep, timeout, poll))

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
					(&testutil.ExpectedWorkload{
						Name:        "system-sidekiq-billing",
						Namespace:   namespace,
						Replicas:    2,
						PDB:         true,
						HPA:         true,
						PodMonitor:  true,
						LastVersion: rvs["deployment/system-sidekiq-billing"],
					}).Assert(k8sClient, dep, timeout, poll))

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
					(&testutil.ExpectedWorkload{
						Name:       "system-sidekiq-billing-canary",
						Replicas:   3,
						Namespace:  namespace,
						PodMonitor: true,
					}).Assert(k8sClient, dep, timeout, poll))

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
					(&testutil.ExpectedWorkload{
						Name:        "system-sidekiq-low",
						Namespace:   namespace,
						Replicas:    2,
						PDB:         true,
						HPA:         true,
						PodMonitor:  true,
						LastVersion: rvs["deployment/system-sidekiq-low"],
					}).Assert(k8sClient, dep, timeout, poll))

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
					(&testutil.ExpectedWorkload{
						Name:       "system-sidekiq-low-canary",
						Replicas:   5,
						Namespace:  namespace,
						PodMonitor: true,
					}).Assert(k8sClient, dep, timeout, poll))

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

				sts := &appsv1.StatefulSet{}
				By("adding a twemproxy sidecar to the system-console statefulset",
					(&testutil.ExpectedResource{Name: "system-console", Namespace: namespace}).
						Assert(k8sClient, sts, timeout, poll))

				Expect(sts.Spec.Template.Spec.Volumes).To(HaveLen(2))
				Expect(sts.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-config"))
				Expect(sts.Spec.Template.Spec.Volumes[0].VolumeSource.Secret.SecretName).To(Equal("system-config"))
				Expect(sts.Spec.Template.Spec.Volumes[1].Name).To(Equal("twemproxy-config"))
				Expect(sts.Spec.Template.Spec.Volumes[1].VolumeSource.ConfigMap.LocalObjectReference.Name).To(Equal("system-twemproxyconfig"))
				Expect(sts.Spec.Template.Spec.Containers).To(HaveLen(2))
				Expect(sts.Spec.Template.Spec.Containers[1].Name).To(Equal("twemproxy"))
				Expect(sts.Spec.Template.Spec.Containers[1].VolumeMounts[0].Name).To(Equal("twemproxy-config"))
				Expect(sts.Spec.Template.Spec.Containers[1].Env[3].Name).To(Equal("TWEMPROXY_LOG_LEVEL"))
				Expect(sts.Spec.Template.Spec.Containers[1].Env[3].Value).To(Equal("2"))

			})
		})

		When("updating system secret properties", func() {

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

					rvs["externalsecret/system-database"] = testutil.GetResourceVersion(
						k8sClient, &externalsecretsv1beta1.ExternalSecret{}, "system-database", namespace, timeout, poll)

					patch := client.MergeFrom(system.DeepCopy())

					system.Spec.Config.ExternalSecret.RefreshInterval = &metav1.Duration{Duration: 1 * time.Second}
					system.Spec.Config.ExternalSecret.SecretStoreRef = &saasv1alpha1.ExternalSecretSecretStoreReferenceSpec{
						Name: util.Pointer("other-store"),
						Kind: util.Pointer("SecretStore"),
					}
					system.Spec.Config.DatabaseDSN.FromVault.Path = "secret/data/updated-path"

					return k8sClient.Patch(context.Background(), system, patch)

				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("updates the system secret properties", func() {

				es := &externalsecretsv1beta1.ExternalSecret{}
				By("updating the system-database external secret",
					(&testutil.ExpectedResource{
						Name:        "system-database",
						Namespace:   namespace,
						LastVersion: rvs["externalsecret/system-database"],
					}).Assert(k8sClient, es, timeout, poll))

				Expect(es.Spec.RefreshInterval.ToUnstructured()).To(Equal("1s"))
				Expect(es.Spec.SecretStoreRef.Name).To(Equal("other-store"))
				Expect(es.Spec.SecretStoreRef.Kind).To(Equal("SecretStore"))

				for _, data := range es.Spec.Data {
					switch data.SecretKey {
					case "DATABASE_URL":
						Expect(data.RemoteRef.Key).To(Equal("updated-path"))
					}
				}
			})
		})

		When("updating a System resource tekton tasks", func() {

			// Resource Versions
			rvs := make(map[string]string)

			BeforeEach(func() {
				Eventually(func() error {
					err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						system,
					)
					Expect(err).ToNot(HaveOccurred())

					for _, tr := range []string{
						"system-db-migrate",
						"system-searchd-reindex",
						"system-backend-sync",
					} {
						rvs[fmt.Sprintf("task/%s", tr)] = testutil.GetResourceVersion(
							k8sClient, &pipelinev1beta1.Task{}, tr, namespace, timeout, poll)
						rvs[fmt.Sprintf("pipeline/%s", tr)] = testutil.GetResourceVersion(
							k8sClient, &pipelinev1beta1.Pipeline{}, tr, namespace, timeout, poll)
					}

					patch := client.MergeFrom(system.DeepCopy())

					system.Spec.Config.Rails = &saasv1alpha1.SystemRailsSpec{
						LogLevel: util.Pointer("debug"),
					}

					system.Spec.Tasks = []saasv1alpha1.SystemTektonTaskSpec{
						{
							Name:    util.Pointer("system-db-migrate"),
							Enabled: util.Pointer(false),
						},
						{
							Name: util.Pointer("system-searchd-reindex"),
							Config: &saasv1alpha1.SystemTektonTaskConfig{
								Image: &saasv1alpha1.ImageSpec{
									Name: util.Pointer("newImage"),
									Tag:  util.Pointer("newTag"),
								},
								Command: []string{"cmd"},
								Args:    []string{"arg1", "arg1"},
								Timeout: &metav1.Duration{Duration: 1 * time.Hour},
								ExtraEnv: []corev1.EnvVar{
									{Name: "test", Value: "test"},
									{Name: "THINKING_SPHINX_BATCH_SIZE", Value: "50"},
								},
							},
						},
						{
							Name:        util.Pointer("test-task"),
							Description: util.Pointer("Test task"),
							Config: &saasv1alpha1.SystemTektonTaskConfig{
								Command: []string{"cmd"},
								Args:    []string{"arg1", "arg1"},
								ExtraEnv: []corev1.EnvVar{
									{Name: "test", Value: "test"},
									{Name: "RAILS_LOG_LEVEL", Value: "debug"},
								},
							},
						},
					}
					return k8sClient.Patch(context.Background(), system, patch)
				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("updates the required tekton resources", func() {

				task := &pipelinev1beta1.Task{}
				pipeline := &pipelinev1beta1.Pipeline{}

				By("keeping the system-backend-sync task",
					(&testutil.ExpectedResource{
						Name: "system-backend-sync", Namespace: namespace,
					}).Assert(k8sClient, task, timeout, poll))

				Expect(task.Spec.Params[0].Default.StringVal).To(Equal("quay.io/3scale/porta"))
				Expect(task.Spec.Params[1].Default.StringVal).To(Equal("nightly"))
				Expect(task.Spec.Steps[0].Timeout).
					To(Equal(&metav1.Duration{Duration: 3 * time.Hour}))

				By("updating the system-searchd-reindex task",
					(&testutil.ExpectedResource{
						Name: "system-searchd-reindex", Namespace: namespace,
						LastVersion: rvs["task/system-searchd-reindex"],
					}).Assert(k8sClient, task, timeout, poll))

				Expect(task.Spec.Params[0].Default.StringVal).To(Equal("newImage"))
				Expect(task.Spec.Params[1].Default.StringVal).To(Equal("newTag"))

				Expect(task.Spec.Steps[0].Command[0]).To(Equal("cmd"))
				Expect(task.Spec.Steps[0].Args[0]).To(Equal("arg1"))
				Expect(task.Spec.Steps[0].Timeout).
					To(Equal(&metav1.Duration{Duration: 1 * time.Hour}))

				for _, env := range task.Spec.StepTemplate.Env {
					switch env.Name {
					case "test":
						Expect(env.Value).To(Equal("test"))
					case "RAILS_LOG_TO_STDOUT":
						Expect(env.Value).To(Equal("true"))
					case "RAILS_LOG_LEVEL":
						Expect(env.Value).To(Equal("debug"))
					case "THINKING_SPHINX_BATCH_SIZE":
						Expect(env.Value).To(Equal("50"))
					}
				}

				By("updating the system-searchd-reindex pipeline",
					(&testutil.ExpectedResource{
						Name: "system-searchd-reindex", Namespace: namespace,
						LastVersion: rvs["pipeline/system-searchd-reindex"],
					}).Assert(k8sClient, pipeline, timeout, poll))

				Expect(pipeline.Spec.Params[0].Default.StringVal).To(Equal("newImage"))
				Expect(pipeline.Spec.Params[1].Default.StringVal).To(Equal("newTag"))

				By("adding the new test-task",
					(&testutil.ExpectedResource{
						Name: "test-task", Namespace: namespace,
					}).Assert(k8sClient, task, timeout, poll))

				Expect(task.Spec.DisplayName).To(Equal("test-task"))
				Expect(task.Spec.Description).To(Equal("Test task"))
				Expect(task.Spec.Steps[0].Command[0]).To(Equal("cmd"))
				Expect(task.Spec.Steps[0].Args[0]).To(Equal("arg1"))

				for _, env := range task.Spec.StepTemplate.Env {
					switch env.Name {
					case "test":
						Expect(env.Value).To(Equal("test"))
					case "RAILS_LOG_LEVEL":
						Expect(env.Value).To(Equal("debug"))
					}
				}

				By("removing the system-db-migrate task",
					(&testutil.ExpectedResource{
						Name: "system-db-migrate", Namespace: namespace, Missing: true,
					}).Assert(k8sClient, task, timeout, poll))

				By("removing the system-db-migrate pipeline",
					(&testutil.ExpectedResource{
						Name: "system-db-migrate", Namespace: namespace, Missing: true,
					}).Assert(k8sClient, pipeline, timeout, poll))

			})
		})
	})
})
