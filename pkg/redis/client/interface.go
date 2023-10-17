package client

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// TestableInterface is an interface that both the go-redis and the fake client implement. It's not intended to
// support client implementations other than go-redis, it just exists to be able to inject redis server
// responses through the use of the Fake client, for testing purposes.
type TestableInterface interface {
	SentinelMaster(context.Context, string) (*SentinelMasterCmdResult, error)
	SentinelMasters(context.Context) ([]interface{}, error)
	SentinelGetMasterAddrByName(ctx context.Context, shard string) ([]string, error)
	SentinelSlaves(context.Context, string) ([]interface{}, error)
	SentinelMonitor(context.Context, string, string, string, int) error
	SentinelSet(context.Context, string, string, string) error
	SentinelPSubscribe(context.Context, ...string) (<-chan *redis.Message, func() error)
	SentinelInfoCache(context.Context) (interface{}, error)
	SentinelDo(context.Context, ...interface{}) (interface{}, error)
	SentinelPing(ctx context.Context) error
	RedisRole(context.Context) (interface{}, error)
	RedisConfigGet(context.Context, string) ([]interface{}, error)
	RedisConfigSet(context.Context, string, string) error
	RedisSlaveOf(context.Context, string, string) error
	RedisDebugSleep(context.Context, time.Duration) error
	RedisDo(context.Context, ...interface{}) (interface{}, error)
	RedisBGSave(context.Context) error
	RedisLastSave(context.Context) (int64, error)
	RedisSet(context.Context, string, interface{}) error
	RedisInfo(ctx context.Context, section string) (string, error)
	Close() error
}

// check that GoRedisClient implements Client interface
var _ TestableInterface = &GoRedisClient{}

// check that FakeClient implements Client interface
var _ TestableInterface = &FakeClient{}
