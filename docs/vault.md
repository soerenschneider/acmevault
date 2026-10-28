# Vault Setup

Acmevault needs the following from Vault:

- a [K/V version 2](https://developer.hashicorp.com/vault/docs/secrets/kv/kv-v2) secrets engine to store the ACME account and the certificates
- an [AWS secrets engine](https://developer.hashicorp.com/vault/docs/secrets/aws) role to issue short-lived credentials for Route53
- an auth method for acmevault itself: AppRole, Kubernetes or a plain token

The examples below use the default mount paths `secret` and `aws` and the path prefix `preprod`.

## AWS Secrets Engine

The AWS role needs the permissions to solve DNS01 challenges in Route53:

```shell
$ vault secrets enable aws
$ vault write aws/config/root access_key=... secret_key=... region=us-east-1
$ vault write aws/roles/acmevault credential_type=iam_user policy_document=-<<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["route53:GetChange", "route53:ListHostedZonesByName"],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": ["route53:ChangeResourceRecordSets", "route53:ListResourceRecordSets"],
      "Resource": "arn:aws:route53:::hostedzone/<your-hosted-zone-id>"
    }
  ]
}
EOF
```

## Policies

### Acmevault
```hcl
# ACME account and certificates
path "secret/data/acmevault/preprod/*" {
  capabilities = ["create", "update", "read"]
}

# Dynamic Route53 credentials
path "aws/creds/acmevault" {
  capabilities = ["read"]
}
```

### Clients
Each client should only be able to read the certificate and private key of its own domain:
```hcl
path "secret/data/acmevault/preprod/client/domain1.tld/*" {
  capabilities = ["read"]
}
```

## Stored Data

All data is written to the K/V version 2 secrets engine below `acmevault/<pathPrefix>`:

| Path                                               | Content                                        |
|----------------------------------------------------|------------------------------------------------|
| `acmevault/<pathPrefix>/server/account/<email>`    | ACME account registration and its private key  |
| `acmevault/<pathPrefix>/client/<domain>/certificate` | Certificate data                             |
| `acmevault/<pathPrefix>/client/<domain>/privatekey`  | Private key of the certificate               |

The `<domain>` segment can be customized with `domainPathFormat`, see [configuration](configuration.md).

The certificate secret contains the keys `cert` (certificate), `issuer` (issuer certificate), `domain`, `url`, `stable_url` and `version`.
The private key secret contains the key `private_key`. Certificates and private keys are PEM encoded and additionally base64 encoded.

`cert` contains the certificate including the issuer chain.

## Consuming Certificates

Clients read their certificates directly from Vault, for example using a [Vault Agent](https://developer.hashicorp.com/vault/docs/agent-and-proxy/agent/template) template:

```hcl
template {
  destination = "/etc/ssl/domain1.tld.crt"
  contents    = <<EOT
{{ with secret "secret/data/acmevault/preprod/client/domain1.tld/certificate" }}{{ base64Decode .Data.data.cert }}{{ end }}
EOT
}

template {
  destination = "/etc/ssl/domain1.tld.key"
  perms       = "0600"
  contents    = <<EOT
{{ with secret "secret/data/acmevault/preprod/client/domain1.tld/privatekey" }}{{ base64Decode .Data.data.private_key }}{{ end }}
EOT
}
```

Or using the Vault CLI:
```shell
$ vault kv get -mount=secret -field=cert acmevault/preprod/client/domain1.tld/certificate | base64 -d
```
