/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"github.com/3scale-ops/basereconciler/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	backendDefaultImage defaultImageSpec = defaultImageSpec{
		Name:       util.Pointer("quay.io/3scale/apisonator"),
		Tag:        util.Pointer("nightly"),
		PullPolicy: (*corev1.PullPolicy)(util.Pointer(string(corev1.PullIfNotPresent))),
	}
	backendDefaultGrafanaDashboard defaultGrafanaDashboardSpec = defaultGrafanaDashboardSpec{
		SelectorKey:   util.Pointer("monitoring-key"),
		SelectorValue: util.Pointer("middleware"),
	}
	backendDefaultConfigRackEnv         string                             = "dev"
	backendDefaultConfigMasterServiceID int32                              = 6
	backendDefaultListenerHPA           defaultHorizontalPodAutoscalerSpec = defaultHorizontalPodAutoscalerSpec{
		MinReplicas:         util.Pointer[int32](2),
		MaxReplicas:         util.Pointer[int32](4),
		ResourceUtilization: util.Pointer[int32](90),
		ResourceName:        util.Pointer("cpu"),
	}
	backendDefaultListenerPDB defaultPodDisruptionBudgetSpec = defaultPodDisruptionBudgetSpec{
		MaxUnavailable: util.Pointer(intstr.FromInt(1)),
	}
	backendDefaultListenerNLBLoadBalancer defaultNLBLoadBalancerSpec = defaultNLBLoadBalancerSpec{
		ProxyProtocol:                 util.Pointer(true),
		CrossZoneLoadBalancingEnabled: util.Pointer(true),
		TerminationProtection:         util.Pointer(false),
	}
	backendDefaultListenerReplicas  int32                           = 2
	backendDefaultListenerResources defaultResourceRequirementsSpec = defaultResourceRequirementsSpec{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("500m"),
			corev1.ResourceMemory: resource.MustParse("550Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("1"),
			corev1.ResourceMemory: resource.MustParse("700Mi"),
		},
	}
	backendDefaultListenerLivenessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: util.Pointer[int32](30),
		TimeoutSeconds:      util.Pointer[int32](1),
		PeriodSeconds:       util.Pointer[int32](10),
		SuccessThreshold:    util.Pointer[int32](1),
		FailureThreshold:    util.Pointer[int32](3),
	}
	backendDefaultListenerReadinessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: util.Pointer[int32](30),
		TimeoutSeconds:      util.Pointer[int32](5),
		PeriodSeconds:       util.Pointer[int32](10),
		SuccessThreshold:    util.Pointer[int32](1),
		FailureThreshold:    util.Pointer[int32](3),
	}
	backendDefaultListenerMarin3rSpec                 defaultMarin3rSidecarSpec          = defaultMarin3rSidecarSpec{}
	backendDefaultListenerConfigLogFormat             string                             = "json"
	backendDefaultListenerConfigRedisAsync            bool                               = false
	backendDefaultListenerConfigListenerWorkers       int32                              = 16
	backendDefaultListenerConfigLegacyReferrerFilters bool                               = true
	backendDefaultWorkerHPA                           defaultHorizontalPodAutoscalerSpec = defaultHorizontalPodAutoscalerSpec{
		MinReplicas:         util.Pointer[int32](2),
		MaxReplicas:         util.Pointer[int32](4),
		ResourceUtilization: util.Pointer[int32](90),
		ResourceName:        util.Pointer("cpu"),
	}
	backendDefaultWorkerPDB defaultPodDisruptionBudgetSpec = defaultPodDisruptionBudgetSpec{
		MaxUnavailable: util.Pointer(intstr.FromInt(1)),
	}
	backendDefaultWorkerReplicas  int32                           = 2
	backendDefaultWorkerResources defaultResourceRequirementsSpec = defaultResourceRequirementsSpec{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("150m"),
			corev1.ResourceMemory: resource.MustParse("50Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("1"),
			corev1.ResourceMemory: resource.MustParse("300Mi"),
		},
	}
	backendDefaultWorkerLivenessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: util.Pointer[int32](10),
		TimeoutSeconds:      util.Pointer[int32](3),
		PeriodSeconds:       util.Pointer[int32](15),
		SuccessThreshold:    util.Pointer[int32](1),
		FailureThreshold:    util.Pointer[int32](3),
	}
	backendDefaultWorkerReadinessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: util.Pointer[int32](10),
		TimeoutSeconds:      util.Pointer[int32](5),
		PeriodSeconds:       util.Pointer[int32](30),
		SuccessThreshold:    util.Pointer[int32](1),
		FailureThreshold:    util.Pointer[int32](3),
	}
	backendDefaultWorkerConfigLogFormat  string                          = "json"
	backendDefaultWorkerConfigRedisAsync bool                            = false
	backendDefaultCronReplicas           int32                           = 1
	backendDefaultCronResources          defaultResourceRequirementsSpec = defaultResourceRequirementsSpec{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("50m"),
			corev1.ResourceMemory: resource.MustParse("50Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("150m"),
			corev1.ResourceMemory: resource.MustParse("150Mi"),
		},
	}
)

// BackendSpec defines the desired state of Backend
type BackendSpec struct {
	// Image specification for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Image *ImageSpec `json:"image,omitempty"`
	// Application specific configuration options for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Config BackendConfig `json:"config"`
	// Configures the Grafana Dashboard for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	GrafanaDashboard *GrafanaDashboardSpec `json:"grafanaDashboard,omitempty"`
	// Configures the backend listener
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Listener ListenerSpec `json:"listener"`
	// Configures the backend worker
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Worker *WorkerSpec `json:"worker,omitempty"`
	// Configures the backend cron
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Cron *CronSpec `json:"cron,omitempty"`
	// Configures twemproxy
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Twemproxy *TwemproxySpec `json:"twemproxy,omitempty"`
}

// Default implements defaulting for BackendSpec
func (spec *BackendSpec) Default() {

	spec.Image = InitializeImageSpec(spec.Image, backendDefaultImage)
	spec.Config.Default()
	spec.Listener.Default()
	if spec.Worker == nil {
		spec.Worker = &WorkerSpec{}
	}
	spec.Worker.Default()
	if spec.Cron == nil {
		spec.Cron = &CronSpec{}
	}
	spec.Cron.Default()
	spec.GrafanaDashboard = InitializeGrafanaDashboardSpec(spec.GrafanaDashboard, backendDefaultGrafanaDashboard)
	if spec.Twemproxy != nil {
		spec.Twemproxy.Default()
	}
}

// ResolveCanarySpec modifies the BackendSpec given the provided canary configuration
func (spec *BackendSpec) ResolveCanarySpec(canary *Canary) (*BackendSpec, error) {
	canarySpec := &BackendSpec{}
	if err := canary.PatchSpec(spec, canarySpec); err != nil {
		return nil, err
	}

	if canary.ImageName != nil {
		canarySpec.Image.Name = canary.ImageName
	}
	if canary.ImageTag != nil {
		canarySpec.Image.Tag = canary.ImageTag
	}
	canarySpec.Listener.Replicas = canary.Replicas
	canarySpec.Worker.Replicas = canary.Replicas
	canarySpec.Cron.Replicas = canary.Replicas

	// Call Default() on the resolved canary spec to apply
	// defaulting to potentially added fields
	canarySpec.Default()

	return canarySpec, nil
}

// ListenerSpec is the configuration for Backend Listener
type ListenerSpec struct {
	// Listener specific configuration options for the component element
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Config *ListenerConfig `json:"config,omitempty"`
	// Pod Disruption Budget for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	PDB *PodDisruptionBudgetSpec `json:"pdb,omitempty"`
	// Horizontal Pod Autoscaler for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	HPA *HorizontalPodAutoscalerSpec `json:"hpa,omitempty"`
	// Number of replicas (ignored if hpa is enabled) for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// Resource requirements for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Resources *ResourceRequirementsSpec `json:"resources,omitempty"`
	// Liveness probe for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	LivenessProbe *ProbeSpec `json:"livenessProbe,omitempty"`
	// Readiness probe for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ReadinessProbe *ProbeSpec `json:"readinessProbe,omitempty"`
	// The external endpoint/s for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Endpoint Endpoint `json:"endpoint"`
	// Marin3r configures the Marin3r sidecars for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Marin3r *Marin3rSidecarSpec `json:"marin3r,omitempty"`
	// Configures the AWS Network load balancer for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	LoadBalancer *NLBLoadBalancerSpec `json:"loadBalancer,omitempty"`
	// Describes node affinity scheduling rules for the pod.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	NodeAffinity *corev1.NodeAffinity `json:"nodeAffinity,omitempty" protobuf:"bytes,1,opt,name=nodeAffinity"`
	// If specified, the pod's tolerations.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty" protobuf:"bytes,22,opt,name=tolerations"`
	// Canary defines spec changes for the canary Deployment. If
	// left unset the canary Deployment wil not be created.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Canary *Canary `json:"canary,omitempty"`
}

// Default implements defaulting for the each backend listener
func (spec *ListenerSpec) Default() {

	spec.HPA = InitializeHorizontalPodAutoscalerSpec(spec.HPA, backendDefaultListenerHPA)
	spec.Replicas = intOrDefault(spec.Replicas, &backendDefaultListenerReplicas)
	spec.PDB = InitializePodDisruptionBudgetSpec(spec.PDB, backendDefaultListenerPDB)
	spec.Resources = InitializeResourceRequirementsSpec(spec.Resources, backendDefaultListenerResources)
	spec.LivenessProbe = InitializeProbeSpec(spec.LivenessProbe, backendDefaultListenerLivenessProbe)
	spec.ReadinessProbe = InitializeProbeSpec(spec.ReadinessProbe, backendDefaultListenerReadinessProbe)
	spec.LoadBalancer = InitializeNLBLoadBalancerSpec(spec.LoadBalancer, backendDefaultListenerNLBLoadBalancer)
	spec.Marin3r = InitializeMarin3rSidecarSpec(spec.Marin3r, backendDefaultListenerMarin3rSpec)
	if spec.Config == nil {
		spec.Config = &ListenerConfig{}
	}
	spec.Config.Default()
}

// WorkerSpec is the configuration for Backend Worker
type WorkerSpec struct {
	// Listener specific configuration options for the component element
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Config *WorkerConfig `json:"config,omitempty"`
	// Pod Disruption Budget for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	PDB *PodDisruptionBudgetSpec `json:"pdb,omitempty"`
	// Horizontal Pod Autoscaler for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	HPA *HorizontalPodAutoscalerSpec `json:"hpa,omitempty"`
	// Number of replicas (ignored if hpa is enabled) for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// Resource requirements for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Resources *ResourceRequirementsSpec `json:"resources,omitempty"`
	// Liveness probe for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	LivenessProbe *ProbeSpec `json:"livenessProbe,omitempty"`
	// Readiness probe for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ReadinessProbe *ProbeSpec `json:"readinessProbe,omitempty"`
	// Describes node affinity scheduling rules for the pod.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	NodeAffinity *corev1.NodeAffinity `json:"nodeAffinity,omitempty" protobuf:"bytes,1,opt,name=nodeAffinity"`
	// If specified, the pod's tolerations.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty" protobuf:"bytes,22,opt,name=tolerations"`
	// Canary defines spec changes for the canary Deployment. If
	// left unset the canary Deployment wil not be created.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Canary *Canary `json:"canary,omitempty"`
}

// Default implements defaulting for the each backend worker
func (spec *WorkerSpec) Default() {

	spec.HPA = InitializeHorizontalPodAutoscalerSpec(spec.HPA, backendDefaultWorkerHPA)
	spec.Replicas = intOrDefault(spec.Replicas, &backendDefaultWorkerReplicas)
	spec.PDB = InitializePodDisruptionBudgetSpec(spec.PDB, backendDefaultWorkerPDB)
	spec.Resources = InitializeResourceRequirementsSpec(spec.Resources, backendDefaultWorkerResources)
	spec.LivenessProbe = InitializeProbeSpec(spec.LivenessProbe, backendDefaultWorkerLivenessProbe)
	spec.ReadinessProbe = InitializeProbeSpec(spec.ReadinessProbe, backendDefaultWorkerReadinessProbe)
	if spec.Config == nil {
		spec.Config = &WorkerConfig{}
	}
	spec.Config.Default()
}

// CronSpec is the configuration for Backend Cron
type CronSpec struct {
	// Number of replicas for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// Resource requirements for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Resources *ResourceRequirementsSpec `json:"resources,omitempty"`
	// Describes node affinity scheduling rules for the pod.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	NodeAffinity *corev1.NodeAffinity `json:"nodeAffinity,omitempty" protobuf:"bytes,1,opt,name=nodeAffinity"`
	// If specified, the pod's tolerations.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty" protobuf:"bytes,22,opt,name=tolerations"`
}

// Default implements defaulting for the each backend cron
func (spec *CronSpec) Default() {

	spec.Replicas = intOrDefault(spec.Replicas, &backendDefaultCronReplicas)
	spec.Resources = InitializeResourceRequirementsSpec(spec.Resources, backendDefaultCronResources)
}

// BackendConfig configures app behavior for Backend
type BackendConfig struct {
	// Rack environment
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	RackEnv *string `json:"rackEnv,omitempty"`
	// Master service account ID in Porta
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	MasterServiceID *int32 `json:"masterServiceID,omitempty"`
	// Redis Storage DSN
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	RedisStorageDSN string `json:"redisStorageDSN"`
	// Redis Queues DSN
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	RedisQueuesDSN string `json:"redisQueuesDSN"`
	// External Secret common configuration
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ExternalSecret ExternalSecret `json:"externalSecret,omitempty"`
	// A reference to the secret holding the backend-system-events-hook URL
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	SystemEventsHookURL SecretReference `json:"systemEventsHookURL"`
	// A reference to the secret holding the backend-system-events-hook password
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	SystemEventsHookPassword SecretReference `json:"systemEventsHookPassword"`
	// A reference to the secret holding the backend-internal-api user
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	InternalAPIUser SecretReference `json:"internalAPIUser"`
	// A reference to the secret holding the backend-internal-api password
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	InternalAPIPassword SecretReference `json:"internalAPIPassword"`
	// A reference to the secret holding the backend-error-monitoring service
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ErrorMonitoringService *SecretReference `json:"errorMonitoringService,omitempty"`
	// A reference to the secret holding the backend-error-monitoring key
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ErrorMonitoringKey *SecretReference `json:"errorMonitoringKey,omitempty"`
}

// Default sets default values for any value not specifically set in the BackendConfig struct
func (cfg *BackendConfig) Default() {
	cfg.RackEnv = stringOrDefault(cfg.RackEnv, util.Pointer(backendDefaultConfigRackEnv))
	cfg.MasterServiceID = intOrDefault(cfg.MasterServiceID, util.Pointer[int32](backendDefaultConfigMasterServiceID))
	cfg.ExternalSecret.SecretStoreRef = InitializeExternalSecretSecretStoreReferenceSpec(cfg.ExternalSecret.SecretStoreRef, defaultExternalSecretSecretStoreReference)
	cfg.ExternalSecret.RefreshInterval = durationOrDefault(cfg.ExternalSecret.RefreshInterval, &defaultExternalSecretRefreshInterval)
}

// ListenerConfig configures app behavior for Backend Listener
type ListenerConfig struct {
	// Listener log format
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:validation:Enum=test;json
	// +optional
	LogFormat *string `json:"logFormat,omitempty"`
	// Enable (true) or disable (false) listener redis async mode
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	RedisAsync *bool `json:"redisAsync,omitempty"`
	// Number of worker processes per listener pod
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ListenerWorkers *int32 `json:"listenerWorkers,omitempty"`
	// Enable (true) or disable (false) Legacy Referrer Filters
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	LegacyReferrerFilters *bool `json:"legacyReferrerFilters,omitempty"`
}

// Default sets default values for any value not specifically set in the ListenerConfig struct
func (cfg *ListenerConfig) Default() {
	cfg.LogFormat = stringOrDefault(cfg.LogFormat, util.Pointer(backendDefaultListenerConfigLogFormat))
	cfg.RedisAsync = boolOrDefault(cfg.RedisAsync, util.Pointer(backendDefaultListenerConfigRedisAsync))
	cfg.ListenerWorkers = intOrDefault(cfg.ListenerWorkers, util.Pointer[int32](backendDefaultListenerConfigListenerWorkers))
	cfg.LegacyReferrerFilters = boolOrDefault(cfg.LegacyReferrerFilters, util.Pointer(backendDefaultListenerConfigLegacyReferrerFilters))
}

// WorkerConfig configures app behavior for Backend Worker
type WorkerConfig struct {
	// Worker log format
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:validation:Enum=test;json
	// +optional
	LogFormat *string `json:"logFormat,omitempty"`
	// Enable (true) or disable (false) worker redis async mode
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	RedisAsync *bool `json:"redisAsync,omitempty"`
}

// Default sets default values for any value not specifically set in the WorkerConfig struct
func (cfg *WorkerConfig) Default() {
	cfg.LogFormat = stringOrDefault(cfg.LogFormat, util.Pointer(backendDefaultWorkerConfigLogFormat))
	cfg.RedisAsync = boolOrDefault(cfg.RedisAsync, util.Pointer(backendDefaultWorkerConfigRedisAsync))
}

// BackendStatus defines the observed state of Backend
type BackendStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Backend is the Schema for the backends API
type Backend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackendSpec   `json:"spec,omitempty"`
	Status BackendStatus `json:"status,omitempty"`
}

// Defaults impletements defaulting for the Apicast resource
func (b *Backend) Default() {
	b.Spec.Default()
}

// +kubebuilder:object:root=true

// BackendList contains a list of Backend
type BackendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Backend `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Backend{}, &BackendList{})
}
