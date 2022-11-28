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

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	externalsecretsv1beta1 "github.com/3scale/saas-operator/pkg/apis/externalsecrets/v1beta1"
	grafanav1alpha1 "github.com/3scale/saas-operator/pkg/apis/grafana/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators/system"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	"github.com/3scale/saas-operator/pkg/reconcilers/workloads"
	"github.com/go-logr/logr"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// SystemReconciler reconciles a System object
type SystemReconciler struct {
	workloads.WorkloadReconciler
	Log logr.Logger
}

// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=systems,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=systems/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=systems/finalizers,verbs=update
// +kubebuilder:rbac:groups="core",namespace=placeholder,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="apps",namespace=placeholder,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="apps",namespace=placeholder,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="monitoring.coreos.com",namespace=placeholder,resources=podmonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="autoscaling",namespace=placeholder,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="policy",namespace=placeholder,resources=poddisruptionbudgets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="integreatly.org",namespace=placeholder,resources=grafanadashboards,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="external-secrets.io",namespace=placeholder,resources=externalsecrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *SystemReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("name", req.Name, "namespace", req.Namespace)
	ctx = log.IntoContext(ctx, logger)

	instance := &saasv1alpha1.System{}
	key := types.NamespacedName{Name: req.Name, Namespace: req.Namespace}
	result, err := r.GetInstance(ctx, key, instance, saasv1alpha1.Finalizer, nil)
	if result != nil || err != nil {
		return *result, err
	}

	// Apply defaults for reconcile but do not store them in the API
	instance.Default()

	gen, err := system.NewGenerator(
		instance.GetName(),
		instance.GetNamespace(),
		instance.Spec,
	)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Shared resources
	resources := append(gen.ExternalSecrets(), gen.GrafanaDashboard())

	// System APP
	var app_resources []basereconciler.Resource
	if instance.Spec.App.Canary != nil {
		app_resources, err = r.NewDeploymentWorkloadWithTraffic(ctx, instance, &gen.App, &gen.App, gen.CanaryApp)
	} else {
		app_resources, err = r.NewDeploymentWorkloadWithTraffic(ctx, instance, &gen.App, &gen.App)
	}
	if err != nil {
		return ctrl.Result{}, err
	}

	resources = append(resources, app_resources...)

	// Sidekiq Default resources
	var sidekiq_default_resources []basereconciler.Resource
	if instance.Spec.SidekiqDefault.Canary != nil {
		sidekiq_default_resources, err = r.NewDeploymentWorkload(ctx, instance, &gen.SidekiqDefault, gen.CanarySidekiqDefault)
	} else {
		sidekiq_default_resources, err = r.NewDeploymentWorkload(ctx, instance, &gen.SidekiqDefault)
	}
	if err != nil {
		return ctrl.Result{}, err
	}
	resources = append(resources, sidekiq_default_resources...)

	// Sidekiq Billing resources
	var sidekiq_billing_resources []basereconciler.Resource
	if instance.Spec.SidekiqBilling.Canary != nil {
		sidekiq_billing_resources, err = r.NewDeploymentWorkload(ctx, instance, &gen.SidekiqBilling, gen.CanarySidekiqBilling)
	} else {
		sidekiq_billing_resources, err = r.NewDeploymentWorkload(ctx, instance, &gen.SidekiqBilling)
	}
	if err != nil {
		return ctrl.Result{}, err
	}
	resources = append(resources, sidekiq_billing_resources...)

	// Sidekiq Low resources
	var sidekiq_low_resources []basereconciler.Resource
	if instance.Spec.SidekiqLow.Canary != nil {
		sidekiq_low_resources, err = r.NewDeploymentWorkload(ctx, instance, &gen.SidekiqLow, gen.CanarySidekiqLow)
	} else {
		sidekiq_low_resources, err = r.NewDeploymentWorkload(ctx, instance, &gen.SidekiqLow)
	}
	if err != nil {
		return ctrl.Result{}, err
	}
	resources = append(resources, sidekiq_low_resources...)

	// Sphinx resources
	resources = append(resources, gen.Sphinx.StatefulSetWithTraffic()...)

	// Console resources
	resources = append(resources, gen.Console.StatefulSet())

	err = r.ReconcileOwnedResources(ctx, instance, resources)
	if err != nil {
		logger.Error(err, "unable to update owned resources")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SystemReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&saasv1alpha1.System{}).
		Owns(&appsv1.Deployment{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Owns(&policyv1.PodDisruptionBudget{}).
		Owns(&autoscalingv2beta2.HorizontalPodAutoscaler{}).
		Owns(&monitoringv1.PodMonitor{}).
		Owns(&externalsecretsv1beta1.ExternalSecret{}).
		Owns(&grafanav1alpha1.GrafanaDashboard{}).
		Watches(&source.Kind{Type: &corev1.Secret{TypeMeta: metav1.TypeMeta{Kind: "Secret"}}},
			r.SecretEventHandler(&saasv1alpha1.SystemList{}, r.Log)).
		Complete(r)
}
