# Configuration

Acmevault is configured using a YAML file, passed via `-config <path>` or the environment variable `ACMEVAULT_CONFIG_FILE`.
Every option can be overridden by an environment variable, see the reference below.

## Configuration example
```yaml
email: my@email.tld
acmeUrl: https://acme-staging-v02.api.letsencrypt.org/directory
intervalSeconds: 43200
metricsAddr: 127.0.0.1:9112
domains:
  - domain: domain1.tld
    sans:
      - domain3.tld
      - domain4.tld
  - domain: domain2.tld
vault:
  addr: https://vault:8200
  authMethod: approle
  roleId: my-role-id
  secretIdFile: /run/secrets/acmevault-secret-id
  pathPrefix: preprod
  kv2MountPath: secret
  awsMountPath: aws
  awsRole: acmevault
```

## Configuration reference

### General

| Keyword              | Environment variable                | Description                                                                                                  | Default                                          | Mandatory |
|----------------------|-------------------------------------|--------------------------------------------------------------------------------------------------------------|--------------------------------------------------|-----------|
| email                | ACMEVAULT_ACME_EMAIL                | Email address used to register the ACME account                                                              |                                                  | Y         |
| acmeUrl              | ACMEVAULT_ACME_URL                  | ACME directory URL. Either Let's Encrypt production or `https://acme-staging-v02.api.letsencrypt.org/directory` | https://acme-v02.api.letsencrypt.org/directory   | N         |
| acmeCustomDnsServers | ACMEVAULT_ACME_CUSTOM_DNS_SERVERS   | List of IP addresses of recursive nameservers used to check DNS01 challenge propagation                      |                                                  | N         |
| acmeDnsProvider      | ACMEVAULT_ACME_DNS_PROVIDER         | Currently ignored, Route53 is always used                                                                    |                                                  | N         |
| intervalSeconds      | ACMEVAULT_INTERVAL_SECONDS          | Interval between checks of all certificates, between 3600 and 86400                                          | 43200                                            | N         |
| metricsAddr          | ACMEVAULT_METRICS_ADDR              | Address to serve Prometheus metrics on (`/metrics`). Set to an empty string to disable                       | 127.0.0.1:9112                                   | N         |
| verbose              | ACMEVAULT_VERBOSE                   | Enable debug logging                                                                                         | false                                            | N         |
| domains              |                                     | List of certificates to issue, see [Domains](#domains)                                                       |                                                  | Y         |
| vault                |                                     | Vault configuration, see [Vault](#vault)                                                                     |                                                  | Y         |

### Domains

| Keyword | Description                                         | Example                      | Mandatory |
|---------|-----------------------------------------------------|------------------------------|-----------|
| domain  | FQDN of the certificate's common name               | domain1.tld                  | Y         |
| sans    | List of FQDNs to add as Subject Alternative Names   | [domain3.tld, domain4.tld]   | N         |

### Vault

| Keyword          | Environment variable                     | Description                                                                                                                    | Default                                   | Mandatory                  |
|------------------|------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------|----------------------------|
| addr             | ACMEVAULT_VAULT_ADDR                     | Address of Vault                                                                                                               | value of `VAULT_ADDR`                     | Y                          |
| authMethod       | ACMEVAULT_VAULT_AUTH_METHOD              | One of `token`, `approle`, `kubernetes`, `implicit` (reads the token from `~/.vault-token`)                                    |                                           | Y                          |
| token            | ACMEVAULT_VAULT_TOKEN                    | Vault token                                                                                                                    | value of `VAULT_TOKEN`                    | if authMethod is `token`   |
| roleId           | ACMEVAULT_VAULT_APPROLE_ROLE_ID          | [AppRole](https://developer.hashicorp.com/vault/docs/auth/approle) role id, the AppRole auth method must be mounted at `approle` |                                           | if authMethod is `approle` |
| secretId         | ACMEVAULT_VAULT_APPROLE_SECRET_ID        | AppRole secret id. Mutually exclusive with `secretIdFile`                                                                      |                                           | if authMethod is `approle` |
| secretIdFile     | ACMEVAULT_VAULT_APPROLE_SECRET_ID_FILE   | File to read the AppRole secret id from. Mutually exclusive with `secretId`                                                    |                                           | if authMethod is `approle` |
| k8sRoleId        | ACMEVAULT_VAULT_K8S_ROLE_ID              | Role of the [Kubernetes auth method](https://developer.hashicorp.com/vault/docs/auth/kubernetes)                               |                                           | if authMethod is `kubernetes` |
| k8sMountPath     | ACMEVAULT_VAULT_K8S_MOUNT                | Mount path of the Kubernetes auth method                                                                                       | kubernetes                                | N                          |
| pathPrefix       | ACMEVAULT_VAULT_PATH_PREFIX              | Prefix for all data written to the K/V store, allows running multiple instances. Must not start or end with `/`                | host of `acmeUrl`, e.g. acme-v02.api.letsencrypt.org | N               |
| domainPathFormat | ACMEVAULT_VAULT_DOMAIN_PATH_FORMAT       | `fmt` format string to build the per-domain path segment, must contain `%`, e.g. `cert-%s`                                     | domain as is                              | N                          |
| kv2MountPath     | ACMEVAULT_VAULT_KV2_MOUNT                | Mount path of the K/V version 2 secrets engine                                                                                 | secret                                    | Y                          |
| awsMountPath     | ACMEVAULT_VAULT_AWS_MOUNT                | Mount path of the AWS secrets engine                                                                                           | aws                                       | Y                          |
| awsRole          | ACMEVAULT_VAULT_AWS_ROLE                 | Role of the AWS secrets engine to request Route53 credentials from                                                             | acmevault                                 | Y                          |

See [Vault setup](vault.md) for the required Vault configuration and the layout of the stored data.
