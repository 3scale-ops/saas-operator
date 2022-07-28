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
	"path/filepath"
	"testing"
	"time"

	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	"github.com/3scale/saas-operator/pkg/reconcilers/threads"
	"github.com/3scale/saas-operator/pkg/reconcilers/workloads"
	"github.com/goombaio/namegenerator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	externalsecretsv1beta1 "github.com/3scale/saas-operator/pkg/apis/externalsecrets/v1beta1"
	grafanav1alpha1 "github.com/3scale/saas-operator/pkg/apis/grafana/v1alpha1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	k8sClient     client.Client
	testEnv       *envtest.Environment
	nameGenerator namegenerator.Generator
	timeout       time.Duration = 45 * time.Second
	poll          time.Duration = 5 * time.Second
	ctx           context.Context
	cancel        context.CancelFunc
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
	err = externalsecretsv1beta1.AddToScheme(scheme.Scheme)
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

	ctx, cancel = context.WithCancel(context.Background())

	go func() {
		defer GinkgoRecover()
		err = mgr.Start(ctx)
		Expect(err).ToNot(HaveOccurred())
	}()

	// Add controllers for testing
	err = (&AutoSSLReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr, "AutoSSL", false),
		Log:                ctrl.Log.WithName("controllers").WithName("AutoSSL"),
	}).SetupWithManager(mgr)
	Expect(err).ToNot(HaveOccurred())

	err = (&ApicastReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr, "Apicast", false),
		Log:                ctrl.Log.WithName("controllers").WithName("Apicast"),
	}).SetupWithManager(mgr)
	Expect(err).ToNot(HaveOccurred())

	err = (&EchoAPIReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr, "EchoAPI", false),
		Log:                ctrl.Log.WithName("controllers").WithName("EchoAPI"),
	}).SetupWithManager(mgr)
	Expect(err).ToNot(HaveOccurred())

	err = (&MappingServiceReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr, "MappingService", false),
		Log:                ctrl.Log.WithName("controllers").WithName("MappingService"),
	}).SetupWithManager(mgr)
	Expect(err).ToNot(HaveOccurred())

	err = (&CORSProxyReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr, "CORSProxy", false),
		Log:                ctrl.Log.WithName("controllers").WithName("CORSProxy"),
	}).SetupWithManager(mgr)
	Expect(err).ToNot(HaveOccurred())

	err = (&BackendReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr, "Backend", false),
		Log:                ctrl.Log.WithName("controllers").WithName("Backend"),
	}).SetupWithManager(mgr)
	Expect(err).ToNot(HaveOccurred())

	err = (&RedisShardReconciler{
		Reconciler: basereconciler.NewFromManager(mgr, "RedisShard", false),
		Log:        ctrl.Log.WithName("controllers").WithName("RedisShard"),
	}).SetupWithManager(mgr)
	Expect(err).ToNot(HaveOccurred())

	err = (&SentinelReconciler{
		Reconciler:     basereconciler.NewFromManager(mgr, "Sentinel", false),
		SentinelEvents: threads.NewManager(),
		Metrics:        threads.NewManager(),
		Log:            ctrl.Log.WithName("controllers").WithName("Sentinel"),
	}).SetupWithManager(mgr)
	Expect(err).ToNot(HaveOccurred())

	err = (&SystemReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr, "System", false),
		Log:                ctrl.Log.WithName("controllers").WithName("System"),
	}).SetupWithManager(mgr)
	Expect(err).ToNot(HaveOccurred())

	err = (&ZyncReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr, "Zync", false),
		Log:                ctrl.Log.WithName("controllers").WithName("Zync"),
	}).SetupWithManager(mgr)
	Expect(err).ToNot(HaveOccurred())

}, 60)

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	Eventually(func() error {
		return testEnv.Stop()
	}, timeout, poll).ShouldNot(HaveOccurred())
})
