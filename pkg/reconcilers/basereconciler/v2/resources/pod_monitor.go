package resources

import (
	"context"
	"fmt"

	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var _ basereconciler.Resource = PodMonitorTemplate{}

// PodMonitorTemplate has methods to generate and reconcile a PodMonitor
type PodMonitorTemplate struct {
	Template  func() *monitoringv1.PodMonitor
	IsEnabled bool
}

// Build returns a PodMonitor resource
func (pmt PodMonitorTemplate) Build(ctx context.Context, cl client.Client) (client.Object, []string, error) {

	pm := pmt.Template()
	pm.GetObjectKind().SetGroupVersionKind(monitoringv1.SchemeGroupVersion.WithKind("PodMonitor"))
	return pm.DeepCopy(), DefaultExcludedPaths, nil
}

// Enabled indicates if the resource should be present or not
func (pmt PodMonitorTemplate) Enabled() bool {
	return pmt.IsEnabled
}

// ResourceReconciler implements a generic reconciler for PodMonitor resources
func (pmt PodMonitorTemplate) ResourceReconciler(ctx context.Context, cl client.Client, obj client.Object) error {
	logger := log.FromContext(ctx, "ResourceReconciler", "PodMonitor")

	needsUpdate := false
	desired := obj.(*monitoringv1.PodMonitor)

	instance := &monitoringv1.PodMonitor{}
	err := cl.Get(ctx, types.NamespacedName{Name: desired.GetName(), Namespace: desired.GetNamespace()}, instance)
	if err != nil {
		if errors.IsNotFound(err) {

			if pmt.Enabled() {
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
	if !pmt.Enabled() {
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
		logger.Info("Resource updated")
	}

	return nil
}
