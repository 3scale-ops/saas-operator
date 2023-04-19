package controllers

import (
	"context"

	marin3rv1alpha1 "github.com/3scale-ops/marin3r/apis/marin3r/v1alpha1"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	testutil "github.com/3scale/saas-operator/test/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("EchoAPI controller", func() {
	var namespace string
	var echoapi *saasv1alpha1.EchoAPI

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

	When("deploying a defaulted EchoAPI instance", func() {

		BeforeEach(func() {
			By("creating an EchoAPI simple resource")
			echoapi = &saasv1alpha1.EchoAPI{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance",
					Namespace: namespace,
				},
				Spec: saasv1alpha1.EchoAPISpec{
					Endpoint: saasv1alpha1.Endpoint{
						DNS: []string{"echo-api.example.com"},
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
								},
							},
						},
					},
				},
			}
			err := k8sClient.Create(context.Background(), echoapi)
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, echoapi)
			}, timeout, poll).ShouldNot(HaveOccurred())
		})

		It("creates the required EchoAPI resources", func() {

			dep := &appsv1.Deployment{}
			By("deploying an echo-api workload",
				(&testutil.ExpectedWorkload{
					Name:          "echo-api",
					Namespace:     namespace,
					Replicas:      2,
					ContainerName: "echo-api",
					PDB:           true,
					HPA:           true,
					PodMonitor:    true,
				}).Assert(k8sClient, dep, timeout, poll))

			Expect(dep.Spec.Template.Spec.Volumes).To(HaveLen(0))

			svc := &corev1.Service{}
			By("deploying an echo-api service",
				(&testutil.ExpectedResource{Name: "echo-api", Namespace: namespace}).
					Assert(k8sClient, svc, timeout, poll))

			Expect(svc.Spec.Selector["deployment"]).To(Equal("echo-api"))
			Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("echo-api"))

			hpa := &autoscalingv2.HorizontalPodAutoscaler{}
			By("updates echo-api hpa behaviour",
				(&testutil.ExpectedResource{Name: "echo-api", Namespace: namespace}).
					Assert(k8sClient, hpa, timeout, poll))
			Expect(hpa.Spec.Behavior.ScaleUp.Policies).To(Not(BeEmpty()))
			Expect(hpa.Spec.Behavior.ScaleUp.Policies[0].Type).To(Equal(autoscalingv2.PodsScalingPolicy))
			Expect(hpa.Spec.Behavior.ScaleUp.Policies[0].PeriodSeconds).To(Equal(int32(60)))

		})

		It("doesn't create the non-default resources", func() {

			ec := &marin3rv1alpha1.EnvoyConfig{}
			By("ensuring an echo-api envoyconfig is not created",
				(&testutil.ExpectedResource{Name: "echo-api", Namespace: namespace, Missing: true}).
					Assert(k8sClient, ec, timeout, poll))

		})

		When("updating a EchoAPI resource with customizations", func() {

			// Resource Versions
			rvs := make(map[string]string)

			BeforeEach(func() {
				Eventually(func() error {

					echoapi := &saasv1alpha1.EchoAPI{}
					if err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						echoapi,
					); err != nil {
						return err
					}

					rvs["deployment/echoapi"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "echo-api", namespace, timeout, poll)

					patch := client.MergeFrom(echoapi.DeepCopy())
					echoapi.Spec.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{
						MinReplicas: pointer.Int32(3),
					}
					echoapi.Spec.LivenessProbe = &saasv1alpha1.ProbeSpec{}
					echoapi.Spec.ReadinessProbe = &saasv1alpha1.ProbeSpec{}

					echoapi.Spec.Marin3r = &saasv1alpha1.Marin3rSidecarSpec{
						NodeID: pointer.String("echo-api"),
						EnvoyDynamicConfig: saasv1alpha1.MapOfEnvoyDynamicConfig{
							"http": {
								GeneratorVersion: pointer.String("v1"),
								ListenerHttp: &saasv1alpha1.ListenerHttp{
									Port:            8080,
									RouteConfigName: "route",
								},
							}},
					}

					return k8sClient.Patch(context.Background(), echoapi, patch)

				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("updates the echo-api resources", func() {

				dep := &appsv1.Deployment{}
				By("updating the EchoAPI workload",
					(&testutil.ExpectedWorkload{
						Name:          "echo-api",
						Namespace:     namespace,
						Replicas:      3,
						ContainerName: "echo-api",
						PDB:           true,
						HPA:           true,
						PodMonitor:    true,
						EnvoyConfig:   true,
						LastVersion:   rvs["deployment/echoapi"],
					}).Assert(k8sClient, dep, timeout, poll))

				Expect(dep.Spec.Template.Spec.Containers[0].LivenessProbe).To(BeNil())
				Expect(dep.Spec.Template.Spec.Containers[0].ReadinessProbe).To(BeNil())

			})

		})

		When("removing the PDB and HPA from a EchoAPI instance", func() {

			// Resource Versions
			rvs := make(map[string]string)

			BeforeEach(func() {
				Eventually(func() error {

					echoapi := &saasv1alpha1.EchoAPI{}
					if err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						echoapi,
					); err != nil {
						return err
					}

					rvs["deployment/echoapi"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "echo-api", namespace, timeout, poll)

					patch := client.MergeFrom(echoapi.DeepCopy())
					echoapi.Spec.Replicas = pointer.Int32(0)
					echoapi.Spec.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{}
					echoapi.Spec.PDB = &saasv1alpha1.PodDisruptionBudgetSpec{}

					return k8sClient.Patch(context.Background(), echoapi, patch)

				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("removes the EchoAPI disabled resources", func() {

				dep := &appsv1.Deployment{}
				By("updating the EchoAPI workload",
					(&testutil.ExpectedWorkload{
						Name:        "echo-api",
						Namespace:   namespace,
						Replicas:    0,
						HPA:         false,
						PDB:         false,
						PodMonitor:  true,
						LastVersion: rvs["deployment/echoapi"],
					}).Assert(k8sClient, dep, timeout, poll))

			})

		})

	})

})
