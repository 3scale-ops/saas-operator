package backup

import (
	"context"
	"fmt"
	"net"
	"path"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	backupFilePrefix    string = "redis-backup"
	backupFileExtension string = "rdb"
	sshKeyFile          string = "ssh-private-key"
	awsAccessKeyEnvvar  string = "AWS_ACCESS_KEY_ID"
	awsSecretKeyEnvvar  string = "AWS_SECRET_ACCESS_KEY"
	awsRegionEnvvar     string = "AWS_REGION"
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

	var commands = []string{
		// mv /data/dump.rdb /data/redis-backup-<shard>-<server>-<timestamp>.rdb
		fmt.Sprintf("mv %s %s/%s",
			br.RedisDBFile,
			path.Dir(br.RedisDBFile), br.BackupFile(),
		),
		// gzip /data/redis-backup-<shard>-<server>-<timestamp>.rdb
		fmt.Sprintf("gzip %s/%s", path.Dir(br.RedisDBFile), br.BackupFile()),
		// TODO: use awscli instead
		// AWS_ACCESS_KEY_ID=*** AWS_SECRET_ACCESS_KEY=*** s3cmd put /data/redis-backup-<shard>-<server>-<timestamp>.rdb s3://<bucket>/<path>/redis-backup-<shard>-<server>-<timestamp>.rdb
		fmt.Sprintf("%s=%s %s=%s %s=%s s3 cp %s/%s s3://%s/%s/%s",
			awsRegionEnvvar, br.AWSRegion,
			awsAccessKeyEnvvar, br.AWSAccessKeyID,
			awsSecretKeyEnvvar, br.AWSSecretAccessKey,
			path.Dir(br.RedisDBFile), br.BackupFileCompressed(),
			br.S3Bucket, br.S3Path, br.BackupFileCompressed(),
		),
	}

	for _, command := range commands {
		logger.V(1).Info(fmt.Sprintf("running command '%s' on %s:%d", command, br.Server.GetHost(), br.SSHPort))
		output, err := remoteRun(ctx, br.SSHUser, br.Server.GetHost(), strconv.Itoa(int(br.SSHPort)), br.SSHKey, command)
		if output != "" {
			logger.V(1).Info(fmt.Sprintf("remote ssh command output: %s", output))
		}
		if err != nil {
			logger.V(1).Info(fmt.Sprintf("remote ssh command error: %s", err.Error()))
			return fmt.Errorf("remote ssh command failed: %w (%s)", err, output)
		}
	}

	return nil
}

// e.g. output, err := remoteRun(ctx, "root", "MY_IP", "MY_PORT", "PRIVATE_KEY", "ls")
func remoteRun(ctx context.Context, user, addr, port, privateKey, cmd string) (string, error) {

	key, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return "", err
	}
	// Authentication
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}
	client, err := ssh.Dial("tcp", net.JoinHostPort(addr, port), config)
	if err != nil {
		return "", err
	}
	// Create a session. It is one session per command.
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	return string(output), err
}
