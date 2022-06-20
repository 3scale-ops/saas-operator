package resources

import (
	"context"
	"fmt"

	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var _ basereconciler.Resource = PodDisruptionBudgetTemplate{}

// PodDisruptionBudget specifies a PodDisruptionBudget resource
type PodDisruptionBudgetTemplate struct {
	Template  func() *policyv1.PodDisruptionBudget
	IsEnabled bool
}

func (pdbt PodDisruptionBudgetTemplate) Build(ctx context.Context, cl client.Client) (client.Object, []string, error) {

	pdb := pdbt.Template()
	pdb.GetObjectKind().SetGroupVersionKind(policyv1.SchemeGroupVersion.WithKind("PodDisruptionBudget"))
	return pdb.DeepCopy(), DefaultExcludedPaths, nil
}

func (pdbt PodDisruptionBudgetTemplate) Enabled() bool {
	return pdbt.IsEnabled
}

func (st PodDisruptionBudgetTemplate) ResourceReconciler(ctx context.Context, cl client.Client, obj client.Object) error {
	logger := log.FromContext(ctx, "ResourceReconciler", "PodDisruptionBudget")

	needsUpdate := false
	desired := obj.(*policyv1.PodDisruptionBudget)

	instance := &policyv1.PodDisruptionBudget{}
	err := cl.Get(ctx, types.NamespacedName{Name: desired.GetName(), Namespace: desired.GetNamespace()}, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			err = cl.Create(ctx, desired)
			if err != nil {
				return fmt.Errorf("unable to create object: " + err.Error())
			}
			logger.Info("Resource created")
			return nil
		}
		return err
	}

	/* Reconcile metadata */
	if !equality.Semantic.DeepEqual(instance.GetAnnotations(), desired.GetAnnotations()) {
		instance.ObjectMeta.Annotations = desired.GetAnnotations()
		needsUpdate = true
	}
	if !equality.Semantic.DeepEqual(instance.GetLabels(), desired.GetLabels()) {
		instance.ObjectMeta.Labels = desired.GetLabels()
		needsUpdate = true
	}

	/* Reconcile the maxUnavaliable and minAvaliable properties */
	if !equality.Semantic.DeepEqual(instance.Spec.MaxUnavailable, desired.Spec.MaxUnavailable) {
		instance.Spec.MaxUnavailable = desired.Spec.MaxUnavailable
		needsUpdate = true
	}
	if !equality.Semantic.DeepEqual(instance.Spec.MinAvailable, desired.Spec.MinAvailable) {
		instance.Spec.MinAvailable = desired.Spec.MinAvailable
		needsUpdate = true
	}

	/* Reconcile label selector */
	if !equality.Semantic.DeepEqual(instance.Spec.Selector, desired.Spec.Selector) {
		instance.Spec.Selector = desired.Spec.Selector
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
