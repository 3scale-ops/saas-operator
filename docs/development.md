# Development

## Images

## Controller image

The operator controller image lifecycle is managed with the `docker-build` and `docker-push`
targets using the `IMG` and `VERSION` environment variables.

By default, the image is [quay.io/3scale/saas-operator](https://quay.io/3scale/saas-operator) and the version is the one hardcoded in the Makefile, which should be changed for release of new versions.

## Bundle image

The operator bundle image lifecycle is managed with the `bundle-publish` target using the following variables:

* `VERSION`: the version of the olm release. Should always be the same as the operator version.
* `BUNDLE_IMG`: the bundle image for the olm release. Default is `quay.io/3scaleops/go-saas-operator-catalog:$(VERSION)`.
* `CATALOG_IMAGE`: the catalog image. Default is `quay.io/3scaleops/go-saas-operator-catalog:latest`.

## Running the operator

### Run the operator locally against a kind cluster

```bash
# create a kind cluster
make kind-create
# export the kubeconfig
export KUBECONFIG=$PWD/kubeconfig
# install operator CRDs
make install
# install required APIs CRDs
kubectl apply -f config/test/external-apis/
# run the operator
make run
```

### Running the inside a kind cluster

You can then run the operator with the following command:

```bash
# build the operator image
make docker-build
# create a kind cluster
make kind-create
# export the kubeconfig
export KUBECONFIG=$PWD/kubeconfig
# install operator CRDs
make kind-deploy
```
