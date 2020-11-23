# CORSProxy Custom Resource Reference

## Simple CR Example

```yaml
apiVersion: saas.3scale.net/v1alpha1
kind: CORSProxy
metadata:
  name: simple-example
spec:
  image:
    tag: v1.1.0
    pullSecretName: quay-pull-secret
  secret:
    systemDatabaseVaultPath: secret/data/openshift/cluster-example/3scale/cors-proxy-system-database
  replicas: 1
  ingress:
    host: "*.example.3sca.net"
```

## Full CR Example

Most of the fields do not need to be specified (can use default values), this is just an example of everything that can be overriden under your own risk:

```yaml
apiVersion: saas.3scale.net/v1alpha1
kind: CORSProxy
metadata:
  name: full-example
spec:
  image:
    name: quay.io/3scale/cors-proxy
    tag: v1.1.0
    pullSecretName: quay-pull-secret
  secret:
    systemDatabaseVaultPath: secret/data/openshift/cluster-example/3scale/cors-proxy-system-database
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
    initialDelaySeconds: 3
    timeoutSeconds: 1
    periodSeconds: 10
    successThreshold: 1
    failureThreshold: 3
  readinessProbe:
    initialDelaySeconds: 3
    timeoutSeconds: 1
    periodSeconds: 10
    successThreshold: 1
    failureThreshold: 3
  resources:
    requests:
      cpu: "250m"
      memory: "250Mi"
    limits:
      cpu: "1"
      memory: "512Mi"
  ingress:
    path: /api_docs/proxy
    host: "*.example.3sca.net"
  grafanaDashboard:
    label:
      key: discovery
      value: enabled
```

## CR Spec

| **Field** | **Type** | **Required** | **Default value** | **Description** |
|:---:|:---:|:---:|:---:|:---:|
| `image.name` | `string` | No | `quay.io/3scale/cors-proxy` | Image name (docker repository) |
| `image.tag` | `string` | No | `v1.1.0` | Image tag |
| `image.pullSecretName` | `string` | No | - | Quay pull secret for private repository |
| `secret.systemDatabaseVaultPath` | `string` | Yes | - | Vault Path with cors-proxy-system-database secret definition |
| `pdb.enabled` | `boolean` | No | `true` | Enable (`true`) or disable (`false`) PodDisruptionBudget |
| `pdb.maxUnavailable` | `string` | No | `1` | Maximum number of unavailable pods (number or percentage of pods) ** |
| `pdb.minAvailable` | `string` | No | - | Minimum number of available pods (number or percentage of pods), overrides maxUnavailable ** |
| `hpa.enabled` | `boolean` | No | `true` | Enable (`true`) or disable (`false`) HoritzontalPodAutoscaler |
| `hpa.minReplicas` | `int` | No | `2` | Minimum number of replicas |
| `hpa.maxReplicas` | `int` | No | `4` | Maximum number of replicas |
| `hpa.resourceName` | `string` | No | `cpu` | Resource used for autoscale (cpu/memory) |
| `hpa.resourceUtilization` | `int` | No | `90` | Percentage usage of the resource used for autoscale |
| `replicas` | `int` | No | `2` | Number of replicas (ignored if hpa is enabled) |
| `replicas` | `int` | No | `2` | Number of replicas |
| `resources.requests.cpu` | `string` | No | `250m` | Override CPU requests |
| `resources.requests.memory` | `string` | No | `250Mi` | Override Memory requests |
| `resources.limits.cpu` | `string` | No | `1` | Override CPU limits |
| `resources.limits.memory` | `string` | No | `512Mi` | Override Memory limits |
| `livenessProbe.initialDelaySeconds` | `int` | No | `3` | Override liveness initial delay (seconds) |
| `livenessProbe.timeoutSeconds` | `int` | No | `1` | Override liveness timeout (seconds) |
| `livenessProbe.periodSeconds` | `int` | No | `10` | Override liveness period (seconds) |
| `livenessProbe.successThreshold` | `int` | No | `1` | Override liveness success threshold |
| `livenessProbe.failureThreshold` | `int` | No | `3` | Override liveness failure threshold |
| `readinessProbe.initialDelaySeconds` | `int` | No | `3` | Override readiness initial delay (seconds) |
| `readinessProbe.timeoutSeconds` | `int` | No | `1` | Override readiness timeout (seconds) |
| `readinessProbe.periodSeconds` | `int` | No | `10` | Override readiness period (seconds) |
| `readinessProbe.successThreshold` | `int` | No | `1` | Override readiness success threshold |
| `readinessProbe.failureThreshold` | `int` | No | `3` | Override readiness failure threshold |
| `ingress.path` | `string` | No | `/api_docs/proxy` | Proxy path to configure on Nginx Ingress |
| `ingress.host` | `string` | Yes | - | Host to configure on Nginx Ingress |
| `grafanaDashboard.label.key` | `string` | No | `monitoring-key` | Label `key` used by grafana-operator for dashboard discovery |
| `grafanaDashboard.label.value` | `string` | No | `middleware` | Label `value` used by grafana-operator for dashboard discovery |

** If you are already using `pdb.maxUnavailable` and want to use `pdb.minAvailable` (or the other way around), due to ansible operator limitation of doing patch operation (if objects already exist), operator will receive an error when managing PDB object because although the spec of the PDB resource it creates is correct, operator will try to patch an existing object which already has the other variable, and these two variables `pdb.maxUnavailable`/`pdb.minAvailable` are mutually exclusive and cannot coexists on the same PDB. To solve that situation:
  - Configure `pdb.enabled=false` (so operator will delete associated PDB, and then re-enable it with `pdb.enabled=true` setting desired PDB field `pdb.minAvailable` or `pdb.maxUnavailable`. so operator will create it from scratch on next reconcile
  - Or, delete manually associated PDB object, and operator will create it from scratch on next reconcile