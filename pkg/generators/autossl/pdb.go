package autossl

import (
	"github.com/3scale/saas-operator/pkg/basereconciler"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PDB returns a basereconciler.GeneratorFunction funtion that will return a PDB
// resource when called
func (opts *Options) PDB() basereconciler.GeneratorFunction {

	return func() client.Object {

		return &policyv1beta1.PodDisruptionBudget{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PodDisruptionBudget",
				APIVersion: policyv1beta1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      Component,
				Namespace: opts.Namespace,
				Labels:    opts.labels(),
			},
			Spec: func() policyv1beta1.PodDisruptionBudgetSpec {
				spec := policyv1beta1.PodDisruptionBudgetSpec{
					Selector: opts.selector(),
				}
				if opts.Spec.PDB.MinAvailable != nil {
					spec.MinAvailable = *&opts.Spec.PDB.MinAvailable
				} else {
					spec.MaxUnavailable = *&opts.Spec.PDB.MaxUnavailable
				}
				return spec
			}(),
		}
	}
}
