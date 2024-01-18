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

func (gen *AppGenerator) deployment() *appsv1.Deployment {
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
					InitContainers: []corev1.Container{
						{
							Name:  fmt.Sprintf("%s-k8s-deploy", gen.GetComponent()),
							Image: fmt.Sprintf("%s:%s", *gen.Image.Name, *gen.Image.Tag),
							Args: []string{
								"bundle", "exec", "rake", "k8s:deploy",
							},
							Env:             gen.Options.BuildEnvironment(),
							ImagePullPolicy: *gen.Image.PullPolicy,
						},
					},
					Containers: []corev1.Container{
						{
							Name:  strings.Join([]string{component, app}, "-"),
							Image: fmt.Sprintf("%s:%s", *gen.Image.Name, *gen.Image.Tag),
							Args: []string{
								"env",
								"PORT=3000",
								"container-entrypoint",
								"bundle",
								"exec",
								"unicorn",
								"-c",
								"config/unicorn.rb",
							},
							Env: gen.Options.BuildEnvironment(),
							Ports: pod.ContainerPorts(
								pod.ContainerPortTCP("ui-api", 3000),
								pod.ContainerPortTCP("metrics", 9394),
							),
							Resources:     corev1.ResourceRequirements(*gen.Spec.Resources),
							LivenessProbe: pod.TCPProbe(intstr.FromString("ui-api"), *gen.Spec.LivenessProbe),
							ReadinessProbe: pod.HTTPProbeWithHeaders("/check.txt", intstr.FromString("ui-api"),
								corev1.URISchemeHTTP, *gen.Spec.ReadinessProbe, map[string]string{"X-Forwarded-Proto": "https"}),
							ImagePullPolicy: *gen.Image.PullPolicy,
						},
					},
					Affinity:                      pod.Affinity(gen.GetSelector(), gen.Spec.NodeAffinity),
					Tolerations:                   gen.Spec.Tolerations,
					TerminationGracePeriodSeconds: gen.Spec.TerminationGracePeriodSeconds,
				},
			},
		},
	}

	dep.Spec.Template.Spec.Volumes = append(
		dep.Spec.Template.Spec.Volumes,
		corev1.Volume{
			Name: "system-config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: util.Pointer[int32](420),
					SecretName:  gen.ConfigFilesSecret,
				},
			},
		},
	)

	dep.Spec.Template.Spec.Containers[0].VolumeMounts = append(
		dep.Spec.Template.Spec.Containers[0].VolumeMounts,
		corev1.VolumeMount{
			Name:      "system-config",
			ReadOnly:  true,
			MountPath: "/opt/system-extra-configs",
		},
	)

	if gen.TwemproxySpec != nil {
		dep.Spec.Template = twemproxy.AddTwemproxySidecar(dep.Spec.Template, gen.TwemproxySpec)
	}
	return dep
}
