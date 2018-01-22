// Transport is an interface which is used for communication between
// services. It uses NATS implementations

package transport

import (
	"strings"
	"time"

	"github.com/nats-io/go-nats"
)

type Transport interface {
	Options() Options
	Init() error
	Request(string, []byte, ResponseHandler) error
	Publish(string, []byte) error
	Subscribe() error
	SetHandler(nats.MsgHandler)
	Close() error
}

type transport struct {
	conn *nats.Conn
	sub  *nats.Subscription
	opts Options
	handler nats.MsgHandler
}

type ResponseHandler func([]byte) error

var (
	DefaultTransport   = NewTransport()
	DefaultTimeout     = time.Second * 15
	DefaultDialTimeout = time.Second * 5
)

func (n *transport) Options() Options {
	return n.opts
}

func (n *transport) Request(sub string, req []byte, handler ResponseHandler) error {

	rsp, respErr := n.conn.Request(sub, req, n.opts.Timeout)
	if respErr != nil {
		return respErr
	}

	if handler != nil {
		return handler(rsp.Data)
	}

	return nil
}

func (n *transport) Publish(sub string, b []byte) error {

	// no deadline
	if n.opts.Timeout == time.Duration(0) {
		return n.conn.Publish(sub, b)
	}

	return n.conn.Publish(sub, b)
}

func (n *transport) Subscribe() error {
	var err error
	n.sub, err = n.conn.Subscribe(n.opts.Subject, n.handler)
	return err
}

func (n *transport) SetHandler(handler nats.MsgHandler) {
	n.handler = handler
}

func (n *transport) Close() error {
	//n.sub.Unsubscribe()
	n.conn.Close()
	return nil
}

func (n *transport) Init() error {
	options := n.opts

	var cAddrs []string

	for _, addr := range options.Addrs {
		if len(addr) == 0 {
			continue
		}
		if !strings.HasPrefix(addr, "nats://") {
			addr = "nats://" + addr
		}
		cAddrs = append(cAddrs, addr)
	}

	if len(cAddrs) == 0 {
		cAddrs = []string{nats.DefaultURL}
	}

	client_opts := nats.GetDefaultOptions()
	client_opts.Servers = cAddrs
	client_opts.Timeout = options.Timeout

	c, err := client_opts.Connect()
	if err != nil {
		return err
	}

	options.Timeout = DefaultDialTimeout

	if err != nil {
		return err
	}

	n.conn = c
	n.opts = options
	//n.sub = sub

	return nil
}

func NewTransport(opts ...Option) *transport {
	options := Options{
		Timeout: DefaultTimeout,
	}

	for _, o := range opts {
		o(&options)
	}

	return &transport{
		opts: options,
	}
}
