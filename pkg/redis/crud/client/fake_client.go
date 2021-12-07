package client

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type FakeClient struct {
	InjectResponse func() interface{}
	InjectError    func() error
}

func (fc *FakeClient) SentinelMaster(ctx context.Context, shard string) (*SentinelMasterCmdResult, error) {
	return fc.InjectResponse().(*SentinelMasterCmdResult), fc.InjectError()
}

func (fc *FakeClient) SentinelMasters(ctx context.Context) ([]interface{}, error) {
	return fc.InjectResponse().([]interface{}), fc.InjectError()
}

func (fc *FakeClient) SentinelSlaves(ctx context.Context, shard string) ([]interface{}, error) {
	return fc.InjectResponse().([]interface{}), fc.InjectError()
}

func (fc *FakeClient) SentinelMonitor(ctx context.Context, name, host string, port string, quorum int) error {
	return fc.InjectError()
}

func (fc *FakeClient) SentinelSet(ctx context.Context, shard, parameter, value string) error {
	return fc.InjectError()
}

func (fc *FakeClient) SentinelPSubscribe(ctx context.Context, events ...string) (<-chan *redis.Message, func() error) {
	return fc.InjectResponse().(<-chan *redis.Message), nil
}

func (fc *FakeClient) RedisRole(ctx context.Context) (interface{}, error) {
	return fc.InjectResponse(), fc.InjectError()
}

func (fc *FakeClient) RedisConfigGet(ctx context.Context, parameter string) ([]interface{}, error) {
	return fc.InjectResponse().([]interface{}), fc.InjectError()
}

func (fc *FakeClient) RedisSlaveOf(ctx context.Context, host, port string) error {
	return fc.InjectError()
}
