package redisshard

import (
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"
)

const (
	component string = "redis-shard"
)

// Generator configures the generators for RedisShard
type Generator struct {
	generators.BaseOptions
	Image       saasv1alpha1.ImageSpec
	MasterIndex int32
}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.RedisShardSpec) Generator {
	return Generator{
		BaseOptions: generators.BaseOptions{
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
