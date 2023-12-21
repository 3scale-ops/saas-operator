package backend

import (
	"strings"

	"github.com/3scale-ops/basereconciler/util"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (gen *CronGenerator) deployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Replicas: gen.CronSpec.Replicas,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: util.Pointer(intstr.FromInt(0)),
					MaxSurge:       util.Pointer(intstr.FromInt(1)),
				},
			},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					ImagePullSecrets: pod.ImagePullSecrets(gen.Image.PullSecretName),
					Containers: []corev1.Container{
						{
							Name:            strings.Join([]string{component, cron}, "-"),
							Image:           pod.Image(gen.Image),
							Args:            []string{"backend-cron"},
							Env:             gen.Options.BuildEnvironment(),
							Resources:       corev1.ResourceRequirements(*gen.CronSpec.Resources),
							ImagePullPolicy: *gen.Image.PullPolicy,
						},
					},
					Affinity:                      pod.Affinity(gen.GetSelector(), gen.CronSpec.NodeAffinity),
					Tolerations:                   gen.CronSpec.Tolerations,
					TerminationGracePeriodSeconds: util.Pointer[int64](30),
				},
			},
		},
	}
}
