package basereconciler

import (
	"context"
	"fmt"
	"hash/fnv"

	secretsmanagerv1alpha1 "github.com/3scale/saas-operator/pkg/apis/secrets-manager/v1alpha1"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-logr/logr"
	"github.com/redhat-cop/operator-utils/pkg/util"
	"github.com/redhat-cop/operator-utils/pkg/util/lockedresourcecontroller"
	"github.com/redhat-cop/operator-utils/pkg/util/lockedresourcecontroller/lockedresource"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	// DefaultExcludedPaths is a list of jsonpaths paths to ignore during reconciliation
	DefaultExcludedPaths []string = []string{".metadata", ".status"}
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

// CalculateSecretHash claculates the hash of a Secret's contents from the SecretDefinition generator function
func (r *Reconciler) CalculateSecretHash(ctx context.Context, fn GeneratorFunction) (string, error) {
	sd := fn().(*secretsmanagerv1alpha1.SecretDefinition)
	key := types.NamespacedName{
		Name:      sd.Spec.Name,
		Namespace: sd.GetNamespace(),
	}
	secret := &corev1.Secret{}
	err := r.GetClient().Get(ctx, key, secret)
	if err != nil {
		if errors.IsNotFound(err) {
			// The secret hasn't been created yet
			return "", nil
		}
		return "", err
	}
	return hash(secret.Data), nil
}

// CalculateConfigMapHash ...
func (r *Reconciler) CalculateConfigMapHash(ctx context.Context, fn GeneratorFunction) string {
	cm := fn().(*corev1.ConfigMap)
	return hash(cm.Data)
}

func hash(o interface{}) string {
	hasher := fnv.New32a()
	hasher.Reset()
	printer := spew.ConfigState{
		Indent:         " ",
		SortKeys:       true,
		DisableMethods: true,
		SpewKeys:       true,
	}
	printer.Fprintf(hasher, "%#v", o)
	return rand.SafeEncodeString(fmt.Sprint(hasher.Sum32()))
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
				logger.Error(err, "unable to retrieve the list of mappingservices")
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
