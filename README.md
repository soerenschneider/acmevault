# acmevault

## Problem Statement

Rolling out TLS encryption shouldn't need to be pitched anymore (even for internal services). Using the DNS01 ACME challenge is proven and allows issuing certs non-public routable machines. On the other hand, you need to have access to either highly-privileged/narrowly-scoped credentials of your DNS provider to solve these DNS01 challenges.

In the case of Route53, if you don't want to end up creating dozens of hosted zones, one for each of your subdomains, you're at risk of leaking highly-privileged IAM credentials.

Acmevault requests short-lived IAM credentials for Route53 and uses them to perform DNS01 challenges for the configured domains and writes the issued X509 certificates to Hashicorp Vault's K/V secret store - only readable by the appropriate AppRole.

Its client mode reads the respective written certificates from Vault and installs them to a preconfigured location, optionally invoking post-installation hooks.


## Overview
![Overview](overview.png)