package dynamic_config

import (
	"testing"

	"github.com/3scale-ops/marin3r/pkg/envoy"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/resource_builders/envoyconfig/templates"
	"github.com/MakeNowJust/heredoc"
	"github.com/davecgh/go-spew/spew"
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_service_runtime_v3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
)

type unregisteredType struct{ opts *opts }
type opts struct{}

func (x *unregisteredType) GetOptions() interface{}     { return x.opts }
func (x *unregisteredType) GetGeneratorVersion() string { return "" }
func (x *unregisteredType) GetName() string             { return "" }

var testFactory = EnvoyDynamicConfigFactory{
	"ListenerHttp_v1":       RegisterTemplate(templates.ListenerHTTP_v1, &envoy_config_listener_v3.Listener{}),
	"Cluster_v1":            RegisterTemplate(templates.Cluster_v1, &envoy_config_cluster_v3.Cluster{}),
	"RouteConfiguration_v1": RegisterTemplate(templates.RouteConfiguration_v1, &envoy_config_route_v3.RouteConfiguration{}),
	"Runtime_v1":            RegisterTemplate(templates.Runtime_v1, &envoy_service_runtime_v3.Runtime{}),
	"RawConfig_v1":          RegisterTemplate(templates.RawConfig_v1, nil),
}

func TestEnvoyDynamicConfigFactory_NewResource(t *testing.T) {
	type args struct {
		descriptor EnvoyDynamicConfigDescriptor
	}
	tests := []struct {
		name    string
		factory EnvoyDynamicConfigFactory
		args    args
		want    envoy.Resource
		wantErr bool
	}{
		{
			name:    "Generates a cluster proto",
			factory: testFactory,
			args: args{
				descriptor: &saasv1alpha1.EnvoyDynamicConfig{
					EnvoyDynamicConfigMeta: saasv1alpha1.EnvoyDynamicConfigMeta{
						Name:             "test",
						GeneratorVersion: pointer.String("v1"),
					},
					Cluster: &saasv1alpha1.Cluster{
						Host:    "127.0.0.1",
						Port:    8080,
						IsHttp2: pointer.Bool(true),
					},
				},
			},
			want: func() envoy.Resource {
				c, _ := templates.Cluster_v1("test", &saasv1alpha1.Cluster{
					Host:    "127.0.0.1",
					Port:    8080,
					IsHttp2: pointer.Bool(true),
				})
				return c
			}(),
			wantErr: false,
		},
		{
			name:    "Generates a cluster from a RawConfig",
			factory: testFactory,
			args: args{
				descriptor: &saasv1alpha1.EnvoyDynamicConfig{
					EnvoyDynamicConfigMeta: saasv1alpha1.EnvoyDynamicConfigMeta{
						Name:             "test",
						GeneratorVersion: pointer.String("v1"),
					},
					RawConfig: &saasv1alpha1.RawConfig{
						Type: "cluster",
						Value: runtime.RawExtension{
							Raw: []byte(heredoc.Doc(`
								{
								  "load_assignment": {
								    "cluster_name": "cluster1"
								  },
								  "name": "cluster1"
								}
							`)),
						},
					},
				},
			},
			want: &envoy_config_cluster_v3.Cluster{
				Name: "cluster1",
				LoadAssignment: &envoy_config_endpoint_v3.ClusterLoadAssignment{
					ClusterName: "cluster1",
					Endpoints:   []*envoy_config_endpoint_v3.LocalityLbEndpoints{},
				},
			},
			wantErr: false,
		},
		{
			name:    "Generates a runtime proto",
			factory: testFactory,
			args: args{
				descriptor: &saasv1alpha1.EnvoyDynamicConfig{
					EnvoyDynamicConfigMeta: saasv1alpha1.EnvoyDynamicConfigMeta{
						Name:             "test",
						GeneratorVersion: pointer.String("v1"),
					},
					Runtime: &saasv1alpha1.Runtime{
						ListenerNames: []string{"http", "https"},
					},
				},
			},
			want: func() envoy.Resource {
				l, _ := structpb.NewStruct(map[string]interface{}{
					"envoy": map[string]interface{}{
						"resource_limits": map[string]interface{}{
							"listener": map[string]interface{}{
								"http": map[string]interface{}{
									"connection_limit": 10000,
								},
								"https": map[string]interface{}{
									"connection_limit": 10000,
								},
							},
						},
					},
					"overload": map[string]interface{}{
						"global_downstream_max_connections": 50000,
					},
				})
				return &envoy_service_runtime_v3.Runtime{
					Name:  "test",
					Layer: l,
				}
			}(),
			wantErr: false,
		},
		{
			name:    "Unregistered class",
			factory: testFactory,
			args: args{
				descriptor: &unregisteredType{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.factory.NewResource(tt.args.descriptor)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnvoyDynamicConfigFactory.NewResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(got, tt.want) {
				spew.Dump(got)
				spew.Dump(tt.want)
				t.Errorf("EnvoyDynamicConfigFactory.NewResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
