package router

import (
	"errors"
	"log"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/nzgogo/micro/codec"
)

type Handler func(*codec.Message, string) *Error

type Error struct {
	StatusCode int
	Message    string
}

type Router interface {
	Init(opts ...Option) error
	Add(*Node)
	Dispatch(*codec.Message) (Handler, error)
	HttpMatch(*codec.Message) error
	Register() error
	Deregister() error
}

type router struct {
	routes []*Node
	opts   Options
}

type Node struct {
	Method  string  `json:"Method,omitempty"`
	Path    string  `json:"Path,omitempty"`
	ID      string  `json:"ID"`
	Handler Handler `json:"-"`
}

var (
	ErrInvalidPath      = errors.New("invalid path cannot process")
	ErrNotFound         = errors.New("service not found")
	ErrMethodNotAllowed = errors.New("method not allowed")
)

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
		log.Println("this is router client failure")
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

// Based on reqeust.path  and method (e.g GET /gogox/v1/greeter/hello),
// this method will download all relavent nodes from consul KV store
// according to parsed key (/gogox/v1/greeter) and find matching service (/hello)
func (r *router) HttpMatch(req *codec.Message) error {
	srvPath, subPath, err := r.splitPath(req.Path)
	if err != nil {
		return err
	}

	routes, err := r.loadRemoteRoutes(srvPath)
	if err != nil {
		return err
	}

	if len(routes) == 0 || routes == nil {
		return ErrNotFound
	}

	// for _, node := range routes {
	// 	if node.Path == subPath {
	// 		if node.Method == req.Method {
	// 			req.Node = node.ID
	// 			return nil
	// 		} else {
	// 			return ErrMethodNotAllowed
	// 		}
	// 	}
	// }

	if paths := r.pathMatch(subPath); len(paths) > 0 {
		log.Printf("%d path(s) matched!\n", len(paths))
		if node := r.methodMatch(paths, req.Method); node != nil {
			req.Node = node.ID
		} else {
			return ErrMethodNotAllowed
		}
	}

	return ErrNotFound
}

//Once server decoded a request, it can use this method to dispatch tasks to relevant handlers
func (r *router) Dispatch(req *codec.Message) (Handler, error) {
	if req.Node == "" {
		return nil, ErrNotFound
	}

	for _, node := range r.routes {
		if node.ID == req.Node {
			return node.Handler, nil
		}
	}

	return nil, ErrNotFound
}

func (r *router) pathMatch(path string) (matched []*Node) {
	matched = make([]*Node, 0)
	for _, node := range r.routes {
		if node.Path == path {
			matched = append(matched, node)
			log.Printf("%v\n", node)
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
		err = ErrInvalidPath
		return
	}

	srvPath = "gogo/" + results[1] + "/" + results[2] + "/" + results[0]
	//fmt.Println("srvpath: " + srvPath)
	for i := 3; i < len(results); i++ {
		subPath += "/" + results[i]
	}
	//fmt.Println("subpath: " + subPath)
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
