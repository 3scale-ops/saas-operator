package sharded

import (
	"testing"

	"github.com/3scale/saas-operator/pkg/redis/client"
	redis "github.com/3scale/saas-operator/pkg/redis/server"
	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-test/deep"
)

func init() {
	deep.CompareUnexportedFields = true
}

func TestNewRedisServerFromParams(t *testing.T) {
	type args struct {
		connectionString string
		alias            *string
		pool             *redis.ServerPool
	}
	tests := []struct {
		name    string
		args    args
		want    *RedisServer
		wantErr bool
	}{
		{
			name: "Retuns a RedisServer",
			args: args{
				connectionString: "redis://127.0.0.1:1000",
				alias:            util.Pointer("host1"),
				pool:             redis.NewServerPool(redis.NewServerFromParams("host1", "127.0.0.1", "1000", client.MustNewFromConnectionString("redis://127.0.0.1:1000"))),
			},
			want: &RedisServer{
				Server: redis.NewServerFromParams("host1", "127.0.0.1", "1000", client.MustNewFromConnectionString("redis://127.0.0.1:1000")),
				Role:   client.Unknown,
				Config: map[string]string{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRedisServerFromPool(tt.args.connectionString, tt.args.alias, tt.args.pool)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServerFromParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("NewServerFromParams() = got diff %v", diff)
			}
		})
	}
}
