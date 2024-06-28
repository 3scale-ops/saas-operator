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
	"github.com/3scale-ops/saas-operator/pkg/generators/backend"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BackendReconciler reconciles a Backend object
type BackendReconciler struct {
	*reconciler.Reconciler
}

// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=backends,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=backends/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=backends/finalizers,verbs=update
// +kubebuilder:rbac:groups="core",namespace=placeholder,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="core",namespace=placeholder,resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="apps",namespace=placeholder,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="monitoring.coreos.com",namespace=placeholder,resources=podmonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="autoscaling",namespace=placeholder,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="policy",namespace=placeholder,resources=poddisruptionbudgets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="grafana.integreatly.org",namespace=placeholder,resources=grafanadashboards,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="external-secrets.io",namespace=placeholder,resources=externalsecrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="marin3r.3scale.net",namespace=placeholder,resources=envoyconfigs,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *BackendReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	ctx, _ = r.Logger(ctx, "name", req.Name, "namespace", req.Namespace)
	instance := &saasv1alpha1.Backend{}
	result := r.ManageResourceLifecycle(ctx, req, instance,
		reconciler.WithInMemoryInitializationFunc(util.ResourceDefaulter(instance)),
		reconciler.WithInitializationFunc(BackendResourceUpgrader))
	if result.ShouldReturn() {
		return result.Values()
	}

	gen, err := backend.NewGenerator(instance.GetName(), instance.GetNamespace(), instance.Spec)
	if err != nil {
		return ctrl.Result{}, err
	}

	resources, err := gen.Resources()
	if err != nil {
		return ctrl.Result{}, err
	}

	// Reconcile all resources
	result = r.ReconcileOwnedResources(ctx, instance, resources)
	if result.ShouldReturn() {
		return result.Values()
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BackendReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return reconciler.SetupWithDynamicTypeWatches(r,
		ctrl.NewControllerManagedBy(mgr).
			For(&saasv1alpha1.Backend{}).
			Watches(&corev1.Secret{}, r.FilteredEventHandler(&saasv1alpha1.BackendList{}, nil, r.Log)),
	)
}

func BackendResourceUpgrader(ctx context.Context, cl client.Client, o client.Object) error {
	instance := o.(*saasv1alpha1.Backend)

	if instance.Spec.Listener.PublishingStrategies == nil {
		pss, err := saasv1alpha1.UpgradeCR2PublishingStrategies(ctx, cl,
			saasv1alpha1.WorkloadPublishingStrategyUpgrader{
				EndpointName: "HTTP",
				ServiceName:  "backend-listener-nlb",
				Namespace:    instance.GetNamespace(),
				ServiceType:  saasv1alpha1.ServiceTypeNLB,
				Endpoint:     instance.Spec.Listener.Endpoint,
				Marin3r:      instance.Spec.Listener.Marin3r,
				NLBSpec:      instance.Spec.Listener.LoadBalancer,
				ServicePortOverrides: []corev1.ServicePort{
					{
						Name:       "http",
						Protocol:   corev1.ProtocolTCP,
						Port:       80,
						TargetPort: intstr.FromString("backend-http"),
					},
					{
						Name:       "https",
						Protocol:   corev1.ProtocolTCP,
						Port:       443,
						TargetPort: intstr.FromString("backend-https"),
					},
				},
			},
			saasv1alpha1.WorkloadPublishingStrategyUpgrader{
				EndpointName: "Internal",
				ServiceName:  "backend-listener-internal",
				Namespace:    instance.GetNamespace(),
				ServiceType:  saasv1alpha1.ServiceTypeClusterIP,
				ServicePortOverrides: []corev1.ServicePort{{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.FromString("http-internal"),
				}},
				// Use 'Create' mode to add the internal endpoint
				Create: true,
			},
		)
		if err != nil {
			return err
		}

		if len(pss.Endpoints) > 0 {
			instance.Spec.Listener.PublishingStrategies = pss
			instance.Spec.Listener.Marin3r = nil
			instance.Spec.Listener.Endpoint = nil
			instance.Spec.Listener.LoadBalancer = nil
		}

	}

	return nil
}
