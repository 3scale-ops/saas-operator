package controllers

import (
	"context"
	"time"

	"github.com/3scale-ops/basereconciler/util"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	testutil "github.com/3scale-ops/saas-operator/test/util"
	externalsecretsv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	grafanav1alpha1 "github.com/grafana-operator/grafana-operator/v4/api/integreatly/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Zync controller", func() {
	var namespace string
	var zync *saasv1alpha1.Zync

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

	When("deploying a defaulted Zync instance", func() {

		BeforeEach(func() {
			By("creating a Zync simple resource")
			zync = &saasv1alpha1.Zync{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance",
					Namespace: namespace,
				},
				Spec: saasv1alpha1.ZyncSpec{
					Config: saasv1alpha1.ZyncConfig{
						DatabaseDSN: saasv1alpha1.SecretReference{
							FromVault: &saasv1alpha1.VaultSecretReference{
								Path: "some-path-db",
								Key:  "some-key-db",
							},
						},
						SecretKeyBase: saasv1alpha1.SecretReference{
							FromVault: &saasv1alpha1.VaultSecretReference{
								Path: "some-path-base",
								Key:  "some-key-base",
							},
						},
						ZyncAuthToken: saasv1alpha1.SecretReference{
							FromVault: &saasv1alpha1.VaultSecretReference{
								Path: "some-path-token",
								Key:  "some-key-token",
							},
						},
						Bugsnag: &saasv1alpha1.BugsnagSpec{
							ReleaseStage: util.Pointer("staging"),
							APIKey: saasv1alpha1.SecretReference{
								FromVault: &saasv1alpha1.VaultSecretReference{
									Path: "some-path-bugsnag",
									Key:  "some-key-bugsnag",
								},
							},
						},
					},
				},
			}
			err := k8sClient.Create(context.Background(), zync)
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{Name: "instance", Namespace: namespace}, zync)
			}, timeout, poll).ShouldNot(HaveOccurred())
		})

		It("creates the required Zync resources", func() {

			dep := &appsv1.Deployment{}
			By("deploying a Zync workload",
				(&testutil.ExpectedWorkload{
					Name:          "zync",
					Namespace:     namespace,
					Replicas:      2,
					ContainerName: "zync",
					PDB:           true,
					HPA:           true,
					PodMonitor:    true,
				}).Assert(k8sClient, dep, timeout, poll))

			for _, env := range dep.Spec.Template.Spec.Containers[0].Env {
				switch env.Name {
				case "RAILS_ENV":
					Expect(env.Value).To(Equal("development"))
				case "RAILS_MAX_THREADS":
					Expect(env.Value).To(Equal("10"))
				case "ZYNC_AUTHENTICATION_TOKEN":
					Expect(env.ValueFrom.SecretKeyRef.LocalObjectReference.Name).To(Equal("zync"))
				}
			}
			Expect(dep.Spec.Template.Spec.Volumes).To(HaveLen(0))
			Expect(dep.Spec.Template.Spec.Containers[0].LivenessProbe).ToNot(BeNil())
			Expect(dep.Spec.Template.Spec.Containers[0].ReadinessProbe).ToNot(BeNil())

			By("deploying a Zync-Que workload",
				(&testutil.ExpectedWorkload{
					Name:          "zync-que",
					Namespace:     namespace,
					Replicas:      2,
					ContainerName: "zync-que",
					PDB:           true,
					HPA:           true,
					PodMonitor:    true,
				}).Assert(k8sClient, dep, timeout, poll))

			for _, env := range dep.Spec.Template.Spec.Containers[0].Env {
				switch env.Name {
				case "RAILS_ENV":
					Expect(env.Value).To(Equal("development"))
				case "RAILS_LOG_LEVEL":
					Expect(env.Value).To(Equal("info"))
				case "ZYNC_AUTHENTICATION_TOKEN":
					Expect(env.ValueFrom.SecretKeyRef.LocalObjectReference.Name).To(Equal("zync"))
				}
			}
			Expect(dep.Spec.Template.Spec.Volumes).To(HaveLen(0))
			Expect(dep.Spec.Template.Spec.Containers[0].LivenessProbe).ToNot(BeNil())
			Expect(dep.Spec.Template.Spec.Containers[0].ReadinessProbe).ToNot(BeNil())

			svc := &corev1.Service{}
			By("deploying a Zync service",
				(&testutil.ExpectedResource{Name: "zync", Namespace: namespace}).
					Assert(k8sClient, svc, timeout, poll))

			Expect(svc.Spec.Selector["deployment"]).To(Equal("zync"))
			Expect(svc.Spec.Selector["saas.3scale.net/traffic"]).To(Equal("zync"))

			es := &externalsecretsv1beta1.ExternalSecret{}
			By("deploying the Zync external secret",
				(&testutil.ExpectedResource{Name: "zync", Namespace: namespace}).
					Assert(k8sClient, es, timeout, poll))

			Expect(es.Spec.RefreshInterval.ToUnstructured()).To(Equal("1m0s"))
			Expect(es.Spec.SecretStoreRef.Name).To(Equal("vault-mgmt"))
			Expect(es.Spec.SecretStoreRef.Kind).To(Equal("ClusterSecretStore"))

			for _, data := range es.Spec.Data {
				switch data.SecretKey {
				case "DATABASE_URL":
					Expect(data.RemoteRef.Property).To(Equal("some-key-db"))
					Expect(data.RemoteRef.Key).To(Equal("some-path-db"))
				case "SECRET_KEY_BASE":
					Expect(data.RemoteRef.Property).To(Equal("some-key-base"))
					Expect(data.RemoteRef.Key).To(Equal("some-path-base"))
				case "ZYNC_AUTHENTICATION_TOKEN":
					Expect(data.RemoteRef.Property).To(Equal("some-key-token"))
					Expect(data.RemoteRef.Key).To(Equal("some-path-token"))
				case "BUGSNAG_API_KEY":
					Expect(data.RemoteRef.Property).To(Equal("some-key-bugsnag"))
					Expect(data.RemoteRef.Key).To(Equal("some-path-bugsnag"))
				}
			}

			By("deploying the Zync grafana dashboard",
				(&testutil.ExpectedResource{Name: "zync", Namespace: namespace}).
					Assert(k8sClient, &grafanav1alpha1.GrafanaDashboard{}, timeout, poll))

		})

		It("doesn't creates the non-default resources", func() {

			sts := &appsv1.StatefulSet{}
			By("ensuring the zync-console statefulset",
				(&testutil.ExpectedResource{Name: "zync-console", Namespace: namespace, Missing: true}).
					Assert(k8sClient, sts, timeout, poll))
		})

		When("updating a Zync resource with customizations", func() {

			// Resource Versions
			rvs := make(map[string]string)

			BeforeEach(func() {
				Eventually(func() error {

					zync := &saasv1alpha1.Zync{}
					if err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						zync,
					); err != nil {
						return err
					}

					rvs["deployment/zync"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "zync", namespace, timeout, poll)
					rvs["deployment/zync-que"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "zync-que", namespace, timeout, poll)
					rvs["externalsecret/zync"] = testutil.GetResourceVersion(
						k8sClient, &externalsecretsv1beta1.ExternalSecret{}, "zync", namespace, timeout, poll)

					patch := client.MergeFrom(zync.DeepCopy())
					zync.Spec.API = &saasv1alpha1.APISpec{
						HPA: &saasv1alpha1.HorizontalPodAutoscalerSpec{
							MinReplicas: util.Pointer[int32](3),
						},
						LivenessProbe:  &saasv1alpha1.ProbeSpec{},
						ReadinessProbe: &saasv1alpha1.ProbeSpec{},
					}
					zync.Spec.Que = &saasv1alpha1.QueSpec{
						HPA: &saasv1alpha1.HorizontalPodAutoscalerSpec{
							MinReplicas: util.Pointer[int32](3),
						},
						LivenessProbe:  &saasv1alpha1.ProbeSpec{},
						ReadinessProbe: &saasv1alpha1.ProbeSpec{},
					}
					zync.Spec.Config.Rails = &saasv1alpha1.ZyncRailsSpec{
						Environment: util.Pointer("production"),
						MaxThreads:  util.Pointer[int32](12),
						LogLevel:    util.Pointer("debug"),
					}
					zync.Spec.Config.ExternalSecret.RefreshInterval = &metav1.Duration{Duration: 1 * time.Second}
					zync.Spec.Config.ExternalSecret.SecretStoreRef = &saasv1alpha1.ExternalSecretSecretStoreReferenceSpec{
						Name: util.Pointer("other-store"),
						Kind: util.Pointer("SecretStore"),
					}
					zync.Spec.Config.SecretKeyBase.FromVault.Path = "secret/data/updated-path"

					zync.Spec.GrafanaDashboard = &saasv1alpha1.GrafanaDashboardSpec{}

					return k8sClient.Patch(context.Background(), zync, patch)

				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("updates the Zync resources", func() {

				dep := &appsv1.Deployment{}
				By("updating the Zync workload",
					(&testutil.ExpectedWorkload{
						Name:          "zync",
						Namespace:     namespace,
						Replicas:      3,
						ContainerName: "zync",
						PDB:           true,
						HPA:           true,
						PodMonitor:    true,
						LastVersion:   rvs["deployment/zync"],
					}).Assert(k8sClient, dep, timeout, poll))

				for _, env := range dep.Spec.Template.Spec.Containers[0].Env {
					switch env.Name {
					case "RAILS_ENV":
						Expect(env.Value).To(Equal("production"))
					case "RAILS_MAX_THREADS":
						Expect(env.Value).To(Equal("12"))
					case "ZYNC_AUTHENTICATION_TOKEN":
						Expect(env.ValueFrom.SecretKeyRef.LocalObjectReference.Name).To(Equal("zync"))
					}
				}
				Expect(dep.Spec.Template.Spec.Containers[0].LivenessProbe).To(BeNil())
				Expect(dep.Spec.Template.Spec.Containers[0].ReadinessProbe).To(BeNil())

				By("updating the Zync-Que workload",
					(&testutil.ExpectedWorkload{
						Name:          "zync-que",
						Namespace:     namespace,
						Replicas:      3,
						ContainerName: "zync-que",
						ContainterCmd: []string{
							"/usr/bin/bash",
							"-c",
							"bundle exec rake 'que[--worker-count 10]'",
						},
						PDB:         true,
						HPA:         true,
						PodMonitor:  true,
						LastVersion: rvs["deployment/zync-que"],
					}).Assert(k8sClient, dep, timeout, poll))

				for _, env := range dep.Spec.Template.Spec.Containers[0].Env {
					switch env.Name {
					case "RAILS_ENV":
						Expect(env.Value).To(Equal("production"))
					case "RAILS_LOG_LEVEL":
						Expect(env.Value).To(Equal("debug"))
					case "ZYNC_AUTHENTICATION_TOKEN":
						Expect(env.ValueFrom.SecretKeyRef.LocalObjectReference.Name).To(Equal("zync"))
					}
				}
				Expect(dep.Spec.Template.Spec.Containers[0].LivenessProbe).To(BeNil())
				Expect(dep.Spec.Template.Spec.Containers[0].ReadinessProbe).To(BeNil())

				es := &externalsecretsv1beta1.ExternalSecret{}
				By("updating the Zync external secret",
					(&testutil.ExpectedResource{
						Name:        "zync",
						Namespace:   namespace,
						LastVersion: rvs["externalsecret/zync"],
					}).Assert(k8sClient, es, timeout, poll))

				Expect(es.Spec.RefreshInterval.ToUnstructured()).To(Equal("1s"))
				Expect(es.Spec.SecretStoreRef.Name).To(Equal("other-store"))
				Expect(es.Spec.SecretStoreRef.Kind).To(Equal("SecretStore"))

				for _, data := range es.Spec.Data {
					switch data.SecretKey {
					case "SECRET_KEY_BASE":
						Expect(data.RemoteRef.Key).To(Equal("updated-path"))
					}
				}

				By("ensuring the Zync grafana dashboard is gone",
					(&testutil.ExpectedResource{
						Name:      "zync",
						Namespace: namespace,
						Missing:   true,
					}).Assert(k8sClient, &grafanav1alpha1.GrafanaDashboard{}, timeout, poll))

			})

		})

		When("updating a Zync resource with console", func() {

			BeforeEach(func() {
				Eventually(func() error {
					err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						zync,
					)
					Expect(err).ToNot(HaveOccurred())
					patch := client.MergeFrom(zync.DeepCopy())
					zync.Spec.Console = &saasv1alpha1.ZyncRailsConsoleSpec{
						Enabled: util.Pointer(true),
						Image: &saasv1alpha1.ImageSpec{
							Name: util.Pointer("newImage"),
							Tag:  util.Pointer("newTag"),
						},
					}
					return k8sClient.Patch(context.Background(), zync, patch)
				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("creates the required console resources", func() {

				sts := &appsv1.StatefulSet{}
				By("deploying the zync-console StatefulSet",
					(&testutil.ExpectedResource{Name: "zync-console", Namespace: namespace}).
						Assert(k8sClient, sts, timeout, poll))

				Expect(sts.Spec.Template.Spec.Containers[0].Image).Should((Equal("newImage:newTag")))

				pdb := &policyv1.PodDisruptionBudget{}
				By("ensuring the zync-console PDB",
					(&testutil.ExpectedResource{Name: "zync-console", Namespace: namespace, Missing: true}).
						Assert(k8sClient, pdb, timeout, poll))

				hpa := &autoscalingv2.HorizontalPodAutoscaler{}
				By("ensuring the zync-console HPA",
					(&testutil.ExpectedResource{Name: "zync-console", Namespace: namespace, Missing: true}).
						Assert(k8sClient, hpa, timeout, poll))

			})

		})

		When("removing the PDB and HPA from a Zync instance", func() {

			// Resource Versions
			rvs := make(map[string]string)

			BeforeEach(func() {
				Eventually(func() error {

					zync := &saasv1alpha1.Zync{}
					if err := k8sClient.Get(
						context.Background(),
						types.NamespacedName{Name: "instance", Namespace: namespace},
						zync,
					); err != nil {
						return err
					}

					rvs["deployment/zync"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "zync", namespace, timeout, poll)
					rvs["deployment/zync-que"] = testutil.GetResourceVersion(
						k8sClient, &appsv1.Deployment{}, "zync-que", namespace, timeout, poll)

					patch := client.MergeFrom(zync.DeepCopy())

					zync.Spec.API = &saasv1alpha1.APISpec{
						Replicas: util.Pointer[int32](0),
						HPA:      &saasv1alpha1.HorizontalPodAutoscalerSpec{},
						PDB:      &saasv1alpha1.PodDisruptionBudgetSpec{},
					}

					zync.Spec.Que = &saasv1alpha1.QueSpec{
						Replicas: util.Pointer[int32](0),
						HPA:      &saasv1alpha1.HorizontalPodAutoscalerSpec{},
						PDB:      &saasv1alpha1.PodDisruptionBudgetSpec{},
					}

					return k8sClient.Patch(context.Background(), zync, patch)

				}, timeout, poll).ShouldNot(HaveOccurred())
			})

			It("removes the Zync disabled resources", func() {

				By("updating the Zync workload",
					(&testutil.ExpectedWorkload{
						Name:        "zync",
						Namespace:   namespace,
						Replicas:    0,
						HPA:         false,
						PDB:         false,
						PodMonitor:  true,
						LastVersion: rvs["deployment/zync"],
					}).Assert(k8sClient, &appsv1.Deployment{}, timeout, poll))

				By("updating the Zync-Que workload",
					(&testutil.ExpectedWorkload{
						Name:        "zync-que",
						Namespace:   namespace,
						Replicas:    0,
						HPA:         false,
						PDB:         false,
						PodMonitor:  true,
						LastVersion: rvs["deployment/zync-que"],
					}).Assert(k8sClient, &appsv1.Deployment{}, timeout, poll))

			})

		})

	})

})
