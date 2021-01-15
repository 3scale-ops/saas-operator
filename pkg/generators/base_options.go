package generators

// BaseOptions configures the generators for a component
type BaseOptions struct {
	Component    string
	InstanceName string
	Namespace    string
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
