package service

import (
	"fmt"
	"strings"

	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceDescriptor struct {
	PortDef corev1.ServicePort
	saasv1alpha1.PublishingStrategy
	namePrefix string
}

func (sd *ServiceDescriptor) SetNamePrefix(prefix string) {
	sd.namePrefix = prefix
}

func (sd *ServiceDescriptor) Service() *corev1.Service {
	opts := ServiceOptions{}

	strategy := sd.PublishingStrategy.Strategy
	if strategy == saasv1alpha1.SimpleStrategy || strategy == saasv1alpha1.Marin3rStrategy {
		var simple *saasv1alpha1.Simple
		if strategy == saasv1alpha1.SimpleStrategy {
			simple = sd.PublishingStrategy.Simple
		} else {
			simple = sd.PublishingStrategy.Marin3rSidecar.Simple
		}

		// service name
		if simple.ServiceNameOverride != nil {
			opts.Name = *simple.ServiceNameOverride
		} else {
			opts.Name = fmt.Sprintf("%s-%s", sd.namePrefix, strings.ToLower(sd.EndpointName))
		}
		// service annotations
		switch *simple.ServiceType {
		case saasv1alpha1.ServiceTypeNLB:
			simple.NetworkLoadBalancerConfig = saasv1alpha1.InitializeNetworkLoadBalancerSpec(simple.NetworkLoadBalancerConfig, saasv1alpha1.DefaultNetworkLoadBalancerSpec)
			opts.Annotations = NLBServiceAnnotations(*simple.NetworkLoadBalancerConfig, simple.ExternalDnsHostnames)
		case saasv1alpha1.ServiceTypeELB:
			simple.ElasticLoadBalancerConfig = saasv1alpha1.InitializeElasticLoadBalancerSpec(simple.ElasticLoadBalancerConfig, saasv1alpha1.DefaultElasticLoadBalancerSpec)
			opts.Annotations = ELBServiceAnnotations(*simple.ElasticLoadBalancerConfig, simple.ExternalDnsHostnames)
		default:
			opts.Annotations = map[string]string{}
		}
		// service type
		switch *simple.ServiceType {
		case saasv1alpha1.ServiceTypeNLB:
			opts.Type = corev1.ServiceTypeLoadBalancer
		case saasv1alpha1.ServiceTypeELB:
			opts.Type = corev1.ServiceTypeLoadBalancer
		default:
			opts.Type = corev1.ServiceTypeClusterIP
		}
		// service ports
		opts.Ports = []corev1.ServicePort{sd.PortDef}
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
