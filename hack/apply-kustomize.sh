#!/bin/bash

set -eux

function generate_resources() {
    local KPATH=${1}
    local RESFILE=$(mktemp)
    ${KUSTOMIZE_BIN} build ${KUSTOMIZE_OPTIONS} ${KPATH} > ${RESFILE}
    echo ${RESFILE}
}

function filter_resources() {
    local RESFILE=${1}
    local FILTER=${2}
    cat ${RESFILE} | ${YQ_BIN} ". | select(${FILTER})"
}

function resource_names() {
    local RESFILE=${1}
    local FILTER=${2}
    filter_resources ${RESFILE} "${FILTER}" | ${YQ_BIN} -N '[.metadata.namespace,.metadata.name] | join("/")'
}

function deploy_crds() {
    local RESFILE=${1}
    local FILTER=".kind == \"CustomResourceDefinition\""
    if [[ $(resource_names ${RESFILE} "${FILTER}") != "/" ]]; then
        echo; echo "#################### > Deploying CRDs for ${NAME}"
        filter_resources ${RESFILE} "${FILTER}" | kubectl apply -f -
        resource_names ${RESFILE} "${FILTER}" | cut -f2 -d/ | xargs kubectl wait --for condition=established --timeout=60s crd
    fi
}

function wait_for() {
    local KIND=${1}
    # local NS=${2}
    FILTER=".kind == \"${KIND}\""
    if [[ $(resource_names ${RESFILE} "${FILTER}") != "/" ]]; then
        for ITEM in $(resource_names ${RESFILE} "${FILTER}"); do
            local NAME=${ITEM#*/}
            local NS=${ITEM%/*}
            echo; echo "#################### > Waiting for ${KIND} ${NAME} in namespace ${NS}"
            local SELECTOR=$(kubectl -n ${NS} describe ${KIND} ${NAME} | awk '/^Selector:/{print $2}')
            kubectl -n ${NS} get pods -l ${SELECTOR} --no-headers -o name | xargs kubectl -n ${NS} wait --for condition=ready
        done
    fi
}

function deploy_controller() {
    local RESFILE=${1}
    local FILTER=".kind != \"CustomResourceDefinition\" and .apiVersion != \"*${NAME}*\""
    if [[ $(resource_names ${RESFILE} "${FILTER}") != "/" ]]; then
        echo; echo "#################### > Deploying controller for ${NAME}"
        filter_resources ${RESFILE} "${FILTER}" | kubectl apply -f -
        for KIND in "Deployment" "StatefulSet"; do wait_for ${KIND}; done
    fi
}

function deploy_custom_resources() {
    local RESFILE=${1}
    local FILTER=".kind != \"CustomResourceDefinition\" and .apiVersion == \"*${NAME}*\""
    if [[ $(resource_names ${RESFILE} "${FILTER}") != "/" ]]; then
        echo; echo "#################### > Deploying custom resources for ${NAME}"
        filter_resources ${RESFILE} "${FILTER}" | kubectl apply -f -
    fi
}


test -n "${KUSTOMIZE_BIN}" || (echo "KUSTOMIZE_BIN envvar must be set" && exit -1)
test -n "${YQ_BIN}" || (echo "YQ_BIN envvar must be set" && exit -1)
test -n "${BASE_PATH}" || (echo "BASE_PATH envvar must be set" && exit -1)

KUSTOMIZE_OPTIONS="--enable-helm"
NAME=${1}
RESFILE=$(generate_resources ${BASE_PATH}/${NAME})
# resource_names release.yaml ".kind == \"StatefulSet\""
deploy_crds ${RESFILE}
deploy_controller ${RESFILE}
deploy_custom_resources ${RESFILE}
rm -f ${RESFILE}
