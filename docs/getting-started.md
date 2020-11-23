# Getting started

## Deploy the controller using this repo

All the resouces requires to run the operator are available at [config](config/default)
and use kustomize to generate the final YAMLs.

By default rendered and applied on the fly by the different makefile targets,
but you can generate the final YAMLs using:

```
bin/kustomize build config/default
```

### Install the crds

Installs the custom resource definitions into the configured cluster in ~/.kube/config.

```
make install
```

### Deploy the controller

Deploys the controller, rbac and crds into the configured cluster in ~/.kube/config.

```
make deploy
```

### Uninstall the crds

Deletes the custom resource definitions from the configured cluster in ~/.kube/config.

```
make uninstall
```

### Undeploy the controller

Deletes the controller, rbac and crds from the configured cluster in ~/.kube/config.

```
make undeploy
```