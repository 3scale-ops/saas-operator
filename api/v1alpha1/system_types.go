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
	"github.com/3scale/saas-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type systemSidekiqType string

const (
	Default systemSidekiqType = "default"
	Billing systemSidekiqType = "billing"
	Low     systemSidekiqType = "low"
)

var (
	// Common
	systemDefaultSandboxProxyOpensslVerifyMode string           = "VERIFY_NONE"
	systemDefaultForceSSL                      bool             = true
	systemDefaultSSLCertsDir                   string           = "/etc/pki/tls/certs"
	systemDefaultThreescaleProviderPlan        string           = "enterprise"
	systemDefaultThreescaleSuperdomain         string           = "localhost"
	systemDefaultRailsConsole                  bool             = false
	systemDefaultRailsEnvironment              string           = "preview"
	systemDefaultRailsLogLevel                 string           = "info"
	systemDefaultConfigFilesSecret             string           = "system-config"
	systemDefaultBugsnagSpec                   BugsnagSpec      = BugsnagSpec{}
	systemDefaultImage                         defaultImageSpec = defaultImageSpec{
		Name:       pointer.StringPtr("quay.io/3scale/porta"),
		Tag:        pointer.StringPtr("nightly"),
		PullPolicy: (*corev1.PullPolicy)(pointer.StringPtr(string(corev1.PullIfNotPresent))),
	}
	systemDefaultGrafanaDashboard defaultGrafanaDashboardSpec = defaultGrafanaDashboardSpec{
		SelectorKey:   pointer.StringPtr("monitoring-key"),
		SelectorValue: pointer.StringPtr("middleware"),
	}

	// App
	systemDefaultAppReplicas     int32                   = 2
	systemDefaultAppLoadBalancer defaultLoadBalancerSpec = defaultLoadBalancerSpec{
		ProxyProtocol:                 pointer.BoolPtr(true),
		CrossZoneLoadBalancingEnabled: pointer.BoolPtr(true),
		ConnectionDrainingEnabled:     pointer.BoolPtr(true),
		ConnectionDrainingTimeout:     pointer.Int32Ptr(60),
		HealthcheckHealthyThreshold:   pointer.Int32Ptr(2),
		HealthcheckUnhealthyThreshold: pointer.Int32Ptr(2),
		HealthcheckInterval:           pointer.Int32Ptr(5),
		HealthcheckTimeout:            pointer.Int32Ptr(3),
	}
	systemDefaultAppResources defaultResourceRequirementsSpec = defaultResourceRequirementsSpec{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("200m"),
			corev1.ResourceMemory: resource.MustParse("1Gi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("400m"),
			corev1.ResourceMemory: resource.MustParse("2Gi"),
		},
	}
	systemDefaultAppHPA defaultHorizontalPodAutoscalerSpec = defaultHorizontalPodAutoscalerSpec{
		MinReplicas:         pointer.Int32Ptr(2),
		MaxReplicas:         pointer.Int32Ptr(4),
		ResourceUtilization: pointer.Int32Ptr(90),
		ResourceName:        pointer.StringPtr("cpu"),
	}
	systemDefaultAppLivenessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: pointer.Int32Ptr(30),
		TimeoutSeconds:      pointer.Int32Ptr(1),
		PeriodSeconds:       pointer.Int32Ptr(10),
		SuccessThreshold:    pointer.Int32Ptr(1),
		FailureThreshold:    pointer.Int32Ptr(3),
	}
	systemDefaultAppReadinessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: pointer.Int32Ptr(30),
		TimeoutSeconds:      pointer.Int32Ptr(5),
		PeriodSeconds:       pointer.Int32Ptr(10),
		SuccessThreshold:    pointer.Int32Ptr(1),
		FailureThreshold:    pointer.Int32Ptr(3),
	}
	systemDefaultAppPDB defaultPodDisruptionBudgetSpec = defaultPodDisruptionBudgetSpec{
		MaxUnavailable: util.IntStrPtr(intstr.FromInt(1)),
	}
	systemDefaultAppMarin3rSpec defaultMarin3rSidecarSpec = defaultMarin3rSidecarSpec{}

	// Sidekiq
	systemDefaultSidekiqReplicas  int32                           = 2
	systemDefaultSidekiqResources defaultResourceRequirementsSpec = defaultResourceRequirementsSpec{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("500m"),
			corev1.ResourceMemory: resource.MustParse("1Gi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("1"),
			corev1.ResourceMemory: resource.MustParse("2Gi"),
		},
	}
	systemDefaultSidekiqHPA defaultHorizontalPodAutoscalerSpec = defaultHorizontalPodAutoscalerSpec{
		MinReplicas:         pointer.Int32Ptr(2),
		MaxReplicas:         pointer.Int32Ptr(4),
		ResourceUtilization: pointer.Int32Ptr(90),
		ResourceName:        pointer.StringPtr("cpu"),
	}
	systemDefaultSidekiqLivenessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: pointer.Int32Ptr(10),
		TimeoutSeconds:      pointer.Int32Ptr(3),
		PeriodSeconds:       pointer.Int32Ptr(15),
		SuccessThreshold:    pointer.Int32Ptr(1),
		FailureThreshold:    pointer.Int32Ptr(5),
	}
	systemDefaultSidekiqReadinessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: pointer.Int32Ptr(10),
		TimeoutSeconds:      pointer.Int32Ptr(5),
		PeriodSeconds:       pointer.Int32Ptr(30),
		SuccessThreshold:    pointer.Int32Ptr(1),
		FailureThreshold:    pointer.Int32Ptr(5),
	}
	systemDefaultSidekiqPDB defaultPodDisruptionBudgetSpec = defaultPodDisruptionBudgetSpec{
		MaxUnavailable: util.IntStrPtr(intstr.FromInt(1)),
	}

	systemDefaultSidekiqConfigDefault defaultSidekiqConfig = defaultSidekiqConfig{
		Queues: []string{
			"critical", "backend_sync", "events", "zync,40",
			"priority,25", "default,15", "web_hooks,10", "deletion,5",
		},
		MaxThreads: pointer.Int32Ptr(15),
	}
	systemDefaultSidekiqConfigBilling defaultSidekiqConfig = defaultSidekiqConfig{
		Queues:     []string{"billing"},
		MaxThreads: pointer.Int32Ptr(15),
	}
	systemDefaultSidekiqConfigLow defaultSidekiqConfig = defaultSidekiqConfig{
		Queues: []string{
			"mailers", "low", "bulk_indexing",
		},
		MaxThreads: pointer.Int32Ptr(15),
	}

	// Sphinx
	systemDefaultSphinxPort                int32                           = 9306
	systemDefaultSphinxBindAddress         string                          = "0.0.0.0"
	systemDefaultSphinxConfigFile          string                          = "/opt/system/db/sphinx/sphinx.conf"
	systemDefaultSphinxDBPath              string                          = "/opt/system/db/sphinx"
	systemDefaultSphinxDatabaseStorageSize string                          = "30Gi"
	systemDefaultSphinxPIDFile             string                          = "/opt/system/tmp/pids/searchd.pid"
	systemDefaultSphinxResources           defaultResourceRequirementsSpec = defaultResourceRequirementsSpec{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("250m"),
			corev1.ResourceMemory: resource.MustParse("4Gi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("750m"),
			corev1.ResourceMemory: resource.MustParse("5Gi"),
		},
	}
	systemDefaultSphinxLivenessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: pointer.Int32Ptr(60),
		TimeoutSeconds:      pointer.Int32Ptr(3),
		PeriodSeconds:       pointer.Int32Ptr(15),
		SuccessThreshold:    pointer.Int32Ptr(1),
		FailureThreshold:    pointer.Int32Ptr(5),
	}
	systemDefaultSphinxReadinessProbe defaultProbeSpec = defaultProbeSpec{
		InitialDelaySeconds: pointer.Int32Ptr(60),
		TimeoutSeconds:      pointer.Int32Ptr(5),
		PeriodSeconds:       pointer.Int32Ptr(30),
		SuccessThreshold:    pointer.Int32Ptr(1),
		FailureThreshold:    pointer.Int32Ptr(5),
	}
	systemDefaultRailsConsoleResources defaultResourceRequirementsSpec = defaultResourceRequirementsSpec{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("200m"),
			corev1.ResourceMemory: resource.MustParse("1Gi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("400m"),
			corev1.ResourceMemory: resource.MustParse("2Gi"),
		},
	}
)

// SystemSpec defines the desired state of System
type SystemSpec struct {
	// Application specific configuration options for System components
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Config SystemConfig `json:"config"`
	// Image specification for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Image *ImageSpec `json:"image,omitempty"`
	// Application specific configuration options
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	App *SystemAppSpec `json:"app,omitempty"`
	// Sidekiq Default specific configuration options
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	SidekiqDefault *SystemSidekiqSpec `json:"sidekiqDefault,omitempty"`
	// Sidekiq Billing specific configuration options
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	SidekiqBilling *SystemSidekiqSpec `json:"sidekiqBilling,omitempty"`
	// Sidekiq Low specific configuration options
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	SidekiqLow *SystemSidekiqSpec `json:"sidekiqLow,omitempty"`
	// Sphinx specific configuration options
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Sphinx *SystemSphinxSpec `json:"sphinx,omitempty"`
	// Console specific configuration options
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Console *SystemRailsConsoleSpec `json:"console,omitempty"`
	// Configures the Grafana Dashboard for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	GrafanaDashboard *GrafanaDashboardSpec `json:"grafanaDashboard,omitempty"`
	// Configures twemproxy
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Twemproxy *TwemproxySpec `json:"twemproxy,omitempty"`
}

// Default implements defaulting for SystemSpec
func (spec *SystemSpec) Default() {

	spec.Config.Default()
	spec.Image = InitializeImageSpec(spec.Image, systemDefaultImage)
	spec.GrafanaDashboard = InitializeGrafanaDashboardSpec(spec.GrafanaDashboard, systemDefaultGrafanaDashboard)
	if spec.App == nil {
		spec.App = &SystemAppSpec{}
	}
	spec.App.Default()

	if spec.SidekiqDefault == nil {
		spec.SidekiqDefault = &SystemSidekiqSpec{}
	}
	spec.SidekiqDefault.Default(Default)

	if spec.SidekiqBilling == nil {
		spec.SidekiqBilling = &SystemSidekiqSpec{}
	}
	spec.SidekiqBilling.Default(Billing)

	if spec.SidekiqLow == nil {
		spec.SidekiqLow = &SystemSidekiqSpec{}
	}
	spec.SidekiqLow.Default(Low)

	if spec.Sphinx == nil {
		spec.Sphinx = &SystemSphinxSpec{}
	}
	spec.Sphinx.Default(spec.Image)

	if spec.Console == nil {
		spec.Console = &SystemRailsConsoleSpec{}
	}
	spec.Console.Default()

	if spec.Twemproxy != nil {
		spec.Twemproxy.Default()
	}

}

// SystemConfig holds configuration for SystemApp component
type SystemConfig struct {
	// Rails configuration options for system components
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Rails *SystemRailsSpec `json:"rails,omitempty"`
	// OpenSSL verification mode for sandbox proxy
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	SandboxProxyOpensslVerifyMode *string `json:"sandboxProxyOpensslVerifyMode,omitempty"`
	// Enable (true) or disable (false) enforcing SSL
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ForceSSL *bool `json:"forceSSL,omitempty"`
	// SSL certificates path
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	SSLCertsDir *string `json:"sslCertsDir,omitempty"`
	// 3scale provider plan
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ThreescaleProviderPlan *string `json:"threescaleProviderPlan,omitempty"`
	// 3scale superdomain
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ThreescaleSuperdomain *string `json:"threescaleSuperdomain,omitempty"`
	// Secret containging system configuration files to be mounted in the pods
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ConfigFilesSecret *string `json:"configFilesSecret,omitempty"`
	// External Secret common configuration
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ExternalSecret ExternalSecret `json:"externalSecret,omitempty"`
	// DSN of system's main database
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	DatabaseDSN SecretReference `json:"databaseDSN"`
	// EventsSharedSecret
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	EventsSharedSecret SecretReference `json:"eventsSharedSecret"`
	// Holds recaptcha configuration options
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Recaptcha SystemRecaptchaSpec `json:"recaptcha"`
	// SecretKeyBase
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	SecretKeyBase SecretReference `json:"secretKeyBase"`
	// AccessCode to protect admin urls
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	AccessCode *SecretReference `json:"accessCode,omitempty"`
	// Options for Segment integration
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Segment SegmentSpec `json:"segment"`
	// Options for Github integration
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Github GithubSpec `json:"github"`
	// Options for configuring RH Customer Portal integration
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	RedHatCustomerPortal RedHatCustomerPortalSpec `json:"redhatCustomerPortal"`
	// Options for configuring Bugsnag integration
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Bugsnag *BugsnagSpec `json:"bugsnag,omitempty"`
	// Database secret
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	DatabaseSecret SecretReference `json:"databaseSecret"`
	// Memcached servers
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	MemcachedServers string `json:"memcachedServers"`
	// Redis configuration options
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Redis RedisSpec `json:"redis"`
	// SMTP configuration options
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	SMTP SMTPSpec `json:"smtp"`
	// Mapping Service access token
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	MappingServiceAccessToken SecretReference `json:"mappingServiceAccessToken"`
	// Zync authentication token
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ZyncAuthToken SecretReference `json:"zyncAuthToken"`
	// Backend has configuration options for system to contact backend
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Backend SystemBackendSpec `json:"backend"`
	// Assets has configuration to access assets in AWS s3
	Assets AssetsSpec `json:"assets"`
}

// Default applies default values to a SystemConfig struct
func (sc *SystemConfig) Default() {
	if sc.Rails == nil {
		sc.Rails = &SystemRailsSpec{}
	}
	sc.Rails.Default()

	sc.ConfigFilesSecret = stringOrDefault(sc.ConfigFilesSecret, pointer.StringPtr(systemDefaultConfigFilesSecret))

	if sc.Bugsnag == nil {
		sc.Bugsnag = &systemDefaultBugsnagSpec
	}

	sc.SandboxProxyOpensslVerifyMode = stringOrDefault(sc.SandboxProxyOpensslVerifyMode, pointer.StringPtr(systemDefaultSandboxProxyOpensslVerifyMode))
	sc.ForceSSL = boolOrDefault(sc.ForceSSL, pointer.BoolPtr(systemDefaultForceSSL))
	sc.SSLCertsDir = stringOrDefault(sc.SSLCertsDir, pointer.StringPtr(systemDefaultSSLCertsDir))
	sc.ThreescaleProviderPlan = stringOrDefault(sc.ThreescaleProviderPlan, pointer.StringPtr(systemDefaultThreescaleProviderPlan))
	sc.ThreescaleSuperdomain = stringOrDefault(sc.ThreescaleSuperdomain, pointer.StringPtr(systemDefaultThreescaleSuperdomain))
	sc.ExternalSecret.SecretStoreRef = InitializeExternalSecretSecretStoreReferenceSpec(sc.ExternalSecret.SecretStoreRef, defaultExternalSecretSecretStoreReference)
	sc.ExternalSecret.RefreshInterval = durationOrDefault(sc.ExternalSecret.RefreshInterval, &defaultExternalSecretRefreshInterval)
}

// ResolveCanarySpec modifies the SystemSpec given the provided canary configuration
func (spec *SystemSpec) ResolveCanarySpec(canary *Canary) (*SystemSpec, error) {
	canarySpec := &SystemSpec{}
	if err := canary.PatchSpec(spec, canarySpec); err != nil {
		return nil, err
	}

	if canary.ImageName != nil {
		canarySpec.Image.Name = canary.ImageName
	}
	if canary.ImageTag != nil {
		canarySpec.Image.Tag = canary.ImageTag
	}
	canarySpec.App.Replicas = canary.Replicas
	canarySpec.SidekiqDefault.Replicas = canary.Replicas
	canarySpec.SidekiqLow.Replicas = canary.Replicas
	canarySpec.SidekiqBilling.Replicas = canary.Replicas

	// Call Default() on the resolved canary spec to apply
	// defaulting to potentially added fields
	canarySpec.Default()

	return canarySpec, nil
}

// SystemRecaptchaSpec holds recaptcha configurations
type SystemRecaptchaSpec struct {
	// Public key
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	PublicKey SecretReference `json:"publicKey"`
	// Private key
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	PrivateKey SecretReference `json:"privateKey"`
}

// SegmentSpec has configuration for Segment integration
type SegmentSpec struct {
	// Deletion workspace
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	DeletionWorkspace string `json:"deletionWorkspace"`
	// Deletion token
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	DeletionToken SecretReference `json:"deletionToken"`
	// Write key
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	WriteKey SecretReference `json:"writeKey"`
}

// GithubSpec has configuration for Github integration
type GithubSpec struct {
	// Client ID
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ClientID SecretReference `json:"clientID"`
	// Client secret
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ClientSecret SecretReference `json:"clientSecret"`
}

// RedHatCustomerPortalSpec has configuration for integration with
// Red Hat Customer Portal
type RedHatCustomerPortalSpec struct {
	// Client ID
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ClientID SecretReference `json:"clientID"`
	// Client secret
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ClientSecret SecretReference `json:"clientSecret"`
	// Realm
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Realm *SecretReference `json:"realm,omitempty"`
}

// RedisSpec holds redis configuration
type RedisSpec struct {
	// Data source name
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	QueuesDSN string `json:"queuesDSN"`
}

// SMTPSpec has options to configure system's SMTP
type SMTPSpec struct {
	// Address
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Address string `json:"address"`
	// User
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	User SecretReference `json:"user"`
	// Password
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Password SecretReference `json:"password"`
	// Port
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Port int32 `json:"port"`
	// Authentication protocol
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	AuthProtocol string `json:"authProtocol"`
	// OpenSSL verify mode
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	OpenSSLVerifyMode string `json:"opensslVerifyMode"`
	// Enable/disable auto STARTTLS
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	STARTTLSAuto bool `json:"starttlsAuto"`
}

// SystemBackendSpec has configuration options for backend
type SystemBackendSpec struct {
	// External endpoint
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ExternalEndpoint string `json:"externalEndpoint"`
	// Internal endpoint
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	InternalEndpoint string `json:"internalEndpoint"`
	// Internal API user
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	InternalAPIUser SecretReference `json:"internalAPIUser"`
	// Internal API password
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	InternalAPIPassword SecretReference `json:"internalAPIPassword"`
	// Redis data source name
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	RedisDSN string `json:"redisDSN"`
}

// AssetsSpec has configuration to access assets in AWS s3
type AssetsSpec struct {
	// AWS S3 bucket name
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Bucket string `json:"bucket"`
	// AWS S3 region
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Region string `json:"region"`
	// AWS access key
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	AccessKey SecretReference `json:"accessKey"`
	// AWS secret access key
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	SecretKey SecretReference `json:"secretKey"`
	// Assets host (CDN)
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Host *string `json:"host,omitempty"`
}

// SystemRailsSpec configures rails for system components
type SystemRailsSpec struct {
	// Rails Console
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Console *bool `json:"console,omitempty"`
	// Rails environment
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Environment *string `json:"environment,omitempty"`
	// Rails log level (debug, info, warn, error, fatal or unknown)
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:validation:Enum=debug;info;warn;error;fatal;unknown
	// +optional
	LogLevel *string `json:"logLevel,omitempty"`
}

// Default applies defaults for SystemRailsSpec
func (srs *SystemRailsSpec) Default() {
	srs.Console = boolOrDefault(srs.Console, pointer.BoolPtr(systemDefaultRailsConsole))
	srs.Environment = stringOrDefault(srs.Environment, pointer.StringPtr(systemDefaultRailsEnvironment))
	srs.LogLevel = stringOrDefault(srs.LogLevel, pointer.StringPtr(systemDefaultRailsLogLevel))
}

// SystemAppSpec configures the App component of System
type SystemAppSpec struct {
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
	// Marin3r configures the Marin3r sidecars for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Marin3r *Marin3rSidecarSpec `json:"marin3r,omitempty"`
	// Describes node affinity scheduling rules for the pod.
	// +optional
	NodeAffinity *corev1.NodeAffinity `json:"nodeAffinity,omitempty" protobuf:"bytes,1,opt,name=nodeAffinity"`
	// If specified, the pod's tolerations.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty" protobuf:"bytes,22,opt,name=tolerations"`
	// Canary defines spec changes for the canary Deployment. If
	// left unset the canary Deployment wil not be created.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Canary *Canary `json:"canary,omitempty"`
}

// Default implements defaulting for the system App component
func (spec *SystemAppSpec) Default() {
	spec.HPA = InitializeHorizontalPodAutoscalerSpec(spec.HPA, systemDefaultAppHPA)
	spec.Replicas = intOrDefault(spec.Replicas, &systemDefaultAppReplicas)
	spec.PDB = InitializePodDisruptionBudgetSpec(spec.PDB, systemDefaultAppPDB)
	spec.Resources = InitializeResourceRequirementsSpec(spec.Resources, systemDefaultAppResources)
	spec.LivenessProbe = InitializeProbeSpec(spec.LivenessProbe, systemDefaultAppLivenessProbe)
	spec.ReadinessProbe = InitializeProbeSpec(spec.ReadinessProbe, systemDefaultAppReadinessProbe)
	spec.Marin3r = InitializeMarin3rSidecarSpec(spec.Marin3r, systemDefaultAppMarin3rSpec)
}

// SystemSidekiqSpec configures the Sidekiq component of System
type SystemSidekiqSpec struct {
	// Sidekiq specific configuration options for the component element
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Config *SidekiqConfig `json:"config,omitempty"`
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
	// Marin3r configures the Marin3r sidecars for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Marin3r *Marin3rSidecarSpec `json:"marin3r,omitempty"`
	// Describes node affinity scheduling rules for the pod.
	// +optional
	NodeAffinity *corev1.NodeAffinity `json:"nodeAffinity,omitempty" protobuf:"bytes,1,opt,name=nodeAffinity"`
	// If specified, the pod's tolerations.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty" protobuf:"bytes,22,opt,name=tolerations"`
	// Canary defines spec changes for the canary Deployment. If
	// left unset the canary Deployment wil not be created.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Canary *Canary `json:"canary,omitempty"`
}

// SidekiqConfig configures app behavior for System Sidekiq
type SidekiqConfig struct {
	// List of queues to be consumed by sidekiq. Format: queue[,Priority]
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Queues []string `json:"queues,omitempty"`
	// Number of rails max threads per sidekiq pod
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	MaxThreads *int32 `json:"maxThreads,omitempty"`
}

type defaultSidekiqConfig struct {
	Queues     []string
	MaxThreads *int32
}

// Default sets default values for any value not specifically set in the SidekiqConfig struct
func (cfg *SidekiqConfig) Default(def defaultSidekiqConfig) {
	if cfg.Queues == nil {
		cfg.Queues = def.Queues
	}
	cfg.MaxThreads = intOrDefault(cfg.MaxThreads, pointer.Int32Ptr(*def.MaxThreads))
}

// Default implements defaulting for the system Sidekiq component
func (spec *SystemSidekiqSpec) Default(sidekiqType systemSidekiqType) {
	spec.HPA = InitializeHorizontalPodAutoscalerSpec(spec.HPA, systemDefaultSidekiqHPA)
	spec.Replicas = intOrDefault(spec.Replicas, &systemDefaultSidekiqReplicas)
	spec.PDB = InitializePodDisruptionBudgetSpec(spec.PDB, systemDefaultSidekiqPDB)
	spec.Resources = InitializeResourceRequirementsSpec(spec.Resources, systemDefaultSidekiqResources)
	spec.LivenessProbe = InitializeProbeSpec(spec.LivenessProbe, systemDefaultSidekiqLivenessProbe)
	spec.ReadinessProbe = InitializeProbeSpec(spec.ReadinessProbe, systemDefaultSidekiqReadinessProbe)
	spec.Marin3r = InitializeMarin3rSidecarSpec(spec.Marin3r, systemDefaultAppMarin3rSpec)
	if spec.Config == nil {
		spec.Config = &SidekiqConfig{}
	}

	if sidekiqType == Billing {
		spec.Config.Default(systemDefaultSidekiqConfigBilling)
	} else if sidekiqType == Low {
		spec.Config.Default(systemDefaultSidekiqConfigLow)
	} else {
		spec.Config.Default(systemDefaultSidekiqConfigDefault)
	}
}

// SystemSphinxSpec configures the App component of System
type SystemSphinxSpec struct {
	// Image specification for the Sphinx component.
	// Defaults to system image if not defined.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Image *ImageSpec `json:"image,omitempty"`
	// Configuration options for System's sphinx
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Config *SphinxConfig `json:"config,omitempty"`
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
	// +optional
	NodeAffinity *corev1.NodeAffinity `json:"nodeAffinity,omitempty" protobuf:"bytes,1,opt,name=nodeAffinity"`
	// If specified, the pod's tolerations.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty" protobuf:"bytes,22,opt,name=tolerations"`
}

// Default implements defaulting for the system sphinx component
func (spec *SystemSphinxSpec) Default(systemDefaultImage *ImageSpec) {

	spec.Image = InitializeImageSpec(spec.Image, defaultImageSpec(*systemDefaultImage))
	spec.Resources = InitializeResourceRequirementsSpec(spec.Resources, systemDefaultSphinxResources)
	spec.LivenessProbe = InitializeProbeSpec(spec.LivenessProbe, systemDefaultSphinxLivenessProbe)
	spec.ReadinessProbe = InitializeProbeSpec(spec.ReadinessProbe, systemDefaultSphinxReadinessProbe)
	if spec.Config == nil {
		spec.Config = &SphinxConfig{}
	}
	spec.Config.Default()
}

// SphinxConfig has configuration options for System's sphinx
type SphinxConfig struct {
	// Thinking configuration for sphinx
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Thinking *ThinkingSpec `json:"thinking,omitempty"`
}

// Default implements defaulting for SphinxConfig
func (sc *SphinxConfig) Default() {
	if sc.Thinking == nil {
		sc.Thinking = &ThinkingSpec{}
	}
	sc.Thinking.Default()
}

// ThinkingSpec configures the thinking library for sphinx
type ThinkingSpec struct {
	// The TCP port Sphinx will run its daemon on
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Port *int32 `json:"port,omitempty"`
	// Allows setting the TCP host for Sphinx to a different address
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	BindAddress *string `json:"bindAddress,omitempty"`
	// Sphinx configuration file path
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ConfigFile *string `json:"configFile,omitempty"`
	// Sphinx database path
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	DatabasePath *string `json:"databasePath,omitempty"`
	// Sphinx database storage size
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	DatabaseStorageSize *resource.Quantity `json:"databaseStorageSize,omitempty"`
	// Sphinx database storage type
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	DatabaseStorageClass *string `json:"databaseStorageClass,omitempty"`
	// Sphinx PID file path
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	PIDFile *string `json:"pidFile,omitempty"`
}

// Default implements defaulting for ThinkingSpec
func (tc *ThinkingSpec) Default() {
	tc.Port = intOrDefault(tc.Port, pointer.Int32Ptr(systemDefaultSphinxPort))
	tc.BindAddress = stringOrDefault(tc.BindAddress, pointer.StringPtr(systemDefaultSphinxBindAddress))
	tc.ConfigFile = stringOrDefault(tc.ConfigFile, pointer.StringPtr(systemDefaultSphinxConfigFile))
	tc.DatabasePath = stringOrDefault(tc.DatabasePath, pointer.StringPtr(systemDefaultSphinxDBPath))
	tc.PIDFile = stringOrDefault(tc.PIDFile, pointer.StringPtr(systemDefaultSphinxPIDFile))
	if tc.DatabaseStorageSize == nil {
		size := resource.MustParse(systemDefaultSphinxDatabaseStorageSize)
		tc.DatabaseStorageSize = &size
	}
}

// SystemRailsConsoleSpec configures the App component of System
type SystemRailsConsoleSpec struct {
	// Resource requirements for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Resources *ResourceRequirementsSpec `json:"resources,omitempty"`
	// Describes node affinity scheduling rules for the pod.
	// +optional
	NodeAffinity *corev1.NodeAffinity `json:"nodeAffinity,omitempty" protobuf:"bytes,1,opt,name=nodeAffinity"`
	// If specified, the pod's tolerations.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty" protobuf:"bytes,22,opt,name=tolerations"`
}

// Default implements defaulting for the system App component
func (spec *SystemRailsConsoleSpec) Default() {
	spec.Resources = InitializeResourceRequirementsSpec(spec.Resources, systemDefaultRailsConsoleResources)
}

// SystemStatus defines the observed state of System
type SystemStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// System is the Schema for the systems API
type System struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SystemSpec   `json:"spec,omitempty"`
	Status SystemStatus `json:"status,omitempty"`
}

// Default implements defaulting for the System resource
func (s *System) Default() {
	s.Spec.Default()
}

// +kubebuilder:object:root=true

// SystemList contains a list of System
type SystemList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []System `json:"items"`
}

// GetItem returns a client.Objectfrom a SystemList
func (sl *SystemList) GetItem(idx int) client.Object {
	return &sl.Items[idx]
}

// CountItems returns the item count in SystemList.Items
func (sl *SystemList) CountItems() int {
	return len(sl.Items)
}

func init() {
	SchemeBuilder.Register(&System{}, &SystemList{})
}
