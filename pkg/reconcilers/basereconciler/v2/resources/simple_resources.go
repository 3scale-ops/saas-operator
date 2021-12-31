package resources

import (
	"context"

	grafanav1alpha1 "github.com/3scale/saas-operator/pkg/apis/grafana/v1alpha1"
	secretsmanagerv1alpha1 "github.com/3scale/saas-operator/pkg/apis/secrets-manager/v1alpha1"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ basereconciler.Resource = SecretDefinitionTemplate{}

// SecretDefinition specifies a SecretDefinition resource
type SecretDefinitionTemplate struct {
	Template  func() *secretsmanagerv1alpha1.SecretDefinition
	IsEnabled bool
}

func (sdt SecretDefinitionTemplate) Build(ctx context.Context, cl client.Client) (client.Object, []string, error) {

	sd := sdt.Template()
	sd.GetObjectKind().SetGroupVersionKind(secretsmanagerv1alpha1.GroupVersion.WithKind("SecretDefinition"))
	return sd, DefaultExcludedPaths, nil
}

func (sdt SecretDefinitionTemplate) Enabled() bool {
	return sdt.IsEnabled
}

var _ basereconciler.Resource = PodDisruptionBudgetTemplate{}

// PodDisruptionBudget specifies a PodDisruptionBudget resource
type PodDisruptionBudgetTemplate struct {
	Template  func() *policyv1beta1.PodDisruptionBudget
	IsEnabled bool
}

func (pdbt PodDisruptionBudgetTemplate) Build(ctx context.Context, cl client.Client) (client.Object, []string, error) {

	pdb := pdbt.Template()
	pdb.GetObjectKind().SetGroupVersionKind(policyv1beta1.SchemeGroupVersion.WithKind("PodDisruptionBudget"))
	return pdb, DefaultExcludedPaths, nil
}

func (pdbt PodDisruptionBudgetTemplate) Enabled() bool {
	return pdbt.IsEnabled
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
	return hpa, DefaultExcludedPaths, nil
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
	return pm, DefaultExcludedPaths, nil
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
	gd.GetObjectKind().SetGroupVersionKind(grafanav1alpha1.SchemeGroupVersion.WithKind("GrafanaDashboard"))
	return gd, DefaultExcludedPaths, nil
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
	return cm, DefaultExcludedPaths, nil
}

func (cmt ConfigMapTemplate) Enabled() bool {
	return cmt.IsEnabled
}
