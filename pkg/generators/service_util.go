package generators

import (
	"fmt"
	"strings"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// ELBServiceAnnotations returns annotations for services exposed through AWS Classic LoadBalancers
func (bo *BaseOptions) ELBServiceAnnotations(cfg saasv1alpha1.LoadBalancerSpec, hostnames []string) map[string]string {
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

// ServicePortTCP returns a TCP corev1.ServicePort
func (bo *BaseOptions) ServicePortTCP(name string, port int32, targetPort intstr.IntOrString) corev1.ServicePort {
	return corev1.ServicePort{
		Name:       name,
		Port:       port,
		TargetPort: targetPort,
		Protocol:   corev1.ProtocolTCP,
	}
}

// ServicePorts returns a list of corev1.ServicePort
func (bo *BaseOptions) ServicePorts(ports ...corev1.ServicePort) []corev1.ServicePort {
	list := []corev1.ServicePort{}
	return append(list, ports...)
}

// ServiceExcludes generates the list of excluded paths for a Service resource
func (bo *BaseOptions) ServiceExcludes(fn basereconciler.GeneratorFunction) []string {
	svc := fn().(*corev1.Service)
	paths := []string{}
	paths = append(paths, "/spec/clusterIP")
	for idx := range svc.Spec.Ports {
		paths = append(paths, fmt.Sprintf("/spec/ports/%d/nodePort", idx))
	}
	return paths
}
