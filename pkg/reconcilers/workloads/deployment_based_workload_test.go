package workloads

import (
	"reflect"
	"testing"
)

func Test_unwrapNil(t *testing.T) {
	type args struct {
		w DeploymentWorkload
	}
	tests := []struct {
		name string
		args args
		want DeploymentWorkload
	}{
		{
			name: "Detects a nil value",
			args: args{
				w: nil,
			},
			want: nil,
		},
		{
			name: "Detects interface containing nil value",
			args: args{
				w: func() DeploymentWorkload {
					val := (*TestWorkloadGenerator)(nil)
					return val
				}(),
			},
			want: nil,
		},
		{
			name: "Lets an interface containing a non nil value pass through",
			args: args{
				w: &TestWorkloadGenerator{},
			},
			want: &TestWorkloadGenerator{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := unwrapNil(tt.args.w); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unwrapNil() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_toWithTraffic(t *testing.T) {
	type args struct {
		w DeploymentWorkload
	}
	tests := []struct {
		name string
		args args
		want WithTraffic
	}{
		{
			name: "Detects a nil value",
			args: args{
				w: nil,
			},
			want: nil,
		},
		{
			name: "Detects interface containing nil value",
			args: args{
				w: func() DeploymentWorkload {
					val := (*TestWorkloadGenerator)(nil)
					return val
				}(),
			},
			want: nil,
		},
		{
			name: "Converts DeploymentWorkload to WithTraffic",
			args: args{
				w: &TestWorkloadGenerator{},
			},
			want: &TestWorkloadGenerator{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toWithTraffic(tt.args.w); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toWithTraffic() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
