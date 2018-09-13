# micro [![License Apache 2](https://img.shields.io/badge/License-Apache2-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0) [![Go Report Card](https://goreportcard.com/badge/github.com/nzgogo/micro)](https://goreportcard.com/report/github.com/nzgogo/micro) [![GoDoc](https://godoc.org/github.com/nzgogo/micro?status.svg)](http://godoc.org/github.com/nzgogo/micro)

GOGO micro is a set of utilities and libraries that makes it easier to write and manage microservices.

## Getting started
- Service Discovery
- NATS
- Health Check Scrpit
- Examples

## Service Discovery
Service discovery is used to resolve service names to NATS addresses (subject). We use [Consul](https://www.consul.io) as our service discovery system.
[Install Consul](https://www.consul.io/intro/getting-started/install.html)

## NATS
[NATS](https://nats.io) is a messaging system used as internal communication of our distributed services.

### [Installation](https://github.com/nats-io/go-nats)

```bash
# Go client
go get github.com/nats-io/go-nats

# Server
go get github.com/nats-io/gnatsd
```

## Health Check Script
One of the primary roles of the Consul agent is management of system-level and application-level health checks. There are several different kinds of checks, see [Checks Definition](https://www.consul.io/docs/agent/checks.html).
The checks used in micro is **Script + Interval**. 

## Examples
Refer to nzgogo/examples directory for more details, which contains examples such as how to write a service using micro.
