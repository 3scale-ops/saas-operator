package server

import (
	"reflect"
	"testing"

	"github.com/3scale/saas-operator/pkg/util"
	"github.com/davecgh/go-spew/spew"
)

func TestServerPool_GetServer(t *testing.T) {
	type fields struct {
		servers []*Server
	}
	type args struct {
		connectionString string
		alias            *string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Server
		wantErr bool
	}{
		{
			name: "Gets the server by alias",
			fields: fields{
				servers: []*Server{
					{Alias: "host1", Client: nil, Host: "127.0.0.1", Port: "1000"},
					{Alias: "host2", Client: nil, Host: "127.0.0.2", Port: "2000"},
				}},
			args: args{
				connectionString: "redis://127.0.0.2:2000",
				alias:            util.Pointer("host2"),
			},
			want:    &Server{Alias: "host2", Client: nil, Host: "127.0.0.2", Port: "2000"},
			wantErr: false,
		},
		{
			name: "Gets the server by host",
			fields: fields{
				servers: []*Server{
					{Alias: "host1", Client: nil, Host: "127.0.0.1", Port: "1000"},
					{Alias: "host2", Client: nil, Host: "127.0.0.2", Port: "2000"},
				}},
			args: args{
				connectionString: "redis://127.0.0.2:2000",
			},
			want:    &Server{Alias: "host2", Client: nil, Host: "127.0.0.2", Port: "2000"},
			wantErr: false,
		},
		{
			name: "Gets the server by host and sets the alias",
			fields: fields{
				servers: []*Server{
					{Alias: "host1", Client: nil, Host: "127.0.0.1", Port: "1000"},
					{Alias: "", Client: nil, Host: "127.0.0.2", Port: "2000"},
				}},
			args: args{
				connectionString: "redis://127.0.0.2:2000",
				alias:            util.Pointer("host2"),
			},
			want:    &Server{Alias: "host2", Client: nil, Host: "127.0.0.2", Port: "2000"},
			wantErr: false,
		},
		{
			name: "Returns error",
			fields: fields{
				servers: []*Server{}},
			args: args{
				connectionString: "host",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := &ServerPool{
				servers: tt.fields.servers,
			}
			got, err := pool.GetServer(tt.args.connectionString, tt.args.alias)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServerPool.GetServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				spew.Dump(got)
				spew.Dump(tt.want)
				t.Errorf("ServerPool.GetServer() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("Adds a new server to the pool", func(t *testing.T) {
		pool := &ServerPool{
			servers: []*Server{{Alias: "host1", Client: nil, Host: "127.0.0.1", Port: "1000"}},
		}
		new, _ := pool.GetServer("redis://127.0.0.2:2000", util.Pointer("host2"))
		exists, _ := pool.GetServer("redis://127.0.0.2:2000", util.Pointer("host2"))
		if new != exists {
			t.Errorf("ServerPool.GetServer() = %v, want %v", new, exists)
		}
	})
}

func TestServerPool_indexByAlias(t *testing.T) {
	type fields struct {
		servers []*Server
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]*Server
	}{
		{
			name: "Returns a map indexed by host",
			fields: fields{
				servers: []*Server{
					{Alias: "host1", Client: nil, Host: "127.0.0.1", Port: "1000"},
					{Alias: "host2", Client: nil, Host: "127.0.0.2", Port: "2000"},
				},
			},
			want: map[string]*Server{
				"host1": {Alias: "host1", Client: nil, Host: "127.0.0.1", Port: "1000"},
				"host2": {Alias: "host2", Client: nil, Host: "127.0.0.2", Port: "2000"},
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := &ServerPool{
				servers: tt.fields.servers,
			}
			if got := pool.indexByAlias(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServerPool.indexByAlias() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerPool_indexByHost(t *testing.T) {
	type fields struct {
		servers []*Server
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]*Server
	}{
		{
			name: "Returns a map indexed by host",
			fields: fields{
				servers: []*Server{
					{Alias: "host1", Client: nil, Host: "127.0.0.1", Port: "1000"},
					{Alias: "host2", Client: nil, Host: "127.0.0.2", Port: "2000"},
				},
			},
			want: map[string]*Server{
				"127.0.0.1:1000": {Alias: "host1", Client: nil, Host: "127.0.0.1", Port: "1000"},
				"127.0.0.2:2000": {Alias: "host2", Client: nil, Host: "127.0.0.2", Port: "2000"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := &ServerPool{
				servers: tt.fields.servers,
			}
			if got := pool.indexByHost(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServerPool.indexByHost() = %v, want %v", got, tt.want)
			}
		})
	}
}
