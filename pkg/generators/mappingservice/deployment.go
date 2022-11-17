package mappingservice

import (
	"fmt"

	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
	"github.com/3scale/saas-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// deployment returns a function that will return a *appsv1.Deployment for echo-api
func (gen *Generator) deployment() func() *appsv1.Deployment {

	return func() *appsv1.Deployment {

		return &appsv1.Deployment{
			Spec: appsv1.DeploymentSpec{
				Replicas: gen.Spec.Replicas,
				Strategy: appsv1.DeploymentStrategy{
					Type: appsv1.RollingUpdateDeploymentStrategyType,
					RollingUpdate: &appsv1.RollingUpdateDeployment{
						MaxUnavailable: util.IntStrPtr(intstr.FromInt(0)),
						MaxSurge:       util.IntStrPtr(intstr.FromInt(1)),
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
								Name:  "mapping-service",
								Image: fmt.Sprintf("%s:%s", *gen.Spec.Image.Name, *gen.Spec.Image.Tag),
								Ports: pod.ContainerPorts(
									pod.ContainerPortTCP("mapping", 8093),
									pod.ContainerPortTCP("management", 8090),
									pod.ContainerPortTCP("metrics", 9421),
								),
								Env:                      pod.BuildEnvironment(gen.Options),
								Resources:                corev1.ResourceRequirements(*gen.Spec.Resources),
								ImagePullPolicy:          *gen.Spec.Image.PullPolicy,
								LivenessProbe:            pod.TCPProbe(intstr.FromString("mapping"), *gen.Spec.LivenessProbe),
								ReadinessProbe:           pod.HTTPProbe("/status/ready", intstr.FromString("management"), corev1.URISchemeHTTP, *gen.Spec.ReadinessProbe),
								TerminationMessagePath:   corev1.TerminationMessagePathDefault,
								TerminationMessagePolicy: corev1.TerminationMessageReadFile,
							},
						},
						RestartPolicy:   corev1.RestartPolicyAlways,
						SecurityContext: &corev1.PodSecurityContext{},
						Affinity:        pod.Affinity(gen.GetSelector(), gen.Spec.NodeAffinity),
						Tolerations:     gen.Spec.Tolerations,
					},
				},
			},
		}
	}
}
