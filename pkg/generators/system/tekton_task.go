package system

import (
	"fmt"

	"github.com/3scale/saas-operator/pkg/resource_builders/pod"
	"github.com/3scale/saas-operator/pkg/resource_builders/twemproxy"
	pipelinev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

// task returns a basereconciler.GeneratorFunction function that will return a
// Tekton Task resource when called
func (gen *SystemTektonGenerator) task() func() *pipelinev1beta1.Task {

	return func() *pipelinev1beta1.Task {

		task := &pipelinev1beta1.Task{
			ObjectMeta: metav1.ObjectMeta{
				Name:      gen.GetComponent(),
				Namespace: gen.GetNamespace(),
				Labels:    gen.GetLabels(),
			},
			Spec: pipelinev1beta1.TaskSpec{
				DisplayName: gen.GetComponent(),
				Description: *gen.Spec.Description,
				Params: []pipelinev1beta1.ParamSpec{
					{
						Name:        "container-image",
						Description: "Container image for the task",
						Default: &pipelinev1beta1.ParamValue{
							StringVal: fmt.Sprint(*gen.Image.Name),
							Type:      pipelinev1beta1.ParamTypeString,
						},
					},
					{
						Name:        "container-tag",
						Description: "Container tag for the task",
						Default: &pipelinev1beta1.ParamValue{
							StringVal: fmt.Sprint(*gen.Image.Tag),
							Type:      pipelinev1beta1.ParamTypeString,
						},
					},
				},
				StepTemplate: &pipelinev1beta1.StepTemplate{
					Env: append(
						pod.BuildEnvironment(gen.Options),
						gen.Spec.Config.ExtraEnv...,
					),
				},
				Steps: []pipelinev1beta1.Step{
					{
						Name:      "task-command",
						Command:   gen.Spec.Config.Command,
						Args:      gen.Spec.Config.Args,
						Image:     "$(params.container-image):$(params.container-tag)",
						Resources: corev1.ResourceRequirements(*gen.Spec.Resources),
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "system-config",
								ReadOnly:  true,
								MountPath: "/opt/system-extra-configs",
							},
						},
					},
				},
				Volumes: []corev1.Volume{
					{
						Name: "system-config",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								DefaultMode: pointer.Int32(420),
								SecretName:  gen.ConfigFilesSecret,
							},
						},
					},
				},
			},
		}

		if gen.TwemproxySpec != nil {
			task.Spec = twemproxy.AddTwemproxyTaskSidecar(task.Spec, gen.TwemproxySpec)
		}

		return task
	}
}
