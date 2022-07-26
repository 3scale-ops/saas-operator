package controllers

import (
	"context"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	grafanav1alpha1 "github.com/3scale/saas-operator/pkg/apis/grafana/v1alpha1"
	testutil "github.com/3scale/saas-operator/test/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("AutoSSL controller", func() {
	var namespace string
	var autossl *saasv1alpha1.AutoSSL

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

	When("deploying a defaulted AutoSSL instance", func() {

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

		It("creates the required AutoSSL resources", func() {

			dep := &appsv1.Deployment{}
			By("deploying an autossl workload",
				(&testutil.ExpectedWorkload{
					Name:          "autossl",
					Namespace:     namespace,
					Replicas:      2,
					ContainerName: "autossl",
					PDB:           true,
					HPA:           true,
					PodMonitor:    true,
				}).Assert(k8sClient, dep, timeout, poll))

			Expect(dep.Spec.Template.Spec.Volumes).To(HaveLen(2))
			Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("autossl-cache"))
			Expect(dep.Spec.Template.Spec.Volumes[0].VolumeSource.EmptyDir).ShouldNot(BeNil())
			Expect(dep.Spec.Template.Spec.Volumes[1].Name).To(Equal("nginx-cache"))
			Expect(dep.Spec.Template.Spec.Volumes[1].VolumeSource.EmptyDir).ShouldNot(BeNil())
			Expect(dep.Spec.Template.Spec.Containers[0].VolumeMounts).To(HaveLen(2))
			Expect(dep.Spec.Template.Spec.Containers[0].VolumeMounts[0].Name).To(Equal("autossl-cache"))
			Expect(dep.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath).To(Equal("/etc/resty-auto-ssl/"))
			Expect(dep.Spec.Template.Spec.Containers[0].VolumeMounts[1].Name).To(Equal("nginx-cache"))
			Expect(dep.Spec.Template.Spec.Containers[0].VolumeMounts[1].MountPath).To(Equal("/var/lib/nginx"))

			svc := &corev1.Service{}
			By("deploying an autossl service",
				(&testutil.ExpectedResource{
					Name:      "autossl",
					Namespace: namespace,
				}).Assert(k8sClient, svc, timeout, poll))

			Expect(svc.Spec.Selector["deployment"]).To(Equal("autossl"))
			Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("autossl"))

			By("deploying an autossl grafana dashboard",
				(&testutil.ExpectedResource{
					Name:      "autossl",
					Namespace: namespace,
				}).Assert(k8sClient, &grafanav1alpha1.GrafanaDashboard{}, timeout, poll))

		})

		It("doesn't creates the non-default resources", func() {

			dep := &appsv1.Deployment{}
			By("ensuring an autossl-canary workload",
				(&testutil.ExpectedResource{
					Name:      "autossl-canary",
					Namespace: namespace,
					Missing:   true,
				}).Assert(k8sClient, dep, timeout, poll))

		})

		When("updating a AutoSSL resource with customizations", func() {

			// Resource Versions
			rvs := make(map[string]string)

			BeforeEach(func() {
				Eventually(func() error {

					autossl := &saasv1alpha1.AutoSSL{}
					if err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						autossl,
					); err != nil {
						return err
					}

					rvs["autossl"] = testutil.GetResourceVersion(
						k8sClient, autossl, "instance", namespace, timeout, poll)
					rvs["deployment/autossl"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "autossl", namespace, timeout, poll)

					patch := client.MergeFrom(autossl.DeepCopy())
					autossl.Spec.Config.ContactEmail = "updated-example@3scale.net"
					autossl.Spec.Config.VerificationEndpoint = "updated-example.com/verification"
					autossl.Spec.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{
						MinReplicas: pointer.Int32(3),
					}
					autossl.Spec.GrafanaDashboard = &saasv1alpha1.GrafanaDashboardSpec{}
					autossl.Spec.LivenessProbe = &saasv1alpha1.ProbeSpec{}
					autossl.Spec.ReadinessProbe = &saasv1alpha1.ProbeSpec{}

					return k8sClient.Patch(context.Background(), autossl, patch)

				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("updates the AutoSSL resources", func() {

				By("ensuring the AutoSSL grafana dashboard",
					(&testutil.ExpectedResource{
						Name:      "autossl",
						Namespace: namespace,
						Missing:   true,
					}).Assert(k8sClient, &grafanav1alpha1.GrafanaDashboard{}, timeout, poll))

				dep := &appsv1.Deployment{}
				By("updating the AutoSSL workload",
					(&testutil.ExpectedWorkload{
						Name:          "autossl",
						Namespace:     namespace,
						Replicas:      3,
						ContainerName: "autossl",
						HPA:           true,
						PDB:           true,
						PodMonitor:    true,
						LastVersion:   rvs["deployment/autossl"],
					}).Assert(k8sClient, dep, timeout, poll))

				Expect(dep.Spec.Template.Spec.Containers[0].Name).To(Equal("autossl"))
				for _, env := range dep.Spec.Template.Spec.Containers[0].Env {
					switch env.Name {
					case "CONTACT_EMAIL":
						Expect(env.Value).To(Equal("updated-example@3scale.net"))
					case "VERIFICATION_ENDPOINT":
						Expect(env.Value).To(Equal("updated-example.com/verification"))
					}
				}
				Expect(dep.Spec.Template.Spec.Containers[0].LivenessProbe).To(BeNil())
				Expect(dep.Spec.Template.Spec.Containers[0].ReadinessProbe).To(BeNil())
				Expect(dep.Spec.Template.Spec.Volumes).To(HaveLen(2))
				Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("autossl-cache"))
				Expect(dep.Spec.Template.Spec.Volumes[0].VolumeSource.EmptyDir).ShouldNot(BeNil())
				Expect(dep.Spec.Template.Spec.Volumes[1].Name).To(Equal("nginx-cache"))
				Expect(dep.Spec.Template.Spec.Volumes[1].VolumeSource.EmptyDir).ShouldNot(BeNil())
				Expect(dep.Spec.Template.Spec.Containers[0].VolumeMounts).To(HaveLen(2))
				Expect(dep.Spec.Template.Spec.Containers[0].VolumeMounts[0].Name).To(Equal("autossl-cache"))
				Expect(dep.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath).To(Equal("/etc/resty-auto-ssl/"))
				Expect(dep.Spec.Template.Spec.Containers[0].VolumeMounts[1].Name).To(Equal("nginx-cache"))
				Expect(dep.Spec.Template.Spec.Containers[0].VolumeMounts[1].MountPath).To(Equal("/var/lib/nginx"))

			})

		})

		When("updating a AutoSSL resource with canary", func() {

			// Resource Versions
			rvs := make(map[string]string)

			BeforeEach(func() {
				Eventually(func() error {
					autossl := &saasv1alpha1.AutoSSL{}
					if err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						autossl,
					); err != nil {
						return err
					}

					rvs["svc/autossl"] = testutil.GetResourceVersion(
						k8sClient, &corev1.Service{}, "autossl", namespace, timeout, poll)
					rvs["deployment/autossl"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "autossl", namespace, timeout, poll)

					patch := client.MergeFrom(autossl.DeepCopy())
					autossl.Spec.Canary = &saasv1alpha1.Canary{
						ImageName: pointer.StringPtr("newImage"),
						ImageTag:  pointer.StringPtr("newTag"),
						Replicas:  pointer.Int32Ptr(1),
					}

					return k8sClient.Patch(context.Background(), autossl, patch)
				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("creates the required canary resources", func() {

				dep := &appsv1.Deployment{}
				By("deploying a autossl-canary workload",
					(&testutil.ExpectedWorkload{
						Name:          "autossl-canary",
						Namespace:     namespace,
						Replicas:      1,
						ContainerName: "autossl",
						PodMonitor:    true,
					}).Assert(k8sClient, dep, timeout, poll))

				Expect(dep.Spec.Template.Spec.Containers[0].Name).To(Equal("autossl"))
				for _, env := range dep.Spec.Template.Spec.Containers[0].Env {
					switch env.Name {
					case "CONTACT_EMAIL":
						Expect(env.Value).To(Equal("test@example.com"))
					case "VERIFICATION_ENDPOINT":
						Expect(env.Value).To(Equal("example.com/verification"))
					}
				}
				Expect(dep.Spec.Template.Spec.Volumes).To(HaveLen(2))
				Expect(dep.Spec.Template.Spec.Volumes[0].Name).To(Equal("autossl-cache"))
				Expect(dep.Spec.Template.Spec.Volumes[0].VolumeSource.EmptyDir).ShouldNot(BeNil())
				Expect(dep.Spec.Template.Spec.Volumes[1].Name).To(Equal("nginx-cache"))
				Expect(dep.Spec.Template.Spec.Volumes[1].VolumeSource.EmptyDir).ShouldNot(BeNil())
				Expect(dep.Spec.Template.Spec.Containers[0].VolumeMounts).To(HaveLen(2))
				Expect(dep.Spec.Template.Spec.Containers[0].VolumeMounts[0].Name).To(Equal("autossl-cache"))
				Expect(dep.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath).To(Equal("/etc/resty-auto-ssl/"))
				Expect(dep.Spec.Template.Spec.Containers[0].VolumeMounts[1].Name).To(Equal("nginx-cache"))
				Expect(dep.Spec.Template.Spec.Containers[0].VolumeMounts[1].MountPath).To(Equal("/var/lib/nginx"))

				svc := &corev1.Service{}
				By("keeping the autossl service deployment label selector",
					(&testutil.ExpectedResource{
						Name: "autossl", Namespace: namespace,
					}).Assert(k8sClient, svc, timeout, poll))

				Expect(svc.Spec.Selector["deployment"]).To(Equal("autossl"))
				Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("autossl"))

			})

			When("enabling canary traffic", func() {

				BeforeEach(func() {
					Eventually(func() error {
						autossl := &saasv1alpha1.AutoSSL{}
						if err := k8sClient.Get(
							context.Background(),
							types.NamespacedName{Name: "instance", Namespace: namespace},
							autossl,
						); err != nil {
							return err
						}
						rvs["svc/autossl"] = testutil.GetResourceVersion(
							k8sClient, &corev1.Service{}, "autossl", namespace, timeout, poll)

						patch := client.MergeFrom(autossl.DeepCopy())
						autossl.Spec.Canary = &saasv1alpha1.Canary{
							SendTraffic: *pointer.Bool(true),
						}
						return k8sClient.Patch(context.Background(), autossl, patch)
					}, timeout, poll).ShouldNot(HaveOccurred())
				})

				It("updates the autossl service", func() {

					svc := &corev1.Service{}
					By("removing the autossl service deployment label selector",
						(&testutil.ExpectedResource{
							Name: "autossl", Namespace: namespace,
							LastVersion: rvs["svc/autossl"],
						}).Assert(k8sClient, svc, timeout, poll))

					Expect(svc.Spec.Selector).NotTo(HaveKey("deployment"))
					Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("autossl"))

				})

			})

		})

		When("removing the PDB and HPA from a AutoSSL instance", func() {

			// Resource Versions
			rvs := make(map[string]string)

			BeforeEach(func() {
				Eventually(func() error {

					autossl := &saasv1alpha1.AutoSSL{}
					if err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						autossl,
					); err != nil {
						return err
					}

					rvs["deployment/autossl"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "autossl", namespace, timeout, poll)

					patch := client.MergeFrom(autossl.DeepCopy())
					autossl.Spec.Replicas = pointer.Int32(0)
					autossl.Spec.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{}
					autossl.Spec.PDB = &saasv1alpha1.PodDisruptionBudgetSpec{}

					return k8sClient.Patch(context.Background(), autossl, patch)

				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("removes the AutoSSL disabled resources", func() {

				dep := &appsv1.Deployment{}
				By("updating the AutoSSL workload",
					(&testutil.ExpectedWorkload{
						Name:        "autossl",
						Namespace:   namespace,
						Replicas:    0,
						HPA:         false,
						PDB:         false,
						PodMonitor:  true,
						LastVersion: rvs["deployment/autossl"],
					}).Assert(k8sClient, dep, timeout, poll))

			})

		})

	})

})
