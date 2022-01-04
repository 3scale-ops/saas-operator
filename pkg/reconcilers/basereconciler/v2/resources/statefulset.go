package resources

import (
	"context"

	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	// DeploymentExcludedPaths is a list fo path to ignore for Deployment resources
	StatefulSetExcludedPaths []string = []string{
		"/metadata",
		"/status",
		"/spec/revisionHistoryLimit",
		"/spec/template/spec/dnsPolicy",
		"/spec/template/spec/restartPolicy",
		"/spec/template/spec/schedulerName",
		"/spec/template/spec/securityContext",
		"/spec/template/spec/terminationGracePeriodSeconds",
	}
)

var _ basereconciler.Resource = StatefulSetTemplate{}

// StatefulSet specifies a StatefulSet resource and its rollout triggers
type StatefulSetTemplate struct {
	Template        func() *appsv1.StatefulSet
	RolloutTriggers []RolloutTrigger
	IsEnabled       bool
}

func (sst StatefulSetTemplate) Build(ctx context.Context, cl client.Client) (client.Object, []string, error) {

	ss := sst.Template()
	ss.GetObjectKind().SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("StatefulSet"))

	if err := sst.reconcileRolloutTriggers(ctx, cl, ss); err != nil {
		return nil, nil, err
	}

	return ss, StatefulSetExcludedPaths, nil
}

func (sst StatefulSetTemplate) Enabled() bool {
	return sst.IsEnabled
}

// DeploymentWithRolloutTriggers returns the Deployment modified with the appropriate rollout triggers (annotations)
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
