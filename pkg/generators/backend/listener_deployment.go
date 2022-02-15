package backend

import (
	"strings"

	"github.com/3scale/saas-operator/pkg/resource_builders/marin3r"
	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
	"github.com/3scale/saas-operator/pkg/resource_builders/twemproxy"
	"github.com/3scale/saas-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Deployment returns a function that will return a Deployment
// resource when called
func (gen *ListenerGenerator) deployment() func() *appsv1.Deployment {

	return func() *appsv1.Deployment {

		dep := &appsv1.Deployment{
			Spec: appsv1.DeploymentSpec{
				Replicas: gen.ListenerSpec.Replicas,
				Strategy: appsv1.DeploymentStrategy{
					Type: appsv1.RollingUpdateDeploymentStrategyType,
					RollingUpdate: &appsv1.RollingUpdateDeployment{
						MaxUnavailable: util.IntStrPtr(intstr.FromInt(0)),
						MaxSurge:       util.IntStrPtr(intstr.FromInt(1)),
					},
				},
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						ImagePullSecrets: pod.ImagePullSecrets(gen.Image.PullSecretName),
						Containers: []corev1.Container{
							{
								Name:  strings.Join([]string{component, listener}, "-"),
								Image: pod.Image(gen.Image),
								Args: func() (args []string) {
									if *gen.ListenerSpec.Config.RedisAsync {
										args = []string{"bin/3scale_backend", "-s", "falcon", "start"}
									} else {
										args = []string{"bin/3scale_backend", "start"}
									}
									args = append(args, "-e", "production", "-p", "3000", "-x", "/dev/stdout")
									return
								}(),
								Ports: pod.ContainerPorts(
									pod.ContainerPortTCP("http", 3000),
									pod.ContainerPortTCP("metrics", 9394),
								),
								Env:                      pod.BuildEnvironment(gen.Options),
								Resources:                corev1.ResourceRequirements(*gen.ListenerSpec.Resources),
								ImagePullPolicy:          *gen.Image.PullPolicy,
								LivenessProbe:            pod.TCPProbe(intstr.FromString("http"), *gen.ListenerSpec.LivenessProbe),
								ReadinessProbe:           pod.HTTPProbe("/status", intstr.FromString("http"), corev1.URISchemeHTTP, *gen.ListenerSpec.ReadinessProbe),
								TerminationMessagePath:   corev1.TerminationMessagePathDefault,
								TerminationMessagePolicy: corev1.TerminationMessageReadFile,
							},
						},
						Affinity:    pod.Affinity(gen.GetSelector(), gen.ListenerSpec.NodeAffinity),
						Tolerations: gen.ListenerSpec.Tolerations,
					},
				},
			},
		}

		if !gen.ListenerSpec.Marin3r.IsDeactivated() {
			dep = marin3r.EnableSidecar(*dep, *gen.ListenerSpec.Marin3r)
		}

		if gen.TwemproxySpec != nil {
			dep = twemproxy.AddTwemproxySidecar(*dep, gen.TwemproxySpec)
		}

		return dep
	}
}
