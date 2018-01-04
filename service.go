package gogo

import (
	"fmt"
	"micro/transport"
	"micro/registry"
	"strings"
	"github.com/satori/go.uuid"
	"os/signal"
	"syscall"
	"os"
	"github.com/nats-io/go-nats"
	"micro/codec"
	"micro/router"
)

type Service interface {
	Options() Options
	Init(...Option) error
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
	if err:= s.opts.Transport.Init(); err != nil{
		return err
	}

	if err:= s.opts.Registry.Init(); err != nil{
		return err
	}
	if s.opts.Router == nil {
		return nil
	}
	if err:= s.opts.Router.Init(router.Client(s.opts.Registry.Client())); err != nil{
		return err
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

	if config.Router == nil {
		return nil
	}
	if err := config.Router.Register(config.Router.String()); err!=nil {
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

	if config.Router == nil {
		return nil
	}
	//delete all service kv
	if err := config.Router.Deregister(config.Router.String()); err!=nil {
		return err
	}

	return nil
}

func (s *service) Start() error {
	if err := s.Register(); err != nil {
		return err
	}

	tc := s.Options().Transport
	if err := tc.Subscribe(func(msg *nats.Msg){
		req := &codec.Request{}
		s.opts.Codec.Unmarshal(msg.Data, req)
		handler, err1 := s.opts.Router.Dispatch(req)
		if err1 != nil || handler == nil{
			resp, _ := s.opts.Codec.Marshal(codec.Response{
				404,
				make(map[string][]string,0),
				"Page not found",
			})
			tc.Publish(msg.Reply, resp)
		}
		err2 := handler(req, tc, msg.Reply)
		if err2 != nil {
			resp, _ := s.opts.Codec.Marshal(codec.Response{
				500,
				make(map[string][]string,0),
				"Internal Server Error",
			})
			tc.Publish(msg.Reply, resp)
		}

	}); err != nil{
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

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	select {
	// wait on kill signal
	case <-ch:
		// wait on context cancel
	//case <-s.opts.Context.Done():
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
