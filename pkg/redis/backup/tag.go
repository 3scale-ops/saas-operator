package backup

import (
	"context"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (br *Runner) TagBackup(ctx context.Context) error {
	logger := log.FromContext(ctx, "function", "(br *Runner) TagBackup()")

	var cfg aws.Config
	var err error

	if br.AWSS3Endpoint != nil {
		resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               *br.AWSS3Endpoint,
				SigningRegion:     br.AWSRegion,
				HostnameImmutable: true,
			}, nil
		})

		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(br.AWSRegion),
			config.WithEndpointResolverWithOptions(resolver),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(br.AWSAccessKeyID, br.AWSSecretAccessKey, "")),
		)
		if err != nil {
			return err
		}
	} else {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(br.AWSRegion),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(br.AWSAccessKeyID, br.AWSSecretAccessKey, "")),
		)
		if err != nil {
			return err
		}
	}

	client := s3.NewFromConfig(cfg)

	// get backups of current day
	dayResult, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(br.S3Bucket),
		Prefix: aws.String(br.S3Path + "/" + br.BackupFileBaseNameWithTimeSuffix(br.Timestamp.Format("2006-01-02"))),
	})
	if err != nil {
		return err
	}
	sort.SliceStable(dayResult.Contents, func(i, j int) bool {
		return dayResult.Contents[i].LastModified.Before(*dayResult.Contents[j].LastModified)
	})

	// get backups of current hour
	hourResult, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(br.S3Bucket),
		Prefix: aws.String(br.S3Path + "/" + br.BackupFileBaseNameWithTimeSuffix(br.Timestamp.Format("2006-01-02T15"))),
	})
	if err != nil {
		return err
	}
	sort.SliceStable(hourResult.Contents, func(i, j int) bool {
		return hourResult.Contents[i].LastModified.Before(*hourResult.Contents[j].LastModified)
	})

	firstOfDay := dayResult.Contents[0]
	firstOfHour := hourResult.Contents[0]

	last := hourResult.Contents[len(hourResult.Contents)-1]
	if br.BackupFileS3Path() != *last.Key {
		return fmt.Errorf("last backup %s has different key than expected (%s)", *last.Key, br.BackupFileS3Path())
	}
	// store backup size
	br.status.BackupSize = last.Size

	tags := []types.Tag{
		{Key: aws.String("Layer"), Value: aws.String("bck-storage")},
		{Key: aws.String("App"), Value: aws.String("Backend")},
	}

	switch br.BackupFileS3Path() {
	case *firstOfDay.Key:
		tags = append(tags, types.Tag{Key: aws.String("Retention"), Value: aws.String("90d")})
		logger.V(1).Info("backup tagged with 90d retention")
	case *firstOfHour.Key:
		tags = append(tags, types.Tag{Key: aws.String("Retention"), Value: aws.String("7d")})
		logger.V(1).Info("backup tagged with 7d retention")
	default:
		tags = append(tags, types.Tag{Key: aws.String("Retention"), Value: aws.String("24h")})
		logger.V(1).Info("backup tagged with 24h retention")
	}

	_, err = client.PutObjectTagging(ctx, &s3.PutObjectTaggingInput{
		Bucket:  &br.S3Bucket,
		Key:     aws.String(br.BackupFileS3Path()),
		Tagging: &types.Tagging{TagSet: tags},
	})
	if err != nil {
		return err
	}

	return nil
}
