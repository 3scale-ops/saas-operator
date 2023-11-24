package redisshard

import (
	"fmt"

	"github.com/3scale-ops/basereconciler/resource"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	component string = "redis-shard"
)

// Generator configures the generators for RedisShard
type Generator struct {
	generators.BaseOptionsV2
	Image       saasv1alpha1.ImageSpec
	MasterIndex int32
	Replicas    int32
	Command     string
}

// Override the GetSelector function as it needs to be different in this case
// because there can be more than one redis-shard instance in the same namespace
func (gen *Generator) GetSelector() map[string]string {
	return map[string]string{"redis-shard": gen.GetInstanceName()}
}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.RedisShardSpec) Generator {
	return Generator{
		BaseOptionsV2: generators.BaseOptionsV2{
			Component:    component,
			InstanceName: instance,
			Namespace:    namespace,
			Labels: map[string]string{
				"app":     component,
				"part-of": "3scale-saas-testing",
			},
		},
		Image:       *spec.Image,
		MasterIndex: *spec.MasterIndex,
		Replicas:    *spec.SlaveCount + 1,
		Command:     *spec.Command,
	}
}

// Returns the name of the StatefulSet headless Service
func (gen *Generator) ServiceName() string {
	return fmt.Sprintf("%s-%s", gen.GetComponent(), gen.GetInstanceName())
}

// Returns all the resource templates that this generator manages
func (gen *Generator) Resources() []resource.TemplateInterface {
	return []resource.TemplateInterface{
		resource.NewTemplateFromObjectFunction[*appsv1.StatefulSet](gen.statefulSet),
		resource.NewTemplateFromObjectFunction[*corev1.Service](gen.service),
		resource.NewTemplateFromObjectFunction[*corev1.ConfigMap](gen.redisConfigConfigMap),
		resource.NewTemplateFromObjectFunction[*corev1.ConfigMap](gen.redisReadinessScriptConfigMap),
	}
}
