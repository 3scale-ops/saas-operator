package envoyconfig

// EnvoyDynamicConfigDescriptor is a struct that contains
// information to generate an Envoy dynamic configuration
type EnvoyDynamicConfigDescriptor interface {
	GetGeneratorVersion() string
	GetName() string
	GetOptions() interface{}
}
