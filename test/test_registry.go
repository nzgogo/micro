// +build ignore

package main

import (
	"fmt"

	consul "github.com/hashicorp/consul/api"
	"github.com/nzgogo/micro/registry"
)

func main() {
	reg := registry.NewRegistry()

	srv := registry.Service{
		ID:      "01",
		Name:    "Nats",
		Address: "127.0.0.1",
		Port:    4222,
	}

	healthCheck := consul.AgentServiceCheck{
		Args:     []string{"/usr/local/bin/check", "-s", srv.Name + "." + srv.ID},
		Interval: "30s",
		Timeout:  "3s",
	}

	srv.Check = &healthCheck

	err := reg.Register(&srv)

	if err != nil {
		fmt.Printf("[Error]%s\n", err)
	}
}
