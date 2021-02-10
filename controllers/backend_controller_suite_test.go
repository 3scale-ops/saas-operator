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

var _ = Describe("Backend controller", func() {
	var namespace string
	var backend *saasv1alpha1.Backend

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

	Context("All defaults Backend resource", func() {

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
		})

		It("creates the required resources", func() {

			Eventually(func() bool {
				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, backend)
				Expect(err).ToNot(HaveOccurred())
				if len(backend.GetFinalizers()) > 0 {
					return true
				}
				return false
			}, timeout, poll).Should(BeTrue())

			dep := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "backend-listener", Namespace: namespace},
					dep,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "backend-worker", Namespace: namespace},
					dep,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "backend-cron", Namespace: namespace},
					dep,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			svc := &corev1.Service{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "backend-listener", Namespace: namespace},
					svc,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "backend-listener-internal", Namespace: namespace},
					svc,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			pm := &monitoringv1.PodMonitor{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "backend-listener", Namespace: namespace},
					pm,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "backend-worker", Namespace: namespace},
					pm,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "backend-listener", Namespace: namespace},
					hpa,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "backend-worker", Namespace: namespace},
					hpa,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			pdb := &policyv1beta1.PodDisruptionBudget{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "backend-listener", Namespace: namespace},
					pdb,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "backend-worker", Namespace: namespace},
					pdb,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			sd := &secretsmanagerv1alpha1.SecretDefinition{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "backend-internal-api", Namespace: namespace},
					sd,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "backend-system-events-hook", Namespace: namespace},
					sd,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			gd := &grafanav1alpha1.GrafanaDashboard{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "backend", Namespace: namespace},
					gd,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
		})

		It("updates the service with non default vaule", func() {
			svc := &corev1.Service{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "backend-listener", Namespace: namespace},
					svc,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			backend := &saasv1alpha1.Backend{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "instance", Namespace: namespace},
					backend,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			patch := client.MergeFrom(backend.DeepCopy())
			backend.Spec.Listener.LoadBalancer = &saasv1alpha1.NLBLoadBalancerSpec{CrossZoneLoadBalancingEnabled: pointer.BoolPtr(false)}
			err := k8sClient.Patch(context.Background(), backend, patch)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err := k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "backend-listener", Namespace: namespace},
					svc,
				)
				Expect(err).ToNot(HaveOccurred())
				if svc.GetAnnotations()["service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled"] == "false" {
					return true
				}
				return false
			}, timeout, poll).Should(BeTrue())
		})
	})
})
