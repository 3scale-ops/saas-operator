package basereconciler

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/redhat-cop/operator-utils/pkg/util"
	"github.com/redhat-cop/operator-utils/pkg/util/lockedresourcecontroller"
	"github.com/redhat-cop/operator-utils/pkg/util/lockedresourcecontroller/lockedpatch"
	"github.com/redhat-cop/operator-utils/pkg/util/lockedresourcecontroller/lockedresource"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
}

type ResourceWithCustomReconciler interface {
	Resource
	ResourceReconciler(context.Context, client.Client, client.Object) error
}

// Reconciler computes a list of resources that it needs to keep in place
type Reconciler struct {
	lockedresourcecontroller.EnforcingReconciler
}

// NewFromManager constructs a new Reconciler from the given manager
func NewFromManager(mgr manager.Manager, recorderName string, clusterWatchers bool) Reconciler {
	return Reconciler{
		EnforcingReconciler: lockedresourcecontroller.NewFromManager(mgr, recorderName, clusterWatchers, false),
	}
}

// GetInstance tries to retrieve the custom resource instance and perform some standard
// tasks like initialization and cleanup when required.
func (r *Reconciler) GetInstance(ctx context.Context, key types.NamespacedName,
	instance client.Object, finalizer string, cleanupFns []func()) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)

	err := r.GetClient().Get(ctx, key, instance)
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
		err := r.ManageCleanUpLogic(instance, cleanupFns, logger)
		if err != nil {
			logger.Error(err, "unable to delete instance")
			result, err := r.ManageError(ctx, instance, err)
			return &result, err
		}
		util.RemoveFinalizer(instance, finalizer)
		err = r.GetClient().Update(ctx, instance)
		if err != nil {
			logger.Error(err, "unable to update instance")
			result, err := r.ManageError(ctx, instance, err)
			return &result, err
		}
		return &ctrl.Result{}, nil
	}

	if ok := r.IsInitialized(instance, finalizer); !ok {
		err := r.GetClient().Update(ctx, instance)
		if err != nil {
			logger.Error(err, "unable to initialize instance")
			result, err := r.ManageError(ctx, instance, err)
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
	if !util.HasFinalizer(instance, finalizer) {
		util.AddFinalizer(instance, finalizer)
		ok = false
	}
	return ok
}

// ManageCleanUpLogic contains finalization logic for the LockedResourcesReconciler
// Functionality can be extended by passing extra cleanup functions
func (r *Reconciler) ManageCleanUpLogic(instance client.Object, fns []func(), log logr.Logger) error {

	// Call any extra cleanup functions passed
	for _, fn := range fns {
		fn()
	}

	err := r.Terminate(instance, true)
	if err != nil {
		log.Error(err, "unable to terminate locked resources reconciler")
		return err
	}
	return nil
}

// ReconcileOwnedResources handles generalized resource reconcile logic for
// all controllers
func (r *Reconciler) ReconcileOwnedResources(ctx context.Context, owner client.Object, resources []Resource) error {

	lr := make([]lockedresource.LockedResource, 0, len(resources))

	for _, res := range resources {

		// If the resource implements a custom reconciler, call it and
		// avoid the generic resource processing using operator-utils
		if custom, ok := res.(ResourceWithCustomReconciler); ok {

			object, _, err := res.Build(ctx, r.GetClient())
			if err != nil {
				return err
			}

			if err := controllerutil.SetControllerReference(owner, object, r.GetScheme()); err != nil {
				return err
			}

			if err := custom.ResourceReconciler(ctx, r.GetClient(), object); err != nil {
				return err
			}

		} else {

			if res.Enabled() {

				object, exclude, err := res.Build(ctx, r.GetClient())
				if err != nil {
					return err
				}

				if err := controllerutil.SetControllerReference(owner, object, r.GetScheme()); err != nil {
					return err
				}

				u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
				if err != nil {
					return err
				}

				lr = append(lr, lockedresource.LockedResource{
					Unstructured:  unstructured.Unstructured{Object: u},
					ExcludedPaths: exclude,
				})
			}

		}
	}

	// Call UpdateLockedResources() to reconcile resource types controlled by operator-utils
	if err := r.UpdateLockedResources(ctx, owner, lr, []lockedpatch.LockedPatch{}); err != nil {
		return err
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
			if err := r.GetClient().List(context.TODO(), ol); err != nil {
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
