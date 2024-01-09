package backup

import (
	"context"
	"fmt"
	"sort"

	operatorutils "github.com/3scale-ops/saas-operator/pkg/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (br *Runner) CheckBackup(ctx context.Context) error {
	logger := log.FromContext(ctx, "function", "(br *Runner) CheckBackup()")

	awsconfig, err := operatorutils.AWSConfig(ctx, br.AWSAccessKeyID, br.AWSSecretAccessKey, br.AWSRegion, br.AWSS3Endpoint)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(*awsconfig)

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

	latest := hourResult.Contents[len(hourResult.Contents)-1]
	if br.BackupFileS3Path() != *latest.Key {
		err := fmt.Errorf("latest backup %s has different key than expected (%s)", *latest.Key, br.BackupFileS3Path())
		logger.Error(err, "unable to find backup s3")
		return err
	}
	// store backup size
	br.status.BackupSize = latest.Size

	return nil
}
