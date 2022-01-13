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
	"time"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators/redisshard"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	"github.com/3scale/saas-operator/pkg/redis"
	"github.com/3scale/saas-operator/pkg/redis/crud/client"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// RedisShardReconciler reconciles a RedisShard object
type RedisShardReconciler struct {
	basereconciler.Reconciler
	Log logr.Logger
}

// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=redisshards,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=redisshards/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=redisshards/finalizers,verbs=update
// +kubebuilder:rbac:groups="core",namespace=placeholder,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="core",namespace=placeholder,resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups="apps",namespace=placeholder,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *RedisShardReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("name", req.Name, "namespace", req.Namespace)

	instance := &saasv1alpha1.RedisShard{}
	key := types.NamespacedName{Name: req.Name, Namespace: req.Namespace}
	result, err := r.GetInstance(ctx, key, instance, saasv1alpha1.Finalizer, nil, log)
	if result != nil || err != nil {
		return *result, err
	}

	// Apply defaults for reconcile but do not store them in the API
	instance.Default()

	gen := redisshard.NewGenerator(
		instance.GetName(),
		instance.GetNamespace(),
		instance.Spec,
	)

	if err := r.ReconcileOwnedResources(ctx, instance, r.GetScheme(), gen.Resources()); err != nil {
		log.Error(err, "unable to update owned resources")
		return r.ManageError(ctx, instance, err)
	}

	shard, result, err := r.setRedisRoles(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace},
		*instance.Spec.MasterIndex, *instance.Spec.SlaveCount+1, gen.ServiceName(), log)
	if result != nil || err != nil {
		return *result, err
	}

	if err = r.updateStatus(ctx, shard, instance, log); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RedisShardReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&saasv1alpha1.RedisShard{}).
		Watches(&source.Channel{Source: r.GetStatusChangeChannel()}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}

func (r *RedisShardReconciler) setRedisRoles(ctx context.Context, key types.NamespacedName, masterIndex, replicas int32, serviceName string, log logr.Logger) (*redis.Shard, *ctrl.Result, error) {

	redisURLs := make([]string, replicas)
	for i := 0; i < int(replicas); i++ {
		pod := &corev1.Pod{}
		key := types.NamespacedName{Name: fmt.Sprintf("%s-%d", serviceName, i), Namespace: key.Namespace}
		err := r.GetClient().Get(ctx, key, pod)
		if err != nil {
			return nil, &ctrl.Result{}, err
		}

		redisURLs[i] = fmt.Sprintf("redis://%s:%d", pod.Status.PodIP, 6379)
	}

	shard, err := redis.NewShard(key.Name, redisURLs)
	if err != nil {
		return nil, &ctrl.Result{}, err
	}

	_, err = shard.Init(ctx, masterIndex, log)
	if err != nil {
		log.Info("waiting for redis shard init")
		return nil, &ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
	}

	return shard, nil, nil
}

func (r *RedisShardReconciler) updateStatus(ctx context.Context, shard *redis.Shard, instance *saasv1alpha1.RedisShard, log logr.Logger) error {

	status := saasv1alpha1.RedisShardStatus{
		ShardNodes: &saasv1alpha1.RedisShardNodes{Master: nil, Slaves: []string{}},
	}

	for _, server := range shard.Servers {
		if server.Role == client.Master {
			status.ShardNodes.Master = pointer.StringPtr(server.Name)
		} else if server.Role == client.Slave {
			status.ShardNodes.Slaves = append(status.ShardNodes.Slaves, server.Name)
		}
	}
	if !equality.Semantic.DeepEqual(status, instance.Status) {
		instance.Status = status
		if err := r.GetClient().Status().Update(ctx, instance); err != nil {
			return err
		}
	}

	return nil
}
