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

	"github.com/go-logr/logr"
	"github.com/redhat-cop/operator-utils/pkg/util"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	"github.com/3scale/saas-operator/pkg/generators/apicast"
)

// ApicastReconciler reconciles a Apicast object
type ApicastReconciler struct {
	basereconciler.Reconciler
	Log logr.Logger
}

// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=apicasts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=apicasts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=apicasts/finalizers,verbs=update
// +kubebuilder:rbac:groups="core",namespace=placeholder,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="apps",namespace=placeholder,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="monitoring.coreos.com",namespace=placeholder,resources=podmonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="autoscaling",namespace=placeholder,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="policy",namespace=placeholder,resources=poddisruptionbudgets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="integreatly.org",namespace=placeholder,resources=grafanadashboards,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ApicastReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("name", req.Name, "namespace", req.Namespace)

	instance := &saasv1alpha1.Apicast{}
	key := types.NamespacedName{Name: req.Name, Namespace: req.Namespace}
	result, err := r.GetInstance(ctx, key, instance, saasv1alpha1.Finalizer, nil, log)
	if result != nil || err != nil {
		return *result, err
	}

	// Apply defaults for reconcile but do not store them in the API
	instance.Default()
	// json, _ := json.Marshal(instance.Spec)
	// log.V(1).Info("Apply defaults before resolving templates", "JSON", string(json))

	gen := apicast.NewGenerator(
		instance.GetName(),
		instance.GetNamespace(),
		instance.Spec,
	)

	err = r.ReconcileOwnedResources(ctx, instance, basereconciler.ControlledResources{
		Deployments: []basereconciler.Deployment{
			{
				Template: gen.Staging.Deployment(),
				HasHPA:   !instance.Spec.Staging.HPA.IsDeactivated(),
			},
			{
				Template: gen.Production.Deployment(),
				HasHPA:   !instance.Spec.Production.HPA.IsDeactivated(),
			},
		},
		SecretDefinitions: []basereconciler.SecretDefinition{},
		Services: []basereconciler.Service{
			{
				Template: gen.Staging.GatewayService(),
				Enabled:  true,
			},
			{
				Template: gen.Staging.MgmtService(),
				Enabled:  true,
			},
			{
				Template: gen.Production.GatewayService(),
				Enabled:  true,
			},
			{
				Template: gen.Production.MgmtService(),
				Enabled:  true,
			},
		},
		PodDisruptionBudgets: []basereconciler.PodDisruptionBudget{
			{
				Template: gen.Staging.PDB(),
				Enabled:  !instance.Spec.Staging.PDB.IsDeactivated(),
			},
			{
				Template: gen.Production.PDB(),
				Enabled:  !instance.Spec.Production.PDB.IsDeactivated(),
			},
		},
		HorizontalPodAutoscalers: []basereconciler.HorizontalPodAutoscaler{
			{
				Template: gen.Staging.HPA(),
				Enabled:  !instance.Spec.Staging.HPA.IsDeactivated(),
			},
			{
				Template: gen.Production.HPA(),
				Enabled:  !instance.Spec.Production.HPA.IsDeactivated(),
			},
		},
		PodMonitors: []basereconciler.PodMonitor{
			{
				Template: gen.Staging.PodMonitor(),
				Enabled:  true,
			},
			{
				Template: gen.Production.PodMonitor(),
				Enabled:  true,
			},
		},
		GrafanaDashboards: []basereconciler.GrafanaDashboard{
			{
				Template: gen.ApicastDashboard(),
				Enabled:  !instance.Spec.GrafanaDashboard.IsDeactivated(),
			},
			{
				Template: gen.ApicastServicesDashboard(),
				Enabled:  !instance.Spec.GrafanaDashboard.IsDeactivated(),
			},
		},
	})

	if err != nil {
		log.Error(err, "unable to update locked resources")
		return r.ManageError(ctx, instance, err)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApicastReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&saasv1alpha1.Apicast{}, builder.WithPredicates(util.ResourceGenerationOrFinalizerChangedPredicate{})).
		Watches(&source.Channel{Source: r.GetStatusChangeChannel()}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
