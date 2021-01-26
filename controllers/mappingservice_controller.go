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

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/service"
	"github.com/3scale/saas-operator/pkg/generators/mappingservice"
	"github.com/go-logr/logr"
	"github.com/redhat-cop/operator-utils/pkg/util"
	"github.com/redhat-cop/operator-utils/pkg/util/lockedresourcecontroller/lockedpatch"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// MappingServiceReconciler reconciles a MappingService object
type MappingServiceReconciler struct {
	basereconciler.Reconciler
	Log logr.Logger
}

// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=mappingservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=mappingservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=mappingservices/finalizers,verbs=update
// +kubebuilder:rbac:groups="core",namespace=placeholder,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="apps",namespace=placeholder,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="monitoring.coreos.com",namespace=placeholder,resources=podmonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="autoscaling",namespace=placeholder,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="policy",namespace=placeholder,resources=poddisruptionbudgets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="integreatly.org",namespace=placeholder,resources=grafanadashboards,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="secrets-manager.tuenti.io",namespace=placeholder,resources=secretdefinitions,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *MappingServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("name", req.Name, "namespace", req.Namespace)

	instance := &saasv1alpha1.MappingService{}
	key := types.NamespacedName{Name: req.Name, Namespace: req.Namespace}
	err := r.GetClient().Get(ctx, key, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if util.IsBeingDeleted(instance) {
		if !util.HasFinalizer(instance, saasv1alpha1.Finalizer) {
			return ctrl.Result{}, nil
		}
		err := r.ManageCleanUpLogic(instance, log)
		if err != nil {
			log.Error(err, "unable to delete instance")
			return r.ManageError(ctx, instance, err)
		}
		util.RemoveFinalizer(instance, saasv1alpha1.Finalizer)
		err = r.GetClient().Update(ctx, instance)
		if err != nil {
			log.Error(err, "unable to update instance")
			return r.ManageError(ctx, instance, err)
		}
		return ctrl.Result{}, nil
	}

	if ok := r.IsInitialized(instance, saasv1alpha1.Finalizer); !ok {
		err := r.GetClient().Update(ctx, instance)
		if err != nil {
			log.Error(err, "unable to initialize instance")
			return r.ManageError(ctx, instance, err)
		}
		return ctrl.Result{}, nil
	}

	// Apply defaults for reconcile but do not store them in the API
	instance.Default()
	json, _ := json.Marshal(instance.Spec)
	log.V(1).Info("Apply defaults before resolving templates", "JSON", string(json))

	gen := mappingservice.NewGenerator(
		instance.GetName(),
		instance.GetNamespace(),
		instance.Spec,
	)

	// Calculate resources to enforce
	resources := []basereconciler.LockedResource{}

	resources = append(resources,
		basereconciler.LockedResource{
			GeneratorFn: gen.Deployment(),
			ExcludePaths: func() []string {
				if instance.Spec.HPA.IsDeactivated() {
					return basereconciler.DefaultExcludedPaths
				}
				return append(basereconciler.DefaultExcludedPaths, "/spec/replicas")
			}(),
		})

	resources = append(resources,
		basereconciler.LockedResource{
			GeneratorFn:  gen.SecretDefinition(),
			ExcludePaths: basereconciler.DefaultExcludedPaths,
		})

	resources = append(resources,
		basereconciler.LockedResource{
			GeneratorFn:  gen.Service(),
			ExcludePaths: append(basereconciler.DefaultExcludedPaths, service.Excludes(gen.Service())...),
		})

	resources = append(resources,
		basereconciler.LockedResource{
			GeneratorFn:  gen.PodMonitor(),
			ExcludePaths: basereconciler.DefaultExcludedPaths,
		})

	if !instance.Spec.HPA.IsDeactivated() {
		resources = append(resources,
			basereconciler.LockedResource{
				GeneratorFn:  gen.HPA(),
				ExcludePaths: basereconciler.DefaultExcludedPaths,
			},
		)
	}

	if !instance.Spec.PDB.IsDeactivated() {
		resources = append(resources,
			basereconciler.LockedResource{
				GeneratorFn:  gen.PDB(),
				ExcludePaths: basereconciler.DefaultExcludedPaths,
			},
		)
	}

	if !instance.Spec.GrafanaDashboard.IsDeactivated() {
		resources = append(resources,
			basereconciler.LockedResource{
				GeneratorFn:  gen.GrafanaDashboard(),
				ExcludePaths: basereconciler.DefaultExcludedPaths,
			},
		)
	}

	lockedResources, err := r.NewLockedResources(resources, instance)
	err = r.UpdateLockedResources(ctx, instance, lockedResources, []lockedpatch.LockedPatch{})
	if err != nil {
		log.Error(err, "unable to update locked resources")
		return r.ManageError(ctx, instance, err)
	}

	return r.ManageSuccess(ctx, instance)
}

// SetupWithManager sets up the controller with the Manager.
func (r *MappingServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&saasv1alpha1.MappingService{}).
		Watches(&source.Channel{Source: r.GetStatusChangeChannel()}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}