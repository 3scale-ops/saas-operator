package backup

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"text/template"
	"time"

	"github.com/3scale-ops/saas-operator/pkg/ssh"
	operatorutils "github.com/3scale-ops/saas-operator/pkg/util"
	"github.com/MakeNowJust/heredoc"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	backupFilePrefix    string = "redis-backup"
	backupFileExtension string = "rdb"
)

type Retention string

const (
	Retention90d Retention = "90d"
	Retention7d  Retention = "7d"
	Retention24h Retention = "24h"
)

func (br *Runner) BackupFileBaseName() string {
	return fmt.Sprintf("%s_%s", backupFilePrefix, br.ShardName)
}

func (br *Runner) BackupFileBaseNameWithTimeSuffix(timeSuffix string) string {
	return fmt.Sprintf("%s_%s", br.BackupFileBaseName(), timeSuffix)
}

// BackupFile returns the backup file as "redis-backup-<shard>-<server>-<timestamp>.rdb"
func (br *Runner) BackupFile() string {
	return fmt.Sprintf("%s.%s",
		br.BackupFileBaseNameWithTimeSuffix(br.Timestamp.Format(time.RFC3339)),
		backupFileExtension)
}

func (br *Runner) BackupFileCompressed() string {
	return fmt.Sprintf("%s.gz", br.BackupFile())
}

func (br *Runner) BackupFileS3Path() string {
	return fmt.Sprintf("%s/%s", br.S3Path, br.BackupFileCompressed())
}

func (br *Runner) UploadBackup(ctx context.Context) error {
	logger := log.FromContext(ctx, "function", "(br *Runner) UploadBackup()")

	uploadScript, err := br.uploadScript(ctx)
	if err != nil {
		return err
	}

	remoteExec := ssh.RemoteExecutor{
		Host:       br.Server.GetHost(),
		User:       br.SSHUser,
		Port:       br.SSHPort,
		PrivateKey: br.SSHKey,
		Logger:     logger,
		CmdTimeout: 0,
		Commands: []ssh.Runnable{
			ssh.NewCommand(fmt.Sprintf("mv %s %s/%s", br.RedisDBFile, path.Dir(br.RedisDBFile), br.BackupFile())),
			ssh.NewCommand(fmt.Sprintf("gzip -1 %s/%s", path.Dir(br.RedisDBFile), br.BackupFile())),
			ssh.NewScript(fmt.Sprintf("%s=%s %s=%s %s=%s python -",
				operatorutils.AWSRegionEnvvar, br.AWSRegion,
				operatorutils.AWSAccessKeyEnvvar, br.AWSAccessKeyID,
				operatorutils.AWSSecretKeyEnvvar, br.AWSSecretAccessKey),
				uploadScript,
				br.AWSSecretAccessKey,
			),
			ssh.NewCommand(fmt.Sprintf("rm -f %s/%s*", path.Dir(br.RedisDBFile), br.BackupFileBaseName())),
		},
	}

	err = remoteExec.Run()
	if err != nil {
		return err
	}

	return nil
}

func (br *Runner) resolveTags(ctx context.Context) (string, error) {
	logger := log.FromContext(ctx, "function", "(br *Runner) ResolveTags()")
	var retention Retention

	awsconfig, err := operatorutils.AWSConfig(ctx, br.AWSAccessKeyID, br.AWSSecretAccessKey, br.AWSRegion, br.AWSS3Endpoint)
	if err != nil {
		return "{}", err
	}

	client := s3.NewFromConfig(*awsconfig)

	// get backups of current day
	dayResult, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(br.S3Bucket),
		Prefix: aws.String(br.S3Path + "/" + br.BackupFileBaseNameWithTimeSuffix(br.Timestamp.Format("2006-01-02"))),
	})
	if err != nil {
		return "{}", err
	}

	// get backups of current hour
	hourResult, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(br.S3Bucket),
		Prefix: aws.String(br.S3Path + "/" + br.BackupFileBaseNameWithTimeSuffix(br.Timestamp.Format("2006-01-02T15"))),
	})
	if err != nil {
		return "{}", err
	}

	if len(dayResult.Contents) == 0 {
		retention = Retention90d
		logger.V(1).Info("backup tagged with 90d retention")
	} else if len(hourResult.Contents) == 0 {
		retention = Retention7d
		logger.V(1).Info("backup tagged with 7d retention")
	} else {
		retention = Retention24h
		logger.V(1).Info("backup tagged with 24h retention")
	}

	tags := url.Values{
		"Layer":       []string{"bck-storage"},
		"App":         []string{"Backend"},
		"Shard":       []string{br.ShardName},
		"HostAddress": []string{br.Server.ID()},
		"HostAlias":   []string{br.Server.GetAlias()},
		"Retention":   []string{string(retention)},
	}

	return tags.Encode(), nil
}

func (br *Runner) uploadScript(ctx context.Context) (string, error) {
	tags, err := br.resolveTags(ctx)
	if err != nil {
		return "", err
	}

	scriptTemplate := heredoc.Doc(`
		import boto3
		session = boto3.session.Session()
		s3 = session.client(service_name="s3"{{if .Endpoint}},endpoint_url="{{.Endpoint}}"{{end}})
		s3.upload_file(
			"{{.File}}",
			"{{.Bucket}}",
			"{{.Key}}",
			ExtraArgs={"Tagging": "{{.Tags}}"},
		)
	`)

	templateVars := struct {
		File, Bucket, Key, Endpoint, Tags string
	}{
		File:   filepath.Join(path.Dir(br.RedisDBFile), br.BackupFileCompressed()),
		Bucket: br.S3Bucket,
		Key:    filepath.Join(br.S3Path, br.BackupFileCompressed()),
		Tags:   tags,
	}
	if br.AWSS3Endpoint != nil {
		templateVars.Endpoint = *br.AWSS3Endpoint
	}

	t := template.Must(template.New("script").Parse(scriptTemplate))
	script := new(bytes.Buffer)
	err = t.Execute(script, templateVars)
	if err != nil {
		return "", err
	}

	return script.String(), nil
}
