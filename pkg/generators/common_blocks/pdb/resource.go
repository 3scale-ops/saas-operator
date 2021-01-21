package pdb

import (
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// New returns a basereconciler.GeneratorFunction funtion that will return a PodDisruptionBudget
// resource when called
func New(key types.NamespacedName, labels map[string]string, selector map[string]string,
	cfg saasv1alpha1.PodDisruptionBudgetSpec) basereconciler.GeneratorFunction {

	return func() client.Object {

		return &policyv1beta1.PodDisruptionBudget{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PodDisruptionBudget",
				APIVersion: policyv1beta1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
				Labels:    labels,
			},
			Spec: func() policyv1beta1.PodDisruptionBudgetSpec {
				spec := policyv1beta1.PodDisruptionBudgetSpec{
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
		}
	}
}
