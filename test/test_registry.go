package main

import (
	"fmt"

	"github.com/nzgogo/micro/registry"
)

func main() {
	reg := registry.NewRegistry()

	var srv registry.Service
	srv.ID = "01"
	srv.Name = "Nats"
	srv.Address = "127.0.0.1"
	srv.Port = 4222

	err := reg.Register(&srv)

	if err != nil {
		fmt.Printf("[Error]%s\n", err)
	}
}
