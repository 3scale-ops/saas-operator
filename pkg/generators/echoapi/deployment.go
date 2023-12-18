package echoapi

import (
	"fmt"

	"github.com/3scale-ops/basereconciler/util"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/marin3r"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// deployment returns a function that will return a *appsv1.Deployment for echo-api
func (gen *Generator) deployment() *appsv1.Deployment {

	dep := &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Replicas: gen.Spec.Replicas,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: util.Pointer(intstr.FromInt(0)),
					MaxSurge:       util.Pointer(intstr.FromInt(1)),
				},
			},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					ImagePullSecrets: func() []corev1.LocalObjectReference {
						if gen.Spec.Image.PullSecretName != nil {
							return []corev1.LocalObjectReference{{Name: *gen.Spec.Image.PullSecretName}}
						}
						return nil
					}(),
					Containers: []corev1.Container{
						{
							Name:  "echo-api",
							Image: fmt.Sprintf("%s:%s", *gen.Spec.Image.Name, *gen.Spec.Image.Tag),
							Ports: pod.ContainerPorts(
								pod.ContainerPortTCP("http", 9292),
							),
							Resources:       corev1.ResourceRequirements(*gen.Spec.Resources),
							ImagePullPolicy: *gen.Spec.Image.PullPolicy,
							LivenessProbe:   pod.TCPProbe(intstr.FromInt(9292), *gen.Spec.LivenessProbe),
							ReadinessProbe:  pod.HTTPProbe("/", intstr.FromInt(9292), corev1.URISchemeHTTP, *gen.Spec.ReadinessProbe),
						},
					},
					Affinity:                      pod.Affinity(gen.GetSelector(), gen.Spec.NodeAffinity),
					Tolerations:                   gen.Spec.Tolerations,
					TerminationGracePeriodSeconds: util.Pointer[int64](30),
				},
			},
		},
	}

	if !gen.Spec.Marin3r.IsDeactivated() {
		dep = marin3r.EnableSidecar(*dep, *gen.Spec.Marin3r)
	}

	return dep
}
