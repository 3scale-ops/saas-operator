package corsproxy

import (
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pod"
	"github.com/3scale/saas-operator/pkg/generators/corsproxy/config"
	"github.com/3scale/saas-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Deployment returns a basereconciler.GeneratorFunction function that will return a Deployment
// resource when called
func (gen *Generator) Deployment(hash string) basereconciler.GeneratorFunction {

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
						Annotations: map[string]string{
							saasv1alpha1.RolloutTriggerAnnotationKeyPrefix + "config.systemDatabaseDSN.hash": hash,
						},
					},
					Spec: corev1.PodSpec{
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
									pod.ContainerPortTCP("http", 8080),
									pod.ContainerPortTCP("metrics", 9145),
								),
								Env: pod.GenerateEnvironment(config.Default,
									map[string]pod.EnvVarValue{
										config.DatabaseURL: &pod.SecretRef{SecretName: config.SecretDefinitions.LookupSecretName(config.DatabaseURL)},
									}),
								Resources:              corev1.ResourceRequirements(*gen.Spec.Resources),
								TerminationMessagePath: corev1.TerminationMessagePathDefault,
								ImagePullPolicy:        *gen.Spec.Image.PullPolicy,
								LivenessProbe:          pod.HTTPProbe("/healthz", intstr.FromString("metrics"), corev1.URISchemeHTTP, *gen.Spec.LivenessProbe),
								ReadinessProbe:         pod.HTTPProbe("/healthz", intstr.FromString("metrics"), corev1.URISchemeHTTP, *gen.Spec.ReadinessProbe),
							},
						},
						Affinity: pod.Affinity(gen.Selector().MatchLabels),
					},
				},
			},
		}
	}
}
