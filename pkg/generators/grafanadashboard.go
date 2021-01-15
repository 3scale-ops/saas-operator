package generators

import (
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	grafanav1alpha1 "github.com/3scale/saas-operator/pkg/apis/grafana/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GrafanaDashboard returns a basereconciler.GeneratorFunction funtion that will return a GrafanaDashboard
// resource when called
func (bo *BaseOptions) GrafanaDashboard(cfg saasv1alpha1.GrafanaDashboardSpec, dashboard []byte) basereconciler.GeneratorFunction {

	return func() client.Object {

		return &grafanav1alpha1.GrafanaDashboard{
			TypeMeta: metav1.TypeMeta{
				Kind:       "GrafanaDashboard",
				APIVersion: grafanav1alpha1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      bo.GetComponent(),
				Namespace: bo.GetNamespace(),
				Labels: func() map[string]string {
					labels := bo.Labels()
					labels[*cfg.SelectorKey] = *cfg.SelectorValue
					return labels
				}(),
			},
			Spec: grafanav1alpha1.GrafanaDashboardSpec{
				Name: fmt.Sprintf("%s-%s-%s.json", bo.GetNamespace(), "threescale", bo.GetComponent()),
				Json: "",
			},
		}
	}
}
