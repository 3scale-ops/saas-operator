package crud

import (
	"context"
	"fmt"
	"strings"

	"github.com/3scale/saas-operator/pkg/redis/crud/client"
	"github.com/go-redis/redis/v8"
)

type CRUD struct {
	client Client
	ip     string
	port   string
}

type Client interface {
	SentinelMaster(context.Context, string) (*client.SentinelMasterCmdResult, error)
	SentinelMasters(context.Context) ([]interface{}, error)
	SentinelSlaves(context.Context, string) ([]interface{}, error)
	SentinelMonitor(context.Context, string, string, string, int) error
	SentinelSet(context.Context, string, string, string) error
	SentinelPSubscribe(context.Context, ...string) (<-chan *redis.Message, func() error)
	RedisRole(context.Context) (interface{}, error)
	RedisConfigGet(context.Context, string) ([]interface{}, error)
	RedisSlaveOf(context.Context, string, string) error
}

// check that RedisGoCLient implements Client interface
var _ Client = &client.RedisGoClient{}

// check that FakeClient implements Client interface
var _ Client = &client.FakeClient{}

func NewRedisCRUD(connectionString string) (*CRUD, error) {

	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(opt.Addr, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("error parsing redis/sentinel address")
	}

	return &CRUD{
		ip:     parts[0],
		port:   parts[1],
		client: client.NewFromOptions(opt),
	}, nil

}

func (crud *CRUD) GetIP() string {
	return crud.ip
}

func (sc *CRUD) GetPort() string {
	return sc.port
}

func (crud *CRUD) SentinelMaster(ctx context.Context, shard string) (*client.SentinelMasterCmdResult, error) {

	result, err := crud.client.SentinelMaster(ctx, shard)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (crud *CRUD) SentinelMasters(ctx context.Context) ([]client.SentinelMasterCmdResult, error) {

	values, err := crud.client.SentinelMasters(ctx)
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

	values, err := crud.client.SentinelSlaves(ctx, shard)
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
	return crud.client.SentinelMonitor(ctx, name, host, port, quorum)
}

func (crud *CRUD) SentinelSet(ctx context.Context, shard, parameter, value string) error {
	return crud.client.SentinelSet(ctx, shard, parameter, value)
}

func (crud *CRUD) SentinelPSubscribe(ctx context.Context, events ...string) (<-chan *redis.Message, func() error) {
	return crud.client.SentinelPSubscribe(ctx, events...)
}

func (crud *CRUD) RedisRole(ctx context.Context) (client.Role, string, error) {
	val, err := crud.client.RedisRole(ctx)
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
	val, err := crud.client.RedisConfigGet(ctx, parameter)
	if err != nil {
		return "", err
	}
	return val[1].(string), nil
}

func (sc *CRUD) RedisSlaveOf(ctx context.Context, host, port string) error {
	return sc.client.RedisSlaveOf(ctx, host, port)
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
