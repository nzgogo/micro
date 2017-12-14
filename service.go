package gogo

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/transport"

	natsBroker "github.com/micro/go-plugins/broker/nats"
	natsTransport "github.com/micro/go-plugins/transport/nats"

	"github.com/nzgogo/micro/codec/json"
)

var natsAddr = "nats://dev.gogox.co.nz:4222"

// NewService returns a go-micro compatible service using nats as broker and transport
func NewService(opts ...micro.Option) micro.Service {
	bOptions := []broker.Option{
		broker.Addrs(natsAddr),
		broker.Codec(json.NewBrokerCodec()),
	}
	b := natsBroker.NewBroker(bOptions...)

	tOptions := []transport.Option{
		transport.Addrs(natsAddr),
		transport.Codec(json.NewTransportCodec()),
	}
	t := natsTransport.NewTransport(tOptions...)

	srvOptions := []micro.Option{
		micro.Broker(b),
		micro.Transport(t),
	}

	srvOptions = append(srvOptions, opts...)

	return micro.NewService(srvOptions...)
}
