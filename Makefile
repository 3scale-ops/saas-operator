SHELL:=/bin/bash
# Current Operator version
VERSION ?= 0.10.1
# Default catalog image
CATALOG_IMG ?= quay.io/3scaleops/saas-operator-bundle:catalog
# Default bundle image tag
BUNDLE_IMG ?= quay.io/3scaleops/saas-operator-bundle:$(VERSION)
# Options for 'bundle-build'
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

# Image URL to use all building/pushing image targets
IMG ?= quay.io/3scale/saas-operator:$(VERSION)
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false"

all: manager

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

#############################
### Makefile requirements ###
#############################

OS=$(shell uname | awk '{print tolower($$0)}')
ARCH = $(shell arch)
ifeq ($(shell arch),x86_64)
ARCH := amd64
endif
ifeq ($(shell arch),aarch64)
ARCH := arm64
endif
ifeq ($(shell uname -m),x86_64)
ARCH := amd64
endif

# Download operator-sdk binary if necesasry
OPERATOR_SDK_RELEASE = v1.3.2
OPERATOR_SDK = $(shell pwd)/bin/operator-sdk-$(OPERATOR_SDK_RELEASE)
OPERATOR_SDK_DL_URL = https://github.com/operator-framework/operator-sdk/releases/download/$(OPERATOR_SDK_RELEASE)/operator-sdk_$(OS)_$(ARCH)
$(OPERATOR_SDK):
	curl -sL -o $(OPERATOR_SDK) $(OPERATOR_SDK_DL_URL)
	chmod +x $(OPERATOR_SDK)

# Download operator package manager if necessary
OPM_RELEASE = v1.17.0
OPM = $(shell pwd)/bin/opm-$(OPM_RELEASE)
OPM_DL_URL = https://github.com/operator-framework/operator-registry/releases/download/$(OPM_RELEASE)/$(OS)-$(ARCH)-opm
$(OPM):
	curl -sL -o $(OPM) $(OPM_DL_URL)
	chmod +x $(OPM)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

# Download controller-gen locally if necessary
CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen:
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1)

# Download kustomize locally if necessary
KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize:
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

# Download ginkgo locally if necessary
GINKGO = $(shell pwd)/bin/ginkgo
ginkgo:
	$(call go-get-tool,$(GINKGO),github.com/onsi/ginkgo/ginkgo)

GOBINDATA=$(shell pwd)/bin/go-bindata
go-bindata:
	$(call go-get-tool,$(GOBINDATA),github.com/go-bindata/go-bindata/...)

# Download kind locally if necessary
KIND = $(shell pwd)/bin/kind
kind:
	$(call go-get-tool,$(KIND),sigs.k8s.io/kind@v0.9.0)

# Download crd-ref-docs locally if necessary
CRD_REFDOCS = $(shell pwd)/bin/crd-ref-docs
$(CRD_REFDOCS):
	$(call go-get-tool,$(CRD_REFDOCS),github.com/elastic/crd-ref-docs@v0.0.6)

###########################
### Kubebuilder targets ###
###########################

# Run tests
ENVTEST_ASSETS_DIR=$(shell pwd)/testbin
test: generate fmt vet manifests assets ginkgo
	mkdir -p ${ENVTEST_ASSETS_DIR}
	test -f ${ENVTEST_ASSETS_DIR}/setup-envtest.sh || curl -sSLo ${ENVTEST_ASSETS_DIR}/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.7.0/hack/setup-envtest.sh
	source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh; \
		fetch_envtest_tools $(ENVTEST_ASSETS_DIR); \
		setup_envtest_env $(ENVTEST_ASSETS_DIR); \
		$(GINKGO) -p -r ./

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests assets
	go run ./main.go

# Install CRDs into a cluster
install: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests kustomize
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

# UnDeploy controller from the configured Kubernetes cluster in ~/.kube/config
undeploy:
	$(KUSTOMIZE) build config/default | kubectl delete -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths=./api/... output:crd:artifacts:config=config/crd/bases
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths=./controllers/... output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./api/..."

## assets: Generate embedded assets
# assets: export PATH=$(PATH):$(shell pwd)/bin
assets: go-bindata
	@echo Generate Go embedded assets files by processing source
	PATH=$$PATH:$$PWD/bin go generate github.com/3scale/saas-operator/pkg/assets

# Build the docker image
docker-build: test
	docker build -t ${IMG} .

# Push the docker image
docker-push:
	docker push ${IMG}

# Generate bundle manifests and metadata, then validate generated files.
.PHONY: bundle
bundle: $(OPERATOR_SDK) manifests kustomize
	$(OPERATOR_SDK) generate kustomize manifests -q
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	$(KUSTOMIZE) build config/manifests | $(OPERATOR_SDK) generate bundle -q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)
	$(OPERATOR_SDK) bundle validate ./bundle

# Build the bundle image.
.PHONY: bundle-build
bundle-build:
	docker build -f bundle.Dockerfile -t $(BUNDLE_IMG) .

tmp:
	mkdir -p $@

#########################
#### Release targets ####
#########################

prepare-alpha-release: bump-release generate fmt vet manifests bundle

prepare-stable-release: bump-release generate fmt vet manifests bundle
	$(MAKE) bundle CHANNELS=alpha,stable DEFAULT_CHANNEL=alpha

bump-release:
	sed -i 's/version string = "v\(.*\)"/version string = "v$(VERSION)"/g' pkg/version/version.go

bundle-push:
	docker push $(BUNDLE_IMG)

catalog-build: $(OPM)
	$(OPM) index add \
		--build-tool docker \
		--mode semver \
		--bundles $(BUNDLE_IMG) \
		--from-index $(CATALOG_IMG) \
		--tag $(CATALOG_IMG)

catalog-push:
	docker push $(CATALOG_IMG)

bundle-publish: bundle-build bundle-push catalog-build catalog-push


get-new-release:
	@hack/new-release.sh v$(VERSION)

############################################
#### Targets to manually test with Kind ####
############################################

kind-create: ## runs a k8s kind cluster for testing
kind-create: export KUBECONFIG = ${PWD}/kubeconfig
kind-create: tmp $(KIND)
	$(KIND) create cluster --wait 5m

kind-delete: ## deletes the kind cluster
kind-delete: $(KIND)
	$(KIND) delete cluster

kind-deploy: ## Deploys the operator in the kind cluster for testing
kind-deploy: export KUBECONFIG = ${PWD}/kubeconfig
kind-deploy: manifests kustomize kind
	$(KIND) load docker-image $(IMG)
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/test | kubectl apply -f -

kind-refresh-operator: ## Reloads the operator image into the cluster and deletes the old Pod
kind-refresh-operator: export KUBECONFIG = ${PWD}/kubeconfig
kind-refresh-operator: manifests kind
	$(KIND) load docker-image $(IMG)
	kubectl delete pod -l control-plane=controller-manager

kind-undeploy: ## Removes the operator from the kind cluster
kind-undeploy: export KUBECONFIG = ${PWD}/kubeconfig
kind-undeploy: manifests kustomize
	$(KUSTOMIZE) build config/test | kubectl delete -f -

############################
#### refdocs generation ####
############################

refdocs: $(CRD_REFDOCS) ## Generates api reference documentation from code
	$(CRD_REFDOCS) \
		--source-path=api \
		--config=docs/api-reference/config.yaml \
		--templates-dir=docs/api-reference/templates/asciidoctor \
		--renderer=asciidoctor \
		--output-path=docs/api-reference/reference.asciidoc
