package system

import (
	"fmt"
	"strings"

	"github.com/3scale-ops/basereconciler/util"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pod"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/twemproxy"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StatefulSet returns a basereconciler.GeneratorFunction function that will return
// a StatefulSet resource when called
func (gen *ConsoleGenerator) statefulset() *appsv1.StatefulSet {
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gen.GetComponent(),
			Namespace: gen.Namespace,
			Labels:    gen.GetLabels(),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: util.Pointer[int32](1),
			Selector: &metav1.LabelSelector{MatchLabels: gen.GetSelector()},
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: appsv1.RollingUpdateStatefulSetStrategyType,
			},
			PodManagementPolicy: appsv1.ParallelPodManagement,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: util.MergeMaps(map[string]string{}, gen.GetLabels(), gen.GetSelector()),
				},
				Spec: corev1.PodSpec{
					ImagePullSecrets: func() []corev1.LocalObjectReference {
						if gen.Image.PullSecretName != nil {
							return []corev1.LocalObjectReference{{Name: *gen.Image.PullSecretName}}
						}
						return nil
					}(),
					Containers: []corev1.Container{
						{
							Name:  strings.Join([]string{component, console}, "-"),
							Image: fmt.Sprintf("%s:%s", *gen.Image.Name, *gen.Image.Tag),
							Args: []string{
								"sleep",
								"infinity",
							},
							Env:             pod.BuildEnvironment(gen.Options),
							Ports:           nil,
							Resources:       corev1.ResourceRequirements(*gen.Spec.Resources),
							ImagePullPolicy: *gen.Image.PullPolicy,
						},
					},
					Tolerations:                   gen.Spec.Tolerations,
					TerminationGracePeriodSeconds: util.Pointer[int64](30),
				},
			},
		},
	}

	sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes,
		corev1.Volume{
			Name: "system-config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: util.Pointer[int32](420),
					SecretName:  gen.ConfigFilesSecret,
				},
			},
		})
	sts.Spec.Template.Spec.Containers[0].VolumeMounts = append(sts.Spec.Template.Spec.Containers[0].VolumeMounts,
		corev1.VolumeMount{
			Name:      "system-config",
			ReadOnly:  true,
			MountPath: "/opt/system-extra-configs",
		})

	if gen.TwemproxySpec != nil {
		sts.Spec.Template = twemproxy.AddTwemproxySidecar(sts.Spec.Template, gen.TwemproxySpec)
	}

	return sts
}
