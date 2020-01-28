package gogo

import (
	"flag"

	"github.com/nzgogo/go-nats"
	"github.com/nzgogo/micro/constant"
)

func parseFlags(s *service) {
	if s.config[constant.CONFIG_NATS_ADDRESS] == "" {
		s.config[constant.CONFIG_NATS_ADDRESS] = nats.DefaultURL
	}

	consulAddr := flag.String("consul", s.config[constant.CONFIG_CONSUL_ADDRRESS], "Consul server address")
	natsAddr := flag.String("nats", s.config[constant.CONFIG_NATS_ADDRESS], "Nats server address")

	flag.Parse()

	s.config[constant.CONFIG_CONSUL_ADDRRESS] = *consulAddr
	s.config[constant.CONFIG_NATS_ADDRESS] = *natsAddr
}
