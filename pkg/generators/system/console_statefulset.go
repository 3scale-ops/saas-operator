package system

import (
	"fmt"
	"strings"

	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
	"github.com/3scale/saas-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

// StatefulSet returns a basereconciler.GeneratorFunction function that will return
// a StatefulSet resource when called
func (gen *ConsoleGenerator) statefulset() func() *appsv1.StatefulSet {

	return func() *appsv1.StatefulSet {

		sts := &appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      gen.GetComponent(),
				Namespace: gen.Namespace,
				Labels:    gen.GetLabels(),
			},
			Spec: appsv1.StatefulSetSpec{
				Replicas: pointer.Int32Ptr(1),
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
						RestartPolicy:   "Always",
						SecurityContext: &corev1.PodSecurityContext{},
						Containers: []corev1.Container{
							{
								Name:                     strings.Join([]string{component, console}, "-"),
								Image:                    fmt.Sprintf("%s:%s", *gen.Image.Name, *gen.Image.Tag),
								Command:                  []string{"sleep"},
								Args:                     []string{"infinity"},
								Env:                      pod.BuildEnvironment(gen.Options),
								Ports:                    nil,
								Resources:                corev1.ResourceRequirements(*gen.Spec.Resources),
								ImagePullPolicy:          *gen.Image.PullPolicy,
								TerminationMessagePath:   corev1.TerminationMessagePathDefault,
								TerminationMessagePolicy: corev1.TerminationMessageReadFile,
							},
						},
						Tolerations:                   gen.Spec.Tolerations,
						TerminationGracePeriodSeconds: pointer.Int64(30),
					},
				},
			},
		}

		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes,
			corev1.Volume{
				Name: "system-config",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						DefaultMode: pointer.Int32Ptr(420),
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

		return sts
	}
}
