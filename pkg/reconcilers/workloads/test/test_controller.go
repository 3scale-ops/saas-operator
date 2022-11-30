package test

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

import (
	"context"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	externalsecretsv1beta1 "github.com/3scale/saas-operator/pkg/apis/externalsecrets/v1beta1"
	"github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2/resources"
	"github.com/3scale/saas-operator/pkg/reconcilers/workloads"
	"github.com/3scale/saas-operator/pkg/reconcilers/workloads/test/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-logr/logr"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// WorkloadReconciler reconciles a Test object
// +kubebuilder:object:generate=false
type Reconciler struct {
	workloads.WorkloadReconciler
	Log logr.Logger
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("name", req.Name, "namespace", req.Namespace)
	ctx = log.IntoContext(ctx, logger)

	instance := &v1alpha1.Test{}
	key := types.NamespacedName{Name: req.Name, Namespace: req.Namespace}
	result, err := r.GetInstance(ctx, key, instance, nil, nil)
	if result != nil || err != nil {
		return *result, err
	}

	main := &TestWorkloadGenerator{
		TName:            instance.Spec.Main.Name,
		TNamespace:       req.Namespace,
		TTraffic:         instance.Spec.Main.Traffic,
		TLabels:          instance.Spec.Main.Labels,
		TSelector:        instance.Spec.Main.Selector,
		TTrafficSelector: instance.Spec.TrafficSelector,
	}

	canary := &TestWorkloadGenerator{
		TName:            instance.Spec.Canary.Name,
		TNamespace:       req.Namespace,
		TTraffic:         instance.Spec.Canary.Traffic,
		TLabels:          instance.Spec.Canary.Labels,
		TSelector:        instance.Spec.Canary.Selector,
		TTrafficSelector: instance.Spec.TrafficSelector,
	}

	deployments, err := r.NewDeploymentWorkload(main, canary)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Reconcile all resources
	err = r.ReconcileOwnedResources(ctx, instance, deployments)
	if err != nil {
		logger.Error(err, "unable to reconcile owned resources")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Test{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&policyv1.PodDisruptionBudget{}).
		Owns(&autoscalingv2beta2.HorizontalPodAutoscaler{}).
		Owns(&externalsecretsv1beta1.ExternalSecret{}).
		Watches(&source.Kind{Type: &corev1.Secret{TypeMeta: metav1.TypeMeta{Kind: "Secret"}}},
			r.SecretEventHandler(&v1alpha1.TestList{}, r.Log)).
		Complete(r)
}

var _ workloads.WithTraffic = &TestWorkloadGenerator{}
var _ workloads.DeploymentWorkload = &TestWorkloadGenerator{}

type TestWorkloadGenerator struct {
	TName            string
	TNamespace       string
	TTraffic         bool
	TLabels          map[string]string
	TSelector        map[string]string
	TTrafficSelector map[string]string
}

func (gen *TestWorkloadGenerator) Services() []resources.ServiceTemplate {
	return []resources.ServiceTemplate{{
		Template: func() *corev1.Service {
			return &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "service",
					Namespace: gen.TNamespace,
				},
				Spec: corev1.ServiceSpec{
					Type:                  corev1.ServiceTypeLoadBalancer,
					ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeCluster,
					SessionAffinity:       corev1.ServiceAffinityNone,
					Ports: []corev1.ServicePort{{
						Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP}},
				},
			}
		},
		IsEnabled: true,
	},
	}
}

func (gen *TestWorkloadGenerator) TrafficSelector() map[string]string {
	return gen.TTrafficSelector
}

func (gen *TestWorkloadGenerator) Deployment() resources.DeploymentTemplate {
	return resources.DeploymentTemplate{
		Template: func() *appsv1.Deployment {
			return &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Replicas: pointer.Int32Ptr(1),
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"orig-key": "orig-value"},
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
		},
		RolloutTriggers: []resources.RolloutTrigger{{
			Name:       "secret",
			SecretName: pointer.String("secret"),
		}},
		IsEnabled:       true,
		EnforceReplicas: true,
	}
}

func (gen *TestWorkloadGenerator) MonitoredEndpoints() []monitoringv1.PodMetricsEndpoint { return nil }
func (gen *TestWorkloadGenerator) GetKey() types.NamespacedName {
	return types.NamespacedName{Name: gen.TName, Namespace: gen.TNamespace}
}
func (gen *TestWorkloadGenerator) GetLabels() map[string]string { return gen.TLabels }
func (gen *TestWorkloadGenerator) GetSelector() map[string]string {
	return gen.TSelector
}
func (gen *TestWorkloadGenerator) HPASpec() *saasv1alpha1.HorizontalPodAutoscalerSpec {
	return &saasv1alpha1.HorizontalPodAutoscalerSpec{
		MinReplicas:         pointer.Int32Ptr(1),
		MaxReplicas:         pointer.Int32Ptr(2),
		ResourceUtilization: pointer.Int32Ptr(90),
		ResourceName:        pointer.StringPtr("cpu"),
	}
}
func (gen *TestWorkloadGenerator) PDBSpec() *saasv1alpha1.PodDisruptionBudgetSpec {
	return &saasv1alpha1.PodDisruptionBudgetSpec{
		MaxUnavailable: util.IntStrPtr(intstr.FromInt(1)),
	}
}
func (gen *TestWorkloadGenerator) SendTraffic() bool { return gen.TTraffic }
