package system

import (
	"fmt"

	"github.com/3scale-ops/basereconciler/util"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/twemproxy"
	pipelinev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (gen *SystemTektonGenerator) task() *pipelinev1beta1.Task {
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
					Type: pipelinev1beta1.ParamTypeString,
				},
				{
					Name:        "container-tag",
					Description: "Container tag for the task",
					Default: &pipelinev1beta1.ParamValue{
						StringVal: fmt.Sprint(*gen.Image.Tag),
						Type:      pipelinev1beta1.ParamTypeString,
					},
					Type: pipelinev1beta1.ParamTypeString,
				},
			},
			StepTemplate: &pipelinev1beta1.StepTemplate{
				Image: "$(params.container-image):$(params.container-tag)",
				Env:   gen.Options.WithExtraEnv(gen.Spec.Config.ExtraEnv).BuildEnvironment(),
			},
			Steps: []pipelinev1beta1.Step{
				{
					Name:      "task-command",
					Command:   gen.Spec.Config.Command,
					Args:      gen.Spec.Config.Args,
					Resources: corev1.ResourceRequirements(*gen.Spec.Resources),
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "system-config",
							ReadOnly:  true,
							MountPath: "/opt/system-extra-configs",
						},
					},
					Timeout: gen.Spec.Config.Timeout,
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "system-config",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							DefaultMode: util.Pointer[int32](420),
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
