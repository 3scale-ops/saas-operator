package resources

import (
	"context"
	"fmt"

	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	"github.com/go-test/deep"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var _ basereconciler.Resource = StatefulSetTemplate{}

// StatefulSet specifies a StatefulSet resource and its rollout triggers
type StatefulSetTemplate struct {
	Template        func() *appsv1.StatefulSet
	RolloutTriggers []RolloutTrigger
	IsEnabled       bool
}

func (sst StatefulSetTemplate) Build(ctx context.Context, cl client.Client) (client.Object, error) {

	ss := sst.Template()

	if err := sst.reconcileRolloutTriggers(ctx, cl, ss); err != nil {
		return nil, err
	}

	return ss.DeepCopy(), nil
}

func (sst StatefulSetTemplate) Enabled() bool {
	return sst.IsEnabled
}

func (sts StatefulSetTemplate) ResourceReconciler(ctx context.Context, cl client.Client, obj client.Object) error {
	logger := log.FromContext(ctx, "kind", "StatefulSet", "resource", obj.GetName())

	needsUpdate := false
	desired := obj.(*appsv1.StatefulSet)

	instance := &appsv1.StatefulSet{}
	err := cl.Get(ctx, types.NamespacedName{Name: desired.GetName(), Namespace: desired.GetNamespace()}, instance)
	if err != nil {
		if errors.IsNotFound(err) {

			if sts.Enabled() {
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
	if !sts.Enabled() {
		err := cl.Delete(ctx, instance)
		if err != nil {
			return fmt.Errorf("unable to delete object: " + err.Error())
		}
		logger.Info("resource deleted")
		return nil
	}

	/* Reconcile metadata */
	if !equality.Semantic.DeepEqual(instance.GetAnnotations(), desired.GetAnnotations()) {
		logger.Info("resource update required due to differences in Annotations.")
		logger.V(1).Info(
			fmt.Sprintf("metadata.annotations differences: %s",
				deep.Equal(instance.GetAnnotations(), desired.GetAnnotations())),
		)
		instance.ObjectMeta.Annotations = desired.GetAnnotations()
		needsUpdate = true
	}
	if !equality.Semantic.DeepEqual(instance.GetLabels(), desired.GetLabels()) {
		logger.Info("resource update required due to differences in Labels.")
		logger.V(1).Info(
			fmt.Sprintf("metadala.labels differences: %s",
				deep.Equal(instance.GetLabels(), desired.GetLabels())),
		)
		instance.ObjectMeta.Labels = desired.GetLabels()
		needsUpdate = true
	}

	/* Reconcile the MinReadySeconds */
	if !equality.Semantic.DeepEqual(instance.Spec.MinReadySeconds, desired.Spec.MinReadySeconds) {
		logger.Info("resource update required due to differences in Spec.MinReadySeconds.")
		logger.V(1).Info(
			fmt.Sprintf("spec.minReadySeconds differences: %s",
				deep.Equal(instance.Spec.MinReadySeconds, desired.Spec.MinReadySeconds)),
		)
		instance.Spec.MinReadySeconds = desired.Spec.MinReadySeconds
		needsUpdate = true
	}

	/* Reconcile the PersistentVolumeClaimRetentionPolicy */
	if !equality.Semantic.DeepEqual(instance.Spec.PersistentVolumeClaimRetentionPolicy, desired.Spec.PersistentVolumeClaimRetentionPolicy) {
		logger.Info("resource update required due to differences in Spec.PersistentVolumeClaimRetentionPolicy.")
		logger.V(1).Info(
			fmt.Sprintf("spec.persistentVolumeClaimRetentionPolicy differences: %s",
				deep.Equal(instance.Spec.PersistentVolumeClaimRetentionPolicy, desired.Spec.PersistentVolumeClaimRetentionPolicy)),
		)
		instance.Spec.PersistentVolumeClaimRetentionPolicy = desired.Spec.PersistentVolumeClaimRetentionPolicy
		needsUpdate = true
	}

	/* Reconcile the Replicas */
	if !equality.Semantic.DeepEqual(instance.Spec.Replicas, desired.Spec.Replicas) {
		logger.Info("resource update required due to differences in Spec.Replicas.")
		logger.V(1).Info(
			fmt.Sprintf("spec.replicas differences: %s",
				deep.Equal(instance.Spec.Replicas, desired.Spec.Replicas)),
		)
		instance.Spec.Replicas = desired.Spec.Replicas
		needsUpdate = true
	}

	/* Reconcile the Selector */
	if !equality.Semantic.DeepEqual(instance.Spec.Selector, desired.Spec.Selector) {
		logger.Info("resource update required due to differences in Spec.Selector.")
		logger.V(1).Info(
			fmt.Sprintf("spec.selector differences: %s",
				deep.Equal(instance.Spec.Selector, desired.Spec.Selector)),
		)
		instance.Spec.Selector = desired.Spec.Selector
		needsUpdate = true
	}

	/* Reconcile the ServiceName */
	if !equality.Semantic.DeepEqual(instance.Spec.ServiceName, desired.Spec.ServiceName) {
		logger.Info("resource update required due to differences in Spec.ServiceName.")
		logger.V(1).Info(
			fmt.Sprintf("spec.serviceName differences: %s",
				deep.Equal(instance.Spec.ServiceName, desired.Spec.ServiceName)),
		)
		instance.Spec.ServiceName = desired.Spec.ServiceName
		needsUpdate = true
	}

	/* Reconcile the Template Labels */
	if !equality.Semantic.DeepEqual(
		instance.Spec.Template.ObjectMeta.Labels, desired.Spec.Template.ObjectMeta.Labels) {
		logger.Info("resource update required due to differences in Spec.Template.ObjectMeta.Labels.")
		logger.V(1).Info(
			fmt.Sprintf("spec.template.metadata.labels differences: %s",
				deep.Equal(instance.Spec.Template.ObjectMeta.Labels, desired.Spec.Template.ObjectMeta.Labels)),
		)
		instance.Spec.Template.ObjectMeta.Labels = desired.Spec.Template.ObjectMeta.Labels
		needsUpdate = true
	}

	/* Reconcile the Template Spec */

	if desired.Spec.Template.Spec.SchedulerName == "" {
		desired.Spec.Template.Spec.SchedulerName = instance.Spec.Template.Spec.SchedulerName
	}
	if desired.Spec.Template.Spec.DNSPolicy == "" {
		desired.Spec.Template.Spec.DNSPolicy = instance.Spec.Template.Spec.DNSPolicy
	}

	if !equality.Semantic.DeepEqual(instance.Spec.Template.Spec, desired.Spec.Template.Spec) {
		logger.Info("resource update required due to differences in Spec.Template.Spec.")
		logger.V(1).Info(
			fmt.Sprintf("spec.template.spec differences: %s",
				deep.Equal(instance.Spec.Template.Spec, desired.Spec.Template.Spec)),
		)
		instance.Spec.Template.Spec = desired.Spec.Template.Spec
		needsUpdate = true
	}

	/* Reconcile the UpdateStrategy */
	if !equality.Semantic.DeepEqual(instance.Spec.UpdateStrategy, desired.Spec.UpdateStrategy) {
		logger.Info("resource update required due to differences in Spec.UpdateStrategy.")
		logger.V(1).Info(
			fmt.Sprintf("spec.updateStrategy differences: %s",
				deep.Equal(instance.Spec.UpdateStrategy, desired.Spec.UpdateStrategy)),
		)
		instance.Spec.UpdateStrategy = desired.Spec.UpdateStrategy
		needsUpdate = true
	}

	/* Reconcile the VolumeClaimTemplates */
	if !equality.Semantic.DeepEqual(instance.Spec.VolumeClaimTemplates, desired.Spec.VolumeClaimTemplates) {
		logger.Info("resource update required due to differences in Spec.VolumeClaimTemplates.")
		logger.V(1).Info(
			fmt.Sprintf("spec.volumeClaimTemplates differences: %s",
				deep.Equal(instance.Spec.VolumeClaimTemplates, desired.Spec.VolumeClaimTemplates)),
		)
		instance.Spec.VolumeClaimTemplates = desired.Spec.VolumeClaimTemplates
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

// reconcileRolloutTriggers modifies the StatefulSet with the appropriate rollout triggers (annotations)
func (sst StatefulSetTemplate) reconcileRolloutTriggers(ctx context.Context, cl client.Client, ss *appsv1.StatefulSet) error {

	if ss.Spec.Template.ObjectMeta.Annotations == nil {
		ss.Spec.Template.ObjectMeta.Annotations = map[string]string{}
	}

	for _, trigger := range sst.RolloutTriggers {
		hash, err := trigger.GetHash(ctx, cl, ss.GetNamespace())
		if err != nil {
			return err
		}
		ss.Spec.Template.ObjectMeta.Annotations[trigger.GetAnnotationKey()] = hash
	}

	return nil
}
