package basereconciler

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/redhat-cop/operator-utils/pkg/util"
	"github.com/redhat-cop/operator-utils/pkg/util/lockedresourcecontroller"
	"github.com/redhat-cop/operator-utils/pkg/util/lockedresourcecontroller/lockedresource"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var (
	// DefaultExcludedPaths is a list of jsonpaths paths to ignore during reconciliation
	DefaultExcludedPaths []string = []string{
		"/metadata/creationTimestamp",
		"/metadata/deletionGracePeriodSeconds",
		"/metadata/deletionTimestamp",
		"/metadata/finalizers",
		"/metadata/generateName",
		"/metadata/generation",
		"/metadata/managedFields",
		"/metadata/ownerReferences",
		"/metadata/resourceVersion",
		"/metadata/selfLink",
		"/metadata/uid",
		"/status",
	}
	// DeploymentExcludedPaths is a list fo path to ignore for Deployment resources
	DeploymentExcludedPaths []string = []string{
		"/metadata",
		"/status",
		"/spec/progressDeadlineSeconds",
		"/spec/revisionHistoryLimit",
		"/spec/template/metadata/creationTimestamp",
		"/spec/template/spec/dnsPolicy",
		"/spec/template/spec/restartPolicy",
		"/spec/template/spec/schedulerName",
		"/spec/template/spec/securityContext",
		"/spec/template/spec/terminationGracePeriodSeconds",
	}
)

// Reconciler computes a list of resources that it needs to keep in place
type Reconciler struct {
	lockedresourcecontroller.EnforcingReconciler
}

// NewFromManager constructs a new Reconciler from the given manager
func NewFromManager(mgr manager.Manager, recorder record.EventRecorder, clusterWatchers bool) Reconciler {
	return Reconciler{
		EnforcingReconciler: lockedresourcecontroller.NewFromManager(mgr, mgr.GetEventRecorderFor("DiscoveryService"), clusterWatchers),
	}
}

// GeneratorFunction is a function that returns a client.Object
type GeneratorFunction func() client.Object

// LockedResource is a struct that instructs the reconciler how to
// generate and reconcile a resource
type LockedResource struct {
	GeneratorFn  GeneratorFunction
	ExcludePaths []string
}

// GetInstance tries to retrieve the custom resource instance and perform some standard
// tasks like initizalization and cleanup when required.
func (r *Reconciler) GetInstance(ctx context.Context, key types.NamespacedName,
	instance client.Object, finalizer string, log logr.Logger) (ctrl.Result, error) {
	err := r.GetClient().Get(ctx, key, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if util.IsBeingDeleted(instance) {
		if !util.HasFinalizer(instance, finalizer) {
			return ctrl.Result{}, nil
		}
		err := r.ManageCleanUpLogic(instance, log)
		if err != nil {
			log.Error(err, "unable to delete instance")
			return r.ManageError(ctx, instance, err)
		}
		util.RemoveFinalizer(instance, finalizer)
		err = r.GetClient().Update(ctx, instance)
		if err != nil {
			log.Error(err, "unable to update instance")
			return r.ManageError(ctx, instance, err)
		}
		return ctrl.Result{}, nil
	}

	if ok := r.IsInitialized(instance, finalizer); !ok {
		err := r.GetClient().Update(ctx, instance)
		if err != nil {
			log.Error(err, "unable to initialize instance")
			return r.ManageError(ctx, instance, err)
		}
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
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
func (r *Reconciler) ManageCleanUpLogic(instance client.Object, log logr.Logger) error {
	err := r.Terminate(instance, true)
	if err != nil {
		log.Error(err, "unable to terminate locked resources reconciler")
		return err
	}
	return nil
}

// NewLockedResources returns the list of lockedresource.LockedResource that the reconciler needs to enforce
func (r *Reconciler) NewLockedResources(list []LockedResource, owner client.Object) ([]lockedresource.LockedResource, error) {
	resources := []lockedresource.LockedResource{}
	var err error

	for _, res := range list {
		resources, err = add(resources, res.GeneratorFn, res.ExcludePaths, owner, r.GetScheme())
		if err != nil {
			return nil, err
		}
	}
	return resources, nil
}

func add(resources []lockedresource.LockedResource, fn GeneratorFunction, excludedPaths []string,
	owner client.Object, scheme *runtime.Scheme) ([]lockedresource.LockedResource, error) {

	u, err := newUnstructured(fn, owner, scheme)
	if err != nil {
		return nil, err
	}

	res := lockedresource.LockedResource{
		Unstructured:  u,
		ExcludedPaths: excludedPaths,
	}

	return append(resources, res), nil
}

func newUnstructured(fn GeneratorFunction, owner client.Object, scheme *runtime.Scheme) (unstructured.Unstructured, error) {
	o := fn()
	if err := controllerutil.SetControllerReference(owner, o, scheme); err != nil {
		return unstructured.Unstructured{}, err
	}
	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(o)
	if err != nil {
		return unstructured.Unstructured{}, err
	}
	return unstructured.Unstructured{Object: u}, nil
}
