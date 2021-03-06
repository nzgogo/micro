package registry

import (
	"fmt"
	"net"
	"sync"

	consul "github.com/hashicorp/consul/api"
	hash "github.com/mitchellh/hashstructure"
	"github.com/nzgogo/micro/constant"
)

type Registry interface {
	Init() error
	Register(*Service) error
	Deregister(*Service) error
	GetService(string) ([]*Service, error)
	ListServices() ([]*Service, error)
	Options() Options
	Client() *consul.Client
}

type Option func(*Options)

// Registry service
type registry struct {
	conn *consul.Client
	opts Options
	sync.Mutex
	register map[string]uint64
}

// Service struct
type Service struct {
	Name    string  `json:"name"`
	Version string  `json:"version"`
	Nodes   []*Node `json:"nodes"`
}

type Node struct {
	ID string `json:"id"`
}

var (
	DefaultRegistry = NewRegistry()
)

func (r *registry) Init() error {
	config := consul.DefaultConfig()

	// check if there are any addrs
	if len(r.opts.Addrs) > 0 {
		addr, port, err := net.SplitHostPort(r.opts.Addrs[0])
		if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
			port = "8500"
			addr = r.opts.Addrs[0]
			config.Address = fmt.Sprintf("%s:%s", addr, port)
		} else if err == nil {
			config.Address = fmt.Sprintf("%s:%s", addr, port)
		}
	}

	// create the client
	client, err := consul.NewClient(config)
	if err != nil {
		return err
	}
	// set timeout
	if r.opts.Timeout > 0 {
		config.HttpClient.Timeout = r.opts.Timeout
	}

	r.conn = client
	r.register = make(map[string]uint64)

	return nil
}

// Register a service
func (r *registry) Register(s *Service) error {
	if len(s.Nodes) == 0 {
		return constant.ErrRegistryEmptyNode
	}

	// create hash of service; uint64
	h, err := hash.Hash(s, nil)
	if err != nil {
		return err
	}

	// use first node
	node := s.Nodes[0]

	// get existing hash
	r.Lock()
	v, ok := r.register[s.Nodes[0].ID]
	r.Unlock()

	if ok && v == h {
		return nil
	}

	// encode the tags
	tags := []string{s.Version}

	// register the service
	if err := r.conn.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		ID:   node.ID,
		Name: s.Name,
		Tags: tags,
		//Port:    node.Port,
		//Address: node.Address,
		Checks: r.opts.Checks,
	}); err != nil {
		return err
	}

	// save our hash of the service
	r.Lock()
	r.register[s.Nodes[0].ID] = h
	r.Unlock()

	return nil
}

// Deregister a service
func (r *registry) Deregister(s *Service) error {
	if len(s.Nodes) == 0 {
		return constant.ErrRegistryEmptyNode
	}

	// delete our hash of the service
	r.Lock()
	delete(r.register, s.Name)
	r.Unlock()

	node := s.Nodes[0]
	return r.conn.Agent().ServiceDeregister(node.ID)
}

func (r *registry) GetService(name string) ([]*Service, error) {
	// todo make passingOnly configurable
	rsp, _, err := r.conn.Health().Service(name, "", false, nil)
	if err != nil {
		return nil, err
	}

	serviceMap := map[string]*Service{}

	for _, s := range rsp {
		if s.Service.Service != name {
			continue
		}

		if len(s.Service.Tags) <= 0 {
			continue
		}
		version := s.Service.Tags[0]
		id := s.Service.ID
		// key is always the version
		key := version

		svc, ok := serviceMap[key]
		if !ok {
			svc = &Service{
				Name:    s.Service.Service,
				Version: version,
			}
			serviceMap[key] = svc
		}

		var del bool
		for _, check := range s.Checks {
			// delete the node if the status is critical
			if check.Status == "critical" {
				del = true
				break
			}
		}

		// if delete then skip the node
		if del {
			continue
		}

		svc.Nodes = append(svc.Nodes, &Node{
			ID: id,
		})
	}

	var services []*Service
	for _, service := range serviceMap {
		services = append(services, service)
	}
	return services, nil
}

func (r *registry) ListServices() ([]*Service, error) {
	rsp, _, err := r.conn.Catalog().Services(nil)
	if err != nil {
		return nil, err
	}

	var services []*Service

	for service := range rsp {
		services = append(services, &Service{Name: service})
	}

	return services, nil
}

func (r *registry) Options() Options {
	return r.opts
}

func (r *registry) Client() *consul.Client {
	return r.conn
}

// NewRegistry function
func NewRegistry(opts ...Option) *registry {
	var options Options

	for _, o := range opts {
		o(&options)
	}

	return &registry{
		opts: options,
	}
}
