package pdb

import (
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// New returns a basereconciler_types.GeneratorFunction function that will return a PodDisruptionBudget
// resource when called
func New(key types.NamespacedName, labels map[string]string, selector map[string]string,
	cfg saasv1alpha1.PodDisruptionBudgetSpec) func(client.Object) (*policyv1.PodDisruptionBudget, error) {

	return func(client.Object) (*policyv1.PodDisruptionBudget, error) {

		return &policyv1.PodDisruptionBudget{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
				Labels:    labels,
			},
			Spec: func() policyv1.PodDisruptionBudgetSpec {
				spec := policyv1.PodDisruptionBudgetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: selector,
					},
				}
				if cfg.MinAvailable != nil {
					spec.MinAvailable = cfg.MinAvailable
				} else {
					spec.MaxUnavailable = cfg.MaxUnavailable
				}
				return spec
			}(),
		}, nil
	}
}
