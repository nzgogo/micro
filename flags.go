package gogo

import (
	"flag"
	"fmt"
)

func parseFlags(s *service) {
	if s.config["nats_addr"] == "" {
		s.config["nats_addr"] = "nats://127.0.0.1:4222"
	}

	// s.config["consul_addr"] = *flag.String("consul", s.config["consul_addr"], "Consul server address")
	// s.config["nats_addr"] = *flag.String("nats", s.config["nats_addr"], "Nats server address")

	a := *flag.String("consul", s.config["consul_addr"], "Consul server address")
	b := *flag.String("nats", s.config["nats_addr"], "Nats server address")

	fmt.Println(a)
	fmt.Println(b)

	flag.Parse()
}
