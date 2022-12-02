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
	"reflect"
	"sync"
	"time"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	grafanav1alpha1 "github.com/3scale/saas-operator/pkg/apis/grafana/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators/twemproxyconfig"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	"github.com/3scale/saas-operator/pkg/reconcilers/threads"
	"github.com/3scale/saas-operator/pkg/redis/events"
	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
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
// +kubebuilder:rbac:groups="core",namespace=placeholder,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="core",namespace=placeholder,resources=pods,verbs=list;patch
// +kubebuilder:rbac:groups="integreatly.org",namespace=placeholder,resources=grafanadashboards,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *TwemproxyConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("name", req.Name, "namespace", req.Namespace)
	ctx = log.IntoContext(ctx, logger)

	instance := &saasv1alpha1.TwemproxyConfig{}
	key := types.NamespacedName{Name: req.Name, Namespace: req.Namespace}
	result, err := r.GetInstance(ctx, key, instance,
		pointer.String(saasv1alpha1.Finalizer),
		[]func(){r.SentinelEvents.CleanupThreads(instance)})
	if result != nil || err != nil {
		return *result, err
	}

	// Apply defaults for reconcile but do not store them in the API
	instance.Default()

	// Generate the ConfigMap
	gen, err := twemproxyconfig.NewGenerator(
		ctx, instance, r.Client, logger.WithName("generator"),
	)

	if err != nil {
		return ctrl.Result{}, err
	}
	cm, _, err := gen.ConfigMap().Build(ctx, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Reconcile the ConfigMap
	hash, err := r.reconcileConfigMap(ctx, instance, cm.(*corev1.ConfigMap), *instance.Spec.ReconcileServerPools, logger)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Reconcile sync annotations in pods. This is done to force a change in the target
	// Pods annotations so the ConfigMap is re-synced inside the container. Otherwide kubelet
	// would re-sync the file asynchronously depending on its configured refresh time, which might
	// take several seconds.
	if err := r.reconcileSyncAnnotations(ctx, instance, hash, logger); err != nil {
		return ctrl.Result{}, err
	}

	// Reconcile sentinel event watchers
	eventWatchers := make([]threads.RunnableThread, 0, len(gen.Spec.SentinelURIs))
	for _, uri := range gen.Spec.SentinelURIs {
		eventWatchers = append(eventWatchers, &events.SentinelEventWatcher{
			Instance:      instance,
			SentinelURI:   uri,
			ExportMetrics: false,
		})
	}
	r.SentinelEvents.ReconcileThreads(ctx, instance, eventWatchers, logger.WithName("event-watcher"))

	if err := r.ReconcileOwnedResources(
		ctx, instance, []basereconciler.Resource{gen.GrafanaDashboard()},
	); err != nil {
		return ctrl.Result{}, err
	}

	// Reconcile periodically in case some event is lost ...
	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

func (r *TwemproxyConfigReconciler) reconcileConfigMap(ctx context.Context, owner client.Object,
	desired *corev1.ConfigMap, reconcileData bool, log logr.Logger) (string, error) {

	current := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, util.ObjectKey(desired), current)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create
			if err := controllerutil.SetControllerReference(owner, desired, r.Scheme); err != nil {
				return "", err
			}
			if err := r.Client.Create(ctx, desired); err != nil {
				return "", err
			}
			log.Info("created ConfigMap")
			return util.Hash(desired.Data), nil
		}
		return "", err
	}

	var hash string
	if reconcileData {
		// Compare .data field of both ConfigMaps and patch if required.
		// We use patch to avoid failures due to having an older version
		// of the configmap so the config changes are propagated faster.
		if !reflect.DeepEqual(desired.Data, current.Data) {
			if err := r.Client.Patch(ctx, desired, client.MergeFrom(current)); err != nil {
				log.Error(err, "unable to patch ConfigMap")
				return "", err
			}
			log.Info("updated ConfigMap")
		}
		hash = util.Hash(desired.Data)
	} else {
		hash = util.Hash(current.Data)
	}

	return hash, nil
}

func (r *TwemproxyConfigReconciler) reconcileSyncAnnotations(ctx context.Context,
	instance *saasv1alpha1.TwemproxyConfig, hash string, log logr.Logger) error {

	podList := &corev1.PodList{}
	if err := r.Client.List(ctx, podList, instance.PodSyncSelector(),
		client.InNamespace(instance.GetNamespace())); err != nil {
		return err
	}

	failures := util.MultiError{}
	errCh := make(chan error)
	innerCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	var wg sync.WaitGroup

	// listen the error channel for errors
	go func() {
		for {
			select {
			case err := <-errCh:
				failures = append(failures, err)
			case <-innerCtx.Done():
				return
			}
		}
	}()

	// Patch the Pods concurrently
	for _, pod := range podList.Items {
		wg.Add(1)
		go func(pod corev1.Pod) {
			r.syncPod(innerCtx, pod, hash, errCh, log)
			wg.Done()
		}(pod)
	}

	wg.Wait()
	cancel()
	if len(failures) > 0 {
		return failures
	}

	return nil
}

func (r *TwemproxyConfigReconciler) syncPod(ctx context.Context, pod corev1.Pod, hash string, errCh chan<- error, log logr.Logger) {
	annotatedHash, ok := pod.GetAnnotations()[saasv1alpha1.TwemproxySyncAnnotationKey]
	if !ok || annotatedHash != hash {
		patch := client.MergeFrom(pod.DeepCopy())
		if pod.GetAnnotations() != nil {
			pod.ObjectMeta.Annotations[saasv1alpha1.TwemproxySyncAnnotationKey] = hash
		} else {
			pod.ObjectMeta.Annotations = map[string]string{
				saasv1alpha1.TwemproxySyncAnnotationKey: hash,
			}
		}

		if err := r.Client.Patch(ctx, &pod, patch); err != nil {
			errCh <- err
		}
		log.V(1).Info(fmt.Sprintf("configmap re-sync forced in target pod %s", util.ObjectKey(&pod)))
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *TwemproxyConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&saasv1alpha1.TwemproxyConfig{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&grafanav1alpha1.GrafanaDashboard{}).
		Watches(&source.Channel{Source: r.SentinelEvents.GetChannel()}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
