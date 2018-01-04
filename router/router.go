package router

import (
	"errors"

	//micro "micro"
	"micro/codec"
	"micro/transport"
	"github.com/hashicorp/consul/api"
	"strings"
)

type Handler func(*codec.Request, transport.Transport, string) error

type Router interface {
	Init(opts ...Option) error
	Add(*Node)
	Dispatch(*codec.Request) (Handler, error)
	HttpMatch(*codec.Request) error
	Register(key string) error
	Deregister(key string) error
	String() string
}

type router struct {
	routes      []*Node
	opts		Options
}

type Node struct {
	Method  string	`json:"Method,omitempty"`
	Path    string	`json:"Path,omitempty"`
	ID      string	`json:"ID"`
	Handler Handler	`json:"-"`
}

var(
	ErrNotFound = errors.New("Service not found")
	ErrmethodNotAllowed = errors.New("Method not allowed")
	Codec = codec.NewCodec()
)

func (r *router) Init(opts ...Option) error {
	for _, o := range opts {
		o(&r.opts)
	}
	if *r.opts.Client == (api.Client{}) {
		var err error
		r.opts.Client,err = api.NewClient(api.DefaultConfig())
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
func (r *router) Register(key string) error {
	if len(r.routes) == 0 {
		return nil
	}
	value, err := Codec.Marshal(&r.routes)
	if err != nil {
		return err
	}

	kv := r.opts.Client.KV()

	p := &api.KVPair{Key: key, Value: value}
	if _, err := kv.Put(p, nil); err != nil {
		return err
	}

	return nil
}

// Delete all previously registered nodes from consul kv store
func (r *router) Deregister(key string) error {
	if len(r.routes) == 0 {
		return nil
	}

	kv := r.opts.Client.KV()

	// Delete all
	if _, err := kv.DeleteTree(key, nil); err != nil {
		return err
	}

	return nil
}

// Based on reqeust.path  and method (e.g GET /gogox/v1/greeter/hello),
// this method will download all relavent nodes from consul KV store
// according to parsed key (/gogox/v1/greeter) and find matching service (/hello)
func (r *router) HttpMatch(req *codec.Request) error {
	key, subpath := PathToKeySubpath(req.Path)
	if err:=r.loadRouterNodes(key);err!=nil{
		return err
	}

	if len(r.routes) == 0{
		if r.opts.notFound != nil{
			return r.opts.notFound(req)
		}
		return ErrNotFound
	}

	for _, node := range r.routes {
		//service search
		if node.Path == subpath {
			//found service, check if request method supported
			if node.Method == req.Method {
				var err error
				req.Node, err = Codec.Marshal(node)
				if err != nil{
					return err
				}
				return nil	//found service and method supported
			}
			//call reject request handler if any, otherwise return error ErrmethodNotAllowed
			if r.opts.methodNotAllowed != nil{
				return r.opts.methodNotAllowed(req)
			}
			return ErrmethodNotAllowed
		}
	}
	//call reject handler if any, otherwise return error ErrNotFound
	if r.opts.notFound != nil{
		return r.opts.notFound(req)
	}
	return ErrNotFound
}

//Once server decoded a request, it can use this method to dispatch tasks to relevant handlers
func (r *router) Dispatch(req *codec.Request) (Handler, error) {
	reqNode := &Node{}
	Codec.Unmarshal(req.Node, reqNode)
	if reqNode ==nil {
		return nil, ErrNotFound
	}
	for _, node := range r.routes {
		if node.ID == reqNode.ID {
			return node.Handler, nil
		}
	}
	return nil, ErrNotFound
}
//TODO may need to move this shit to somewhere else
func (r *router) loadRouterNodes(key string) error {
	r.routes = make([]*Node,0)
	kv := r.opts.Client.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return nil
	}

	if pair != nil {
		if err := Codec.Unmarshal(pair.Value, &r.routes); err != nil{
			return err
		}
	}
	return nil
}

func PathToKeySubpath(path string) (key, subpath string){
	i := 0
	for m:=0;m<3;m++{
		x := strings.Index(path[i+1:],"/")
		if x < 0 {
			break
		}
		i += x
		i++
	}

	return path[:i], path[i:]
}

// Return routes key used by service kv store
func (r *router) String() string{
	return r.opts.name
}

func NewRouter(opts ...Option) *router {
	options := Options{}

	for _, o := range opts {
		o(&options)
	}

	return &router{
		routes: make([]*Node,0),
		opts: options,
	}
}