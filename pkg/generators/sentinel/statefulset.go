package sentinel

import (
	"fmt"
	"strings"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
	"github.com/3scale/saas-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

var (
	healthCommand string = fmt.Sprintf("redis-cli -p %d PING", saasv1alpha1.SentinelPort)
)

// statefulSet returns a function that will return
// a StatefulSet resource when called
func (gen *Generator) statefulSet() func() *appsv1.StatefulSet {

	return func() *appsv1.StatefulSet {
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
				PodManagementPolicy: appsv1.ParallelPodManagement,
				Replicas:            gen.Spec.Replicas,
				Selector:            &metav1.LabelSelector{MatchLabels: gen.GetSelector()},
				ServiceName:         gen.GetComponent() + "-headless",
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: util.MergeMaps(gen.GetLabels(), gen.GetSelector()),
					},
					Spec: corev1.PodSpec{
						Affinity:                     pod.Affinity(gen.GetSelector(), gen.Spec.NodeAffinity),
						AutomountServiceAccountToken: pointer.Bool(false),
						DNSPolicy:                    corev1.DNSClusterFirst,
						ImagePullSecrets: func() []corev1.LocalObjectReference {
							if gen.Spec.Image.PullSecretName != nil {
								return []corev1.LocalObjectReference{{Name: *gen.Spec.Image.PullSecretName}}
							}
							return nil
						}(),
						Containers: []corev1.Container{
							{
								Command:         []string{"redis-server", "/redis/sentinel.conf", "--sentinel"},
								Image:           fmt.Sprintf("%s:%s", *gen.Spec.Image.Name, *gen.Spec.Image.Tag),
								ImagePullPolicy: *gen.Spec.Image.PullPolicy,
								LivenessProbe: &corev1.Probe{
									ProbeHandler: corev1.ProbeHandler{Exec: &corev1.ExecAction{
										Command: strings.Split(healthCommand, " ")}},
									FailureThreshold:    3,
									InitialDelaySeconds: 30,
									PeriodSeconds:       10,
									SuccessThreshold:    1,
									TimeoutSeconds:      5,
								},
								Name: gen.GetComponent(),
								Ports: pod.ContainerPorts(
									pod.ContainerPortTCP(gen.GetComponent(), int32(saasv1alpha1.SentinelPort)),
								),
								ReadinessProbe: &corev1.Probe{
									ProbeHandler: corev1.ProbeHandler{Exec: &corev1.ExecAction{
										Command: strings.Split(healthCommand, " ")}},
									FailureThreshold:    3,
									InitialDelaySeconds: 30,
									PeriodSeconds:       10,
									SuccessThreshold:    1,
									TimeoutSeconds:      5,
								},
								Resources:                corev1.ResourceRequirements(*gen.Spec.Resources),
								TerminationMessagePath:   corev1.TerminationMessagePathDefault,
								TerminationMessagePolicy: corev1.TerminationMessageReadFile,
								VolumeMounts: []corev1.VolumeMount{
									{Name: gen.GetComponent() + "-config-rw", MountPath: "/redis"},
								},
							},
						},
						InitContainers: []corev1.Container{
							{
								Command: strings.Split("sh /redis-ro/generate-config.sh /redis/sentinel.conf", " "),
								Env: []corev1.EnvVar{{
									Name: "POD_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath:  "status.podIP",
											APIVersion: corev1.SchemeGroupVersion.Version,
										},
									},
								}},
								Image:                    fmt.Sprintf("%s:%s", *gen.Spec.Image.Name, *gen.Spec.Image.Tag),
								ImagePullPolicy:          *gen.Spec.Image.PullPolicy,
								Name:                     gen.GetComponent() + "-gen-config",
								TerminationMessagePath:   corev1.TerminationMessagePathDefault,
								TerminationMessagePolicy: corev1.TerminationMessageReadFile,
								VolumeMounts: []corev1.VolumeMount{
									{Name: gen.GetComponent() + "-gen-config", MountPath: "/redis-ro"},
									{Name: gen.GetComponent() + "-config-rw", MountPath: "/redis"},
								},
							},
						},
						Tolerations:                   gen.Spec.Tolerations,
						TerminationGracePeriodSeconds: pointer.Int64(30),
						Volumes: []corev1.Volume{
							{
								Name: gen.GetComponent() + "-gen-config",
								VolumeSource: corev1.VolumeSource{
									ConfigMap: &corev1.ConfigMapVolumeSource{
										DefaultMode:          pointer.Int32(484),
										LocalObjectReference: corev1.LocalObjectReference{Name: gen.GetComponent() + "-gen-config"}},
								}},
						}},
				},
				UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
					Type: appsv1.RollingUpdateStatefulSetStrategyType,
					RollingUpdate: &appsv1.RollingUpdateStatefulSetStrategy{
						Partition: pointer.Int32(0),
					},
				},
				VolumeClaimTemplates: []corev1.PersistentVolumeClaim{{
					ObjectMeta: metav1.ObjectMeta{
						Name: gen.GetComponent() + "-config-rw",
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
						Resources:        corev1.ResourceRequirements{Requests: corev1.ResourceList{corev1.ResourceStorage: *gen.Spec.Config.StorageSize}},
						StorageClassName: gen.Spec.Config.StorageClass,
						VolumeMode:       (*corev1.PersistentVolumeMode)(pointer.StringPtr(string(corev1.PersistentVolumeFilesystem))),
						DataSource:       &corev1.TypedLocalObjectReference{},
					},
					Status: corev1.PersistentVolumeClaimStatus{
						Phase: corev1.ClaimPending,
					},
				}},
			},
		}
	}
}
