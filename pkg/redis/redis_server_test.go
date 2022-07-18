package redis

import (
	"context"
	"errors"
	"testing"

	"github.com/3scale/saas-operator/pkg/redis/crud"
	"github.com/3scale/saas-operator/pkg/redis/crud/client"
	"github.com/go-test/deep"
)

func TestNewRedisServerFromConnectionString(t *testing.T) {
	type args struct {
		name             string
		connectionString string
	}
	tests := []struct {
		name    string
		args    args
		want    *RedisServer
		wantErr bool
	}{
		{
			name: "Returns a new RedisServer object for the given connection string",
			args: args{
				name:             "test",
				connectionString: "redis://127.0.0.1:3333",
			},
			want: &RedisServer{
				Name: "test",
				Host: "127.0.0.1",
				Port: "3333",
				Role: client.Unknown,
				CRUD: func() *crud.CRUD { c, _ := crud.NewRedisCRUDFromConnectionString("redis://127.0.0.1:3333"); return c }(),
			},
			wantErr: false,
		},
		{
			name: "Returns error",
			args: args{
				name:             "test",
				connectionString: "127.0.0.1:3333",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRedisServerFromConnectionString(tt.args.name, tt.args.connectionString)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRedisServerFromConnectionString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("NewSentinelServer() got diff: %v", diff)
			}
		})
	}
}

func TestRedisServer_Discover(t *testing.T) {
	type fields struct {
		Name     string
		IP       string
		Port     string
		Role     client.Role
		ReadOnly bool
		CRUD     *crud.CRUD
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantRole     client.Role
		wantReadOnly bool
		wantErr      bool
	}{
		{
			name: "Discovers characteristics of the redis server: master/rw",
			fields: fields{
				CRUD: crud.NewFakeCRUD(
					client.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{"master", ""}
						},
						InjectError: func() error { return nil },
					},
				),
			},
			args:         args{ctx: context.TODO()},
			wantRole:     client.Master,
			wantReadOnly: false,
			wantErr:      false,
		},
		{
			name: "Discovers characteristics of the redis server: slave/ro",
			fields: fields{
				CRUD: crud.NewFakeCRUD(
					client.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{"slave", "127.0.0.1:3333"}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{"read-only", "yes"}
						},
						InjectError: func() error { return nil },
					},
				),
			},
			args:         args{ctx: context.TODO()},
			wantRole:     client.Slave,
			wantReadOnly: true,
			wantErr:      false,
		},
		{
			name: "Discovers characteristics of the redis server: slave/rw",
			fields: fields{
				CRUD: crud.NewFakeCRUD(
					client.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{"slave", "127.0.0.1:3333"}
						},
						InjectError: func() error { return nil },
					},
					client.FakeResponse{
						InjectResponse: func() interface{} {
							return []interface{}{"read-only", "no"}
						},
						InjectError: func() error { return nil },
					},
				),
			},
			args:         args{ctx: context.TODO()},
			wantRole:     client.Slave,
			wantReadOnly: false,
			wantErr:      false,
		},
		{
			name: "'role' command fails, returns an error",
			fields: fields{
				CRUD: crud.NewFakeCRUD(
					client.FakeResponse{
						InjectResponse: func() interface{} { return []interface{}{} },
						InjectError:    func() error { return errors.New("error") },
					},
				),
			},
			args:         args{ctx: context.TODO()},
			wantRole:     client.Unknown,
			wantReadOnly: false,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &RedisServer{
				Name: tt.fields.Name,
				Host: tt.fields.IP,
				Port: tt.fields.Port,
				Role: tt.fields.Role,
				CRUD: tt.fields.CRUD,
			}
			if err := srv.Discover(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("RedisServer.Discover() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantRole != srv.Role {
				t.Errorf("RedisServer.Discover() got = %v, want %v", srv.Role, tt.wantRole)
			}
		})
	}
}
