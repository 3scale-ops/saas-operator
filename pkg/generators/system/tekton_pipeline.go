package system

import (
	"fmt"

	pipelinev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (gen *SystemTektonGenerator) pipeline() *pipelinev1beta1.Pipeline {
	pipeline := &pipelinev1beta1.Pipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gen.GetComponent(),
			Namespace: gen.GetNamespace(),
			Labels:    gen.GetLabels(),
		},
		Spec: pipelinev1beta1.PipelineSpec{
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
			Tasks: []pipelinev1beta1.PipelineTask{
				{
					Name: *gen.Spec.Name,
					Params: pipelinev1beta1.Params{
						pipelinev1beta1.Param{
							Name: "container-image",
							Value: pipelinev1beta1.ParamValue{
								StringVal: "$(params.container-image)",
								Type:      pipelinev1beta1.ParamTypeString,
							},
						},
						pipelinev1beta1.Param{
							Name: "container-tag",
							Value: pipelinev1beta1.ParamValue{
								StringVal: "$(params.container-tag)",
								Type:      pipelinev1beta1.ParamTypeString,
							},
						},
					},
					TaskRef: &pipelinev1beta1.TaskRef{
						Name: gen.GetComponent(),
						Kind: pipelinev1beta1.NamespacedTaskKind,
					},
				},
			},
		},
	}
	return pipeline
}
