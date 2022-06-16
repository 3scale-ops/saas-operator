package resources

import (
	"context"

	externalsecretsv1beta1 "github.com/3scale/saas-operator/pkg/apis/externalsecrets/v1beta1"
	grafanav1alpha1 "github.com/3scale/saas-operator/pkg/apis/grafana/v1alpha1"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ basereconciler.Resource = ExternalSecretTemplate{}

// ExternalSecretTemplate specifies a ExternalSecret resource
type ExternalSecretTemplate struct {
	Template  func() *externalsecretsv1beta1.ExternalSecret
	IsEnabled bool
}

func (est ExternalSecretTemplate) Build(ctx context.Context, cl client.Client) (client.Object, []string, error) {

	es := est.Template()
	es.GetObjectKind().SetGroupVersionKind(externalsecretsv1beta1.SchemeGroupVersion.WithKind(externalsecretsv1beta1.ExtSecretKind))
	return es.DeepCopy(), DefaultExcludedPaths, nil
}

func (est ExternalSecretTemplate) Enabled() bool {
	return est.IsEnabled
}

var _ basereconciler.Resource = HorizontalPodAutoscalerTemplate{}

// HorizontalPodAutoscaler specifies a HorizontalPodAutoscaler resource
type HorizontalPodAutoscalerTemplate struct {
	Template  func() *autoscalingv2beta2.HorizontalPodAutoscaler
	IsEnabled bool
}

func (hpat HorizontalPodAutoscalerTemplate) Build(ctx context.Context, cl client.Client) (client.Object, []string, error) {

	hpa := hpat.Template()
	hpa.GetObjectKind().SetGroupVersionKind(autoscalingv2beta2.SchemeGroupVersion.WithKind("HorizontalPodAutoscaler"))
	return hpa.DeepCopy(), DefaultExcludedPaths, nil
}

func (hpat HorizontalPodAutoscalerTemplate) Enabled() bool {
	return hpat.IsEnabled
}

var _ basereconciler.Resource = PodMonitorTemplate{}

// PodMonitor specifies a PodMonitor resource
type PodMonitorTemplate struct {
	Template  func() *monitoringv1.PodMonitor
	IsEnabled bool
}

func (pmt PodMonitorTemplate) Build(ctx context.Context, cl client.Client) (client.Object, []string, error) {

	pm := pmt.Template()
	pm.GetObjectKind().SetGroupVersionKind(monitoringv1.SchemeGroupVersion.WithKind("PodMonitor"))
	return pm.DeepCopy(), DefaultExcludedPaths, nil
}

func (pmt PodMonitorTemplate) Enabled() bool {
	return pmt.IsEnabled
}

var _ basereconciler.Resource = GrafanaDashboardTemplate{}

// GrafanaDashboard specifies a GrafanaDashboard resource
type GrafanaDashboardTemplate struct {
	Template  func() *grafanav1alpha1.GrafanaDashboard
	IsEnabled bool
}

func (gdt GrafanaDashboardTemplate) Build(ctx context.Context, cl client.Client) (client.Object, []string, error) {

	gd := gdt.Template()
	gd.GetObjectKind().SetGroupVersionKind(grafanav1alpha1.GroupVersion.WithKind("GrafanaDashboard"))
	return gd.DeepCopy(), DefaultExcludedPaths, nil
}

func (gdt GrafanaDashboardTemplate) Enabled() bool {
	return gdt.IsEnabled
}

var _ basereconciler.Resource = ConfigMapTemplate{}

// ConfigMaps specifies a ConfigMap resource
type ConfigMapTemplate struct {
	Template  func() *corev1.ConfigMap
	IsEnabled bool
}

func (cmt ConfigMapTemplate) Build(ctx context.Context, cl client.Client) (client.Object, []string, error) {

	cm := cmt.Template()
	cm.GetObjectKind().SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))
	return cm.DeepCopy(), DefaultExcludedPaths, nil
}

func (cmt ConfigMapTemplate) Enabled() bool {
	return cmt.IsEnabled
}
