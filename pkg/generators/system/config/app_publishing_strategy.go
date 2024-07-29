package config

import (
	"github.com/3scale-ops/basereconciler/util"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/service"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func DefaultAppPublishingStrategy() []service.ServiceDescriptor {
	return []service.ServiceDescriptor{
		{
			PublishingStrategy: saasv1alpha1.PublishingStrategy{
				Strategy:     saasv1alpha1.SimpleStrategy,
				EndpointName: "HTTP",
				Simple:       &saasv1alpha1.Simple{ServiceType: util.Pointer(saasv1alpha1.ServiceTypeClusterIP)},
			},
			PortDefinitions: []corev1.ServicePort{{
				Name:       "http",
				Protocol:   corev1.ProtocolTCP,
				Port:       3000,
				TargetPort: intstr.FromString("ui-api"),
			}},
		},
	}
}
