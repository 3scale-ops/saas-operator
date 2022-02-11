/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/goombaio/namegenerator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	externalsecretsv1alpha1 "github.com/3scale/saas-operator/pkg/apis/externalsecrets/v1alpha1"
	grafanav1alpha1 "github.com/3scale/saas-operator/pkg/apis/grafana/v1alpha1"
	secretsmanagerv1alpha1 "github.com/3scale/saas-operator/pkg/apis/secrets-manager/v1alpha1"
	"github.com/3scale/saas-operator/pkg/reconcilers/workloads"
	// +kubebuilder:scaffold:imports
)

type expectedWorkload struct {
	Namespace      string
	Name           string
	Replicas       int32
	ContainerName  string
	ContainerImage string
	ContainterArgs []string
	HPA            bool
	PDB            bool
	PodMonitor     bool
}

func checkWorkloadResources(dep *appsv1.Deployment, ew expectedWorkload) func() {
	return func() {
		Eventually(func() error {
			return k8sClient.Get(
				context.Background(),
				types.NamespacedName{Name: ew.Name, Namespace: ew.Namespace},
				dep,
			)
		}, timeout, poll).ShouldNot(HaveOccurred())

		Expect(dep.Spec.Replicas).To(Equal(pointer.Int32Ptr(ew.Replicas)))

		if ew.ContainerName != "" {
			Expect(dep.Spec.Template.Spec.Containers[0].Name).To(Equal(ew.ContainerName))
		}

		if ew.ContainerImage != "" {
			Expect(dep.Spec.Template.Spec.Containers[0].Image).To(Equal(ew.ContainerImage))
		}

		if ew.ContainterArgs != nil {
			Expect(dep.Spec.Template.Spec.Containers[0].Args).To(Equal(ew.ContainterArgs))
		}

		hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{}
		By(fmt.Sprintf("%s workload HPA", ew.Name),
			checkResource(hpa, expectedResource{
				Name:      ew.Name,
				Namespace: ew.Namespace, Missing: !ew.HPA,
			}),
		)
		if ew.HPA {
			Expect(hpa.Spec.ScaleTargetRef.Kind).Should(Equal("Deployment"))
			Expect(hpa.Spec.ScaleTargetRef.Name).Should(Equal(ew.Name))
		}

		pdb := &policyv1beta1.PodDisruptionBudget{}
		By(fmt.Sprintf("%s workload PDB", ew.Name),
			checkResource(pdb, expectedResource{
				Name:      ew.Name,
				Namespace: ew.Namespace, Missing: !ew.PDB,
			}),
		)
		if ew.PDB {
			Expect(pdb.Spec.Selector.MatchLabels["deployment"]).Should(Equal(ew.Name))
		}

		pm := &monitoringv1.PodMonitor{}
		By(fmt.Sprintf("%s workload PodMonitor", ew.Name),
			checkResource(pm, expectedResource{
				Name:      ew.Name,
				Namespace: ew.Namespace, Missing: !ew.PodMonitor,
			}),
		)
		if ew.PodMonitor {
			Expect(pm.Spec.Selector.MatchLabels["deployment"]).Should(Equal(ew.Name))
		}

	}
}

type expectedResource struct {
	Namespace string
	Name      string
	Missing   bool
}

func checkResource(r client.Object, er expectedResource) func() {

	if er.Missing {
		return func() {
			By(fmt.Sprintf("%s object does NOT exist", er.Name))
			Eventually(func() error {
				return k8sClient.Get(context.Background(),
					types.NamespacedName{Name: er.Name, Namespace: er.Namespace}, r,
				)
			}, timeout, poll).Should(HaveOccurred())
		}
	}

	return func() {
		By(fmt.Sprintf("%s object does exists", er.Name))
		Eventually(func() error {
			return k8sClient.Get(context.Background(),
				types.NamespacedName{Name: er.Name, Namespace: er.Namespace}, r,
			)
		}, timeout, poll).ShouldNot(HaveOccurred())
	}

}

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	cfg           *rest.Config
	k8sClient     client.Client
	testEnv       *envtest.Environment
	nameGenerator namegenerator.Generator
	timeout       time.Duration = 45 * time.Second
	poll          time.Duration = 5 * time.Second
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(false)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "config", "crd", "bases"),
			filepath.Join("..", "config", "test", "external-apis"),
		},
	}

	seed := time.Now().UTC().UnixNano()
	nameGenerator = namegenerator.NewNameGenerator(seed)

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = saasv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())
	err = monitoringv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())
	err = grafanav1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())
	err = secretsmanagerv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())
	err = externalsecretsv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
		// Disable the metrics port to allow running the
		// test suite in parallel
		MetricsBindAddress: "0",
	})
	Expect(err).ToNot(HaveOccurred())

	k8sClient = mgr.GetClient()
	Expect(k8sClient).ToNot(BeNil())

	go func() {
		defer GinkgoRecover()
		err = mgr.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred())
	}()

	// Add controllers for testing
	err = (&AutoSSLReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr, mgr.GetEventRecorderFor("AutoSSL"), false),
		Log:                ctrl.Log.WithName("controllers").WithName("AutoSSL"),
	}).SetupWithManager(mgr)
	Expect(err).ToNot(HaveOccurred())

	err = (&ApicastReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr, mgr.GetEventRecorderFor("Apicast"), false),
		Log:                ctrl.Log.WithName("controllers").WithName("Apicast"),
	}).SetupWithManager(mgr)
	Expect(err).ToNot(HaveOccurred())

	err = (&EchoAPIReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr, mgr.GetEventRecorderFor("EchoAPI"), false),
		Log:                ctrl.Log.WithName("controllers").WithName("EchoAPI"),
	}).SetupWithManager(mgr)
	Expect(err).ToNot(HaveOccurred())

	err = (&MappingServiceReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr, mgr.GetEventRecorderFor("MappingService"), false),
		Log:                ctrl.Log.WithName("controllers").WithName("MappingService"),
	}).SetupWithManager(mgr)
	Expect(err).ToNot(HaveOccurred())

	err = (&CORSProxyReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr, mgr.GetEventRecorderFor("CORSProxy"), false),
		Log:                ctrl.Log.WithName("controllers").WithName("CORSProxy"),
	}).SetupWithManager(mgr)
	Expect(err).ToNot(HaveOccurred())

	err = (&BackendReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr, mgr.GetEventRecorderFor("Backend"), false),
		Log:                ctrl.Log.WithName("controllers").WithName("Backend"),
	}).SetupWithManager(mgr)
	Expect(err).ToNot(HaveOccurred())

	err = (&SystemReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr, mgr.GetEventRecorderFor("System"), false),
		Log:                ctrl.Log.WithName("controllers").WithName("System"),
	}).SetupWithManager(mgr)
	Expect(err).ToNot(HaveOccurred())

	err = (&ZyncReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr, mgr.GetEventRecorderFor("Zync"), false),
		Log:                ctrl.Log.WithName("controllers").WithName("Zync"),
	}).SetupWithManager(mgr)
	Expect(err).ToNot(HaveOccurred())

}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
