# 3scale SaaS Operator

[![build status](https://circleci.com/gh/3scale/saas-operator.svg?style=shield)](https://circleci.com/gh/3scale/saas-operator)
[![release](https://badgen.net/github/release/3scale/saas-operator)](https://github.com/3scale/saas-operator/releases)
[![license](https://badgen.net/github/license/3scale/saas-operator)](https://github.com/3scale/saas-operator/blob/master/LICENSE)

A Kubernetes Operator based on the Operator SDK to manage 3scale SaaS (hosted version) on **Kubernetes/OpenShift**.

3scale SaaS controllers supported:

* Apicast
* AutoSSL
* Backend
* CORSProxy
* EchoAPI
* MappingService
* System
* Zync

## Requirements

* [prometheus-operator](https://github.com/coreos/prometheus-operator) v0.17.0+
* [grafana-operator](https://github.com/integr8ly/grafana-operator) v3.0.0+
* [secrets-manager](https://github.com/tuenti/secrets-manager) v1.0.2+
* [marin3r](https://github.com/3scale/marin3r) v0.4.0+
* [aws-nlb-helper-operator](https://github.com/3scale/aws-nlb-helper-operator) v0.2.0+

## Documentation

* [Getting started](docs/getting-started.md)
* [Apicast Custom Resource Reference](docs/apicast-crd-reference.md)
* [AutoSSL Custom Resource Reference](docs/autossl-crd-reference.md)
* [Backend Custom Resource Reference](docs/backend-crd-reference.md)
* [CORSProxy Custom Resource Reference](docs/corsproxy-crd-reference.md)
* [EchoAPI Custom Resource Reference](docs/echoapi-crd-reference.md)
* [MappingService Custom Resource Reference](docs/mappingservice-crd-reference.md)
* [System Custom Resource Reference](docs/system-crd-reference.md)
* [Zync Custom Resource Reference](docs/zync-crd-reference.md)

## Development

To run the operator locally you need to install some ansible dependencies first:

* ansible-runner: `sudo dnf install python-ansible-runner`
* ansible-runner-http: `pip install python-ansible-runner`
* openshift ansible module: `pip install openshift`

You can then run the operator with the following command:

```bash
operator-sdk run --local --watch-namespace <namespace>
```

## License

3scale SaaS Operator is under Apache 2.0 license. See the [LICENSE](LICENSE) file for details.
