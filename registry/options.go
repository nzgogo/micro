package registry

import (
	"time"

	consul "github.com/hashicorp/consul/api"
)

type Options struct {
	//registry options
	Addrs   []string
	Timeout time.Duration

	//consul agent check options
	Checks consul.AgentServiceChecks
}

// Addrs is the registry addresses to use
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

//specify consul Agent check args
func Checks(checks ...*consul.AgentServiceCheck) Option {
	return func(o *Options) {
		for _, c := range checks {
			o.Checks = append(o.Checks, c)
		}
	}
}
