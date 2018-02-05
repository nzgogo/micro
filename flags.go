package gogo

import (
	"flag"
)

func parseFlags(s *service) {
	if s.config["nats_addr"] == "" {
		s.config["nats_addr"] = "nats://127.0.0.1:4222"
	}

	consul_addr := flag.String("consul", s.config["consul_addr"], "Consul server address")
	nats_addr := flag.String("nats", s.config["nats_addr"], "Nats server address")

	flag.Parse()

	s.config["consul_addr"] = *consul_addr
	s.config["nats_addr"] = *nats_addr
}
