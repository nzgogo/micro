package registry

import (
	"errors"

	consul "github.com/hashicorp/consul/api"
)

// Registry service
type Registry struct {
	Client *consul.Client
	Config *consul.Config
}

// Service struct
type Service struct {
	ID      string
	Name    string
	Tags    []string
	Port    int
	Address string
	check   *consul.AgentServiceCheck
}

// NewRegistry function
func NewRegistry() *Registry {
	config := consul.DefaultConfig()

	// create the client
	client, _ := consul.NewClient(config)

	cr := &Registry{
		Client: client,
		Config: config,
	}

	return cr
}

// Deregister a service
func (r *Registry) Deregister(s *Service) error {
	if len(s.ID) == 0 {
		return errors.New("Service ID is required")
	}

	return r.Client.Agent().ServiceDeregister(s.ID)
}

// Register a service
func (r *Registry) Register(s *Service) error {
	// register the service
	err := r.Client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		ID:      s.ID,
		Name:    s.Name,
		Tags:    s.Tags,
		Port:    s.Port,
		Address: s.Address,
		Check:   s.check,
	})
	if err != nil {
		return err
	}

	return nil
}
