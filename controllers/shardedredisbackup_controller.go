/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/3scale-ops/basereconciler/reconciler"
	"github.com/3scale-ops/basereconciler/util"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/reconcilers/threads"
	"github.com/3scale-ops/saas-operator/pkg/redis/backup"
	redis "github.com/3scale-ops/saas-operator/pkg/redis/server"
	"github.com/3scale-ops/saas-operator/pkg/redis/sharded"
	operatorutils "github.com/3scale-ops/saas-operator/pkg/util"
	"github.com/go-logr/logr"
	"github.com/robfig/cron/v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// ShardedRedisBackupReconciler reconciles a ShardedRedisBackup object
type ShardedRedisBackupReconciler struct {
	*reconciler.Reconciler
	Log          logr.Logger
	BackupRunner threads.Manager
	Pool         *redis.ServerPool
}

//+kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=shardedredisbackups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=shardedredisbackups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=saas.3scale.net,namespace=placeholder,resources=shardedredisbackups/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ShardedRedisBackup object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *ShardedRedisBackupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("name", req.Name, "namespace", req.Namespace)
	ctx = log.IntoContext(ctx, logger)
	now := time.Now()

	// ----------------------------------
	// ----- Phase 1: get instances -----
	// ----------------------------------

	instance := &saasv1alpha1.ShardedRedisBackup{}
	result := r.ManageResourceLifecycle(ctx, req, instance,
		reconciler.WithInMemoryInitializationFunc(util.ResourceDefaulter(instance)),
		reconciler.WithFinalizer(saasv1alpha1.Finalizer),
		reconciler.WithFinalizationFunc(r.BackupRunner.CleanupThreads(instance)),
	)
	if result.ShouldReturn() {
		return result.Values()
	}

	// Get Sentinel status
	sentinel := &saasv1alpha1.Sentinel{ObjectMeta: metav1.ObjectMeta{Name: instance.Spec.SentinelRef, Namespace: req.Namespace}}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(sentinel), sentinel); err != nil {
		return ctrl.Result{}, err
	}

	cluster, err := sentinel.Status.ShardedCluster(ctx, r.Pool)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Get SSH key
	sshPrivateKey := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{
		Name: instance.Spec.SSHOptions.PrivateKeySecretRef.Name, Namespace: req.Namespace}}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(sshPrivateKey), sshPrivateKey); err != nil {
		return ctrl.Result{}, err
	}
	if sshPrivateKey.Type != corev1.SecretTypeSSHAuth {
		return ctrl.Result{}, fmt.Errorf("secret %s must be of 'kubernetes.io/ssh-auth' type", sshPrivateKey.GetName())
	}
	if _, ok := sshPrivateKey.Data[corev1.SSHAuthPrivateKey]; !ok {
		return ctrl.Result{}, fmt.Errorf("secret %s is missing %s key", sshPrivateKey.GetName(), corev1.SSHAuthPrivateKey)
	}

	// Get AWS credentials
	awsCredentials := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{
		Name: instance.Spec.S3Options.CredentialsSecretRef.Name, Namespace: req.Namespace}}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(awsCredentials), awsCredentials); err != nil {
		return ctrl.Result{}, err
	}
	if _, ok := awsCredentials.Data[operatorutils.AWSAccessKeyEnvvar]; !ok {
		return ctrl.Result{}, fmt.Errorf("secret %s is missing %s key", awsCredentials.GetName(), operatorutils.AWSAccessKeyEnvvar)
	}
	if _, ok := awsCredentials.Data[operatorutils.AWSSecretKeyEnvvar]; !ok {
		return ctrl.Result{}, fmt.Errorf("secret %s is missing %s key", awsCredentials.GetName(), operatorutils.AWSSecretKeyEnvvar)
	}

	// ----------------------------------------
	// ----- Phase 2: run pending backups -----
	// ----------------------------------------

	statusChanged := false
	requeue := false
	runners := make([]threads.RunnableThread, 0, len(cluster.Shards))
	for _, shard := range cluster.Shards {
		scheduledBackup, _ := instance.Status.FindLastBackup(shard.Name, saasv1alpha1.BackupPendingState)
		if scheduledBackup != nil && scheduledBackup.ScheduledFor.Time.Before(now) {
			// hanlde error when no available RO slaves
			var roSlaves []*sharded.RedisServer
			if roSlaves = shard.GetSlavesRO(); len(roSlaves) == 0 {
				logger.Error(fmt.Errorf("no available RO slaves in shard"), fmt.Sprintf("skipped shard %s, will be retried", shard.Name))
				requeue = true
				continue
			}

			// add the backup runner thread
			runners = append(runners, &backup.Runner{
				ShardName:          shard.Name,
				Server:             roSlaves[0],
				ScheduledFor:       scheduledBackup.ScheduledFor.Time,
				Timestamp:          now,
				Timeout:            instance.Spec.Timeout.Duration,
				PollInterval:       instance.Spec.PollInterval.Duration,
				RedisDBFile:        instance.Spec.DBFile,
				Instance:           instance,
				SSHUser:            instance.Spec.SSHOptions.User,
				SSHKey:             string(sshPrivateKey.Data[corev1.SSHAuthPrivateKey]),
				SSHPort:            *instance.Spec.SSHOptions.Port,
				SSHSudo:            *instance.Spec.SSHOptions.Sudo,
				S3Bucket:           instance.Spec.S3Options.Bucket,
				S3Path:             instance.Spec.S3Options.Path,
				AWSAccessKeyID:     string(awsCredentials.Data[saasv1alpha1.AWSAccessKeyID_SecretKey]),
				AWSSecretAccessKey: string(awsCredentials.Data[saasv1alpha1.AWSSecretAccessKey_SecretKey]),
				AWSRegion:          instance.Spec.S3Options.Region,
				AWSS3Endpoint:      instance.Spec.S3Options.ServiceEndpoint,
			})
			scheduledBackup.ServerAlias = util.Pointer(roSlaves[0].GetAlias())
			scheduledBackup.ServerID = util.Pointer(roSlaves[0].ID())
			scheduledBackup.StartedAt = &metav1.Time{Time: now}
			scheduledBackup.Message = "backup is running"
			scheduledBackup.State = saasv1alpha1.BackupRunningState
			statusChanged = true
		}
	}

	if err := r.BackupRunner.ReconcileThreads(ctx, instance, runners, logger.WithName("backup-runner")); err != nil {
		return ctrl.Result{}, err
	}
	if statusChanged {
		err := r.Client.Status().Update(ctx, instance)
		return ctrl.Result{}, err
	}
	// requeue if any of the shards had no available RO slaves
	if requeue {
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	// --------------------------------------------------------
	// ----- Phase 3: reconcile status of running backups -----
	// --------------------------------------------------------

	for _, b := range instance.Status.GetRunningBackups() {
		var thread *backup.Runner
		var srv *sharded.RedisServer

		if srv = cluster.LookupServerByID(*b.ServerID); srv == nil {
			b.State = saasv1alpha1.BackupUnknownState
			b.Message = "server not found in cluster"
			statusChanged = true
			continue
		}

		if t := r.BackupRunner.GetThread(backup.ID(b.Shard, srv.GetAlias(), b.ScheduledFor.Time), instance, logger); t != nil {
			thread = t.(*backup.Runner)
		} else {
			b.State = saasv1alpha1.BackupUnknownState
			b.Message = "runner not found"
			statusChanged = true
			continue
		}

		if status := thread.Status(); status.Finished {
			if err := status.Error; err != nil {
				b.State = saasv1alpha1.BackupFailedState
				b.Message = err.Error()
			} else {
				b.State = saasv1alpha1.BackupCompletedState
				b.Message = "backup complete"
				b.BackupFile = &status.BackupFile
				b.BackupSize = &status.BackupSize
				b.FinishedAt = &metav1.Time{Time: status.FinishedAt}
			}
			statusChanged = true
		}
	}

	if statusChanged {
		err := r.Client.Status().Update(ctx, instance)
		return ctrl.Result{}, err
	}

	// -------------------------------------
	// ----- Phase 4: schedule backups -----
	// -------------------------------------

	schedule, err := cron.ParseStandard(instance.Spec.Schedule)
	if err != nil {
		return ctrl.Result{}, err
	}
	nextRun := schedule.Next(now)

	// only actually add the schedule if pause == false
	if !*instance.Spec.Pause {
		statusChanged, err = r.reconcileBackupList(ctx, instance, nextRun, cluster.GetShardNames())
		if err != nil {
			return reconcile.Result{}, err
		}

		if statusChanged {
			err := r.Client.Status().Update(ctx, instance)
			return ctrl.Result{}, err
		}
	}

	// requeue for next schedule
	return ctrl.Result{RequeueAfter: time.Until(nextRun.Add(1 * time.Second))}, nil
}

func (r *ShardedRedisBackupReconciler) reconcileBackupList(ctx context.Context, instance *saasv1alpha1.ShardedRedisBackup, nextRun time.Time, shards []string) (bool, error) {
	logger := log.FromContext(ctx, "function", "(r *ShardedRedisBackupReconciler) reconcileBackupList")
	changed := false

	for _, shard := range shards {
		// don't schedule if a backup is already running
		if runningbackup, _ := instance.Status.FindLastBackup(shard, saasv1alpha1.BackupRunningState); runningbackup != nil {
			continue
		}
		if lastbackup, pos := instance.Status.FindLastBackup(shard, saasv1alpha1.BackupPendingState); lastbackup != nil {
			// found a pending backup for this shard
			if nextRun == lastbackup.ScheduledFor.Time {
				// already scheduled, do nothing
				continue
			} else {
				// already scheduled for a different time, replace with new schedule
				instance.Status.DeleteBackup(pos)
			}
		}

		// add a new backup to the list
		instance.Status.AddBackup(saasv1alpha1.BackupStatus{
			Shard:        shard,
			ScheduledFor: metav1.NewTime(nextRun),
			Message:      "backup scheduled",
			State:        saasv1alpha1.BackupPendingState,
		})
		logger.V(1).Info("scheduled backup", "shard", shard, "scheduledFor", nextRun)
		changed = true
	}

	// apply historyLimit
	if instance.Status.ApplyHistoryLimit(*instance.Spec.HistoryLimit, shards) {
		changed = true
	}

	return changed, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ShardedRedisBackupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&saasv1alpha1.ShardedRedisBackup{}).
		Watches(&source.Channel{Source: r.BackupRunner.GetChannel()}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
