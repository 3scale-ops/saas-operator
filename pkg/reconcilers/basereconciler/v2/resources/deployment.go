package resources

import (
	"context"
	"fmt"

	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-test/deep"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var _ basereconciler.Resource = DeploymentTemplate{}

// DeploymentTemplate specifies a Deployment resource and its rollout triggers
type DeploymentTemplate struct {
	Template        func() *appsv1.Deployment
	RolloutTriggers []RolloutTrigger
	EnforceReplicas bool
	IsEnabled       bool
}

func (dt DeploymentTemplate) Build(ctx context.Context, cl client.Client) (client.Object, error) {

	dep := dt.Template()

	if err := dt.reconcileDeploymentReplicas(ctx, cl, dep); err != nil {
		return nil, err
	}

	if err := dt.reconcileRolloutTriggers(ctx, cl, dep); err != nil {
		return nil, err
	}

	return dep.DeepCopy(), nil
}

func (dt DeploymentTemplate) Enabled() bool {
	return dt.IsEnabled
}

func (dep DeploymentTemplate) ResourceReconciler(ctx context.Context, cl client.Client, obj client.Object) error {
	logger := log.FromContext(ctx, "kind", "Deployment", "resource", obj.GetName())

	needsUpdate := false
	desired := obj.(*appsv1.Deployment)

	instance := &appsv1.Deployment{}
	err := cl.Get(ctx, types.NamespacedName{Name: desired.GetName(), Namespace: desired.GetNamespace()}, instance)
	if err != nil {
		if errors.IsNotFound(err) {

			if dep.Enabled() {
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
	if !dep.Enabled() {
		err := cl.Delete(ctx, instance)
		if err != nil {
			return fmt.Errorf("unable to delete object: " + err.Error())
		}
		logger.Info("resource deleted")
		return nil
	}

	/* Reconcile metadata */

	desired.ObjectMeta.Annotations = util.MergeMaps(
		map[string]string{},
		desired.GetAnnotations(),
		map[string]string{"deployment.kubernetes.io/revision": instance.GetAnnotations()["deployment.kubernetes.io/revision"]},
	)

	if !equality.Semantic.DeepEqual(instance.GetAnnotations(), desired.GetAnnotations()) {
		logger.Info("resource update required due to differences in metadata.annotations.")
		logger.V(1).Info(
			fmt.Sprintf("metadata.annotations differences: %s",
				deep.Equal(instance.GetAnnotations(), desired.GetAnnotations())),
		)
		instance.ObjectMeta.Annotations = desired.GetAnnotations()
		needsUpdate = true
	}
	if !equality.Semantic.DeepEqual(instance.GetLabels(), desired.GetLabels()) {
		logger.Info("resource update required due to differences in metadata.labels.")
		logger.V(1).Info(
			fmt.Sprintf("metadata.labels differences: %s",
				deep.Equal(instance.GetLabels(), desired.GetLabels())),
		)
		instance.ObjectMeta.Labels = desired.GetLabels()
		needsUpdate = true
	}

	/* Reconcile the MinReadySeconds */
	if !equality.Semantic.DeepEqual(instance.Spec.MinReadySeconds, desired.Spec.MinReadySeconds) {
		logger.Info("resource update required due to differences in spec.minReadySeconds.")
		logger.V(1).Info(
			fmt.Sprintf("spec.minReadySeconds differences: %s",
				deep.Equal(instance.Spec.MinReadySeconds, desired.Spec.MinReadySeconds)),
		)
		instance.Spec.MinReadySeconds = desired.Spec.MinReadySeconds
		needsUpdate = true
	}

	/* Reconcile the Replicas */
	if !equality.Semantic.DeepEqual(instance.Spec.Replicas, desired.Spec.Replicas) {
		logger.Info("resource update required due to differences in spec.replicas.")
		logger.V(1).Info(
			fmt.Sprintf("spec.replicas differences: %s",
				deep.Equal(instance.Spec.Replicas, desired.Spec.Replicas)),
		)
		instance.Spec.Replicas = desired.Spec.Replicas
		needsUpdate = true
	}

	/* Reconcile the Selector */
	if !equality.Semantic.DeepEqual(instance.Spec.Selector, desired.Spec.Selector) {
		logger.Info("resource update required due to differences in spec.selector.")
		logger.V(1).Info(
			fmt.Sprintf("spec.selector differences: %s",
				deep.Equal(instance.Spec.Selector, desired.Spec.Selector)),
		)
		instance.Spec.Selector = desired.Spec.Selector
		needsUpdate = true
	}

	/* Reconcile the Strategy */
	if !equality.Semantic.DeepEqual(instance.Spec.Strategy, desired.Spec.Strategy) {
		logger.Info("resource update required due to differences in spec.strategy.")
		logger.V(1).Info(
			fmt.Sprintf("spec.strategy differences: %s",
				deep.Equal(instance.Spec.Strategy, desired.Spec.Strategy)),
		)
		instance.Spec.Strategy = desired.Spec.Strategy
		needsUpdate = true
	}

	/* Reconcile the Template Labels */
	if !equality.Semantic.DeepEqual(
		instance.Spec.Template.ObjectMeta.Labels, desired.Spec.Template.ObjectMeta.Labels) {
		logger.Info("resource update required due to differences in spec.template.metadata.labels.")
		logger.V(1).Info(
			fmt.Sprintf("spec.template.metadata.labels differences: %s",
				deep.Equal(instance.Spec.Template.ObjectMeta.Labels, desired.Spec.Template.ObjectMeta.Labels)),
		)
		instance.Spec.Template.ObjectMeta.Labels = desired.Spec.Template.ObjectMeta.Labels
		needsUpdate = true
	}

	/* Reconcile the Template Annotations */
	if !equality.Semantic.DeepEqual(
		instance.Spec.Template.ObjectMeta.Annotations, desired.Spec.Template.ObjectMeta.Annotations) {
		logger.Info("resource update required due differences in spec.template.metadata.annotations.")
		logger.V(1).Info(
			fmt.Sprintf("spec.template.metadata.annotations differences: %s",
				deep.Equal(instance.Spec.Template.ObjectMeta.Annotations, desired.Spec.Template.ObjectMeta.Annotations)),
		)
		instance.Spec.Template.ObjectMeta.Annotations = desired.Spec.Template.ObjectMeta.Annotations
		needsUpdate = true
	}

	/* Inherit some values usually defaulted by the cluster if not defined on the template */
	if desired.Spec.Template.Spec.DNSPolicy == "" {
		desired.Spec.Template.Spec.DNSPolicy = instance.Spec.Template.Spec.DNSPolicy
	}
	if desired.Spec.Template.Spec.SchedulerName == "" {
		desired.Spec.Template.Spec.SchedulerName = instance.Spec.Template.Spec.SchedulerName
	}

	/* Reconcile the Template Spec */
	if !equality.Semantic.DeepEqual(instance.Spec.Template.Spec, desired.Spec.Template.Spec) {
		logger.Info("resource update required due to differences in spec.template.spec.")
		logger.V(1).Info(
			fmt.Sprintf("spec.template.spec differences: %s",
				deep.Equal(instance.Spec.Template.Spec, desired.Spec.Template.Spec)),
		)
		instance.Spec.Template.Spec = desired.Spec.Template.Spec
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

// reconcileDeploymentReplicas reconciles the number of replicas of a Deployment
func (dt DeploymentTemplate) reconcileDeploymentReplicas(ctx context.Context, cl client.Client, dep *appsv1.Deployment) error {

	if dt.EnforceReplicas {
		// Let the value in the template
		// override the runtime value
		return nil
	}

	key := types.NamespacedName{
		Name:      dep.GetName(),
		Namespace: dep.GetNamespace(),
	}
	instance := &appsv1.Deployment{}
	err := cl.Get(ctx, key, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// override the value in the template with the
	// runtime value
	dep.Spec.Replicas = instance.Spec.Replicas
	return nil
}

// reconcileRolloutTriggers modifies the Deployment with the appropriate rollout triggers (annotations)
func (dt DeploymentTemplate) reconcileRolloutTriggers(ctx context.Context, cl client.Client, dep *appsv1.Deployment) error {

	if dep.Spec.Template.ObjectMeta.Annotations == nil {
		dep.Spec.Template.ObjectMeta.Annotations = map[string]string{}
	}

	for _, trigger := range dt.RolloutTriggers {
		hash, err := trigger.GetHash(ctx, cl, dep.GetNamespace())
		if err != nil {
			return err
		}
		dep.Spec.Template.ObjectMeta.Annotations[trigger.GetAnnotationKey()] = hash
	}

	return nil
}
