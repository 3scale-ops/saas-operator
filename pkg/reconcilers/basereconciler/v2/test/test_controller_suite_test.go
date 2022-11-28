package test

import (
	"context"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	externalsecretsv1beta1 "github.com/3scale/saas-operator/pkg/apis/externalsecrets/v1beta1"
	"github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2/test/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Test controller", func() {
	var namespace string
	var instance *v1alpha1.Test

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

			es := &externalsecretsv1beta1.ExternalSecret{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "secret", Namespace: namespace},
					es,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			pdb := &policyv1.PodDisruptionBudget{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "pdb", Namespace: namespace},
					pdb,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

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
				return value == util.Hash(secret.Data)
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
				return value == util.Hash(secret.Data)
			}, timeout, poll).ShouldNot(BeTrue())
		})

		It("Deletes specific resources when disabled", func() {
			// Wait for resources to be created
			pdb := &policyv1.PodDisruptionBudget{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "pdb", Namespace: namespace},
					pdb,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "hpa", Namespace: namespace},
					hpa,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			// disable pdb and hpa
			instance = &v1alpha1.Test{}
			Eventually(func() error {
				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, instance)
				if err != nil {
					return err
				}
				instance.Spec.PDB = pointer.Bool(false)
				instance.Spec.HPA = pointer.Bool(false)
				err = k8sClient.Update(context.Background(), instance)
				return err

			}, timeout, poll).ShouldNot(HaveOccurred())

			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "pdb", Namespace: namespace},
					pdb,
				)
			}, timeout, poll).Should(HaveOccurred())

			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "hpa", Namespace: namespace},
					hpa,
				)
			}, timeout, poll).Should(HaveOccurred())

		})

		It("Deletes all owned resources when custom resource is deleted", func() {
			// Wait for all resources to be created

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

			es := &externalsecretsv1beta1.ExternalSecret{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "secret", Namespace: namespace},
					es,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			// Delete the custom resource
			err := k8sClient.Delete(context.Background(), instance)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, instance)
				if err != nil && errors.IsNotFound(err) {
					return true
				}
				return false
			}, timeout, poll).Should(BeTrue())

		})

		It("updates service annotations", func() {
			svc := &corev1.Service{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "service", Namespace: namespace},
					svc,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "instance", Namespace: namespace},
					instance,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			patch := client.MergeFrom(instance.DeepCopy())
			instance.Spec.ServiceAnnotations = map[string]string{"key": "value"}
			err := k8sClient.Patch(context.Background(), instance, patch)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err := k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "service", Namespace: namespace},
					svc,
				)
				Expect(err).ToNot(HaveOccurred())
				return svc.GetAnnotations()["key"] == "value"
			}, timeout, poll).Should(BeTrue())
		})
	})

	Context("Marin3r enabled Deployments", func() {

		BeforeEach(func() {
			By("creating a marin3r enabled Test resource")
			instance = &v1alpha1.Test{
				ObjectMeta: metav1.ObjectMeta{Name: "instance", Namespace: namespace},
				Spec: v1alpha1.TestSpec{
					Marin3r: &saasv1alpha1.Marin3rSidecarSpec{
						Ports: []saasv1alpha1.SidecarPort{
							{
								Name: "test",
								Port: 9999,
							},
						},
						Resources: &saasv1alpha1.ResourceRequirementsSpec{
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("200m"),
								corev1.ResourceMemory: resource.MustParse("200Mi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("100m"),
								corev1.ResourceMemory: resource.MustParse("100Mi"),
							},
						},
						ExtraPodAnnotations: map[string]string{
							"extra-key": "extra-value",
						},
					},
				},
			}
			err := k8sClient.Create(context.Background(), instance)
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, instance)
			}, timeout, poll).ShouldNot(HaveOccurred())
		})

		It("creates the required deployment with proper labels and annotations", func() {

			dep := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "deployment", Namespace: namespace},
					dep,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())

			Expect(dep.Spec.Template.ObjectMeta.Annotations["marin3r.3scale.net/resources.limits.cpu"]).To(Equal("200m"))
			Expect(dep.Spec.Template.ObjectMeta.Annotations["marin3r.3scale.net/resources.limits.memory"]).To(Equal("200Mi"))
			Expect(dep.Spec.Template.ObjectMeta.Annotations["marin3r.3scale.net/resources.requests.cpu"]).To(Equal("100m"))
			Expect(dep.Spec.Template.ObjectMeta.Annotations["marin3r.3scale.net/resources.requests.memory"]).To(Equal("100Mi"))
			Expect(dep.Spec.Template.ObjectMeta.Annotations["marin3r.3scale.net/ports"]).To(Equal("test:9999"))
			Expect(dep.Spec.Template.ObjectMeta.Annotations["extra-key"]).To(Equal("extra-value"))
			Expect(dep.Spec.Template.ObjectMeta.Labels["marin3r.3scale.net/status"]).To(Equal("enabled"))
		})
	})

})
