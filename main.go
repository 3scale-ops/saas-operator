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

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.

	_ "k8s.io/client-go/plugin/pkg/client/auth"

	marin3rv1alpha1 "github.com/3scale-ops/marin3r/apis/marin3r/v1alpha1"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/controllers"
	externalsecretsv1beta1 "github.com/3scale/saas-operator/pkg/apis/externalsecrets/v1beta1"
	grafanav1alpha1 "github.com/3scale/saas-operator/pkg/apis/grafana/v1alpha1"
	"github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	"github.com/3scale/saas-operator/pkg/reconcilers/threads"
	"github.com/3scale/saas-operator/pkg/reconcilers/workloads"
	"github.com/3scale/saas-operator/pkg/util"
	"github.com/3scale/saas-operator/pkg/version"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	apimachineryruntime "k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	// +kubebuilder:scaffold:imports
)

// Change below variables to serve metrics on different host or port.
const (
	// WatchNamespaceEnvVar is the constant for env variable WATCH_NAMESPACE
	// which specifies the Namespace to watch.
	// An empty value means the operator is running with cluster scope.
	watchNamespaceEnvVar string = "WATCH_NAMESPACE"
	debugLevelEnvVar     string = "DEBUG_LEVEL"
)

var (
	scheme   = apimachineryruntime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(saasv1alpha1.AddToScheme(scheme))
	utilruntime.Must(monitoringv1.AddToScheme(scheme))
	utilruntime.Must(grafanav1alpha1.AddToScheme(scheme))
	utilruntime.Must(externalsecretsv1beta1.AddToScheme(scheme))
	utilruntime.Must(marin3rv1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	ctrl.SetLogger((util.Logger{}).New())

	printVersion()

	if err := (&util.ProfilerConfig{
		Log: ctrl.Log.WithName("profiler"),
	}).Setup(); err != nil {
		setupLog.Error(err, "unable to start the Profiler")
	}

	watchNamespace, err := getWatchNamespace()
	if err != nil {
		setupLog.Error(err, "unable to get WatchNamespace, "+
			"the manager will watch and manage resources in all Namespaces")
	}

	options := ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "256bc75b.3scale.net",
		Namespace:              watchNamespace, // namespaced-scope when the value is not an empty string
	}

	if strings.Contains(watchNamespace, ",") {
		setupLog.Info(fmt.Sprintf("manager in MultiNamespaced mode will be watching namespaces %q", watchNamespace))
		options.NewCache = cache.MultiNamespacedCacheBuilder(strings.Split(watchNamespace, ","))
	} else if watchNamespace != "" {
		setupLog.Info(fmt.Sprintf("manager in Namespaced mode will be watching namespace %q", watchNamespace))
		options.Namespace = watchNamespace
	} else {
		setupLog.Info("manager in cluster mode will be watching all namespaces")
		options.Namespace = ""
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), options)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	/* BASERECONCILER_V2 BASED CONTROLLERS*/

	if err = (&controllers.SentinelReconciler{
		Reconciler:     basereconciler.NewFromManager(mgr),
		SentinelEvents: threads.NewManager(),
		Metrics:        threads.NewManager(),
		Log:            ctrl.Log.WithName("controllers").WithName("Sentinel"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Sentinel")
		os.Exit(1)
	}
	if err = (&controllers.RedisShardReconciler{
		Reconciler: basereconciler.NewFromManager(mgr),
		Log:        ctrl.Log.WithName("controllers").WithName("RedisShard"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "RedisShard")
		os.Exit(1)
	}
	if err = (&controllers.TwemproxyConfigReconciler{
		Reconciler:     basereconciler.NewFromManager(mgr),
		SentinelEvents: threads.NewManager(),
		Log:            ctrl.Log.WithName("controllers").WithName("TwemproxyConfig"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "TwemproxyConfig")
		os.Exit(1)
	}

	/* WORKLOADS RECONCILER BASED CONTROLLERS*/

	if err = (&controllers.ApicastReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr),
		Log:                ctrl.Log.WithName("controllers").WithName("Apicast"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Apicast")
		os.Exit(1)
	}

	if err = (&controllers.ZyncReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr),
		Log:                ctrl.Log.WithName("controllers").WithName("Zync"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Zync")
		os.Exit(1)
	}

	if err = (&controllers.MappingServiceReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr),
		Log:                ctrl.Log.WithName("controllers").WithName("MappingService"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "MappingService")
		os.Exit(1)
	}

	if err = (&controllers.CORSProxyReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr),
		Log:                ctrl.Log.WithName("controllers").WithName("CORSProxy"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CORSProxy")
		os.Exit(1)
	}

	if err = (&controllers.AutoSSLReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr),
		Log:                ctrl.Log.WithName("controllers").WithName("AutoSSL"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AutoSSL")
		os.Exit(1)
	}

	if err = (&controllers.EchoAPIReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr),
		Log:                ctrl.Log.WithName("controllers").WithName("EchoAPI"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "EchoAPI")
		os.Exit(1)
	}

	if err = (&controllers.BackendReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr),
		Log:                ctrl.Log.WithName("controllers").WithName("Backend"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Backend")
		os.Exit(1)
	}

	if err = (&controllers.SystemReconciler{
		WorkloadReconciler: workloads.NewFromManager(mgr),
		Log:                ctrl.Log.WithName("controllers").WithName("System"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "System")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

// getWatchNamespace returns the Namespace the operator should be watching for changes
func getWatchNamespace() (string, error) {

	ns, found := os.LookupEnv(watchNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", watchNamespaceEnvVar)
	}
	return ns, nil
}

func printVersion() {
	setupLog.Info(fmt.Sprintf("SaaS Operator Version: %s", version.Current()))
	setupLog.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	setupLog.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
}
