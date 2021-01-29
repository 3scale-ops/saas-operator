package autossl

import (
	"fmt"
	"strings"

	"github.com/3scale/saas-operator/pkg/basereconciler"
	"github.com/3scale/saas-operator/pkg/generators/autossl/config"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pod"
	"github.com/3scale/saas-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	leACMEStagingEndpoint = "https://acme-staging-v02.api.letsencrypt.org/directory"
)

// Deployment returns a basereconciler.GeneratorFunction function that will return a Deployment
// resource when called
func (gen *Generator) Deployment() basereconciler.GeneratorFunction {

	return func() client.Object {

		return &appsv1.Deployment{
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
								Name:  gen.GetComponent(),
								Image: fmt.Sprintf("%s:%s", *gen.Spec.Image.Name, *gen.Spec.Image.Tag),
								Ports: pod.ContainerPorts(
									pod.ContainerPortTCP("http", 8081),
									pod.ContainerPortTCP("https", 8444),
									pod.ContainerPortTCP("http-no-pp", 8080),
									pod.ContainerPortTCP("https-no-pp", 8443),
									pod.ContainerPortTCP("metrics", 9145),
								),
								Env: pod.GenerateEnvironment(config.Default,
									func() map[string]pod.EnvVarValue {
										m := map[string]pod.EnvVarValue{
											config.ContactEmail:         &pod.DirectValue{Value: gen.Spec.Config.ContactEmail},
											config.ProxyEndpoint:        &pod.DirectValue{Value: gen.Spec.Config.ProxyEndpoint},
											config.RedisHost:            &pod.DirectValue{Value: gen.Spec.Config.RedisHost},
											config.RedisPort:            &pod.DirectValue{Value: fmt.Sprintf("%v", *gen.Spec.Config.RedisPort)},
											config.VerificationEndpoint: &pod.DirectValue{Value: gen.Spec.Config.VerificationEndpoint},
											config.LogLevel:             &pod.DirectValue{Value: *gen.Spec.Config.LogLevel},
											config.DomainWhitelist:      &pod.DirectValue{Value: strings.Join(gen.Spec.Config.DomainWhitelist, ",")},
											config.DomainBlacklist:      &pod.DirectValue{Value: strings.Join(gen.Spec.Config.DomainBlacklist, ",")},
										}
										if *gen.Spec.Config.ACMEStaging {
											m[config.ACMEStaging] = &pod.DirectValue{Value: leACMEStagingEndpoint}
										}
										return m
									}(),
								),
								Resources:              corev1.ResourceRequirements(*gen.Spec.Resources),
								TerminationMessagePath: corev1.TerminationMessagePathDefault,
								ImagePullPolicy:        *gen.Spec.Image.PullPolicy,
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
						Affinity: pod.Affinity(gen.Selector().MatchLabels),
					},
				},
			},
		}
	}
}
