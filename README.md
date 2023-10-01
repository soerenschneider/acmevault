# acmevault
[![Go Report Card](https://goreportcard.com/badge/github.com/soerenschneider/acmevault)](https://goreportcard.com/report/github.com/soerenschneider/acmevault)
![test-workflow](https://github.com/soerenschneider/acmevault/actions/workflows/test.yaml/badge.svg)
![release-workflow](https://github.com/soerenschneider/acmevault/actions/workflows/release-container.yaml/badge.svg)
![golangci-lint-workflow](https://github.com/soerenschneider/acmevault/actions/workflows/golangci-lint.yaml/badge.svg)


## Features

🔐 Issues certificates from any [ACME provider](https://datatracker.ietf.org/doc/html/rfc8555), such as Let's Encrypt<br/>
⏰ Automatically renews certificates before they expire<br/>
🔌 Stores all data inside Vault and thus decouples from clients<br/>

## Why would I need this?

## Problem Statement

Rolling out TLS encryption shouldn't need to be pitched anymore (even for internal services). Using the DNS01 ACME challenge is proven and allows issuing certs non-public routable machines. On the other hand, you need to have access to either highly-privileged/narrowly-scoped credentials of your DNS provider to solve these DNS01 challenges.

In the case of Route53, if you don't want to end up creating dozens of hosted zones, one for each of your subdomains, you're at risk of leaking highly-privileged IAM credentials.

Acmevault requests short-lived IAM credentials for Route53 and uses them to perform DNS01 challenges for the configured domains and writes the issued X509 certificates to Hashicorp Vault's K/V secret store - only readable by the appropriate AppRole.

Its client mode reads the respective written certificates from Vault and installs them to a preconfigured location, optionally invoking post-installation hooks.


## Overview
![Overview](overview.png)

## Installation

### Docker / Podman
```
$ git clone https://github.com/soerenschneider/acmevault
$ cd acmevault
$ docker run -v $(pwd)/contrib:/config ghcr.io/soerenschneider/acmevault -conf /config/server.json
```
### Binaries
Download a prebuilt binary from the [releases section](https://github.com/soerenschneider/acmevault/releases) for your system.

### From Source
As a prerequisite, you need to have [Golang SDK](https://go.dev/dl/) installed. Then you can install acmevault from source by invoking:
```shell
$ go install github.com/soerenschneider/acmevault@latest
```

## Configuration
See the [configuration section](docs/configuration.md) for examples and configuration reference.

## Observability
See the [metrics section](docs/configuration.md) for an overview of exposed metrics.

## Changelog
See the full changelog [here](CHANGELOG.md)
