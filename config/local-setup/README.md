# Instructions

1. Create the file `config/local-setup/secrets/pull-secrets.json` with the registry auths required (for private repositories used). Example:

```json
{
    "auths": {
        "quay.io": {
            "auth": "token1"
        },
        "brew.registry.redhat.io": {
            "auth": "token2"
        }
    }
}
```

1. Issue the following commands

```bash
make kind-create
make kind-local-setup
```

## URLs

* **backend**: http://backend-172-27-27-100.nip.io
* **apicast**:
  * http://*.production-172-27-27-101.nip.io
  * http://*.staging-172-27-27-102.nip.io
* **echo-api**: http://echo-api-172-27-27-103.nip.io
* **autossl**: http://autossl-172-27-27-104.nip.io
* **system**: https://<tenant>-admin.system-172-27-27-105.nip.io


## TODO

* System SMPTP?? maybe just disable altogether
