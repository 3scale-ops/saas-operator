package redisshard

import (
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	basereconciler_resources "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2/resources"
)

const (
	component string = "redis-shard"
)

// Generator configures the generators for RedisShard
type Generator struct {
	generators.BaseOptionsV2
	Image       saasv1alpha1.ImageSpec
	MasterIndex int32
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
	}
}

func (gen *Generator) ServiceName() string {
	return fmt.Sprintf("%s-%s", gen.GetComponent(), gen.GetInstanceName())
}

func (gen *Generator) Resources() []basereconciler.Resource {
	return []basereconciler.Resource{
		basereconciler_resources.StatefulSetTemplate{
			Template:        gen.statefulSet(),
			IsEnabled:       true,
			RolloutTriggers: nil,
		},
		basereconciler_resources.ServiceTemplate{
			Template:  gen.service(),
			IsEnabled: true,
		},
		basereconciler_resources.ConfigMapTemplate{
			Template:  gen.redisConfigConfigMap(),
			IsEnabled: true,
		},
		basereconciler_resources.ConfigMapTemplate{
			Template:  gen.redisReadinessScriptConfigMap(),
			IsEnabled: true,
		},
	}
}
