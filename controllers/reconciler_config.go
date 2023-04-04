package controllers

import (
	basereconciler "github.com/3scale-ops/basereconciler/reconciler"
	marin3rv1alpha1 "github.com/3scale-ops/marin3r/apis/marin3r/v1alpha1"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	externalsecretsv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	grafanav1alpha1 "github.com/grafana-operator/grafana-operator/v4/api/integreatly/v1alpha1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
)

func init() {
	basereconciler.Config.AnnotationsDomain = saasv1alpha1.AnnotationsDomain
	basereconciler.Config.ResourcePruner = true
	basereconciler.Config.ManagedTypes = basereconciler.NewManagedTypes().
		Register(&corev1.ServiceList{}).
		Register(&corev1.ConfigMapList{}).
		Register(&appsv1.DeploymentList{}).
		Register(&appsv1.StatefulSetList{}).
		Register(&externalsecretsv1beta1.ExternalSecretList{}).
		Register(&grafanav1alpha1.GrafanaDashboardList{}).
		Register(&autoscalingv2.HorizontalPodAutoscalerList{}).
		Register(&policyv1.PodDisruptionBudgetList{}).
		Register(&monitoringv1.PodMonitorList{}).
		Register(&marin3rv1alpha1.EnvoyConfigList{})
}
