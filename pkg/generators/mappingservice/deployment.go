package mappingservice

import (
	"fmt"

	"github.com/3scale/saas-operator/pkg/basereconciler"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pod"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/secrets"
	"github.com/3scale/saas-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
									pod.ContainerPortTCP("mapping", 8093),
									pod.ContainerPortTCP("management", 8090),
									pod.ContainerPortTCP("metrics", 9421),
								),
								Env: func() []corev1.EnvVar {
									vars := []corev1.EnvVar{
										{Name: "APICAST_CONFIGURATION_LOADER", Value: "lazy"},
										{Name: "API_HOST", Value: gen.Spec.Config.APIHost},
										{Name: "APICAST_LOG_LEVEL", Value: *gen.Spec.Config.LogLevel},
										secrets.EnvVarFromSecret("MASTER_ACCESS_TOKEN",
											masterTokenSecretName, gen.Spec.Config.SystemAdminToken.FromVault.Key),
									}
									if gen.Spec.Config.PreviewBaseDomain != nil {
										vars = append(vars, corev1.EnvVar{Name: "PREVIEW_BASE_DOMAIN", Value: *gen.Spec.Config.PreviewBaseDomain})
									}
									return vars
								}(),
								Resources:              corev1.ResourceRequirements(*gen.Spec.Resources),
								TerminationMessagePath: corev1.TerminationMessagePathDefault,
								ImagePullPolicy:        *gen.Spec.Image.PullPolicy,
								LivenessProbe:          pod.TCPProbe(intstr.FromString("mapping"), *gen.Spec.LivenessProbe),
								ReadinessProbe:         pod.HTTPProbe("/status/ready", intstr.FromString("management"), corev1.URISchemeHTTP, *gen.Spec.ReadinessProbe),
							},
						},
						Affinity: pod.Affinity(gen.Selector().MatchLabels),
					},
				},
			},
		}
	}
}
