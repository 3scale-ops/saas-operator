# AutoSSL Custom Resource Reference

## Full CR Example

Most of the fields do not need to be specified (can use default values), this is just an example of everything that can be overriden under your own risk:

```yaml
apiVersion: saas.3scale.net/v1alpha1
kind: AutoSSL
metadata:
  name: dev
spec:
  image:
    version: v1.0.0
    pullSecretName: quay-pull-secret
  replicas: 2
  acmeStaging: 1
  contactEmail: 3scale-operations@redhat.com
  proxyEndpoint: https://multitenant-admin.dev.3sca.net
  storageAdapter: redis
  redisHost: dev-saas-autossl-ec-redis.ng.0001.use1.cache.amazonaws.com
  redisPort: 6379
  verificationEndpoint: https://multitenant-admin.dev.3sca.net/swagger/spec.json
  domainWhitelist: autossl.dev.3sca.net
  domainBlacklist: blacklisted.dev.3sca.net
  logLevel: debug
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
      cpu: 75m
      memory: 64Mi
    limits:
      cpu: 150m
      memory: 128Mi
  externalDnsHostname: mtssl-edge-a.dev.3sca.net,mtssl-edge-b.dev.3sca.net,autossl.dev.3sca.net
  loadBalancer:
    proxyProtocol: "*"
    crossZoneLoadBalancingEnabled: true
    connectionDrainingEnabled: true
    connectionDrainingTimeout: 60
    connectionHealthcheckHealthyThreshold: 2
    connectionHealthcheckUnhealthyThreshold: 2
    connectionHealthcheckInterval: 5
    connectionHealthcheckTimeout: 3
    accessLogEnabled: true
    accessLogEmitInterval: 5
    accessLogS3BucketName: dev-s3-logs
    accessLogS3BucketPrefix: dev-ocp4-3/autossl
  grafanaDashboard:
    label:
      key: discovery
      value: enabled
```
