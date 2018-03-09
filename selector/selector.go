package selector

import (
	"strings"

	"github.com/nzgogo/micro/constant"
	"github.com/nzgogo/micro/registry"
)

type Selector interface {
	Init() error
	Select(service, version string) (string, error)
}

type selector struct {
	registry registry.Registry
	opts     Options
}

// Next is a function that returns the next node based on the selector's strategy
type Next func() (*registry.Node, error)

// Strategy is a selection strategy e.g random, round robin
type Strategy func([]*registry.Service) Next

func (r *selector) Init() error {
	if r.opts.Strategy == nil {
		r.opts.Strategy = Random
	}
	if r.registry == nil {
		return constant.ErrRegistryEmptyNode
	}

	return nil
}

// return an available service
func (r *selector) Select(service, version string) (string, error) {
	services, err := r.registry.GetService(service)
	if err != nil {
		return "", err
	}

	filterVersion(services, version)

	if len(services) == 0 {
		return "", constant.ErrSelectNoneAvailable
	}

	next := r.opts.Strategy(services)
	node, err := next()
	if err != nil {
		return "", err
	}

	return strings.Replace(service, "-", ".", -1) + "." + version + "." + node.ID, nil
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

func NewSelector(registry registry.Registry, opts ...Option) Selector {
	sOpts := Options{}

	for _, opt := range opts {
		opt(&sOpts)
	}

	return &selector{
		registry: registry,
		opts:     sOpts,
	}
}
