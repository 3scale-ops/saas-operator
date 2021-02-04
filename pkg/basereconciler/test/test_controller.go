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

package test

import (
	"context"

	secretsmanagerv1alpha1 "github.com/3scale/saas-operator/pkg/apis/secrets-manager/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	"github.com/3scale/saas-operator/pkg/basereconciler/test/api/v1alpha1"
	"github.com/go-logr/logr"
	"github.com/redhat-cop/operator-utils/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Reconciler reconciles a Test object
// +kubebuilder:object:generate=false
type Reconciler struct {
	basereconciler.Reconciler
	Log logr.Logger
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("name", req.Name, "namespace", req.Namespace)

	instance := &v1alpha1.Test{}
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
		if !util.HasFinalizer(instance, "finalizer.example.com") {
			return ctrl.Result{}, nil
		}
		err := r.ManageCleanUpLogic(instance, log)
		if err != nil {
			log.Error(err, "unable to delete instance")
			return r.ManageError(ctx, instance, err)
		}
		util.RemoveFinalizer(instance, "finalizer.example.com")
		err = r.GetClient().Update(ctx, instance)
		if err != nil {
			log.Error(err, "unable to update instance")
			return r.ManageError(ctx, instance, err)
		}
		return ctrl.Result{}, nil
	}

	if ok := r.IsInitialized(instance, "finalizer.example.com"); !ok {
		err := r.GetClient().Update(ctx, instance)
		if err != nil {
			log.Error(err, "unable to initialize instance")
			return r.ManageError(ctx, instance, err)
		}
		return ctrl.Result{}, nil
	}

	// triggerName := gen.SecretDefinition()().GetName()
	// secret, err := r.SecretFromSecretDef(ctx, gen.SecretDefinition())
	// if err != nil {
	// 	return r.ManageError(ctx, instance, err)
	// }
	// var trigger basereconciler.RolloutTrigger
	// if secret != nil {
	// 	trigger = basereconciler.NewRolloutTrigger(triggerName, secret)
	// } else {
	// 	trigger = basereconciler.NewRolloutTrigger(triggerName, &corev1.Secret{})
	// }

	err = r.ReconcileOwnedResources(ctx, instance, basereconciler.ControlledResources{
		Deployments: []basereconciler.Deployment{{
			Template:        deployment(req.Namespace),
			RolloutTriggers: []basereconciler.RolloutTrigger{},
			HasHPA:          false,
		}},
		SecretDefinitions: []basereconciler.SecretDefinition{{
			Template: secretDefinition(req.Namespace),
			Enabled:  true,
		}},
		Services: []basereconciler.Service{{
			Template: service(req.Namespace),
			Enabled:  true,
		}},
		PodDisruptionBudgets: []basereconciler.PodDisruptionBudget{{
			Template: nil,
			Enabled:  false,
		}},
		HorizontalPodAutoscalers: []basereconciler.HorizontalPodAutoscaler{{
			Template: nil,
			Enabled:  false,
		}},
		PodMonitors: []basereconciler.PodMonitor{{
			Template: nil,
			Enabled:  false}},
		GrafanaDashboards: []basereconciler.GrafanaDashboard{{
			Template: nil,
			Enabled:  false,
		}},
	})
	if err != nil {
		log.Error(err, "unable to reconcile owned resources")
		return r.ManageError(ctx, instance, err)
	}

	return r.ManageSuccess(ctx, instance)
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Test{}).
		Watches(&source.Channel{Source: r.GetStatusChangeChannel()}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &corev1.Secret{TypeMeta: metav1.TypeMeta{Kind: "Secret"}}},
			r.SecretEventHandler(&v1alpha1.TestList{}, r.Log)).
		Complete(r)
}

func deployment(namespace string) basereconciler.GeneratorFunction {
	return func() client.Object {
		return &appsv1.Deployment{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Deployment",
				APIVersion: appsv1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "deployment",
				Namespace: namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: pointer.Int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"selector": "deployment"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"selector": "deployment"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:      "container",
								Image:     "example.com:latest",
								Resources: corev1.ResourceRequirements{},
							},
						},
					},
				},
			},
		}
	}
}

func service(namespace string) basereconciler.GeneratorFunction {
	return func() client.Object {
		return &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: corev1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service",
				Namespace: namespace,
			},
			Spec: corev1.ServiceSpec{
				Type:                  corev1.ServiceTypeLoadBalancer,
				ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeCluster,
				SessionAffinity:       corev1.ServiceAffinityNone,
				Ports:                 []corev1.ServicePort{{Name: "port", Port: 80, TargetPort: intstr.FromInt(80)}},
				Selector:              map[string]string{"selector": "deployment"},
			},
		}
	}
}

func secretDefinition(namespace string) basereconciler.GeneratorFunction {

	return func() client.Object {
		return &secretsmanagerv1alpha1.SecretDefinition{
			TypeMeta: metav1.TypeMeta{
				Kind:       "SecretDefinition",
				APIVersion: secretsmanagerv1alpha1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret-definition",
				Namespace: namespace,
			},
			Spec: secretsmanagerv1alpha1.SecretDefinitionSpec{
				Name: "secret-definition",
				Type: "opaque",
				KeysMap: map[string]secretsmanagerv1alpha1.DataSource{
					"KEY": {Key: "vault-key", Path: "vault-path"},
				},
			},
		}
	}
}
