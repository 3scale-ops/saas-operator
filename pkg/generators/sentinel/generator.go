package sentinel

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/3scale-ops/basereconciler/mutators"
	"github.com/3scale-ops/basereconciler/resource"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/generators"
	"github.com/3scale-ops/saas-operator/pkg/generators/sentinel/config"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/grafanadashboard"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/pdb"
	operatorutils "github.com/3scale-ops/saas-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
)

const (
	component string = "redis-sentinel"
)

// Generator configures the generators for Sentinel
type Generator struct {
	generators.BaseOptionsV2
	Spec    saasv1alpha1.SentinelSpec
	Options config.Options
}

// NewGenerator returns a new Options struct
func NewGenerator(instance, namespace string, spec saasv1alpha1.SentinelSpec) Generator {
	return Generator{
		BaseOptionsV2: generators.BaseOptionsV2{
			Component:    component,
			InstanceName: instance,
			Namespace:    namespace,
			Labels: map[string]string{
				"app":     component,
				"part-of": "3scale-saas",
			},
		},
		Spec:    spec,
		Options: config.NewOptions(spec),
	}
}

// Resources returns a list of templates
func (gen *Generator) Resources() []resource.TemplateInterface {
	resources := []resource.TemplateInterface{
		resource.NewTemplateFromObjectFunction(gen.statefulSet),
		resource.NewTemplateFromObjectFunction(gen.statefulSetService).WithMutation(mutators.SetServiceLiveValues()),
		resource.NewTemplate(pdb.New(gen.GetKey(), gen.GetLabels(), gen.GetSelector(), *gen.Spec.PDB)),
		resource.NewTemplateFromObjectFunction(gen.configMap),
		resource.NewTemplate(grafanadashboard.New(gen.GetKey(), gen.GetLabels(), *gen.Spec.GrafanaDashboard, "dashboards/redis-sentinel.json.gtpl")).
			WithEnabled(!gen.Spec.GrafanaDashboard.IsDeactivated()),
	}

	for idx := 0; idx < int(*gen.Spec.Replicas); idx++ {
		i := idx
		resources = append(resources,
			resource.NewTemplateFromObjectFunction(
				func() *corev1.Service { return gen.podServices(i) }).
				WithMutation(mutators.SetServiceLiveValues()),
		)
	}

	return resources
}

func (gen *Generator) ClusterTopology(ctx context.Context) (map[string]map[string]string, error) {

	clustermap := map[string]map[string]string{}

	if gen.Spec.Config.ClusterTopology != nil {
		for shard, serversdef := range gen.Spec.Config.ClusterTopology {
			shardmap := map[string]string{}
			for alias, server := range serversdef {
				// the redis servers must be defined using IP
				// addresses, so this tries to resolve a hostname
				// if present in the connection string.
				u, err := url.Parse(server)
				if err != nil {
					return nil, err
				}
				ip, err := operatorutils.LookupIPv4(ctx, u.Hostname())
				if err != nil {
					return nil, err
				}
				u.Host = net.JoinHostPort(ip, u.Port())
				if err != nil {
					return nil, err
				}
				shardmap[alias] = u.String()
			}
			clustermap[shard] = shardmap
		}

	} else if gen.Spec.Config.MonitoredShards != nil {
		for shard, servers := range gen.Spec.Config.MonitoredShards {
			shardmap := map[string]string{}
			for _, server := range servers {
				// the redis servers must be defined using IP
				// addresses, so this tries to resolve a hostname
				// if present in the connection string.
				u, err := url.Parse(server)
				if err != nil {
					return nil, err
				}
				alias := u.Host
				ip, err := operatorutils.LookupIPv4(ctx, u.Hostname())
				if err != nil {
					return nil, err
				}
				u.Host = net.JoinHostPort(ip, u.Port())
				if err != nil {
					return nil, err
				}
				shardmap[alias] = u.String()
			}
			clustermap[shard] = shardmap
		}

	} else {
		return nil, fmt.Errorf("either 'spec.config.clusterTopology' or 'spec.cluster.MonitoredShards' must be set")
	}

	clustermap["sentinel"] = make(map[string]string, int(*gen.Spec.Replicas))
	for _, uri := range gen.SentinelURIs() {
		u, err := url.Parse(uri)
		if err != nil {
			return nil, err
		}
		alias := strings.Split(u.Hostname(), ".")[0]
		clustermap["sentinel"][alias] = u.String()
	}

	return clustermap, nil
}
