package util

import (
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/redis"
	"github.com/3scale/saas-operator/pkg/redis/crud"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

func RedisClient(cfg *rest.Config, podKey types.NamespacedName) (*crud.CRUD, chan struct{}, error) {
	localPort, stopCh, err := PortForward(cfg, podKey, 6379)
	if err != nil {
		return nil, nil, err
	}

	rs, err := redis.NewRedisServerFromConnectionString("", fmt.Sprintf("redis://localhost:%d", localPort))
	if err != nil {
		return nil, nil, err
	}

	return rs.CRUD, stopCh, nil
}

func SentinelClient(cfg *rest.Config, podKey types.NamespacedName) (*crud.CRUD, chan struct{}, error) {
	localPort, stopCh, err := PortForward(cfg, podKey, saasv1alpha1.SentinelPort)
	if err != nil {
		return nil, nil, err
	}

	ss, err := redis.NewSentinelServerFromConnectionString("", fmt.Sprintf("redis://localhost:%d", localPort))
	if err != nil {
		return nil, nil, err
	}

	return ss.CRUD, stopCh, nil
}
