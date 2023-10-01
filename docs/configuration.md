# Configuration

## Configuration example
```json
{
  "vaultAddr": "https://vault:8200",
  "vaultRoleId": "my_role_id",
  "vaultSecretId": "my_secret_id",
  "metricsPath": "/var/lib/node_exporter/acmevault_server.prom",
  "email": "my-acme-email@domain.tld"
}
```

## Configuration reference

### Vault

| Keyword     | Description                                                                                           | Example                               | Mandatory |
|-------------|-------------------------------------------------------------------------------------------------------|---------------------------------------|-----------|
| vaultAddr        | Connection string for vault                                                                      | https://vault:8200                    | Y         |
| vaultRoleId      | [AppRole role id](https://www.vaultproject.io/docs/auth/approle) to login                        | 988a9dfd-ea69-4a53-6cb6-9d6b86474bba  | Y         |
| vaultSecretId    | [AppRole secret id](https://www.vaultproject.io/docs/auth/approle) to authenticate against vault | 37b74931-c4cd-d49a-9246-ccc62d682a25  | Y         |
| vaultPathPrefix  | Path prefix for the K/V path in vault for this instance running acmevault                        | production                            | N         |
| email            | Email to register at ACME server                                                                 | your@email.tld                        | Y         |
| metricsPath      | Path to write metrics to on filesystem                                                           | /var/lib/node_exporter/acmevault.prom | N         |
| acmeUrl          | URL of the acme provider                                                                         | /var/lib/node_exporter/acmevault.prom | N         |
