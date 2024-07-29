package grafanadashboard

import (
	"github.com/3scale-ops/basereconciler/util"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/assets"
	grafanav1beta1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// New returns a basereconciler_types.GeneratorFunction function that will return a GrafanaDashboard
// resource when called
func New(key types.NamespacedName, labels map[string]string, cfg saasv1alpha1.GrafanaDashboardSpec,
	template string) func(client.Object) (*grafanav1beta1.GrafanaDashboard, error) {

	return func(client.Object) (*grafanav1beta1.GrafanaDashboard, error) {
		return &grafanav1beta1.GrafanaDashboard{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
				Labels: func() map[string]string {
					if cfg.SelectorKey != nil && cfg.SelectorValue != nil {
						labels[*cfg.SelectorKey] = *cfg.SelectorValue
					}
					return labels
				}(),
			},
			Spec: grafanav1beta1.GrafanaDashboardSpec{
				AllowCrossNamespaceImport: util.Pointer(true),
				InstanceSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"dashboards": "grafana",
					},
				},
				ResyncPeriod: "5m",
				Json:         assets.TemplateAsset(template, key),
			},
		}, nil
	}
}
