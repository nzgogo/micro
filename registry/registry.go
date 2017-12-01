package registry

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	consul "github.com/hashicorp/consul/api"
	hash "github.com/mitchellh/hashstructure"
)

// Registry service
type Registry struct {
	Client  *consul.Client
	Config  *consul.Config
	Address string
	Options Options

	sync.Mutex
	register map[string]uint64
}

// Option passed
type Option func(*Options)

// RegisterOption passed
type RegisterOption func(*RegisterOptions)

func newRegistry(opts ...Option) Registry {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	config := consul.DefaultConfig()

	if len(options.Addrs) > 0 {
		addr, port, err := net.SplitHostPort(options.Addrs[0])
		if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
			port = "8500"
			addr = options.Addrs[0]
			config.Address = fmt.Sprintf("%s:%s", addr, port)
		} else if err == nil {
			config.Address = fmt.Sprintf("%s:%s", addr, port)
		}
	}

	// create the client
	client, _ := consul.NewClient(config)

	// set timeout
	if options.Timeout > 0 {
		config.HttpClient.Timeout = options.Timeout
	}

	cr := &Registry{
		Address:  config.Address,
		Client:   client,
		Options:  options,
		register: make(map[string]uint64),
	}

	return *cr
}

// Deregister a service
func (c *Registry) Deregister(s *Service) error {
	if len(s.Nodes) == 0 {
		return errors.New("Require at least one node")
	}

	// delete our hash of the service
	c.Lock()
	delete(c.register, s.Name)
	c.Unlock()

	node := s.Nodes[0]
	return c.Client.Agent().ServiceDeregister(node.Id)
}

// Register a service
func (c *Registry) Register(s *Service, opts ...RegisterOption) error {
	if len(s.Nodes) == 0 {
		return errors.New("Require at least one node")
	}

	var options RegisterOptions
	for _, o := range opts {
		o(&options)
	}

	// create hash of service; unit64
	h, err := hash.Hash(s, nil)
	if err != nil {
		return err
	}

	// use first node
	node := s.Nodes[0]

	// get existing hash
	c.Lock()
	v, ok := c.register[s.Name]
	c.Unlock()

	if ok && v == h {
		if err := c.Client.Agent().PassTTL("service:"+node.Id, ""); err == nil {
			return nil
		}
	}

	// encode the tags
	tags := encodeMetadata(node.Metadata)
	tags = append(tags, encodeEndpoints(s.Endpoints)...)
	tags = append(tags, encodeVersion(s.Version)...)

	var check *consul.AgentServiceCheck

	// if the TTL is greater than 0 create an associated check
	if options.TTL > time.Duration(0) {
		// splay slightly for the watcher?
		splay := time.Second * 5
		deregTTL := options.TTL + splay
		// consul has a minimum timeout on deregistration of 1 minute.
		if options.TTL < time.Minute {
			deregTTL = time.Minute + splay
		}

		check = &consul.AgentServiceCheck{
			TTL: fmt.Sprintf("%v", options.TTL),
			DeregisterCriticalServiceAfter: fmt.Sprintf("%v", deregTTL),
		}
	}

	// register the service
	if err := c.Client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
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
	c.Lock()
	c.register[s.Name] = h
	c.Unlock()

	// if the TTL is 0 we don't mess with the checks
	if options.TTL == time.Duration(0) {
		return nil
	}

	// pass the healthcheck
	return c.Client.Agent().PassTTL("service:"+node.Id, "")
}

// GetService by name
func (c *Registry) GetService(name string) ([]*Service, error) {
	rsp, _, err := c.Client.Health().Service(name, "", false, nil)
	if err != nil {
		return nil, err
	}

	serviceMap := map[string]*Service{}

	for _, s := range rsp {
		if s.Service.Service != name {
			continue
		}

		// version is now a tag
		version, found := decodeVersion(s.Service.Tags)
		// service ID is now the node id
		id := s.Service.ID
		// key is always the version
		key := version
		// address is service address
		address := s.Service.Address

		// if we can't get the version we bail
		// use old the old ways
		if !found {
			continue
		}

		svc, ok := serviceMap[key]
		if !ok {
			svc = &Service{
				Endpoints: decodeEndpoints(s.Service.Tags),
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
			Metadata: decodeMetadata(s.Service.Tags),
		})
	}

	var services []*Service
	for _, service := range serviceMap {
		services = append(services, service)
	}
	return services, nil
}

// ListServices show service list
func (c *Registry) ListServices() ([]*Service, error) {
	rsp, _, err := c.Client.Catalog().Services(nil)
	if err != nil {
		return nil, err
	}

	var services []*Service

	for service := range rsp {
		services = append(services, &Service{Name: service})
	}

	return services, nil
}

// // Watch call
// func (c *Registry) Watch() (Watcher, error) {
// 	return newConsulWatcher(c)
// }

// func (c *Registry) String() string {
// 	return "consul"
// }
