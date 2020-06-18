# Getting started

## Controllers management

* Operator initial creation:

```bash
$ operator-sdk new threescale-saas-operator --api-version=saas.3scale.net/v1alpha1 --kind=AutoSSL --type=ansible
```

* New API addition:
```bash
$ operator-sdk add api --api-version=saas.3scale.net/v1alpha1 --kind=Backend
```

## Operator image

* Apply changes on Operator ([ansible roles](../roles/)), then create a new operator image and push it to the registry with:
```bash
$ make operator-image-update
```
* Operator images are available [here](https://quay.io/repository/3scale/3scale-saas-operator?tab=tags)

## Operator deploy

* Deploy operator (namespace, CRD, service account, role, role binding and operator deployment):
```bash
$ make operator-deploy
```
* Create any CR resource
* Once tested, delete created operator objects (except CRD/namespace for caution):
```bash
$ make operator-delete
```
