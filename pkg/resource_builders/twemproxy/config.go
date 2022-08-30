package twemproxy

import (
	"fmt"
	"strconv"
	"strings"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
)

const (
	HealthPoolName    string = "health"
	HealthBindAddress string = "127.0.0.1:22333"
)

type Server struct {
	Address  string
	Priority int
	Name     string
}

func (srv *Server) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s:%d %s\"", srv.Address, srv.Priority, srv.Name)), nil
}

func (srv *Server) UnmarshalJSON(data []byte) error {
	parts := strings.Split(strings.Trim(string(data), "\""), " ")
	srv.Name = parts[1]
	parts = strings.Split(parts[0], ":")
	p, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return err
	}
	srv.Priority = p
	srv.Address = strings.Join(parts[0:len(parts)-1], ":")
	return nil
}

type ServerPoolConfig struct {
	Listen             string   `json:"listen"`
	Hash               string   `json:"hash,omitempty"`
	HashTag            string   `json:"hash_tag,omitempty"`
	Distribution       string   `json:"distribution,omitempty"`
	Timeout            int      `json:"timeout,omitempty"`
	Backlog            int      `json:"backlog,omitempty"`
	PreConnect         bool     `json:"preconnect"`
	Redis              bool     `json:"redis"`
	AutoEjectHosts     bool     `json:"auto_eject_hosts"`
	ServerFailureLimit int      `json:"server_failure_limit,omitempty"`
	Servers            []Server `json:"servers"`
}

func GenerateServerPool(pool saasv1alpha1.TwemproxyServerPool, targets map[string]Server) ServerPoolConfig {

	servers := make([]Server, 0, len(pool.Topology))
	for _, s := range pool.Topology {
		srv := targets[s.PhysicalShard]
		srv.Name = s.ShardName
		servers = append(servers, srv)
	}

	return ServerPoolConfig{
		// The following parameters cannot be changed
		Hash:           "fnv1a_64",
		HashTag:        "{}",
		Distribution:   "ketama",
		AutoEjectHosts: false,
		Redis:          true,
		// The following parameters could be safely modified or exposed in the CR
		Listen:     pool.BindAddress,
		Backlog:    pool.TCPBacklog,
		PreConnect: pool.PreConnect,
		Timeout:    pool.Timeout,
		// The list of servers is generated from the
		// list fo shards provided by the user in the Backend spec
		Servers: servers,
	}
}
