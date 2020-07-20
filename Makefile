.PHONY: operator-image-update operator-deploy operator-delete help

.DEFAULT_GOAL := help

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
THISDIR_PATH := $(patsubst %/,%,$(abspath $(dir $(MKFILE_PATH))))
UNAME := $(shell uname)

ifeq (${UNAME}, Linux)
  INPLACE_SED=sed -i
else ifeq (${UNAME}, Darwin)
  INPLACE_SED=sed -i ""
endif

CURRENT_GIT_REF := $(shell git describe --always --dirty)
TAG ?= $(CURRENT_GIT_REF)
REGISTRY ?= quay.io
ORG ?= 3scale
PROJECT ?= 3scale-saas-operator
IMAGE ?= $(REGISTRY)/$(ORG)/$(PROJECT)
KUBE_CLIENT ?= kubectl # It can be used "oc" or "kubectl"
NAMESPACE ?= 3scale-example

## Operator ##
operator-image-build: ## OPERATOR IMAGE - Build operator Docker image
	operator-sdk build $(IMAGE):$(TAG)

operator-image-push: ## OPERATOR IMAGE - Push operator Docker image to remote registry
	docker push $(IMAGE):$(TAG)

operator-image-update: operator-image-build operator-image-push ## OPERATOR IMAGE - Build and Push Operator Docker image to remote registry

namespace-create: # NAMESPACE MANAGEMENT - Create namespace for the operator
	$(KUBE_CLIENT) create namespace $(NAMESPACE) || true
	$(KUBE_CLIENT) label namespace $(NAMESPACE) monitoring-key=middleware || true

operator-deploy: namespace-create ## OPERATOR DEPLOY - Deploy Operator objects (namespace, CRDs, service account, role, role binding and operator deployment)
	$(KUBE_CLIENT) apply -f deploy/crds/saas.3scale.net_autossls_crd.yaml --validate=false || true
	$(KUBE_CLIENT) apply -f deploy/crds/saas.3scale.net_backends_crd.yaml --validate=false || true
	$(KUBE_CLIENT) apply -f deploy/crds/saas.3scale.net_zyncs_crd.yaml --validate=false || true
	$(KUBE_CLIENT) apply -f deploy/crds/saas.3scale.net_corsproxies_crd.yaml --validate=false || true
	$(KUBE_CLIENT) apply -f deploy/crds/saas.3scale.net_echoapis_crd.yaml --validate=false || true
	$(KUBE_CLIENT) apply -f deploy/service_account.yaml -n $(NAMESPACE)
	$(KUBE_CLIENT) apply -f deploy/role.yaml -n $(NAMESPACE)
	$(KUBE_CLIENT) apply -f deploy/role_binding.yaml -n $(NAMESPACE)
	$(INPLACE_SED) 's|REPLACE_IMAGE|$(IMAGE):$(TAG)|g' deploy/operator.yaml
	$(KUBE_CLIENT) apply -f deploy/operator.yaml -n $(NAMESPACE)
	$(INPLACE_SED) 's|$(IMAGE):$(TAG)|REPLACE_IMAGE|g' deploy/operator.yaml

operator-delete: ## OPERATOR DEPLOY - Delete Operator objects (except CRD/namespace for caution)
	$(KUBE_CLIENT) delete -f deploy/operator.yaml -n $(NAMESPACE) || true
	$(KUBE_CLIENT) delete -f deploy/role_binding.yaml -n $(NAMESPACE) || true
	$(KUBE_CLIENT) delete -f deploy/role.yaml -n $(NAMESPACE) || true
	$(KUBE_CLIENT) delete -f deploy/service_account.yaml -n $(NAMESPACE) || true

help: ## Print this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-33s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
