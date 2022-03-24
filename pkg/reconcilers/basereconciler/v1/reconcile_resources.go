package basereconciler

import (
	"context"
	"fmt"
	"hash/fnv"
	"reflect"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	externalsecretsv1alpha1 "github.com/3scale/saas-operator/pkg/apis/externalsecrets/v1alpha1"
	"github.com/davecgh/go-spew/spew"
	"github.com/redhat-cop/operator-utils/pkg/util/lockedresourcecontroller/lockedpatch"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ControlledResources defines the resources that each of the
// Custom Resources own
type ControlledResources struct {
	Deployments              []Deployment
	StatefulSets             []StatefulSet
	ExternalSecrets          []ExternalSecret
	Services                 []Service
	PodDisruptionBudgets     []PodDisruptionBudget
	HorizontalPodAutoscalers []HorizontalPodAutoscaler
	PodMonitors              []PodMonitor
	GrafanaDashboards        []GrafanaDashboard
	ConfigMaps               []ConfigMaps
}

func (cm *ControlledResources) Add(resources *ControlledResources) *ControlledResources {
	cm.Deployments = append(cm.Deployments, resources.Deployments...)
	cm.StatefulSets = append(cm.StatefulSets, resources.StatefulSets...)
	cm.Services = append(cm.Services, resources.Services...)
	cm.PodDisruptionBudgets = append(cm.PodDisruptionBudgets, resources.PodDisruptionBudgets...)
	cm.HorizontalPodAutoscalers = append(cm.HorizontalPodAutoscalers, resources.HorizontalPodAutoscalers...)
	cm.PodMonitors = append(cm.PodMonitors, resources.PodMonitors...)
	cm.GrafanaDashboards = append(cm.GrafanaDashboards, resources.GrafanaDashboards...)
	cm.ConfigMaps = append(cm.ConfigMaps, resources.ConfigMaps...)
	return cm
}

// RolloutTrigger defines a configuration source that should trigger a
// rollout whenever the data within that configuration source changes
type RolloutTrigger struct {
	name      string
	configMap *corev1.ConfigMap
	secret    *corev1.Secret
}

// GetHash returns the hash of the data container in the RolloutTrigger
// config source
func (rt *RolloutTrigger) GetHash() string {
	if rt.secret != nil {
		if reflect.DeepEqual(rt.secret, &corev1.Secret{}) {
			return ""
		}
		return Hash(rt.secret.Data)
	}
	if rt.configMap != nil {
		if reflect.DeepEqual(rt.secret, &corev1.ConfigMap{}) {
			return ""
		}
		return Hash(rt.configMap.Data)
	}
	return ""
}

// GetAnnotationKey returns the annotation key to be used in the Pods that read
// from the config source defined in the RolloutTrigger
func (rt *RolloutTrigger) GetAnnotationKey() string {
	if rt.secret != nil {
		return fmt.Sprintf("%s/%s.%s", saasv1alpha1.AnnotationsDomain, rt.name, "secret-hash")
	}
	return fmt.Sprintf("%s/%s.%s", saasv1alpha1.AnnotationsDomain, rt.name, "configmap-hash")
}

// NewRolloutTrigger returns a new RolloutTrigger from a Secret or ConfigMap
// It panics if the passed client.Object is not a Secret or ConfigMap
func NewRolloutTrigger(name string, o client.Object) RolloutTrigger {
	switch trigger := o.(type) {
	case *corev1.Secret:
		return RolloutTrigger{name: name, secret: trigger}
	case *corev1.ConfigMap:
		return RolloutTrigger{name: name, configMap: trigger}
	default:
		panic("unsupported rollout trigger")
	}
}

// TriggersFromExternalSecrets generates a list of RolloutTrigger from the given ExternalSecret generator functions
func (r *Reconciler) TriggersFromExternalSecrets(ctx context.Context, es ...GeneratorFunction) ([]RolloutTrigger, error) {

	triggers := []RolloutTrigger{}

	for _, externalSecret := range es {
		es := externalSecret().(*externalsecretsv1alpha1.ExternalSecret)
		strgs, err := r.TriggersFromSecret(ctx, es.GetNamespace(), es.GetName())
		if err != nil {
			return nil, err
		}
		triggers = append(triggers, strgs...)
	}

	return triggers, nil
}

// TriggersFromSecret generates a list of RolloutTriggers from the given Secrets name list
func (r *Reconciler) TriggersFromSecret(ctx context.Context, namespace string, secrets ...string) ([]RolloutTrigger, error) {

	strgs := []RolloutTrigger{}

	for _, secretName := range secrets {

		key := types.NamespacedName{
			Name:      secretName,
			Namespace: namespace,
		}
		secret := &corev1.Secret{}
		err := r.GetClient().Get(ctx, key, secret)

		if err != nil {
			if errors.IsNotFound(err) {
				strgs = append(strgs, NewRolloutTrigger(key.Name, &corev1.Secret{}))
				continue
			}
			return nil, err
		}

		strgs = append(strgs, NewRolloutTrigger(key.Name, secret))
	}

	return strgs, nil
}

// Deployment specifies a Deployment resource and its rollout triggers
type Deployment struct {
	Template        GeneratorFunction
	RolloutTriggers []RolloutTrigger
	HasHPA          bool
}

// StatefulSet specifies a StatefulSet resource and its rollout triggers
type StatefulSet struct {
	Template        GeneratorFunction
	RolloutTriggers []RolloutTrigger
	Enabled         bool
}

// ExternalSecret specifies a ExternalSecret resource
type ExternalSecret struct {
	Template GeneratorFunction
	Enabled  bool
}

// Service specifies a Service resource
type Service struct {
	Template GeneratorFunction
	Enabled  bool
}

func (s *Service) PopulateSpecRuntimeValues(ctx context.Context, cl client.Client) (GeneratorFunction, error) {

	svc := s.Template().(*corev1.Service)
	instance := &corev1.Service{}
	if err := cl.Get(ctx, types.NamespacedName{
		Name:      svc.GetName(),
		Namespace: svc.GetNamespace(),
	}, instance); err != nil {
		if errors.IsNotFound(err) {
			// Resource not found, return the template as is
			// because there are not runtime values yet
			return s.Template, nil
		}
		return s.Template, err
	}

	// Set runtime values in the resource:
	// "/spec/clusterIP", "/spec/clusterIPs", "/spec/ipFamilies", "/spec/ipFamilyPolicy", "/spec/ports/*/nodePort"
	svc.Spec.ClusterIP = instance.Spec.ClusterIP
	svc.Spec.ClusterIPs = instance.Spec.ClusterIPs
	svc.Spec.IPFamilies = instance.Spec.IPFamilies
	svc.Spec.IPFamilyPolicy = instance.Spec.IPFamilyPolicy

	// For services that are not ClusterIP we need to populate the runtime values
	// of NodePort for each port
	if svc.Spec.Type != corev1.ServiceTypeClusterIP {
		for idx, port := range svc.Spec.Ports {
			runtimePort := findPort(port.Port, port.Protocol, instance.Spec.Ports)
			if runtimePort != nil {
				svc.Spec.Ports[idx].NodePort = runtimePort.NodePort
			}
		}
	}

	return func() client.Object {
		return svc
	}, nil
}

func findPort(pNumber int32, pProtocol corev1.Protocol, ports []corev1.ServicePort) *corev1.ServicePort {
	// Ports within a svc are uniquely identified by
	// the "port" and "protocol" fields. This is documented in
	// k8s API reference
	for _, port := range ports {
		if pNumber == port.Port && pProtocol == port.Protocol {
			return &port
		}
	}
	// not found
	return nil
}

// PodDisruptionBudget specifies a PodDisruptionBudget resource
type PodDisruptionBudget struct {
	Template GeneratorFunction
	Enabled  bool
}

// HorizontalPodAutoscaler specifies a HorizontalPodAutoscaler resource
type HorizontalPodAutoscaler struct {
	Template GeneratorFunction
	Enabled  bool
}

// PodMonitor specifies a PodMonitor resource
type PodMonitor struct {
	Template GeneratorFunction
	Enabled  bool
}

// GrafanaDashboard specifies a GrafanaDashboard resource
type GrafanaDashboard struct {
	Template GeneratorFunction
	Enabled  bool
}

// ConfigMaps specifies a ConfigMap resource
type ConfigMaps struct {
	Template GeneratorFunction
	Enabled  bool
}

// GetDeploymentReplicas returns the number of replicas for a deployment,
// current value if HPA is enabled.
func (r *Reconciler) GetDeploymentReplicas(ctx context.Context, d Deployment) (*int32, error) {
	dep := d.Template().(*appsv1.Deployment)
	if !d.HasHPA {
		return dep.Spec.Replicas, nil
	}
	key := types.NamespacedName{
		Name:      dep.GetName(),
		Namespace: dep.GetNamespace(),
	}
	instance := &appsv1.Deployment{}
	err := r.GetClient().Get(ctx, key, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return dep.Spec.Replicas, nil
		}
		return dep.Spec.Replicas, err
	}
	return instance.Spec.Replicas, nil
}

// ReconcileOwnedResources handles generalized resource reconcile logic for
// all controllers
func (r *Reconciler) ReconcileOwnedResources(ctx context.Context, owner client.Object, crs ControlledResources) error {
	// Calculate resources to enforce
	resources := []LockedResource{}

	for _, dep := range crs.Deployments {

		if dep.HasHPA {
			currentReplicas, err := r.GetDeploymentReplicas(ctx, dep)
			if err != nil {
				return err
			}
			resources = append(resources,
				LockedResource{
					GeneratorFn:  r.DeploymentWithRolloutTriggers(dep.Template, dep.RolloutTriggers, currentReplicas),
					ExcludePaths: append(DeploymentExcludedPaths, "/spec/replicas"),
				})

		} else {
			resources = append(resources,
				LockedResource{
					GeneratorFn:  r.DeploymentWithRolloutTriggers(dep.Template, dep.RolloutTriggers, nil),
					ExcludePaths: DeploymentExcludedPaths,
				})
		}

	}

	for _, ss := range crs.StatefulSets {
		if ss.Enabled {
			resources = append(resources,
				LockedResource{
					GeneratorFn:  r.StatefulSetWithRolloutTriggers(ss.Template, ss.RolloutTriggers),
					ExcludePaths: DefaultExcludedPaths,
				})
		}
	}

	for _, es := range crs.ExternalSecrets {
		if es.Enabled {
			resources = append(resources,
				LockedResource{
					GeneratorFn:  es.Template,
					ExcludePaths: DefaultExcludedPaths,
				})
		}
	}

	for _, svc := range crs.Services {
		if svc.Enabled {
			template, err := svc.PopulateSpecRuntimeValues(ctx, r.GetClient())
			if err != nil {
				return err
			}
			resources = append(resources,
				LockedResource{
					GeneratorFn:  template,
					ExcludePaths: append(DefaultExcludedPaths, ServiceExcludes(svc.Template)...),
				})
		}
	}

	for _, pm := range crs.PodMonitors {
		if pm.Enabled {
			resources = append(resources,
				LockedResource{
					GeneratorFn:  pm.Template,
					ExcludePaths: DefaultExcludedPaths,
				})
		}
	}

	for _, hpa := range crs.HorizontalPodAutoscalers {
		if hpa.Enabled {
			resources = append(resources,
				LockedResource{
					GeneratorFn:  hpa.Template,
					ExcludePaths: DefaultExcludedPaths,
				})
		}
	}

	for _, pdb := range crs.PodDisruptionBudgets {
		if pdb.Enabled {
			resources = append(resources,
				LockedResource{
					GeneratorFn:  pdb.Template,
					ExcludePaths: DefaultExcludedPaths,
				})
		}
	}

	for _, gd := range crs.GrafanaDashboards {
		if gd.Enabled {
			resources = append(resources,
				LockedResource{
					GeneratorFn:  gd.Template,
					ExcludePaths: DefaultExcludedPaths,
				})
		}
	}

	for _, cm := range crs.ConfigMaps {
		if cm.Enabled {
			resources = append(resources,
				LockedResource{
					GeneratorFn:  cm.Template,
					ExcludePaths: DefaultExcludedPaths,
				})
		}
	}

	lockedResources, err := r.NewLockedResources(resources, owner)
	if err != nil {
		return err
	}
	err = r.UpdateLockedResources(ctx, owner, lockedResources, []lockedpatch.LockedPatch{})
	if err != nil {
		return err
	}

	return nil
}

// ServiceExcludes generates the list of excluded paths for a Service resource
func ServiceExcludes(fn GeneratorFunction) []string {
	svc := fn().(*corev1.Service)
	paths := []string{}
	paths = append(paths, "/spec/clusterIP", "/spec/clusterIPs", "/spec/ipFamilies", "/spec/ipFamilyPolicy")
	for idx := range svc.Spec.Ports {
		paths = append(paths, fmt.Sprintf("/spec/ports/%d/nodePort", idx))
	}
	return paths
}

// DeploymentWithRolloutTriggers returns the Deployment modified with the appropriate rollout triggers (annotations)
func (r *Reconciler) DeploymentWithRolloutTriggers(deployment GeneratorFunction,
	triggers []RolloutTrigger, replicas *int32) GeneratorFunction {

	return func() client.Object {
		dep := deployment().(*appsv1.Deployment)
		if dep.Spec.Template.ObjectMeta.Annotations == nil {
			dep.Spec.Template.ObjectMeta.Annotations = map[string]string{}
		}
		for _, trigger := range triggers {
			dep.Spec.Template.ObjectMeta.Annotations[trigger.GetAnnotationKey()] = trigger.GetHash()
		}

		if replicas != nil {
			dep.Spec.Replicas = replicas
		}

		return dep
	}
}

// StatefulSetWithRolloutTriggers returns the StatefulSet modified with the appropriate rollout triggers (annotations)
func (r *Reconciler) StatefulSetWithRolloutTriggers(statefulset GeneratorFunction,
	triggers []RolloutTrigger) GeneratorFunction {

	return func() client.Object {
		ss := statefulset().(*appsv1.StatefulSet)
		if ss.Spec.Template.ObjectMeta.Annotations == nil {
			ss.Spec.Template.ObjectMeta.Annotations = map[string]string{}
		}
		for _, trigger := range triggers {
			ss.Spec.Template.ObjectMeta.Annotations[trigger.GetAnnotationKey()] = trigger.GetHash()
		}
		return ss
	}
}

// Hash returns a hash of the passed object
func Hash(o interface{}) string {
	hasher := fnv.New32a()
	hasher.Reset()
	printer := spew.ConfigState{
		Indent:         " ",
		SortKeys:       true,
		DisableMethods: true,
		SpewKeys:       true,
	}
	printer.Fprintf(hasher, "%#v", o)
	return rand.SafeEncodeString(fmt.Sprint(hasher.Sum32()))
}
