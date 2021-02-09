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

// +kubebuilder:rbac:groups=saas.3scale.net,resources=systems,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=saas.3scale.net,resources=systems/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=saas.3scale.net,resources=systems/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the System object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
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

	err = r.ReconcileOwnedResources(ctx, instance, basereconciler.ControlledResources{
		Deployments: []basereconciler.Deployment{
			{
				Template:        gen.App.Deployment(),
				RolloutTriggers: nil,
				HasHPA:          !instance.Spec.App.HPA.IsDeactivated(),
			},
		},
		SecretDefinitions: []basereconciler.SecretDefinition{
			{Template: gen.ConfigFilesSecretDefinition(), Enabled: true},
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
		// Services: []basereconciler.Service{{
		// 	Template: gen.Service(),
		// 	Enabled:  true,
		// }},
		// PodDisruptionBudgets: []basereconciler.PodDisruptionBudget{{
		// 	Template: gen.PDB(),
		// 	Enabled:  !instance.Spec.PDB.IsDeactivated(),
		// }},
		// HorizontalPodAutoscalers: []basereconciler.HorizontalPodAutoscaler{{
		// 	Template: gen.HPA(),
		// 	Enabled:  !instance.Spec.HPA.IsDeactivated(),
		// }},
		// PodMonitors: []basereconciler.PodMonitor{{
		// 	Template: gen.PodMonitor(),
		// 	Enabled:  true,
		// }},
		// GrafanaDashboards: []basereconciler.GrafanaDashboard{{
		// 	Template: gen.GrafanaDashboard(),
		// 	Enabled:  !instance.Spec.GrafanaDashboard.IsDeactivated(),
		// }},
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
