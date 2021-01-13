package autossl

import (
	"fmt"
	"strings"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	Component string = "autossl"
)

// Options configures the generators for AutoSSL
type Options struct {
	InstanceName string
	Namespace    string
	Spec         saasv1alpha1.AutoSSLSpec
}

func (opts *Options) labels() map[string]string {
	return map[string]string{
		"app":     Component,
		"part-of": "3scale-saas",
	}
}

func (opts *Options) labelsWithSelector() map[string]string {
	return map[string]string{
		"app":        Component,
		"part-of":    "3scale-saas",
		"deployment": Component,
	}
}

func (opts *Options) selector() *metav1.LabelSelector {
	return &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"deployment": Component,
		},
	}
}

func (opts *Options) serviceAnnotations() map[string]string {
	annotations := map[string]string{
		"external-dns.alpha.kubernetes.io/hostname":                                      strings.Join(opts.Spec.Endpoint.DNS, ","),
		"service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled": fmt.Sprintf("%t", *opts.Spec.LoadBalancer.CrossZoneLoadBalancingEnabled),
		"service.beta.kubernetes.io/aws-load-balancer-connection-draining-enabled":       fmt.Sprintf("%t", *opts.Spec.LoadBalancer.ConnectionDrainingEnabled),
		"service.beta.kubernetes.io/aws-load-balancer-connection-draining-timeout":       fmt.Sprintf("%d", *opts.Spec.LoadBalancer.ConnectionDrainingTimeout),
		"service.beta.kubernetes.io/aws-load-balancer-healthcheck-healthy-threshold":     fmt.Sprintf("%d", *opts.Spec.LoadBalancer.HealthcheckHealthyThreshold),
		"service.beta.kubernetes.io/aws-load-balancer-healthcheck-unhealthy-threshold":   fmt.Sprintf("%d", *opts.Spec.LoadBalancer.HealthcheckUnhealthyThreshold),
		"service.beta.kubernetes.io/aws-load-balancer-healthcheck-interval":              fmt.Sprintf("%d", *opts.Spec.LoadBalancer.HealthcheckInterval),
		"service.beta.kubernetes.io/aws-load-balancer-healthcheck-timeout":               fmt.Sprintf("%d", *opts.Spec.LoadBalancer.HealthcheckTimeout),
	}

	if *opts.Spec.LoadBalancer.ProxyProtocol {
		annotations["service.beta.kubernetes.io/aws-load-balancer-proxy-protocol"] = "*"
	}

	return annotations

}

func httpProbe(path string, port intstr.IntOrString, scheme corev1.URIScheme, opts saasv1alpha1.HTTPProbeSpec) *corev1.Probe {
	if opts.IsDeactivated() {
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
		InitialDelaySeconds: *opts.InitialDelaySeconds,
		TimeoutSeconds:      *opts.TimeoutSeconds,
		PeriodSeconds:       *opts.PeriodSeconds,
		SuccessThreshold:    *opts.SuccessThreshold,
		FailureThreshold:    *opts.FailureThreshold,
	}
}

func (opts *Options) affinity() *corev1.Affinity {
	return &corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
				{
					Weight: 100,
					PodAffinityTerm: corev1.PodAffinityTerm{
						TopologyKey:   corev1.LabelHostname,
						LabelSelector: opts.selector(),
					},
				},
				{
					Weight: 99,
					PodAffinityTerm: corev1.PodAffinityTerm{
						TopologyKey:   corev1.LabelTopologyZone,
						LabelSelector: opts.selector(),
					},
				},
			},
		},
	}
}
