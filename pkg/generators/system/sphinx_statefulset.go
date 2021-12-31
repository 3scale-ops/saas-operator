package system

import (
	"fmt"

	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pod"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StatefulSet returns a basereconciler.GeneratorFunction function that will return
// a StatefulSet resource when called
func (gen *SphinxGenerator) StatefulSet() basereconciler.GeneratorFunction {

	return func() client.Object {
		return &appsv1.StatefulSet{
			TypeMeta: metav1.TypeMeta{
				Kind:       "StatefulSet",
				APIVersion: appsv1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      gen.GetComponent(),
				Namespace: gen.Namespace,
				Labels:    gen.GetLabels(),
			},
			Spec: appsv1.StatefulSetSpec{
				Replicas: pointer.Int32Ptr(1),
				Selector: gen.Selector(),
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
									"rake",
									"openshift:thinking_sphinx:start",
								},
								Env: pod.BuildEnvironment(gen.Options),
								Ports: pod.ContainerPorts(
									pod.ContainerPortTCP("sphinx", 9306),
								),
								Resources:                corev1.ResourceRequirements(*gen.Spec.Resources),
								LivenessProbe:            pod.TCPProbe(intstr.FromString("sphinx"), *gen.Spec.LivenessProbe),
								ReadinessProbe:           pod.TCPProbe(intstr.FromString("sphinx"), *gen.Spec.ReadinessProbe),
								ImagePullPolicy:          *gen.ImageSpec.PullPolicy,
								TerminationMessagePath:   corev1.TerminationMessagePathDefault,
								TerminationMessagePolicy: corev1.TerminationMessageReadFile,
								VolumeMounts: []corev1.VolumeMount{{
									Name:      "system-sphinx-database",
									MountPath: gen.DatabasePath,
								}},
							},
						},
						Affinity:    pod.Affinity(gen.Selector().MatchLabels, gen.Spec.NodeAffinity),
						Tolerations: gen.Spec.Tolerations,
					},
				},
				VolumeClaimTemplates: []corev1.PersistentVolumeClaim{{
					TypeMeta: metav1.TypeMeta{
						Kind:       "PersistentVolumeClaim",
						APIVersion: corev1.SchemeGroupVersion.String(),
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: "system-sphinx-database",
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
						Resources:        corev1.ResourceRequirements{Requests: corev1.ResourceList{corev1.ResourceStorage: gen.DatabaseStorageSize}},
						StorageClassName: gen.DatabaseStorageClass,
						VolumeMode:       (*corev1.PersistentVolumeMode)(pointer.StringPtr(string(corev1.PersistentVolumeFilesystem))),
						DataSource:       &corev1.TypedLocalObjectReference{},
					},
				}},
			},
		}
	}
}
