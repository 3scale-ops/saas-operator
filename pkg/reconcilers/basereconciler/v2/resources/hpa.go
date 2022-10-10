package resources

import (
	"context"
	"fmt"

	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var _ basereconciler.Resource = HorizontalPodAutoscalerTemplate{}

// HorizontalPodAutoscalerTemplate has methods to generate and reconcile a HorizontalPodAutoscaler
type HorizontalPodAutoscalerTemplate struct {
	Template  func() *autoscalingv2beta2.HorizontalPodAutoscaler
	IsEnabled bool
}

// Build returns a HorizontalPodAutoscaler resource
func (hpat HorizontalPodAutoscalerTemplate) Build(ctx context.Context, cl client.Client) (client.Object, []string, error) {

	hpa := hpat.Template()
	hpa.GetObjectKind().SetGroupVersionKind(autoscalingv2beta2.SchemeGroupVersion.WithKind("HorizontalPodAutoscaler"))
	return hpa.DeepCopy(), []string{}, nil
}

// Enabled indicates if the resource should be present or not
func (hpat HorizontalPodAutoscalerTemplate) Enabled() bool {
	return hpat.IsEnabled
}

// ResourceReconciler implements a generic reconciler for HorizontalPodAutoscaler resources
func (hpat HorizontalPodAutoscalerTemplate) ResourceReconciler(ctx context.Context, cl client.Client, obj client.Object) error {
	logger := log.FromContext(ctx, "ResourceReconciler", "HorizontalPodAutoscaler")

	needsUpdate := false
	desired := obj.(*autoscalingv2beta2.HorizontalPodAutoscaler)

	instance := &autoscalingv2beta2.HorizontalPodAutoscaler{}
	err := cl.Get(ctx, types.NamespacedName{Name: desired.GetName(), Namespace: desired.GetNamespace()}, instance)
	if err != nil {
		if errors.IsNotFound(err) {

			if hpat.Enabled() {
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
	if !hpat.Enabled() {
		err := cl.Delete(ctx, instance)
		if err != nil {
			return fmt.Errorf("unable to delete object: " + err.Error())
		}
		logger.Info("Resource deleted")
		return nil
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

	/* Reconcile the ScaleTargetRef, MinReplicas, MaxReplicas and Metrics properties */
	if !equality.Semantic.DeepEqual(instance.Spec.ScaleTargetRef, desired.Spec.ScaleTargetRef) {
		instance.Spec.ScaleTargetRef = desired.Spec.ScaleTargetRef
		needsUpdate = true
	}

	if !equality.Semantic.DeepEqual(instance.Spec.MinReplicas, desired.Spec.MinReplicas) {
		instance.Spec.MinReplicas = desired.Spec.MinReplicas
		needsUpdate = true
	}
	if !equality.Semantic.DeepEqual(instance.Spec.MaxReplicas, desired.Spec.MaxReplicas) {
		instance.Spec.MaxReplicas = desired.Spec.MaxReplicas
		needsUpdate = true
	}

	if !equality.Semantic.DeepEqual(instance.Spec.Metrics, desired.Spec.Metrics) {
		instance.Spec.Metrics = desired.Spec.Metrics
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
