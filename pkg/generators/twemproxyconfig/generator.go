package twemproxyconfig

import (
	"context"
	"fmt"
	"reflect"
	"sort"

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

	responses := make([]saasv1alpha1.MonitoredShards, 0, len(gen.Spec.SentinelURIs))

	for _, uri := range gen.Spec.SentinelURIs {
		sentinel, err := redis.NewSentinelServerFromConnectionString("sentinel", uri)
		if err != nil {
			return nil, err
		}

		resp, err := sentinel.MonitoredShards(ctx)
		if err != nil {
			// jump to next sentinel if error occurs
			continue
		}
		responses = append(responses, resp)
	}

	monitoredShards, err := applyQuorum(responses, 2)
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

func applyQuorum(responses []saasv1alpha1.MonitoredShards, quorum int) (saasv1alpha1.MonitoredShards, error) {

	for _, r := range responses {
		// Sort each of the MonitoredShards responses to
		// avoid diffs due to unordered responses from redis
		sort.Sort(r)
	}

	for idx, a := range responses {
		count := 0
		for _, b := range responses {
			if reflect.DeepEqual(a, b) {
				count++
			}
		}

		// check if this response has quorum
		if count >= quorum {
			return responses[idx], nil
		}
	}

	return nil, fmt.Errorf("unable to get monitored shards from sentinel")
}
