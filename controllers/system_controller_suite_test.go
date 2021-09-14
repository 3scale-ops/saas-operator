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

	Context("All defaults System resource", func() {

		BeforeEach(func() {
			By("creating an System simple resource")
			system = &saasv1alpha1.System{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance",
					Namespace: namespace,
				},
				Spec: saasv1alpha1.SystemSpec{
					Config: saasv1alpha1.SystemConfig{
						ConfigFiles: &saasv1alpha1.ConfigFilesSpec{
							VaultPath: "some-path",
							Files:     []string{"some-file"},
						},
						Seed: saasv1alpha1.SystemSeedSpec{
							MasterAccessToken: saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
							MasterDomain:      "value",
							MasterUser:        saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
							MasterPassword:    saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
							AdminAccessToken:  saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
							AdminUser:         saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
							AdminPassword:     saasv1alpha1.SecretReference{Override: pointer.StringPtr("override")},
							AdminEmail:        "value",
							TenantName:        "value",
						},
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
					},
				},
			}
			err := k8sClient.Create(context.Background(), system)
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, system)
			}, timeout, poll).ShouldNot(HaveOccurred())
		})

		It("creates the required resources", func() {

			Eventually(func() bool {
				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, system)
				Expect(err).ToNot(HaveOccurred())
				if len(system.GetFinalizers()) > 0 {
					return true
				}
				return false
			}, timeout, poll).Should(BeTrue())

			dep := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-app", Namespace: namespace},
					dep,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sidekiq", Namespace: namespace},
					dep,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

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
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "system-sphinx", Namespace: namespace},
					svc,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

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
					types.NamespacedName{Name: "system-sidekiq", Namespace: namespace},
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
					types.NamespacedName{Name: "system-sidekiq", Namespace: namespace},
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
					types.NamespacedName{Name: "system-sidekiq", Namespace: namespace},
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
				"system-config",
				"system-database",
				"system-seed",
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
})
