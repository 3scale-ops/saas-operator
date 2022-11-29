package test

import (
	"context"
	"reflect"

	"github.com/3scale/saas-operator/pkg/reconcilers/workloads/test/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
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
				Spec: v1alpha1.TestSpec{
					TrafficSelector: map[string]string{"traffic": "yes"},
					Alice: v1alpha1.Workload{
						Name:     "alice",
						Traffic:  true,
						Selector: map[string]string{"deployment": "alice"},
						Labels:   map[string]string{"alice-lkey1": "alice-lvalue1"},
					},
					Bob: v1alpha1.Workload{
						Name:     "bob",
						Traffic:  true,
						Selector: map[string]string{"deployment": "bob"},
						Labels:   map[string]string{"bob-lkey1": "bob-lvalue1"},
					},
				},
			}
			err := k8sClient.Create(context.Background(), instance)
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, instance)
			}, timeout, poll).ShouldNot(HaveOccurred())
		})

		It("creates the required resources", func() {

			// alice Deployment
			{
				alice := &appsv1.Deployment{}
				Eventually(func() error {
					return k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "alice", Namespace: namespace},
						alice,
					)
				}, timeout, poll).ShouldNot(HaveOccurred())
				Expect(alice.GetLabels()).To(Equal(map[string]string{"alice-lkey1": "alice-lvalue1"}))
				Expect(alice.Spec.Selector.MatchLabels).To(Equal(map[string]string{"deployment": "alice"}))
				Expect(alice.Spec.Template.GetLabels()).
					To(Equal(map[string]string{"alice-lkey1": "alice-lvalue1", "deployment": "alice", "traffic": "yes", "orig-key": "orig-value"}))

				hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{}
				Eventually(func() error {
					return k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "alice", Namespace: namespace},
						hpa,
					)
				}, timeout, poll).ShouldNot(HaveOccurred())
				Expect(hpa.GetLabels()).To(Equal(map[string]string{"alice-lkey1": "alice-lvalue1"}))

				pdb := &policyv1.PodDisruptionBudget{}
				Eventually(func() error {
					return k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "alice", Namespace: namespace},
						pdb,
					)
				}, timeout, poll).ShouldNot(HaveOccurred())
				Expect(pdb.GetLabels()).To(Equal(map[string]string{"alice-lkey1": "alice-lvalue1"}))
			}

			// bob Deployment
			{
				bob := &appsv1.Deployment{}
				Eventually(func() error {
					return k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "bob", Namespace: namespace},
						bob,
					)
				}, timeout, poll).ShouldNot(HaveOccurred())
				Expect(bob.GetLabels()).To(Equal(map[string]string{"bob-lkey1": "bob-lvalue1"}))
				Expect(bob.Spec.Selector.MatchLabels).To(Equal(map[string]string{"deployment": "bob"}))
				Expect(bob.Spec.Template.GetLabels()).
					To(Equal(map[string]string{"bob-lkey1": "bob-lvalue1", "deployment": "bob", "traffic": "yes", "orig-key": "orig-value"}))

				hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{}
				Eventually(func() error {
					return k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "bob", Namespace: namespace},
						hpa,
					)
				}, timeout, poll).ShouldNot(HaveOccurred())
				Expect(hpa.GetLabels()).To(Equal(map[string]string{"bob-lkey1": "bob-lvalue1"}))

				pdb := &policyv1.PodDisruptionBudget{}
				Eventually(func() error {
					return k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "bob", Namespace: namespace},
						pdb,
					)
				}, timeout, poll).ShouldNot(HaveOccurred())
				Expect(pdb.GetLabels()).To(Equal(map[string]string{"bob-lkey1": "bob-lvalue1"}))
			}

			// Service
			svc := &corev1.Service{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "service", Namespace: namespace},
					svc,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			Expect(svc.Spec.Selector).To(Equal(instance.Spec.TrafficSelector))

		})

		It("Triggers Deployment rollouts on Secret contents change", func() {

			alice := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "alice", Namespace: namespace},
					alice,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			// Annotations should be empty when Secret does not exists
			value, ok := alice.Spec.Template.ObjectMeta.Annotations["saas.3scale.net/secret.secret-hash"]
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(""))

			bob := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "bob", Namespace: namespace},
					bob,
				)
			}, timeout, poll).ShouldNot(HaveOccurred())
			// Annotations should be empty when Secret does not exists
			value, ok = bob.Spec.Template.ObjectMeta.Annotations["saas.3scale.net/secret.secret-hash"]
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
					types.NamespacedName{Name: "alice", Namespace: namespace},
					alice,
				)
				Expect(err).ToNot(HaveOccurred())
				value, ok := alice.Spec.Template.ObjectMeta.Annotations["saas.3scale.net/secret.secret-hash"]
				Expect(ok).To(BeTrue())
				// Value of the annotation should be the hash of the Secret contents
				return value == util.Hash(secret.Data)
			}, timeout, poll).ShouldNot(BeTrue())

			Eventually(func() bool {
				err := k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "bob", Namespace: namespace},
					bob,
				)
				Expect(err).ToNot(HaveOccurred())
				value, ok := bob.Spec.Template.ObjectMeta.Annotations["saas.3scale.net/secret.secret-hash"]
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
					types.NamespacedName{Name: "alice", Namespace: namespace},
					alice,
				)
				Expect(err).ToNot(HaveOccurred())
				value, ok := alice.Spec.Template.ObjectMeta.Annotations["saas.3scale.net/secret.secret-hash"]
				Expect(ok).To(BeTrue())
				// Value of the annotation should be the hash of the Secret new contents
				return value == util.Hash(secret.Data)
			}, timeout, poll).ShouldNot(BeTrue())

			Eventually(func() bool {
				err := k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "bob", Namespace: namespace},
					bob,
				)
				Expect(err).ToNot(HaveOccurred())
				value, ok := bob.Spec.Template.ObjectMeta.Annotations["saas.3scale.net/secret.secret-hash"]
				Expect(ok).To(BeTrue())
				// Value of the annotation should be the hash of the Secret new contents
				return value == util.Hash(secret.Data)
			}, timeout, poll).ShouldNot(BeTrue())
		})

		It("Modifies the Service selector when traffic configuration is changed", func() {
			Eventually(func() error {
				if err := k8sClient.Get(context.Background(),
					types.NamespacedName{Name: "instance", Namespace: namespace}, instance); err != nil {
					return err
				}
				instance.Spec.Alice.Traffic = false
				return k8sClient.Update(context.Background(), instance)
			}, timeout, poll).ShouldNot(HaveOccurred())

			svc := &corev1.Service{}
			Eventually(func() bool {
				err := k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "service", Namespace: namespace},
					svc)
				if err != nil {
					return false
				}
				return reflect.DeepEqual(svc.Spec.Selector, map[string]string{"deployment": "bob", "traffic": "yes"})
			}, timeout, poll).Should(BeTrue())

			// store the nodePort to check later that it does not change
			// in the reconciles
			nodePort := svc.Spec.Ports[0].NodePort

			Eventually(func() error {
				if err := k8sClient.Get(context.Background(),
					types.NamespacedName{Name: "instance", Namespace: namespace}, instance); err != nil {
					return err
				}
				instance.Spec.Alice.Traffic = true
				return k8sClient.Update(context.Background(), instance)
			}, timeout, poll).ShouldNot(HaveOccurred())

			Eventually(func() bool {
				err := k8sClient.Get(
					context.Background(),
					types.NamespacedName{Name: "service", Namespace: namespace},
					svc)
				if err != nil {
					return false
				}
				return reflect.DeepEqual(svc.Spec.Selector, map[string]string{"traffic": "yes"})
			}, timeout, poll).Should(BeTrue())

			Expect(svc.Spec.Ports[0].NodePort).To(Equal(nodePort))
		})
	})
})
