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
	"github.com/davecgh/go-spew/spew"
	"github.com/redhat-cop/operator-utils/pkg/util/lockedresourcecontroller/lockedpatch"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client" // policyv1beta1 "k8s.io/api/policy/v1beta1"
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
}

type RolloutTrigger struct {
	name      string
	configMap *corev1.ConfigMap
	secret    *corev1.Secret
}

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

func (rt *RolloutTrigger) GetAnnotationKey() string {
	if rt.secret != nil {
		return fmt.Sprintf("%s/%s.%s", saasv1alpha1.AnnotationsDomain, rt.name, "secret-hash")
	}
	return fmt.Sprintf("%s/%s.%s", saasv1alpha1.AnnotationsDomain, rt.name, "configmap-hash")
}

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
func (r *Reconciler) TriggersFromSecretDefs(ctx context.Context, sd ...GeneratorFunction) ([]RolloutTrigger, error) {

	triggers := []RolloutTrigger{}

	for _, secretDef := range sd {
		sd := secretDef().(*secretsmanagerv1alpha1.SecretDefinition)
		key := types.NamespacedName{
			Name:      sd.GetName(),
			Namespace: sd.GetNamespace(),
		}
		secret := &corev1.Secret{}
		err := r.GetClient().Get(ctx, key, secret)
		if err != nil {
			if errors.IsNotFound(err) {
				triggers = append(triggers, NewRolloutTrigger(sd.GetName(), &corev1.Secret{}))
				continue
			}
			return nil, err
		}

		triggers = append(triggers, NewRolloutTrigger(sd.GetName(), secret))

	}

	return triggers, nil
}

// Deployment specifies a Deployment resources and its rollout triggers
type Deployment struct {
	Template        GeneratorFunction
	RolloutTriggers []RolloutTrigger
	HasHPA          bool
}

// StatefulSet ...
type StatefulSet struct {
	Template        GeneratorFunction
	RolloutTriggers []RolloutTrigger
	Enabled         bool
}

// SecretDefinition ...
type SecretDefinition struct {
	Template GeneratorFunction
	Enabled  bool
}

// Service ...
type Service struct {
	Template GeneratorFunction
	Enabled  bool
}

// PodDisruptionBudget ...
type PodDisruptionBudget struct {
	Template GeneratorFunction
	Enabled  bool
}

// HorizontalPodAutoscaler ...
type HorizontalPodAutoscaler struct {
	Template GeneratorFunction
	Enabled  bool
}

// PodMonitor ...
type PodMonitor struct {
	Template GeneratorFunction
	Enabled  bool
}

// GrafanaDashboard ...
type GrafanaDashboard struct {
	Template GeneratorFunction
	Enabled  bool
}

// ReconcileOwnedResources handles generalized resource reconcile logic for
// all controllers
func (r *Reconciler) ReconcileOwnedResources(ctx context.Context, owner client.Object, crs ControlledResources) error {
	// Calculate resources to enforce
	resources := []LockedResource{}

	for _, dep := range crs.Deployments {

		resources = append(resources,
			LockedResource{
				GeneratorFn: r.DeploymentWithRolloutTriggers(dep.Template, dep.RolloutTriggers),
				ExcludePaths: func() []string {
					if dep.HasHPA {
						return append(DeploymentExcludedPaths, "/spec/replicas")
					}
					return DeploymentExcludedPaths
				}(),
			})
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

	for _, sd := range crs.SecretDefinitions {
		if sd.Enabled {
			resources = append(resources,
				LockedResource{
					GeneratorFn:  sd.Template,
					ExcludePaths: DefaultExcludedPaths,
				})
		}
	}

	for _, svc := range crs.Services {
		if svc.Enabled {
			resources = append(resources,
				LockedResource{
					GeneratorFn:  svc.Template,
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

	lockedResources, err := r.NewLockedResources(resources, owner)
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
	paths = append(paths, "/spec/clusterIP")
	for idx := range svc.Spec.Ports {
		paths = append(paths, fmt.Sprintf("/spec/ports/%d/nodePort", idx))
	}
	return paths
}

func (r *Reconciler) DeploymentWithRolloutTriggers(deployment GeneratorFunction, triggers []RolloutTrigger) GeneratorFunction {

	return func() client.Object {
		dep := deployment().(*appsv1.Deployment)
		if dep.GetAnnotations() == nil {
			dep.Spec.Template.ObjectMeta.Annotations = map[string]string{}
		}
		for _, trigger := range triggers {
			dep.Spec.Template.ObjectMeta.Annotations[trigger.GetAnnotationKey()] = trigger.GetHash()
		}
		return dep
	}
}

func (r *Reconciler) StatefulSetWithRolloutTriggers(statefulset GeneratorFunction, triggers []RolloutTrigger) GeneratorFunction {

	return func() client.Object {
		ss := statefulset().(*appsv1.StatefulSet)
		if ss.GetAnnotations() == nil {
			ss.Spec.Template.ObjectMeta.Annotations = map[string]string{}
		}
		for _, trigger := range triggers {
			ss.Spec.Template.ObjectMeta.Annotations[trigger.GetAnnotationKey()] = trigger.GetHash()
		}
		return ss
	}
}

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
