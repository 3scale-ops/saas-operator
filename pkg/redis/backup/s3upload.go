package backup

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/3scale/saas-operator/pkg/ssh"
	"github.com/3scale/saas-operator/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	backupFilePrefix    string = "redis-backup"
	backupFileExtension string = "rdb"
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

	var awsBaseCommand string
	if br.AWSS3Endpoint != nil {
		awsBaseCommand = strings.Join([]string{"aws", "--endpoint-url", *br.AWSS3Endpoint}, " ")
	} else {
		awsBaseCommand = "aws"
	}

	var commands = []ssh.Runnable{
		// mv /data/dump.rdb /data/redis-backup-<shard>-<server>-<timestamp>.rdb
		ssh.NewCommand(fmt.Sprintf("mv %s %s/%s",
			br.RedisDBFile,
			path.Dir(br.RedisDBFile), br.BackupFile(),
		)),
		// gzip /data/redis-backup-<shard>-<server>-<timestamp>.rdb
		// ssh.NewCommand(fmt.Sprintf("gzip -1 %s/%s", path.Dir(br.RedisDBFile), br.BackupFile())),
		ssh.NewCommand(fmt.Sprintf("gzip -1 %s/%s", path.Dir(br.RedisDBFile), br.BackupFile())),
		// AWS_ACCESS_KEY_ID=*** AWS_SECRET_ACCESS_KEY=*** s3cmd put /data/redis-backup-<shard>-<server>-<timestamp>.rdb s3://<bucket>/<path>/redis-backup-<shard>-<server>-<timestamp>.rdb
		ssh.NewScript(
			fmt.Sprintf("%s=%s %s=%s %s=%s sh -s", util.AWSRegionEnvvar, br.AWSRegion, util.AWSAccessKeyEnvvar, br.AWSAccessKeyID, util.AWSSecretKeyEnvvar, br.AWSSecretAccessKey),
			fmt.Sprintf("%s s3 cp --only-show-errors %s/%s s3://%s/%s/%s",
				awsBaseCommand,
				path.Dir(br.RedisDBFile), br.BackupFileCompressed(),
				br.S3Bucket, br.S3Path, br.BackupFileCompressed(),
			)),
		ssh.NewCommand(fmt.Sprintf("rm -f %s/%s*", path.Dir(br.RedisDBFile), br.BackupFileBaseName())),
	}

	remoteExec := ssh.RemoteExecutor{
		Host:       br.Server.GetHost(),
		User:       br.SSHUser,
		Port:       br.SSHPort,
		PrivateKey: br.SSHKey,
		Logger:     logger,
		CmdTimeout: 0,
		Commands:   commands,
	}

	err := remoteExec.Run()
	if err != nil {
		return err
	}

	// for _, command := range commands {
	// 	if br.SSHSudo {
	// 		command = "sudo " + command
	// 	}
	// 	logger.V(1).Info(br.hideSensitive(fmt.Sprintf("running command '%s' on %s:%d", command, br.Server.GetHost(), br.SSHPort)))
	// 	output, err := remoteRun(ctx, br.SSHUser, br.Server.GetHost(), strconv.Itoa(int(br.SSHPort)), br.SSHKey, command)
	// 	if output != "" {
	// 		logger.V(1).Info(fmt.Sprintf("remote ssh command output: %s", output))
	// 	}
	// 	if err != nil {
	// 		logger.V(1).Info(fmt.Sprintf("remote ssh command error: %s", err.Error()))
	// 		return fmt.Errorf("remote ssh command failed: %w (%s)", err, output)
	// 	}
	// }

	return nil
}
