# System Custom Resource Reference

## Simple CR Example

```yaml
apiVersion: saas.3scale.net/v1alpha1
kind: System
metadata:
  name: simple-example
spec:
  env:
    threescaleSuperdomain: example.com
  secret:
    multitenantIngressTLS: wildcard-multitenant-example-com-certificate
    appVaultPath: secret/data/openshift/cluster-example/3scale/system-app
    backendListenerVaultPath: secret/data/openshift/cluster-example/3scale/system-backend-listener
    configVaultPath: secret/data/openshift/dev-eng-ocp/3sscale-saas/system-config
    databaseVaultPath: secret/data/openshift/cluster-example/3scale/system-database
    eventsHookVaultPath: secret/data/openshift/cluster-example/3scale/system-events-hook
    masterApicastVaultPath: secret/data/openshift/cluster-example/3scale/system-master-apicast
    memcachedVaultPath: secret/data/openshift/cluster-example/3scale/system-memcached
    multitenantAssetsS3VaultPath: secret/data/openshift/cluster-example/3scale/system-multitenant-assets-s3
    recaptchaVaultPath: secret/data/openshift/cluster-example/3scale/system-recaptcha
    redisVaultPath: secret/data/openshift/cluster-example/3scale/system-redis
    seedVaultPath: secret/data/openshift/cluster-example/3scale/system-seed
    smtpVaultPath: secret/data/openshift/cluster-example/3scale/system-smtp
```

## Full CR Example

Most of the fields do not need to be specified (can use default values), this is just an example of everything that can be overriden under your own risk:

```yaml
apiVersion: saas.3scale.net/v1alpha1
kind: System
metadata:
  name: full-example
spec:
  image:
    name: "quay.io/3scale/porta"
    tag: "nightly"
    pullPolicy: "Always"
  grafanaDashboard:
    label:
      key: "monitoring-key"
      value: "middleware"
  env:
    ampRelease: 2.7.1
    rails:
      env: production
      logLevel: info
      logToStdout: true
    sandboxProxyOpensslVerifyMode: VERIFY_NONE
    forceSSL: true
    sslCertDir: /etc/pki/tls/cert
    threescaleProviderPlan: enterprise
    threescaleSuperdomain: example.com
  secret:
    multitenantIngressTLS: wildcard-multitenant-example-com-certificate
    appVaultPath: secret/data/openshift/cluster-example/3scale/system-app
    backendListenerVaultPath: secret/data/openshift/cluster-example/3scale/system-backend-listener
    configVaultPath: secret/data/openshift/dev-eng-ocp/3sscale-saas/system-config
    databaseVaultPath: secret/data/openshift/cluster-example/3scale/system-database
    eventsHookVaultPath: secret/data/openshift/cluster-example/3scale/system-events-hook
    masterApicastVaultPath: secret/data/openshift/cluster-example/3scale/system-master-apicast
    memcachedVaultPath: secret/data/openshift/cluster-example/3scale/system-memcached
    multitenantAssetsS3VaultPath: secret/data/openshift/cluster-example/3scale/system-multitenant-assets-s3
    recaptchaVaultPath: secret/data/openshift/cluster-example/3scale/system-recaptcha
    redisVaultPath: secret/data/openshift/cluster-example/3scale/system-redis
    seedVaultPath: secret/data/openshift/cluster-example/3scale/system-seed
    smtpVaultPath: secret/data/openshift/cluster-example/3scale/system-smtp
  app:
    replicas: 2
    resources:
      requests:
        cpu: "200m"
        memory: "1Gi"
      limits:
        cpu: "400m"
        memory: "2Gi"
    livenessProbe:
      initialDelaySeconds: 30
      timeoutSeconds: 1
      periodSeconds: 10
      successThreshold: 1
      failureThreshold: 3
    readinessProbe:
      initialDelaySeconds: 30
      timeoutSeconds: 5
      periodSeconds: 10
      successThreshold: 1
      failureThreshold: 3
    pdb:
      enabled: true
      maxUnavailable: "1"
    hpa:
      enabled: true
      minReplicas: 3
      maxReplicas: 6
      resourceName: "cpu"
      resourceUtilization: 90
  sidekiq:
    replicas: 2
    resources:
      requests:
        cpu: "500m"
        memory: "1Gi"
      limits:
        cpu: "1"
        memory: "2Gi"
    livenessProbe:
      initialDelaySeconds: 10
      timeoutSeconds: 3
      periodSeconds: 15
      successThreshold: 1
      failureThreshold: 5
    readinessProbe:
      initialDelaySeconds: 10
      timeoutSeconds: 5
      periodSeconds: 30
      successThreshold: 1
      failureThreshold: 5
    pdb:
      enabled: false
      maxUnavailable: "1"
    hpa:
      enabled: false
      minReplicas: 2
      maxReplicas: 4
      resourceName: "cpu"
      resourceUtilization: 90
  sphinx:
    env:
      thinking:
        port: 9306
        bindAddress: "0.0.0.0"
        configFile: "/opt/sphinx/sphinx.conf"
        dbPath: "/opt/sphinx/db"
        pidFile: "/var/run/searchd.pid"
      deltaIndexInterval: 5
      fullReindexInterval: 60
    resources:
      requests:
        cpu: "250m"
        memory: "4Gi"
      limits:
        cpu: "750m"
        memory: "5Gi"
      storage: "30Gi"
    livenessProbe:
      initialDelaySeconds: 60
      timeoutSeconds: 3
      periodSeconds: 15
      successThreshold: 1
      failureThreshold: 5
    readinessProbe:
      initialDelaySeconds: 60
      timeoutSeconds: 5
      periodSeconds: 30
      successThreshold: 1
      failureThreshold: 5
```

## CR Spec

|                  **Field**                   | **Type**  | **Required** |          **Default value**           |                                                            **Description**                                                            |
| :------------------------------------------: | :-------: | :----------: | :----------------------------------: | :-----------------------------------------------------------------------------------------------------------------------------------: |
|                 `image.name`                 | `string`  |      No      |        `quay.io/3scale/porta`        |                                                    Image name (docker repository)                                                     |
|                 `image.tag`                  | `string`  |      No      |              `nightly`               |                                                               Image tag                                                               |
|              `image.pullPolicy`              | `string`  |      No      |                  -                   |                                Determine if the image should be pulled prior to starting the container                                |
|         `grafanaDashboard.label.key`         | `string`  |      No      |           `monitoring-key`           |                                     Label `key` used by grafana-operator for dashboard discovery                                      |
|        `grafanaDashboard.label.value`        | `string`  |      No      |             `middleware`             |                                    Label `value` used by grafana-operator for dashboard discovery                                     |
|               `env.ampRelease`               | `string`  |      No      |               `2.7.1`                |                                                          AMP release number                                                           |
|                `env.forceSSL`                | `string`  |      No      |                `true`                |                                            Enable (true) or disable (false) enforcing SSL                                             |
|               `env.sslCertDir`               | `string`  |      No      |         `/etc/pki/tls/certs`         |                                                         SSL certificates path                                                         |
|     `env.sandboxProxyOpensslVerifyMode`      | `string`  |      No      |            `VERIFY_NONE`             |                        OpenSSL verification mode for sandbox proxy OpenSSL verification mode for sandbox proxy                        |
|         `env.threescaleSuperdomain`          | `string`  |      No      |             `localhost`              |                                                          3scale superdomain                                                           |
|         `env.threescaleProviderPlan`         | `string`  |      No      |             `enterprise`             |                                                         3scale provider plan                                                          |
|                `env.rail.env`                | `string`  |      No      |              `preview`               |                                                           Rails environment                                                           |
|             `env.rail.logLevel`              | `string`  |      No      |                `info`                |                                     Rails log level (debug, info, warn, error, fatal or unknown)                                      |
|            `env.rail.logToStdout`            | `string`  |      No      |                `true`                |                                     Enable (true) or disable (false) writting logs to the stdout                                      |
|        `secret.multitenantIngressTLS`        | `string`  |      No      |                  -                   |                                                  Multitenant Ingress TLS secret name                                                  |
|            `secret.appVaultPath`             | `string`  |     Yes      |                  -                   |                                             Vault Path with system-app secret definition                                              |
|      `secret.backendListenerVaultPath`       | `string`  |     Yes      |                  -                   |                                       Vault Path with system-backend-listener secret definition                                       |
|           `secret.configVaultPath`           | `string`  |     Yes      |                  -                   |                                            Vault Path with system-config secret definition                                            |
|          `secret.databaseVaultPath`          | `string`  |     Yes      |                  -                   |                                           Vault Path with system-database secret definition                                           |
|         `secret.eventsHookVaultPath`         | `string`  |     Yes      |                  -                   |                                         Vault Path with system-events-hook secret definition                                          |
|       `secret.masterApicastVaultPath`        | `string`  |     Yes      |                  -                   |                                        Vault Path with system-master-apicast secret definition                                        |
|         `secret.memcachedVaultPath`          | `string`  |     Yes      |                  -                   |                                          Vault Path with system-memcached secret definition                                           |
|    `secret.multitenantAssetsS3VaultPath`     | `string`  |     Yes      |                  -                   |                                    Vault Path with system-multitenant-assets-s3 secret definition                                     |
|         `secret.recaptchaVaultPath`          | `string`  |     Yes      |                  -                   |                                          Vault Path with system-recaptcha secret definition                                           |
|           `secret.redisVaultPath`            | `string`  |     Yes      |                  -                   |                                            Vault Path with system-redis secret definition                                             |
|            `secret.seedVaultPath`            | `string`  |     Yes      |                  -                   |                                             Vault Path with system-seed secret definition                                             |
|            `secret.smtpVaultPath`            | `string`  |     Yes      |                  -                   |                                             Vault Path with system-smtp secret definition                                             |
|              `app.pdb.enabled`               | `boolean` |      No      |                `true`                |                                       Enable (`true`) or disable (`false`) PodDisruptionBudget                                        |
|           `app.pdb.maxUnavailable`           | `string`  |      No      |                 `1`                  |                                  Maximum number of unavailable pods (number or percentage of pods)**                                  |
|            `app.pdb.minAvailable`            | `string`  |      No      |                  -                   |                      Minimum number of available pods (number or percentage of pods), overrides maxUnavailable**                      |
|              `app.hpa.enabled`               | `boolean` |      No      |                `true`                |                                     Enable (`true`) or disable (`false`) HoritzontalPodAutoscaler                                     |
|            `app.hpa.minReplicas`             |   `int`   |      No      |                 `2`                  |                                                      Minimum number of replicas                                                       |
|            `app.hpa.maxReplicas`             |   `int`   |      No      |                 `4`                  |                                                      Maximum number of replicas                                                       |
|            `app.hpa.resourceName`            | `string`  |      No      |                `cpu`                 |                                               Resource used for autoscale (cpu/memory)                                                |
|        `app.hpa.resourceUtilization`         |   `int`   |      No      |                 `90`                 |                                          Percentage usage of the resource used for autoscale                                          |
|                `app.replicas`                |   `int`   |      No      |                 `2`                  |                                            Number of replicas (ignored if hpa is enabled)                                             |
|             `app.env.logFormat`              | `string`  |      No      |                `json`                |                                                      Log format (`text`/`json`)                                                       |
|          `app.env.listenerWorkers`           |   `int`   |      No      |                 `16`                 |                                              Number of worker processes per listener pod                                              |
|             `app.env.redisAsync`             |  `bool`   |      No      |               `false`                |                                         Enable (`true`) or disable (`false`) redis async mode                                         |
|         `app.resources.requests.cpu`         | `string`  |      No      |                `500m`                |                                                         Override CPU requests                                                         |
|       `app.resources.requests.memory`        | `string`  |      No      |               `550Mi`                |                                                       Override Memory requests                                                        |
|          `app.resources.limits.cpu`          | `string`  |      No      |                 `1`                  |                                                          Override CPU limits                                                          |
|        `app.resources.limits.memory`         | `string`  |      No      |               `700Mi`                |                                                        Override Memory limits                                                         |
|   `app.livenessProbe.initialDelaySeconds`    |   `int`   |      No      |                 `30`                 |                                               Override liveness initial delay (seconds)                                               |
|      `app.livenessProbe.timeoutSeconds`      |   `int`   |      No      |                 `1`                  |                                                  Override liveness timeout (seconds)                                                  |
|      `app.livenessProbe.periodSeconds`       |   `int`   |      No      |                 `10`                 |                                                  Override liveness period (seconds)                                                   |
|     `app.livenessProbe.successThreshold`     |   `int`   |      No      |                 `1`                  |                                                  Override liveness success threshold                                                  |
|     `app.livenessProbe.failureThreshold`     |   `int`   |      No      |                 `3`                  |                                                  Override liveness failure threshold                                                  |
|   `app.readinessProbe.initialDelaySeconds`   |   `int`   |      No      |                 `30`                 |                                              Override readiness initial delay (seconds)                                               |
|     `app.readinessProbe.timeoutSeconds`      |   `int`   |      No      |                 `5`                  |                                                 Override readiness timeout (seconds)                                                  |
|      `app.readinessProbe.periodSeconds`      |   `int`   |      No      |                 `10`                 |                                                  Override readiness period (seconds)                                                  |
|    `app.readinessProbe.successThreshold`     |   `int`   |      No      |                 `1`                  |                                                 Override readiness success threshold                                                  |
|    `app.readinessProbe.failureThreshold`     |   `int`   |      No      |                 `3`                  |                                                 Override readiness failure threshold                                                  |
|            `sidekiq.pdb.enabled`             | `boolean` |      No      |                `true`                |                                       Enable (`true`) or disable (`false`) PodDisruptionBudget                                        |
|         `sidekiq.pdb.maxUnavailable`         | `string`  |      No      |                 `1`                  |                                  Maximum number of unavailable pods (number or percentage of pods)**                                  |
|          `sidekiq.pdb.minAvailable`          | `string`  |      No      |                  -                   |                      Minimum number of available pods (number or percentage of pods), overrides maxUnavailable**                      |
|            `sidekiq.hpa.enabled`             | `boolean` |      No      |                `true`                |                                     Enable (`true`) or disable (`false`) HoritzontalPodAutoscaler                                     |
|          `sidekiq.hpa.minReplicas`           |   `int`   |      No      |                 `2`                  |                                                      Minimum number of replicas                                                       |
|          `sidekiq.hpa.maxReplicas`           |   `int`   |      No      |                 `4`                  |                                                      Maximum number of replicas                                                       |
|          `sidekiq.hpa.resourceName`          | `string`  |      No      |                `cpu`                 |                                               Resource used for autoscale (cpu/memory)                                                |
|      `sidekiq.hpa.resourceUtilization`       |   `int`   |      No      |                 `90`                 |                                          Percentage usage of the resource used for autoscale                                          |
|              `sidekiq.replicas`              |   `int`   |      No      |                 `2`                  |                                            Number of replicas (ignored if hpa is enabled)                                             |
|       `sidekiq.resources.requests.cpu`       | `string`  |      No      |                `500m`                |                                                         Override CPU requests                                                         |
|     `sidekiq.resources.requests.memory`      | `string`  |      No      |                `1Gi`                 |                                                       Override Memory requests                                                        |
|        `sidekiq.resources.limits.cpu`        | `string`  |      No      |                 `1`                  |                                                          Override CPU limits                                                          |
|      `sidekiq.resources.limits.memory`       | `string`  |      No      |                `2Gi`                 |                                                        Override Memory limits                                                         |
| `sidekiq.livenessProbe.initialDelaySeconds`  |   `int`   |      No      |                 `10`                 |                                               Override liveness initial delay (seconds)                                               |
|    `sidekiq.livenessProbe.timeoutSeconds`    |   `int`   |      No      |                 `3`                  |                                                  Override liveness timeout (seconds)                                                  |
|    `sidekiq.livenessProbe.periodSeconds`     |   `int`   |      No      |                 `15`                 |                                                  Override liveness period (seconds)                                                   |
|   `sidekiq.livenessProbe.successThreshold`   |   `int`   |      No      |                 `1`                  |                                                  Override liveness success threshold                                                  |
|   `sidekiq.livenessProbe.failureThreshold`   |   `int`   |      No      |                 `5`                  |                                                  Override liveness failure threshold                                                  |
| `sidekiq.readinessProbe.initialDelaySeconds` |   `int`   |      No      |                 `10`                 |                                              Override readiness initial delay (seconds)                                               |
|   `sidekiq.readinessProbe.timeoutSeconds`    |   `int`   |      No      |                 `5`                  |                                                 Override readiness timeout (seconds)                                                  |
|    `sidekiq.readinessProbe.periodSeconds`    |   `int`   |      No      |                 `30`                 |                                                  Override readiness period (seconds)                                                  |
|  `sidekiq.readinessProbe.successThreshold`   |   `int`   |      No      |                 `1`                  |                                                 Override readiness success threshold                                                  |
|  `sidekiq.readinessProbe.failureThreshold`   |   `int`   |      No      |                 `5`                  |                                                 Override readiness failure threshold                                                  |
|          `sphinx.env.thinking.port`          |   `int`   |      No      |                `9306`                |                                              The TCP port Sphinx will run its daemon on                                               |
|      `sphinx.env.thinking.bindAddress`       | `string`  |      No      |              `0.0.0.0`               |                                     Allows setting the TCP host for Sphinx to a different address                                     |
|        `sphinx.env.thinking.endpoint`        | `string`  |      No      |           `system-sphinx`            |                                                        The Sphinx DNS endpoint                                                        |
|       `sphinx.env.thinking.configFile`       | `string`  |      No      | `/opt/system/db/sphinx/preview.conf` |                                                    Sphinx configuration file path                                                     |
|         `sphinx.env.thinking.dbPath`         | `string`  |      No      |       `/opt/system/db/sphinx`        |                                                         Sphinx database path                                                          |
|        `sphinx.env.thinking.pidFile`         | `string`  |      No      |  `/opt/system/tmp/pids/searchd.pid`  |                                                         Sphinx PID file path                                                          |
|       `sphinx.env.deltaIndexInterval`        |   `int`   |      No      |                 `5`                  | Interval used for adding chunks of brand new documents to the primary index at certain intervals without having to do a full re-index |
|       `sphinx.env.fullReindexInterval`       |   `int`   |      No      |                 `60`                 |                                                  Interval used to do a full re-index                                                  |
|       `sphinx.resources.requests.cpu`        | `string`  |      No      |                `250m`                |                                                         Override CPU requests                                                         |
|      `sphinx.resources.requests.memory`      | `string`  |      No      |                `4Gi`                 |                                                       Override Memory requests                                                        |
|        `sphinx.resources.limits.cpu`         | `string`  |      No      |                `750m`                |                                                          Override CPU limits                                                          |
|       `sphinx.resources.limits.memory`       | `string`  |      No      |                `5Gi`                 |                                                        Override Memory limits                                                         |
|          `sphinx.resources.storage`          | `string`  |      No      |                `30Gi`                |                                                       Override Memory requests                                                        |
|  `sphinx.livenessProbe.initialDelaySeconds`  |   `int`   |      No      |                 `10`                 |                                               Override liveness initial delay (seconds)                                               |
|    `sphinx.livenessProbe.timeoutSeconds`     |   `int`   |      No      |                 `3`                  |                                                  Override liveness timeout (seconds)                                                  |
|     `sphinx.livenessProbe.periodSeconds`     |   `int`   |      No      |                 `15`                 |                                                  Override liveness period (seconds)                                                   |
|   `sphinx.livenessProbe.successThreshold`    |   `int`   |      No      |                 `1`                  |                                                  Override liveness success threshold                                                  |
|   `sphinx.livenessProbe.failureThreshold`    |   `int`   |      No      |                 `5`                  |                                                  Override liveness failure threshold                                                  |
| `sphinx.readinessProbe.initialDelaySeconds`  |   `int`   |      No      |                 `10`                 |                                              Override readiness initial delay (seconds)                                               |
|    `sphinx.readinessProbe.timeoutSeconds`    |   `int`   |      No      |                 `5`                  |                                                 Override readiness timeout (seconds)                                                  |
|    `sphinx.readinessProbe.periodSeconds`     |   `int`   |      No      |                 `30`                 |                                                  Override readiness period (seconds)                                                  |
|   `sphinx.readinessProbe.successThreshold`   |   `int`   |      No      |                 `1`                  |                                                 Override readiness success threshold                                                  |
|   `sphinx.readinessProbe.failureThreshold`   |   `int`   |      No      |                 `5`                  |                                                 Override readiness failure threshold                                                  |


** If you are already using `pdb.maxUnavailable` and want to use `pdb.minAvailable` (or the other way around), due to ansible operator limitation of doing patch operation (if objects already exist), operator will receive an error when managing PDB object because although the spec of the PDB resource it creates is correct, operator will try to patch an existing object which already has the other variable, and these two variables `pdb.maxUnavailable`/`pdb.minAvailable` are mutually exclusive and cannot coexists on the same PDB. To solve that situation:
  - Configure `pdb.enabled=false` (so operator will delete associated PDB, and then re-enable it with `pdb.enabled=true` setting desired PDB field `pdb.minAvailable` or `pdb.maxUnavailable`. so operator will create it from scratch on next reconcile
  - Or, delete manually associated PDB object, and operator will create it from scratch on next reconcile
