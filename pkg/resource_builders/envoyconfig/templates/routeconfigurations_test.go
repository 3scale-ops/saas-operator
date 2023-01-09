package templates

import (
	"testing"

	envoy_serializer_v3 "github.com/3scale-ops/marin3r/pkg/envoy/serializer/v3"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/MakeNowJust/heredoc"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

func TestRouteConfiguration_v1(t *testing.T) {
	type args struct {
		name string
		opts *saasv1alpha1.RouteConfiguration
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Generate a route with the given virtual hosts",
			args: args{
				name: "my_route",
				opts: &saasv1alpha1.RouteConfiguration{
					VirtualHosts: []runtime.RawExtension{
						{
							Raw: []byte(`{"name":"example","domains":["example.com"],"routes":[{"route":{"cluster":"example_cluster"},"match":{"prefix":"/"}}],"rate_limits":[{"actions":[{"remote_address":{}}]}]}`),
						},
						{
							Raw: []byte(`{"name":"example2","domains":["example2.com"],"routes":[{"direct_response":{"status":"403"},"match":{"prefix":"/forbidden"}}]}`),
						},
					},
				},
			},
			want: heredoc.Doc(`
                name: my_route
                virtual_hosts:
                - domains:
                  - example.com
                  name: example
                  rate_limits:
                  - actions:
                    - remote_address: {}
                  routes:
                  - match:
                      prefix: /
                    route:
                      cluster: example_cluster
                - domains:
                  - example2.com
                  name: example2
                  routes:
                  - direct_response:
                      status: 403
                    match:
                      prefix: /forbidden
			`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RouteConfiguration_v1(tt.args.name, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouteConfiguration_v1() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			j, err := envoy_serializer_v3.JSON{}.Marshal(got)
			if err != nil {
				t.Error(err)
			}
			y, err := yaml.JSONToYAML([]byte(j))
			if err != nil {
				t.Error(err)
			}
			if string(y) != tt.want {
				t.Errorf("RouteConfiguration_v1():\n# got:\n%v\n# want:\n%v", string(y), tt.want)
			}
		})
	}
}
