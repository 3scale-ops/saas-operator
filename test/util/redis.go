package util

import (
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	redis "github.com/3scale/saas-operator/pkg/redis_v2/server"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

func RedisClient(cfg *rest.Config, podKey types.NamespacedName) (*redis.Server, chan struct{}, error) {
	localPort, stopCh, err := PortForward(cfg, podKey, 6379)
	if err != nil {
		return nil, nil, err
	}

	rs, err := redis.NewServer(fmt.Sprintf("redis://localhost:%d", localPort), nil)
	if err != nil {
		return nil, nil, err
	}

	return rs, stopCh, nil
}

func SentinelClient(cfg *rest.Config, podKey types.NamespacedName) (*redis.Server, chan struct{}, error) {
	localPort, stopCh, err := PortForward(cfg, podKey, saasv1alpha1.SentinelPort)
	if err != nil {
		return nil, nil, err
	}

	ss, err := redis.NewServer(fmt.Sprintf("redis://localhost:%d", localPort), nil)
	if err != nil {
		return nil, nil, err
	}

	return ss, stopCh, nil
}
