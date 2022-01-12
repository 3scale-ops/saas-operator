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
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators/sentinel"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	"github.com/3scale/saas-operator/pkg/reconcilers/threads"
	"github.com/3scale/saas-operator/pkg/redis"
	"github.com/3scale/saas-operator/pkg/redis/events"
	"github.com/3scale/saas-operator/pkg/redis/metrics"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// SentinelReconciler reconciles a Sentinel object
type SentinelReconciler struct {
	basereconciler.Reconciler
	Log            logr.Logger
	SentinelEvents threads.Manager
	Metrics        threads.Manager
}

// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=sentinels,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=sentinels/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=sentinels/finalizers,verbs=update
// +kubebuilder:rbac:groups="core",namespace=placeholder,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="apps",namespace=placeholder,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="monitoring.coreos.com",namespace=placeholder,resources=podmonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="policy",namespace=placeholder,resources=poddisruptionbudgets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="integreatly.org",namespace=placeholder,resources=grafanadashboards,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *SentinelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("name", req.Name, "namespace", req.Namespace)

	instance := &saasv1alpha1.Sentinel{}
	key := types.NamespacedName{Name: req.Name, Namespace: req.Namespace}
	result, err := r.GetInstance(ctx,
		key,
		instance,
		saasv1alpha1.Finalizer,
		[]func(){r.SentinelEvents.CleanupThreads(instance), r.Metrics.CleanupThreads(instance)},
		log)
	if result != nil || err != nil {
		return *result, err
	}

	// Apply defaults for reconcile but do not store them in the API
	instance.Default()

	gen := sentinel.NewGenerator(
		instance.GetName(),
		instance.GetNamespace(),
		instance.Spec,
	)

	if err := r.ReconcileOwnedResources(ctx, instance, r.GetScheme(), gen.Resources()); err != nil {
		log.Error(err, "unable to update owned resources")
		return r.ManageError(ctx, instance, err)
	}

	// Create the redis-sentinel server pool
	sentinelPool, err := redis.NewSentinelPool(ctx, r.GetClient(),
		types.NamespacedName{Name: gen.GetComponent(), Namespace: gen.GetNamespace()}, int(*instance.Spec.Replicas))
	if err != nil {
		return r.ManageError(ctx, instance, err)
	}

	// Create the ShardedCluster objects that represents the redis servers to be monitored by sentinel
	shardedCluster, err := redis.NewShardedCluster(ctx, instance.Spec.Config.MonitoredShards, log)
	if err != nil {
		return r.ManageError(ctx, instance, err)
	}

	// Ensure all shards are being monitored
	allMonitored, err := sentinelPool.IsMonitoringShards(ctx, shardedCluster.GetShardNames())
	if err != nil {
		return r.ManageError(ctx, instance, err)
	}
	if !allMonitored {
		if err := shardedCluster.Discover(ctx, log); err != nil {
			return r.ManageError(ctx, instance, err)
		}
		if _, err := sentinelPool.Monitor(ctx, shardedCluster); err != nil {
			return r.ManageError(ctx, instance, err)
		}
	}

	// Reconcile sentinel the event watchers and metrics gatherers
	eventWatchers := make([]threads.RunnableThread, 0, len(gen.SentinelURIs()))
	metricsGatherers := make([]threads.RunnableThread, 0, len(gen.SentinelURIs()))
	for _, uri := range gen.SentinelURIs() {
		eventWatchers = append(eventWatchers, &events.SentinelEventWatcher{
			Instance:      instance,
			SentinelURI:   uri,
			ExportMetrics: true,
		})
		metricsGatherers = append(metricsGatherers, &metrics.SentinelMetricsGatherer{
			RefreshInterval: *gen.Spec.Config.MetricsRefreshInterval,
			SentinelURI:     uri,
		})
	}
	if err := r.SentinelEvents.ReconcileThreads(ctx, instance, eventWatchers, log.WithName("event-watcher")); err != nil {
		return r.ManageError(ctx, instance, err)
	}
	if err := r.Metrics.ReconcileThreads(ctx, instance, metricsGatherers, log.WithName("metrics-gatherer")); err != nil {
		return r.ManageError(ctx, instance, err)
	}

	// Reconcile status of the Sentinel resource
	if err := r.reconcileStatus(ctx, instance, &gen, log); err != nil {
		return r.ManageError(ctx, instance, err)
	}

	return r.ManageSuccess(ctx, instance)
}

func (r *SentinelReconciler) reconcileStatus(ctx context.Context, instance *saasv1alpha1.Sentinel, gen *sentinel.Generator, log logr.Logger) error {

	sentinel, err := redis.NewSentinelServerFromConnectionString("sentinel", gen.SentinelServiceEndpoint())
	if err != nil {
		return err
	}

	monitoredShards, err := sentinel.MonitoredShards(ctx)
	if err != nil {
		return err
	}

	replicas := int(*gen.Spec.Replicas)
	addressList := make([]string, 0, replicas)

	for i := 0; i < replicas; i++ {
		key := types.NamespacedName{Name: gen.PodServiceName(i), Namespace: instance.GetNamespace()}
		svc := &corev1.Service{}
		if err := r.GetClient().Get(ctx, key, svc); err != nil {
			return err
		}
		addressList = append(addressList, fmt.Sprintf("%s:%d", svc.Spec.ClusterIP, saasv1alpha1.SentinelPort))
	}

	status := saasv1alpha1.SentinelStatus{
		Sentinels:       addressList,
		MonitoredShards: monitoredShards,
	}

	if !equality.Semantic.DeepEqual(status, instance.Status) {
		instance.Status = status
		if err := r.GetClient().Status().Update(ctx, instance); err != nil {
			return err
		}
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SentinelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&saasv1alpha1.Sentinel{}).
		Watches(&source.Channel{Source: r.GetStatusChangeChannel()}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Channel{Source: r.SentinelEvents.GetChannel()}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
