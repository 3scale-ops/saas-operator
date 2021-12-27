package system

import (
	"fmt"

	basereconciler_types "github.com/3scale/saas-operator/pkg/basereconciler/types"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/marin3r"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pod"
	"github.com/3scale/saas-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Deployment returns a basereconciler_types.GeneratorFunction function that will return a Deployment
// resource when called
func (gen *SidekiqGenerator) Deployment() basereconciler_types.GeneratorFunction {

	return func() client.Object {

		dep := &appsv1.Deployment{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Deployment",
				APIVersion: appsv1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      gen.GetComponent(),
				Namespace: gen.Namespace,
				Labels:    gen.GetLabels(),
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: gen.Spec.Replicas,
				Selector: gen.Selector(),
				Strategy: appsv1.DeploymentStrategy{
					Type: appsv1.RollingUpdateDeploymentStrategyType,
					RollingUpdate: &appsv1.RollingUpdateDeployment{
						MaxUnavailable: util.IntStrPtr(intstr.FromInt(0)),
						MaxSurge:       util.IntStrPtr(intstr.FromInt(1)),
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: gen.LabelsWithSelector(),
					},
					Spec: corev1.PodSpec{
						ImagePullSecrets: func() []corev1.LocalObjectReference {
							if gen.ImageSpec.PullSecretName != nil {
								return []corev1.LocalObjectReference{{Name: *gen.ImageSpec.PullSecretName}}
							}
							return nil
						}(),
						Containers: []corev1.Container{
							{
								Name:  gen.GetComponent(),
								Image: fmt.Sprintf("%s:%s", *gen.ImageSpec.Name, *gen.ImageSpec.Tag),
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
								ImagePullPolicy:          *gen.ImageSpec.PullPolicy,
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
						Affinity:    pod.Affinity(gen.Selector().MatchLabels, gen.Spec.NodeAffinity),
						Tolerations: gen.Spec.Tolerations,
					},
				},
			},
		}

		dep.Spec.Template.Spec.Volumes = append(dep.Spec.Template.Spec.Volumes,
			corev1.Volume{
				Name: "system-config",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: gen.ConfigFilesSecret,
					},
				},
			})
		dep.Spec.Template.Spec.Containers[0].VolumeMounts = append(dep.Spec.Template.Spec.Containers[0].VolumeMounts,
			corev1.VolumeMount{
				Name:      "system-config",
				ReadOnly:  true,
				MountPath: "/opt/system-extra-configs",
			})

		if !gen.Spec.Marin3r.IsDeactivated() {
			dep = marin3r.EnableSidecar(*dep, *gen.Spec.Marin3r)
		}

		return dep
	}
}
