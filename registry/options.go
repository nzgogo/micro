package registry

import (
	"crypto/tls"
	"time"
	consul "github.com/hashicorp/consul/api"
)

type Options struct {
	//registry options
	Addrs     []string
	Timeout   time.Duration
	Secure    bool
	TLSConfig *tls.Config

	//consul agent check options
	//CheckArgs     []string
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

// Secure communication with the registry
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

// Specify TLS Config
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

//specify consul Agent check args
func Checks(checks ...*consul.AgentServiceCheck) Option {
	return func(o *Options) {
		for _, c := range checks{
			o.Checks = append(o.Checks, c)
		}
	}
}
