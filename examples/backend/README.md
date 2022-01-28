# NOTES

* Create a cluster

```bash
kubectl kind-create
kubectl kind-deploy
export KUBECONFIG=$PWD/kubeconfig
```

* Deploy first the RedisShards for redis-storage and redis-resque

```bash
kubectl apply -f examples/backend/redis.yaml
```

* Wait until the redis Pods are ready, otherwise the sentinel controller will start giving a lot of errors and will end up taking a very long time to properly reconcile due to controller reconcile throttling. Once the redis Pods are ready, apply the rest of the manifests.

```bash
kubectl apply -f examples/backend/
```

* In the end you should have something like this:

```bash
❯ kubectl get pods
NAME                                                READY   STATUS    RESTARTS   AGE
backend-cron-c89c5f5d5-s4xqx                        1/1     Running   0          101s
backend-listener-7684f5dbd8-vgf96                   2/2     Running   0          102s
backend-worker-6fb65b7c6b-2gwn5                     2/2     Running   0          101s
redis-sentinel-0                                    1/1     Running   0          108s
redis-sentinel-1                                    1/1     Running   0          108s
redis-sentinel-2                                    1/1     Running   0          107s
redis-shard-resque-0                                1/1     Running   0          3m32s
redis-shard-shard01-0                               1/1     Running   0          3m37s
redis-shard-shard01-1                               1/1     Running   0          3m37s
redis-shard-shard01-2                               1/1     Running   0          3m37s
redis-shard-shard02-0                               1/1     Running   0          3m36s
redis-shard-shard02-1                               1/1     Running   0          3m36s
redis-shard-shard02-2                               1/1     Running   0          3m36s
saas-operator-controller-manager-57474ff9c4-fwlss   1/1     Running   1          46m
```
