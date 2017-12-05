package main

import (
	"micro/discovery"
	"fmt"
)
func main() {
	disc := discovery.NewDiscovery()

	disc.DiscoverServices(true,"Nats")

	fmt.Println("found service: \n", disc.ServList)

}