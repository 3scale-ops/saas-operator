package server

import (
	"bufio"
	"context"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/3scale/saas-operator/pkg/redis/client"
	"github.com/go-redis/redis/v8"
)

// Server is a host that talks the redis protocol
// Contains methods to use a subset of redis commands
type Server struct {
	alias  string
	client client.TestableInterface
	host   string
	port   string
	mu     sync.Mutex
}

// NewServer returns a new client for this redis server from the given connection
// string. It can optionally be passed an alias to identify the server.
func NewServer(connectionString string, alias *string) (*Server, error) {

	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}

	host, port, err := net.SplitHostPort(opt.Addr)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		host:   host,
		port:   port,
		client: client.NewFromOptions(opt),
	}

	if alias != nil {
		srv.SetAlias(*alias)
	}

	return srv, nil
}

func MustNewServer(connectionString string, alias *string) *Server {
	srv, err := NewServer(connectionString, alias)
	if err != nil {
		panic(err)
	}
	return srv
}

func NewServerFromParams(alias, host, port string, c client.TestableInterface) *Server {
	return &Server{
		alias:  alias,
		host:   host,
		port:   port,
		client: c,
	}
}

func (srv *Server) CloseClient() error {
	return srv.client.Close()
}

func (srv *Server) GetClient() client.TestableInterface {
	return srv.client
}

func (srv *Server) GetHost() string {
	return srv.host
}

func (srv *Server) GetPort() string {
	return srv.port
}

func (srv *Server) GetAlias() string {
	if srv.alias != "" {
		return srv.alias
	}
	return srv.ID()
}

func (srv *Server) SetAlias(alias string) {
	srv.mu.Lock()
	srv.alias = alias
	srv.mu.Unlock()
}

// ID returns the ID of the server, which takes the form "host:port"
func (srv *Server) ID() string {
	return net.JoinHostPort(srv.host, srv.port)
}

func (srv *Server) SentinelMaster(ctx context.Context, shard string) (*client.SentinelMasterCmdResult, error) {

	result, err := srv.client.SentinelMaster(ctx, shard)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (srv *Server) SentinelMasters(ctx context.Context) ([]client.SentinelMasterCmdResult, error) {

	values, err := srv.client.SentinelMasters(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]client.SentinelMasterCmdResult, len(values))
	for i, val := range values {
		masterResult := &client.SentinelMasterCmdResult{}
		err := sliceCmdToStruct(val, masterResult)
		if err != nil {
			return nil, err
		}
		result[i] = *masterResult
	}

	return result, nil
}

func (srv *Server) SentinelSlaves(ctx context.Context, shard string) ([]client.SentinelSlaveCmdResult, error) {

	values, err := srv.client.SentinelSlaves(ctx, shard)
	if err != nil {
		return nil, err
	}

	result := make([]client.SentinelSlaveCmdResult, len(values))
	for i, val := range values {
		slaveResult := &client.SentinelSlaveCmdResult{}
		err := sliceCmdToStruct(val, slaveResult)
		if err != nil {
			return nil, err
		}
		result[i] = *slaveResult
	}

	return result, nil
}

func (srv *Server) SentinelMonitor(ctx context.Context, name, host string, port string, quorum int) error {
	return srv.client.SentinelMonitor(ctx, name, host, port, quorum)
}

func (srv *Server) SentinelSet(ctx context.Context, shard, parameter, value string) error {
	return srv.client.SentinelSet(ctx, shard, parameter, value)
}

func (srv *Server) SentinelPSubscribe(ctx context.Context, events ...string) (<-chan *redis.Message, func() error) {
	return srv.client.SentinelPSubscribe(ctx, events...)
}

func (srv *Server) SentinelInfoCache(ctx context.Context) (client.SentinelInfoCache, error) {
	result := client.SentinelInfoCache{}

	raw, err := srv.client.SentinelInfoCache(ctx)
	mval := islice2imap(raw)

	for shard, servers := range mval {
		result[shard] = make(map[string]client.RedisServerInfoCache, len(servers.([]interface{})))

		for _, server := range servers.([]interface{}) {
			// When sentinel is unable to reach the redis slave the info field can be nil
			// so we have to check this to avoid panics
			if server.([]interface{})[1] != nil {
				info := InfoStringToMap(server.([]interface{})[1].(string))
				result[shard][info["run_id"]] = client.RedisServerInfoCache{
					CacheAge: time.Duration(server.([]interface{})[0].(int64)) * time.Millisecond,
					Info:     info,
				}
			}
		}
	}

	return result, err
}

func (srv *Server) SentinelPing(ctx context.Context) error {
	return srv.client.SentinelPing(ctx)
}

func (srv *Server) RedisRole(ctx context.Context) (client.Role, string, error) {
	val, err := srv.client.RedisRole(ctx)
	if err != nil {
		return client.Unknown, "", err
	}

	if client.Role(val.([]interface{})[0].(string)) == client.Master {
		return client.Master, "", nil
	} else {
		return client.Slave, val.([]interface{})[1].(string), nil
	}
}

func (srv *Server) RedisConfigGet(ctx context.Context, parameter string) (string, error) {
	val, err := srv.client.RedisConfigGet(ctx, parameter)
	if err != nil {
		return "", err
	}
	return val[1].(string), nil
}

func (srv *Server) RedisConfigSet(ctx context.Context, parameter, value string) error {
	return srv.client.RedisConfigSet(ctx, parameter, value)
}

func (srv *Server) RedisSlaveOf(ctx context.Context, host, port string) error {
	return srv.client.RedisSlaveOf(ctx, host, port)
}

func (srv *Server) RedisDebugSleep(ctx context.Context, duration time.Duration) error {
	return srv.client.RedisDebugSleep(ctx, duration)
}

func (srv *Server) RedisBGSave(ctx context.Context) error {
	return srv.client.RedisBGSave(ctx)
}

func (srv *Server) RedisLastSave(ctx context.Context) (int64, error) {
	return srv.client.RedisLastSave(ctx)
}

func (srv *Server) RedisSet(ctx context.Context, key string, value interface{}) error {
	return srv.client.RedisSet(ctx, key, value)
}

// This is a horrible function to parse the horrible structs that the go-redis
// client returns for administrative commands. I swear it's not my fault ...
func sliceCmdToStruct(in interface{}, out interface{}) error {
	m := map[string]string{}
	for i := range in.([]interface{}) {
		if i%2 != 0 {
			continue
		}
		m[in.([]interface{})[i].(string)] = in.([]interface{})[i+1].(string)
	}

	err := redis.NewStringStringMapResult(m, nil).Scan(out)
	if err != nil {
		return err
	}
	return nil
}

func islice2imap(in interface{}) map[string]interface{} {
	m := map[string]interface{}{}
	for i := range in.([]interface{}) {
		if i%2 != 0 {
			continue
		}
		m[in.([]interface{})[i].(string)] = in.([]interface{})[i+1].([]interface{})
	}
	return m
}

func InfoStringToMap(in string) map[string]string {

	m := map[string]string{}
	scanner := bufio.NewScanner(strings.NewReader(in))
	for scanner.Scan() {
		// do not add empty lines or section headings (see the test for more info)
		if line := scanner.Text(); line != "" && !strings.HasPrefix(line, "# ") {
			kv := strings.SplitN(line, ":", 2)
			m[kv[0]] = kv[1]
		}
	}

	return m
}
