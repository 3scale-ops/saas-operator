package basereconciler

import (
	"context"
	"strconv"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	externalsecretsv1beta1 "github.com/3scale/saas-operator/pkg/apis/externalsecrets/v1beta1"
	grafanav1alpha1 "github.com/3scale/saas-operator/pkg/apis/grafana/v1alpha1"
	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-logr/logr"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	operatorutils "github.com/redhat-cop/operator-utils/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var SupportedListTypes = []client.ObjectList{
	&corev1.ServiceAccountList{},
	&corev1.ConfigMapList{},
	&appsv1.DeploymentList{},
	&appsv1.StatefulSetList{},
	&externalsecretsv1beta1.ExternalSecretList{},
	&grafanav1alpha1.GrafanaDashboardList{},
	&autoscalingv2beta2.HorizontalPodAutoscalerList{},
	&policyv1.PodDisruptionBudgetList{},
	&monitoringv1.PodMonitorList{},
}

type Resource interface {
	Build(ctx context.Context, cl client.Client) (client.Object, []string, error)
	Enabled() bool
	ResourceReconciler(context.Context, client.Client, client.Object) error
}

// Reconciler computes a list of resources that it needs to keep in place
type Reconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func NewFromManager(mgr manager.Manager) Reconciler {
	return Reconciler{Client: mgr.GetClient(), Scheme: mgr.GetScheme()}
}

// GetInstance tries to retrieve the custom resource instance and perform some standard
// tasks like initialization and cleanup when required.
func (r *Reconciler) GetInstance(ctx context.Context, key types.NamespacedName,
	instance client.Object, finalizer *string, cleanupFns []func()) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)

	err := r.Client.Get(ctx, key, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Return and don't requeue
			return &ctrl.Result{}, nil
		}
		return &ctrl.Result{}, err
	}

	if operatorutils.IsBeingDeleted(instance) {

		// finalizer logic is only triggered if the controller
		// sets a finalizer, otherwise there's notihng to be done
		if finalizer != nil {

			if !operatorutils.HasFinalizer(instance, *finalizer) {
				return &ctrl.Result{}, nil
			}
			err := r.ManageCleanupLogic(instance, cleanupFns, logger)
			if err != nil {
				logger.Error(err, "unable to delete instance")
				result, err := ctrl.Result{}, err
				return &result, err
			}
			operatorutils.RemoveFinalizer(instance, *finalizer)
			err = r.Client.Update(ctx, instance)
			if err != nil {
				logger.Error(err, "unable to update instance")
				result, err := ctrl.Result{}, err
				return &result, err
			}

		}
		return &ctrl.Result{}, nil
	}

	if ok := r.IsInitialized(instance, finalizer); !ok {
		err := r.Client.Update(ctx, instance)
		if err != nil {
			logger.Error(err, "unable to initialize instance")
			result, err := ctrl.Result{}, err
			return &result, err
		}
		return &ctrl.Result{}, nil
	}
	return nil, nil
}

// IsInitialized can be used to check if instance is correctly initialized.
// Returns false if it isn't.
func (r *Reconciler) IsInitialized(instance client.Object, finalizer *string) bool {
	ok := true
	if finalizer != nil && !operatorutils.HasFinalizer(instance, *finalizer) {
		operatorutils.AddFinalizer(instance, *finalizer)
		ok = false
	}

	// this is temporary code to update existent custom resources:
	//    ensure the finalizers are removed, they are no longer required
	//    for most custom resources as we have dropped usage the lockedresources
	//    reconciler
	if finalizer == nil && operatorutils.HasFinalizer(instance, saasv1alpha1.Finalizer) {
		operatorutils.RemoveFinalizer(instance, saasv1alpha1.Finalizer)
		ok = false
	}
	return ok
}

// ManageCleanupLogic contains finalization logic for the LockedResourcesReconciler
// Functionality can be extended by passing extra cleanup functions
func (r *Reconciler) ManageCleanupLogic(instance client.Object, fns []func(), log logr.Logger) error {

	// Call any cleanup functions passed
	for _, fn := range fns {
		fn()
	}

	return nil
}

// ReconcileOwnedResources handles generalized resource reconcile logic for
// all controllers
func (r *Reconciler) ReconcileOwnedResources(ctx context.Context, owner client.Object, resources []Resource) error {

	managedResources := []types.NamespacedName{}

	for _, res := range resources {

		object, _, err := res.Build(ctx, r.Client)
		if err != nil {
			return err
		}

		if err := controllerutil.SetControllerReference(owner, object, r.Scheme); err != nil {
			return err
		}

		if err := res.ResourceReconciler(ctx, r.Client, object); err != nil {
			return err
		}

		managedResources = append(managedResources, util.ObjectKey(object))
	}

	if value, ok := owner.GetAnnotations()["saas.3scale.net/prune"]; ok {
		prune, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		if !prune {
			return nil
		}
	}

	for _, list := range SupportedListTypes {
		r.PruneOrphaned(ctx, owner, list, managedResources)
	}

	return nil
}

func (r *Reconciler) PruneOrphaned(ctx context.Context, owner client.Object, list client.ObjectList, managed []types.NamespacedName) error {

	err := r.Client.List(ctx, list, client.InNamespace(owner.GetNamespace()))
	if err != nil {
		return err
	}

	for _, obj := range util.GetItems(list) {

		if isOwned(owner, obj) && !operatorutils.IsBeingDeleted(obj) && !isManaged(obj.GetName(), obj.GetNamespace(), managed) {
			err := r.Client.Delete(ctx, obj)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isOwned(owner client.Object, owned client.Object) bool {
	refs := owned.GetOwnerReferences()
	for _, ref := range refs {
		if ref.Kind == owner.GetObjectKind().GroupVersionKind().Kind && ref.Name == owner.GetName() {
			return true
		}
	}
	return false
}

func isManaged(name, namespace string, managed []types.NamespacedName) bool {
	for _, m := range managed {
		if m.Name == name && m.Namespace == namespace {
			return true
		}
	}
	return false
}

// SecretEventHandler returns an EventHandler for the specific client.ObjectList
// list object passed as parameter
func (r *Reconciler) SecretEventHandler(ol client.ObjectList, logger logr.Logger) handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(
		func(o client.Object) []reconcile.Request {
			if err := r.Client.List(context.TODO(), ol); err != nil {
				logger.Error(err, "unable to retrieve the list of resources")
				return []reconcile.Request{}
			}
			items := util.GetItems(ol)
			if len(items) == 0 {
				return []reconcile.Request{}
			}

			return []reconcile.Request{{NamespacedName: util.ObjectKey(items[0])}}
		},
	)
}
