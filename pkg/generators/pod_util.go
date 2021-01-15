package generators

import (
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// HTTPProbe returns an HTTP corev1.Probe struct
func (bo *BaseOptions) HTTPProbe(path string, port intstr.IntOrString, scheme corev1.URIScheme, cfg saasv1alpha1.HTTPProbeSpec) *corev1.Probe {
	if cfg.IsDeactivated() {
		return nil
	}
	return &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   path,
				Port:   port,
				Scheme: scheme,
			},
		},
		InitialDelaySeconds: *cfg.InitialDelaySeconds,
		TimeoutSeconds:      *cfg.TimeoutSeconds,
		PeriodSeconds:       *cfg.PeriodSeconds,
		SuccessThreshold:    *cfg.SuccessThreshold,
		FailureThreshold:    *cfg.FailureThreshold,
	}
}

// Affinity returns a corev1.Affinity struct
func (bo *BaseOptions) Affinity() *corev1.Affinity {
	return &corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
				{
					Weight: 100,
					PodAffinityTerm: corev1.PodAffinityTerm{
						TopologyKey:   corev1.LabelHostname,
						LabelSelector: bo.Selector(),
					},
				},
				{
					Weight: 99,
					PodAffinityTerm: corev1.PodAffinityTerm{
						TopologyKey:   corev1.LabelTopologyZone,
						LabelSelector: bo.Selector(),
					},
				},
			},
		},
	}
}

// ContainerPortTCP returns a TCP corev1.ContainerPort
func (bo *BaseOptions) ContainerPortTCP(name string, port int32) corev1.ContainerPort {
	return corev1.ContainerPort{
		Name:          name,
		ContainerPort: port,
		Protocol:      corev1.ProtocolTCP,
	}
}

// ContainerPorts returns a list of corev1.ContainerPort
func (bo *BaseOptions) ContainerPorts(ports ...corev1.ContainerPort) []corev1.ContainerPort {
	list := []corev1.ContainerPort{}
	return append(list, ports...)
}
