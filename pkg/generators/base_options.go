package generators

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	// PodSelectorKey is the label key used as Pod selector
	PodSelectorKey = "deployment"
)

// BaseOptions configures the generators for a component
type BaseOptions struct {
	Component    string
	InstanceName string
	Namespace    string
	Labels       map[string]string
}

// GetComponent returns the name of the component
func (bo *BaseOptions) GetComponent() string {
	return bo.Component
}

// GetInstanceName returns the name of the custom resource instance
func (bo *BaseOptions) GetInstanceName() string {
	return bo.InstanceName
}

// GetNamespace returns the custom resource namespace
func (bo *BaseOptions) GetNamespace() string {
	return bo.Namespace
}

// GetLabels returns metadata labels
func (bo *BaseOptions) GetLabels() map[string]string {
	return bo.Labels
}

// LabelsWithSelector returns Labels() with the addition of the Pod
// selector label
func (bo *BaseOptions) LabelsWithSelector() map[string]string {
	labels := bo.GetLabels()
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
