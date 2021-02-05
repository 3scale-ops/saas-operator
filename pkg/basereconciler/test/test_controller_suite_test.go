package test

import (
	"context"

	secretsmanagerv1alpha1 "github.com/3scale/saas-operator/pkg/apis/secrets-manager/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	"github.com/3scale/saas-operator/pkg/basereconciler/test/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Test controller", func() {
	var namespace string
	var instance *v1alpha1.Test

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

	Context("Creates resources", func() {

		BeforeEach(func() {
			By("creating a Test simple resource")
			instance = &v1alpha1.Test{
				ObjectMeta: metav1.ObjectMeta{Name: "instance", Namespace: namespace},
				Spec:       v1alpha1.TestSpec{},
			}
			err := k8sClient.Create(context.Background(), instance)
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, instance)
			}, timeout, poll).ShouldNot(HaveOccurred())
		})

		It("creates the required resources", func() {

			Eventually(func() bool {
				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, instance)
				Expect(err).ToNot(HaveOccurred())
				if len(instance.GetFinalizers()) > 0 {
					return true
				}
				return false
			}, timeout, poll).Should(BeTrue())

			dep := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "deployment", Namespace: namespace},
					dep,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			svc := &corev1.Service{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "service", Namespace: namespace},
					svc,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			sd := &secretsmanagerv1alpha1.SecretDefinition{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "secret", Namespace: namespace},
					sd,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			// pm := &monitoringv1.PodMonitor{}
			// Eventually(func() error {
			// 	return k8sClient.Get(
			// 		context.Background(),
			// 		types.NamespacedName{Name: "autossl", Namespace: namespace},
			// 		pm,
			// 	)
			// }, timeout, poll).ShouldNot(HaveOccurred())

			// hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{}
			// Eventually(func() error {
			// 	return k8sClient.Get(
			// 		context.Background(),
			// 		types.NamespacedName{Name: "autossl", Namespace: namespace},
			// 		hpa,
			// 	)
			// }, timeout, poll).ShouldNot(HaveOccurred())

			// pdb := &policyv1beta1.PodDisruptionBudget{}
			// Eventually(func() error {
			// 	return k8sClient.Get(
			// 		context.Background(),
			// 		types.NamespacedName{Name: "autossl", Namespace: namespace},
			// 		pdb,
			// 	)
			// }, timeout, poll).ShouldNot(HaveOccurred())

			// gd := &grafanav1alpha1.GrafanaDashboard{}
			// Eventually(func() error {
			// 	return k8sClient.Get(
			// 		context.Background(),
			// 		types.NamespacedName{Name: "autossl", Namespace: namespace},
			// 		gd,
			// 	)
			// }, timeout, poll).ShouldNot(HaveOccurred())
		})

		It("Triggers a Deployment rollout on Secret contents change", func() {

			dep := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "deployment", Namespace: namespace},
					dep,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			// Annotations should be empty when Secret does not exists
			value, ok := dep.Spec.Template.ObjectMeta.Annotations["saas.3scale.net/secret.secret-hash"]
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(""))

			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "secret", Namespace: namespace},
				Type:       corev1.SecretTypeOpaque,
				Data:       map[string][]byte{"KEY": []byte("value")},
			}
			err := k8sClient.Create(context.Background(), secret)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err := k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "deployment", Namespace: namespace},
					dep,
				)
				Expect(err).ToNot(HaveOccurred())
				value, ok := dep.Spec.Template.ObjectMeta.Annotations["saas.3scale.net/secret.secret-hash"]
				Expect(ok).To(BeTrue())
				// Value of the annotation should be the hash of the Secret contents
				if value == basereconciler.Hash(secret.Data) {
					return true
				}
				return false
			}, timeout, poll).ShouldNot(BeTrue())

			patch := client.MergeFrom(secret.DeepCopy())
			secret.Data = map[string][]byte{"KEY": []byte("new-value")}
			err = k8sClient.Patch(context.Background(), secret, patch)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err := k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "deployment", Namespace: namespace},
					dep,
				)
				Expect(err).ToNot(HaveOccurred())
				value, ok := dep.Spec.Template.ObjectMeta.Annotations["saas.3scale.net/secret.secret-hash"]
				Expect(ok).To(BeTrue())
				// Value of the annotation should be the hash of the Secret new contents
				if value == basereconciler.Hash(secret.Data) {
					return true
				}
				return false
			}, timeout, poll).ShouldNot(BeTrue())
		})
	})

})
