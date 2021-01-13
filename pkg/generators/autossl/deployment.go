package autossl

import (
	"fmt"
	"strings"

	"github.com/3scale/saas-operator/pkg/basereconciler"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	leACMEStagingEndpoint = "https://acme-staging-v02.api.letsencrypt.org/directory"
)

// Deployment returns a basereconciler.GeneratorFunction funtion that will return a Deployment
// resource when called
func (opts *Options) Deployment() basereconciler.GeneratorFunction {

	return func() client.Object {

		return &appsv1.Deployment{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Deployment",
				APIVersion: appsv1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      Component,
				Namespace: opts.Namespace,
				Labels:    opts.labels(),
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: opts.Spec.Replicas,
				Selector: opts.selector(),
				Strategy: appsv1.DeploymentStrategy{
					Type: appsv1.RollingUpdateDeploymentStrategyType,
					RollingUpdate: &appsv1.RollingUpdateDeployment{
						MaxUnavailable: &intstr.IntOrString{
							Type:   intstr.Int,
							IntVal: 0,
						},
						MaxSurge: &intstr.IntOrString{
							Type:   intstr.Int,
							IntVal: 1,
						},
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: opts.labelsWithSelector(),
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
							if opts.Spec.Image.PullSecretName != nil {
								return []corev1.LocalObjectReference{{Name: *opts.Spec.Image.PullSecretName}}
							}
							return nil
						}(),
						Containers: []corev1.Container{
							{
								Name:  Component,
								Image: fmt.Sprintf("%s:%s", *opts.Spec.Image.Name, *opts.Spec.Image.Tag),
								Ports: []corev1.ContainerPort{
									{
										Name:          "http",
										ContainerPort: 8081,
										Protocol:      corev1.ProtocolTCP,
									},
									{
										Name:          "https",
										ContainerPort: 8444,
										Protocol:      corev1.ProtocolTCP,
									},
									{
										Name:          "http-no-pp",
										ContainerPort: 8080,
										Protocol:      corev1.ProtocolTCP,
									},
									{
										Name:          "https-no-pp",
										ContainerPort: 8443,
										Protocol:      corev1.ProtocolTCP,
									},
									{
										Name:          "metrics",
										ContainerPort: 9145,
										Protocol:      corev1.ProtocolTCP,
									},
								},
								Env: []corev1.EnvVar{
									{
										Name: "ACME_STAGING",
										Value: func() string {
											if *opts.Spec.Config.ACMEStaging {
												return leACMEStagingEndpoint
											}
											return ""
										}(),
									},
									{Name: "CONTACT_EMAIL", Value: opts.Spec.Config.ContactEmail},
									{Name: "PROXY_ENDPOINT", Value: opts.Spec.Config.ProxyEndpoint},
									{Name: "STORAGE_ADAPTER", Value: "redis"},
									{Name: "REDIS_HOST", Value: opts.Spec.Config.RedisHost},
									{Name: "REDIS_PORT", Value: fmt.Sprintf("%v", *opts.Spec.Config.RedisPort)},
									{Name: "VERIFICATION_ENDPOINT", Value: opts.Spec.Config.VerificationEndpoint},
									{Name: "LOG_LEVEL", Value: *opts.Spec.Config.LogLevel},
									{Name: "DOMAIN_WHITELIST", Value: strings.Join(opts.Spec.Config.DomainWhitelist, ",")},
									{Name: "DOMAIN_BLACKLIST", Value: strings.Join(opts.Spec.Config.DomainBlacklist, ",")},
								},
								Resources:              corev1.ResourceRequirements(*opts.Spec.Resources),
								TerminationMessagePath: corev1.TerminationMessagePathDefault,
								ImagePullPolicy:        corev1.PullAlways,
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
								LivenessProbe:  httpProbe("/health", intstr.FromInt(9145), corev1.URISchemeHTTP, *opts.Spec.LivenessProbe),
								ReadinessProbe: httpProbe("/health", intstr.FromInt(9145), corev1.URISchemeHTTP, *opts.Spec.ReadinessProbe),
							},
						},
						Affinity: opts.affinity(),
					},
				},
			},
		}
	}
}
