package basereconciler

import (
	"context"
	"fmt"
	"hash/fnv"
	"reflect"

	// grafanav1alpha1 "github.com/3scale/saas-operator/pkg/apis/grafana/v1alpha1"
	// secretsmanagerv1alpha1 "github.com/3scale/saas-operator/pkg/apis/secrets-manager/v1alpha1"
	// monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	// autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	secretsmanagerv1alpha1 "github.com/3scale/saas-operator/pkg/apis/secrets-manager/v1alpha1"
	basereconciler_types "github.com/3scale/saas-operator/pkg/basereconciler/types"
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
	SecretDefinitions        []SecretDefinition
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
	cm.SecretDefinitions = append(cm.SecretDefinitions, resources.SecretDefinitions...)
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

// TriggersFromSecretDefs generates a list of RolloutTrigger from the given SecretDefinition generator functions
func (r *Reconciler) TriggersFromSecretDefs(ctx context.Context, sd ...basereconciler_types.GeneratorFunction) ([]RolloutTrigger, error) {

	triggers := []RolloutTrigger{}

	for _, secretDef := range sd {
		sd := secretDef().(*secretsmanagerv1alpha1.SecretDefinition)
		strgs, err := r.TriggersFromSecret(ctx, sd.GetNamespace(), sd.GetName())
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
	Template        basereconciler_types.GeneratorFunction
	RolloutTriggers []RolloutTrigger
	HasHPA          bool
}

// StatefulSet specifies a StatefulSet resource and its rollout triggers
type StatefulSet struct {
	Template        basereconciler_types.GeneratorFunction
	RolloutTriggers []RolloutTrigger
	Enabled         bool
}

// SecretDefinition specifies a SecretDefinition resource
type SecretDefinition struct {
	Template basereconciler_types.GeneratorFunction
	Enabled  bool
}

// Service specifies a Service resource
type Service struct {
	Template basereconciler_types.GeneratorFunction
	Enabled  bool
}

// PodDisruptionBudget specifies a PodDisruptionBudget resource
type PodDisruptionBudget struct {
	Template basereconciler_types.GeneratorFunction
	Enabled  bool
}

// HorizontalPodAutoscaler specifies a HorizontalPodAutoscaler resource
type HorizontalPodAutoscaler struct {
	Template basereconciler_types.GeneratorFunction
	Enabled  bool
}

// PodMonitor specifies a PodMonitor resource
type PodMonitor struct {
	Template basereconciler_types.GeneratorFunction
	Enabled  bool
}

// GrafanaDashboard specifies a GrafanaDashboard resource
type GrafanaDashboard struct {
	Template basereconciler_types.GeneratorFunction
	Enabled  bool
}

// ConfigMaps specifies a ConfigMap resource
type ConfigMaps struct {
	Template basereconciler_types.GeneratorFunction
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
	resources := []basereconciler_types.LockedResource{}

	for _, dep := range crs.Deployments {

		if dep.HasHPA {
			currentReplicas, err := r.GetDeploymentReplicas(ctx, dep)
			if err != nil {
				return err
			}
			resources = append(resources,
				basereconciler_types.LockedResource{
					GeneratorFn:  r.DeploymentWithRolloutTriggers(dep.Template, dep.RolloutTriggers, currentReplicas),
					ExcludePaths: append(DeploymentExcludedPaths, "/spec/replicas"),
				})

		} else {
			resources = append(resources,
				basereconciler_types.LockedResource{
					GeneratorFn:  r.DeploymentWithRolloutTriggers(dep.Template, dep.RolloutTriggers, nil),
					ExcludePaths: DeploymentExcludedPaths,
				})
		}

	}

	for _, ss := range crs.StatefulSets {
		if ss.Enabled {
			resources = append(resources,
				basereconciler_types.LockedResource{
					GeneratorFn:  r.StatefulSetWithRolloutTriggers(ss.Template, ss.RolloutTriggers),
					ExcludePaths: DefaultExcludedPaths,
				})
		}
	}

	for _, sd := range crs.SecretDefinitions {
		if sd.Enabled {
			resources = append(resources,
				basereconciler_types.LockedResource{
					GeneratorFn:  sd.Template,
					ExcludePaths: DefaultExcludedPaths,
				})
		}
	}

	for _, svc := range crs.Services {
		if svc.Enabled {
			resources = append(resources,
				basereconciler_types.LockedResource{
					GeneratorFn:  svc.Template,
					ExcludePaths: append(DefaultExcludedPaths, ServiceExcludes(svc.Template)...),
				})
		}
	}

	for _, pm := range crs.PodMonitors {
		if pm.Enabled {
			resources = append(resources,
				basereconciler_types.LockedResource{
					GeneratorFn:  pm.Template,
					ExcludePaths: DefaultExcludedPaths,
				})
		}
	}

	for _, hpa := range crs.HorizontalPodAutoscalers {
		if hpa.Enabled {
			resources = append(resources,
				basereconciler_types.LockedResource{
					GeneratorFn:  hpa.Template,
					ExcludePaths: DefaultExcludedPaths,
				})
		}
	}

	for _, pdb := range crs.PodDisruptionBudgets {
		if pdb.Enabled {
			resources = append(resources,
				basereconciler_types.LockedResource{
					GeneratorFn:  pdb.Template,
					ExcludePaths: DefaultExcludedPaths,
				})
		}
	}

	for _, gd := range crs.GrafanaDashboards {
		if gd.Enabled {
			resources = append(resources,
				basereconciler_types.LockedResource{
					GeneratorFn:  gd.Template,
					ExcludePaths: DefaultExcludedPaths,
				})
		}
	}

	for _, cm := range crs.ConfigMaps {
		if cm.Enabled {
			resources = append(resources,
				basereconciler_types.LockedResource{
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
func ServiceExcludes(fn basereconciler_types.GeneratorFunction) []string {
	svc := fn().(*corev1.Service)
	paths := []string{}
	paths = append(paths, "/spec/clusterIP", "/spec/clusterIPs", "/spec/ipFamilies", "/spec/ipFamilyPolicy")
	for idx := range svc.Spec.Ports {
		paths = append(paths, fmt.Sprintf("/spec/ports/%d/nodePort", idx))
	}
	return paths
}

// DeploymentWithRolloutTriggers returns the Deployment modified with the appropriate rollout triggers (annotations)
func (r *Reconciler) DeploymentWithRolloutTriggers(deployment basereconciler_types.GeneratorFunction,
	triggers []RolloutTrigger, replicas *int32) basereconciler_types.GeneratorFunction {

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
func (r *Reconciler) StatefulSetWithRolloutTriggers(statefulset basereconciler_types.GeneratorFunction,
	triggers []RolloutTrigger) basereconciler_types.GeneratorFunction {

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
