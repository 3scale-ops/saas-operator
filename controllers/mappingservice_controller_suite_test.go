package controllers

import (
	"context"
	"time"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	externalsecretsv1alpha1 "github.com/3scale/saas-operator/pkg/apis/externalsecrets/v1alpha1"
	grafanav1alpha1 "github.com/3scale/saas-operator/pkg/apis/grafana/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("MappingService controller", func() {
	var namespace string
	var mappingservice *saasv1alpha1.MappingService

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

	When("deploying a defaulted MappingService instance", func() {

		BeforeEach(func() {

			By("creating a MappingService simple resource", func() {
				mappingservice = &saasv1alpha1.MappingService{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "instance",
						Namespace: namespace,
					},
					Spec: saasv1alpha1.MappingServiceSpec{
						Config: saasv1alpha1.MappingServiceConfig{
							APIHost: "example.com",
							SystemAdminToken: saasv1alpha1.SecretReference{
								FromVault: &saasv1alpha1.VaultSecretReference{
									Path: "some-path",
									Key:  "some-key",
								},
							},
						},
					},
				}
				err := k8sClient.Create(context.Background(), mappingservice)
				Expect(err).ToNot(HaveOccurred())
				Eventually(func() error {
					return k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, mappingservice)
				}, timeout, poll).ShouldNot(HaveOccurred())
			})

		})

		It("creates the required MappingService resources", func() {

			dep := &appsv1.Deployment{}
			By("deploying a MappingService workload",
				checkWorkloadResources(dep,
					expectedWorkload{
						Name:          "mapping-service",
						Namespace:     namespace,
						Replicas:      2,
						ContainerName: "mapping-service",
						PDB:           true,
						HPA:           true,
						PodMonitor:    true,
					},
				),
			)
			for _, env := range dep.Spec.Template.Spec.Containers[0].Env {
				switch env.Name {
				case "MASTER_ACCESS_TOKEN":
					Expect(env.ValueFrom.SecretKeyRef.LocalObjectReference.Name).To(Equal("mapping-service-system-master-access-token"))
				case "API_HOST":
					Expect(env.Value).To(Equal("example.com"))
				}
			}
			Expect(dep.Spec.Template.Spec.Volumes).To(HaveLen(0))

			svc := &corev1.Service{}
			By("deploying a MappingService service",
				checkResource(svc,
					expectedResource{
						Name:      "mapping-service",
						Namespace: namespace,
					},
				),
			)
			Expect(svc.Spec.Selector["deployment"]).To(Equal("mapping-service"))
			Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("mapping-service"))

			es := &externalsecretsv1alpha1.ExternalSecret{}
			By("deploying the MappingService System Token external secret",
				checkResource(
					es,
					expectedResource{
						Name:      "mapping-service-system-master-access-token",
						Namespace: namespace,
					},
				),
			)

			Expect(es.Spec.RefreshInterval.ToUnstructured()).To(Equal("1m0s"))
			Expect(es.Spec.SecretStoreRef.Name).To(Equal("vault-mgmt"))
			Expect(es.Spec.SecretStoreRef.Kind).To(Equal("ClusterSecretStore"))

			for _, data := range es.Spec.Data {
				switch data.SecretKey {
				case "MASTER_ACCESS_TOKEN":
					Expect(data.RemoteRef.Property).To(Equal("some-key"))
					Expect(data.RemoteRef.Key).To(Equal("some-path"))
				}
			}

			By("deploying the MappingService grafana dashboard",
				checkResource(
					&grafanav1alpha1.GrafanaDashboard{},
					expectedResource{
						Name:      "mapping-service",
						Namespace: namespace,
					},
				),
			)

		})

		When("updating a MappingService resource with customizations", func() {

			// Resource Versions
			rvs := make(map[string]string)

			BeforeEach(func() {
				Eventually(func() error {

					mappingservice := &saasv1alpha1.MappingService{}
					if err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						mappingservice,
					); err != nil {
						return err
					}

					rvs["mapping-service"] = getResourceVersion(
						mappingservice, "instance", namespace,
					)
					rvs["deployment/mappingservice"] = getResourceVersion(
						&appsv1.Deployment{}, "mapping-service", namespace,
					)

					patch := client.MergeFrom(mappingservice.DeepCopy())
					mappingservice.Spec.Config.APIHost = "updated-example.com"
					mappingservice.Spec.Config.SystemAdminToken.FromVault.Path = "secret/data/updated-path"
					mappingservice.Spec.Config.SystemAdminToken.FromVault.RefreshInterval = &metav1.Duration{Duration: 1 * time.Second}
					mappingservice.Spec.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{
						MinReplicas: pointer.Int32(3),
					}
					mappingservice.Spec.LivenessProbe = &saasv1alpha1.ProbeSpec{}
					mappingservice.Spec.ReadinessProbe = &saasv1alpha1.ProbeSpec{}
					mappingservice.Spec.GrafanaDashboard = &saasv1alpha1.GrafanaDashboardSpec{}

					return k8sClient.Patch(context.Background(), mappingservice, patch)

				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("updates the MappingService resources", func() {

				dep := &appsv1.Deployment{}
				By("updating the MappingService workload",
					checkWorkloadResources(dep,
						expectedWorkload{
							Name:          "mapping-service",
							Namespace:     namespace,
							Replicas:      3,
							ContainerName: "mapping-service",
							PDB:           true,
							HPA:           true,
							PodMonitor:    true,
							LastVersion:   rvs["deployment/mappingservice"],
						},
					),
				)
				for _, env := range dep.Spec.Template.Spec.Containers[0].Env {
					switch env.Name {
					case "MASTER_ACCESS_TOKEN":
						Expect(dep.Spec.Template.Spec.Containers[0].Env[0].ValueFrom.SecretKeyRef.LocalObjectReference.Name).To(Equal("mapping-service-system-master-access-token"))
					case "API_HOST":
						Expect(dep.Spec.Template.Spec.Containers[0].Env[1].Value).To(Equal("updated-example.com"))
					}
				}
				Expect(dep.Spec.Template.Spec.Containers[0].LivenessProbe).To(BeNil())
				Expect(dep.Spec.Template.Spec.Containers[0].ReadinessProbe).To(BeNil())

				es := &externalsecretsv1alpha1.ExternalSecret{}
				By("updating the MappingService System Token external secret",
					checkResource(
						es,
						expectedResource{
							Name:      "mapping-service-system-master-access-token",
							Namespace: namespace,
						},
					),
				)

				Expect(es.Spec.RefreshInterval.ToUnstructured()).To(Equal("1s"))

				for _, data := range es.Spec.Data {
					switch data.SecretKey {
					case "MASTER_ACCESS_TOKEN":
						Expect(data.RemoteRef.Key).To(Equal("updated-path"))
					}
				}

				By("ensuring the MappingService grafana dashboard is gone",
					checkResource(
						&grafanav1alpha1.GrafanaDashboard{},
						expectedResource{
							Name:      "mapping-service",
							Namespace: namespace,
							Missing:   true,
						},
					),
				)

			})

		})

		// Disabled due to https://github.com/3scale-ops/saas-operator/issues/126
		if flag_executeRemoveTests {

			When("removing the PDB and HPA from a MappingService instance", func() {

				// Resource Versions
				rvs := make(map[string]string)

				BeforeEach(func() {
					Eventually(func() error {

						mappingservice := &saasv1alpha1.MappingService{}
						if err := k8sClient.Get(
							context.Background(),
							types.NamespacedName{Name: "instance", Namespace: namespace},
							mappingservice,
						); err != nil {
							return err
						}

						rvs["deployment/mappingservice"] = getResourceVersion(
							&appsv1.Deployment{}, "mapping-service", namespace,
						)

						patch := client.MergeFrom(mappingservice.DeepCopy())
						mappingservice.Spec.Replicas = pointer.Int32(0)
						mappingservice.Spec.HPA = &saasv1alpha1.HorizontalPodAutoscalerSpec{}
						mappingservice.Spec.PDB = &saasv1alpha1.PodDisruptionBudgetSpec{}

						return k8sClient.Patch(context.Background(), mappingservice, patch)

					}, timeout, poll).ShouldNot(HaveOccurred())
				})

				It("removes the MappingService disabled resources", func() {

					dep := &appsv1.Deployment{}
					By("updating the MappingService workload",
						checkWorkloadResources(dep,
							expectedWorkload{
								Name:        "mapping-service",
								Namespace:   namespace,
								Replicas:    0,
								HPA:         false,
								PDB:         false,
								PodMonitor:  true,
								LastVersion: rvs["deployment/mappingservice"],
							},
						),
					)

				})

			})

		}

	})

})
