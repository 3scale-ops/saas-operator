# Instructions

1. Create the file `config/local-setup/env-inputs/pull-secrets.json` with the registry auths required (for private repositories used). Example:

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

2. Create the file `config/local-setup/env-inputs/seed.env` with the following contents. Change the values to your heart's content:

```bash
MYSQL_ROOT_PASSWORD=password
MYSQL_DATABASE=system_enterprise
MYSQL_USER=app
MYSQL_PASSWORD=password
POSTGRES_USER=app
POSTGRES_PASSWORD=password
POSTGRES_DB=zync
BACKEND_INTERNAL_API_USER=user
BACKEND_INTERNAL_API_PASSWORD=password
SYSTEM_MASTER_USER=admin
SYSTEM_MASTER_PASSWORD=master-pass
SYSTEM_MASTER_ACCESS_TOKEN=mtoken
SYSTEM_TENANT_USER=admin
SYSTEM_TENANT_PASSWORD=provider-pass
SYSTEM_TENANT_TOKEN=ptoken
SYSTEM_APICAST_TOKEN=atoken
SYSTEM_EVENTS_SHARED_SECRET=password
SYSTEM_ASSETS_S3_ACCESS_KEY=admin
SYSTEM_ASSETS_S3_SECRET_KEY=admin1234
SYSTEM_SECRET_KEY_BASE=xxxxx
SYSTEM_DATABASE_SECRET=xxxxx
SYSTEM_SMTP_USER=""
SYSTEM_SMTP_PASSWORD=""
SYSTEM_ACCESS_CODE=""
SYSTEM_SEGMENT_DELETION_TOKEN=""
SYSTEM_SEGMENT_WRITE_KEY=""
SYSTEM_GITHUB_CLIENT_ID=""
SYSTEM_GITHUB_CLIENT_SECRET=""
SYSTEM_RH_CUSTOMER_PORTAL_CLIENT_ID=""
SYSTEM_RH_CUSTOMER_PORTAL_CLIENT_SECRET=""
SYSTEM_BUGSNAG_API_KEY=""
SYSTEM_RECAPTCHA_PUBLIC_KEY=""
SYSTEM_RECAPTCHA_PRIVATE_KEY=""
ZYNC_SECRET_KEY_BASE=xxxxx
ZYNC_AUTH_TOKEN=ztoken
ZYNC_BUGSNAG_API_KEY=""
```

3. You can tweak configurations in `config/local-setup/env-inputs/configuration.yaml`.


4. Issue the following commands

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
