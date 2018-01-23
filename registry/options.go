package registry

import (
	"crypto/tls"
	"time"
)

type Options struct {
	//registry options
	Addrs     []string
	Timeout   time.Duration
	Secure    bool
	TLSConfig *tls.Config

	//consul agent check options
	CheckArgs     []string
	CheckInterval string //todo
	CheckTimeout  string //todo
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
func Args(a []string) Option {
	return func(o *Options) {
		o.CheckArgs = a
	}
}
