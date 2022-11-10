package envoyconfig

import (
	"testing"

	envoy_serializer_v3 "github.com/3scale-ops/marin3r/pkg/envoy/serializer/v3"
	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/MakeNowJust/heredoc"
	"sigs.k8s.io/yaml"
)

func TestRuntime_v1(t *testing.T) {
	type args struct {
		name string
		opts *saasv1alpha1.Runtime
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Generates runtime",
			args: args{
				name: "runtime",
				opts: &saasv1alpha1.Runtime{ListenerNames: []string{"listener1", "listener2"}},
			},
			want: heredoc.Doc(`
                layer:
                  envoy:
                    resource_limits:
                      listener:
                        listener1:
                          connection_limit: 10000
                        listener2:
                          connection_limit: 10000
                  overload:
                    global_downstream_max_connections: 50000
                name: runtime
			`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Runtime_v1(tt.args.name, tt.args.opts)
			j, err := envoy_serializer_v3.JSON{}.Marshal(got)
			if err != nil {
				t.Error(err)
			}
			y, err := yaml.JSONToYAML([]byte(j))
			if err != nil {
				t.Error(err)
			}
			if string(y) != tt.want {
				t.Errorf("Runtime_v1():\n# got:\n%v\n# want:\n%v", string(y), tt.want)
			}
		})
	}
}
