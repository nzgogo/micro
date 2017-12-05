// +build ignore

package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"micro/registry"
)

func main() {
	reg := registry.NewRegistry()

	var srv registry.Service
	srv.ID = "01"
	srv.Name = "Nats"
	srv.Address = "127.0.0.1"
	srv.Port = 4222
	srv.Check = &api.AgentServiceCheck{
		Script:     "ls",
		Interval: "5s",
		Timeout:  "1s",
	}
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
