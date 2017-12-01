package registry

import (
	consul "github.com/hashicorp/consul/api"
)

type Registry struct {
	Client *consul.Client
	Config *consul.Config
}
