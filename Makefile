# VERSION defines the project version for the bundle.
# Update this value when you upgrade the version of your project.
# To re-generate a bundle for another specific version without changing the standard setup, you can:
# - use the VERSION as arg of the bundle target (e.g make bundle VERSION=0.0.2)
# - use environment variables to overwrite this value (e.g export VERSION=0.0.2)
VERSION ?= 0.22.0
# CHANNELS define the bundle channels used in the bundle.
# Add a new line here if you would like to change its default config. (E.g CHANNELS = "candidate,fast,stable")
# To re-generate a bundle for other specific channels without changing the standard setup, you can:
# - use the CHANNELS as arg of the bundle target (e.g make bundle CHANNELS=candidate,fast,stable)
# - use environment variables to overwrite this value (e.g export CHANNELS="candidate,fast,stable")
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif

# DEFAULT_CHANNEL defines the default channel used in the bundle.
# Add a new line here if you would like to change its default config. (E.g DEFAULT_CHANNEL = "stable")
# To re-generate a bundle for any other default channel without changing the default setup, you can:
# - use the DEFAULT_CHANNEL as arg of the bundle target (e.g make bundle DEFAULT_CHANNEL=stable)
# - use environment variables to overwrite this value (e.g export DEFAULT_CHANNEL="stable")
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

# IMAGE_TAG_BASE defines the docker.io namespace and part of the image name for remote images.
# This variable is used to construct full image tags for bundle and catalog images.
#
# For example, running 'make bundle-build bundle-push catalog-build catalog-push' will build and push both
# 3scale.net/operator-sdk-barebones-bundle:$VERSION and 3scale.net/operator-sdk-barebones-catalog:$VERSION.
IMAGE_TAG_BASE ?= quay.io/3scale/saas-operator

# BUNDLE_IMG defines the image:tag used for the bundle.
# You can use it as an arg. (E.g make bundle-build BUNDLE_IMG=<some-registry>/<project-name-bundle>:<tag>)
BUNDLE_IMG ?= $(IMAGE_TAG_BASE)-bundle:v$(VERSION)

# BUNDLE_GEN_FLAGS are the flags passed to the operator-sdk generate bundle command
BUNDLE_GEN_FLAGS ?= -q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)

# USE_IMAGE_DIGESTS defines if images are resolved via tags or digests
# You can enable this value if you would like to use SHA Based Digests
# To enable set flag to true
USE_IMAGE_DIGESTS ?= false
ifeq ($(USE_IMAGE_DIGESTS), true)
    BUNDLE_GEN_FLAGS += --use-image-digests
endif

# Image URL to use all building/pushing image targets
IMG ?= $(IMAGE_TAG_BASE):v$(VERSION)

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.27

# KIND_K8S_VERSION refers to the version of the kind k8s cluster for e2e testing.
KIND_K8S_VERSION = v1.27.0

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec


OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m | sed 's/x86_64/amd64/')

all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./api/..." output:crd:artifacts:config=config/crd/bases
	$(CONTROLLER_GEN)  rbac:roleName=manager-role crd webhook paths="./controllers/..." output:crd:artifacts:config=config/crd/bases

generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

TEST_PKG = ./api/... ./controllers/... ./pkg/...
KUBEBUILDER_ASSETS = "$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)"

test/assets/external-apis/crds.yaml: kustomize
	mkdir -p $(@D)
	$(KUSTOMIZE) build config/dependencies/external-secrets-crds > $@
	echo "---" >> $@ && $(KUSTOMIZE) build config/dependencies/grafana-crds >> $@
	echo "---" >> $@ && $(KUSTOMIZE) build config/dependencies/marin3r-crds >> $@
	echo "---" >> $@ && $(KUSTOMIZE) build config/dependencies/prometheus-crds >> $@
	echo "---" >> $@ && $(KUSTOMIZE) build config/dependencies/tekton-crds >> $@

test: manifests generate fmt vet envtest assets ginkgo test/assets/external-apis/crds.yaml ## Run tests.
	KUBEBUILDER_ASSETS=$(KUBEBUILDER_ASSETS) $(GINKGO) -p -r $(TEST_PKG)  -coverprofile cover.out

test-debug: manifests generate fmt vet envtest assets ginkgo test/assets/external-apis/crds.yaml ## Run tests.
	KUBEBUILDER_ASSETS=$(KUBEBUILDER_ASSETS) $(GINKGO) -v -r $(TEST_PKG)  -coverprofile cover.out

TEST_E2E_DEPLOY = marin3r-crds prometheus-crds tekton-crds grafana-crds external-secrets-crds minio
test-e2e: export KUBECONFIG = $(PWD)/kubeconfig
test-e2e: manifests ginkgo kind-create $(foreach elem,$(TEST_E2E_DEPLOY),install-$(elem)) kind-deploy-controller kind-load-redis-with-ssh ## Runs e2e tests
	$(GINKGO) -p -r ./test/e2e
	$(MAKE) kind-delete

assets: go-bindata ## assets: Generate embedded assets
	@echo Generate Go embedded assets files by processing source
	PATH=$$PATH:$$PWD/bin go generate github.com/3scale-ops/saas-operator/pkg/assets

##@ Build

build: generate fmt vet assets ## Build manager binary.
	go build -o bin/manager main.go

run: manifests generate fmt vet assets ## Run a controller from your host.
	LOG_MODE="development" go run ./main.go

docker-build: ## Build docker image with the manager.
	docker build -t ${IMG} .

docker-push: ## Push docker image with the manager.
	docker push ${IMG}

# PLATFORMS defines the target platforms for  the manager image be build to provide support to multiple
# architectures. (i.e. make docker-buildx IMG=myregistry/mypoperator:0.0.1). To use this option you need to:
# - able to use docker buildx . More info: https://docs.docker.com/build/buildx/
# - have enable BuildKit, More info: https://docs.docker.com/develop/develop-images/build_enhancements/
# - be able to push the image for your registry (i.e. if you do not inform a valid value via IMG=<myregistry/image:<tag>> than the export will fail)
# To properly provided solutions that supports more than one platform you should use this option.
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
.PHONY: docker-buildx
docker-buildx: test ## Build and push docker image for the manager for cross-platform support
	# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- docker buildx create --name project-v3-builder
	docker buildx use project-v3-builder
	- docker buildx build --push --platform=$(PLATFORMS) --tag ${IMG} -f Dockerfile.cross
	- docker buildx rm project-v3-builder
	rm Dockerfile.cross

##@ Deployment

install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | kubectl delete -f -

.PHONY: bundle
bundle: manifests kustomize operator-sdk ## Generate bundle manifests and metadata, then validate generated files.
	$(OPERATOR_SDK) generate kustomize manifests -q
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	$(KUSTOMIZE) build config/manifests | $(OPERATOR_SDK) generate bundle $(BUNDLE_GEN_FLAGS)
	$(OPERATOR_SDK) bundle validate ./bundle

.PHONY: bundle-build
bundle-build: ## Build the bundle image.
	docker build -f bundle.Dockerfile -t $(BUNDLE_IMG) .

.PHONY: bundle-push
bundle-push: ## Push the bundle image.
	$(MAKE) docker-push IMG=$(BUNDLE_IMG)

.PHONY: opm
OPM = ./bin/opm
opm: ## Download opm locally if necessary.
ifeq (,$(wildcard $(OPM)))
ifeq (,$(shell which opm 2>/dev/null))
	@{ \
	set -e ;\
	mkdir -p $(dir $(OPM)) ;\
	OS=$(shell go env GOOS) && ARCH=$(shell go env GOARCH) && \
	curl -sSLo $(OPM) https://github.com/operator-framework/operator-registry/releases/download/v1.23.0/$${OS}-$${ARCH}-opm ;\
	chmod +x $(OPM) ;\
	}
else
OPM = $(shell which opm)
endif
endif

# A comma-separated list of bundle images (e.g. make catalog-build BUNDLE_IMGS=example.com/operator-bundle:v0.1.0,example.com/operator-bundle:v0.2.0).
# These images MUST exist in a registry and be pull-able.
BUNDLE_IMGS ?= $(BUNDLE_IMG)

# The image tag given to the resulting catalog image (e.g. make catalog-build CATALOG_IMG=example.com/operator-catalog:v0.2.0).
CATALOG_IMG ?= $(IMAGE_TAG_BASE)-catalog:v$(VERSION)

# Default catalog base image to append bundles to
CATALOG_BASE_IMG ?= $(IMAGE_TAG_BASE)-catalog:latest

# Set CATALOG_BASE_IMG to an existing catalog image tag to add $BUNDLE_IMGS to that image.
ifneq ($(origin CATALOG_BASE_IMG), undefined)
FROM_INDEX_OPT := --from-index $(CATALOG_BASE_IMG)
endif

# Build a catalog image by adding bundle images to an empty catalog using the operator package manager tool, 'opm'.
# This recipe invokes 'opm' in 'semver' bundle add mode. For more information on add modes, see:
# https://github.com/operator-framework/community-operators/blob/7f1438c/docs/packaging-operator.md#updating-your-existing-operator
.PHONY: catalog-build
catalog-build: opm ## Build a catalog image.
	$(OPM) index add --container-tool docker --mode semver --tag $(CATALOG_IMG) --bundles $(BUNDLE_IMGS) $(FROM_INDEX_OPT)

# Push the catalog image.
.PHONY: catalog-push
catalog-push: ## Push a catalog image.
	$(MAKE) docker-push IMG=$(CATALOG_IMG)

##@ Release

prepare-alpha-release: bump-release generate fmt vet manifests assets bundle ## Generates bundle manifests for alpha channel release

prepare-stable-release: bump-release generate fmt vet manifests assets bundle refdocs ## Generates bundle manifests for stable channel release
	$(MAKE) bundle CHANNELS=alpha,stable DEFAULT_CHANNEL=alpha

bump-release: ## Write release name to "pkg/version" package
	sed -i 's/version string = "v\(.*\)"/version string = "v$(VERSION)"/g' pkg/version/version.go

bundle-publish: bundle-build bundle-push catalog-build catalog-push catalog-retag-latest ## Generates and pushes all required images for a release

get-new-release:
	@hack/new-release.sh v$(VERSION)

catalog-retag-latest:
	docker tag $(CATALOG_IMG) $(IMAGE_TAG_BASE)-catalog:latest
	$(MAKE) docker-push IMG=$(IMAGE_TAG_BASE)-catalog:latest

##@ Kind Deployment

kind-create: export KUBECONFIG = $(PWD)/kubeconfig
kind-create: kind ## Runs a k8s kind cluster
	docker inspect kind-saas-operator > /dev/null || docker network create -d bridge --subnet 172.27.27.0/24 kind-saas-operator
	KIND_EXPERIMENTAL_DOCKER_NETWORK=kind-saas-operator $(KIND) create cluster --wait 5m --image kindest/node:$(KIND_K8S_VERSION)

install-%: export KUBECONFIG = $(PWD)/kubeconfig
install-%: kustomize yq helm
	echo
	KUSTOMIZE_BIN=$(KUSTOMIZE) YQ_BIN=$(YQ) BASE_PATH=config/dependencies hack/apply-kustomize.sh $*

kind-delete: ## Deletes the kind cluster and the registry
kind-delete: kind
	$(KIND) delete cluster

CONTROLLER_DEPS = prometheus-crds grafana-crds
kind-deploy-controller: export KUBECONFIG = $(PWD)/kubeconfig
kind-deploy-controller: manifests kustomize docker-build $(foreach elem,$(CONTROLLER_DEPS),install-$(elem)) ## Deploy operator to the Kind K8s cluster
	$(KIND) load docker-image $(IMG)
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/test --load-restrictor LoadRestrictionsNone | kubectl apply -f -

kind-refresh-controller: export KUBECONFIG = ${PWD}/kubeconfig
kind-refresh-controller: manifests kind docker-build ## Reloads the controller image into the K8s cluster and deletes the old Pod
	$(KIND) load docker-image $(IMG)
	kubectl delete pod -l control-plane=controller-manager

kind-deploy-databases: export KUBECONFIG = $(PWD)/kubeconfig
kind-deploy-databases: kind-deploy-controller
	$(KUSTOMIZE) build config/local-setup/databases | kubectl apply -f -
	sleep 10
	kubectl wait --for condition=ready --timeout=300s pod --all

kind-undeploy: export KUBECONFIG = $(PWD)/kubeconfig
kind-undeploy: ## Undeploy controller from the Kind K8s cluster
	$(KUSTOMIZE) build config/test | kubectl delete -f -

REDIS_WITH_SSH_IMG = redis-with-ssh:6.2.13-alpine
kind-load-redis-with-ssh:
	docker build -t $(REDIS_WITH_SSH_IMG) test/assets/redis-with-ssh
	$(KIND) load docker-image $(REDIS_WITH_SSH_IMG)

kind-deploy-saas: export KUBECONFIG = ${PWD}/kubeconfig
kind-deploy-saas: kind-load-redis-with-ssh ## Deploys a 3scale SaaS dev environment
	$(KUSTOMIZE) build config/local-setup | kubectl apply -f -
	sleep 5
	kubectl wait --for condition=ready --timeout=300s pod system-console-0
	kubectl get pods --no-headers -o name | grep -v system | xargs kubectl wait --for condition=ready --timeout=300s
	kubectl -ti exec system-console-0 -c system-console -- bash -c '\
		MASTER_DOMAIN=multitenant-admin \
		MASTER_ACCESS_TOKEN=mtoken \
		MASTER_PASSWORD=mpass \
		MASTER_USER=admin \
		TENANT_NAME=provider \
		PROVIDER_NAME="3scale SaaS Dev Provider" \
		USER_LOGIN=admin \
		USER_PASSWORD=ppass \
		ADMIN_ACCESS_TOKEN=ptoken \
		USER_EMAIL="admin@cluster.local" \
		DISABLE_DATABASE_ENVIRONMENT_CHECK=1 \
		bundle exec rake db:setup'
	kubectl get pods --no-headers -o name | grep system | xargs kubectl wait --for condition=ready --timeout=300s

kind-cleanup-saas: export KUBECONFIG = ${PWD}/kubeconfig
kind-cleanup-saas:
	-$(KUSTOMIZE) build config/local-setup/databases | kubectl delete -f -
	-$(KUSTOMIZE) build config/local-setup | kubectl delete -f -
	-kubectl get pod --no-headers -o name | grep -v saas-operator | xargs kubectl delete --grace-period=0 --force
	-kubectl get pvc --no-headers -o name | xargs kubectl delete

LOCAL_SETUP_DEPS = metallb cert-manager marin3r prometheus-crds tekton-crds grafana-crds external-secrets-crds minio
kind-local-setup: export KUBECONFIG = ${PWD}/kubeconfig
kind-local-setup: $(foreach elem,$(LOCAL_SETUP_DEPS),install-$(elem)) kind-deploy-controller kind-deploy-databases kind-deploy-saas

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

export PATH := $(LOCALBIN):$(PATH)

## Tool Binaries
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest
GINKGO ?= $(LOCALBIN)/ginkgo
CRD_REFDOCS ?= $(LOCALBIN)/crd-ref-docs
KIND ?= $(LOCALBIN)/kind
GOBINDATA ?= $(LOCALBIN)/go-bindata
YQ ?= $(LOCALBIN)/yq
HELM ?= $(LOCALBIN)/helm

## Tool Versions
KUSTOMIZE_VERSION ?= v5.1.1
CONTROLLER_TOOLS_VERSION ?= v0.11.0
GINKGO_VERSION ?= v2.17.0
CRD_REFDOCS_VERSION ?= v0.0.8
KIND_VERSION ?= v0.16.0
ENVTEST_VERSION ?= latest
GOBINDATA_VERSION ?= latest
TEKTON_VERSION ?= v0.49.0
YQ_VERSION ?= v4.40.5
HELM_VERSION ?= v3.14.0

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	test -s $(KUSTOMIZE) || curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN)

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(CONTROLLER_GEN) || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	test -s $(ENVTEST) || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@$(ENVTEST_VERSION)

.PHONY: ginkgo
ginkgo: $(GINKGO) ## Download ginkgo locally if necessary
$(GINKGO):
	test -s $(GINKGO) || GOBIN=$(LOCALBIN) go install -mod=mod github.com/onsi/ginkgo/v2/ginkgo@$(GINKGO_VERSION)

.PHONY: crd-ref-docs
crd-ref-docs: ## Download crd-ref-docs locally if necessary
	test -s $(CRD_REFDOCS) || GOBIN=$(LOCALBIN) go install github.com/elastic/crd-ref-docs@$(CRD_REFDOCS_VERSION)

.PHONY: kind
kind: $(KIND) ## Download kind locally if necessary
$(KIND):
	test -s $(KIND) || GOBIN=$(LOCALBIN) go install sigs.k8s.io/kind@$(KIND_VERSION)

go-bindata: $(GOBINDATA) ## Download go-bindata locally if necessary.
$(GOBINDATA):
	test -s $(GOBINDATA) || GOBIN=$(LOCALBIN) go install github.com/go-bindata/go-bindata/...@$(GOBINDATA_VERSION)

.PHONY: yq
yq: $(YQ)
$(YQ):
	test -s $(YQ) || GOBIN=$(LOCALBIN) go install github.com/mikefarah/yq/v4@$(YQ_VERSION)

HELM_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3"
.PHONY: helm
helm: $(HELM)
$(HELM):
	curl -s $(HELM_INSTALL_SCRIPT) | HELM_INSTALL_DIR=$(LOCALBIN) bash -s -- --no-sudo --version $(HELM_VERSION)

##@ Other

.PHONY: operator-sdk
OPERATOR_SDK_RELEASE = v1.27.0
OPERATOR_SDK = bin/operator-sdk-$(OPERATOR_SDK_RELEASE)
operator-sdk: ## Download operator-sdk locally if necessary.
ifeq (,$(wildcard $(OPERATOR_SDK)))
ifeq (,$(shell which $(OPERATOR_SDK) 2>/dev/null))
	@{ \
	set -e ;\
	mkdir -p $(dir $(OPERATOR_SDK)) ;\
	OS=$(shell go env GOOS) && ARCH=$(shell go env GOARCH) && \
	curl -sSLo $(OPERATOR_SDK) https://github.com/operator-framework/operator-sdk/releases/download/$(OPERATOR_SDK_RELEASE)/operator-sdk_$${OS}_$${ARCH};\
	chmod +x $(OPERATOR_SDK) ;\
	}
else
OPERATOR_SDK = $(shell which $(OPERATOR_SDK))
endif
endif


refdocs: ## Generates api reference documentation from code
refdocs: crd-ref-docs
	$(CRD_REFDOCS) \
		--source-path=api \
		--max-depth=10 \
		--config=docs/api-reference/config.yaml \
		--renderer=asciidoctor \
		--output-path=docs/api-reference/reference.asciidoc
