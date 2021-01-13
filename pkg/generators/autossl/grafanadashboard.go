package autossl

import (
	"fmt"

	grafanav1alpha1 "github.com/3scale/saas-operator/pkg/apis/grafana/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GrafanaDashboard returns a basereconciler.GeneratorFunction funtion that will return a GrafanaDashboard
// resource when called
func (opts *Options) GrafanaDashboard() basereconciler.GeneratorFunction {

	return func() client.Object {

		return &grafanav1alpha1.GrafanaDashboard{
			TypeMeta: metav1.TypeMeta{
				Kind:       "GrafanaDashboard",
				APIVersion: grafanav1alpha1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      Component,
				Namespace: opts.Namespace,
				Labels:    opts.labels(),
			},
			Spec: grafanav1alpha1.GrafanaDashboardSpec{
				Name: fmt.Sprintf("%s-%s-%s.json", opts.Namespace, "thresscale", Component),
				Json: "",
			},
		}
	}
}
