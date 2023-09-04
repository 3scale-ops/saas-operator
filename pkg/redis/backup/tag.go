package backup

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (br *Runner) TagBackup(ctx context.Context) error {
	logger := log.FromContext(ctx, "function", "(br *Runner) TagBackup()")

	// set AWS credentials
	os.Setenv(awsAccessKeyEnvvar, br.AWSAccessKeyID)
	os.Setenv(awsSecretKeyEnvvar, br.AWSSecretAccessKey)
	os.Setenv(awsRegionEnvvar, br.AWSRegion)

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
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

	// check backup size of last (given a size passed as threshold in the CR)
	if last.Size < int64(br.MinSize) {
		return fmt.Errorf("last backup %s is smaller that declared min size of %d", *last.Key, br.MinSize)
	}

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