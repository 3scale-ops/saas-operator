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

	"github.com/3scale-ops/basereconciler/reconciler"
	"github.com/3scale-ops/basereconciler/util"
	marin3rv1alpha1 "github.com/3scale-ops/marin3r/apis/marin3r/v1alpha1"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/generators/backend"
	"github.com/3scale-ops/saas-operator/pkg/reconcilers/workloads"
	externalsecretsv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	"github.com/go-logr/logr"
	grafanav1alpha1 "github.com/grafana-operator/grafana-operator/v4/api/integreatly/v1alpha1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// BackendReconciler reconciles a Backend object
type BackendReconciler struct {
	workloads.WorkloadReconciler
	Log logr.Logger
}

// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=backends,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=backends/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=backends/finalizers,verbs=update
// +kubebuilder:rbac:groups="core",namespace=placeholder,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="core",namespace=placeholder,resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="apps",namespace=placeholder,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="monitoring.coreos.com",namespace=placeholder,resources=podmonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="autoscaling",namespace=placeholder,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="policy",namespace=placeholder,resources=poddisruptionbudgets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="integreatly.org",namespace=placeholder,resources=grafanadashboards,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="external-secrets.io",namespace=placeholder,resources=externalsecrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="marin3r.3scale.net",namespace=placeholder,resources=envoyconfigs,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *BackendReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("name", req.Name, "namespace", req.Namespace)
	ctx = log.IntoContext(ctx, logger)

	instance := &saasv1alpha1.Backend{}
	result := r.ManageResourceLifecycle(ctx, req, instance,
		reconciler.WithInMemoryInitializationFunc(util.ResourceDefaulter(instance)))
	if result.ShouldReturn() {
		return result.Values()
	}

	gen, err := backend.NewGenerator(instance.GetName(), instance.GetNamespace(), instance.Spec)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Shared resources
	resources := gen.Resources()

	// Listener resources
	listener_resources, err := r.NewDeploymentWorkload(&gen.Listener, gen.CanaryListener)
	if err != nil {
		return ctrl.Result{}, err
	}
	resources = append(resources, listener_resources...)

	// Worker resources
	worker_resources, err := r.NewDeploymentWorkload(&gen.Worker, gen.CanaryWorker)
	if err != nil {
		return ctrl.Result{}, err
	}
	resources = append(resources, worker_resources...)

	// Cron resources
	cron_resources, err := r.NewDeploymentWorkload(&gen.Cron, nil)
	if err != nil {
		return ctrl.Result{}, err
	}
	resources = append(resources, cron_resources...)

	// Reconcile all resources
	result = r.ReconcileOwnedResources(ctx, instance, resources)
	if result.ShouldReturn() {
		return result.Values()
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BackendReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&saasv1alpha1.Backend{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&policyv1.PodDisruptionBudget{}).
		Owns(&autoscalingv2.HorizontalPodAutoscaler{}).
		Owns(&monitoringv1.PodMonitor{}).
		Owns(&externalsecretsv1beta1.ExternalSecret{}).
		Owns(&grafanav1alpha1.GrafanaDashboard{}).
		Owns(&marin3rv1alpha1.EnvoyConfig{}).
		Watches(&source.Kind{Type: &corev1.Secret{TypeMeta: metav1.TypeMeta{Kind: "Secret"}}},
			r.FilteredEventHandler(&saasv1alpha1.BackendList{}, nil, r.Log)).
		Complete(r)
}
