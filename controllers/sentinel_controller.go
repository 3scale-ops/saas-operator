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
	"errors"
	"time"

	basereconciler "github.com/3scale-ops/basereconciler/reconciler"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators/sentinel"
	"github.com/3scale/saas-operator/pkg/reconcilers/threads"
	"github.com/3scale/saas-operator/pkg/redis/events"
	"github.com/3scale/saas-operator/pkg/redis/metrics"
	redis "github.com/3scale/saas-operator/pkg/redis/server"
	"github.com/3scale/saas-operator/pkg/redis/sharded"
	"github.com/go-logr/logr"
	grafanav1alpha1 "github.com/grafana-operator/grafana-operator/v4/api/integreatly/v1alpha1"
	"golang.org/x/time/rate"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/ratelimiter"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// SentinelReconciler reconciles a Sentinel object
type SentinelReconciler struct {
	basereconciler.Reconciler
	Log            logr.Logger
	SentinelEvents threads.Manager
	Metrics        threads.Manager
	Pool           *redis.ServerPool
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
	logger := r.Log.WithValues("name", req.Name, "namespace", req.Namespace)
	ctx = log.IntoContext(ctx, logger)

	instance := &saasv1alpha1.Sentinel{}
	key := types.NamespacedName{Name: req.Name, Namespace: req.Namespace}
	result, err := r.GetInstance(ctx,
		key,
		instance,
		pointer.String(saasv1alpha1.Finalizer),
		[]func(){r.SentinelEvents.CleanupThreads(instance), r.Metrics.CleanupThreads(instance)})
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

	if err := r.ReconcileOwnedResources(ctx, instance, gen.Resources()); err != nil {
		logger.Error(err, "unable to update owned resources")
		return ctrl.Result{}, err
	}

	clustermap, err := gen.ClusterTopology(ctx)
	if err != nil {
		return ctrl.Result{}, err
	}
	shardedCluster, err := sharded.NewShardedCluster(ctx, clustermap, r.Pool)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Ensure all shards are being monitored
	for _, sentinel := range shardedCluster.Sentinels {
		allMonitored, err := sentinel.IsMonitoringShards(ctx, shardedCluster.GetShardNames())
		if err != nil {
			return ctrl.Result{}, err
		}
		if !allMonitored {
			if err := shardedCluster.Discover(ctx); err != nil {
				return ctrl.Result{}, err
			}
			if _, err := sentinel.Monitor(ctx, shardedCluster); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	// Reconcile sentinel the event watchers and metrics gatherers
	eventWatchers := make([]threads.RunnableThread, 0, len(gen.SentinelURIs()))
	metricsGatherers := make([]threads.RunnableThread, 0, len(gen.SentinelURIs()))
	for _, uri := range gen.SentinelURIs() {
		watcher, err := events.NewSentinelEventWatcher(uri, instance, shardedCluster, true, r.Pool)
		if err != nil {
			return ctrl.Result{}, err
		}
		gatherer, err := metrics.NewSentinelMetricsGatherer(uri, *gen.Spec.Config.MetricsRefreshInterval, r.Pool)
		if err != nil {
			return ctrl.Result{}, err
		}
		eventWatchers = append(eventWatchers, watcher)
		metricsGatherers = append(metricsGatherers, gatherer)
	}
	if err := r.SentinelEvents.ReconcileThreads(ctx, instance, eventWatchers, logger.WithName("event-watcher")); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.Metrics.ReconcileThreads(ctx, instance, metricsGatherers, logger.WithName("metrics-gatherer")); err != nil {
		return ctrl.Result{}, err
	}

	// Reconcile status of the Sentinel resource
	if err := r.reconcileStatus(ctx, instance, shardedCluster, logger); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

func (r *SentinelReconciler) reconcileStatus(ctx context.Context, instance *saasv1alpha1.Sentinel, cluster *sharded.Cluster,
	log logr.Logger) error {

	// sentinels info to the status
	sentinels := make([]string, len(cluster.Sentinels))
	for idx, srv := range cluster.Sentinels {
		sentinels[idx] = srv.ID()
	}

	// redis shards info to the status
	merr := cluster.SentinelDiscover(ctx, sharded.SlaveReadOnlyDiscoveryOpt, sharded.SaveConfigDiscoveryOpt)
	// if the failure occurred calling sentinel discard the result and return error
	// otherwise keep going on and use the information that was returned, even if there were some
	// other errors
	sentinelError := &sharded.DiscoveryError_Sentinel_Failure{}
	if errors.As(merr, sentinelError) {
		return merr
	}
	// We don't want the controller to keep failing while things reconfigure as
	// this makes controller throttling to kick in. Instead, just log the errors
	// and rely on reconciles triggered by sentinel events to correct the situation.
	masterError := &sharded.DiscoveryError_Master_SingleServerFailure{}
	slaveError := &sharded.DiscoveryError_Slave_SingleServerFailure{}
	if errors.As(merr, masterError) || errors.As(merr, slaveError) {
		log.Error(merr, "DiscoveryError")
	}

	shards := make(saasv1alpha1.MonitoredShards, len(cluster.Shards))
	for idx, shard := range cluster.Shards {
		shards[idx] = saasv1alpha1.MonitoredShard{
			Name:    shard.Name,
			Servers: make(map[string]saasv1alpha1.RedisServerDetails, len(shard.Servers)),
		}
		for _, srv := range shard.Servers {
			shards[idx].Servers[srv.GetAlias()] = saasv1alpha1.RedisServerDetails{
				Role:    srv.Role,
				Address: srv.ID(),
				Config:  srv.Config,
			}
		}
	}

	status := saasv1alpha1.SentinelStatus{
		Sentinels:       sentinels,
		MonitoredShards: shards,
	}

	if !equality.Semantic.DeepEqual(status, instance.Status) {
		instance.Status = status
		if err := r.Client.Status().Update(ctx, instance); err != nil {
			return err
		}
		log.Info("status updated")
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SentinelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&saasv1alpha1.Sentinel{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Owns(&policyv1.PodDisruptionBudget{}).
		Owns(&grafanav1alpha1.GrafanaDashboard{}).
		Owns(&corev1.ConfigMap{}).
		Watches(&source.Channel{Source: r.SentinelEvents.GetChannel()}, &handler.EnqueueRequestForObject{}).
		WithOptions(controller.Options{
			RateLimiter: AggressiveRateLimiter(),
		}).
		Complete(r)
}

func AggressiveRateLimiter() ratelimiter.RateLimiter {
	// return workqueue.DefaultControllerRateLimiter()
	return workqueue.NewMaxOfRateLimiter(
		// First retries are more spaced that default
		// Max retry time is limited to 10 seconds
		workqueue.NewItemExponentialFailureRateLimiter(5*time.Millisecond, 10*time.Second),
		// 10 qps, 100 bucket size.  This is only for retry speed and its only the overall factor (not per item)
		&workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(10), 100)},
	)
}
