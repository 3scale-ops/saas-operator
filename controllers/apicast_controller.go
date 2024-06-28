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
	"github.com/3scale-ops/saas-operator/pkg/generators/apicast"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ApicastReconciler reconciles a Apicast object
type ApicastReconciler struct {
	*reconciler.Reconciler
}

// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=apicasts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=apicasts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=apicasts/finalizers,verbs=update
// +kubebuilder:rbac:groups="core",namespace=placeholder,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="apps",namespace=placeholder,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="monitoring.coreos.com",namespace=placeholder,resources=podmonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="autoscaling",namespace=placeholder,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="policy",namespace=placeholder,resources=poddisruptionbudgets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="grafana.integreatly.org",namespace=placeholder,resources=grafanadashboards,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="marin3r.3scale.net",namespace=placeholder,resources=envoyconfigs,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ApicastReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	ctx, _ = r.Logger(ctx, "name", req.Name, "namespace", req.Namespace)
	instance := &saasv1alpha1.Apicast{}
	result := r.ManageResourceLifecycle(ctx, req, instance,
		reconciler.WithInMemoryInitializationFunc(util.ResourceDefaulter(instance)),
		reconciler.WithInitializationFunc(ApicastResourceUpgrader),
	)
	if result.ShouldReturn() {
		return result.Values()
	}

	gen, err := apicast.NewGenerator(instance.GetName(), instance.GetNamespace(), instance.Spec)
	if err != nil {
		return ctrl.Result{}, err
	}
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
func (r *ApicastReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return reconciler.SetupWithDynamicTypeWatches(r,
		ctrl.NewControllerManagedBy(mgr).
			For(&saasv1alpha1.Apicast{}),
	)
}

func ApicastResourceUpgrader(ctx context.Context, cl client.Client, o client.Object) error {
	instance := o.(*saasv1alpha1.Apicast)

	if instance.Spec.Production.PublishingStrategies == nil {
		pss, err := saasv1alpha1.UpgradeCR2PublishingStrategies(ctx, cl,
			saasv1alpha1.WorkloadPublishingStrategyUpgrader{
				EndpointName: "Gateway",
				ServiceName:  "apicast-production",
				Namespace:    instance.GetNamespace(),
				ServiceType:  saasv1alpha1.ServiceTypeELB,
				Endpoint:     instance.Spec.Production.Endpoint,
				Marin3r:      instance.Spec.Production.Marin3r,
				ELBSpec:      instance.Spec.Production.LoadBalancer,
				NLBSpec:      nil,
				ServicePortOverrides: []corev1.ServicePort{
					{
						Name:       "gateway-http",
						Protocol:   corev1.ProtocolTCP,
						Port:       80,
						TargetPort: intstr.FromString("gateway-http"),
					},
					{
						Name:       "gateway-https",
						Protocol:   corev1.ProtocolTCP,
						Port:       443,
						TargetPort: intstr.FromString("gateway-https"),
					},
				},
			},
			saasv1alpha1.WorkloadPublishingStrategyUpgrader{
				EndpointName: "Management",
				ServiceName:  "apicast-production-management",
				Namespace:    instance.GetNamespace(),
				ServiceType:  saasv1alpha1.ServiceTypeClusterIP,
			},
		)

		if err != nil {
			return err
		}

		if len(pss.Endpoints) > 0 {
			instance.Spec.Production.PublishingStrategies = pss
			instance.Spec.Production.Endpoint = nil
			instance.Spec.Production.Marin3r = nil
			instance.Spec.Production.LoadBalancer = nil
		}
	}

	if instance.Spec.Staging.PublishingStrategies == nil {
		pss, err := saasv1alpha1.UpgradeCR2PublishingStrategies(ctx, cl,
			saasv1alpha1.WorkloadPublishingStrategyUpgrader{
				EndpointName: "Gateway",
				ServiceName:  "apicast-staging",
				Namespace:    instance.GetNamespace(),
				ServiceType:  saasv1alpha1.ServiceTypeELB,
				Endpoint:     instance.Spec.Staging.Endpoint,
				Marin3r:      instance.Spec.Staging.Marin3r,
				ELBSpec:      instance.Spec.Staging.LoadBalancer,
				NLBSpec:      nil,
				ServicePortOverrides: []corev1.ServicePort{
					{
						Name:       "gateway-http",
						Protocol:   corev1.ProtocolTCP,
						Port:       80,
						TargetPort: intstr.FromString("gateway-http"),
					},
					{
						Name:       "gateway-https",
						Protocol:   corev1.ProtocolTCP,
						Port:       443,
						TargetPort: intstr.FromString("gateway-https"),
					},
				},
			},
			saasv1alpha1.WorkloadPublishingStrategyUpgrader{
				EndpointName: "Management",
				ServiceName:  "apicast-staging-management",
				Namespace:    instance.GetNamespace(),
				ServiceType:  saasv1alpha1.ServiceTypeClusterIP,
			},
		)

		if err != nil {
			return err
		}

		if len(pss.Endpoints) > 0 {
			instance.Spec.Staging.PublishingStrategies = pss
			instance.Spec.Staging.Endpoint = nil
			instance.Spec.Staging.Marin3r = nil
			instance.Spec.Staging.LoadBalancer = nil
		}
	}

	return nil
}
