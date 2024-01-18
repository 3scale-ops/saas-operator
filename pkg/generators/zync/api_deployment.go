package zync

import (
	"fmt"

	"github.com/3scale-ops/basereconciler/util"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (gen *APIGenerator) deployment() *appsv1.Deployment {
	dep := &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Replicas: gen.APISpec.Replicas,
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
						if gen.Image.PullSecretName != nil {
							return []corev1.LocalObjectReference{{Name: *gen.Image.PullSecretName}}
						}
						return nil
					}(),
					Containers: []corev1.Container{
						{
							Name:  "zync",
							Image: fmt.Sprintf("%s:%s", *gen.Image.Name, *gen.Image.Tag),
							Ports: pod.ContainerPorts(
								pod.ContainerPortTCP("http", 8080),
								pod.ContainerPortTCP("metrics", 9393),
							),
							Env: gen.Options.WithExtraEnv(
								[]corev1.EnvVar{
									{
										Name: "POD_NAME",
										ValueFrom: &corev1.EnvVarSource{
											FieldRef: &corev1.ObjectFieldSelector{
												FieldPath:  "metadata.name",
												APIVersion: "v1",
											},
										},
									},
									{
										Name: "POD_NAMESPACE",
										ValueFrom: &corev1.EnvVarSource{
											FieldRef: &corev1.ObjectFieldSelector{
												FieldPath:  "metadata.namespace",
												APIVersion: "v1",
											},
										},
									},
								},
							).BuildEnvironment(),
							Resources:       corev1.ResourceRequirements(*gen.APISpec.Resources),
							ImagePullPolicy: *gen.Image.PullPolicy,
							LivenessProbe:   pod.HTTPProbe("/status/live", intstr.FromString("http"), corev1.URISchemeHTTP, *gen.APISpec.LivenessProbe),
							ReadinessProbe:  pod.HTTPProbe("/status/ready", intstr.FromString("http"), corev1.URISchemeHTTP, *gen.APISpec.ReadinessProbe),
						},
					},
					Affinity:                      pod.Affinity(gen.GetSelector(), gen.APISpec.NodeAffinity),
					Tolerations:                   gen.APISpec.Tolerations,
					TerminationGracePeriodSeconds: util.Pointer[int64](30),
				},
			},
		},
	}
	return dep
}
