package system

import (
	"fmt"

	"github.com/3scale/saas-operator/pkg/basereconciler"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/marin3r"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pod"
	"github.com/3scale/saas-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Deployment returns a basereconciler.GeneratorFunction function that will return a Deployment
// resource when called
func (gen *AppGenerator) Deployment() basereconciler.GeneratorFunction {

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
								Env: pod.BuildEnvironment(gen.Options),
								Ports: pod.ContainerPorts(
									pod.ContainerPortTCP("ui-api", 3000),
									pod.ContainerPortTCP("metrics", 9394),
								),
								Resources:     corev1.ResourceRequirements(*gen.Spec.Resources),
								LivenessProbe: pod.TCPProbe(intstr.FromString("ui-api"), *gen.Spec.LivenessProbe),
								ReadinessProbe: pod.HTTPProbeWithHeaders("/check.txt", intstr.FromString("ui-api"),
									corev1.URISchemeHTTP, *gen.Spec.ReadinessProbe, map[string]string{"X-Forwarded-Proto": "https"}),
								ImagePullPolicy:          *gen.ImageSpec.PullPolicy,
								TerminationMessagePath:   corev1.TerminationMessagePathDefault,
								TerminationMessagePolicy: corev1.TerminationMessageReadFile,
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "system-config",
										ReadOnly:  true,
										MountPath: "/opt/system-extra-configs",
									},
								},
							},
						},
						Affinity: pod.Affinity(gen.Selector().MatchLabels),
						Volumes: []corev1.Volume{
							{
								Name: "system-config",
								VolumeSource: corev1.VolumeSource{
									Secret: &corev1.SecretVolumeSource{
										SecretName: systemConfigSecret,
									},
								},
							},
						},
					},
				},
			},
		}

		if !gen.Spec.Marin3r.IsDeactivated() {
			dep = marin3r.EnableSidecar(*dep, *gen.Spec.Marin3r)
		}

		return dep
	}
}
