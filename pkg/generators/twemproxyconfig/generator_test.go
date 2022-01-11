package twemproxyconfig

import (
	"reflect"
	"testing"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
)

func Test_applyQuorum(t *testing.T) {
	type args struct {
		responses []saasv1alpha1.MonitoredShards
		quorum    int
	}
	tests := []struct {
		name    string
		args    args
		want    saasv1alpha1.MonitoredShards
		wantErr bool
	}{
		{
			name: "Should return the accepted response",
			args: args{
				responses: []saasv1alpha1.MonitoredShards{
					{
						{Name: "shard01", Master: "127.0.0.1:1111"},
						{Name: "shard02", Master: "127.0.0.2:2222"},
						{Name: "shard03", Master: "127.0.0.3:3333"},
					},
					{
						{Name: "shard03", Master: "127.0.0.3:3333"},
						{Name: "shard02", Master: "127.0.0.2:2222"},
						{Name: "shard01", Master: "127.0.0.1:1111"},
					},
				},
				quorum: 2,
			},
			want: []saasv1alpha1.MonitoredShard{
				{Name: "shard01", Master: "127.0.0.1:1111"},
				{Name: "shard02", Master: "127.0.0.2:2222"},
				{Name: "shard03", Master: "127.0.0.3:3333"},
			},
			wantErr: false,
		},
		{
			name: "Should fail, no quorum",
			args: args{
				responses: []saasv1alpha1.MonitoredShards{
					{
						{Name: "shard01", Master: "127.0.0.1:1111"},
						{Name: "shard02", Master: "127.0.0.2:2222"},
						{Name: "shard03", Master: "127.0.0.3:3333"},
					},
					{
						{Name: "shard03", Master: "127.0.0.3:3333"},
						{Name: "shard02", Master: "127.0.0.2:2222"},
						{Name: "shard01", Master: "127.0.0.4:4444"},
					},
				},
				quorum: 2,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := applyQuorum(tt.args.responses, tt.args.quorum)
			if (err != nil) != tt.wantErr {
				t.Errorf("applyQuorum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("applyQuorum() = %v, want %v", got, tt.want)
			}
		})
	}
}
