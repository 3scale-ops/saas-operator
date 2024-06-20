package service

import (
	"fmt"
	"strings"

	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceDescriptor struct {
	PortDefinitions []corev1.ServicePort
	saasv1alpha1.PublishingStrategy
}

func (sd *ServiceDescriptor) Service(prefix, suffix string) *corev1.Service {
	opts := ServiceOptions{}

	strategy := sd.PublishingStrategy.Strategy

	// generate Service resource
	if strategy == saasv1alpha1.SimpleStrategy || strategy == saasv1alpha1.Marin3rSidecarStrategy {

		var spec *saasv1alpha1.Simple
		if strategy == saasv1alpha1.SimpleStrategy {
			spec = sd.PublishingStrategy.Simple
		} else {
			spec = sd.PublishingStrategy.Marin3rSidecar.Simple
		}

		switch *spec.ServiceType {

		case saasv1alpha1.ServiceTypeClusterIP:
			opts.Type = corev1.ServiceTypeClusterIP
			opts.Name = fmt.Sprintf("%s-%s-%s", prefix, strings.ToLower(sd.EndpointName), suffix)
			opts.Annotations = map[string]string{}

		case saasv1alpha1.ServiceTypeELB:
			opts.Type = corev1.ServiceTypeLoadBalancer
			opts.Name = fmt.Sprintf("%s-%s-%s-elb", prefix, strings.ToLower(sd.EndpointName), suffix)
			spec.ElasticLoadBalancerConfig = saasv1alpha1.InitializeElasticLoadBalancerSpec(spec.ElasticLoadBalancerConfig, saasv1alpha1.DefaultElasticLoadBalancerSpec)
			opts.Annotations = ELBServiceAnnotations(*spec.ElasticLoadBalancerConfig, spec.ExternalDnsHostnames)

		case saasv1alpha1.ServiceTypeNLB:
			opts.Type = corev1.ServiceTypeLoadBalancer
			opts.Name = fmt.Sprintf("%s-%s-%s-nlb", prefix, strings.ToLower(sd.EndpointName), suffix)
			spec.NetworkLoadBalancerConfig = saasv1alpha1.InitializeNetworkLoadBalancerSpec(spec.NetworkLoadBalancerConfig, saasv1alpha1.DefaultNetworkLoadBalancerSpec)
			opts.Annotations = NLBServiceAnnotations(*spec.NetworkLoadBalancerConfig, spec.ExternalDnsHostnames)
		}

		// service name override
		if spec.ServiceNameOverride != nil {
			opts.Name = *spec.ServiceNameOverride
		}

		// Add service ports
		if spec.ServicePortsOverride != nil {
			opts.Ports = spec.ServicePortsOverride
		} else {
			opts.Ports = sd.PortDefinitions
		}
	}

	return opts.Service()
}

type ServiceOptions struct {
	Name        string
	Namespace   string
	Annotations map[string]string
	Type        corev1.ServiceType
	Ports       []corev1.ServicePort
}

func (opts ServiceOptions) Service() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        opts.Name,
			Namespace:   opts.Namespace,
			Annotations: opts.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Type:  opts.Type,
			Ports: opts.Ports,
		},
	}
}
