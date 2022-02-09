package system

import (
	"fmt"
	"strings"

	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
	"github.com/3scale/saas-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

// StatefulSet returns a basereconciler.GeneratorFunction function that will return
// a StatefulSet resource when called
func (gen *SphinxGenerator) statefulset() func() *appsv1.StatefulSet {

	return func() *appsv1.StatefulSet {

		return &appsv1.StatefulSet{
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
				PodManagementPolicy: appsv1.OrderedReadyPodManagement,
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
								Name:  strings.Join([]string{component, sphinx}, "-"),
								Image: fmt.Sprintf("%s:%s", *gen.Image.Name, *gen.Image.Tag),
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
								ImagePullPolicy:          *gen.Image.PullPolicy,
								TerminationMessagePath:   corev1.TerminationMessagePathDefault,
								TerminationMessagePolicy: corev1.TerminationMessageReadFile,
								VolumeMounts: []corev1.VolumeMount{{
									Name:      "system-sphinx-database",
									MountPath: gen.DatabasePath,
								}},
							},
						},
						Affinity:    pod.Affinity(gen.GetSelector(), gen.Spec.NodeAffinity),
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
					Status: corev1.PersistentVolumeClaimStatus{
						Phase: corev1.ClaimPending,
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
