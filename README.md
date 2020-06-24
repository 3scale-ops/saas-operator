# 3scale SaaS Operator

[![build status](https://circleci.com/gh/3scale/saas-operator.svg?style=shield)](https://circleci.com/gh/3scale/saas-operator)
[![release](https://badgen.net/github/release/3scale/saas-operator)](https://github.com/3scale/saas-operator/releases)
[![license](https://badgen.net/github/license/3scale/saas-operator)](https://github.com/3scale/saas-operator/blob/master/LICENSE)

A Kubernetes Operator based on the Operator SDK to manage 3scale SaaS (hosted version) on **Kubernetes/OpenShift**.

Current 3scale SaaS controllers supported:
* AutoSSL
* Backend
* Zync

Future 3scale SaaS controllers to be added:
* System
* Apicast
* MappingService
* CORSProxy
* PostFix

## Requirements

* [prometheus-operator](https://github.com/coreos/prometheus-operator) v0.17.0+
* [grafana-operator](https://github.com/integr8ly/grafana-operator) v3.0.0+
* [secrets-manager](https://github.com/tuenti/secrets-manager) v1.0.2+

## Documentation

* [Getting started](docs/getting-started.md)
* [AutoSSL Custom Resource Reference](docs/autossl-crd-reference.md)
* [Backend Custom Resource Reference](docs/backend-crd-reference.md)
* [Zync Custom Resource Reference](docs/zync-crd-reference.md)

## License

3scale SaaS Operator is under Apache 2.0 license. See the [LICENSE](LICENSE) file for details.
