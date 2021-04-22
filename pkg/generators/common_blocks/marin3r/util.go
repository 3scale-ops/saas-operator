package marin3r

import (
	"fmt"
	"strings"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	sidecarEnabledLabelKey   string = "marin3r.3scale.net/status"
	sidecarEnabledLabelValue string = "enabled"
)

var (
	defaultAnnotations map[string]string = map[string]string{
		"marin3r.3scale.net/shutdown-manager.enabled": "true",
	}
)

// EnableSidecar adds the apporopriates labels and annotations for marin3r sidecar
// injection to work for this Deployment
func EnableSidecar(dep appsv1.Deployment, spec saasv1alpha1.Marin3rSidecarSpec) *appsv1.Deployment {

	if dep.Spec.Template.ObjectMeta.Labels == nil {
		dep.Spec.Template.ObjectMeta.Labels = map[string]string{}
	}
	if dep.Spec.Template.ObjectMeta.Annotations == nil {
		dep.Spec.Template.ObjectMeta.Annotations = map[string]string{}
	}
	dep.Spec.Template.ObjectMeta.Labels[sidecarEnabledLabelKey] = sidecarEnabledLabelValue
	dep.Spec.Template.ObjectMeta.Annotations = mergeMaps(
		dep.Spec.Template.ObjectMeta.Annotations,
		resourcesAnnotations(spec.Resources),
		portsAnnotation(spec.Ports),
		defaultAnnotations,
		spec.ExtraPodAnnotations,
	)

	return &dep
}

// resourcesAnnotations generates the corresponding annotations for marin3r sidecar resources
// requests configuration
func resourcesAnnotations(resources *saasv1alpha1.ResourceRequirementsSpec) map[string]string {
	annotations := map[string]string{}
	if resources == nil {
		return annotations
	}

	if resources.Requests != nil {
		if value, ok := resources.Requests[corev1.ResourceCPU]; ok {
			annotations["marin3r.3scale.net/resources.requests.cpu"] = value.String()
		}
		if value, ok := resources.Requests[corev1.ResourceMemory]; ok {
			annotations["marin3r.3scale.net/resources.requests.memory"] = value.String()
		}
	}

	if resources.Limits != nil {
		if value, ok := resources.Limits[corev1.ResourceCPU]; ok {
			annotations["marin3r.3scale.net/resources.limits.cpu"] = value.String()
		}
		if value, ok := resources.Limits[corev1.ResourceMemory]; ok {
			annotations["marin3r.3scale.net/resources.limits.memory"] = value.String()
		}
	}

	return annotations
}

// podAnnotations generates the annotations value for the marin3r sidecar ports
// annotation
func portsAnnotation(ports []saasv1alpha1.SidecarPort) map[string]string {
	// marin3r syntax for port specification is 'name:port[:protocol]'
	portSpec := []string{}
	for _, port := range ports {
		portSpec = append(portSpec, strings.Join([]string{port.Name, fmt.Sprintf("%d", port.Port)}, ":"))
	}
	return map[string]string{
		"marin3r.3scale.net/ports": strings.Join(portSpec, ","),
	}
}

// mergeMaps merges two maps. B overrides A if keys collide.
func mergeMaps(base map[string]string, merges ...map[string]string) map[string]string {
	for _, m := range merges {
		for key, value := range m {
			base[key] = value
		}
	}
	return base
}
