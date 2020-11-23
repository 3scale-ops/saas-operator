# Development

## Help

The default makefile is `help` and will show all available targets documented.

```
$ make
help                           Print this help
run                            Run against the configured Kubernetes cluster in ~/.kube/config
install                        Install CRDs into a cluster
uninstall                      Uninstall CRDs from a cluster
deploy                         Deploy controller in the configured Kubernetes cluster in ~/.kube/config
undeploy                       Undeploy controller in the configured Kubernetes cluster in ~/.kube/config
docker-build                   Build the docker image
docker-push                    Push the docker image
kustomize                      Install kustomize binary if missing
ansible-operator               Install ansible-operator binary if missing
bundle                         Generate bundle manifests and metadata, then validate generated files.
bundle-build                   Build the bundle image.
```

## Prerequisites

### Kustomize

`kustomize` lets you customize raw, template-free YAML files for multiple purposes,
leaving the original YAML untouched and usable as is.

The binary can be installed using:

```
make kustomize
```

More info at https://kustomize.io

### ansible-operator

To run the operator locally you may need to install some ansible dependencies,
check the [Additional Prerequisites](https://sdk.operatorframework.io/docs/building-operators/ansible/installation/#additional-prerequisites)
section in the oficial documentation.

The binary can be installed using:

```
make ansible-operator
```

## Images

## Controller image

The operator controller image lifecycle is managed with the `docker-build` and `docker-push`
targets using the `IMG` and `VERSION` environment variables.

By default, the image is [quay.io/3scale/saas-operator](https://quay.io/3scale/saas-operator) the version `build`.

## Bundle image

The operator bundle image lifecycle is managed with the `bundle-build` and `bundle-push`
targets using the `BUNDLE_IMG` and `VERSION` environment variables.

By default, the image is [quay.io/3scaleops/saas-operator-bundle](https://quay.io/3scaleops/saas-operator-bundle) and the version `build`.

## Running the operator

## Running the operator locally

You can then run the operator with the following command:

```bash
make run
```

By default, will watch all Namespaces. To watch a specific namespace,
define the WATCH_NAMESPACE environment variable.

```
WATCH_NAMESPACE=3scale make run
```