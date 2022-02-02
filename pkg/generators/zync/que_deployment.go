package zync

import (
	"fmt"

	"github.com/3scale/saas-operator/pkg/resource_builders/pod"

	"github.com/3scale/saas-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// deployment returns a function that will return a *appsv1.Deployment for zync-que
func (gen *QueGenerator) deployment() func() *appsv1.Deployment {

	return func() *appsv1.Deployment {

		dep := &appsv1.Deployment{
			Spec: appsv1.DeploymentSpec{
				Replicas: gen.QueSpec.Replicas,
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
							if gen.Image.PullSecretName != nil {
								return []corev1.LocalObjectReference{{Name: *gen.Image.PullSecretName}}
							}
							return nil
						}(),
						Containers: []corev1.Container{
							{
								Name:  "zync-que",
								Image: fmt.Sprintf("%s:%s", *gen.Image.Name, *gen.Image.Tag),
								Command: []string{
									"/usr/bin/bash",
									"-c",
									"bundle exec rake 'que[--worker-count 10]'",
								},
								Ports: pod.ContainerPorts(
									pod.ContainerPortTCP("metrics", 9394),
								),
								Env: func() []corev1.EnvVar {
									envVars := pod.BuildEnvironment(gen.Options)
									envVars = append(envVars,
										corev1.EnvVar{
											Name: "POD_NAME",
											ValueFrom: &v1.EnvVarSource{
												FieldRef: &v1.ObjectFieldSelector{
													FieldPath:  "metadata.name",
													APIVersion: "v1",
												},
											},
										},
										corev1.EnvVar{
											Name: "POD_NAMESPACE",
											ValueFrom: &v1.EnvVarSource{
												FieldRef: &v1.ObjectFieldSelector{
													FieldPath:  "metadata.namespace",
													APIVersion: "v1",
												},
											},
										},
									)
									return envVars
								}(),
								Resources:                corev1.ResourceRequirements(*gen.QueSpec.Resources),
								ImagePullPolicy:          *gen.Image.PullPolicy,
								LivenessProbe:            pod.HTTPProbe("/metrics", intstr.FromString("metrics"), corev1.URISchemeHTTP, *gen.QueSpec.LivenessProbe),
								ReadinessProbe:           pod.HTTPProbe("/metrics", intstr.FromString("metrics"), corev1.URISchemeHTTP, *gen.QueSpec.ReadinessProbe),
								TerminationMessagePath:   corev1.TerminationMessagePathDefault,
								TerminationMessagePolicy: corev1.TerminationMessageReadFile,
							},
						},
						Affinity:    pod.Affinity(gen.GetSelector(), gen.QueSpec.NodeAffinity),
						Tolerations: gen.QueSpec.Tolerations,
					},
				},
			},
		}
		return dep
	}
}
