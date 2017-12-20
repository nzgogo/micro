package registry

import (
	"errors"
	"net"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"sync"
	hash "github.com/mitchellh/hashstructure"
)

type Registry interface {
	Init() error
	Register(*Service) error
	Deregister(*Service) error
	GetService(string) ([]*Service, error)
	ListServices() ([]*Service, error)
	Options() Options
	//Watch() (Watcher, error)
}

type Option func(*Options)

// Registry service
type registry struct {
	Client *consul.Client
	opts Options

	sync.Mutex
	register map[string]uint64
}

// Service struct
type Service struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Nodes     []*Node           `json:"nodes"`
}

type Node struct {
	Id       string            `json:"id"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
}

// NewRegistry function
func NewRegistry(opts ...Option) *registry {
	var options Options

	for _, o := range opts {
		o(&options)
	}

	return &registry{
		opts:options,
	}
}

func (r *registry) Init() error{
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
	if err !=nil{
		return err
	}
	// set timeout
	if r.opts.Timeout > 0 {
		config.HttpClient.Timeout = r.opts.Timeout
	}

	r.Client = client
	r.register = make(map[string]uint64)

	return nil
}

// Register a service
func (r *registry) Register(s *Service) error {
	if len(s.Nodes) == 0 {
		return errors.New("Require at least one node")
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
	v, ok := r.register[s.Name]
	r.Unlock()

	// if it's already registered and matches then just pass the check
	if ok && v == h {
		//// if the err is nil we're all good, bail out
		//// if not, we don't know what the state is, so full re-register
		//if err := r.Client.Agent().PassTTL("service:"+node.Id, ""); err == nil {
			return nil
		//}
	}

	// encode the tags
	tags := []string{s.Version}

	var check *consul.AgentServiceCheck


	if len(r.opts.CheckArgs) > 0{
		check = &consul.AgentServiceCheck{
			Args:r.opts.CheckArgs,
		}
	}

	// register the service
	if err := r.Client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		ID:      node.Id,
		Name:    s.Name,
		Tags:    tags,
		Port:    node.Port,
		Address: node.Address,
		Check:   check,
	}); err != nil {
		return err
	}

	// save our hash of the service
	r.Lock()
	r.register[s.Name] = h
	r.Unlock()

	return nil
}

// Deregister a service
func (r *registry) Deregister(s *Service) error {
	if len(s.Nodes) == 0 {
		return errors.New("Service ID is required")
	}

	// delete our hash of the service
	r.Lock()
	delete(r.register, s.Name)
	r.Unlock()

	node := s.Nodes[0]
	return r.Client.Agent().ServiceDeregister(node.Id)
}

func (r *registry) GetService(name string) ([]*Service, error) {
	rsp, _, err := r.Client.Health().Service(name, "", false, nil)
	if err != nil {
		return nil, err
	}

	serviceMap := map[string]*Service{}

	for _, s := range rsp {
		if s.Service.Service != name {
			continue
		}

		if len(s.Service.Tags) <=0 {
			continue
		}
		// version is now a tag
		version := s.Service.Tags[0]
		// service ID is now the node id
		id := s.Service.ID
		// key is always the version
		key := version
		// address is service address
		address := s.Service.Address

		svc, ok := serviceMap[key]
		if !ok {
			svc = &Service{
				Name:      s.Service.Service,
				Version:   version,
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
			Id:       id,
			Address:  address,
			Port:     s.Service.Port,
		})
	}

	var services []*Service
	for _, service := range serviceMap {
		services = append(services, service)
	}
	return services, nil
}

func (r *registry) ListServices() ([]*Service, error) {
	rsp, _, err := r.Client.Catalog().Services(nil)
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
