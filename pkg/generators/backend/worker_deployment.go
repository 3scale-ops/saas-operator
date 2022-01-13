package backend

import (
	"strings"

	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
	"github.com/3scale/saas-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Deployment returns a function that will return a Deployment
// resource when called
func (gen *WorkerGenerator) deployment() func() *appsv1.Deployment {

	return func() *appsv1.Deployment {

		dep := &appsv1.Deployment{
			Spec: appsv1.DeploymentSpec{
				Replicas: gen.WorkerSpec.Replicas,
				Strategy: appsv1.DeploymentStrategy{
					Type: appsv1.RollingUpdateDeploymentStrategyType,
					RollingUpdate: &appsv1.RollingUpdateDeployment{
						MaxUnavailable: util.IntStrPtr(intstr.FromInt(0)),
						MaxSurge:       util.IntStrPtr(intstr.FromInt(1)),
					},
				},
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						ImagePullSecrets: pod.ImagePullSecrets(gen.Image.PullSecretName),
						Containers: []corev1.Container{
							{
								Name:  strings.Join([]string{component, worker}, "-"),
								Image: pod.Image(gen.Image),
								Args:  []string{"bin/3scale_backend_worker", "run"},
								Ports: pod.ContainerPorts(
									pod.ContainerPortTCP("metrics", 9421),
								),
								Env:                      pod.BuildEnvironment(gen.Options),
								Resources:                corev1.ResourceRequirements(*gen.WorkerSpec.Resources),
								ImagePullPolicy:          *gen.Image.PullPolicy,
								LivenessProbe:            pod.HTTPProbe("/metrics", intstr.FromString("metrics"), corev1.URISchemeHTTP, *gen.WorkerSpec.LivenessProbe),
								ReadinessProbe:           pod.HTTPProbe("/metrics", intstr.FromString("metrics"), corev1.URISchemeHTTP, *gen.WorkerSpec.ReadinessProbe),
								TerminationMessagePath:   corev1.TerminationMessagePathDefault,
								TerminationMessagePolicy: corev1.TerminationMessageReadFile,
							},
						},
						Affinity:    pod.Affinity(gen.GetSelector(), gen.WorkerSpec.NodeAffinity),
						Tolerations: gen.WorkerSpec.Tolerations,
					},
				},
			},
		}
		return dep
	}
}
