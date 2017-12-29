package gogo

import (
	"fmt"
	"micro/transport"
	"strings"

	"github.com/satori/go.uuid"
)

type Service interface {
	Options() Options
	Init(...Options) error
	Run() error
	Close()
}

type service struct {
	opts    Options
	name    string
	version string
	id      string

}

type Option func(*Options)

func (s *service) Options() Options {
	return s.opts
}

func (s *service) Init(opts ...Option) error {
	for _, o := range opts {
		o(&s.opts)
	}
	return nil
}

func (s *service) Run() error {
	return nil
}

func (s *service) Close() {
}

func NewService(n string, v string) *service {
	id := strings.Replace(uuid.NewV4().String(), "-", "", -1)

	s := &service{
		name:    n,
		version: v,
		id:      id,
	}

	fmt.Printf("[Service][Name] %s\n", s.name)
	fmt.Printf("[Service][Version] %s\n", s.version)
	fmt.Printf("[Service][ID] %s\n", s.id)

	parseFlags()

	t := transport.NewTransport(
		transport.Subject(s.name+"."+s.version+"."+s.id),
		transport.Addrs(*transportFlags["nats_addr"]),
	)

	s.opts = newOptions(
		Transport(t),
	)
	return s
}
