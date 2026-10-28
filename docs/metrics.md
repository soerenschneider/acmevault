# Metrics

Metrics are served in Prometheus format on `http://<metricsAddr>/metrics`, see [configuration](configuration.md).

All metrics are prefixed with the namespace `acmevault_`.

## Server

| Metric Name                                       | Description                                                       | Type          | Labels       |
|---------------------------------------------------|-------------------------------------------------------------------|---------------|--------------|
| server_latest_iteration_time_seconds              | Timestamp of the latest check of all certificates                 | Gauge         |              |
| server_certificates_retrieved_total               | Total number of certificates retrieved from the ACME provider      | Counter       |              |
| server_certificate_retrieve_errors_total          | Total errors while retrieving certificates from the ACME provider | Counter       |              |
| server_certificates_renewals_total                | Total number of renewed certificates                              | Counter       |              |
| server_certificates_renewal_errors_total          | Total errors while renewing certificates                          | Counter       |              |
| server_certificates_written_total                 | Total number of certificates written to Vault                     | Counter (Vec) | subsystem    |
| server_certificates_write_errors_total            | Total errors while writing certificates to Vault                  | Counter (Vec) | subsystem    |
| server_certificate_expiry_time                    | Timestamp of certificate expiry                                   | Gauge (Vec)   | domain       |
| server_certificate_errors_total                   | Total number of errors while handling certificates                | Counter (Vec) | domain, desc |
| server_vault_aws_credentials_requested_total      | Total number of dynamic AWS credentials requested                 | Counter       |              |
| server_vault_aws_credentials_request_errors_total | Total errors while requesting dynamic AWS credentials             | Counter       |              |

## Vault token renewal

Only updated for the auth methods `approle` and `kubernetes`, which log in and renew their token automatically.

| Metric Name                      | Description                                     | Type    | Labels |
|----------------------------------|-------------------------------------------------|---------|--------|
| vault_renewal_logins_total       | Total number of logins to Vault                 | Counter |        |
| vault_renewal_login_errors_total | Total errors while logging in to Vault          | Counter |        |
| vault_renewal_renewals_total     | Total number of token renewals                  | Counter |        |
| vault_renewal_renew_errors_total | Total errors while renewing the token           | Counter |        |
| vault_renewal_token_ttl_seconds  | Lease duration of the token after renewal       | Gauge   |        |
