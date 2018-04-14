package router

import (
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/nzgogo/micro/codec"
	"github.com/nzgogo/micro/constant"
	"github.com/thedevsaddam/govalidator"
)

type Handler func(*codec.Message, string) *Error

type Error struct {
	StatusCode int
	Message    string
}

type Router interface {
	Init(opts ...Option) error
	Routes() []*Node
	Add(*Node)
	Dispatch(*codec.Message) (Handler, error)
	HttpMatch(*codec.Message) (*Node, error)
	Register() error
	Deregister() error
}

type router struct {
	routes []*Node
	opts   Options
}

type Node struct {
	Method             string              `json:"Method,omitempty"`
	Path               string              `json:"Path,omitempty"`
	ID                 string              `json:"ID"`
	Handler            Handler             `json:"-"`
	ValidationRules    govalidator.MapData `json:"ValidationRules,omitempty"`
	ValidationMessages govalidator.MapData `json:"ValidationMessages,omitempty"`
}

func (r *router) Init(opts ...Option) error {
	for _, o := range opts {
		o(&r.opts)
	}
	if *r.opts.Client == (api.Client{}) {
		var err error
		r.opts.Client, err = api.NewClient(api.DefaultConfig())
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *router) Routes() []*Node {
	return r.routes
}

// Pack supported service into a node struct and add to routes
func (r *router) Add(n *Node) {
	for _, route := range r.routes {
		if route.ID == n.ID {
			route = n
			return
		}
	}
	r.routes = append(r.routes, n)
}

// Put all nodes to consul key value store
func (r *router) Register() error {
	if len(r.routes) == 0 {
		return nil
	}

	if r.opts.Client == nil {
		return nil
	}

	publicRoutes := make([]*Node, 0)

	for _, v := range r.routes {
		if v.Method != "" {
			publicRoutes = append(publicRoutes, v)
		}
	}

	value, err := codec.Marshal(publicRoutes)
	if err != nil {
		return err
	}

	kv := r.opts.Client.KV()

	p := &api.KVPair{Key: r.opts.name, Value: value}
	if _, err := kv.Put(p, nil); err != nil {
		return err
	}

	return nil
}

// Delete all previously registered nodes from consul kv store
func (r *router) Deregister() error {
	if len(r.routes) == 0 {
		return nil
	}

	kv := r.opts.Client.KV()

	// TBC: check if this is the last service before remove KV
	// Delete all
	if _, err := kv.DeleteTree(r.opts.name, nil); err != nil {
		return err
	}

	return nil
}

// Based on reqeust.path and method (e.g GET /gogox/v1/greeter/hello),
// this method will download all relavent nodes from consul KV store
// according to parsed key (/gogox/v1/greeter) and find matching service (/hello)
func (r *router) HttpMatch(req *codec.Message) (*Node, error) {
	if req == nil {
		return nil, constant.ErrEmptyMsg
	}
	srvPath, subPath, err := r.splitPath(req.Path)
	if err != nil {
		return nil, err
	}

	routes, err := r.loadRemoteRoutes(srvPath)
	if err != nil {
		return nil, err
	}

	if len(routes) == 0 || routes == nil {
		return nil, constant.ErrResourceNotFound
	}

	if paths := r.pathMatch(routes, subPath); len(paths) > 0 {
		if node := r.methodMatch(paths, req.Method); node != nil {
			req.Node = node.ID
			return node, nil
		} else {
			return nil, constant.ErrMethodNotAllowed
		}
	}

	return nil, constant.ErrResourceNotFound
}

// dispatch to route handler
func (r *router) Dispatch(req *codec.Message) (Handler, error) {
	if req == nil {
		return nil, constant.ErrEmptyMsg
	}
	if req.Node == "" {
		return nil, constant.ErrResourceNotFound
	}

	for _, node := range r.routes {
		if node.ID == req.Node {
			return node.Handler, nil
		}
	}

	return nil, constant.ErrResourceNotFound
}

func (r *router) pathMatch(routes []*Node, path string) (matched []*Node) {
	matched = make([]*Node, 0)
	for _, node := range routes {
		if node.Path == path {
			matched = append(matched, node)
		}
	}
	return
}

func (r *router) methodMatch(paths []*Node, method string) (matched *Node) {
	for _, node := range paths {
		if node.Method == method {
			matched = node
		}
	}
	return
}

func (r *router) splitPath(path string) (srvPath, subPath string, err error) {
	processed := strings.Split(path, "/")
	results := make([]string, 0)
	for _, str := range processed {
		if str != "" {
			results = append(results, str)
		}
	}

	if len(results) <= 3 {
		err = constant.ErrRouterInvalidPath
		return
	}

	srvPath = constant.ORGANIZATION + "/" + results[1] + "/" + results[2] + "/" + results[0]
	for i := 3; i < len(results); i++ {
		subPath += "/" + results[i]
	}
	return
}

func (r *router) loadRemoteRoutes(key string) ([]*Node, error) {
	routes := make([]*Node, 0)
	kv := r.opts.Client.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return nil, err
	}

	if pair != nil {
		if err := codec.Unmarshal(pair.Value, &routes); err != nil {
			return nil, err
		}
		return routes, nil
	}
	return nil, nil
}

func NewRouter(opts ...Option) *router {
	options := Options{}

	for _, o := range opts {
		o(&options)
	}

	return &router{
		routes: make([]*Node, 0),
		opts:   options,
	}
}
