package resources

import (
	"context"
	"fmt"

	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var _ basereconciler.Resource = ConfigMapTemplate{}

// ConfigMapsTemplate has methods to generate and reconcile a ConfigMap
type ConfigMapTemplate struct {
	Template  func() *corev1.ConfigMap
	IsEnabled bool
}

// Build returns a ConfigMap resource
func (cmt ConfigMapTemplate) Build(ctx context.Context, cl client.Client) (client.Object, error) {
	return cmt.Template().DeepCopy(), nil
}

// Enabled indicates if the resource should be present or not
func (cmt ConfigMapTemplate) Enabled() bool {
	return cmt.IsEnabled
}

// ResourceReconciler implements a generic reconciler for ConfigMap resources
func (cmt ConfigMapTemplate) ResourceReconciler(ctx context.Context, cl client.Client, obj client.Object) error {
	logger := log.FromContext(ctx, "ResourceReconciler", "ConfigMap")

	needsUpdate := false
	desired := obj.(*corev1.ConfigMap)

	instance := &corev1.ConfigMap{}
	err := cl.Get(ctx, types.NamespacedName{Name: desired.GetName(), Namespace: desired.GetNamespace()}, instance)
	if err != nil {
		if errors.IsNotFound(err) {

			if cmt.Enabled() {
				err = cl.Create(ctx, desired)
				if err != nil {
					return fmt.Errorf("unable to create object: " + err.Error())
				}
				logger.Info("Resource created")
				return nil

			} else {
				return nil
			}
		}

		return err
	}

	/* Delete and return if not enabled */
	if !cmt.Enabled() {
		err := cl.Delete(ctx, instance)
		if err != nil {
			return fmt.Errorf("unable to delete object: " + err.Error())
		}
		logger.Info("Resource deleted")
		return nil
	}

	/* Reconcile metadata */
	if !equality.Semantic.DeepEqual(instance.GetLabels(), desired.GetLabels()) {
		instance.ObjectMeta.Labels = desired.GetLabels()
		needsUpdate = true
	}

	/* Reconcile the data */
	if !equality.Semantic.DeepEqual(instance.Data, desired.Data) {
		instance.Data = desired.Data
		needsUpdate = true
	}

	if needsUpdate {
		err := cl.Update(ctx, instance)
		if err != nil {
			return err
		}
		logger.Info("Resource updated")
	}

	return nil
}
