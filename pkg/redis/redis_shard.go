package redis

import (
	"context"
	"fmt"
	"net"
	"sort"

	"github.com/3scale/saas-operator/pkg/redis/crud"
	"github.com/3scale/saas-operator/pkg/redis/crud/client"
	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-logr/logr"
)

// RedisServer represent a redis server and its characteristics
type RedisServer struct {
	Name     string
	Host     string
	Port     string
	Role     client.Role
	ReadOnly bool
	CRUD     *crud.CRUD
}

func (rs *RedisServer) IP() (string, error) {
	var ip string
	if r := net.ParseIP(rs.Host); r != nil {
		ip = r.String()
	} else {
		// if it is not an IP, try to resolve a DNS
		ips, err := net.LookupIP(rs.Host)
		if err != nil {
			return "", err
		}
		if len(ips) > 1 {
			return "", fmt.Errorf("dns resolves to more than 1 IP")
		}
		ip = ips[0].String()
	}

	return ip, nil
}

func NewRedisServerFromConnectionString(name, connectionString string) (*RedisServer, error) {

	crud, err := crud.NewRedisCRUDFromConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	return &RedisServer{Name: name, Host: crud.GetIP(), Port: crud.GetPort(), CRUD: crud, ReadOnly: false, Role: client.Unknown}, nil
}

// Discover returns the Role and the IsReadOnly flag for a given
// redis Server
func (srv *RedisServer) Discover(ctx context.Context) error {

	role, _, err := srv.CRUD.RedisRole(ctx)
	if err != nil {
		srv.Role = client.Unknown
		return util.WrapError("redis-autodiscovery", err)
	}
	srv.Role = role

	if srv.Role == client.Slave {
		ro, err := srv.CRUD.RedisConfigGet(ctx, "slave-read-only")
		if err != nil {
			return util.WrapError("redis-autodiscovery", err)
		}
		if ro == "yes" {
			srv.ReadOnly = true
		}
	}
	return nil
}

// Shard is a list of the redis Server objects that compose a redis shard
type Shard struct {
	Name    string
	Servers []RedisServer
}

// NewShard returns a Shard object given the passed redis server URLs
func NewShard(name string, connectionStrings []string) (*Shard, error) {
	shard := &Shard{Name: name}
	servers := make([]RedisServer, len(connectionStrings))
	for i, cs := range connectionStrings {
		rs, err := NewRedisServerFromConnectionString(cs, cs)
		if err != nil {
			return nil, err
		}
		servers[i] = *rs
	}

	shard.Servers = servers
	return shard, nil
}

// Discover retrieves the role and read-only flag for all the servers in the shard
func (s *Shard) Discover(ctx context.Context, log logr.Logger) error {

	for idx := range s.Servers {
		if err := s.Servers[idx].Discover(ctx); err != nil {
			return err
		}
	}

	masters := 0
	for _, server := range s.Servers {
		if server.Role == client.Master {
			masters++
		}
	}

	if masters != 1 {
		err := fmt.Errorf("[redis-autodiscovery/Shard.Discover] expected 1 master but got %d", masters)
		log.Error(err, "error discovering shard server roles")
		return err
	}

	return nil
}

// GetMasterAddr returns the URL of the master server in a shard or error if zero
// or more than one master is found
func (s *Shard) GetMasterAddr() (string, string, error) {
	for _, srv := range s.Servers {
		if srv.Role == client.Master {
			ip, err := srv.IP()
			if err != nil {
				return "", "", util.WrapError("redis-autodiscovery/Shard.GetMasterAddr", err)
			}
			return ip, srv.Port, nil
		}
	}
	return "", "", fmt.Errorf("[redis-autodiscovery/Shard.GetMasterAddr] master not found")
}

// Init initializes this shard if not already initialized
func (s *Shard) Init(ctx context.Context, masterIndex int32, log logr.Logger) ([]string, error) {
	changed := []string{}

	for idx, srv := range s.Servers {
		role, slaveof, err := srv.CRUD.RedisRole(ctx)
		if err != nil {
			return changed, err
		}

		if role == client.Slave {

			if slaveof == "127.0.0.1" {

				if idx == int(masterIndex) {
					if err := srv.CRUD.RedisSlaveOf(ctx, "NO", "ONE"); err != nil {
						return changed, err
					}
					log.Info(fmt.Sprintf("[@redis-setup] Configured %s as master", srv.Name))
					changed = append(changed, srv.Name)
				} else {
					if err := srv.CRUD.RedisSlaveOf(ctx, s.Servers[masterIndex].Host, s.Servers[masterIndex].Port); err != nil {
						return changed, err
					}
					log.Info(fmt.Sprintf("[@redis-setup] Configured %s as slave", srv.Name))
					changed = append(changed, srv.Name)
				}

			} else {
				s.Servers[idx].Role = client.Slave
			}

		} else if role == client.Master {
			s.Servers[idx].Role = client.Master
		} else {
			return changed, fmt.Errorf("[@redis-setup] unable to get role for server %s", srv.Name)
		}
	}

	return changed, nil
}

// ShardedCluster represents a sharded redis cluster, composed by several Shards
type ShardedCluster []Shard

// NewShardedCluster returns a new ShardedCluster given the shard structure passed as a map[string][]string
func NewShardedCluster(ctx context.Context, serverList map[string][]string, log logr.Logger) (ShardedCluster, error) {

	sc := make([]Shard, 0, len(serverList))

	for shardName, shardServers := range serverList {

		shard, err := NewShard(shardName, shardServers)
		if err != nil {
			return nil, err
		}
		sc = append(sc, *shard)
	}

	return sc, nil
}

func (sc ShardedCluster) Discover(ctx context.Context, log logr.Logger) error {
	for _, shard := range sc {
		if err := shard.Discover(ctx, log); err != nil {
			return err
		}
	}
	return nil
}

func (sc ShardedCluster) GetShardNames() []string {
	shards := make([]string, len(sc))
	for i, shard := range sc {
		shards[i] = shard.Name
	}
	sort.Strings(shards)
	return shards
}

func (sc ShardedCluster) GetShardByName(name string) *Shard {
	for _, shard := range sc {
		if shard.Name == name {
			return &shard
		}
	}
	return nil
}
