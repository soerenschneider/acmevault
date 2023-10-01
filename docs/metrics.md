# Metrics

All metrics are prefixed with the namespace `acmevault_`.

| Metric Name                                       | Description                                                  | Type          | Labels       |
|---------------------------------------------------|--------------------------------------------------------------|---------------|--------------|
| server_latest_iteration_time_seconds              | Latest invocation of the server                              | Gauge         |              |
| server_certificates_retrieved_total               | Total amount of certificates retrieved                       | Counter       |              |
| server_certificate_retrieve_errors_total          | Total errors while trying to retrieve certificates           | Counter       |              |
| server_certificates_renewals_total                | Total number of renewed certificates                         | Counter       |              |
| server_certificates_renewal_errors_total          | Total errors while trying to renew certificates              | Counter       |              |
| server_certificates_written_total                 | Total number of certificates written total                   | Counter (Vec) | subsystem    |
| server_certificates_write_errors_total            | Total errors while writing the certificate                   | Counter (Vec) | subsystem    |
| server_certificate_expiry_time                    | Timestamp of certificate expiry                              | Gauge (Vec)   | domain       |
| server_certificate_errors_total                   | Total number of errors while handling certificates           | Counter (Vec) | domain, desc |
| server_vault_aws_credentials_requested_total      | Total amount of dynamic AWS credentials requested            | Counter       |              |
| server_vault_aws_credentials_request_errors_total | Total errors while trying to acquire dynamic AWS credentials | Counter       |              |
