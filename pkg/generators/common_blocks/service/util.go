package service

import (
	"fmt"
	"strings"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// ELBServiceAnnotations returns annotations for services exposed through AWS Classic LoadBalancers
func ELBServiceAnnotations(cfg saasv1alpha1.LoadBalancerSpec, hostnames []string) map[string]string {
	annotations := map[string]string{
		"external-dns.alpha.kubernetes.io/hostname":                                      strings.Join(hostnames, ","),
		"service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled": fmt.Sprintf("%t", *cfg.CrossZoneLoadBalancingEnabled),
		"service.beta.kubernetes.io/aws-load-balancer-connection-draining-enabled":       fmt.Sprintf("%t", *cfg.ConnectionDrainingEnabled),
		"service.beta.kubernetes.io/aws-load-balancer-connection-draining-timeout":       fmt.Sprintf("%d", *cfg.ConnectionDrainingTimeout),
		"service.beta.kubernetes.io/aws-load-balancer-healthcheck-healthy-threshold":     fmt.Sprintf("%d", *cfg.HealthcheckHealthyThreshold),
		"service.beta.kubernetes.io/aws-load-balancer-healthcheck-unhealthy-threshold":   fmt.Sprintf("%d", *cfg.HealthcheckUnhealthyThreshold),
		"service.beta.kubernetes.io/aws-load-balancer-healthcheck-interval":              fmt.Sprintf("%d", *cfg.HealthcheckInterval),
		"service.beta.kubernetes.io/aws-load-balancer-healthcheck-timeout":               fmt.Sprintf("%d", *cfg.HealthcheckTimeout),
	}

	if *cfg.ProxyProtocol {
		annotations["service.beta.kubernetes.io/aws-load-balancer-proxy-protocol"] = "*"
	}

	return annotations

}

// NLBServiceAnnotations returns annotations for services exposed through AWS Network LoadBalancers
func NLBServiceAnnotations(cfg saasv1alpha1.NLBLoadBalancerSpec, hostnames []string) map[string]string {
	annotations := map[string]string{
		"external-dns.alpha.kubernetes.io/hostname":                                      strings.Join(hostnames, ","),
		"service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled": fmt.Sprintf("%t", *cfg.CrossZoneLoadBalancingEnabled),
	}

	if *cfg.ProxyProtocol {
		annotations["aws-nlb-helper.3scale.net/enable-targetgroups-proxy-protocol"] = "true"
	}

	if cfg.EIPAllocations != nil {
		annotations["service.beta.kubernetes.io/aws-load-balancer-eip-allocations"] = strings.Join(cfg.EIPAllocations, ",")
	}
	return annotations

}

// TCPPort returns a TCP corev1.ServicePort
func TCPPort(name string, port int32, targetPort intstr.IntOrString) corev1.ServicePort {
	return corev1.ServicePort{
		Name:       name,
		Port:       port,
		TargetPort: targetPort,
		Protocol:   corev1.ProtocolTCP,
	}
}

// Ports returns a list of corev1.ServicePort
func Ports(ports ...corev1.ServicePort) []corev1.ServicePort {
	list := []corev1.ServicePort{}
	return append(list, ports...)
}

// Excludes generates the list of excluded paths for a Service resource
func Excludes(fn basereconciler.GeneratorFunction) []string {
	svc := fn().(*corev1.Service)
	paths := []string{}
	paths = append(paths, "/spec/clusterIP")
	for idx := range svc.Spec.Ports {
		paths = append(paths, fmt.Sprintf("/spec/ports/%d/nodePort", idx))
	}
	return paths
}
