package router

import (
	"github.com/hashicorp/consul/api"
	"micro/codec"
)

type Options struct {
	name    string
	Client	 *api.Client	//consul client
	notFound         RejectHandler
	methodNotAllowed RejectHandler
}

type Option func(*Options)
type RejectHandler func(*codec.Request) error

//config consul client
func Client(c *api.Client) Option{
	return func(options *Options) {
		options.Client = c
	}
}

func Name(n string) Option {
	return func(options *Options) {
		options.name = n
	}
}