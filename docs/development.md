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

### Running the operator locally

You can then run the operator with the following command:

```bash
make run
```

By default, will watch all Namespaces. To watch a specific namespace,
define the WATCH_NAMESPACE environment variable.

```
WATCH_NAMESPACE=3scale make run
```

## Publishing a new release

To build and publish a new release, update the `VERSION` env var and use the
`build-and-publish` Makefile target. This target is just a helper to execute
all the existing targets to publish a new release.

- Build a new image for the operator (`docker-build` and `docker-push` targets)
- Build a new image for the operator bundle (`bundle`, `bundle-buiild` and `bundle-push` targets)
- Add the new bundle to the 3scale Ops Bundle Catalog (`bundle-publish` targets)
-
```
VERSION=0.8.1 make build-and-publish
```

As example, the previous command will output:

``
docker build . -t quay.io/3scale/saas-operator:v0.8.1
Sending build context to Docker daemon  55.73MB
Step 1/5 : FROM quay.io/operator-framework/ansible-operator:v1.2.0
 ---> 43d6b2eb8daf
 ....
Successfully built 2290104fe359
Successfully tagged quay.io/3scale/saas-operator:v0.8.1
docker push quay.io/3scale/saas-operator:v0.8.1
The push refers to repository [quay.io/3scale/saas-operator]
51bed5bc6bd2: Layer already exists
...
v0.8.1: digest: sha256:6f911ceee1e969bee710c243766d52e87b7dfa20718424cd16cbda64fb286953 size: 2615
operator-sdk generate kustomize manifests -q
cd config/manager && /Users/rael/Code/gh/3scale/saas-operator/bin/kustomize edit set image controller=quay.io/3scale/saas-operator:v0.8.1
/Users/rael/Code/gh/3scale/saas-operator/bin/kustomize build config/manifests | operator-sdk generate bundle -q --overwrite --version 0.8.1 --channels=alpha --default-channel=alpha
INFO[0000] Building annotations.yaml
INFO[0000] Writing annotations.yaml in /Users/rael/Code/gh/3scale/saas-operator/bundle/metadata
INFO[0000] Building Dockerfile
...
Successfully tagged quay.io/3scaleops/saas-operator-bundle:v0.8.1
docker push quay.io/3scaleops/saas-operator-bundle:v0.8.1
The push refers to repository [quay.io/3scaleops/saas-operator-bundle]
51531c2e1381: Pushed
ec6d85c02179: Pushed
9b3ef2034ec2: Pushed
v0.8.1: digest: sha256:8213d7ba3e8e6458a34b2aa610dcca0f0376acdd435327a5dde6055c73133b3e size: 940
opm index add \
                --build-tool docker \
                --mode replaces \
                --bundles quay.io/3scaleops/saas-operator-bundle:v0.8.1 \
                --from-index quay.io/3scaleops/olm-catalog:bundle \
                --tag quay.io/3scaleops/olm-catalog:bundle
INFO[0000] building the index                            bundles="[quay.io/3scaleops/saas-operator-bundle:v0.8.1]"
INFO[0000] Pulling previous image quay.io/3scaleops/olm-catalog:bundle to get metadata  bundles="[quay.io/3scaleops/saas-operator-bundle:v0.8.1]"
INFO[0001] resolved name: quay.io/3scaleops/olm-catalog:bundle  bundles="[quay.io/3scaleops/saas-operator-bundle:v0.8.1]"
...
INFO[0038] loading bundle file                           dir=bundle_tmp213655740/manifests file=saas.3scale.net_autossls.yaml load=bundle
INFO[0038] loading bundle file                           dir=bundle_tmp213655740/manifests file=saas.3scale.net_backends.yaml load=bundle
INFO[0038] loading bundle file                           dir=bundle_tmp213655740/manifests file=saas.3scale.net_corsproxies.yaml load=bundle
INFO[0038] loading bundle file                           dir=bundle_tmp213655740/manifests file=saas.3scale.net_echoapis.yaml load=bundle
INFO[0038] loading bundle file                           dir=bundle_tmp213655740/manifests file=saas.3scale.net_mappingservices.yaml load=bundle
INFO[0038] loading bundle file                           dir=bundle_tmp213655740/manifests file=saas.3scale.net_systems.yaml load=bundle
INFO[0038] loading bundle file                           dir=bundle_tmp213655740/manifests file=saas.3scale.net_zyncs.yaml load=bundle
INFO[0038] Generating dockerfile                         bundles="[quay.io/3scaleops/saas-operator-bundle:v0.8.1]"
INFO[0038] writing dockerfile: index.Dockerfile221858170  bundles="[quay.io/3scaleops/saas-operator-bundle:v0.8.1]"
INFO[0038] running docker build                          bundles="[quay.io/3scaleops/saas-operator-bundle:v0.8.1]"
INFO[0038] [docker build -f index.Dockerfile221858170 -t quay.io/3scaleops/olm-catalog:bundle .]  bundles="[quay.io/3scaleops/saas-operator-bundle:v0.8.1]"
docker push quay.io/3scaleops/olm-catalog:bundle
The push refers to repository [quay.io/3scaleops/olm-catalog]
075ed80ec6da: Pushed
...
ace0eda3e3be: Layer already exists
bundle: digest: sha256:4dc02bcbfd229e736725612ce852a5755f77a897560011aa9beba611cf282e9c size: 1578
```
