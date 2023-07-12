package server

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/3scale/saas-operator/pkg/redis_v2/client"
	"github.com/go-redis/redis/v8"
)

type Server struct {
	Alias  string
	Client client.TestableInterface
	Host   string
	Port   string
	Role   client.Role
	Config map[string]string
}

// NewServer returns a new client for this redis server from the given connection
// string. It can optionally be passed an alias to identify the server.
func NewServer(connectionString string, alias *string) (*Server, error) {

	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(opt.Addr, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("error parsing redis/sentinel address")
	}

	srv := &Server{
		Host:   parts[0],
		Port:   parts[1],
		Client: client.NewFromOptions(opt),
		Role:   client.Unknown,
		Config: map[string]string{},
	}

	if alias != nil {
		srv.Alias = *alias
	}

	return srv, nil
}

func (srv *Server) CloseClient() error {
	return srv.Client.Close()
}

// NewFakeServer returns a fake client that will return the provided responses
// when called. This is only intended for testing.
func NewFakeServer(responses ...client.FakeResponse) *Server {

	return &Server{
		Host:   "fake-ip",
		Port:   "fake-port",
		Client: &client.FakeClient{Responses: responses},
	}
}

func (srv *Server) GetHost() string {
	return srv.Host
}

func (srv *Server) GetPort() string {
	return srv.Port
}

func (srv *Server) ID() string {
	if srv.Alias != "" {
		return srv.Alias
	} else {
		return net.JoinHostPort(srv.Host, srv.Port)
	}
}

func (rs *Server) IP() (string, error) {
	var ip string
	if r := net.ParseIP(rs.GetHost()); r != nil {
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

func (srv *Server) SentinelMaster(ctx context.Context, shard string) (*client.SentinelMasterCmdResult, error) {

	result, err := srv.Client.SentinelMaster(ctx, shard)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (srv *Server) SentinelMasters(ctx context.Context) ([]client.SentinelMasterCmdResult, error) {

	values, err := srv.Client.SentinelMasters(ctx)
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

	values, err := srv.Client.SentinelSlaves(ctx, shard)
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
	return srv.Client.SentinelMonitor(ctx, name, host, port, quorum)
}

func (srv *Server) SentinelSet(ctx context.Context, shard, parameter, value string) error {
	return srv.Client.SentinelSet(ctx, shard, parameter, value)
}

func (srv *Server) SentinelPSubscribe(ctx context.Context, events ...string) (<-chan *redis.Message, func() error) {
	return srv.Client.SentinelPSubscribe(ctx, events...)
}

func (srv *Server) SentinelInfoCache(ctx context.Context) (client.SentinelInfoCache, error) {
	result := client.SentinelInfoCache{}

	raw, err := srv.Client.SentinelInfoCache(ctx)
	mval := islice2imap(raw)

	for shard, servers := range mval {
		result[shard] = make(map[string]client.RedisServerInfoCache, len(servers.([]interface{})))

		for _, server := range servers.([]interface{}) {
			// When sentinel is unable to reach the redis slave the info field can be nil
			// so we have to check this to avoid panics
			if server.([]interface{})[1] != nil {
				info := infoStringToMap(server.([]interface{})[1].(string))
				result[shard][info["run_id"]] = client.RedisServerInfoCache{
					CacheAge: time.Duration(server.([]interface{})[0].(int64)) * time.Millisecond,
					Info:     info,
				}
			}
		}
	}

	return result, err
}

func (srv *Server) RedisRole(ctx context.Context) (client.Role, string, error) {
	val, err := srv.Client.RedisRole(ctx)
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
	val, err := srv.Client.RedisConfigGet(ctx, parameter)
	if err != nil {
		return "", err
	}
	return val[1].(string), nil
}

func (srv *Server) RedisConfigSet(ctx context.Context, parameter, value string) error {
	return srv.Client.RedisConfigSet(ctx, parameter, value)
}

func (srv *Server) RedisSlaveOf(ctx context.Context, host, port string) error {
	return srv.Client.RedisSlaveOf(ctx, host, port)
}

func (srv *Server) RedisDebugSleep(ctx context.Context, duration time.Duration) error {
	return srv.Client.RedisDebugSleep(ctx, duration)
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

func infoStringToMap(in string) map[string]string {

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
