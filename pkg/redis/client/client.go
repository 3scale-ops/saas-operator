package client

import (
	"context"
	"strconv"

	redistypes "github.com/3scale/saas-operator/pkg/redis/types"
	"github.com/go-redis/redis/v8"
)

type Client interface {
	SentinelMaster(context.Context, string) (*redistypes.SentinelMasterCmdResult, error)
	SentinelMasters(context.Context) ([]redistypes.SentinelMasterCmdResult, error)
	SentinelSlaves(context.Context, string) ([]redistypes.SentinelSlaveCmdResult, error)
	SentinelMonitor(context.Context, string, string, string, int) error
	SentinelSet(context.Context, string, string, string) error
	SentinelPSubscribe(context.Context, ...string) (<-chan *redis.Message, func() error)
	RedisRole(context.Context) (redistypes.Role, string, error)
	RedisConfigGet(context.Context) (string, error)
	RedisSlaveOf(context.Context, string, string) error
}

type SimpleClient struct {
	config *redis.Options
}

// check that SimpleClient implements Client interface
var _ Client = &SimpleClient{}

func NewSimpleClient(connectionString string) (*SimpleClient, error) {
	var err error
	sclient := &SimpleClient{}

	sclient.config, err = redis.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}

	return sclient, nil
}

func (sc *SimpleClient) RedisClient(ctx context.Context) *redis.Client {
	return redis.NewClient(sc.config).WithContext(ctx)
}

func (sc *SimpleClient) SentinelClient(ctx context.Context) *redis.SentinelClient {
	return redis.NewSentinelClient(sc.config).WithContext(ctx)
}

func (sc *SimpleClient) SentinelMaster(ctx context.Context, shard string) (*redistypes.SentinelMasterCmdResult, error) {

	result := &redistypes.SentinelMasterCmdResult{}
	if err := sc.SentinelClient(ctx).Master(ctx, shard).Scan(result); err != nil {
		return nil, err
	}
	return result, nil
}

func (sc *SimpleClient) SentinelMasters(ctx context.Context) ([]redistypes.SentinelMasterCmdResult, error) {

	values, err := sc.SentinelClient(ctx).Masters(ctx).Result()
	if err != nil {
		return nil, err
	}

	result := make([]redistypes.SentinelMasterCmdResult, len(values))
	for i, val := range values {
		masterResult := &redistypes.SentinelMasterCmdResult{}
		sliceCmdToStruct(val, masterResult)
		if err != nil {
			return nil, err
		}
		result[i] = *masterResult
	}

	return result, nil
}

func (sc *SimpleClient) SentinelSlaves(ctx context.Context, shard string) ([]redistypes.SentinelSlaveCmdResult, error) {

	values, err := sc.SentinelClient(ctx).Slaves(ctx, shard).Result()
	if err != nil {
		return nil, err
	}

	result := make([]redistypes.SentinelSlaveCmdResult, len(values))
	for i, val := range values {
		slaveResult := &redistypes.SentinelSlaveCmdResult{}
		sliceCmdToStruct(val, slaveResult)
		result[i] = *slaveResult
	}

	return result, nil
}

func (sc *SimpleClient) SentinelMonitor(ctx context.Context, name, host string, port string, quorum int) error {
	_, err := sc.SentinelClient(ctx).Monitor(ctx, name, host, port, strconv.Itoa(quorum)).Result()
	return err
}

func (sc *SimpleClient) SentinelSet(ctx context.Context, shard, parameter, value string) error {
	_, err := sc.SentinelClient(ctx).Set(ctx, shard, parameter, value).Result()
	return err
}

func (sc *SimpleClient) SentinelPSubscribe(ctx context.Context, events ...string) (<-chan *redis.Message, func() error) {

	pubsub := sc.SentinelClient(ctx).PSubscribe(ctx, events...)
	return pubsub.Channel(), pubsub.Close
}

func (sc *SimpleClient) RedisRole(ctx context.Context) (redistypes.Role, string, error) {
	val, err := sc.RedisClient(ctx).Do(ctx, "role").Result()
	if err != nil {
		return redistypes.Unknown, "", err
	}

	if redistypes.Role(val.([]interface{})[0].(string)) == redistypes.Master {
		return redistypes.Master, "", nil
	} else {
		return redistypes.Slave, val.([]interface{})[1].(string), nil
	}
}

func (sc *SimpleClient) RedisConfigGet(ctx context.Context) (string, error) {
	val, err := sc.RedisClient(ctx).ConfigGet(ctx, "slave-read-only").Result()
	if err != nil {
		return "", err
	}
	return val[1].(string), nil
}

func (sc *SimpleClient) RedisSlaveOf(ctx context.Context, host, port string) error {
	_, err := sc.RedisClient(ctx).SlaveOf(ctx, host, port).Result()
	return err
}

// This is a horrible function to parse the horrible structs that the redis-go
// client returns for administrative commands. I swear it's not my fault ...
func sliceCmdToStruct(in interface{}, out interface{}) (interface{}, error) {
	m := map[string]string{}
	for i := range in.([]interface{}) {
		if i%2 != 0 {
			continue
		}
		m[in.([]interface{})[i].(string)] = in.([]interface{})[i+1].(string)
	}

	err := redis.NewStringStringMapResult(m, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
