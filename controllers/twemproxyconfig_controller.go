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
	"time"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators/twemproxyconfig"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	"github.com/3scale/saas-operator/pkg/reconcilers/threads"
	"github.com/3scale/saas-operator/pkg/redis/events"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// TwemproxyConfigReconciler reconciles a TwemproxyConfig object
type TwemproxyConfigReconciler struct {
	basereconciler.Reconciler
	Log            logr.Logger
	SentinelEvents threads.Manager
}

// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=twemproxyconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=twemproxyconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=twemproxyconfigs/finalizers,verbs=update
// +kubebuilder:rbac:groups="core",namespace=placeholder,resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *TwemproxyConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("name", req.Name, "namespace", req.Namespace)

	instance := &saasv1alpha1.TwemproxyConfig{}
	key := types.NamespacedName{Name: req.Name, Namespace: req.Namespace}
	result, err := r.GetInstance(ctx, key, instance, saasv1alpha1.Finalizer,
		[]func(){r.SentinelEvents.CleanupThreads(instance)}, log)
	if result != nil || err != nil {
		return *result, err
	}

	gen, err := twemproxyconfig.NewGenerator(ctx, instance)
	if err != nil {
		return r.ManageError(ctx, instance, err)
	}

	if err := r.ReconcileOwnedResources(ctx, instance, r.GetScheme(), gen.Resources()); err != nil {
		log.Error(err, "unable to update owned resources")
		return r.ManageError(ctx, instance, err)
	}

	// Reconcile sentinel the event watchers
	eventWatchers := make([]threads.RunnableThread, 0, len(gen.Spec.SentinelURIs))
	for _, uri := range gen.Spec.SentinelURIs {
		eventWatchers = append(eventWatchers, &events.SentinelEventWatcher{
			Instance:      instance,
			SentinelURI:   uri,
			ExportMetrics: false,
		})
	}
	r.SentinelEvents.ReconcileThreads(ctx, instance, eventWatchers, log.WithName("event-watcher"))

	// Reconcile periodically in case some event is lost ...
	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TwemproxyConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&saasv1alpha1.TwemproxyConfig{}).
		Watches(&source.Channel{Source: r.GetStatusChangeChannel()}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Channel{Source: r.SentinelEvents.GetChannel()}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
