package util

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

const (
	AWSAccessKeyEnvvar string = "AWS_ACCESS_KEY_ID"
	AWSSecretKeyEnvvar string = "AWS_SECRET_ACCESS_KEY"
	AWSRegionEnvvar    string = "AWS_REGION"
)

func AWSConfig(ctx context.Context, accessKeyID, secretAccessKey, region string, serviceEndpoint *string) (*aws.Config, error) {

	var cfg aws.Config
	var err error

	if serviceEndpoint != nil {
		resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               *serviceEndpoint,
				SigningRegion:     region,
				HostnameImmutable: true,
			}, nil
		})

		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithEndpointResolverWithOptions(resolver),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
		)
		if err != nil {
			return nil, err
		}
	} else {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
		)
		if err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}
