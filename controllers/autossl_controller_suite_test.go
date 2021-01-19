package controllers

import (
	"context"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	grafanav1alpha1 "github.com/3scale/saas-operator/pkg/apis/grafana/v1alpha1"
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

var _ = Describe("AutoSSL controller", func() {
	var namespace string
	var autossl *saasv1alpha1.AutoSSL

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

	Context("All defaults AutoSSL resource", func() {

		BeforeEach(func() {
			By("creating an AutoSSL simple resource")
			autossl = &saasv1alpha1.AutoSSL{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance",
					Namespace: namespace,
				},
				Spec: saasv1alpha1.AutoSSLSpec{
					Config: saasv1alpha1.AutoSSLConfig{
						ContactEmail:         "test@example.com",
						ProxyEndpoint:        "example.com",
						VerificationEndpoint: "example.com/verification",
						RedisHost:            "redis.example.com",
					},
					Endpoint: saasv1alpha1.Endpoint{
						DNS: []string{"autossl.example.com"},
					},
				},
			}
			err := k8sClient.Create(context.Background(), autossl)
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, autossl)
			}, timeout, poll).ShouldNot(HaveOccurred())
		})

		It("creates the required resources", func() {

			Eventually(func() bool {
				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, autossl)
				Expect(err).ToNot(HaveOccurred())
				if len(autossl.GetFinalizers()) > 0 {
					return true
				}
				return false
			}, timeout, poll).Should(BeTrue())

			dep := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "autossl", Namespace: namespace},
					dep,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			svc := &corev1.Service{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "autossl", Namespace: namespace},
					svc,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			pm := &monitoringv1.PodMonitor{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "autossl", Namespace: namespace},
					pm,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "autossl", Namespace: namespace},
					hpa,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			pdb := &policyv1beta1.PodDisruptionBudget{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "autossl", Namespace: namespace},
					pdb,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			gd := &grafanav1alpha1.GrafanaDashboard{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "autossl", Namespace: namespace},
					gd,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
		})
	})

	Context("AutoSSL resource with deactivated features", func() {

		BeforeEach(func() {
			By("creating an AutoSSL simple resource")
			autossl = &saasv1alpha1.AutoSSL{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance",
					Namespace: namespace,
				},
				Spec: saasv1alpha1.AutoSSLSpec{
					Config: saasv1alpha1.AutoSSLConfig{
						ContactEmail:         "test@example.com",
						ProxyEndpoint:        "example.com",
						VerificationEndpoint: "example.com/verification",
						RedisHost:            "redis.example.com",
					},
					Endpoint: saasv1alpha1.Endpoint{
						DNS: []string{"autossl.example.com"},
					},
					PDB:              &saasv1alpha1.PodDisruptionBudgetSpec{},
					HPA:              &saasv1alpha1.HorizontalPodAutoscalerSpec{},
					Replicas:         pointer.Int32Ptr(1),
					Resources:        &saasv1alpha1.ResourceRequirementsSpec{},
					GrafanaDashboard: &saasv1alpha1.GrafanaDashboardSpec{},
					LivenessProbe:    &saasv1alpha1.HTTPProbeSpec{},
					ReadinessProbe:   &saasv1alpha1.HTTPProbeSpec{},
				},
			}
			err := k8sClient.Create(context.Background(), autossl)
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, autossl)
			}, timeout, poll).ShouldNot(HaveOccurred())
		})

		It("does not create deactivated resources/blocks", func() {

			dep := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "autossl", Namespace: namespace},
					dep,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			Expect(dep.Spec.Template.Spec.Containers[0].Resources).To(Equal(corev1.ResourceRequirements{}))
			Expect(dep.Spec.Template.Spec.Containers[0].LivenessProbe).To(BeNil())
			Expect(dep.Spec.Template.Spec.Containers[0].ReadinessProbe).To(BeNil())
			Expect(dep.Spec.Replicas).To(Equal(pointer.Int32Ptr(1)))

			hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "autossl", Namespace: namespace},
					hpa,
				)
			}, timeout, poll).Should(HaveOccurred())

			pdb := &policyv1beta1.PodDisruptionBudget{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "autossl", Namespace: namespace},
					pdb,
				)
			}, timeout, poll).Should(HaveOccurred())

			gd := &grafanav1alpha1.GrafanaDashboard{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "autossl", Namespace: namespace},
					gd,
				)
			}, timeout, poll).Should(HaveOccurred())
		})
	})
})
