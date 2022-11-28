package basereconciler

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/redhat-cop/operator-utils/pkg/util"
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
	instance client.Object, finalizer string, cleanupFns []func()) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)

	err := r.Client.Get(ctx, key, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Return and don't requeue
			return &ctrl.Result{}, nil
		}
		return &ctrl.Result{}, err
	}

	if util.IsBeingDeleted(instance) {
		if !util.HasFinalizer(instance, finalizer) {
			return &ctrl.Result{}, nil
		}
		err := r.ManageCleanupLogic(instance, cleanupFns, logger)
		if err != nil {
			logger.Error(err, "unable to delete instance")
			result, err := ctrl.Result{}, err
			return &result, err
		}
		util.RemoveFinalizer(instance, finalizer)
		err = r.Client.Update(ctx, instance)
		if err != nil {
			logger.Error(err, "unable to update instance")
			result, err := ctrl.Result{}, err
			return &result, err
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
func (r *Reconciler) IsInitialized(instance client.Object, finalizer string) bool {
	ok := true
	// ensure the finalizers are removed, they are no longer required
	// as we have dropped usage the lockedresources reconciler
	if util.HasFinalizer(instance, finalizer) {
		util.RemoveFinalizer(instance, finalizer)
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

	}

	return nil
}

// ExtendedObjectList is an extension of client.ObjectList with methods
// to manipulate generically the objects in the list
type ExtendedObjectList interface {
	client.ObjectList
	GetItem(int) client.Object
	CountItems() int
}

// SecretEventHandler returns an EventHandler for the specific ExtendedObjectList
// list object passed as parameter
func (r *Reconciler) SecretEventHandler(ol ExtendedObjectList, logger logr.Logger) handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(
		func(o client.Object) []reconcile.Request {
			if err := r.Client.List(context.TODO(), ol); err != nil {
				logger.Error(err, "unable to retrieve the list of resources")
				return []reconcile.Request{}
			}
			if ol.CountItems() == 0 {
				return []reconcile.Request{}
			}

			key := types.NamespacedName{
				Name:      ol.GetItem(0).GetName(),
				Namespace: ol.GetItem(0).GetNamespace(),
			}
			return []reconcile.Request{{NamespacedName: key}}
		},
	)
}
