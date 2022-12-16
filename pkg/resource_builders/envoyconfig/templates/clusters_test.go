package templates

import (
	"testing"

	envoy_serializer_v3 "github.com/3scale-ops/marin3r/pkg/envoy/serializer/v3"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/MakeNowJust/heredoc"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/yaml"
)

func TestCluster_v1(t *testing.T) {
	type args struct {
		name string
		opts interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Generates http 1.1 cluster",
			args: args{
				name: "my_cluster",
				opts: &saasv1alpha1.Cluster{
					Host:    "localhost",
					Port:    8080,
					IsHttp2: pointer.Bool(false),
				},
			},
			want: heredoc.Doc(`
                connect_timeout: 1s
                dns_lookup_family: V4_ONLY
                load_assignment:
                  cluster_name: my_cluster
                  endpoints:
                  - lb_endpoints:
                    - endpoint:
                        address:
                          socket_address:
                            address: localhost
                            port_value: 8080
                name: my_cluster
                type: STRICT_DNS
			`),
		},
		{
			name: "Generates http 1.1 cluster",
			args: args{
				name: "my_cluster",
				opts: &saasv1alpha1.Cluster{
					Host:    "localhost",
					Port:    8080,
					IsHttp2: pointer.Bool(true),
				},
			},
			want: heredoc.Doc(`
                connect_timeout: 1s
                dns_lookup_family: V4_ONLY
                load_assignment:
                  cluster_name: my_cluster
                  endpoints:
                  - lb_endpoints:
                    - endpoint:
                        address:
                          socket_address:
                            address: localhost
                            port_value: 8080
                name: my_cluster
                type: STRICT_DNS
                typed_extension_protocol_options:
                  envoy.extensions.upstreams.http.v3.HttpProtocolOptions:
                    '@type': type.googleapis.com/envoy.extensions.upstreams.http.v3.HttpProtocolOptions
                    explicit_http_config:
                      http2_protocol_options:
                        initial_connection_window_size: 1048576
                        initial_stream_window_size: 65536
			`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Cluster_v1(tt.args.name, tt.args.opts)
			j, err := envoy_serializer_v3.JSON{}.Marshal(got)
			if err != nil {
				t.Error(err)
			}
			y, err := yaml.JSONToYAML([]byte(j))
			if err != nil {
				t.Error(err)
			}
			if string(y) != tt.want {
				t.Errorf("Cluster_v1():\n# got:\n%v\n# want:\n%v", string(y), tt.want)
			}

		})
	}
}
