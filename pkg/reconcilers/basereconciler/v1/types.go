package basereconciler

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GeneratorFunction is a function that returns a client.Object
type GeneratorFunction func() client.Object

// LockedResource is a struct that instructs the reconciler how to
// generate and reconcile a resource
type LockedResource struct {
	GeneratorFn  GeneratorFunction
	ExcludePaths []string
}

// ExtendedObjectList is an extension of client.ObjectList with methods
// to manipulate generically the objects in the list
type ExtendedObjectList interface {
	client.ObjectList
	GetItem(int) client.Object
	CountItems() int
}
