package gogo

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"log"

	"github.com/nzgogo/micro/codec"
	"github.com/nzgogo/micro/registry"
	"github.com/nzgogo/micro/router"
	"github.com/nzgogo/micro/selector"
	"github.com/nzgogo/micro/transport"
	"github.com/satori/go.uuid"
)

const (
	ORGANIZATION = "gogo"

	// Service configs
	CONFIG_NATS_ADDRESS                         = "nats_addr"
	CONFIG_CONSUL_ADDRRESS                      = "consul_addr"
	CONFIG_HC_SCRIPT                            = "hc_script"
	CONFIG_HC_INTERVAL                          = "hc_interval"
	CONFIG_HC_DEREGISTER_CRITICAL_SERVICE_AFTER = "hc_deregister_critical_service_after"
	CONFIG_HC_LOAD_CRITICAL_THRESHOLD           = "hc_load_critical_threshold"
	CONFIG_HC_LOAD_WARNING_THRESHOLD            = "hc_load_warning_threshold"
	CONFIG_HC_MEMORY_CRITICAL_THRESHOLD         = "hc_memory_critical_threshold"
	CONFIG_HC_MEMORY_WARNING_THRESHOLD          = "hc_memory_warning_threshold"
	CONFIG_HC_CPU_CRITICAL_THRESHOLD            = "hc_cpu_critical_threshold"
	CONFIG_HC_CPU_WARNING_THRESHOLD             = "hc_cpu_warning_threshold"

	// Default value for health checks
	DEFAULT_HC_SCRITP                           = "gghc"
	DEFAULT_HC_INTERVAL                         = "1m"
	DEFALT_HC_DEREGISTER_CRITICAL_SERVICE_AFTER = "5m"
	DEFALT_HC_LOAD_CRITICAL_THRESHOLD           = "0.9"
	DEFALT_HC_LOAD_WARNING_THRESHOLD            = "0.8"
	DEFALT_HC_MEMORY_CRITICAL_THRESHOLD         = "5"
	DEFALT_HC_MEMORY_WARNING_THRESHOLD          = "15"
	DEFALT_HC_CPU_CRITICAL_THRESHOLD            = "5"
	DEFALT_HC_CPU_WARNING_THRESHOLD             = "15"
)

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

func (s *service) Name() string {
	return s.name
}

func (s *service) Version() string {
	return s.version
}

func (s *service) ID() string {
	return s.id
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

	log.Printf("[Service][Name] %s\n", s.name)
	log.Printf("[Service][Version] %s\n", s.version)
	log.Printf("[Service][ID] %s\n", s.id)

	s.config = readConfigFile(strings.Replace(s.name, "-", ".", -1) + "." + s.version)

	parseFlags(s)
	trans := transport.NewTransport(
		transport.Subject(strings.Replace(s.name, "-", ".", -1)+"."+s.version+"."+s.id),
		transport.Addrs(s.config[CONFIG_NATS_ADDRESS]),
	)
	var reg registry.Registry
	check := packHealthCheck(s.config, trans.Options().Subject)
	if check == nil {
		log.Println("NO HEALTH CHECK REGISTERED !!!")
		reg = registry.NewRegistry(
			registry.Addrs(s.config[CONFIG_CONSUL_ADDRRESS]),
		)
	} else {
		reg = registry.NewRegistry(
			registry.Addrs(s.config[CONFIG_CONSUL_ADDRRESS]),
			registry.Checks(check),
		)
	}
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
