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

	"github.com/3scale-ops/basereconciler/reconciler"
	"github.com/3scale-ops/basereconciler/util"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/generators/echoapi"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EchoAPIReconciler reconciles a EchoAPI object
type EchoAPIReconciler struct {
	*reconciler.Reconciler
}

// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=echoapis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=echoapis/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=echoapis/finalizers,verbs=update
// +kubebuilder:rbac:groups="core",namespace=placeholder,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="apps",namespace=placeholder,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="monitoring.coreos.com",namespace=placeholder,resources=podmonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="autoscaling",namespace=placeholder,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="policy",namespace=placeholder,resources=poddisruptionbudgets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="marin3r.3scale.net",namespace=placeholder,resources=envoyconfigs,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *EchoAPIReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	ctx, _ = r.Logger(ctx, "name", req.Name, "namespace", req.Namespace)
	instance := &saasv1alpha1.EchoAPI{}
	result := r.ManageResourceLifecycle(ctx, req, instance,
		reconciler.WithInMemoryInitializationFunc(util.ResourceDefaulter(instance)),
		reconciler.WithInitializationFunc(EchoapiResourceUpgrader))
	if result.ShouldReturn() {
		return result.Values()
	}

	gen := echoapi.NewGenerator(instance.GetName(), instance.GetNamespace(), instance.Spec)

	resources, err := gen.Resources()
	if err != nil {
		return ctrl.Result{}, err
	}

	result = r.ReconcileOwnedResources(ctx, instance, resources)
	if result.ShouldReturn() {
		return result.Values()
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EchoAPIReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return reconciler.SetupWithDynamicTypeWatches(r,
		ctrl.NewControllerManagedBy(mgr).
			For(&saasv1alpha1.EchoAPI{}),
	)
}

func EchoapiResourceUpgrader(ctx context.Context, cl client.Client, o client.Object) error {
	instance := o.(*saasv1alpha1.EchoAPI)

	if instance.Spec.PublishingStrategies == nil {
		pss, err := saasv1alpha1.UpgradeCR2PublishingStrategies(ctx, cl,
			saasv1alpha1.WorkloadPublishingStrategyUpgrader{
				EndpointName: "HTTP",
				ServiceName:  "echo-api-nlb",
				Namespace:    instance.GetNamespace(),
				ServiceType:  saasv1alpha1.ServiceTypeNLB,
				Endpoint:     instance.Spec.Endpoint,
				Marin3r:      instance.Spec.Marin3r,
				NLBSpec:      instance.Spec.LoadBalancer,
				ServicePortOverrides: []corev1.ServicePort{
					{
						Name:       "http",
						Protocol:   corev1.ProtocolTCP,
						Port:       80,
						TargetPort: intstr.FromString("echo-api-http"),
					},
					{
						Name:       "https",
						Protocol:   corev1.ProtocolTCP,
						Port:       443,
						TargetPort: intstr.FromString("echo-api-https"),
					},
				},
			},
		)

		if err != nil {
			return err
		}

		if len(pss.Endpoints) > 0 {
			instance.Spec.PublishingStrategies = pss
			instance.Spec.Endpoint = nil
			instance.Spec.LoadBalancer = nil
			instance.Spec.Marin3r = nil
		}
	}

	return nil
}
