# Backend Custom Resource Reference

## Simple CR Example

```yaml
apiVersion: saas.3scale.net/v1alpha1
kind: Backend
metadata:
  name: simple-example
spec:
  image:
    tag: v2.101.1
  config:
    rackEnv: example
  secret:
    redisVaultPath: secret/data/openshift/cluster-example/3scale/backend-redis
    systemEventsHookVaultPath: secret/data/openshift/cluster-example/3scale/backend-system-events-hook
    internalApiVaultPath: secret/data/openshift/cluster-example/3scale/backend-internal-api
  listener:
    externalDnsHostname: backend.example.3scale.net
    marin3r:
      enabled: true
      annotations:
        marin3r.3scale.net/node-id: backend-listener
        marin3r.3scale.net/ports: backend-listener-http:38080,backend-listener-https:38443,envoy-metrics:9901
    replicas: 1
  worker:
    replicas: 1
  cron:
    replicas: 1
```

## Full CR Example

Most of the fields do not need to be specified (can use default values), this is just an example of everything that can be overriden under your own risk:

```yaml
apiVersion: saas.3scale.net/v1alpha1
kind: Backend
metadata:
  name: full-example
spec:
  image:
    name: quay.io/3scale/apisonator
    tag: v2.101.1
    pullSecretName: quay-pull-secret
  grafanaDashboard:
    label:
      key: discovery
      value: enabled
  config:
    rackEnv: example
    masterServiceId: 6
    oauthMaxTokenSize: 7888
    legacyReferrerFilters: true
  secret:
    redisVaultPath: secret/data/openshift/cluster-example/3scale/backend-redis
    systemEventsHookVaultPath: secret/data/openshift/cluster-example/3scale/backend-system-events-hook
    internalApiVaultPath: secret/data/openshift/cluster-example/3scale/backend-internal-api
    errorMonitoringVaultPath: secret/data/openshift/cluster-example/3scale/backend-error-monitoring
  errorMonitoringEnabled: false
  listener:
    externalDnsHostname: backend.example.3scale.net
    marin3r:
      enabled: true
      annotations:
        marin3r.3scale.net/node-id: backend-listener
        marin3r.3scale.net/ports: backend-listener-http:38080,backend-listener-https:38443,envoy-metrics:9901
    loadBalancer:
      proxyProtocol: true
      crossZoneLoadBalancingEnabled: true
      eipAllocations: "eipalloc-080ecfaf74a799b24,eipalloc-098963e814413a5d1,eipalloc-02bd497572f4321a0"
    replicas: 2
    env:
      logFormat: json
      redisAsync: false
      listenerWorkers: 16
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
    resources:
      requests:
        cpu: "500m"
        memory: "550Mi"
      limits:
        cpu: "1"
        memory: "700Mi"
  worker:
    replicas: 2
    env:
      logFormat: json
      redisAsync: false
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
    resources:
      requests:
        cpu: "150m"
        memory: "50Mi"
      limits:
        cpu: "1"
        memory: "300Mi"
  cron:
    replicas: 1
    resources:
      requests:
        cpu: "50m"
        memory: "40Mi"
      limits:
        cpu: "150m"
        memory: "80Mi"
```

## CR Spec

| **Field** | **Type** | **Required** | **Default value** | **Description** |
|:---:|:---:|:---:|:---:|:---:|
| `image.name` | `string` | No | `quay.io/3scale/apisonator` | Image name (docker repository) |
| `image.tag` | `string` | No | `nightly` | Image tag |
| `image.pullSecretName` | `string` | No | - | Quay pull secret for private repository |
| `grafanaDashboard.label.key` | `string` | No | `monitoring-key` | Label `key` used by grafana-operator for dashboard discovery |
| `grafanaDashboard.label.value` | `string` | No | `middleware` | Label `value` used by grafana-operator for dashboard discovery |
| `config.rackEnv` | `string` | No | `dev` | Rack environment (used for example for error-monitoring ID) |
| `config.masterServiceId` | `int` | No | `6` | Master service account ID in Porta |
| `config.oauthMaxTokenSize` | `int` | No | `7888` | Oauth Max token size (bytes) |
| `config.legacyReferrerFilters` | `bool` | No | `true` | Enable (`true`) or disable (`false`) Legacy Referrer Filters |
| `secret.redisVaultPath` | `string` | Yes | - | Vault Path with backend-redis secret definition |
| `secret.systemEventsHookVaultPath` | `string` | Yes | - | Vault Path with backend-system-events-hook secret definition |
| `secret.internalApiVaultPath` | `string` | Yes | - | Vault Path with backend-internal-api secret definition |
| `secret.errorMonitoringVaultPath` | `string` | No | - | Vault Path with backend-error-monitoring secret definition |
| `errorMonitoringEnabled` | `bool` | No | `false` | Mount (`true`) or not (`false`) backend-error-monitoring Secret on deployments |
| `listener.externalDnsHostname` | `string` | Yes | - | DNS hostnames to manage on AWS Route53 by external-dns |
| `listener.marin3r.enabled` | `boolean` | Yes | - | Enable (`true`) or disable (`false`) marin3r |
| `listener.marin3r.anotations.{}` | `map` | No | - | Map of marin3r annotations |
| `listener.loadBalancer.proxyProtocol` | `boolean` | No | `true` | Enable (`true`) or disable (`false`) proxy protocol with aws-nlb-helper-operator |
| `listener.loadBalancer.crossZoneLoadBalancingEnabled` | `bool` | No | `true` | Enable (`true`) or disable (`false`) cross zone load balancing |
| `listener.loadBalancer.eipAllocations` | `string` | No | - | Optional Elastic IPs allocations |
| `listener.replicas` | `int` | No | `1` | Number of replicas |
| `listener.env.logFormat` | `string` | No | `json` | Log format (`text`/`json`) |
| `listener.env.listenerWorkers` | `int` | No | `16` | Number of worker processes per listener pod |
| `listener.env.redisAsync` | `bool` | No | `false` | Enable (`true`) or disable (`false`) redis async mode |
| `listener.resources.requests.cpu` | `string` | No | `500m` | Override CPU requests |
| `listener.resources.requests.memory` | `string` | No | `550Mi` | Override Memory requests |
| `listener.resources.limits.cpu` | `string` | No | `1` | Override CPU limits |
| `listener.resources.limits.memory` | `string` | No | `700Mi` | Override Memory limits |
| `listener.livenessProbe.initialDelaySeconds` | `int` | No | `30` | Override liveness initial delay (seconds) |
| `listener.livenessProbe.timeoutSeconds` | `int` | No | `1` | Override liveness timeout (seconds) |
| `listener.livenessProbe.periodSeconds` | `int` | No | `10` | Override liveness period (seconds) |
| `listener.livenessProbe.successThreshold` | `int` | No | `1` | Override liveness success threshold |
| `listener.livenessProbe.failureThreshold` | `int` | No | `3` | Override liveness failure threshold |
| `listener.readinessProbe.initialDelaySeconds` | `int` | No | `30` | Override readiness initial delay (seconds) |
| `listener.readinessProbe.timeoutSeconds` | `int` | No | `5` | Override readiness timeout (seconds) |
| `listener.readinessProbe.periodSeconds` | `int` | No | `10` | Override readiness period (seconds) |
| `listener.readinessProbe.successThreshold` | `int` | No | `1` | Override readiness success threshold |
| `listener.readinessProbe.failureThreshold` | `int` | No | `3` | Override readiness failure threshold |
| `worker.replicas` | `int` | No | `1` | Number of replicas |
| `worker.env.logFormat` | `string` | No | `json` | Log format (`text`/`json`) |
| `worker.env.redisAsync` | `bool` | No | `false` | Enable (`true`) or disable (`false`) redis async mode |
| `worker.resources.requests.cpu` | `string` | No | `150m` | Override CPU requests |
| `worker.resources.requests.memory` | `string` | No | `50Mi` | Override Memory requests |
| `worker.resources.limits.cpu` | `string` | No | `1` | Override CPU limits |
| `worker.resources.limits.memory` | `string` | No | `300Mi` | Override Memory limits |
| `worker.livenessProbe.initialDelaySeconds` | `int` | No | `10` | Override liveness initial delay (seconds) |
| `worker.livenessProbe.timeoutSeconds` | `int` | No | `3` | Override liveness timeout (seconds) |
| `worker.livenessProbe.periodSeconds` | `int` | No | `15` | Override liveness period (seconds) |
| `worker.livenessProbe.successThreshold` | `int` | No | `1` | Override liveness success threshold |
| `worker.livenessProbe.failureThreshold` | `int` | No | `5` | Override liveness failure threshold |
| `worker.readinessProbe.initialDelaySeconds` | `int` | No | `10` | Override readiness initial delay (seconds) |
| `worker.readinessProbe.timeoutSeconds` | `int` | No | `5` | Override readiness timeout (seconds) |
| `worker.readinessProbe.periodSeconds` | `int` | No | `30` | Override readiness period (seconds) |
| `worker.readinessProbe.successThreshold` | `int` | No | `1` | Override readiness success threshold |
| `worker.readinessProbe.failureThreshold` | `int` | No | `5` | Override readiness failure threshold |
| `cron.replicas` | `int` | No | `1` | Number of replicas |
| `cron.resources.requests.cpu` | `string` | No | `50m` | Override CPU requests |
| `cron.resources.requests.memory` | `string` | No | `40Mi` | Override Memory requests |
| `cron.resources.limits.cpu` | `string` | No | `150m` | Override CPU limits |
| `cron.resources.limits.memory` | `string` | No | `80Mi` | Override Memory limits |