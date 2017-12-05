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


	var srv1 registry.Service
	srv1.ID = "02"
	srv1.Name = "Nats"
	srv1.Address = "127.0.0.1"
	srv1.Port = 4223
	srv1.Check = &api.AgentServiceCheck{
		Script:     "fuck up",
		Interval: "5s",
		Timeout:  "1s",
	}
	err1 := reg.Register(&srv1)

	if err1 != nil {
		fmt.Printf("[Error]%s\n", err1)
	}
}
