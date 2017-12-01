package registry

import (
	consul "github.com/hashicorp/consul/api"
)

// Registry service
type Registry struct {
	Client *consul.Client
	Config *consul.Config
}

type Service struct {
	ID      string
	Name    string
	Tags    []string
	Port    int
	Address string
	check   *consul.AgentServiceCheck
}

func newRegistry() *Registry {
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
// func (c *Registry) Deregister(s *Service) error {
// 	if len(s.Nodes) == 0 {
// 		return errors.New("Require at least one node")
// 	}
//
// 	// delete our hash of the service
// 	c.Lock()
// 	delete(c.register, s.Name)
// 	c.Unlock()
//
// 	node := s.Nodes[0]
// 	return c.Client.Agent().ServiceDeregister(node.Id)
// }

// Register a service
func (r *Registry) Register(s *Service) error {
	// register the service
	if err := r.Client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		ID:      s.ID,
		Name:    s.Name,
		Tags:    s.Tags,
		Port:    s.Port,
		Address: s.Address,
		Check:   s.check,
	}); err != nil {
		return err
	}

	return nil
}
