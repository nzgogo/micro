package gogo

import (
	"flag"

	"github.com/nats-io/go-nats"
)

func parseFlags(s *service) {
	if s.config[CONFIG_NATS_ADDRESS] ==  ""{
		s.config[CONFIG_NATS_ADDRESS] = nats.DefaultURL
	}

	consul_addr := flag.String("consul", s.config[CONFIG_CONSUL_ADDRRESS], "Consul server address")
	nats_addr := flag.String("nats", s.config[CONFIG_NATS_ADDRESS], "Nats server address")

	flag.Parse()

	s.config[CONFIG_CONSUL_ADDRRESS] = *consul_addr
	s.config[CONFIG_NATS_ADDRESS] = *nats_addr
}
