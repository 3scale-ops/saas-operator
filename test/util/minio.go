package util

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

const (
	minioPort uint32 = 9000
	region    string = "us-east-1"
)

func MinioClient(ctx context.Context, cfg *rest.Config, podKey types.NamespacedName, user, passwd string) (*s3.Client, chan struct{}, error) {

	// set credentials
	os.Setenv("AWS_ACCESS_KEY_ID", user)
	os.Setenv("AWS_SECRET_ACCESS_KEY", passwd)
	os.Setenv("AWS_REGION", region)

	localPort, stopCh, err := PortForward(cfg, podKey, minioPort)
	if err != nil {
		return nil, nil, err
	}

	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               fmt.Sprintf("http://localhost:%d", localPort),
			SigningRegion:     region,
			HostnameImmutable: true,
		}, nil
	})

	awsconfig, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(resolver),
		config.WithClientLogMode(aws.LogResponseWithBody),
	)

	if err != nil {
		return nil, nil, err
	}
	client := s3.NewFromConfig(awsconfig)
	return client, stopCh, nil
}
