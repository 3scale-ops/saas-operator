package config

import (
	"github.com/3scale-ops/basereconciler/util"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/service"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func DefaultPublishingStrategy() []service.ServiceDescriptor {
	return []service.ServiceDescriptor{
		{
			PublishingStrategy: saasv1alpha1.PublishingStrategy{
				Strategy:     saasv1alpha1.SimpleStrategy,
				EndpointName: "Gateway",
				Simple: &saasv1alpha1.Simple{
					ServiceType: util.Pointer(saasv1alpha1.ServiceTypeELB),
					ElasticLoadBalancerConfig: &saasv1alpha1.LoadBalancerSpec{
						ProxyProtocol:                 util.Pointer(true),
						CrossZoneLoadBalancingEnabled: util.Pointer(true),
						ConnectionDrainingEnabled:     util.Pointer(true),
						ConnectionDrainingTimeout:     util.Pointer[int32](60),
						HealthcheckHealthyThreshold:   util.Pointer[int32](2),
						HealthcheckUnhealthyThreshold: util.Pointer[int32](2),
						HealthcheckInterval:           util.Pointer[int32](5),
						HealthcheckTimeout:            util.Pointer[int32](3),
					},
				},
			},
			PortDef: corev1.ServicePort{
				Name:       "gateway-http",
				Protocol:   corev1.ProtocolTCP,
				Port:       80,
				TargetPort: intstr.FromString("gateway-http"),
			},
		},
		{
			PublishingStrategy: saasv1alpha1.PublishingStrategy{
				Strategy:     saasv1alpha1.SimpleStrategy,
				EndpointName: "Management",
				Simple: &saasv1alpha1.Simple{
					ServiceType: util.Pointer(saasv1alpha1.ServiceTypeClusterIP),
				},
			},
			PortDef: corev1.ServicePort{
				Name:       "management",
				Protocol:   corev1.ProtocolTCP,
				Port:       80,
				TargetPort: intstr.FromString("management"),
			},
		},
	}
}
