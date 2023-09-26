package util

import (
	"context"
	"fmt"

	"github.com/3scale/saas-operator/pkg/util"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

const (
	minioPort uint32 = 9000
	region    string = "us-east-1"
)

func MinioClient(ctx context.Context, cfg *rest.Config, podKey types.NamespacedName, user, passwd string) (*s3.Client, chan struct{}, error) {

	localPort, stopCh, err := PortForward(cfg, podKey, minioPort)
	if err != nil {
		return nil, nil, err
	}

	awsconfig, err := util.AWSConfig(ctx, user, passwd, region, util.Pointer(fmt.Sprintf("http://localhost:%d", localPort)))
	if err != nil {
		return nil, nil, err
	}
	client := s3.NewFromConfig(*awsconfig)
	return client, stopCh, nil
}
