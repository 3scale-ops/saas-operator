package marin3r

import (
	"fmt"
	"strings"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	marin3rDomain            string = "marin3r.3scale.net"
	sidecarEnabledLabelKey   string = marin3rDomain + "/status"
	sidecarEnabledLabelValue string = "enabled"
)

var (
	defaultAnnotations map[string]string = map[string]string{
		"marin3r.3scale.net/shutdown-manager.enabled": "true",
	}
)

// EnableSidecar adds the appropriate labels and annotations for marin3r sidecar
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
		nodeIDAnnotation(spec.NodeID),
		imageAnnotation(spec.EnvoyImage),
		apiVersionAnnotation(spec.EnvoyAPIVersion),
		shtdnmgrPortAnnotation(spec.ShutdownManagerPort),
		shtdnmgrExtraLifecycleHooksAnnotation(spec.ShutdownManagerExtraLifecycleHooks),
		resourcesAnnotations(spec.Resources),
		portsAnnotation(spec.Ports),
		defaultAnnotations,
		spec.ExtraPodAnnotations,
	)

	return &dep
}

func nodeIDAnnotation(nodeID *string) map[string]string {
	if nodeID != nil {
		return map[string]string{marin3rDomain + "/node-id": *nodeID}
	}
	return nil
}

func imageAnnotation(image *string) map[string]string {
	if image != nil {
		return map[string]string{marin3rDomain + "/envoy-image": *image}
	}
	return nil
}

func apiVersionAnnotation(version *string) map[string]string {
	if version != nil {
		return map[string]string{marin3rDomain + "/envoy-api-version": *version}
	}
	return nil
}

func shtdnmgrPortAnnotation(port *uint32) map[string]string {
	if port != nil {
		return map[string]string{marin3rDomain + "/shutdown-manager.port": fmt.Sprintf("%d", *port)}
	}
	return nil
}

func shtdnmgrExtraLifecycleHooksAnnotation(hooks []string) map[string]string {
	if len(hooks) > 0 {
		return map[string]string{marin3rDomain + "/shutdown-manager.extra-lifecycle-hooks": strings.Join(hooks, ",")}
	}
	return nil
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
			annotations[marin3rDomain+"/resources.requests.cpu"] = value.String()
		}
		if value, ok := resources.Requests[corev1.ResourceMemory]; ok {
			annotations[marin3rDomain+"/resources.requests.memory"] = value.String()
		}
	}

	if resources.Limits != nil {
		if value, ok := resources.Limits[corev1.ResourceCPU]; ok {
			annotations[marin3rDomain+"/resources.limits.cpu"] = value.String()
		}
		if value, ok := resources.Limits[corev1.ResourceMemory]; ok {
			annotations[marin3rDomain+"/resources.limits.memory"] = value.String()
		}
	}

	return annotations
}

// podAnnotations generates the annotations value for the marin3r sidecar ports
// annotation
func portsAnnotation(ports []saasv1alpha1.SidecarPort) map[string]string {
	// marin3r syntax for port specification is 'name:port[:protocol]'
	if len(ports) > 0 {
		portSpec := []string{}
		for _, port := range ports {
			portSpec = append(portSpec, strings.Join([]string{port.Name, fmt.Sprintf("%d", port.Port)}, ":"))
		}
		return map[string]string{
			marin3rDomain + "/ports": strings.Join(portSpec, ","),
		}
	}

	return nil
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
