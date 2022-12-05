package resources

import (
	"context"
	"fmt"

	externalsecretsv1beta1 "github.com/3scale/saas-operator/pkg/apis/externalsecrets/v1beta1"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var _ basereconciler.Resource = ExternalSecretTemplate{}

// ExternalSecretTemplate has methods to generate and reconcile an ExternalSecret
type ExternalSecretTemplate struct {
	Template  func() *externalsecretsv1beta1.ExternalSecret
	IsEnabled bool
}

// Build returns an ExternalSecret resource
func (est ExternalSecretTemplate) Build(ctx context.Context, cl client.Client) (client.Object, error) {
	return est.Template().DeepCopy(), nil
}

// Enabled indicates if the resource should be present or not
func (est ExternalSecretTemplate) Enabled() bool {
	return est.IsEnabled
}

// ResourceReconciler implements a generic reconciler for ExternalSecret resources
func (est ExternalSecretTemplate) ResourceReconciler(ctx context.Context, cl client.Client, obj client.Object) error {
	logger := log.FromContext(ctx, "kind", "ExternalSecret", "resource", obj.GetName())

	needsUpdate := false
	desired := obj.(*externalsecretsv1beta1.ExternalSecret)

	instance := &externalsecretsv1beta1.ExternalSecret{}
	err := cl.Get(ctx, types.NamespacedName{Name: desired.GetName(), Namespace: desired.GetNamespace()}, instance)
	if err != nil {
		if errors.IsNotFound(err) {

			if est.Enabled() {
				err = cl.Create(ctx, desired)
				if err != nil {
					return fmt.Errorf("unable to create object: " + err.Error())
				}
				logger.Info("resource created")
				return nil

			} else {
				return nil
			}
		}

		return err
	}

	/* Delete and return if not enabled */
	if !est.Enabled() {
		err := cl.Delete(ctx, instance)
		if err != nil {
			return fmt.Errorf("unable to delete object: " + err.Error())
		}
		logger.Info("resource deleted")
		return nil
	}

	/* Reconcile metadata */
	if !equality.Semantic.DeepEqual(instance.GetLabels(), desired.GetLabels()) {
		instance.ObjectMeta.Labels = desired.GetLabels()
		needsUpdate = true
	}

	/* Reconcile the spec */
	if !equality.Semantic.DeepEqual(instance.Spec, desired.Spec) {
		instance.Spec = desired.Spec
		needsUpdate = true
	}

	if needsUpdate {
		err := cl.Update(ctx, instance)
		if err != nil {
			return err
		}
		logger.Info("resource updated")
	}

	return nil
}
