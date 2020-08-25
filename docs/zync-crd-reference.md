# Zync Custom Resource Reference

## Simple CR Example

```yaml
apiVersion: saas.3scale.net/v1alpha1
kind: Zync
metadata:
  name: simple-example
spec:
  image:
    tag: new-feature
  secret:
    vaultPath: secret/data/openshift/dev-example-4-3/3scale-zync
  zync:
    replicas: 1
    resources:
      limits:
        cpu: "1"
        memory: "1G"
```

## Full CR Example

Most of the fields do not need to be specified (can use default values), this is just an example of everything that can be overriden under your own risk:

```yaml
apiVersion: saas.3scale.net/v1alpha1
kind: Zync
metadata:
  name: full-example
spec:
  image:
    name: quay.io/3scale/zync
    tag: nightly
    pullSecretName: quay-pull-secret
  secret:
    vaultPath: secret/data/openshift/dev-example-4-3/3scale-zync
  zync:
    pdb:
      enabled: true
      maxUnavailable: "1"
    hpa:
      enabled: true
      minReplicas: 2
      maxReplicas: 4
      resourceName: cpu
      resourceUtilization: 90
    env:
      railsEnv: development
    resources:
      requests:
        cpu: "300m"
        memory: "128Mi"
      limits:
        cpu: "500m"
        memory: "196Mi"
    livenessProbe:
      initialDelaySeconds: 25
      timeoutSeconds: 2
      periodSeconds: 20
      successThreshold: 1
      failureThreshold: 5
    readinessProbe:
      initialDelaySeconds: 25
      timeoutSeconds: 2
      periodSeconds: 20
      successThreshold: 1
      failureThreshold: 5
  que:
    pdb:
      enabled: true
      minAvailable: "80%"
    hpa:
      enabled: true
      minReplicas: 2
      maxReplicas: 4
      resourceName: cpu
      resourceUtilization: 90
    env:
      railsEnv: development
    resources:
      requests:
        cpu: "250m"
        memory: "256Mi"
      limits:
        cpu: "300m"
        memory: "320Mi"
    livenessProbe:
      initialDelaySeconds: 25
      timeoutSeconds: 2
      periodSeconds: 20
      successThreshold: 1
      failureThreshold: 5
    readinessProbe:
      initialDelaySeconds: 60
      timeoutSeconds: 2
      periodSeconds: 20
      successThreshold: 1
      failureThreshold: 3
  grafanaDashboard:
    label:
      key: discovery
      value: enabled
```

## CR Spec

|                 **Field**                 | **Type**  | **Required** |   **Default value**   |                                      **Description**                                      |
| :---------------------------------------: | :-------: | :----------: | :-------------------: | :---------------------------------------------------------------------------------------: |
|               `image.name`                | `string`  |      No      | `quay.io/3scale/zync` |                          Image name (docker repository) for zync                          |
|                `image.tag`                | `string`  |      No      |       `nightly`       |                                    Image tag for zync                                     |
|          `image.pullSecretName`           | `string`  |      No      |           -           |                 Pull secret for private container repository if required                  |
|            `secret.vaultPath`             | `string`  |     Yes      |           -           |                             Vault path with the zync secrets                              |
|            `zync.pdb.enabled`             | `boolean` |      No      |        `true`         |                 Enable (`true`) or disable (`false`) PodDisruptionBudget                  |
|         `zync.pdb.maxUnavailable`         | `string`  |      No      |          `1`          |             Maximum number of unavailable pods (number or percentage of pods)             |
|          `zync.pdb.minAvailable`          | `string`  |      No      |           -           | Minimum number of available pods (number or percentage of pods), overrides maxUnavailable |
|            `zync.hpa.enabled`             | `boolean` |      No      |        `true`         |               Enable (`true`) or disable (`false`) HoritzontalPodAutoscaler               |
|          `zync.hpa.minReplicas`           |   `int`   |      No      |          `2`          |                                Minimum number of replicas                                 |
|          `zync.hpa.maxReplicas`           |   `int`   |      No      |          `4`          |                                Maximum number of replicas                                 |
|          `zync.hpa.resourceName`          | `string`  |      No      |         `cpu`         |                         Resource used for autoscale (cpu/memory)                          |
|      `zync.hpa.resourceUtilization`       |   `int`   |      No      |         `90`          |                    Percentage usage of the resource used for autoscale                    |
|              `zync.replicas`              |   `int`   |      No      |          `2`          |                  Number of replicas for zync (ignored if hpa is enabled)                  |
|            `zync.env.railsEnv`            | `string`  |      No      |     `development`     |                 Rails environment for zync (test/development/production)                  |
|       `zync.resources.requests.cpu`       | `string`  |      No      |        `250m`         |                              Override CPU requests for zync                               |
|     `zync.resources.requests.memory`      | `string`  |      No      |        `250Mi`        |                             Override Memory requests for zync                             |
|        `zync.resources.limits.cpu`        | `string`  |      No      |        `750m`         |                               Override CPU limits for zync                                |
|      `zync.resources.limits.memory`       | `string`  |      No      |        `512Mi`        |                              Override Memory limits for zync                              |
| `zync.livenessProbe.initialDelaySeconds`  |   `int`   |      No      |         `10`          |                    Override liveness initial delay (seconds) for zync                     |
|    `zync.livenessProbe.timeoutSeconds`    |   `int`   |      No      |         `30`          |                       Override liveness timeout (seconds) for zync                        |
|    `zync.livenessProbe.periodSeconds`     |   `int`   |      No      |         `10`          |                        Override liveness period (seconds) for zync                        |
|   `zync.livenessProbe.successThreshold`   |   `int`   |      No      |          `1`          |                       Override liveness success threshold for zync                        |
|   `zync.livenessProbe.failureThreshold`   |   `int`   |      No      |          `3`          |                       Override liveness failure threshold for zync                        |
| `zync.readinessProbe.initialDelaySeconds` |   `int`   |      No      |         `30`          |                    Override readiness initial delay (seconds) for zync                    |
|   `zync.readinessProbe.timeoutSeconds`    |   `int`   |      No      |         `10`          |                       Override readiness timeout (seconds) for zync                       |
|    `zync.readinessProbe.periodSeconds`    |   `int`   |      No      |         `10`          |                       Override readiness period (seconds) for zync                        |
|  `zync.readinessProbe.successThreshold`   |   `int`   |      No      |          `1`          |                       Override readiness success threshold for zync                       |
|  `zync.readinessProbe.failureThreshold`   |   `int`   |      No      |          `3`          |                       Override readiness failure threshold for zync                       |
|             `que.pdb.enabled`             | `boolean` |      No      |        `true`         |                 Enable (`true`) or disable (`false`) PodDisruptionBudget                  |
|         `que.pdb.maxUnavailable`          | `string`  |      No      |          `1`          |             Maximum number of unavailable pods (number or percentage of pods)             |
|          `que.pdb.minAvailable`           | `string`  |      No      |           -           | Minimum number of available pods (number or percentage of pods), overrides maxUnavailable |
|             `que.hpa.enabled`             | `boolean` |      No      |        `true`         |               Enable (`true`) or disable (`false`) HoritzontalPodAutoscaler               |
|           `que.hpa.minReplicas`           |   `int`   |      No      |          `2`          |                                Minimum number of replicas                                 |
|           `que.hpa.maxReplicas`           |   `int`   |      No      |          `4`          |                                Maximum number of replicas                                 |
|          `que.hpa.resourceName`           | `string`  |      No      |         `cpu`         |                         Resource used for autoscale (cpu/memory)                          |
|       `que.hpa.resourceUtilization`       |   `int`   |      No      |         `90`          |                    Percentage usage of the resource used for autoscale                    |
|              `que.replicas`               |   `int`   |      No      |          `2`          |                Number of replicas for zync-que (ignored if hpa is enabled)                |
|            `que.env.railsEnv`             | `string`  |      No      |     `development`     |               Rails environment for zync-que (test/development/production)                |
|       `que.resources.requests.cpu`        | `string`  |      No      |        `250m`         |                            Override CPU requests for zync-que                             |
|      `que.resources.requests.memory`      | `string`  |      No      |        `250Mi`        |                           Override Memory requests for zync-que                           |
|        `que.resources.limits.cpu`         | `string`  |      No      |        `750m`         |                             Override CPU limits for zync-que                              |
|       `que.resources.limits.memory`       | `string`  |      No      |        `512Mi`        |                            Override Memory limits for zync-que                            |
|  `que.livenessProbe.initialDelaySeconds`  |   `int`   |      No      |         `10`          |                  Override liveness initial delay (seconds) for zync-que                   |
|    `que.livenessProbe.timeoutSeconds`     |   `int`   |      No      |         `30`          |                     Override liveness timeout (seconds) for zync-que                      |
|     `que.livenessProbe.periodSeconds`     |   `int`   |      No      |         `10`          |                      Override liveness period (seconds) for zync-que                      |
|   `que.livenessProbe.successThreshold`    |   `int`   |      No      |          `1`          |                     Override liveness success threshold for zync-que                      |
|   `que.livenessProbe.failureThreshold`    |   `int`   |      No      |          `3`          |                     Override liveness failure threshold for zync-que                      |
| `que.readinessProbe.initialDelaySeconds`  |   `int`   |      No      |         `30`          |                  Override readiness initial delay (seconds) for zync-que                  |
|    `que.readinessProbe.timeoutSeconds`    |   `int`   |      No      |         `10`          |                     Override readiness timeout (seconds) for zync-que                     |
|    `que.readinessProbe.periodSeconds`     |   `int`   |      No      |         `10`          |                     Override readiness period (seconds) for zync-que                      |
|   `que.readinessProbe.successThreshold`   |   `int`   |      No      |          `1`          |                     Override readiness success threshold for zync-que                     |
|   `que.readinessProbe.failureThreshold`   |   `int`   |      No      |          `3`          |                     Override readiness failure threshold for zync-que                     |
|       `grafanaDashboard.label.key`        | `string`  |      No      |   `monitoring-key`    |               Label `key` used by grafana-operator for dashboard discovery                |
|      `grafanaDashboard.label.value`       | `string`  |      No      |     `middleware`      |              Label `value` used by grafana-operator for dashboard discovery               |
