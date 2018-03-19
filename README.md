# micro
GoGo microservice lib

## Getting started
- Service Discovery
- NATS
- Health Check Scrpit

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
