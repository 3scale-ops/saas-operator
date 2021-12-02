package client

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type RedisGoClient struct {
	redis    *redis.Client
	sentinel *redis.SentinelClient
}

// check that SimpleClient implements Client interface
// var _ Client = &SimpleClient{}

func NewFromConnectionString(connectionString string) (*RedisGoClient, error) {
	var err error
	c := &RedisGoClient{}

	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}

	c.redis = redis.NewClient(opt)
	c.sentinel = redis.NewSentinelClient(opt)

	return c, nil
}

func NewFromOptions(opt *redis.Options) *RedisGoClient {
	return &RedisGoClient{
		redis:    redis.NewClient(opt),
		sentinel: redis.NewSentinelClient(opt),
	}
}

func (c *RedisGoClient) SentinelMaster(ctx context.Context, shard string) (*SentinelMasterCmdResult, error) {

	result := &SentinelMasterCmdResult{}
	err := c.sentinel.Master(ctx, shard).Scan(result)
	return result, err
}

func (c *RedisGoClient) SentinelMasters(ctx context.Context) ([]interface{}, error) {

	values, err := c.sentinel.Masters(ctx).Result()
	return values, err
}

func (c *RedisGoClient) SentinelSlaves(ctx context.Context, shard string) ([]interface{}, error) {

	values, err := c.sentinel.Slaves(ctx, shard).Result()
	return values, err
}

func (c *RedisGoClient) SentinelMonitor(ctx context.Context, name, host string, port string, quorum int) error {

	_, err := c.sentinel.Monitor(ctx, name, host, port, strconv.Itoa(quorum)).Result()
	return err
}

func (c *RedisGoClient) SentinelSet(ctx context.Context, shard, parameter, value string) error {

	_, err := c.sentinel.Set(ctx, shard, parameter, value).Result()
	return err
}

func (c *RedisGoClient) SentinelPSubscribe(ctx context.Context, events ...string) (<-chan *redis.Message, func() error) {

	pubsub := c.sentinel.PSubscribe(ctx, events...)
	return pubsub.Channel(), pubsub.Close
}

func (c *RedisGoClient) RedisRole(ctx context.Context) (interface{}, error) {

	val, err := c.redis.Do(ctx, "role").Result()
	return val, err
}

func (c *RedisGoClient) RedisConfigGet(ctx context.Context, parameter string) ([]interface{}, error) {

	val, err := c.redis.ConfigGet(ctx, parameter).Result()
	return val, err
}

func (c *RedisGoClient) RedisSlaveOf(ctx context.Context, host, port string) error {

	_, err := c.redis.SlaveOf(ctx, host, port).Result()
	return err
}
