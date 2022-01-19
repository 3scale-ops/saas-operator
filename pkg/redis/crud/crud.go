package crud

import (
	"bufio"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/3scale/saas-operator/pkg/redis/crud/client"
	"github.com/go-redis/redis/v8"
)

type CRUD struct {
	Client Client
	IP     string
	Port   string
}

type Client interface {
	SentinelMaster(context.Context, string) (*client.SentinelMasterCmdResult, error)
	SentinelMasters(context.Context) ([]interface{}, error)
	SentinelSlaves(context.Context, string) ([]interface{}, error)
	SentinelMonitor(context.Context, string, string, string, int) error
	SentinelSet(context.Context, string, string, string) error
	SentinelPSubscribe(context.Context, ...string) (<-chan *redis.Message, func() error)
	SentinelInfoCache(context.Context) (interface{}, error)
	RedisRole(context.Context) (interface{}, error)
	RedisConfigGet(context.Context, string) ([]interface{}, error)
	RedisSlaveOf(context.Context, string, string) error
}

// check that GoRedisClient implements Client interface
var _ Client = &client.GoRedisClient{}

// check that FakeClient implements Client interface
var _ Client = &client.FakeClient{}

func NewRedisCRUDFromConnectionString(connectionString string) (*CRUD, error) {

	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(opt.Addr, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("error parsing redis/sentinel address")
	}

	return &CRUD{
		IP:     parts[0],
		Port:   parts[1],
		Client: client.NewFromOptions(opt),
	}, nil
}

func NewFakeCRUD(responses ...client.FakeResponse) *CRUD {

	return &CRUD{
		IP:     "fake-ip",
		Port:   "fake-port",
		Client: &client.FakeClient{Responses: responses},
	}
}

func (crud *CRUD) GetIP() string {
	return crud.IP
}

func (sc *CRUD) GetPort() string {
	return sc.Port
}

func (crud *CRUD) SentinelMaster(ctx context.Context, shard string) (*client.SentinelMasterCmdResult, error) {

	result, err := crud.Client.SentinelMaster(ctx, shard)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (crud *CRUD) SentinelMasters(ctx context.Context) ([]client.SentinelMasterCmdResult, error) {

	values, err := crud.Client.SentinelMasters(ctx)
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

func (crud *CRUD) SentinelSlaves(ctx context.Context, shard string) ([]client.SentinelSlaveCmdResult, error) {

	values, err := crud.Client.SentinelSlaves(ctx, shard)
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

func (crud *CRUD) SentinelMonitor(ctx context.Context, name, host string, port string, quorum int) error {
	return crud.Client.SentinelMonitor(ctx, name, host, port, quorum)
}

func (crud *CRUD) SentinelSet(ctx context.Context, shard, parameter, value string) error {
	return crud.Client.SentinelSet(ctx, shard, parameter, value)
}

func (crud *CRUD) SentinelPSubscribe(ctx context.Context, events ...string) (<-chan *redis.Message, func() error) {
	return crud.Client.SentinelPSubscribe(ctx, events...)
}

func (crud *CRUD) SentinelInfoCache(ctx context.Context) (client.SentinelInfoCache, error) {
	result := client.SentinelInfoCache{}

	raw, err := crud.Client.SentinelInfoCache(ctx)
	mval := islice2imap(raw)

	for shard, servers := range mval {
		result[shard] = make(map[string]client.RedisServerInfoCache, len(servers.([]interface{})))

		for _, server := range servers.([]interface{}) {

			info := infoStringToMap(server.([]interface{})[1].(string))
			result[shard][info["run_id"]] = client.RedisServerInfoCache{
				CacheAge: time.Duration(server.([]interface{})[0].(int64)) * time.Millisecond,
				Info:     info,
			}
		}
	}

	return result, err
}

func (crud *CRUD) RedisRole(ctx context.Context) (client.Role, string, error) {
	val, err := crud.Client.RedisRole(ctx)
	if err != nil {
		return client.Unknown, "", err
	}

	if client.Role(val.([]interface{})[0].(string)) == client.Master {
		return client.Master, "", nil
	} else {
		return client.Slave, val.([]interface{})[1].(string), nil
	}
}

func (crud *CRUD) RedisConfigGet(ctx context.Context, parameter string) (string, error) {
	val, err := crud.Client.RedisConfigGet(ctx, parameter)
	if err != nil {
		return "", err
	}
	return val[1].(string), nil
}

func (sc *CRUD) RedisSlaveOf(ctx context.Context, host, port string) error {
	return sc.Client.RedisSlaveOf(ctx, host, port)
}

// This is a horrible function to parse the horrible structs that the redis-go
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
