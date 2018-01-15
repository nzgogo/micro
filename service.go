package gogo

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/nzgogo/micro/registry"
	"github.com/nzgogo/micro/router"
	"github.com/nzgogo/micro/transport"
	"github.com/satori/go.uuid"
)

type Service interface {
	Options() Options
	Init(...Option) error
	Run() error
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
	if err := s.opts.Transport.Init(); err != nil {
		return err
	}

	if err := s.opts.Registry.Init(); err != nil {
		return err
	}

	router := router.NewRouter(
		router.Name(strings.Replace(s.name, "-", "/", -1)+"/"+s.version),
		router.Client(s.opts.Registry.Client()),
	)
	s.opts.Router = router

	return nil
}

func (s *service) Run() error {
	if err := s.start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	select {
	// wait on kill signal
	case <-ch:
		// wait on context cancel
		//case <-s.opts.Context.Done():
	}

	if err := s.stop(); err != nil {
		return err
	}

	return nil
}

func (s *service) start() error {
	if err := s.register(); err != nil {
		return err
	}
	tc := s.Options().Transport

	if err := tc.Subscribe(); err != nil {
		return err
	}

	return nil
}

func (s *service) stop() error {
	if err := s.deregister(); err != nil {
		return err
	}

	return nil
}

func (s *service) register() error {
	config := s.Options()
	// register service
	node := &registry.Node{
		ID: s.id,
	}

	service := &registry.Service{
		Name:    s.name,
		Version: s.version,
		Nodes:   []*registry.Node{node},
	}

	if config.Router != nil {
		if err := config.Router.Register(); err != nil {
			return err
		}
	}

	if err := config.Registry.Register(service); err != nil {
		return err
	}

	return nil
}

func (s *service) deregister() error {
	config := s.Options()

	node := &registry.Node{
		ID: s.id,
	}

	service := &registry.Service{
		Name:    s.name,
		Version: s.version,
		Nodes:   []*registry.Node{node},
	}

	fmt.Printf("Deregistering node: %s", node.ID)
	if err := config.Registry.Deregister(service); err != nil {
		return err
	}

	if config.Router == nil {
		return nil
	}
	//delete all service kv
	if err := config.Router.Deregister(); err != nil {
		return err
	}

	return nil
}

func NewService(n string, v string) *service {
	newUUID, _ := uuid.NewV4()
	id := strings.Replace(newUUID.String(), "-", "", -1)

	s := &service{
		name:    n,
		version: v,
		id:      id,
	}

	fmt.Printf("[Service][Name] %s\n", s.name)
	fmt.Printf("[Service][Version] %s\n", s.version)
	fmt.Printf("[Service][ID] %s\n", s.id)

	parseFlags()

	trans := transport.NewTransport(
		transport.Subject(strings.Replace(s.name, "-", ".", -1)+"."+s.version+"."+s.id),
		transport.Addrs(*transportFlags["nats_addr"]),
	)

	reg := registry.NewRegistry(
		registry.Addrs(*registryFlags["consul_addr"]),
	)

	s.opts = newOptions(
		Transport(trans),
		Registry(reg),
	)
	return s
}
