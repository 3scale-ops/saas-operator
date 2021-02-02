# Getting started

## Deploy the controller using this repo

All the resouces required to run the operator are available at [config](config/default).

The Makefile targets render and apply the kustomize resources for you, but you can also generate the final YAMLs with `bin/kustomize build config/default`.

### Install the crds

Installs the custom resource definitions into the configured cluster in ~/.kube/config.

```bash
make install
```

### Deploy the controller

Deploys the controller, rbac and crds into the configured cluster in ~/.kube/config.

```bash
make deploy
```

### Uninstall the crds

Deletes the custom resource definitions from the configured cluster in ~/.kube/config.

```bash
make uninstall
```

### Undeploy the controller

Deletes the controller, rbac and crds from the configured cluster in ~/.kube/config.

```bash
make undeploy
```
