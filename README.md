[![Go Report Card](https://goreportcard.com/badge/github.com/soerenschneider/acmevault)](https://goreportcard.com/report/github.com/soerenschneider/acmevault)

# acmevault

## Problem Statement

Rolling out TLS encryption shouldn't need to be pitched anymore (even for internal services). Using the DNS01 ACME challenge is proven and allows issuing certs non-public routable machines. On the other hand, you need to have access to either highly-privileged/narrowly-scoped credentials of your DNS provider to solve these DNS01 challenges.

In the case of Route53, if you don't want to end up creating dozens of hosted zones, one for each of your subdomains, you're at risk of leaking highly-privileged IAM credentials.

Acmevault requests short-lived IAM credentials for Route53 and uses them to perform DNS01 challenges for the configured domains and writes the issued X509 certificates to Hashicorp Vault's K/V secret store - only readable by the appropriate AppRole.

Its client mode reads the respective written certificates from Vault and installs them to a preconfigured location, optionally invoking post-installation hooks.


## Overview
![Overview](overview.png)

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
