package generators

import (
	"github.com/3scale/saas-operator/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

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
	// return a copy of the map and not a reference
	m := map[string]string{}
	for k, v := range bo.Labels {
		m[k] = v
	}
	return m
}

// Key returns a types.NamespacedName
func (bo *BaseOptions) Key() types.NamespacedName {
	return types.NamespacedName{Name: bo.GetComponent(), Namespace: bo.GetNamespace()}
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

// BaseOptions configures the generators for a component
type BaseOptionsV2 struct {
	Component    string
	InstanceName string
	Namespace    string
	Labels       map[string]string
}

// GetComponent returns the name of the component
func (bo *BaseOptionsV2) GetComponent() string {
	return bo.Component
}

// GetInstanceName returns the name of the custom resource instance
func (bo *BaseOptionsV2) GetInstanceName() string {
	return bo.InstanceName
}

// GetNamespace returns the custom resource namespace
func (bo *BaseOptionsV2) GetNamespace() string {
	return bo.Namespace
}

// GetLabels returns metadata labels
func (bo *BaseOptionsV2) GetLabels() map[string]string {
	// return a copy of the map and not a reference
	return util.MergeMaps(map[string]string{}, bo.Labels)
}

// Key returns a types.NamespacedName
func (bo *BaseOptionsV2) GetKey() types.NamespacedName {
	return types.NamespacedName{Name: bo.GetComponent(), Namespace: bo.GetNamespace()}
}

// GetSelector returns the LabelSelector struct that matches the labels in
func (bo *BaseOptionsV2) GetSelector() map[string]string {
	return map[string]string{PodSelectorKey: bo.GetComponent()}
}
