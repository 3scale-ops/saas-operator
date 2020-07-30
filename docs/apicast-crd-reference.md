# Apicast Custom Resource Reference

## Simple CR Example

```yaml
apiVersion: saas.3scale.net/v1alpha1
kind: Apicast
metadata:
  name: simple-example
spec:
  staging:
    externalDnsHostname: "apicast-staging.test.stg-saas.3sca.net"
    image:
      tag: apicast-v3.8.0-r3
    replicas: 1
    env:
      apicastConfigurationCache: "30"
      threescalePortalEndpoint: "http://mapping-service/config"
    marin3r:
      enabled: true
      annotations:
        marin3r.3scale.net/node-id: apicast-staging
        marin3r.3scale.net/ports: gateway-http:38080,gateway-https:38443,envoy-metrics:9901
  production:
    externalDnsHostname: "apicast-production.test.stg-saas.3sca.net"
    image:
      tag: apicast-v3.8.0-r3
    replicas: 1
    env:
      apicastConfigurationCache: "300"
      threescalePortalEndpoint: "http://mapping-service/config"
    marin3r:
      enabled: true
      annotations:
        marin3r.3scale.net/node-id: apicast-staging
        marin3r.3scale.net/ports: gateway-http:38080,gateway-https:38443,envoy-metrics:9901
```

## Full CR Example

Most of the fields are not (can use default values), this is just an example of everything that can be overriden under your own risk:

```yaml
apiVersion: saas.3scale.net/v1alpha1
kind: Apicast
metadata:
  name: full-example
spec:
  staging:
    externalDnsHostname: "apicast-staging.test.stg-saas.3sca.net"
    image:
      name: quay.io/3scale/apicast-cloud-hosted
      tag: apicast-v3.8.0-r3
      pullSecretName: quay-pull-secret
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
      apicastConfigurationCache: "30"
      threescalePortalEndpoint: "http://mapping-service/config"
      apicastLogLevel: debug
      apicsatOIDCLogLevel: debug
    marin3r:
      enabled: true
      annotations:
        marin3r.3scale.net/node-id: apicast-staging
        marin3r.3scale.net/ports: gateway-http:38080,gateway-https:38443,envoy-metrics:9901
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
    resources:
      requests:
        cpu: "500m"
        memory: "64Mi"
      limits:
        cpu: "1"
        memory: "128Mi"

  production:
    externalDnsHostname: "apicast-production.test.stg-saas.3sca.net"
    image:
      name: quay.io/3scale/apicast-cloud-hosted
      tag: apicast-v3.8.0-r3
      pullSecretName: quay-pull-secret
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
      apicastConfigurationCache: "300"
      threescalePortalEndpoint: "http://mapping-service/config"
      apicastLogLevel: debug
      apicsatOIDCLogLevel: debug
    marin3r:
      enabled: true
      annotations:
        marin3r.3scale.net/node-id: apicast-staging
        marin3r.3scale.net/ports: gateway-http:38080,gateway-https:38443,envoy-metrics:9901
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
    resources:
      requests:
        cpu: "500m"
        memory: "64Mi"
      limits:
        cpu: "1"
        memory: "128Mi"

  loadBalancer:
    proxyProtocol: "*"
    crossZoneLoadBalancingEnabled: true
    connectionDrainingEnabled: true
    connectionDrainingTimeout: 60
    connectionHealthcheckHealthyThreshold: 2
    connectionHealthcheckUnhealthyThreshold: 2
    connectionHealthcheckInterval: 5
    connectionHealthcheckTimeout: 3

  grafanaDashboard:
    label:
      key: discovery
      value: enabled
```

## CR Spec

|                             **Field**                             | **Type**  | **Required** |           **Default value**           |                                      **Description**                                      |
| :---------------------------------------------------------------: | :-------: | :----------: | :-----------------------------------: | :---------------------------------------------------------------------------------------: |
|                       `staging.image.name`                        | `string`  |      No      | `quay.io/3scale/apicast-cloud-hosted` |                              Image name (docker repository)                               |
|                        `staging.image.tag`                        | `string`  |     Yes      |                   -                   |                                         Image tag                                         |
|                  `staging.image.pullSecretName`                   | `string`  |      No      |                   -                   |                          Quay pull secret for private repository                          |
|                       `staging.pdb.enabled`                       | `boolean` |      No      |                `true`                 |                 Enable (`true`) or disable (`false`) PodDisruptionBudget                  |
|                   `staging.pdb.maxUnavailable`                    | `string`  |      No      |                  `1`                  |             Maximum number of unavailable pods (number or percentage of pods)             |
|                    `staging.pdb.minAvailable`                     | `string`  |      No      |                   -                   | Minimum number of available pods (number or percentage of pods), overrides maxUnavailable |
|                       `staging.hpa.enabled`                       | `boolean` |      No      |                `true`                 |               Enable (`true`) or disable (`false`) HoritzontalPodAutoscaler               |
|                     `staging.hpa.minReplicas`                     |   `int`   |      No      |                  `2`                  |                                Minimum number of replicas                                 |
|                     `staging.hpa.maxReplicas`                     |   `int`   |      No      |                  `4`                  |                                Maximum number of replicas                                 |
|                    `staging.hpa.resourceName`                     | `string`  |      No      |                 `cpu`                 |                         Resource used for autoscale (cpu/memory)                          |
|                 `staging.hpa.resourceUtilization`                 |   `int`   |      No      |                 `90`                  |                    Percentage usage of the resource used for autoscale                    |
|                        `staging.replicas`                         |   `int`   |      No      |                  `2`                  |                      Number of replicas (ignored if hpa is enabled)                       |
|              `staging.env.apicastConfigurationCache`              | `string`  |     Yes      |                   -                   |                             Apicast configurations cache TTL                              |
|              `staging.env.threescalePortalEndpoint`               | `string`  |     Yes      |                   -                   |                        Endpoint to request proxy configurations to                        |
|                   `staging.env.apicastLogLevel`                   | `string`  |      No      |                `warn`                 |                                    Openresty log level                                    |
|                 `staging.env.apicastOIDCLogLevel`                 | `string`  |      No      |               `notice`                |                           OpenID Connect integration log level                            |
|                     `staging.marin3r.enabled`                     | `boolean` |     Yes      |                   -                   |                       Enable (`true`) or disable (`false`) marin3r                        |
|                 `staging.marin3r.annotations.{}`                  |   `map`   |      No      |                   -                   |                                Map of marin3r annotations                                 |
|                 `staging.resources.requests.cpu`                  | `string`  |      No      |                `500m`                 |                                   Override CPU requests                                   |
|                `staging.resources.requests.memory`                | `string`  |      No      |                `64Mi`                 |                                 Override Memory requests                                  |
|                  `staging.resources.limits.cpu`                   | `string`  |      No      |                  `1`                  |                                    Override CPU limits                                    |
|                 `staging.resources.limits.memory`                 | `string`  |      No      |                `128Mi`                |                                  Override Memory limits                                   |
|            `staging.livenessProbe.initialDelaySeconds`            |   `int`   |      No      |                  `5`                  |                         Override liveness initial delay (seconds)                         |
|              `staging.livenessProbe.timeoutSeconds`               |   `int`   |      No      |                  `5`                  |                            Override liveness timeout (seconds)                            |
|               `staging.livenessProbe.periodSeconds`               |   `int`   |      No      |                 `10`                  |                            Override liveness period (seconds)                             |
|             `staging.livenessProbe.successThreshold`              |   `int`   |      No      |                  `1`                  |                            Override liveness success threshold                            |
|             `staging.livenessProbe.failureThreshold`              |   `int`   |      No      |                  `3`                  |                            Override liveness failure threshold                            |
|           `staging.readinessProbe.initialDelaySeconds`            |   `int`   |      No      |                  `5`                  |                        Override readiness initial delay (seconds)                         |
|              `staging.readinessProbe.timeoutSeconds`              |   `int`   |      No      |                  `5`                  |                           Override readiness timeout (seconds)                            |
|              `staging.readinessProbe.periodSeconds`               |   `int`   |      No      |                 `30`                  |                            Override readiness period (seconds)                            |
|             `staging.readinessProbe.successThreshold`             |   `int`   |      No      |                  `1`                  |                           Override readiness success threshold                            |
|             `staging.readinessProbe.failureThreshold`             |   `int`   |      No      |                  `3`                  |                           Override readiness failure threshold                            |
|                   `staging.externalDnsHostname`                   | `string`  |     Yes      |                   -                   |                  DNS hostnames to manage on AWS Route53 by external-dns                   |
|               `staging.loadBalancer.proxyProtocol`                | `string`  |      No      |                  `*`                  |         Proxy protocol enabled (k8s aws provider only accepts `*` by the moment)          |
|       `staging.loadBalancer.crossZoneLoadBalancingEnabled`        |  `bool`   |      No      |                `true`                 |              Enable (`true`) or disable (`false`) cross zone load balancing               |
|         `staging.loadBalancer.connectionDrainingEnabled`          |  `bool`   |      No      |                `true`                 |                 Enable (`true`) or disable (`false`) connection draining                  |
|         `staging.loadBalancer.connectionDrainingTimeout`          |   `int`   |      No      |                 `60`                  |                           Connection draining timeout (seconds)                           |
|   `staging.loadBalancer.connectionHealthcheckHealthyThreshold`    |   `int`   |      No      |                  `2`                  |                         Connection healthcheck healthy threshold                          |
|  `staging.loadBalancer.connectionHealthcheckUnhealthyThreshold`   |   `int`   |      No      |                  `2`                  |                        Connection healthcheck unhealthy threshold                         |
|       `staging.loadBalancer.connectionHealthcheckInterval`        |   `int`   |      No      |                  `5`                  |                         Connection healthcheck interval (seconds)                         |
|        `staging.loadBalancer.connectionHealthcheckTimeout`        |   `int`   |      No      |                  `3`                  |                         Connection healthcheck timeout (seconds)                          |
|                      `production.image.name`                      | `string`  |      No      | `quay.io/3scale/apicast-cloud-hosted` |                              Image name (docker repository)                               |
|                      `production.image.tag`                       | `string`  |     Yes      |                   -                   |                                         Image tag                                         |
|                 `production.image.pullSecretName`                 | `string`  |      No      |                   -                   |                          Quay pull secret for private repository                          |
|                     `production.pdb.enabled`                      | `boolean` |      No      |                `true`                 |                 Enable (`true`) or disable (`false`) PodDisruptionBudget                  |
|                  `production.pdb.maxUnavailable`                  | `string`  |      No      |                  `1`                  |             Maximum number of unavailable pods (number or percentage of pods)             |
|                   `production.pdb.minAvailable`                   | `string`  |      No      |                   -                   | Minimum number of available pods (number or percentage of pods), overrides maxUnavailable |
|                     `production.hpa.enabled`                      | `boolean` |      No      |                `true`                 |               Enable (`true`) or disable (`false`) HoritzontalPodAutoscaler               |
|                   `production.hpa.minReplicas`                    |   `int`   |      No      |                  `2`                  |                                Minimum number of replicas                                 |
|                   `production.hpa.maxReplicas`                    |   `int`   |      No      |                  `4`                  |                                Maximum number of replicas                                 |
|                   `production.hpa.resourceName`                   | `string`  |      No      |                 `cpu`                 |                         Resource used for autoscale (cpu/memory)                          |
|               `production.hpa.resourceUtilization`                |   `int`   |      No      |                 `90`                  |                    Percentage usage of the resource used for autoscale                    |
|                       `production.replicas`                       |   `int`   |      No      |                  `2`                  |                      Number of replicas (ignored if hpa is enabled)                       |
|            `production.env.apicastConfigurationCache`             | `string`  |     Yes      |                   -                   |                             Apicast configurations cache TTL                              |
|             `production.env.threescalePortalEndpoint`             | `string`  |     Yes      |                   -                   |                        Endpoint to request proxy configurations to                        |
|                 `production.env.apicastLogLevel`                  | `string`  |      No      |                `warn`                 |                                    Openresty log level                                    |
|               `production.env.apicastOIDCLogLevel`                | `string`  |      No      |               `notice`                |                           OpenID Connect integration log level                            |
|                   `production.marin3r.enabled`                    | `boolean` |     Yes      |                   -                   |                       Enable (`true`) or disable (`false`) marin3r                        |
|                `production.marin3r.annotations.{}`                |   `map`   |      No      |                   -                   |                                Map of marin3r annotations                                 |
|                `production.resources.requests.cpu`                | `string`  |      No      |                `500m`                 |                                   Override CPU requests                                   |
|              `production.resources.requests.memory`               | `string`  |      No      |                `64Mi`                 |                                 Override Memory requests                                  |
|                 `production.resources.limits.cpu`                 | `string`  |      No      |                  `1`                  |                                    Override CPU limits                                    |
|               `production.resources.limits.memory`                | `string`  |      No      |                `128Mi`                |                                  Override Memory limits                                   |
|          `production.livenessProbe.initialDelaySeconds`           |   `int`   |      No      |                  `5`                  |                         Override liveness initial delay (seconds)                         |
|             `production.livenessProbe.timeoutSeconds`             |   `int`   |      No      |                  `5`                  |                            Override liveness timeout (seconds)                            |
|             `production.livenessProbe.periodSeconds`              |   `int`   |      No      |                 `10`                  |                            Override liveness period (seconds)                             |
|            `production.livenessProbe.successThreshold`            |   `int`   |      No      |                  `1`                  |                            Override liveness success threshold                            |
|            `production.livenessProbe.failureThreshold`            |   `int`   |      No      |                  `3`                  |                            Override liveness failure threshold                            |
|          `production.readinessProbe.initialDelaySeconds`          |   `int`   |      No      |                  `5`                  |                        Override readiness initial delay (seconds)                         |
|            `production.readinessProbe.timeoutSeconds`             |   `int`   |      No      |                  `5`                  |                           Override readiness timeout (seconds)                            |
|             `production.readinessProbe.periodSeconds`             |   `int`   |      No      |                 `30`                  |                            Override readiness period (seconds)                            |
|           `production.readinessProbe.successThreshold`            |   `int`   |      No      |                  `1`                  |                           Override readiness success threshold                            |
|           `production.readinessProbe.failureThreshold`            |   `int`   |      No      |                  `3`                  |                           Override readiness failure threshold                            |
|                 `production.externalDnsHostname`                  | `string`  |     Yes      |                   -                   |                  DNS hostnames to manage on AWS Route53 by external-dns                   |
|              `production.loadBalancer.proxyProtocol`              | `string`  |      No      |                  `*`                  |         Proxy protocol enabled (k8s aws provider only accepts `*` by the moment)          |
|      `production.loadBalancer.crossZoneLoadBalancingEnabled`      |  `bool`   |      No      |                `true`                 |              Enable (`true`) or disable (`false`) cross zone load balancing               |
|        `production.loadBalancer.connectionDrainingEnabled`        |  `bool`   |      No      |                `true`                 |                 Enable (`true`) or disable (`false`) connection draining                  |
|        `production.loadBalancer.connectionDrainingTimeout`        |   `int`   |      No      |                 `60`                  |                           Connection draining timeout (seconds)                           |
|  `production.loadBalancer.connectionHealthcheckHealthyThreshold`  |   `int`   |      No      |                  `2`                  |                         Connection healthcheck healthy threshold                          |
| `production.loadBalancer.connectionHealthcheckUnhealthyThreshold` |   `int`   |      No      |                  `2`                  |                        Connection healthcheck unhealthy threshold                         |
|      `production.loadBalancer.connectionHealthcheckInterval`      |   `int`   |      No      |                  `5`                  |                         Connection healthcheck interval (seconds)                         |
|      `production.loadBalancer.connectionHealthcheckTimeout`       |   `int`   |      No      |                  `3`                  |                         Connection healthcheck timeout (seconds)                          |
|                   `grafanaDashboard.label.key`                    | `string`  |      No      |           `monitoring-key`            |               Label `key` used by grafana-operator for dashboard discovery                |
|                  `grafanaDashboard.label.value`                   | `string`  |      No      |             `middleware`              |              Label `value` used by grafana-operator for dashboard discovery               |
