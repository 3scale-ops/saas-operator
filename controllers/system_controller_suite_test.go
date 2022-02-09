package controllers

import (
	"context"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	grafanav1alpha1 "github.com/3scale/saas-operator/pkg/apis/grafana/v1alpha1"
	secretsmanagerv1alpha1 "github.com/3scale/saas-operator/pkg/apis/secrets-manager/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
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
		testNamespace := &v1.Namespace{
			TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Namespace"},
			ObjectMeta: metav1.ObjectMeta{Name: namespace},
		}

		err := k8sClient.Create(context.Background(), testNamespace)
		Expect(err).ToNot(HaveOccurred())

		n := &v1.Namespace{}
		Eventually(func() error {
			return k8sClient.Get(context.Background(), types.NamespacedName{Name: namespace}, n)
		}, timeout, poll).ShouldNot(HaveOccurred())

	})

	systemDefaultsConfig := saasv1alpha1.SystemConfig{
		DatabaseDSN:        saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
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
	}

	Context("System defaulted resource", func() {

		BeforeEach(func() {
			By("creating an System simple resource")
			system = &saasv1alpha1.System{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance",
					Namespace: namespace,
				},
				Spec: saasv1alpha1.SystemSpec{
					Config: systemDefaultsConfig,
				},
			}
			err := k8sClient.Create(context.Background(), system)
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, system)
			}, timeout, poll).ShouldNot(HaveOccurred())
		})

		It("creates the required default resources", func() {

			Eventually(func() bool {
				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, system)
				Expect(err).ToNot(HaveOccurred())
				return len(system.GetFinalizers()) > 0
			}, timeout, poll).Should(BeTrue())

			dep := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-app", Namespace: namespace},
					dep,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Expect(dep.Spec.Template.Spec.Volumes[0].Secret.SecretName).To(Equal("system-config"))
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-default", Namespace: namespace},
					dep,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Expect(dep.Spec.Template.Spec.Containers[0].Args).To(Equal(
				[]string{"sidekiq",
					"--queue", "critical", "--queue", "backend_sync",
					"--queue", "events", "--queue", "zync,40",
					"--queue", "priority,25", "--queue", "default,15",
					"--queue", "web_hooks,10", "--queue", "deletion,5",
				}))
			Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
			Expect(dep.Spec.Template.Spec.Volumes[1].Secret.SecretName).To(Equal("system-config"))
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-billing", Namespace: namespace},
					dep,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Expect(dep.Spec.Template.Spec.Containers[0].Args).To(Equal(
				[]string{"sidekiq", "--queue", "billing"},
			))
			Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
			Expect(dep.Spec.Template.Spec.Volumes[1].Secret.SecretName).To(Equal("system-config"))
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-low", Namespace: namespace},
					dep,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Expect(dep.Spec.Template.Spec.Containers[0].Args).To(Equal(
				[]string{"sidekiq", "--queue", "mailers", "--queue", "low"},
			))
			Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
			Expect(dep.Spec.Template.Spec.Volumes[1].Secret.SecretName).To(Equal("system-config"))

			ss := &appsv1.StatefulSet{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sphinx", Namespace: namespace},
					ss,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			svc := &corev1.Service{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-app", Namespace: namespace},
					svc,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Expect(svc.Spec.Selector["deployment"]).To(Equal("system-app"))

			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sphinx", Namespace: namespace},
					svc,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Expect(svc.Spec.Selector["deployment"]).To(Equal("system-sphinx"))

			pm := &monitoringv1.PodMonitor{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-app", Namespace: namespace},
					pm,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-default", Namespace: namespace},
					pm,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-billing", Namespace: namespace},
					pm,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-low", Namespace: namespace},
					pm,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-app", Namespace: namespace},
					hpa,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-default", Namespace: namespace},
					hpa,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-billing", Namespace: namespace},
					hpa,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-low", Namespace: namespace},
					hpa,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			pdb := &policyv1beta1.PodDisruptionBudget{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-app", Namespace: namespace},
					pdb,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-default", Namespace: namespace},
					pdb,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-billing", Namespace: namespace},
					pdb,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-low", Namespace: namespace},
					pdb,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			gd := &grafanav1alpha1.GrafanaDashboard{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system", Namespace: namespace},
					gd,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			sd := &secretsmanagerv1alpha1.SecretDefinition{}
			for _, name := range []string{
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
				Eventually(func() error {
					return k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: name, Namespace: namespace},
						sd,
					)
				}, timeout, poll).ShouldNot(HaveOccurred())
			}

		})

	})

	Context("System resource with customizations", func() {

		BeforeEach(func() {
			By("creating a System resource")
			system = &saasv1alpha1.System{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance",
					Namespace: namespace,
				},
				Spec: saasv1alpha1.SystemSpec{
					Config: systemDefaultsConfig,
				},
			}
			err := k8sClient.Create(context.Background(), system)
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, system)
			}, timeout, poll).ShouldNot(HaveOccurred())
		})

		It("creates the System defaulted resource", func() {
			Eventually(func() bool {
				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, system)
				Expect(err).ToNot(HaveOccurred())
				return len(system.GetFinalizers()) > 0
			}, timeout, poll).Should(BeTrue())
		})

		It("doesn't creates the non-default resources", func() {

			ss := &appsv1.StatefulSet{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-console", Namespace: namespace},
					ss,
				)
			}, timeout, poll).Should(HaveOccurred())

			dep := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-console", Namespace: namespace},
					dep,
				)
			}, timeout, poll).Should(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-app-canary", Namespace: namespace},
					dep,
				)
			}, timeout, poll).Should(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-default-canary", Namespace: namespace},
					dep,
				)
			}, timeout, poll).Should(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-billing-canary", Namespace: namespace},
					dep,
				)
			}, timeout, poll).Should(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-low-canary", Namespace: namespace},
					dep,
				)
			}, timeout, poll).Should(HaveOccurred())

		})

		BeforeEach(func() {
			By("Enabling the rails console")
			system := &saasv1alpha1.System{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "instance", Namespace: namespace},
					system,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			patch := client.MergeFrom(system.DeepCopy())
			system.Spec.Config.Rails = &saasv1alpha1.SystemRailsSpec{
				Console: pointer.Bool(true),
			}
			err := k8sClient.Patch(context.Background(), system, patch)
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, system)
			}, timeout, poll).ShouldNot(HaveOccurred())

		})

		It("creates only required console resources", func() {

			Eventually(func() bool {
				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, system)
				Expect(err).ToNot(HaveOccurred())
				return len(system.GetFinalizers()) > 0
			}, timeout, poll).Should(BeTrue())

			hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-console", Namespace: namespace},
					hpa,
				)
			}, timeout, poll).Should(HaveOccurred())

			pdb := &policyv1beta1.PodDisruptionBudget{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-console", Namespace: namespace},
					pdb,
				)
			}, timeout, poll).Should(HaveOccurred())

		})

		BeforeEach(func() {
			By("enabling a canary on all system components")
			system := &saasv1alpha1.System{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "instance", Namespace: namespace},
					system,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			patch := client.MergeFrom(system.DeepCopy())
			system.Spec.App = &saasv1alpha1.SystemAppSpec{
				Canary: &saasv1alpha1.Canary{
					Replicas: pointer.Int32Ptr(1),
				},
			}
			system.Spec.SidekiqDefault = &saasv1alpha1.SystemSidekiqSpec{
				Canary: &saasv1alpha1.Canary{
					Replicas: pointer.Int32Ptr(1),
				},
			}
			system.Spec.SidekiqBilling = &saasv1alpha1.SystemSidekiqSpec{
				Canary: &saasv1alpha1.Canary{
					Replicas: pointer.Int32Ptr(1),
				},
			}
			system.Spec.SidekiqLow = &saasv1alpha1.SystemSidekiqSpec{
				Canary: &saasv1alpha1.Canary{
					Replicas: pointer.Int32Ptr(1),
				},
			}
			err := k8sClient.Patch(context.Background(), system, patch)
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, system)
			}, timeout, poll).ShouldNot(HaveOccurred())
		})

		It("creates only required canary resources", func() {

			dep := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-app-canary", Namespace: namespace},
					dep,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Expect(dep.Spec.Template.Spec.Volumes[0].Secret.SecretName).To(Equal("system-config"))

			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-default-canary", Namespace: namespace},
					dep,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
			Expect(dep.Spec.Template.Spec.Volumes[1].Secret.SecretName).To(Equal("system-config"))

			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-low-canary", Namespace: namespace},
					dep,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
			Expect(dep.Spec.Template.Spec.Volumes[1].Secret.SecretName).To(Equal("system-config"))

			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-billing-canary", Namespace: namespace},
					dep,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("system-tmp"))
			Expect(dep.Spec.Template.Spec.Volumes[1].Secret.SecretName).To(Equal("system-config"))

			svc := &corev1.Service{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-app", Namespace: namespace},
					svc,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Expect(svc.Spec.Selector["deployment"]).To(Equal("system-app"))
			Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("system-app"))

			pdb := &policyv1beta1.PodDisruptionBudget{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-app-canary", Namespace: namespace},
					pdb,
				)
			}, timeout, poll).Should(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-default-canary", Namespace: namespace},
					pdb,
				)
			}, timeout, poll).Should(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-billing-canary", Namespace: namespace},
					pdb,
				)
			}, timeout, poll).Should(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-low-canary", Namespace: namespace},
					pdb,
				)
			}, timeout, poll).Should(HaveOccurred())

			hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{}

			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-canary-app", Namespace: namespace},
					hpa,
				)
			}, timeout, poll).Should(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-default-canary", Namespace: namespace},
					hpa,
				)
			}, timeout, poll).Should(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-billing-canary", Namespace: namespace},
					hpa,
				)
			}, timeout, poll).Should(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq-low-canary", Namespace: namespace},
					hpa,
				)
			}, timeout, poll).Should(HaveOccurred())

		})

		// BeforeEach(func() {
		// 	By("Update System simple resource with canary traffic enabled and replicas bump to 2")
		// 	system := &saasv1alpha1.System{}
		// 	Eventually(func() error {
		// 		return k8sClient.Get(
		// 			context.Background(),
		// 			types.NamespacedName{Name: "instance", Namespace: namespace},
		// 			system,
		// 		)
		// 	}, timeout, poll).ShouldNot(HaveOccurred())

		// 	patch := client.MergeFrom(system.DeepCopy())
		// 	system.Spec.App = &saasv1alpha1.SystemAppSpec{
		// 		Canary: &saasv1alpha1.Canary{
		// 			ImageName:   pointer.StringPtr("newImage"),
		// 			ImageTag:    pointer.StringPtr("newTag"),
		// 			Replicas:    pointer.Int32Ptr(2),
		// 			SendTraffic: true,
		// 		},
		// 	}
		// 	system.Spec.SidekiqDefault = &saasv1alpha1.SystemSidekiqSpec{
		// 		Canary: &saasv1alpha1.Canary{
		// 			ImageName: pointer.StringPtr("newImage"),
		// 			ImageTag:  pointer.StringPtr("newTag"),
		// 			Replicas:  pointer.Int32Ptr(2),
		// 		},
		// 	}
		// 	system.Spec.SidekiqBilling = &saasv1alpha1.SystemSidekiqSpec{
		// 		Canary: &saasv1alpha1.Canary{
		// 			ImageName: pointer.StringPtr("newImage"),
		// 			ImageTag:  pointer.StringPtr("newTag"),
		// 			Replicas:  pointer.Int32Ptr(2),
		// 		},
		// 	}
		// 	system.Spec.SidekiqLow = &saasv1alpha1.SystemSidekiqSpec{
		// 		Canary: &saasv1alpha1.Canary{
		// 			ImageName: pointer.StringPtr("newImage"),
		// 			ImageTag:  pointer.StringPtr("newTag"),
		// 			Replicas:  pointer.Int32Ptr(2),
		// 		},
		// 	}
		// 	err := k8sClient.Patch(context.Background(), system, patch)
		// 	Expect(err).ToNot(HaveOccurred())
		// 	Eventually(func() error {
		// 		return k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, system)
		// 	}, timeout, poll).ShouldNot(HaveOccurred())

		// })

		// It("scales up and enables the canary traffic", func() {

		// 	dep := &appsv1.Deployment{}
		// 	Eventually(func() error {
		// 		return k8sClient.Get(
		// 			context.Background(),
		// 			types.NamespacedName{Name: "system-app", Namespace: namespace},
		// 			dep,
		// 		)
		// 	}, timeout, poll).ShouldNot(HaveOccurred())
		// 	Expect(dep.Spec.Replicas).To(Equal(pointer.Int32Ptr(2)))
		// 	Expect(dep.Spec.Template.Spec.Containers[0].Image).To(Equal("newImage:newTag"))

		// 	Eventually(func() error {
		// 		return k8sClient.Get(
		// 			context.Background(),
		// 			types.NamespacedName{Name: "system-sidekiq-default-canary", Namespace: namespace},
		// 			dep,
		// 		)
		// 	}, timeout, poll).ShouldNot(HaveOccurred())
		// 	Expect(dep.Spec.Replicas).To(Equal(pointer.Int32Ptr(2)))
		// 	Expect(dep.Spec.Template.Spec.Containers[0].Image).To(Equal("newImage:newTag"))

		// 	Eventually(func() error {
		// 		return k8sClient.Get(
		// 			context.Background(),
		// 			types.NamespacedName{Name: "system-sidekiq-low-canary", Namespace: namespace},
		// 			dep,
		// 		)
		// 	}, timeout, poll).ShouldNot(HaveOccurred())
		// 	Expect(dep.Spec.Replicas).To(Equal(pointer.Int32Ptr(2)))
		// 	Expect(dep.Spec.Template.Spec.Containers[0].Image).To(Equal("newImage:newTag"))

		// 	Eventually(func() error {
		// 		return k8sClient.Get(
		// 			context.Background(),
		// 			types.NamespacedName{Name: "system-sidekiq-billing-canary", Namespace: namespace},
		// 			dep,
		// 		)
		// 	}, timeout, poll).ShouldNot(HaveOccurred())
		// 	Expect(dep.Spec.Replicas).To(Equal(pointer.Int32Ptr(2)))
		// 	Expect(dep.Spec.Template.Spec.Containers[0].Image).To(Equal("newImage:newTag"))

		// 	svc := &corev1.Service{}
		// 	Eventually(func() error {
		// 		return k8sClient.Get(
		// 			context.Background(),
		// 			types.NamespacedName{Name: "system-app", Namespace: namespace},
		// 			svc,
		// 		)
		// 	}, timeout, poll).ShouldNot(HaveOccurred())
		// 	Expect(svc.Spec.Selector["deployment"]).To(Equal(""))
		// 	Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("system-app"))

		// })

	})
})
