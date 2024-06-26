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
	"strings"

	"github.com/3scale-ops/basereconciler/reconciler"
	"github.com/3scale-ops/basereconciler/util"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/generators/echoapi"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	ctx, logger := r.Logger(ctx, "name", req.Name, "namespace", req.Namespace)
	instance := &saasv1alpha1.EchoAPI{}
	result := r.ManageResourceLifecycle(ctx, req, instance,
		reconciler.WithInMemoryInitializationFunc(util.ResourceDefaulter(instance)))
	if result.ShouldReturn() {
		return result.Values()
	}

	gen := echoapi.NewGenerator(instance.GetName(), instance.GetNamespace(), instance.Spec)

	// Upgrade NLBs managed by nlb-helper-operator
	svc := &corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: gen.GetComponent(), Namespace: req.Namespace}}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(svc), svc); err == nil {
		// 1. Update parent resource with NLB name for NLB adoption
		// 2. Add termination protection
		nlbName := strings.Split(strings.Split(svc.Status.LoadBalancer.Ingress[0].Hostname, ".")[0], "-")[0]
		if instance.Spec.LoadBalancer.LoadBalancerName == nil ||
			*instance.Spec.LoadBalancer.LoadBalancerName != nlbName {
			patch := client.MergeFrom(instance.DeepCopy())
			instance.Spec.LoadBalancer.LoadBalancerName = &nlbName
			instance.Spec.LoadBalancer.DeletionProtection = util.Pointer(true)
			if err := r.Client.Patch(ctx, instance, patch); err != nil {
				return ctrl.Result{}, err
			}
			logger.Info("resource patched", "kind", "EchoAPI", "key", req)
		}

		// 3. Abandon old Service resource
		if svc.GetOwnerReferences() != nil || svc.GetAnnotations()["service.beta.kubernetes.io/aws-load-balancer-type"] != "external" {
			svc.ObjectMeta.OwnerReferences = nil
			svc.ObjectMeta.Annotations["service.beta.kubernetes.io/aws-load-balancer-type"] = "external"
			if err := r.Client.Update(ctx, svc); err != nil {
				return ctrl.Result{}, err
			}
			logger.Info("resource abandoned", "kind", "Service", "key", req)
		}

	} else if !errors.IsNotFound(err) {
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
func (r *EchoAPIReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return reconciler.SetupWithDynamicTypeWatches(r,
		ctrl.NewControllerManagedBy(mgr).
			For(&saasv1alpha1.EchoAPI{}),
	)
}
