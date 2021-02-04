package basereconciler

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

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
