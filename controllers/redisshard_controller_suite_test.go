package controllers

import (
	"context"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	testutil "github.com/3scale/saas-operator/test/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("TwemproxyConfig controller", func() {
	var namespace string
	var redisshard *saasv1alpha1.RedisShard

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

	When("deploying a defaulted RedisShard instance", func() {

		BeforeEach(func() {

			By("creating a RedisShard simple resource", func() {
				redisshard = &saasv1alpha1.RedisShard{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "shard01",
						Namespace: namespace,
					},
					Spec: saasv1alpha1.RedisShardSpec{},
				}
				err := k8sClient.Create(context.Background(), redisshard)
				Expect(err).ToNot(HaveOccurred())
				Eventually(func() error {
					return k8sClient.Get(context.Background(), types.NamespacedName{Name: "shard01", Namespace: namespace}, redisshard)
				}, timeout, poll).ShouldNot(HaveOccurred())
			})

		})

		It("creates the required RedisShard resources", func() {

			By("deploying a RedisShard statefulset",
				(&testutil.ExpectedResource{Name: "redis-shard-shard01", Namespace: namespace}).
					Assert(k8sClient, &appsv1.StatefulSet{}, timeout, poll))

			By("deploying a RedisShard service",
				(&testutil.ExpectedResource{Name: "redis-shard-shard01", Namespace: namespace}).
					Assert(k8sClient, &corev1.Service{}, timeout, poll))

			By("deploying a RedisShard config configmap",
				(&testutil.ExpectedResource{Name: "redis-config-shard01", Namespace: namespace}).
					Assert(k8sClient, &corev1.ConfigMap{}, timeout, poll))

		})

	})

})
