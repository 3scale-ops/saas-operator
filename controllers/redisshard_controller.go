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
	"time"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	"github.com/3scale/saas-operator/pkg/generators/redisshard"
	"github.com/go-logr/logr"
	"github.com/go-redis/redis/v8"
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
	json, _ := json.Marshal(instance.Spec)
	log.V(1).Info("Apply defaults before resolving templates", "JSON", string(json))

	gen := redisshard.NewGenerator(
		instance.GetName(),
		instance.GetNamespace(),
		instance.Spec,
	)

	err = r.ReconcileOwnedResources(ctx, instance, basereconciler.ControlledResources{
		StatefulSets: []basereconciler.StatefulSet{{
			Template:        gen.StatefulSet(),
			RolloutTriggers: nil,
			Enabled:         true,
		}},
		Services: []basereconciler.Service{{
			Template: gen.Service(),
			Enabled:  true,
		}},
		ConfigMaps: []basereconciler.ConfigMaps{
			{
				Template: gen.RedisConfigConfigMap(),
				Enabled:  true,
			},
			{
				Template: gen.RedisReadinessScriptConfigMap(),
				Enabled:  true,
			},
		},
	})

	if err != nil {
		log.Error(err, "unable to update owned resources")
		return r.ManageError(ctx, instance, err)
	}

	shard := make(RedisShardStatus, saasv1alpha1.RedisShardDefaultReplicas)
	result, err = r.setRedisRoles(ctx, shard, *instance.Spec.MasterIndex, gen.Service()().GetName(), req.Namespace, log)
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

func (r *RedisShardReconciler) setRedisRoles(ctx context.Context, shard RedisShardStatus, masterIndex int32, serviceName, namespace string, log logr.Logger) (*ctrl.Result, error) {
	for i := 0; i < int(saasv1alpha1.RedisShardDefaultReplicas); i++ {
		pod := &corev1.Pod{}
		key := types.NamespacedName{Name: fmt.Sprintf("%s-%d", serviceName, i), Namespace: namespace}
		err := r.GetClient().Get(ctx, key, pod)
		if err != nil {
			return &ctrl.Result{}, err
		}
		shard[i].Hostname = pod.GetName()
		shard[i].IP = pod.Status.PodIP
		shard[i].Port = 6379
		shard[i].Role = Unkown

	}

	err := shard.SetRoles(ctx, masterIndex, log)
	if err != nil {
		return &ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, err
	}

	return nil, nil
}

func (r *RedisShardReconciler) updateStatus(ctx context.Context, shard RedisShardStatus, instance *saasv1alpha1.RedisShard, log logr.Logger) error {

	status := saasv1alpha1.RedisShardStatus{
		ShardNodes: &saasv1alpha1.RedisShardNodes{Master: nil, Slaves: []string{}},
	}

	for _, server := range shard {
		if server.Role == Master {
			status.ShardNodes.Master = pointer.StringPtr(fmt.Sprintf("redis://%s:%d", server.IP, server.Port))
		} else if server.Role == Slave {
			status.ShardNodes.Slaves = append(status.ShardNodes.Slaves, fmt.Sprintf("redis://%s:%d", server.IP, server.Port))
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

type RedisRole string

const (
	Master RedisRole = "master"
	Slave  RedisRole = "slave"
	Unkown RedisRole = "unknown"
)

type RedisShardStatus []RedisServer

type RedisServer struct {
	Hostname string
	IP       string
	Port     uint32
	Role     RedisRole
}

func (rss *RedisShardStatus) SetRoles(ctx context.Context, masterIndex int32, log logr.Logger) error {

	for idx, server := range *rss {
		rdb := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", server.IP, server.Port),
			Password: "",
			DB:       0,
		})

		val, err := rdb.Do(ctx, "role").Result()
		if err != nil {
			return err
		}

		role := val.([]interface{})[0].(string)
		if role == string(Slave) {

			slaveof := val.([]interface{})[1].(string)
			if slaveof == "127.0.0.1" {

				if idx == int(masterIndex) {
					_, err := rdb.SlaveOf(ctx, "NO", "ONE").Result()
					if err != nil {
						return err
					}
					log.Info(fmt.Sprintf("[@redis-setup] Configured %s as master", server.Hostname))
				} else {
					_, err := rdb.SlaveOf(ctx, []RedisServer(*rss)[masterIndex].IP, fmt.Sprintf("%d", []RedisServer(*rss)[masterIndex].Port)).Result()
					if err != nil {
						return err
					}
					log.Info(fmt.Sprintf("[@redis-setup] Configured %s as slave", server.Hostname))
				}

			} else {
				[]RedisServer(*rss)[idx].Role = Slave
			}

		} else if role == string(Master) {
			[]RedisServer(*rss)[idx].Role = Master
		} else {
			return fmt.Errorf("[@redis-setup] unable to get role for server %s:%d", server.IP, server.Port)
		}
	}

	return nil
}
