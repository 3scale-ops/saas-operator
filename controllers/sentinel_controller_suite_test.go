package controllers

import (
	"context"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	testutil "github.com/3scale/saas-operator/test/util"
	grafanav1alpha1 "github.com/grafana-operator/grafana-operator/v4/api/integreatly/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Sentinel controller", func() {
	var namespace string
	var sentinel *saasv1alpha1.Sentinel

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

	When("deploying a defaulted Sentinel instance", func() {

		BeforeEach(func() {

			By("creating a Sentinel simple resource", func() {
				sentinel = &saasv1alpha1.Sentinel{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "instance",
						Namespace: namespace,
					},
					Spec: saasv1alpha1.SentinelSpec{
						Config: &saasv1alpha1.SentinelConfig{
							MonitoredShards: map[string][]string{
								"shard01": {"redis://10.65.0.10:6379", "redis://10.65.0.20:6379", "redis://10.65.0.30:6379"},
								"shard02": {"redis://10.65.0.10:6379", "redis://10.65.0.20:6379", "redis://10.65.0.30:6379"},
							},
						},
					},
				}
				err := k8sClient.Create(context.Background(), sentinel)
				Expect(err).ToNot(HaveOccurred())
				Eventually(func() error {
					return k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, sentinel)
				}, timeout, poll).ShouldNot(HaveOccurred())
			})

		})

		It("creates the required Sentiel resources", func() {

			By("deploying the sentinel statefulset",
				(&testutil.ExpectedResource{Name: "redis-sentinel", Namespace: namespace}).
					Assert(k8sClient, &appsv1.StatefulSet{}, timeout, poll))

			svc := &corev1.Service{}
			By("deploying a Sentinel headless service",
				(&testutil.ExpectedResource{Name: "redis-sentinel-headless", Namespace: namespace}).
					Assert(k8sClient, svc, timeout, poll))

			Expect(svc.Spec.Selector["deployment"]).To(Equal("redis-sentinel"))

			By("deploying a Sentinel redis-0 service",
				(&testutil.ExpectedResource{Name: "redis-sentinel-0", Namespace: namespace}).
					Assert(k8sClient, svc, timeout, poll))

			By("deploying a Sentinel redis-1 service",
				(&testutil.ExpectedResource{Name: "redis-sentinel-1", Namespace: namespace}).
					Assert(k8sClient, svc, timeout, poll))

			By("deploying a Sentinel redis-2 service",
				(&testutil.ExpectedResource{Name: "redis-sentinel-2", Namespace: namespace}).
					Assert(k8sClient, svc, timeout, poll))

			By("deploying a Sentinel gen-config configmap",
				(&testutil.ExpectedResource{Name: "redis-sentinel-gen-config", Namespace: namespace}).
					Assert(k8sClient, &corev1.ConfigMap{}, timeout, poll))

			By("deploying the Sentinel grafana dashboard",
				(&testutil.ExpectedResource{Name: "redis-sentinel", Namespace: namespace}).
					Assert(k8sClient, &grafanav1alpha1.GrafanaDashboard{}, timeout, poll))

		})

	})

})
