package selector

import (
	"micro/registry"
	//"github.com/nzgogo/micro/registry"
)

type Options struct {
	Registry registry.Registry
	Strategy Strategy
	Filters  []Filter
}

// Option used to initialise the selector
type Option func(*Options)

// Registry sets the registry used by the selector
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// SetStrategy sets the default strategy for the selector
func SetStrategy(fn Strategy) Option {
	return func(o *Options) {
		o.Strategy = fn
	}
}

// WithFilter adds a filter function to the list of filters
// used during the Select call.
func WithFilter(fn ...Filter) Option {
	return func(o *Options) {
		o.Filters = append(o.Filters, fn...)
	}
}

// Strategy sets the selector strategy
func WithStrategy(fn Strategy) Option {
	return func(o *Options) {
		o.Strategy = fn
	}
}
