package main

import (
	"fmt"
	"micro/registry"
)

func main() {
	reg := registry.NewRegistry()

	var srv registry.Service
	srv.ID = "01"
	err := reg.Deregister(&srv)

	if err != nil {
		fmt.Printf("[Error]%s \n", err)
	}
}
