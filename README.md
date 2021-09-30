[![Go Report Card](https://goreportcard.com/badge/github.com/soerenschneider/acmevault)](https://goreportcard.com/report/github.com/soerenschneider/acmevault)

# acmevault

## Problem Statement

Rolling out TLS encryption shouldn't need to be pitched anymore (even for internal services). Using the DNS01 ACME challenge is proven and allows issuing certs non-public routable machines. On the other hand, you need to have access to either highly-privileged/narrowly-scoped credentials of your DNS provider to solve these DNS01 challenges.

In the case of Route53, if you don't want to end up creating dozens of hosted zones, one for each of your subdomains, you're at risk of leaking highly-privileged IAM credentials.

Acmevault requests short-lived IAM credentials for Route53 and uses them to perform DNS01 challenges for the configured domains and writes the issued X509 certificates to Hashicorp Vault's K/V secret store - only readable by the appropriate AppRole.

Its client mode reads the respective written certificates from Vault and installs them to a preconfigured location, optionally invoking post-installation hooks.


## Overview
![Overview](overview.png)

# Server component
## Configuration example
```json
{
  "vaultAddr": "https://vault:8200",
  "metricsPath": "/var/lib/node_exporter/acmevault_server.prom",
  "roleId": "my_role_id",
  "secretId": "my_secret_id",
  "email": "my-acme-email@domain.tld"
}
```
## Configuration reference
| Keyword     | Description                                                                                      | Example                               | Mandatory |
|-------------|--------------------------------------------------------------------------------------------------|---------------------------------------|-----------|
| vaultAddr   | Connection string for vault                                                                      | https://vault:8200                    | Y         |
| roleId      | AppRole role id to login                                                                         | 988a9dfd-ea69-4a53-6cb6-9d6b86474bba  | Y         |
| secretId    | [AppRole secret id](https://www.vaultproject.io/docs/auth/approle) to authenticate against vault | 37b74931-c4cd-d49a-9246-ccc62d682a25  | Y         |
| email       | Email to register at ACME server                                                                 | your@email.tld                        | Y         |
| metricsPath | Path to write metrics to on filesystem                                                           | /var/lib/node_exporter/acmevault.prom | N         |

# Client component
## Configuration
```json
{
  "vaultAddr": "https://vault:8200",
  "metricsPath": "/var/lib/node_exporter/acmevault_client.prom",
  "user": "root",
  "group": "root",
  "certFile": "/etc/nginx/my_cert.crt",
  "privateKeyFile": "/etc/nginx/my_private_key.key",
  "roleId": "my_role_id",
  "secretId": "my_secret_id",
  "hooks": [
    "echo",
    "it works"
  ]
}
```

## Configuration reference

| Keyword        | Description                                                                                      | Example                               | Mandatory |
|----------------|--------------------------------------------------------------------------------------------------|---------------------------------------|-----------|
| vaultAddr      | Connection string for vault                                                                      | https://vault:8200                    | Y         |
| roleId         | [AppRole role id](https://www.vaultproject.io/docs/auth/approle) to authenticate against vault   | 988a9dfd-ea69-4a53-6cb6-9d6b86474bba  | Y         |
| secretId       | [AppRole secret id](https://www.vaultproject.io/docs/auth/approle) to authenticate against vault | 37b74931-c4cd-d49a-9246-ccc62d682a25  | Y         |
| user           | User that will own the written certificate and key on disk                                       | root                                  | Y         |
| group          | Group that will own the written certificate and key on disk                                      | root                                  | Y         |
| certFile       | The file path to write the certificate to                                                        | /etc/ssl/ssl-bundle.crt               | Y         |
| privateKeyFile | The file path to write the private key to                                                        | /etc/ssl/ssl-bundle.key               | Y         |
| hooks          | Commands to run after new cert files have been written                                           | ["echo", "it worked"]                 | N         |
| metricsPath    | Path on the disk to write metrics to                                                             | /var/lib/node_exporter/acmevault.prom | N         |

# Metrics

| Subsystem | Metric                                | Type    | Description                                                           | Labels            |
|-----------|---------------------------------------|---------|-----------------------------------------------------------------------|-------------------|
| server    | vault_aws_credentials_requested_total | counter | Total amount of dynamic AWS credentials requested                     |                   |
| server    | latest_iteration_time_seconds         | gauge   | Latest invocation of the server                                       |                   |
| server    | certificates_retrieved_total          | counter | Total amount of certificates retrieved from ACME provider             | domain            |
| server    | certificate_retrieve_errors_total     | counter | Total errors while trying to retrieve certificates from ACME provider | domain            |
| server    | certificates_renewals_total           | counter | Total number of renewed certificates                                  | domain            |
| server    | certificates_renewal_errors_total     | counter | Total errors while trying to renew certificates                       | domain            |
|           | certificates_written_total            | counter | Total number of certificates written total                            | domain, subsystem |
|           | certificates_write_errors_total       | counter | Total number of errors while writing the certificate                  | domain, subsystem |
|           | certificate_errors_total              | counter | Total number of errors while handling certificates                    | domain, desc      |
|           | certificate_expiry_time               | gauge   | Timestamp of certificate expiry                                       | domain            |
| client    | hooks_invocation_errors               | counter | Errors while invoking the hooks                                       |                   |
|           | timestamp                             | gauge   | Date of last measure                                                  |                   |

