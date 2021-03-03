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
	"encoding/json"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	"github.com/3scale/saas-operator/pkg/generators/system"
)

// SystemReconciler reconciles a System object
type SystemReconciler struct {
	basereconciler.Reconciler
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

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *SystemReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("name", req.Name, "namespace", req.Namespace)

	instance := &saasv1alpha1.System{}
	key := types.NamespacedName{Name: req.Name, Namespace: req.Namespace}
	result, err := r.GetInstance(ctx, key, instance, saasv1alpha1.Finalizer, log)
	if result != nil || err != nil {
		return *result, err
	}

	// Apply defaults for reconcile but do not store them in the API
	instance.Default()
	json, _ := json.Marshal(instance.Spec)
	log.V(1).Info("Apply defaults before resolving templates", "JSON", string(json))

	gen := system.NewGenerator(
		instance.GetName(),
		instance.GetNamespace(),
		instance.Spec,
	)

	// Calculate rollout triggers (app & sidekiq)
	triggers, err := r.TriggersFromSecretDefs(ctx,
		gen.ConfigFilesSecretDefinition(),
		gen.SeedSecretDefinition(),
		gen.DatabaseSecretDefinition(),
		gen.RecaptchaSecretDefinition(),
		gen.EventsHookSecretDefinition(),
		gen.SMTPSecretDefinition(),
		gen.MasterApicastSecretDefinition(),
		gen.ZyncSecretDefinition(),
		gen.BackendSecretDefinition(),
		gen.MultitenantAssetsSecretDefinition(),
		gen.AppSecretDefinition(),
	)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Calculate rollout triggers (sphinx)
	sphinxTriggers, err := r.TriggersFromSecretDefs(ctx,
		gen.DatabaseSecretDefinition(),
	)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.ReconcileOwnedResources(ctx, instance, basereconciler.ControlledResources{
		Deployments: []basereconciler.Deployment{
			{
				Template:        gen.App.Deployment(),
				RolloutTriggers: triggers,
				HasHPA:          !instance.Spec.App.HPA.IsDeactivated(),
			},
			{
				Template:        gen.Sidekiq.Deployment(),
				RolloutTriggers: triggers,
				HasHPA:          !instance.Spec.Sidekiq.HPA.IsDeactivated(),
			},
		},
		StatefulSets: []basereconciler.StatefulSet{{
			Template:        gen.Sphinx.StatefulSet(),
			RolloutTriggers: sphinxTriggers,
			Enabled:         true,
		}},
		SecretDefinitions: []basereconciler.SecretDefinition{
			{Template: gen.ConfigFilesSecretDefinition(), Enabled: instance.Spec.Config.ConfigFiles.Enabled()},
			{Template: gen.SeedSecretDefinition(), Enabled: true},
			{Template: gen.DatabaseSecretDefinition(), Enabled: true},
			{Template: gen.RecaptchaSecretDefinition(), Enabled: true},
			{Template: gen.EventsHookSecretDefinition(), Enabled: true},
			{Template: gen.SMTPSecretDefinition(), Enabled: true},
			{Template: gen.MasterApicastSecretDefinition(), Enabled: true},
			{Template: gen.ZyncSecretDefinition(), Enabled: true},
			{Template: gen.BackendSecretDefinition(), Enabled: true},
			{Template: gen.MultitenantAssetsSecretDefinition(), Enabled: true},
			{Template: gen.AppSecretDefinition(), Enabled: true},
		},
		Services: []basereconciler.Service{
			{Template: gen.App.Service(), Enabled: true},
			{Template: gen.Sphinx.Service(), Enabled: true},
		},
		PodDisruptionBudgets: []basereconciler.PodDisruptionBudget{
			{Template: gen.App.PDB(), Enabled: !instance.Spec.App.PDB.IsDeactivated()},
			{Template: gen.Sidekiq.PDB(), Enabled: !instance.Spec.Sidekiq.PDB.IsDeactivated()},
		},
		HorizontalPodAutoscalers: []basereconciler.HorizontalPodAutoscaler{
			{Template: gen.App.HPA(), Enabled: !instance.Spec.App.HPA.IsDeactivated()},
			{Template: gen.Sidekiq.HPA(), Enabled: !instance.Spec.Sidekiq.HPA.IsDeactivated()},
		},
		PodMonitors: []basereconciler.PodMonitor{
			{Template: gen.App.PodMonitor(), Enabled: true},
			{Template: gen.Sidekiq.PodMonitor(), Enabled: true},
		},
		GrafanaDashboards: []basereconciler.GrafanaDashboard{
			{Template: gen.GrafanaDashboard(), Enabled: !instance.Spec.GrafanaDashboard.IsDeactivated()},
		},
	})

	if err != nil {
		log.Error(err, "unable to update owned resources")
		return r.ManageError(ctx, instance, err)
	}

	return r.ManageSuccess(ctx, instance)
}

// SetupWithManager sets up the controller with the Manager.
func (r *SystemReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&saasv1alpha1.System{}).
		Watches(&source.Channel{Source: r.GetStatusChangeChannel()}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &corev1.Secret{TypeMeta: metav1.TypeMeta{Kind: "Secret"}}},
			r.SecretEventHandler(&saasv1alpha1.SystemList{}, r.Log)).
		Complete(r)
}
