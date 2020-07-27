# EchoAPI Custom Resource Reference

## Simple CR Example

```yaml
apiVersion: saas.3scale.net/v1alpha1
kind: EchoAPI
metadata:
  name: simple-example
spec:
  image:
    tag: v1.0.1
  replicas: 1
  externalDnsHostname: echo-api.example.3scale.net
  marin3r:
    enabled: true
    annotations:
      marin3r.3scale.net/node-id: echo-api
      marin3r.3scale.net/ports: echo-api-http:38080,echo-api-https:38443,envoy-metrics:9901
```

## Full CR Example

Most of the fields do not need to be specified (can use default values), this is just an example of everything that can be overriden under your own risk:

```yaml
apiVersion: saas.3scale.net/v1alpha1
kind: EchoAPI
metadata:
  name: full-example
spec:
  image:
    name: quay.io/3scale/echoapi
    tag: v1.0.1
    pullSecretName: quay-pull-secret
  replicas: 2
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
      cpu: "75m"
      memory: "64Mi"
    limits:
      cpu: "150m"
      memory: "128Mi"
  externalDnsHostname: echo-api.example.3scale.net
  marin3r:
    enabled: true
    annotations:
      marin3r.3scale.net/node-id: echo-api
      marin3r.3scale.net/ports: echo-api-http:38080,echo-api-https:38443,envoy-metrics:9901
  loadBalancer:
    proxyProtocol: true
    crossZoneLoadBalancingEnabled: true
```

## CR Spec

| **Field** | **Type** | **Required** | **Default value** | **Description** |
|:---:|:---:|:---:|:---:|:---:|
| `image.name` | `string` | No | `quay.io/3scale/echoapi` | Image name (docker repository) |
| `image.tag` | `string` | No | `latest` | Image tag |
| `image.pullSecretName` | `string` | No | - | Quay pull secret for private repository |
| `replicas` | `int` | No | `2` | Number of replicas |
| `resources.requests.cpu` | `string` | No | `50m` | Override CPU requests |
| `resources.requests.memory` | `string` | No | `40Mi` | Override Memory requests |
| `resources.limits.cpu` | `string` | No | `150m` | Override CPU limits |
| `resources.limits.memory` | `string` | No | `80Mi` | Override Memory limits |
| `livenessProbe.initialDelaySeconds` | `int` | No | `5` | Override liveness initial delay (seconds) |
| `livenessProbe.timeoutSeconds` | `int` | No | `5` | Override liveness timeout (seconds) |
| `livenessProbe.periodSeconds` | `int` | No | `10` | Override liveness period (seconds) |
| `livenessProbe.successThreshold` | `int` | No | `1` | Override liveness success threshold |
| `livenessProbe.failureThreshold` | `int` | No | `3` | Override liveness failure threshold |
| `readinessProbe.initialDelaySeconds` | `int` | No | `5` | Override readiness initial delay (seconds) |
| `readinessProbe.timeoutSeconds` | `int` | No | `5` | Override readiness timeout (seconds) |
| `readinessProbe.periodSeconds` | `int` | No | `30` | Override readiness period (seconds) |
| `readinessProbe.successThreshold` | `int` | No | `1` | Override readiness success threshold |
| `readinessProbe.failureThreshold` | `int` | No | `3` | Override readiness failure threshold |
| `externalDnsHostname` | `string` | Yes | - | DNS hostnames to manage on AWS Route53 by external-dns |
| `marin3r.enabled` | `boolean` | Yes | - | Enable (`true`) or disable (`false`) marin3r |
| `marin3r.annotations.{}` | `map` | No | - | Map of marin3r annotations |
| `loadBalancer.proxyProtocol` | `boolean` | No | `true` | Enable (`true`) or disable (`false`) proxy protocol with aws-nlb-helper-operator |
| `loadBalancer.crossZoneLoadBalancingEnabled` | `bool` | No | `true` | Enable (`true`) or disable (`false`) cross zone load balancing |
