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
  pdb:
    enabled: true
    maxUnavailable: "1"
  hpa:
    enabled: true
    minReplicas: 2
    maxReplicas: 4
    resourceName: cpu
    resourceUtilization: 90
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
| `pdb.enabled` | `boolean` | No | `true` | Enable (`true`) or disable (`false`) PodDisruptionBudget |
| `pdb.maxUnavailable` | `string` | No | `1` | Maximum number of unavailable pods (number or percentage of pods) ** |
| `pdb.minAvailable` | `string` | No | - | Minimum number of available pods (number or percentage of pods), overrides maxUnavailable ** |
| `hpa.enabled` | `boolean` | No | `true` | Enable (`true`) or disable (`false`) HoritzontalPodAutoscaler |
| `hpa.minReplicas` | `int` | No | `2` | Minimum number of replicas |
| `hpa.maxReplicas` | `int` | No | `4` | Maximum number of replicas |
| `hpa.resourceName` | `string` | No | `cpu` | Resource used for autoscale (cpu/memory) |
| `hpa.resourceUtilization` | `int` | No | `90` | Percentage usage of the resource used for autoscale |
| `replicas` | `int` | No | `2` | Number of replicas (ignored if hpa is enabled) |
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

** If you are already using `pdb.maxUnavailable` and want to use `pdb.minAvailable` (or the other way around), due to ansible operator limitation of doing patch operation (if objects already exist), operator will receive an error when managing PDB object because although the spec of the PDB resource it creates is correct, operator will try to patch an existing object which already has the other variable, and these two variables `pdb.maxUnavailable`/`pdb.minAvailable` are mutually exclusive and cannot coexists on the same PDB. To solve that situation:
  - Configure `pdb.enabled=false` (so operator will delete associated PDB, and then re-enable it with `pdb.enabled=true` setting desired PDB field `pdb.minAvailable` or `pdb.maxUnavailable`. so operator will create it from scratch on next reconcile
  - Or, delete manually associated PDB object, and operator will create it from scratch on next reconcile