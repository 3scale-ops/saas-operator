package twemproxy

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
	"github.com/go-test/deep"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func Test_AddTwemproxySidecar(t *testing.T) {
	type args struct {
		dep  appsv1.Deployment
		spec *saasv1alpha1.TwemproxySpec
	}
	tests := []struct {
		name string
		args args
		want corev1.PodTemplateSpec
	}{
		{
			name: "Adds twemproxy sidecar container to a Deployment",
			args: args{
				dep: appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Template: corev1.PodTemplateSpec{
							ObjectMeta: metav1.ObjectMeta{},
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{{
									Name: "test",
								}},
							},
						},
					},
				},
				spec: &saasv1alpha1.TwemproxySpec{
					Image: &saasv1alpha1.ImageSpec{
						Name:       pointer.String("twemproxy"),
						Tag:        pointer.String("latest"),
						PullPolicy: (*corev1.PullPolicy)(pointer.String(string(corev1.PullIfNotPresent))),
					},
					Resources: &saasv1alpha1.ResourceRequirementsSpec{},
					LivenessProbe: &saasv1alpha1.ProbeSpec{
						InitialDelaySeconds: pointer.Int32(1),
						TimeoutSeconds:      pointer.Int32(3),
						PeriodSeconds:       pointer.Int32(5),
						SuccessThreshold:    pointer.Int32(1),
						FailureThreshold:    pointer.Int32(3),
					},
					ReadinessProbe: &saasv1alpha1.ProbeSpec{
						InitialDelaySeconds: pointer.Int32(1),
						TimeoutSeconds:      pointer.Int32(3),
						PeriodSeconds:       pointer.Int32(5),
						SuccessThreshold:    pointer.Int32(1),
						FailureThreshold:    pointer.Int32(3),
					},
					TwemproxyConfigRef: "twem-config",
					Options: &saasv1alpha1.TwemproxyOptions{
						LogLevel:      pointer.Int32(6),
						StatsInterval: &metav1.Duration{Duration: 20 * time.Second},
						MetricsPort:   pointer.Int32(5555),
					},
				},
			},
			want: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"saas.3scale.net/twemproxyconfig.sync": "twem-config",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "test",
						},
						{
							Env: []corev1.EnvVar{
								{Name: "TWEMPROXY_CONFIG_FILE", Value: TwemproxyConfigFile},
								{Name: "TWEMPROXY_METRICS_ADDRESS", Value: ":5555"},
								{Name: "TWEMPROXY_STATS_INTERVAL", Value: "20000"},
								{Name: "TWEMPROXY_LOG_LEVEL", Value: "6"},
							},
							Name:  twemproxy,
							Image: "twemproxy:latest",
							Ports: pod.ContainerPorts(
								pod.ContainerPortTCP(twemproxy, 22121),
								pod.ContainerPortTCP("twem-metrics", 5555),
							),
							Resources:       corev1.ResourceRequirements{},
							ImagePullPolicy: corev1.PullIfNotPresent,
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{Exec: &corev1.ExecAction{
									Command: strings.Split(healthCommand, " "),
								}},
								InitialDelaySeconds: *pointer.Int32(1),
								TimeoutSeconds:      *pointer.Int32(3),
								PeriodSeconds:       *pointer.Int32(5),
								SuccessThreshold:    *pointer.Int32(1),
								FailureThreshold:    *pointer.Int32(3),
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{Exec: &corev1.ExecAction{
									Command: strings.Split(healthCommand, " "),
								}},
								InitialDelaySeconds: *pointer.Int32(1),
								TimeoutSeconds:      *pointer.Int32(3),
								PeriodSeconds:       *pointer.Int32(5),
								SuccessThreshold:    *pointer.Int32(1),
								FailureThreshold:    *pointer.Int32(3),
							},
							Lifecycle: &corev1.Lifecycle{
								PreStop: &corev1.LifecycleHandler{
									Exec: &corev1.ExecAction{
										Command: []string{"pre-stop", TwemproxyConfigFile},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "twemproxy-config",
									MountPath: filepath.Dir(TwemproxyConfigFile),
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: twemproxy + "-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{Name: "twem-config"},
									DefaultMode:          pointer.Int32(420),
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddTwemproxySidecar(tt.args.dep.Spec.Template, tt.args.spec)
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("AddTwemproxySidecar() = diff %v", diff)
			}
		})
	}
}
