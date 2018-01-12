package gogo

import "flag"

var registryFlags map[string]*string
var transportFlags map[string]*string

func parseFlags() {
	registryFlags = make(map[string]*string)
	transportFlags = make(map[string]*string)

	registryFlags["consul_addr"] = flag.String("consul", "", "Consul server address")
	transportFlags["nats_addr"] = flag.String("nats", "nats://127.0.0.1:4222", "Nats server address")

	flag.Parse()
}
