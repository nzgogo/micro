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
The checks we use is **Script + Interval**.  The health check script is on Gitlab repo [gogoexpress / gogoexpress-go-consul-healthcheck-v1](https://gitlab.com/gogoexpress/gogoexpress-go-consul-healthcheck-v1.git).
To get this script program working, a configuration file in json format is required:
```shell
$ vi /etc/gogo/config-healthcheck.json
```
```json
{
  "nats_addr": "",
  "https://hooks.slack.com/services/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"
}
```
