package redisshard

import (
	"fmt"
	"strings"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/pod"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StatefulSet returns a basereconciler.GeneratorFunction function that will return
// a StatefulSet resource when called
func (gen *Generator) StatefulSet() basereconciler.GeneratorFunction {

	return func() client.Object {
		return &appsv1.StatefulSet{
			TypeMeta: metav1.TypeMeta{
				Kind:       "StatefulSet",
				APIVersion: appsv1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", gen.GetComponent(), gen.GetInstanceName()),
				Namespace: gen.Namespace,
				Labels:    gen.GetLabels(),
			},
			Spec: appsv1.StatefulSetSpec{
				PodManagementPolicy:  appsv1.ParallelPodManagement,
				Replicas:             pointer.Int32(saasv1alpha1.RedisShardDefaultReplicas),
				RevisionHistoryLimit: pointer.Int32(1),
				Selector:             gen.Selector(),
				ServiceName:          fmt.Sprintf("%s-%s", gen.GetComponent(), gen.GetInstanceName()),
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: gen.LabelsWithSelector(),
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
								Command: []string{"redis-server", "/redis/redis.conf"},
								Image:   fmt.Sprintf("%s:%s", *gen.Image.Name, *gen.Image.Tag),
								LivenessProbe: &corev1.Probe{
									Handler: corev1.Handler{Exec: &corev1.ExecAction{
										Command: strings.Split("sh -c redis-cli -h $(hostname) ping", " "),
									}},
									FailureThreshold:    3,
									InitialDelaySeconds: 30,
									PeriodSeconds:       10,
									SuccessThreshold:    1,
									TimeoutSeconds:      5,
								},
								Name: "redis-server",
								Ports: pod.ContainerPorts(
									pod.ContainerPortTCP("redis-server", 6379),
								),
								ReadinessProbe: &corev1.Probe{
									Handler: corev1.Handler{Exec: &corev1.ExecAction{
										Command: strings.Split("/bin/sh /redis-readiness/ready.sh", " "),
									}},
									FailureThreshold:    3,
									InitialDelaySeconds: 30,
									PeriodSeconds:       10,
									SuccessThreshold:    1,
									TimeoutSeconds:      5,
								},
								ImagePullPolicy:          *gen.Image.PullPolicy,
								TerminationMessagePath:   corev1.TerminationMessagePathDefault,
								TerminationMessagePolicy: corev1.TerminationMessageReadFile,
								VolumeMounts: []corev1.VolumeMount{
									{Name: "redis-config", MountPath: "/redis"},
									{Name: "redis-readiness-script", MountPath: "/redis-readiness"},
									{Name: "redis-data", MountPath: "/data"},
								},
							},
						},
						SecurityContext: &corev1.PodSecurityContext{
							FSGroup:      pointer.Int64(1000),
							RunAsGroup:   pointer.Int64(1000),
							RunAsNonRoot: pointer.Bool(true),
							RunAsUser:    pointer.Int64(1000),
						},
						TerminationGracePeriodSeconds: pointer.Int64(0),
						Volumes: []corev1.Volume{
							{
								Name: "redis-config",
								VolumeSource: corev1.VolumeSource{
									ConfigMap: &corev1.ConfigMapVolumeSource{
										DefaultMode:          pointer.Int32(420),
										LocalObjectReference: corev1.LocalObjectReference{Name: "redis-config-" + gen.GetInstanceName()}},
								}},
							{
								Name: "redis-readiness-script",
								VolumeSource: corev1.VolumeSource{
									ConfigMap: &corev1.ConfigMapVolumeSource{
										DefaultMode:          pointer.Int32(484),
										LocalObjectReference: corev1.LocalObjectReference{Name: "redis-readiness-script-" + gen.GetInstanceName()}},
								}},
							{
								Name: "redis-data",
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{},
								}},
						},
					},
				},
			},
		}
	}
}
