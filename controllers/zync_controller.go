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
	"github.com/3scale/saas-operator/pkg/generators/zync"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v1"
	"github.com/go-logr/logr"
	"github.com/redhat-cop/operator-utils/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// ZyncReconciler reconciles a Zync object
type ZyncReconciler struct {
	basereconciler.Reconciler
	Log logr.Logger
}

// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=zyncs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=zyncs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=zyncs/finalizers,verbs=update
// +kubebuilder:rbac:groups="core",namespace=placeholder,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="core",namespace=placeholder,resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="apps",namespace=placeholder,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="monitoring.coreos.com",namespace=placeholder,resources=podmonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="autoscaling",namespace=placeholder,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="policy",namespace=placeholder,resources=poddisruptionbudgets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="integreatly.org",namespace=placeholder,resources=grafanadashboards,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="secrets-manager.tuenti.io",namespace=placeholder,resources=secretdefinitions,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ZyncReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("name", req.Name, "namespace", req.Namespace)

	instance := &saasv1alpha1.Zync{}
	key := types.NamespacedName{Name: req.Name, Namespace: req.Namespace}
	result, err := r.GetInstance(ctx, key, instance, saasv1alpha1.Finalizer, nil, log)
	if result != nil || err != nil {
		return *result, err
	}

	// Apply defaults for reconcile but do not store them in the API
	instance.Default()

	gen := zync.NewGenerator(
		instance.GetName(),
		instance.GetNamespace(),
		instance.Spec,
	)

	// Calculate rollout triggers
	triggers, err := r.TriggersFromSecretDefs(ctx,
		gen.ZyncSecretDefinition(),
	)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.ReconcileOwnedResources(ctx, instance, basereconciler.ControlledResources{
		Deployments: []basereconciler.Deployment{
			{
				Template:        gen.API.Deployment(),
				HasHPA:          !instance.Spec.API.HPA.IsDeactivated(),
				RolloutTriggers: triggers,
			},
			{
				Template:        gen.Que.Deployment(),
				HasHPA:          !instance.Spec.Que.HPA.IsDeactivated(),
				RolloutTriggers: triggers,
			},
		},
		SecretDefinitions: []basereconciler.SecretDefinition{
			{
				Template: gen.ZyncSecretDefinition(),
				Enabled:  true,
			},
		},
		Services: []basereconciler.Service{
			{
				Template: gen.API.Service(),
				Enabled:  true,
			}},
		PodDisruptionBudgets: []basereconciler.PodDisruptionBudget{
			{
				Template: gen.API.PDB(),
				Enabled:  !instance.Spec.API.PDB.IsDeactivated(),
			},
			{
				Template: gen.Que.PDB(),
				Enabled:  !instance.Spec.Que.PDB.IsDeactivated(),
			},
		},
		HorizontalPodAutoscalers: []basereconciler.HorizontalPodAutoscaler{
			{
				Template: gen.API.HPA(),
				Enabled:  !instance.Spec.API.HPA.IsDeactivated(),
			},
			{
				Template: gen.Que.HPA(),
				Enabled:  !instance.Spec.Que.HPA.IsDeactivated(),
			},
		},
		PodMonitors: []basereconciler.PodMonitor{
			{
				Template: gen.API.PodMonitor(),
				Enabled:  true,
			},
			{
				Template: gen.Que.PodMonitor(),
				Enabled:  true,
			},
		},
		GrafanaDashboards: []basereconciler.GrafanaDashboard{
			{
				Template: gen.GrafanaDashboard(),
				Enabled:  !instance.Spec.GrafanaDashboard.IsDeactivated(),
			},
		},
	})

	if err != nil {
		log.Error(err, "unable to reconcile owned resources")
		return r.ManageError(ctx, instance, err)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ZyncReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&saasv1alpha1.Zync{}, builder.WithPredicates(util.ResourceGenerationOrFinalizerChangedPredicate{})).
		Watches(&source.Channel{Source: r.GetStatusChangeChannel()}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &corev1.Secret{TypeMeta: metav1.TypeMeta{Kind: "Secret"}}},
			r.SecretEventHandler(&saasv1alpha1.ZyncList{}, r.Log)).
		Complete(r)
}
