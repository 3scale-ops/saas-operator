package server

import (
	"net"
	"sort"
	"sync"

	"github.com/go-redis/redis/v8"
)

// SentinelPool holds a thread safe list of Servers.
// It is intended for client reuse throughout the code.
type ServerPool struct {
	servers []*Server
	mu      sync.Mutex
}

func NewServerPool(servers ...*Server) *ServerPool {
	if len(servers) > 0 {
		return &ServerPool{
			servers: servers,
		}
	}

	return &ServerPool{
		servers: []*Server{},
	}
}

func (pool *ServerPool) GetServer(connectionString string, alias *string) (*Server, error) {
	var srv *Server
	var err error

	// make sure both reads and writes are consistent
	// might cause some contention but grants consistency
	pool.mu.Lock()
	defer pool.mu.Unlock()

	opts, err := redis.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}
	if srv = pool.indexByHostPort()[opts.Addr]; srv != nil {
		// set the alias if it has been passed down
		if alias != nil {
			srv.alias = *alias
		}
		return srv, nil
	}

	// If a Server was not found, create a new one and return it
	if srv, err = NewServer(connectionString, alias); err != nil {
		return nil, err
	}
	pool.servers = append(pool.servers, srv)

	// sort the slice to obtain consistent results
	sort.Slice(pool.servers, func(i, j int) bool {
		return pool.servers[i].ID() < pool.servers[j].ID()
	})

	return srv, nil
}

func (pool *ServerPool) indexByAlias() map[string]*Server {
	index := make(map[string]*Server, len(pool.servers))
	for _, srv := range pool.servers {
		if srv.alias != "" {
			index[srv.alias] = srv
		}
	}

	return index
}

func (pool *ServerPool) indexByHostPort() map[string]*Server {
	index := make(map[string]*Server, len(pool.servers))
	for _, srv := range pool.servers {
		index[net.JoinHostPort(srv.host, srv.port)] = srv
	}

	return index
}
