package twemproxyconfig

import (
	"context"
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	basereconciler_resources "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2/resources"
	"github.com/3scale/saas-operator/pkg/redis"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	component string = "twemproxy"
)

// Generator configures the generators for Sentinel
type Generator struct {
	generators.BaseOptionsV2
	Spec           saasv1alpha1.TwemproxyConfigSpec
	masterTargets  map[string]TwemproxyServer
	slaverwTargets map[string]TwemproxyServer
}

// NewGenerator returns a new Options struct
func NewGenerator(ctx context.Context, instance *saasv1alpha1.TwemproxyConfig, cl client.Client) (Generator, error) {

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
	if gen.Spec.SentinelURIs == nil {
		gen.Spec.SentinelURIs, err = discoverSentinels(ctx, cl, instance.GetNamespace())
		if err != nil {
			return Generator{}, err
		}
	}

	gen.masterTargets, err = gen.getMonitoredMasters(ctx)
	if err != nil {
		return Generator{}, err
	}

	// Check if there are pools in the config that require slave discovery
	discoverSlavesRW := false
	for _, pool := range gen.Spec.ServerPools {
		if *pool.Target == saasv1alpha1.SlavesRW {
			discoverSlavesRW = true
		}
	}
	if discoverSlavesRW {
		gen.slaverwTargets, err = gen.getMonitoredReadWriteSlavesWithFallbackToMasters(ctx)
		if err != nil {
			return Generator{}, err
		}
	}

	return gen, nil
}

func discoverSentinels(ctx context.Context, cl client.Client, namespace string) ([]string, error) {
	sl := &saasv1alpha1.SentinelList{}
	if err := cl.List(ctx, sl, client.InNamespace(namespace)); err != nil {
		return nil, err
	}

	if len(sl.Items) != 1 {
		return nil, fmt.Errorf("unexpected number (%d) of Sentinel resources in namespace", len(sl.Items))
	}

	uris := make([]string, 0, len(sl.Items[0].Status.Sentinels))
	for _, address := range sl.Items[0].Status.Sentinels {
		uris = append(uris, fmt.Sprintf("redis://%s", address))
	}

	return uris, nil
}

func (gen *Generator) getMonitoredMasters(ctx context.Context) (map[string]TwemproxyServer, error) {

	spool := make(redis.SentinelPool, 0, len(gen.Spec.SentinelURIs))

	for _, uri := range gen.Spec.SentinelURIs {
		sentinel, err := redis.NewSentinelServerFromConnectionString("sentinel", uri)
		if err != nil {
			return nil, err
		}

		spool = append(spool, *sentinel)
	}

	monitoredShards, err := spool.MonitoredShards(ctx, 2, false)
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

func (gen *Generator) getMonitoredReadWriteSlavesWithFallbackToMasters(ctx context.Context) (map[string]TwemproxyServer, error) {

	spool := make(redis.SentinelPool, 0, len(gen.Spec.SentinelURIs))

	for _, uri := range gen.Spec.SentinelURIs {
		sentinel, err := redis.NewSentinelServerFromConnectionString("sentinel", uri)
		if err != nil {
			return nil, err
		}

		spool = append(spool, *sentinel)
	}

	monitoredShards, err := spool.MonitoredShards(ctx, 2, true)
	if err != nil {
		return nil, err
	}

	m := make(map[string]TwemproxyServer, len(monitoredShards))
	for _, shard := range monitoredShards {

		if len(shard.SlavesRW) > 0 {
			m[shard.Name] = TwemproxyServer{
				Name:     shard.Name,
				Address:  shard.SlavesRW[0],
				Priority: 1,
			}
		} else {
			// Fall back to masters if there are no
			// available RW slaves
			m[shard.Name] = TwemproxyServer{
				Name:     shard.Name,
				Address:  shard.Master,
				Priority: 1,
			}
		}
	}

	return m, nil
}

// Returns the twemproxy config ConfigMap
func (gen *Generator) ConfigMap() basereconciler.Resource {
	return basereconciler_resources.ConfigMapTemplate{
		Template:  gen.configMap(true),
		IsEnabled: true,
	}
}
