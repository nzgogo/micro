package gogo

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/nzgogo/micro/codec"
	"github.com/nzgogo/micro/registry"
	"github.com/nzgogo/micro/router"
	"github.com/nzgogo/micro/selector"
	"github.com/nzgogo/micro/transport"
	consul "github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
)

//var checkLoads = &consul.AgentServiceCheck{
//	Notes: "check-ping: warning limit 10ms RTA or 2% packet-loss, critical limit is 20ms or 5% packet loss",
//	Script: "/usr/local/Cellar/nagios-plugins/2.2.1/libexec/sbin/check_ping -4 -H 192.168.1.1 -w 10,2% -c 20,5%",
//	Interval: "1m",
//}

type Service interface {
	Options() Options
	Init(...Option) error
	Run() error
	Respond(message *codec.Message, subject string) error
}

type service struct {
	opts    Options
	config  map[string]string
	name    string
	version string
	id      string
}

func (s *service) Options() Options {
	return s.opts
}

func (s *service) Config() map[string]string {
	return s.config
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

	if err := s.opts.Selector.Init(); err != nil {
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
	}

	if err := s.stop(); err != nil {
		return err
	}

	return nil
}

func (s *service) Respond(message *codec.Message, subject string) error {
	s.opts.Context.Delete(message.ContextID)
	resp, err := codec.Marshal(message)
	if err != nil {
		return err
	}
	return s.opts.Transport.Publish(subject, resp)
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

	s.config = readConfigFile()

	parseFlags(s)

	trans := transport.NewTransport(
		transport.Subject(strings.Replace(s.name, "-", ".", -1)+"."+s.version+"."+s.id),
		transport.Addrs(s.config["nats_addr"]),
	)

	var check = &consul.AgentServiceCheck{
		//Notes: "health check",
		Args: []string{s.config["health_check_script"]," -subj="+trans.Options().Subject},
		Interval: "1m",
	}

	reg := registry.NewRegistry(
		registry.Addrs(s.config["consul_addr"]),
		registry.Checks(check),
	)

	sel := selector.NewSelector(
		selector.Registry(reg),
		selector.SetStrategy(selector.RoundRobin),
	)

	s.opts = newOptions(
		Transport(trans),
		Registry(reg),
		Selector(sel),
	)
	return s
}
