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
	"fmt"
	"strings"
	"sync"
	"time"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	"github.com/3scale/saas-operator/pkg/generators/sentinel"
	"github.com/3scale/saas-operator/pkg/redis"
	redismetrics "github.com/3scale/saas-operator/pkg/redis/metrics"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// SentinelMetrics holds a map of SentinelMetricsGatherer to keep track
// of which sentinel pods already have a running exporter
type SentinelMetrics struct {
	mu        sync.Mutex
	exporters map[string]*redismetrics.SentinelMetricsGatherer
}

// RunExporter runs a SentinelMetricsGatherer for the given key. The key should uniquely identify a Pod's exporter.
func (sm *SentinelMetrics) RunExporter(ctx context.Context, key string, sentinelURL string,
	refreshInterval time.Duration, log logr.Logger) {

	// run the exporter for this instance if it is not running, do nothing otherwise
	if _, ok := sm.exporters[key]; !ok {
		sm.mu.Lock()
		sm.exporters[key] = &redismetrics.SentinelMetricsGatherer{
			RefreshInterval: refreshInterval,
			SentinelURL:     sentinelURL,
			Log:             log,
		}
		sm.exporters[key].Start(ctx)
		sm.mu.Unlock()
	}
}

// StopExporter stops the exporter for the given key
func (sm *SentinelMetrics) StopExporter(key string) {
	sm.mu.Lock()
	sm.exporters[key].Stop()
	delete(sm.exporters, key)
	sm.mu.Unlock()
}

// SentinelReconciler reconciles a Sentinel object
type SentinelReconciler struct {
	basereconciler.Reconciler
	Log     logr.Logger
	metrics SentinelMetrics
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
	result, err := r.GetInstance(ctx, key, instance, saasv1alpha1.Finalizer, []func(){r.CleanupExporters(req.NamespacedName)}, log)
	if result != nil || err != nil {
		return *result, err
	}

	// Apply defaults for reconcile but do not store them in the API
	instance.Default()
	json, _ := json.Marshal(instance.Spec)
	log.V(1).Info("Apply defaults before resolving templates", "JSON", string(json))

	gen := sentinel.NewGenerator(
		instance.GetName(),
		instance.GetNamespace(),
		instance.Spec,
	)

	if err := r.ReconcileOwnedResources(ctx, instance, basereconciler.ControlledResources{
		StatefulSets: []basereconciler.StatefulSet{{
			Template:        gen.StatefulSet(),
			RolloutTriggers: nil,
			Enabled:         true,
		}},
		ConfigMaps: []basereconciler.ConfigMaps{{
			Template: gen.ConfigMap(),
			Enabled:  true,
		}},
		Services: func() []basereconciler.Service {
			fns := append(gen.PodServices(), gen.StatefulSetService())
			svcs := make([]basereconciler.Service, 0, len(fns))
			for _, fn := range fns {
				svcs = append(svcs, basereconciler.Service{Template: fn, Enabled: true})
			}
			return svcs
		}(),
		PodDisruptionBudgets: []basereconciler.PodDisruptionBudget{{
			Template: gen.PDB(),
			Enabled:  !instance.Spec.PDB.IsDeactivated(),
		}},
	}); err != nil {
		log.Error(err, "unable to update owned resources")
		return r.ManageError(ctx, instance, err)
	}

	sentinelPool, err := redis.NewSentinelPool(ctx, r.GetClient(),
		types.NamespacedName{Name: gen.GetComponent(), Namespace: gen.GetNamespace()}, int(*instance.Spec.Replicas))
	if err != nil {
		return r.ManageError(ctx, instance, err)
	}

	allMonitored, err := sentinelPool.IsMonitoringShards(ctx,
		func() []string {
			keys := make([]string, 0, len(instance.Spec.Config.MonitoredShards))
			for k := range instance.Spec.Config.MonitoredShards {
				keys = append(keys, k)
			}
			return keys
		}(),
	)
	if err != nil {
		return r.ManageError(ctx, instance, err)
	}

	if !allMonitored {
		shardedCluster, err := redis.NewShardedCluster(ctx, instance.Spec.Config.MonitoredShards, log)
		if err != nil {
			return r.ManageError(ctx, instance, err)
		}

		if err := sentinelPool.Monitor(ctx, shardedCluster); err != nil {
			return r.ManageError(ctx, instance, err)
		}
	}

	r.ReconcileExporters(ctx, gen, log.WithName("exporter"))

	return r.ManageSuccess(ctx, instance)
}

// ReconcileExporters ensures that all Pods within the statefulset have a running exporter. It also stops
// exporters for no longer running replicas (in the case the statefulset number of replicas is reduced)
func (r *SentinelReconciler) ReconcileExporters(ctx context.Context, gen sentinel.Generator, log logr.Logger) {

	if r.metrics.exporters == nil {
		r.metrics.exporters = map[string]*redismetrics.SentinelMetricsGatherer{}
	}

	// Gather metrics for each sentinel replica
	shouldRun := map[string]int{}
	for i, fn := range gen.PodServices() {
		svc := fn()
		key := fmt.Sprintf("%s/%s/%d", gen.Namespace, gen.InstanceName, i)
		sentinelURL := fmt.Sprintf("redis://%s.%s.svc.cluster.local:%d", svc.GetName(), svc.GetNamespace(), saasv1alpha1.SentinelPort)
		shouldRun[key] = 1
		r.metrics.RunExporter(ctx, key, sentinelURL, *gen.Spec.Config.MetricsRefreshInterval, log)
	}

	// Stop gathering metrics for any sentinel replica that does not exist anymore
	for key := range r.metrics.exporters {

		if strings.Contains(key, fmt.Sprintf("%s/%s/", gen.Namespace, gen.InstanceName)) {
			if _, ok := shouldRun[key]; !ok {
				r.metrics.StopExporter(key)
			}
		}
	}
}

// CleanupExporters stops all the exporters for this instance of the Sentinel custom resource
// This is used as a cleanup function in the finalize phase of the controller loop.
func (r *SentinelReconciler) CleanupExporters(instance types.NamespacedName) func() {
	return func() {
		prefix := fmt.Sprintf("%s/%s/", instance.Namespace, instance.Name)
		for key := range r.metrics.exporters {
			if strings.Contains(key, prefix) {
				r.metrics.StopExporter(key)
			}
		}
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *SentinelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&saasv1alpha1.Sentinel{}).
		Watches(&source.Channel{Source: r.GetStatusChangeChannel()}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
