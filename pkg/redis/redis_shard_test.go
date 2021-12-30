package redis

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"testing"

	"github.com/3scale/saas-operator/pkg/redis/crud"
	"github.com/3scale/saas-operator/pkg/redis/crud/client"
	"github.com/go-logr/logr"
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
				Name:     "test",
				IP:       "127.0.0.1",
				Port:     "3333",
				Role:     client.Unknown,
				ReadOnly: false,
				CRUD:     func() *crud.CRUD { c, _ := crud.NewRedisCRUDFromConnectionString("redis://127.0.0.1:3333"); return c }(),
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
		{
			name: "'config get' command fails, returns an error",
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
							return []interface{}{}
						},
						InjectError: func() error { return errors.New("error") },
					},
				),
			},
			args:         args{ctx: context.TODO()},
			wantRole:     client.Slave,
			wantReadOnly: false,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &RedisServer{
				Name:     tt.fields.Name,
				IP:       tt.fields.IP,
				Port:     tt.fields.Port,
				Role:     tt.fields.Role,
				ReadOnly: tt.fields.ReadOnly,
				CRUD:     tt.fields.CRUD,
			}
			if err := srv.Discover(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("RedisServer.Discover() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantRole != srv.Role || tt.wantReadOnly != srv.ReadOnly {
				t.Errorf("RedisServer.Discover() got = %v/%v, want %v/%v", srv.Role, srv.ReadOnly, tt.wantRole, tt.wantReadOnly)
			}
		})
	}
}

func TestNewShard(t *testing.T) {
	type args struct {
		name              string
		connectionStrings []string
	}
	tests := []struct {
		name    string
		args    args
		want    *Shard
		wantErr bool
	}{
		{
			name: "Returns a new Shard object",
			args: args{
				name:              "test",
				connectionStrings: []string{"redis://127.0.0.1:1000", "redis://127.0.0.1:2000", "redis://127.0.0.1:3000"},
			},
			want: &Shard{
				Name: "test",
				Servers: []RedisServer{
					{
						Name:     "redis://127.0.0.1:1000",
						IP:       "127.0.0.1",
						Port:     "1000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD:     func() *crud.CRUD { c, _ := crud.NewRedisCRUDFromConnectionString("redis://127.0.0.1:1000"); return c }(),
					},
					{
						Name:     "redis://127.0.0.1:2000",
						IP:       "127.0.0.1",
						Port:     "2000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD:     func() *crud.CRUD { c, _ := crud.NewRedisCRUDFromConnectionString("redis://127.0.0.1:2000"); return c }(),
					},
					{
						Name:     "redis://127.0.0.1:3000",
						IP:       "127.0.0.1",
						Port:     "3000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD:     func() *crud.CRUD { c, _ := crud.NewRedisCRUDFromConnectionString("redis://127.0.0.1:3000"); return c }(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Returns an error (bad connection string)",
			args: args{
				name:              "test",
				connectionStrings: []string{"redis://127.0.0.1:1000", "127.0.0.1:2000", "redis://127.0.0.1:3000"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewShard(tt.args.name, tt.args.connectionStrings)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewShard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("NewShard() got diff: %v", diff)
			}
		})
	}
}

func TestShard_Discover(t *testing.T) {
	type fields struct {
		Name    string
		Servers []RedisServer
	}
	type args struct {
		ctx context.Context
		log logr.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Discovers characteristics for all servers in the shard",
			fields: fields{
				Name: "test",
				Servers: []RedisServer{
					{
						Name:     "redis://127.0.0.1:1000",
						IP:       "127.0.0.1",
						Port:     "1000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD: crud.NewFakeCRUD(
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"master"}
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
					{
						Name:     "redis://127.0.0.1:2000",
						IP:       "127.0.0.1",
						Port:     "2000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD: crud.NewFakeCRUD(
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"slave", "127.0.0.1"}
								},
								InjectError: func() error { return nil },
							},
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"read-only", "yes"}
								},
								InjectError: func() error { return nil },
							},
						)},
					{
						Name:     "redis://127.0.0.1:3000",
						IP:       "127.0.0.1",
						Port:     "3000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD: crud.NewFakeCRUD(
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"slave", "127.0.0.1"}
								},
								InjectError: func() error { return nil },
							},
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"read-only", "yes"}
								},
								InjectError: func() error { return nil },
							},
						)},
				},
			},
			args:    args{ctx: context.TODO(), log: logr.Discard()},
			wantErr: false,
		},
		{
			name: "second server fails, returns error",
			fields: fields{
				Name: "test",
				Servers: []RedisServer{
					{
						Name:     "redis://127.0.0.1:1000",
						IP:       "127.0.0.1",
						Port:     "1000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD: crud.NewFakeCRUD(
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"master"}
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
					{
						Name:     "redis://127.0.0.1:2000",
						IP:       "127.0.0.1",
						Port:     "2000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD: crud.NewFakeCRUD(
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{}
								},
								InjectError: func() error { return errors.New("error") },
							},
						)},
					{
						Name:     "redis://127.0.0.1:3000",
						IP:       "127.0.0.1",
						Port:     "3000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD: crud.NewFakeCRUD(
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"slave", "127.0.0.1"}
								},
								InjectError: func() error { return nil },
							},
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"read-only", "yes"}
								},
								InjectError: func() error { return nil },
							},
						)},
				},
			},
			args:    args{ctx: context.TODO(), log: logr.Discard()},
			wantErr: true,
		},
		{
			name: "no master, returns error",
			fields: fields{
				Name: "test",
				Servers: []RedisServer{
					{
						Name:     "redis://127.0.0.1:1000",
						IP:       "127.0.0.1",
						Port:     "1000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD: crud.NewFakeCRUD(
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"slave", "no one"}
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
					{
						Name:     "redis://127.0.0.1:2000",
						IP:       "127.0.0.1",
						Port:     "2000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD: crud.NewFakeCRUD(
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"slave", "no one"}
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
					{
						Name:     "redis://127.0.0.1:3000",
						IP:       "127.0.0.1",
						Port:     "3000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD: crud.NewFakeCRUD(
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"slave", "no one"}
								},
								InjectError: func() error { return nil },
							},
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"read-only", "yes"}
								},
								InjectError: func() error { return nil },
							},
						)},
				},
			},
			args:    args{ctx: context.TODO(), log: logr.Discard()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shard{
				Name:    tt.fields.Name,
				Servers: tt.fields.Servers,
			}
			if err := s.Discover(tt.args.ctx, tt.args.log); (err != nil) != tt.wantErr {
				t.Errorf("Shard.Discover() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShard_Init(t *testing.T) {
	type fields struct {
		Name    string
		Servers []RedisServer
	}
	type args struct {
		ctx         context.Context
		masterIndex int32
		log         logr.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				Name: "All redis servers configured",
				Servers: []RedisServer{
					{
						Name:     "redis://127.0.0.1:1000",
						IP:       "127.0.0.1",
						Port:     "1000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD: crud.NewFakeCRUD(
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"slave", "127.0.0.1"}
								},
								InjectError: func() error { return nil },
							},
							client.FakeResponse{
								InjectResponse: func() interface{} { return nil },
								InjectError:    func() error { return nil },
							},
						),
					},
					{
						Name:     "redis://127.0.0.1:2000",
						IP:       "127.0.0.1",
						Port:     "2000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD: crud.NewFakeCRUD(
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"slave", "127.0.0.1"}
								},
								InjectError: func() error { return nil },
							},
							client.FakeResponse{
								InjectResponse: func() interface{} { return nil },
								InjectError:    func() error { return nil },
							},
						)},
					{
						Name:     "redis://127.0.0.1:3000",
						IP:       "127.0.0.1",
						Port:     "3000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD: crud.NewFakeCRUD(
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"slave", "127.0.0.1"}
								},
								InjectError: func() error { return nil },
							},
							client.FakeResponse{
								InjectResponse: func() interface{} { return nil },
								InjectError:    func() error { return nil },
							},
						)},
				},
			},
			args:    args{ctx: context.TODO(), masterIndex: 0, log: logr.Discard()},
			want:    []string{"redis://127.0.0.1:1000", "redis://127.0.0.1:2000", "redis://127.0.0.1:3000"},
			wantErr: false,
		},
		{
			name: "No configuration needed",
			fields: fields{
				Name: "All redis servers configured",
				Servers: []RedisServer{
					{
						Name:     "redis://127.0.0.1:1000",
						IP:       "127.0.0.1",
						Port:     "1000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD: crud.NewFakeCRUD(
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"master"}
								},
								InjectError: func() error { return nil },
							},
							client.FakeResponse{
								InjectResponse: func() interface{} { return nil },
								InjectError:    func() error { return nil },
							},
						),
					},
					{
						Name:     "redis://127.0.0.1:2000",
						IP:       "127.0.0.1",
						Port:     "2000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD: crud.NewFakeCRUD(
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"slave", "10.0.0.1"}
								},
								InjectError: func() error { return nil },
							},
							client.FakeResponse{
								InjectResponse: func() interface{} { return nil },
								InjectError:    func() error { return nil },
							},
						)},
					{
						Name:     "redis://127.0.0.1:3000",
						IP:       "127.0.0.1",
						Port:     "3000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD: crud.NewFakeCRUD(
							client.FakeResponse{
								InjectResponse: func() interface{} {
									return []interface{}{"slave", "10.0.0.1"}
								},
								InjectError: func() error { return nil },
							},
							client.FakeResponse{
								InjectResponse: func() interface{} { return nil },
								InjectError:    func() error { return nil },
							},
						)},
				},
			},
			args:    args{ctx: context.TODO(), masterIndex: 0, log: logr.Discard()},
			want:    []string{},
			wantErr: false,
		},
		{
			name: "Returns error",
			fields: fields{
				Name: "All redis servers configured",
				Servers: []RedisServer{
					{
						Name:     "redis://127.0.0.1:1000",
						IP:       "127.0.0.1",
						Port:     "1000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD: crud.NewFakeCRUD(
							client.FakeResponse{
								InjectResponse: func() interface{} { return []interface{}{} },
								InjectError:    func() error { return errors.New("error") },
							},
						),
					},
				},
			},
			args:    args{ctx: context.TODO(), masterIndex: 0, log: logr.Discard()},
			want:    []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shard{
				Name:    tt.fields.Name,
				Servers: tt.fields.Servers,
			}
			got, err := s.Init(tt.args.ctx, tt.args.masterIndex, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("Shard.Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Shard.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewShardedCluster(t *testing.T) {
	type args struct {
		ctx        context.Context
		serverList map[string][]string
		log        logr.Logger
	}
	tests := []struct {
		name    string
		args    args
		want    ShardedCluster
		wantErr bool
	}{
		{
			name: "Returns a new ShardedCluster object",
			args: args{
				ctx: context.TODO(),
				serverList: map[string][]string{
					"shard00": {"redis://127.0.0.1:1000", "redis://127.0.0.1:2000"},
					"shard01": {"redis://127.0.0.1:3000", "redis://127.0.0.1:4000"},
				},
				log: logr.Discard(),
			},
			want: ShardedCluster{
				{
					Name: "shard00",
					Servers: []RedisServer{
						{
							Name:     "redis://127.0.0.1:1000",
							IP:       "127.0.0.1",
							Port:     "1000",
							Role:     client.Unknown,
							ReadOnly: false,
							CRUD:     func() *crud.CRUD { c, _ := crud.NewRedisCRUDFromConnectionString("redis://127.0.0.1:1000"); return c }(),
						},
						{
							Name:     "redis://127.0.0.1:2000",
							IP:       "127.0.0.1",
							Port:     "2000",
							Role:     client.Unknown,
							ReadOnly: false,
							CRUD:     func() *crud.CRUD { c, _ := crud.NewRedisCRUDFromConnectionString("redis://127.0.0.1:2000"); return c }(),
						},
					},
				},
				{
					Name: "shard01",
					Servers: []RedisServer{
						{
							Name:     "redis://127.0.0.1:3000",
							IP:       "127.0.0.1",
							Port:     "3000",
							Role:     client.Unknown,
							ReadOnly: false,
							CRUD:     func() *crud.CRUD { c, _ := crud.NewRedisCRUDFromConnectionString("redis://127.0.0.1:3000"); return c }(),
						},
						{
							Name:     "redis://127.0.0.1:4000",
							IP:       "127.0.0.1",
							Port:     "4000",
							Role:     client.Unknown,
							ReadOnly: false,
							CRUD:     func() *crud.CRUD { c, _ := crud.NewRedisCRUDFromConnectionString("redis://127.0.0.1:4000"); return c }(),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Returns error",
			args: args{
				ctx: context.TODO(),
				serverList: map[string][]string{
					"shard00": {"redis://127.0.0.1:1000", "redis://127.0.0.1:2000"},
					"shard01": {"127.0.0.1:3000", "redis://127.0.0.1:4000"},
				},
				log: logr.Discard(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewShardedCluster(tt.args.ctx, tt.args.serverList, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewShardedCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			sort.SliceStable(got, func(i, j int) bool {
				return got[i].Name < got[j].Name
			})
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("NewShardedCluster() got diff: %v", diff)
			}
		})
	}
}

func TestShardedCluster_Discover(t *testing.T) {
	type args struct {
		ctx context.Context
		log logr.Logger
	}
	tests := []struct {
		name    string
		sc      ShardedCluster
		args    args
		wantErr bool
	}{
		{
			name: "Discovers characteristics of all servers in the ShardedCluster",
			sc: ShardedCluster{
				{
					Name: "shard00",
					Servers: []RedisServer{
						{
							Name:     "redis://127.0.0.1:1000",
							IP:       "127.0.0.1",
							Port:     "1000",
							Role:     client.Unknown,
							ReadOnly: false,
							CRUD: crud.NewFakeCRUD(
								client.FakeResponse{
									InjectResponse: func() interface{} {
										return []interface{}{"master"}
									},
									InjectError: func() error { return nil },
								},
							)},
					},
				},
				{
					Name: "shard01",
					Servers: []RedisServer{
						{
							Name:     "redis://127.0.0.1:3000",
							IP:       "127.0.0.1",
							Port:     "3000",
							Role:     client.Unknown,
							ReadOnly: false,
							CRUD: crud.NewFakeCRUD(
								client.FakeResponse{
									InjectResponse: func() interface{} {
										return []interface{}{"master"}
									},
									InjectError: func() error { return nil },
								},
							)},
					},
				},
			},
			args:    args{ctx: context.TODO(), log: logr.Discard()},
			wantErr: false,
		},
		{
			name: "Returns error",
			sc: ShardedCluster{
				{
					Name: "shard00",
					Servers: []RedisServer{
						{
							Name:     "redis://127.0.0.1:1000",
							IP:       "127.0.0.1",
							Port:     "1000",
							Role:     client.Unknown,
							ReadOnly: false,
							CRUD: crud.NewFakeCRUD(
								client.FakeResponse{
									InjectResponse: func() interface{} { return []interface{}{} },
									InjectError:    func() error { return errors.New("error") },
								},
							)},
					},
				},
				{
					Name: "shard01",
					Servers: []RedisServer{
						{
							Name:     "redis://127.0.0.1:3000",
							IP:       "127.0.0.1",
							Port:     "3000",
							Role:     client.Unknown,
							ReadOnly: false,
							CRUD:     crud.NewFakeCRUD()},
					},
				},
			},
			args:    args{ctx: context.TODO(), log: logr.Discard()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sc.Discover(tt.args.ctx, tt.args.log); (err != nil) != tt.wantErr {
				t.Errorf("ShardedCluster.Discover() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShardedCluster_GetShardNames(t *testing.T) {
	tests := []struct {
		name string
		sc   ShardedCluster
		want []string
	}{
		{
			name: "Returns the shrard names as a slice of strings",
			sc: ShardedCluster{
				{
					Name:    "shard00",
					Servers: []RedisServer{},
				},
				{
					Name:    "shard01",
					Servers: []RedisServer{},
				},
				{
					Name:    "shard02",
					Servers: []RedisServer{},
				},
			},
			want: []string{"shard00", "shard01", "shard02"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sc.GetShardNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShardedCluster.GetShardNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShardedCluster_GetShardByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		sc   ShardedCluster
		args args
		want *Shard
	}{
		{
			name: "Returns the shard of the given name",
			sc: ShardedCluster{
				{
					Name: "shard00",
					Servers: []RedisServer{
						{
							Name:     "redis://127.0.0.1:1000",
							IP:       "127.0.0.1",
							Port:     "1000",
							Role:     client.Unknown,
							ReadOnly: false,
							CRUD:     func() *crud.CRUD { c, _ := crud.NewRedisCRUDFromConnectionString("redis://127.0.0.1:1000"); return c }(),
						},
					},
				},
				{
					Name: "shard01",
					Servers: []RedisServer{
						{
							Name:     "redis://127.0.0.1:2000",
							IP:       "127.0.0.1",
							Port:     "3000",
							Role:     client.Unknown,
							ReadOnly: false,
							CRUD:     func() *crud.CRUD { c, _ := crud.NewRedisCRUDFromConnectionString("redis://127.0.0.1:2000"); return c }(),
						},
					},
				},
			},
			args: args{
				name: "shard01",
			},
			want: &Shard{
				Name: "shard01",
				Servers: []RedisServer{
					{
						Name:     "redis://127.0.0.1:2000",
						IP:       "127.0.0.1",
						Port:     "3000",
						Role:     client.Unknown,
						ReadOnly: false,
						CRUD:     func() *crud.CRUD { c, _ := crud.NewRedisCRUDFromConnectionString("redis://127.0.0.1:2000"); return c }(),
					},
				},
			},
		},
		{
			name: "Returns nil if not found",
			sc:   ShardedCluster{},
			args: args{
				name: "shard01",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := deep.Equal(tt.sc.GetShardByName(tt.args.name), tt.want); len(diff) > 0 {
				t.Errorf("ShardedCluster.GetShardByName() got diff: %v", diff)
			}
		})
	}
}
