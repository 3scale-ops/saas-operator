package resources

import (
	"context"

	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
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

var _ basereconciler.Resource = DeploymentTemplate{}

// DeploymentTemplate specifies a Deployment resource and its rollout triggers
type DeploymentTemplate struct {
	Template        func() *appsv1.Deployment
	RolloutTriggers []RolloutTrigger
	EnforceReplicas bool
	IsEnabled       bool
}

func (dt DeploymentTemplate) Build(ctx context.Context, cl client.Client) (client.Object, []string, error) {

	dep := dt.Template()
	dep.GetObjectKind().SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("Deployment"))

	if err := dt.reconcileDeploymentReplicas(ctx, cl, dep); err != nil {
		return nil, nil, err
	}

	if err := dt.reconcileRolloutTriggers(ctx, cl, dep); err != nil {
		return nil, nil, err
	}

	if dt.EnforceReplicas {
		return dep.DeepCopy(), DeploymentExcludedPaths, nil

	} else {
		return dep.DeepCopy(), append(DeploymentExcludedPaths, "/spec/replicas"), nil
	}
}

func (dt DeploymentTemplate) Enabled() bool {
	return dt.IsEnabled
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
