package generators

import (
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PDB returns a basereconciler.GeneratorFunction funtion that will return a PDB
// resource when called
func (bo *BaseOptions) PDB(cfg saasv1alpha1.PodDisruptionBudgetSpec) basereconciler.GeneratorFunction {

	return func() client.Object {

		return &policyv1beta1.PodDisruptionBudget{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PodDisruptionBudget",
				APIVersion: policyv1beta1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      bo.GetComponent(),
				Namespace: bo.GetNamespace(),
				Labels:    bo.Labels(),
			},
			Spec: func() policyv1beta1.PodDisruptionBudgetSpec {
				spec := policyv1beta1.PodDisruptionBudgetSpec{
					Selector: bo.Selector(),
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
