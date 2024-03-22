package controllers

import (
	"github.com/3scale-ops/basereconciler/config"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	config.SetAnnotationsDomain(saasv1alpha1.AnnotationsDomain)
	config.EnableResourcePruner()

	config.SetDefaultReconcileConfigForGVK(
		schema.FromAPIVersionAndKind("v1", "Service"),
		config.ReconcileConfigForGVK{
			EnsureProperties: []string{
				"metadata.annotations",
				"metadata.labels",
				"spec.type",
				"spec.ports",
				"spec.selector",
				"spec.clusterIP",
				"spec.clusterIPs",
			},
			IgnoreProperties: []string{
				"metadata.annotations['metallb.universe.tf/ip-allocated-from-pool']",
			},
		})
	config.SetDefaultReconcileConfigForGVK(
		schema.FromAPIVersionAndKind("v1", "ConfigMap"),
		config.ReconcileConfigForGVK{
			EnsureProperties: []string{
				"metadata.annotations",
				"metadata.labels",
				"data",
			},
		})
	config.SetDefaultReconcileConfigForGVK(
		schema.FromAPIVersionAndKind("apps/v1", "Deployment"),
		config.ReconcileConfigForGVK{
			EnsureProperties: []string{
				"metadata.annotations",
				"metadata.labels",
				"spec.minReadySeconds",
				"spec.replicas",
				"spec.selector",
				"spec.strategy",
				"spec.template.metadata.labels",
				"spec.template.metadata.annotations",
				"spec.template.spec",
			},
			IgnoreProperties: []string{
				"metadata.annotations['deployment.kubernetes.io/revision']",
				"spec.template.spec.dnsPolicy",
				"spec.template.spec.schedulerName",
				"spec.template.spec.restartPolicy",
				"spec.template.spec.securityContext",
				"spec.template.spec.containers[*].terminationMessagePath",
				"spec.template.spec.containers[*].terminationMessagePolicy",
				"spec.template.spec.initContainers[*].terminationMessagePath",
				"spec.template.spec.initContainers[*].terminationMessagePolicy",
			},
		})
	config.SetDefaultReconcileConfigForGVK(
		schema.FromAPIVersionAndKind("apps/v1", "StatefulSet"),
		config.ReconcileConfigForGVK{
			EnsureProperties: []string{
				"metadata.annotations",
				"metadata.labels",
				"spec.minReadySeconds",
				"spec.podManagementPolicy",
				"spec.replicas",
				"spec.selector",
				"spec.serviceName",
				"spec.updateStrategy",
				"spec.volumeClaimTemplates",
				"spec.template.metadata.labels",
				"spec.template.metadata.annotations",
				"spec.template.spec",
			},
			IgnoreProperties: []string{
				"metadata.annotations['deployment.kubernetes.io/revision']",
				"spec.template.spec.dnsPolicy",
				"spec.template.spec.restartPolicy",
				"spec.template.spec.schedulerName",
				"spec.template.spec.securityContext",
				"spec.template.spec.containers[*].terminationMessagePath",
				"spec.template.spec.containers[*].terminationMessagePolicy",
				"spec.template.spec.initContainers[*].terminationMessagePath",
				"spec.template.spec.initContainers[*].terminationMessagePolicy",
			},
		})
	config.SetDefaultReconcileConfigForGVK(
		schema.FromAPIVersionAndKind("autoscaling/v2", "HorizontalPodAutoscaler"),
		config.ReconcileConfigForGVK{
			EnsureProperties: []string{
				"metadata.annotations",
				"metadata.labels",
				"spec.scaleTargetRef",
				"spec.minReplicas",
				"spec.maxReplicas",
				"spec.metrics",
			},
		})
	config.SetDefaultReconcileConfigForGVK(
		schema.FromAPIVersionAndKind("policy/v1", "PodDisruptionBudget"),
		config.ReconcileConfigForGVK{
			EnsureProperties: []string{
				"metadata.annotations",
				"metadata.labels",
				"spec.maxUnavailable",
				"spec.minAvailable",
				"spec.selector",
			},
		})
	config.SetDefaultReconcileConfigForGVK(
		schema.FromAPIVersionAndKind("external-secrets.io/v1beta1", "ExternalSecret"),
		config.ReconcileConfigForGVK{
			EnsureProperties: []string{
				"metadata.annotations",
				"metadata.labels",
				"spec",
			},
			IgnoreProperties: []string{
				"spec.data[*].remoteRef.metadataPolicy",
			},
		})
	config.SetDefaultReconcileConfigForGVK(
		schema.FromAPIVersionAndKind("pipeline/v1beta1", "Task"),
		config.ReconcileConfigForGVK{
			EnsureProperties: []string{
				"metadata.annotations",
				"metadata.labels",
				"spec.displayName",
				"spec.description",
				"spec.params",
				"spec.steps",
				"spec.stepTemplate",
				"spec.volumes",
				"spec.sidecars",
				"spec.workspaces",
				"spec.results",
			},
		})
	config.SetDefaultReconcileConfigForGVK(
		schema.FromAPIVersionAndKind("pipeline/v1beta1", "Pipeline"),
		config.ReconcileConfigForGVK{
			EnsureProperties: []string{
				"metadata.annotations",
				"metadata.labels",
				"spec.displayName",
				"spec.description",
				"spec.params",
				"spec.tasks",
				"spec.workspaces",
				"spec.results",
				"spec.finally",
			},
		})
	// default config for any GVK not explicitely declared in the config
	config.SetDefaultReconcileConfigForGVK(
		schema.GroupVersionKind{},
		config.ReconcileConfigForGVK{
			EnsureProperties: []string{
				"metadata.annotations",
				"metadata.labels",
				"spec",
			},
		})
}
