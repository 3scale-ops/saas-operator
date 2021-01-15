package generators

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	// PodSelectorKey is the label key used as Pod selector
	PodSelectorKey = "deployment"
)

// Labels returns metadata labels
func (bo *BaseOptions) Labels() map[string]string {
	return map[string]string{
		"app":     bo.GetComponent(),
		"part-of": "3scale-saas",
	}
}

// LabelsWithSelector returns Labels() with the addition of the Pod
// selector label
func (bo *BaseOptions) LabelsWithSelector() map[string]string {
	labels := bo.Labels()
	labels[PodSelectorKey] = bo.GetComponent()
	return labels
}

// Selector returns the LabelSelector struct that matches the labels in
// LabelsWithSelector()
func (bo *BaseOptions) Selector() *metav1.LabelSelector {
	return &metav1.LabelSelector{
		MatchLabels: map[string]string{PodSelectorKey: bo.GetComponent()},
	}
}
