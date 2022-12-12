package system

import (
	"fmt"
	"strings"

	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
	"github.com/3scale/saas-operator/pkg/resource_builders/twemproxy"
	"github.com/3scale/saas-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

// Deployment returns a basereconciler.GeneratorFunction function that will return a Deployment
// resource when called
func (gen *SidekiqGenerator) deployment() func() *appsv1.Deployment {

	return func() *appsv1.Deployment {

		dep := &appsv1.Deployment{
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
								Resources:                corev1.ResourceRequirements(*gen.Spec.Resources),
								LivenessProbe:            pod.HTTPProbe("/metrics", intstr.FromString("metrics"), corev1.URISchemeHTTP, *gen.Spec.LivenessProbe),
								ReadinessProbe:           pod.HTTPProbe("/metrics", intstr.FromString("metrics"), corev1.URISchemeHTTP, *gen.Spec.LivenessProbe),
								ImagePullPolicy:          *gen.Image.PullPolicy,
								TerminationMessagePath:   corev1.TerminationMessagePathDefault,
								TerminationMessagePolicy: corev1.TerminationMessageReadFile,
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
						RestartPolicy:                 corev1.RestartPolicyAlways,
						SecurityContext:               &corev1.PodSecurityContext{},
						Affinity:                      pod.Affinity(gen.GetSelector(), gen.Spec.NodeAffinity),
						Tolerations:                   gen.Spec.Tolerations,
						TerminationGracePeriodSeconds: pointer.Int64(30),
					},
				},
			},
		}

		dep.Spec.Template.Spec.Volumes = append(dep.Spec.Template.Spec.Volumes,
			corev1.Volume{
				Name: "system-config",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						DefaultMode: pointer.Int32Ptr(420),
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
			dep = twemproxy.AddTwemproxySidecar(*dep, gen.TwemproxySpec)
		}

		return dep
	}
}
