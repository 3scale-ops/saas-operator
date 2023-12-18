package apicast

import (
	"fmt"

	"github.com/3scale-ops/basereconciler/util"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/marin3r"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Deployment returns a function that will return a Deployment
// resource when called
func (gen *EnvGenerator) deployment() *appsv1.Deployment {

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
							Name:  "apicast",
							Image: fmt.Sprintf("%s:%s", *gen.Spec.Image.Name, *gen.Spec.Image.Tag),
							Ports: pod.ContainerPorts(
								pod.ContainerPortTCP("gateway", 8080),
								pod.ContainerPortTCP("management", 8090),
								pod.ContainerPortTCP("metrics", 9421),
							),
							Env:             pod.BuildEnvironment(gen.Options),
							Resources:       corev1.ResourceRequirements(*gen.Spec.Resources),
							LivenessProbe:   pod.TCPProbe(intstr.FromString("gateway"), *gen.Spec.LivenessProbe),
							ReadinessProbe:  pod.HTTPProbe("/status/ready", intstr.FromString("management"), corev1.URISchemeHTTP, *gen.Spec.ReadinessProbe),
							ImagePullPolicy: *gen.Spec.Image.PullPolicy,
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
