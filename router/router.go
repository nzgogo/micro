package router

import (
	"micro/codec"
	"micro/transport"

	"github.com/micro/go-micro/registry"
)

type Router interface {
	Init()
	Add(*Node)
	Dispatch(*codec.Request)
	HttpMatch(*codec.Request, *registry.Registry)
}

type router struct {
	routes           []*Node
	notFound         Handler
	methodNotAllowed Handler
}

type Handler func(*codec.Request, *transport.Transport)

type Node struct {
	Method  string
	Path    string
	ID      string
	Handler Handler
}

func (r *router) Add(n *Node) {
	for _, route := range r.routes {
		if route.ID == n.ID {
			route = n
			return
		}
	}
	r.routes = append(r.routes, n)
}

func (r *router) HttpMatch(req *codec.Request, reg *registry.Registry) {

}

func (r *router) Dispatch(req *codec.Request) {

}

func NewRouter() *router {
	return &router{}
}
