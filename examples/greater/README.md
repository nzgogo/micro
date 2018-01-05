# Greeter

An example Greeter application

## Contents

- **srv** - a greeter service
- **api** - examples of RESTful API

## Deps

- Service discovery is required for all services. Default is Consul.
- Nats-io is defaultly used as all Service internal transport.

### Consul

```
brew install consul
consul agent -dev
```

### Nats
```
go get github.com/nats
gnatsd
```

## Run Service

Start micro.srv.greeter
```shell
go run srv/server.go
```

## API

Micro logically separates API services from backend services. By default the micro API
accepts HTTP requests and converts to *codec.Request and *codec.Response types. Find them here [micro/api/proto](https://github.com/nzgogo/micro/codec).

Run the micro.api.greeter API Service
```shell
go run api/api.go 
```

## Call greeter service
```shell
curl http://localhost:8080/gogox/v1/greeter/hello
```

Browse to http://localhost:8080/gogox/v1/greeter/hello
