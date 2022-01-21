package client

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type FakeClient struct {
	Responses []FakeResponse
}

type FakeResponse struct {
	InjectResponse func() interface{}
	InjectError    func() error
}

func (fc *FakeClient) SentinelMaster(ctx context.Context, shard string) (*SentinelMasterCmdResult, error) {
	rsp := fc.pop()
	return rsp.InjectResponse().(*SentinelMasterCmdResult), rsp.InjectError()
}

func (fc *FakeClient) SentinelMasters(ctx context.Context) ([]interface{}, error) {
	rsp := fc.pop()
	return rsp.InjectResponse().([]interface{}), rsp.InjectError()
}

func (fc *FakeClient) SentinelSlaves(ctx context.Context, shard string) ([]interface{}, error) {
	rsp := fc.pop()
	return rsp.InjectResponse().([]interface{}), rsp.InjectError()
}

func (fc *FakeClient) SentinelMonitor(ctx context.Context, name, host string, port string, quorum int) error {
	rsp := fc.pop()
	return rsp.InjectError()
}

func (fc *FakeClient) SentinelSet(ctx context.Context, shard, parameter, value string) error {
	rsp := fc.pop()
	return rsp.InjectError()
}

func (fc *FakeClient) SentinelPSubscribe(ctx context.Context, events ...string) (<-chan *redis.Message, func() error) {
	rsp := fc.pop()
	return rsp.InjectResponse().(<-chan *redis.Message), nil
}

func (fc *FakeClient) SentinelInfoCache(ctx context.Context) (interface{}, error) {
	rsp := fc.pop()
	return rsp.InjectResponse(), rsp.InjectError()
}

func (fc *FakeClient) RedisRole(ctx context.Context) (interface{}, error) {
	rsp := fc.pop()
	return rsp.InjectResponse(), rsp.InjectError()
}

func (fc *FakeClient) RedisConfigGet(ctx context.Context, parameter string) ([]interface{}, error) {
	rsp := fc.pop()
	return rsp.InjectResponse().([]interface{}), rsp.InjectError()
}

func (fc *FakeClient) RedisSlaveOf(ctx context.Context, host, port string) error {
	rsp := fc.pop()
	return rsp.InjectError()
}

func (fc *FakeClient) pop() (fakeRsp FakeResponse) {
	fakeRsp, fc.Responses = fc.Responses[0], fc.Responses[1:]
	return fakeRsp
}
