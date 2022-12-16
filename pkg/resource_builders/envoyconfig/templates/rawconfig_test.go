package templates

import (
	"testing"

	"github.com/3scale-ops/marin3r/pkg/envoy"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/MakeNowJust/heredoc"
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"google.golang.org/protobuf/proto"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestRawConfig_v1(t *testing.T) {
	type args struct {
		name string
		opts interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    envoy.Resource
		wantErr bool
	}{
		{
			name: "Generates the corresponding proto msg",
			args: args{
				name: "test",
				opts: &saasv1alpha1.RawConfig{
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
			want: &envoy_config_cluster_v3.Cluster{
				Name: "cluster1",
				LoadAssignment: &envoy_config_endpoint_v3.ClusterLoadAssignment{
					ClusterName: "cluster1",
					Endpoints:   []*envoy_config_endpoint_v3.LocalityLbEndpoints{},
				},
			}, wantErr: false,
		},
		{
			name: "Returns an error",
			args: args{
				name: "test",
				opts: &saasv1alpha1.RawConfig{
					Type: "listener",
					Value: runtime.RawExtension{
						Raw: []byte(heredoc.Doc(`
							{
							  "wrong_key": "value",
							  "name": "listener"
							}
						`)),
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RawConfig_v1(tt.args.name, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("RawConfig_v1() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(got, tt.want) {
				t.Errorf("RawConfig_v1() = %v, want %v", got, tt.want)
			}
		})
	}
}
