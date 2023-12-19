package system

import (
	"fmt"
	"strings"

	"github.com/3scale-ops/basereconciler/util"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/twemproxy"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (gen *SidekiqGenerator) deployment() *appsv1.Deployment {
	dep := &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Replicas: gen.Spec.Replicas,
			Strategy: appsv1.DeploymentStrategy(*gen.Spec.DeploymentStrategy),
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
							Name:  strings.Join([]string{component, sidekiq}, "-"),
							Image: fmt.Sprintf("%s:%s", *gen.Image.Name, *gen.Image.Tag),
							Args: func(queues []string) []string {
								var args = []string{"sidekiq"}
								for _, queue := range queues {
									args = append(args, "--queue", queue)
								}
								return args
							}(gen.Spec.Config.Queues),
							Env: func() []corev1.EnvVar {
								envVars := pod.BuildEnvironment(gen.Options)
								envVars = append(envVars,
									corev1.EnvVar{
										Name:  "RAILS_MAX_THREADS",
										Value: fmt.Sprintf("%d", *gen.Spec.Config.MaxThreads),
									},
								)
								return envVars
							}(),
							Ports: pod.ContainerPorts(
								pod.ContainerPortTCP("metrics", 9394),
							),
							Resources:       corev1.ResourceRequirements(*gen.Spec.Resources),
							LivenessProbe:   pod.HTTPProbe("/metrics", intstr.FromString("metrics"), corev1.URISchemeHTTP, *gen.Spec.LivenessProbe),
							ReadinessProbe:  pod.HTTPProbe("/metrics", intstr.FromString("metrics"), corev1.URISchemeHTTP, *gen.Spec.LivenessProbe),
							ImagePullPolicy: *gen.Image.PullPolicy,
							VolumeMounts: []corev1.VolumeMount{{
								Name:      "system-tmp",
								MountPath: "/tmp",
							}},
						},
					},
					Volumes: []corev1.Volume{{
						Name: "system-tmp",
						VolumeSource: corev1.VolumeSource{
							EmptyDir: &corev1.EmptyDirVolumeSource{
								Medium: corev1.StorageMediumMemory,
							},
						},
					}},
					Affinity:                      pod.Affinity(gen.GetSelector(), gen.Spec.NodeAffinity),
					Tolerations:                   gen.Spec.Tolerations,
					TerminationGracePeriodSeconds: gen.Spec.TerminationGracePeriodSeconds,
				},
			},
		},
	}

	dep.Spec.Template.Spec.Volumes = append(dep.Spec.Template.Spec.Volumes,
		corev1.Volume{
			Name: "system-config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: util.Pointer[int32](420),
					SecretName:  gen.ConfigFilesSecret,
				},
			},
		})
	dep.Spec.Template.Spec.Containers[0].VolumeMounts = append(dep.Spec.Template.Spec.Containers[0].VolumeMounts,
		corev1.VolumeMount{
			Name:      "system-config",
			ReadOnly:  true,
			MountPath: "/opt/system-extra-configs",
		})

	if gen.TwemproxySpec != nil {
		dep.Spec.Template = twemproxy.AddTwemproxySidecar(dep.Spec.Template, gen.TwemproxySpec)
	}

	return dep
}
