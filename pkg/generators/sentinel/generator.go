package sentinel

import (
	"context"
	"net"
	"net/url"
	"strings"

	basereconciler "github.com/3scale-ops/basereconciler/reconciler"
	basereconciler_resources "github.com/3scale-ops/basereconciler/resources"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/generators/sentinel/config"
	"github.com/3scale/saas-operator/pkg/resource_builders/grafanadashboard"
	"github.com/3scale/saas-operator/pkg/resource_builders/pdb"
	"github.com/3scale/saas-operator/pkg/util"
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

// Returns all the resource templates that this generator manages
func (gen *Generator) Resources() []basereconciler.Resource {
	resources := []basereconciler.Resource{
		basereconciler_resources.StatefulSetTemplate{
			Template:        gen.statefulSet(),
			RolloutTriggers: []basereconciler_resources.RolloutTrigger{},
			IsEnabled:       true,
		},
		basereconciler_resources.ServiceTemplate{
			Template:  gen.statefulSetService(),
			IsEnabled: true,
		},
		basereconciler_resources.PodDisruptionBudgetTemplate{
			Template:  pdb.New(gen.GetKey(), gen.GetLabels(), gen.GetSelector(), *gen.Spec.PDB),
			IsEnabled: !gen.Spec.PDB.IsDeactivated(),
		},
		basereconciler_resources.ConfigMapTemplate{
			Template:  gen.configMap(),
			IsEnabled: true,
		},
		basereconciler_resources.GrafanaDashboardTemplate{
			Template:  grafanadashboard.New(gen.GetKey(), gen.GetLabels(), *gen.Spec.GrafanaDashboard, "dashboards/redis-sentinel.json.gtpl"),
			IsEnabled: !gen.Spec.GrafanaDashboard.IsDeactivated(),
		},
	}

	for idx := 0; idx < int(*gen.Spec.Replicas); idx++ {
		resources = append(resources,
			basereconciler_resources.ServiceTemplate{Template: gen.podServices(idx), IsEnabled: true})
	}

	return resources
}

func (gen *Generator) ClusterTopology(ctx context.Context) (map[string]map[string]string, error) {

	clustermap := map[string]map[string]string{}

	for shard, servers := range gen.Spec.Config.MonitoredShards {
		shardmap := map[string]string{}
		for _, server := range servers {
			// the redis servesr must be defined using IP
			// addresses, so this tries to resolve a hostname
			// if present in the connection string.
			u, err := url.Parse(server)
			alias := u.Host
			if err != nil {
				return nil, err
			}
			ip, err := util.LookupIPv4(ctx, u.Hostname())
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
