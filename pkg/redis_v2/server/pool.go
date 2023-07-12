package server

import (
	"net"
	"sync"

	"github.com/go-redis/redis/v8"
)

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

	if alias != nil {
		if srv = pool.indexByAlias()[*alias]; srv != nil {
			return srv, nil
		}
	}

	opts, err := redis.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}
	if srv = pool.indexByHost()[opts.Addr]; srv != nil {
		// set the alias if it has been passed down
		if alias != nil {
			srv.Alias = *alias
		}
		return srv, nil
	}

	// If a Server was not found, create a new one and return it
	if srv, err = NewServer(connectionString, alias); err != nil {
		return nil, err
	}
	pool.servers = append(pool.servers, srv)

	return srv, nil
}

func (pool *ServerPool) indexByAlias() map[string]*Server {
	index := make(map[string]*Server, len(pool.servers))
	for _, srv := range pool.servers {
		if srv.Alias != "" {
			index[srv.Alias] = srv
		}
	}

	return index
}

func (pool *ServerPool) indexByHost() map[string]*Server {
	index := make(map[string]*Server, len(pool.servers))
	for _, srv := range pool.servers {
		index[net.JoinHostPort(srv.Host, srv.Port)] = srv
	}

	return index
}
