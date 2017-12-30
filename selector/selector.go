package selector

import (
	"errors"
	"micro/registry"
	//"github.com/nzgogo/micro/registry"
)

type Selector interface {
	Init() error
	Options() Options
	// Select returns a function which should return the next node
	Select(service, version string) (string, error)
	// Mark sets the success/error against a node
	Mark(service string, node *registry.Node, err error)
	// Reset returns state back to zero for a service
	Reset(service string)
	// Close renders the selector unusable
	Close() error
}

type selector struct {
	opts Options
}

// Next is a function that returns the next node based on the selector's strategy
type Next func() (*registry.Node, error)

// Filter is used to filter a service during the selection process
type Filter func([]*registry.Service) []*registry.Service

// Strategy is a selection strategy e.g random, round robin
type Strategy func([]*registry.Service) Next

var (
	ErrNotFound      = errors.New("not found")
	ErrNoneAvailable = errors.New("none available")
)

func NewSelector(opts ...Option) Selector {
	sopts := Options{}

	for _, opt := range opts {
		opt(&sopts)
	}

	return &selector{
		opts: sopts,
	}
}

func (r *selector) Init() error {
	if r.opts.Strategy == nil {
		r.opts.Strategy = Random
	}
	if r.opts.Registry == nil {
		r.opts.Registry = registry.NewRegistry()
	}

	return nil
}

func (r *selector) Options() Options {
	return r.opts
}

func (r *selector) Select(service, version string) (string, error) {
	// get the service
	services, err := r.opts.Registry.GetService(service)
	if err != nil {
		return "", err
	}

	// apply the filters
	//for _, filter := range r.opts.Filters {
	//	services = filter(services)
	//}
	filterVersion(services,version)

	// if there's nothing left, return
	if len(services) == 0 {
		return "", ErrNoneAvailable
	}

	next := r.opts.Strategy(services)
	node, err:= next()
	if err != nil {
		return "", err
	}

	return service+"."+version+"."+node.Id,nil
}

func (r *selector) Mark(service string, node *registry.Node, err error) {
	return
}

func (r *selector) Reset(service string) {
	return
}

func (r *selector) Close() error {
	return nil
}

func filterVersion(old []*registry.Service, version string) []*registry.Service {
	var services []*registry.Service

	for _, service := range old {
		if service.Version == version {
			services = append(services, service)
		}
	}

	return services
}