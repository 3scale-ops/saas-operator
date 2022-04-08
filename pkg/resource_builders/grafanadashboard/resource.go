package grafanadashboard

import (
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	grafanav1alpha1 "github.com/3scale/saas-operator/pkg/apis/grafana/v1alpha1"
	"github.com/3scale/saas-operator/pkg/assets"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// New returns a basereconciler_types.GeneratorFunction function that will return a GrafanaDashboard
// resource when called
func New(key types.NamespacedName, labels map[string]string, cfg saasv1alpha1.GrafanaDashboardSpec,
	template string) func() *grafanav1alpha1.GrafanaDashboard {

	return func() *grafanav1alpha1.GrafanaDashboard {
		data := &struct {
			Namespace string
		}{
			key.Namespace,
		}

		return &grafanav1alpha1.GrafanaDashboard{
			TypeMeta: metav1.TypeMeta{
				Kind:       "GrafanaDashboard",
				APIVersion: grafanav1alpha1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
				Labels: func() map[string]string {
					labels := labels
					labels[*cfg.SelectorKey] = *cfg.SelectorValue
					return labels
				}(),
			},
			Spec: grafanav1alpha1.GrafanaDashboardSpec{
				Json: assets.TemplateAsset(template, data),
			},
		}
	}
}
