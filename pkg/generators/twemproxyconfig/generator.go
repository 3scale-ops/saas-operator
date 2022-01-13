package twemproxyconfig

import (
	"context"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	basereconciler_resources "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2/resources"
	"github.com/3scale/saas-operator/pkg/redis"
)

const (
	component string = "twemproxy"
)

// Generator configures the generators for Sentinel
type Generator struct {
	generators.BaseOptionsV2
	Spec            saasv1alpha1.TwemproxyConfigSpec
	monitoredShards map[string]TwemproxyServer
}

// NewGenerator returns a new Options struct
func NewGenerator(ctx context.Context, instance *saasv1alpha1.TwemproxyConfig) (Generator, error) {

	gen := Generator{
		BaseOptionsV2: generators.BaseOptionsV2{
			Component:    component,
			InstanceName: instance.GetName(),
			Namespace:    instance.GetNamespace(),
			Labels: map[string]string{
				"app":     component,
				"part-of": "3scale-saas",
			},
		},
		Spec: instance.Spec,
	}

	var err error
	gen.monitoredShards, err = gen.getMonitoredShards(ctx)
	if err != nil {
		return Generator{}, err
	}

	return gen, nil
}

func (gen *Generator) getMonitoredShards(ctx context.Context) (map[string]TwemproxyServer, error) {

	spool := make(redis.SentinelPool, 0, len(gen.Spec.SentinelURIs))

	for _, uri := range gen.Spec.SentinelURIs {
		sentinel, err := redis.NewSentinelServerFromConnectionString("sentinel", uri)
		if err != nil {
			return nil, err
		}

		spool = append(spool, *sentinel)
	}

	monitoredShards, err := spool.MonitoredShards(ctx, 2)
	if err != nil {
		return nil, err
	}

	m := make(map[string]TwemproxyServer, len(monitoredShards))
	for _, s := range monitoredShards {
		m[s.Name] = TwemproxyServer{
			Name:     s.Name,
			Address:  s.Master,
			Priority: 1,
		}
	}

	return m, nil
}

// Returns all the resource templates that this generator manages
func (gen *Generator) Resources() []basereconciler.Resource {
	return []basereconciler.Resource{
		basereconciler_resources.ConfigMapTemplate{
			Template:  gen.configMap(true),
			IsEnabled: true,
		},
	}
}
