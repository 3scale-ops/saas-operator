package autossl

import (
	"fmt"

	"github.com/3scale-ops/basereconciler/util"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (gen *Generator) deployment() *appsv1.Deployment {
	return &appsv1.Deployment{
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
					Volumes: []corev1.Volume{
						{
							Name: "autossl-cache",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "nginx-cache",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
					ImagePullSecrets: func() []corev1.LocalObjectReference {
						if gen.Spec.Image.PullSecretName != nil {
							return []corev1.LocalObjectReference{{Name: *gen.Spec.Image.PullSecretName}}
						}
						return nil
					}(),
					Containers: []corev1.Container{
						{
							Name:  "autossl",
							Image: fmt.Sprintf("%s:%s", *gen.Spec.Image.Name, *gen.Spec.Image.Tag),
							Ports: pod.ContainerPorts(
								pod.ContainerPortTCP("http", 8081),
								pod.ContainerPortTCP("https", 8444),
								pod.ContainerPortTCP("http-no-pp", 8080),
								pod.ContainerPortTCP("https-no-pp", 8443),
								pod.ContainerPortTCP("metrics", 9145),
							),
							Env:             gen.Options.BuildEnvironment(),
							Resources:       corev1.ResourceRequirements(*gen.Spec.Resources),
							ImagePullPolicy: *gen.Spec.Image.PullPolicy,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "autossl-cache",
									MountPath: "/etc/resty-auto-ssl/",
								},
								{
									Name:      "nginx-cache",
									MountPath: "/var/lib/nginx",
								},
							},
							LivenessProbe:  pod.HTTPProbe("/healthz", intstr.FromInt(9145), corev1.URISchemeHTTP, *gen.Spec.LivenessProbe),
							ReadinessProbe: pod.HTTPProbe("/healthz", intstr.FromInt(9145), corev1.URISchemeHTTP, *gen.Spec.ReadinessProbe),
						},
					},
					Affinity:                      pod.Affinity(gen.GetSelector(), gen.Spec.NodeAffinity),
					Tolerations:                   gen.Spec.Tolerations,
					TerminationGracePeriodSeconds: util.Pointer[int64](30),
				},
			},
		},
	}
}
