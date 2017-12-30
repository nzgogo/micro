package gogo

import (
	"fmt"
	"micro/transport"
	"micro/registry"
	"strings"
	"github.com/satori/go.uuid"
)

type Service interface {
	Options() Options
	Init(...Options) error
	Start() error
	Stop() error
	Run() error
	Close()
}

type service struct {
	opts    Options
	name    string
	version string
	id      string
}

func (s *service) Options() Options {
	return s.opts
}

func (s *service) Init(opts ...Option) error {
	for _, o := range opts {
		o(&s.opts)
	}

	return nil
}

func (s *service) Register() error {
	config := s.Options()
	// register service
	node := &registry.Node{
		Id:  s.id,
	}

	service := &registry.Service{
		Name:      s.name,
		Version:   s.version,
		Nodes:     []*registry.Node{node},
	}

	if err := config.Registry.Register(service); err != nil {
		return err
	}

	return nil
}

func (s *service) Deregister() error {
	config := s.Options()

	node := &registry.Node{
		Id:  s.id,
	}

	service := &registry.Service{
		Name:      s.name,
		Version:   s.version,
		Nodes:     []*registry.Node{node},
	}

	fmt.Printf("Deregistering node: %s", node.Id)
	if err := config.Registry.Deregister(service); err != nil {
		return err
	}

	return nil
}

func (s *service) Start() error {

	if err := s.Register(); err != nil {
		return err
	}

	return nil
}

func (s *service) Stop() error {
	if err := s.Deregister(); err != nil {
		return err
	}

	return nil
}

func (s *service) Run() error {
	if err := s.Start(); err != nil {
		return err
	}

	if err := s.Stop(); err != nil {
		return err
	}

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

	r := registry.NewRegistry(registry.Addrs(*registryFlags["consul_addr"]))

	s.opts = newOptions(
		Transport(t),
		Registry(r),
	)
	return s
}
