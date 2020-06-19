# Zync Custom Resource Reference

## Full CR Example

Most of the fields do not need to be specified (can use default values), this is just an example of everything that can be overriden under your own risk:

```yaml
apiVersion: saas.3scale.net/v1alpha1
kind: Zync
metadata:
  name: example
spec:
  image:
    name: quay.io/3scale/zync
    tag: nightly
    pullSecretName: quay-pull-secret
  secret:
    zyncDatabaseVaultPath: secret/data/openshift/stg-saas-ocp/stg-saas-3scale-zync
  zync:
    replicas: 3
    env:
      dbWaitSleepSeconds: 10
      railsEnv: "staging"
      railsLogsToStdout: "true"
    resources:
      requests:
        cpu: 300m
        memory: 128Mi
      limits:
        cpu: 500m
        memory: 196Mi
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
    replicas: 2
    env:
      railsEnv: "staging"
      railsLogsToStdout: "true"
    resources:
      requests:
        cpu: 250m
        memory: 256Mi
      limits:
        cpu: 300m
        memory: 320Mi
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
    zyncDatabaseVaultPath: secret/data/openshift/stg-saas-ocp/stg-saas-3scale-zync
  zync:
    replicas: 1
    env:
      dbWaitSleepSeconds: 60
    resources:
      limits:
        cpu: 1
        memory: 1G
```

## CR Spec

|                 **Field**                 | **Type** | **Required** |   **Default value**   |                        **Description**                         |
| :---------------------------------------: | :------: | :----------: | :-------------------: | :------------------------------------------------------------: |
|             `zync.image.name`             | `string` |      No      | `quay.io/3scale/zync` |            Image name (docker repository) for zync             |
|             `zync.image.tag`              | `string` |      No      |       `nightly`       |                       Image tag for zync                       |
|        `zync.image.pullSecretName`        | `string` |      No      |  `quay-pull-secret`   |        Quay pull secret for private repository for zync        |
|    `zync.secret.zyncDatabaseVaultPath`    | `string` |     Yes      |           -           |           Vault path with the zync database secrets            |
|  `zync.env.zync.env.dbWaitSleepSeconds `  |  `int`   |      No      |         `30`          |        Sleep delay while waiting for the zync database         |
|            `zync.env.railsEnv`            | `string` |      No      |       `staging`       |                   Rails environment for zync                   |
|       `zync.env.railsLogsToStdout`        | `string` |      No      |        `false`        |            Rails log to std output toggle for zync             |
|              `zync.replicas`              |  `int`   |      No      |          `3`          |                  Number of replicas for zync                   |
|       `zync.resources.requests.cpu`       | `string` |      No      |        `250m`         |                 Override CPU requests for zync                 |
|     `zync.resources.requests.memory`      | `string` |      No      |        `250Mi`        |               Override Memory requests for zync                |
|        `zync.resources.limits.cpu`        | `string` |      No      |        `750m`         |                  Override CPU limits for zync                  |
|      `zync.resources.limits.memory`       | `string` |      No      |        `512Mi`        |                Override Memory limits for zync                 |
| `zync.livenessProbe.initialDelaySeconds`  |  `int`   |      No      |         `10`          |       Override liveness initial delay (seconds) for zync       |
|    `zync.livenessProbe.timeoutSeconds`    |  `int`   |      No      |         `30`          |          Override liveness timeout (seconds) for zync          |
|    `zync.livenessProbe.periodSeconds`     |  `int`   |      No      |         `10`          |          Override liveness period (seconds) for zync           |
|   `zync.livenessProbe.successThreshold`   |  `int`   |      No      |          `1`          |          Override liveness success threshold for zync          |
|   `zync.livenessProbe.failureThreshold`   |  `int`   |      No      |          `3`          |          Override liveness failure threshold for zync          |
| `zync.readinessProbe.initialDelaySeconds` |  `int`   |      No      |         `30`          |      Override readiness initial delay (seconds) for zync       |
|   `zync.readinessProbe.timeoutSeconds`    |  `int`   |      No      |         `10`          |         Override readiness timeout (seconds) for zync          |
|    `zync.readinessProbe.periodSeconds`    |  `int`   |      No      |         `10`          |          Override readiness period (seconds) for zync          |
|  `zync.readinessProbe.successThreshold`   |  `int`   |      No      |          `1`          |         Override readiness success threshold for zync          |
|  `zync.readinessProbe.failureThreshold`   |  `int`   |      No      |          `3`          |         Override readiness failure threshold for zync          |
|             `que.image.name`              | `string` |      No      | `quay.io/3scale/zync` |          Image name (docker repository) for zync-que           |
|              `que.image.tag`              | `string` |      No      |       `nightly`       |                     Image tag for zync-que                     |
|        `que.image.pullSecretName`         | `string` |      No      |  `quay-pull-secret`   |      Quay pull secret for private repository for zync-que      |
|              `que.replicas`               |  `int`   |      No      |          `3`          |                Number of replicas for zync-que                 |
|            `que.env.railsEnv`             | `string` |      No      |       `staging`       |                 Rails environment for zync-que                 |
|        `que.env.railsLogsToStdout`        | `string` |      No      |        `false`        |           Rail log to std output toggle for zync-que           |
|       `que.resources.requests.cpu`        | `string` |      No      |        `250m`         |               Override CPU requests for zync-que               |
|      `que.resources.requests.memory`      | `string` |      No      |        `250Mi`        |             Override Memory requests for zync-que              |
|        `que.resources.limits.cpu`         | `string` |      No      |        `750m`         |                Override CPU limits for zync-que                |
|       `que.resources.limits.memory`       | `string` |      No      |        `512Mi`        |              Override Memory limits for zync-que               |
|  `que.livenessProbe.initialDelaySeconds`  |  `int`   |      No      |         `10`          |     Override liveness initial delay (seconds) for zync-que     |
|    `que.livenessProbe.timeoutSeconds`     |  `int`   |      No      |         `30`          |        Override liveness timeout (seconds) for zync-que        |
|     `que.livenessProbe.periodSeconds`     |  `int`   |      No      |         `10`          |        Override liveness period (seconds) for zync-que         |
|   `que.livenessProbe.successThreshold`    |  `int`   |      No      |          `1`          |        Override liveness success threshold for zync-que        |
|   `que.livenessProbe.failureThreshold`    |  `int`   |      No      |          `3`          |        Override liveness failure threshold for zync-que        |
| `que.readinessProbe.initialDelaySeconds`  |  `int`   |      No      |         `30`          |    Override readiness initial delay (seconds) for zync-que     |
|    `que.readinessProbe.timeoutSeconds`    |  `int`   |      No      |         `10`          |       Override readiness timeout (seconds) for zync-que        |
|    `que.readinessProbe.periodSeconds`     |  `int`   |      No      |         `10`          |        Override readiness period (seconds) for zync-que        |
|   `que.readinessProbe.successThreshold`   |  `int`   |      No      |          `1`          |       Override readiness success threshold for zync-que        |
|   `que.readinessProbe.failureThreshold`   |  `int`   |      No      |          `3`          |       Override readiness failure threshold for zync-que        |
|       `grafanaDashboard.label.key`        | `string` |      No      |   `monitoring-key`    |  Label `key` used by grafana-operator for dashboard discovery  |
|      `grafanaDashboard.label.value`       | `string` |      No      |     `middleware`      | Label `value` used by grafana-operator for dashboard discovery |
