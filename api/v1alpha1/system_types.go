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
	"reflect"

	"github.com/3scale/saas-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	// Common
	systemDefaultAMPRelease                    string           = "2.7.1"
	systemDefaultSandboxProxyOpensslVerifyMode string           = "VERIFY_NONE"
	systemDefaultForceSSL                      bool             = true
	systemDefaultSSLCertsDir                   string           = "/etc/pki/tls/certs"
	systemDefaultThreescaleProviderPlan        string           = "enterprise"
	systemDefaultThreescaleSuperdomain         string           = "localhost"
	systemDefaultRailsEnvironment              string           = "preview"
	systemDefaultRailsLogLevel                 string           = "info"
	systemDefaultLogToStdout                   bool             = true
	systemDefaultConfigFiles                   ConfigFilesSpec  = ConfigFilesSpec{}
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

	// Sphinx
	systemDefaultSphinxDeltaIndexInterval  int32                           = 5
	systemDefaultSphinxFullReindexInterval int32                           = 60
	systemDefaultSphinxPort                int32                           = 9306
	systemDefaultSphinxBindAddress         string                          = "0.0.0.0"
	systemDefaultSphinxConfigFile          string                          = "/opt/system/db/sphinx/preview.conf"
	systemDefaultSphinxDBPath              string                          = "/opt/system/db/sphinx"
	systemDefaultSphinxPIDFile             string                          = "/opt/system/tmp/pids/searchd.pid"
	systemDefaultSphinxStorage             string                          = "30Gi"
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
	// Sidekiq specific configuration options
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Sidekiq *SystemSidekiqSpec `json:"sidekiq,omitempty"`
	// Sphinx specific configuration options
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	Sphinx *SystemSphinxSpec `json:"sphinx,omitempty"`
	// Configures the Grafana Dashboard for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	GrafanaDashboard *GrafanaDashboardSpec `json:"grafanaDashboard,omitempty"`
}

// Default implements defaulting for the System resource
func (s *System) Default() {

	s.Spec.Config.Default()
	s.Spec.Image = InitializeImageSpec(s.Spec.Image, systemDefaultImage)
	s.Spec.GrafanaDashboard = InitializeGrafanaDashboardSpec(s.Spec.GrafanaDashboard, systemDefaultGrafanaDashboard)
	if s.Spec.App == nil {
		s.Spec.App = &SystemAppSpec{}
	}
	s.Spec.App.Default()

	if s.Spec.Sidekiq == nil {
		s.Spec.Sidekiq = &SystemSidekiqSpec{}
	}
	s.Spec.Sidekiq.Default()

	if s.Spec.Sphinx == nil {
		s.Spec.Sphinx = &SystemSphinxSpec{}
	}
	s.Spec.Sphinx.Default()
}

// SystemConfig holds configuration for SystemApp component
type SystemConfig struct {
	// AMP release number
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	AMPRelease *string `json:"ampRelease,omitempty"`
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
	// Extra configuration files to be mounted in the pods
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	ConfigFiles *ConfigFilesSpec `json:"configFiles,omitempty"`
	// System seed
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Seed SystemSeedSpec `json:"seed"`
	// DSN of system's main database
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	DatabaseDSN SecretReference `json:"databaseDSN"`
	// ???
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	EventsSharedSecret SecretReference `json:"eventsSharedSecret"`
	// Holds recaptcha configuration options
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Recaptcha SystemRecaptchaSpec `json:"recaptcha"`
	// ???
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	SecretKeyBase SecretReference `json:"secretKeyBase"`
	// ???
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	AccessCode SecretReference `json:"accessCode"`
	// Options for Segment integration
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Segment SegmentSpec `json:"segment"`
	// Options for Github integration
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Github GithubSpec `json:"github"`
	// Options for configuring metrics publication
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Metrics MetricsSpec `json:"metrics"`
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
	// Master access token for Apicast
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ApicastAccessToken SecretReference `json:"apicastAccessToken"`
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
	sc.AMPRelease = stringOrDefault(sc.AMPRelease, pointer.StringPtr(systemDefaultAMPRelease))

	if sc.Rails == nil {
		sc.Rails = &SystemRailsSpec{}
	}
	sc.Rails.Default()

	if sc.ConfigFiles == nil {
		sc.ConfigFiles = &systemDefaultConfigFiles
	}

	if sc.Bugsnag == nil {
		sc.Bugsnag = &systemDefaultBugsnagSpec
	}

	sc.SandboxProxyOpensslVerifyMode = stringOrDefault(sc.SandboxProxyOpensslVerifyMode, pointer.StringPtr(systemDefaultSandboxProxyOpensslVerifyMode))
	sc.ForceSSL = boolOrDefault(sc.ForceSSL, pointer.BoolPtr(systemDefaultForceSSL))
	sc.SSLCertsDir = stringOrDefault(sc.SSLCertsDir, pointer.StringPtr(systemDefaultSSLCertsDir))
	sc.ThreescaleProviderPlan = stringOrDefault(sc.ThreescaleProviderPlan, pointer.StringPtr(systemDefaultThreescaleProviderPlan))
	sc.ThreescaleSuperdomain = stringOrDefault(sc.ThreescaleSuperdomain, pointer.StringPtr(systemDefaultThreescaleSuperdomain))
}

// ConfigFilesSpec defines a vault location to
// get system config files from
type ConfigFilesSpec struct {
	VaultPath string   `json:"vaultPath"`
	Files     []string `json:"files"`
}

// Enabled returns a boolean indication whether the
// ConfigFiles config options is enabled or not
func (cfs *ConfigFilesSpec) Enabled() bool {
	if reflect.DeepEqual(cfs, &ConfigFilesSpec{}) {
		return false
	}
	return true
}

// SystemSeedSpec whatever this is
type SystemSeedSpec struct {
	MasterAccessToken SecretReference `json:"masterAccessToken"`
	MasterDomain      string          `json:"masterDomain"`
	MasterUser        SecretReference `json:"masterUser"`
	MasterPassword    SecretReference `json:"masterPassword"`
	AdminAccessToken  SecretReference `json:"adminAccessToken"`
	AdminUser         SecretReference `json:"adminUser"`
	AdminPassword     SecretReference `json:"adminPassword"`
	AdminEmail        string          `json:"adminEmail"`
	TenantName        string          `json:"tenantName"`
}

// SystemRecaptchaSpec holds recaptcha configurations
type SystemRecaptchaSpec struct {
	PublicKey  SecretReference `json:"publicKey"`
	PrivateKey SecretReference `json:"privateKey"`
}

// SegmentSpec has configuration for Segment integration
type SegmentSpec struct {
	DeletionWorkspace string          `json:"deletionWorkspace"`
	DeletionToken     SecretReference `json:"deletionToken"`
	WriteKey          SecretReference `json:"writeKey"`
}

// NewRelicSpec has configuration for NewRelic integration
type NewRelicSpec struct {
	LicenseKey SecretReference `json:"licenseKey"`
}

// GithubSpec has configuration for Github integration
type GithubSpec struct {
	ClientID     SecretReference `json:"clientID"`
	ClientSecret SecretReference `json:"clientSecret"`
}

// MetricsSpec has options to configure prometheus metrics
type MetricsSpec struct {
	User     SecretReference `json:"user"`
	Password SecretReference `json:"password"`
}

// RedHatCustomerPortalSpec has configuration for integration with
// Red Hat Customer Portal
type RedHatCustomerPortalSpec struct {
	ClientID     SecretReference `json:"clientID"`
	ClientSecret SecretReference `json:"clientSecret"`
}

// BugsnagSpec has configuration for Bugsnag integration
type BugsnagSpec struct {
	APIKey SecretReference `json:"apiKey"`
}

// Enabled returns a boolean indication whether the
// Bugsnag integration is enabled or not
func (bs *BugsnagSpec) Enabled() bool {
	if reflect.DeepEqual(bs, &BugsnagSpec{}) {
		return false
	}
	return true
}

// RedisSpec holds redis configuration
type RedisSpec struct {
	DSN           string `json:"dsn"`
	MessageBusDSN string `json:"messageBusDSN"`
}

// SMTPSpec has options to configure system's SMTP
type SMTPSpec struct {
	Address           string          `json:"address"`
	User              SecretReference `json:"user"`
	Password          SecretReference `json:"password"`
	Port              int32           `json:"port"`
	AuthProtocol      string          `json:"authProtocol"`
	OpenSSLVerifyMode string          `json:"opensslVerifyMode"`
	STARTTLSAuto      bool            `json:"starttlsAuto"`
}

// SystemBackendSpec has configuration options for backend
type SystemBackendSpec struct {
	ExternalEndpoint    string          `json:"externalEndpoint"`
	InternalEndpoint    string          `json:"internalEndpoint"`
	InternalAPIUser     SecretReference `json:"internalAPIUser"`
	InternalAPIPassword SecretReference `json:"internalAPIPassword"`
	RedisDSN            string          `json:"redisDSN"`
}

// AssetsSpec has configuration to access assets in AWS s3
type AssetsSpec struct {
	Bucket    string          `json:"bucket"`
	Region    string          `json:"region"`
	AccessKey SecretReference `json:"accessKey"`
	SecretKey SecretReference `json:"secretKey"`
}

// SystemRailsSpec configures rails for system components
type SystemRailsSpec struct {
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
	srs.Environment = pointer.StringPtr(systemDefaultRailsEnvironment)
	srs.LogLevel = pointer.StringPtr(systemDefaultRailsLogLevel)
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
	// Configures the AWS load balancer for the component
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	LoadBalancer *LoadBalancerSpec `json:"loadBalancer,omitempty"`
	// System's app endpoint
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Endpoint Endpoint `json:"endpoint"`
}

// Default implements defaulting for the system App component
func (spec *SystemAppSpec) Default() {
	spec.HPA = InitializeHorizontalPodAutoscalerSpec(spec.HPA, systemDefaultAppHPA)

	if spec.HPA.IsDeactivated() {
		spec.Replicas = intOrDefault(spec.Replicas, &systemDefaultAppReplicas)
	} else {
		spec.Replicas = nil
	}

	spec.PDB = InitializePodDisruptionBudgetSpec(spec.PDB, systemDefaultAppPDB)
	spec.Resources = InitializeResourceRequirementsSpec(spec.Resources, systemDefaultAppResources)
	spec.LivenessProbe = InitializeProbeSpec(spec.LivenessProbe, systemDefaultAppLivenessProbe)
	spec.ReadinessProbe = InitializeProbeSpec(spec.ReadinessProbe, systemDefaultAppReadinessProbe)
	spec.LoadBalancer = InitializeLoadBalancerSpec(spec.LoadBalancer, systemDefaultAppLoadBalancer)
	spec.Marin3r = InitializeMarin3rSidecarSpec(spec.Marin3r, systemDefaultAppMarin3rSpec)
}

// SystemSidekiqSpec configures the App component of System
type SystemSidekiqSpec struct {
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
}

// Default implements defaulting for the system App component
func (spec *SystemSidekiqSpec) Default() {
	spec.HPA = InitializeHorizontalPodAutoscalerSpec(spec.HPA, systemDefaultSidekiqHPA)

	if spec.HPA.IsDeactivated() {
		spec.Replicas = intOrDefault(spec.Replicas, &systemDefaultSidekiqReplicas)
	} else {
		spec.Replicas = nil
	}

	spec.PDB = InitializePodDisruptionBudgetSpec(spec.PDB, systemDefaultSidekiqPDB)
	spec.Resources = InitializeResourceRequirementsSpec(spec.Resources, systemDefaultSidekiqResources)
	spec.LivenessProbe = InitializeProbeSpec(spec.LivenessProbe, systemDefaultSidekiqLivenessProbe)
	spec.ReadinessProbe = InitializeProbeSpec(spec.ReadinessProbe, systemDefaultSidekiqReadinessProbe)
}

// SystemSphinxSpec configures the App component of System
type SystemSphinxSpec struct {
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
}

// Default implements defaulting for the system sphinx component
func (spec *SystemSphinxSpec) Default() {

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
	// Interval used for adding chunks of brand new documents to the primary
	// index at certain intervals without having to do a full re-index
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	DeltaIndexInterval *int32 `json:"deltaIndexInterval,omitempty"`
	// Interval used to do a full re-index
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +optional
	FullReindexInterval *int32 `json:"fullReindexInterval,omitempty"`
}

// Default implements defaulting for SphinxConfig
func (sc *SphinxConfig) Default() {
	if sc.Thinking == nil {
		sc.Thinking = &ThinkingSpec{}
	}
	sc.Thinking.Default()
	sc.DeltaIndexInterval = intOrDefault(sc.DeltaIndexInterval, pointer.Int32Ptr(systemDefaultSphinxDeltaIndexInterval))
	sc.FullReindexInterval = intOrDefault(sc.FullReindexInterval, pointer.Int32Ptr(systemDefaultSphinxFullReindexInterval))
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
	DBPath *string `json:"dbPath,omitempty"`
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
	tc.DBPath = stringOrDefault(tc.DBPath, pointer.StringPtr(systemDefaultSphinxDBPath))
	tc.PIDFile = stringOrDefault(tc.PIDFile, pointer.StringPtr(systemDefaultSphinxPIDFile))
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
