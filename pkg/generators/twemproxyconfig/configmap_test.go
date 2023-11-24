package twemproxyconfig

import (
	"testing"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators"
	"github.com/3scale/saas-operator/pkg/resource_builders/twemproxy"
	"github.com/go-test/deep"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGenerator_configMap(t *testing.T) {
	type fields struct {
		BaseOptionsV2  generators.BaseOptionsV2
		Spec           saasv1alpha1.TwemproxyConfigSpec
		masterTargets  map[string]twemproxy.Server
		slaverwTargets map[string]twemproxy.Server
	}
	type args struct {
		toYAML bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *corev1.ConfigMap
	}{
		{
			name: "Generates the Twemproxy ConfigMap using masters",
			fields: fields{
				BaseOptionsV2: generators.BaseOptionsV2{
					Component:    "twemproxy",
					InstanceName: "test",
					Namespace:    "ns",
					Labels:       map[string]string{},
				},
				Spec: saasv1alpha1.TwemproxyConfigSpec{
					SentinelURIs: []string{"sentinel.example.com"},
					ServerPools: []saasv1alpha1.TwemproxyServerPool{
						{
							Name:   "pool1",
							Target: func() *saasv1alpha1.TargetRedisServers { t := saasv1alpha1.Masters; return &t }(),
							Topology: []saasv1alpha1.ShardedRedisTopology{
								{ShardName: "lshard01", PhysicalShard: "pshard01"},
								{ShardName: "lshard02", PhysicalShard: "pshard01"},
								{ShardName: "lshard03", PhysicalShard: "pshard01"},
								{ShardName: "lshard04", PhysicalShard: "pshard02"},
							},
							BindAddress: "localhost:2000",
							Timeout:     1000,
							TCPBacklog:  500,
							PreConnect:  false,
						},
						{
							Name:   "pool2",
							Target: func() *saasv1alpha1.TargetRedisServers { t := saasv1alpha1.Masters; return &t }(),
							Topology: []saasv1alpha1.ShardedRedisTopology{
								{ShardName: "lshard01", PhysicalShard: "pshard01"},
								{ShardName: "lshard02", PhysicalShard: "pshard02"},
							},
							BindAddress: "localhost:3000",
							Timeout:     1000,
							TCPBacklog:  500,
							PreConnect:  false,
						},
					},
				},
				masterTargets: map[string]twemproxy.Server{
					"pshard01": {Address: "127.0.0.1:6379", Priority: 1, Name: "pshard01"},
					"pshard02": {Address: "127.0.0.2:6379", Priority: 1, Name: "pshard02"},
				},
				slaverwTargets: map[string]twemproxy.Server{},
			},
			args: args{},
			want: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "ns",
					Labels:    map[string]string{},
				},
				Data: map[string]string{
					"nutcracker.yml": `{"health":{"listen":"127.0.0.1:22333","preconnect":false,"redis":true,"auto_eject_hosts":false,"servers":["127.0.0.1:6379:1 dummy"]},"pool1":{"listen":"localhost:2000","hash":"fnv1a_64","hash_tag":"{}","distribution":"ketama","timeout":1000,"backlog":500,"preconnect":false,"redis":true,"auto_eject_hosts":false,"servers":["127.0.0.1:6379:1 lshard01","127.0.0.1:6379:1 lshard02","127.0.0.1:6379:1 lshard03","127.0.0.2:6379:1 lshard04"]},"pool2":{"listen":"localhost:3000","hash":"fnv1a_64","hash_tag":"{}","distribution":"ketama","timeout":1000,"backlog":500,"preconnect":false,"redis":true,"auto_eject_hosts":false,"servers":["127.0.0.1:6379:1 lshard01","127.0.0.2:6379:1 lshard02"]}}`,
				},
			},
		},
		{
			name: "Generates the Twemproxy ConfigMap using rw slaves",
			fields: fields{
				BaseOptionsV2: generators.BaseOptionsV2{
					Component:    "twemproxy",
					InstanceName: "test",
					Namespace:    "ns",
					Labels:       map[string]string{},
				},
				Spec: saasv1alpha1.TwemproxyConfigSpec{
					SentinelURIs: []string{"sentinel.example.com"},
					ServerPools: []saasv1alpha1.TwemproxyServerPool{
						{
							Name:   "pool1",
							Target: func() *saasv1alpha1.TargetRedisServers { t := saasv1alpha1.SlavesRW; return &t }(),
							Topology: []saasv1alpha1.ShardedRedisTopology{
								{ShardName: "lshard01", PhysicalShard: "pshard01"},
								{ShardName: "lshard02", PhysicalShard: "pshard01"},
								{ShardName: "lshard03", PhysicalShard: "pshard01"},
								{ShardName: "lshard04", PhysicalShard: "pshard02"},
							},
							BindAddress: "localhost:2000",
							Timeout:     1000,
							TCPBacklog:  500,
							PreConnect:  false,
						},
					},
				},
				masterTargets: map[string]twemproxy.Server{
					"pshard01": {Address: "127.0.0.1:6379", Priority: 1, Name: "pshard01"},
					"pshard02": {Address: "127.0.0.2:6379", Priority: 1, Name: "pshard02"},
				},
				slaverwTargets: map[string]twemproxy.Server{
					"pshard01": {Address: "127.0.0.3:6379", Priority: 1, Name: "pshard01"},
					"pshard02": {Address: "127.0.0.4:6379", Priority: 1, Name: "pshard02"},
				},
			},
			args: args{},
			want: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "ns",
					Labels:    map[string]string{},
				},
				Data: map[string]string{
					"nutcracker.yml": `{"health":{"listen":"127.0.0.1:22333","preconnect":false,"redis":true,"auto_eject_hosts":false,"servers":["127.0.0.1:6379:1 dummy"]},"pool1":{"listen":"localhost:2000","hash":"fnv1a_64","hash_tag":"{}","distribution":"ketama","timeout":1000,"backlog":500,"preconnect":false,"redis":true,"auto_eject_hosts":false,"servers":["127.0.0.3:6379:1 lshard01","127.0.0.3:6379:1 lshard02","127.0.0.3:6379:1 lshard03","127.0.0.4:6379:1 lshard04"]}}`,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := &Generator{
				BaseOptionsV2:  tt.fields.BaseOptionsV2,
				Spec:           tt.fields.Spec,
				masterTargets:  tt.fields.masterTargets,
				slaverwTargets: tt.fields.slaverwTargets,
			}
			got := gen.configMap(tt.args.toYAML)
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("Generator.configMap() = diff %v", diff)
			}
		})
	}
}
