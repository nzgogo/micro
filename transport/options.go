package transport

import (
	"time"
)

type Options struct {
	Subject string   //  Message destination address
	Addrs   []string // A configured set of nats servers which this client will use when attempting to connect.
	Timeout time.Duration
}

type Option func(*Options)

// subject to use for transport
func Subject(sub string) Option {
	return func(o *Options) {
		o.Subject = sub
	}
}

// Addrs to use for transport
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// Timeout sets the timeout for Send/Recv execution
func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}
