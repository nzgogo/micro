// Transport is an interface which is used for communication between services.

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
	SendFile(string, string, string, []byte) error
	Close() error
}

type transport struct {
	conn    *nats.Conn
	sub     *nats.Subscription
	opts    Options
	handler nats.MsgHandler
}

type ResponseHandler func([]byte) error

var (
	DefaultTransport = NewTransport()
)

const (
	DefaultRequestTimeout = time.Second * 15
)

func (n *transport) Options() Options {
	return n.opts
}

func (n *transport) Request(subject string, req []byte, handler ResponseHandler) error {
	rsp, respErr := n.conn.Request(subject, req, n.opts.Timeout)
	if respErr != nil {
		return respErr
	}

	if handler != nil {
		return handler(rsp.Data)
	}

	return nil
}

func (n *transport) Publish(subject string, data []byte) error {
	return n.conn.Publish(subject, data)
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

	clientOpts := nats.GetDefaultOptions()
	clientOpts.Servers = cAddrs
	if options.Timeout != 0 {
		clientOpts.Timeout = options.Timeout
	}

	conn, err := clientOpts.Connect()
	if err != nil {
		return err
	}

	n.conn = conn
	n.opts = options

	return nil
}

func NewTransport(opts ...Option) *transport {
	options := Options{
		Timeout: DefaultRequestTimeout,
	}

	for _, o := range opts {
		o(&options)
	}

	return &transport{
		opts: options,
	}
}
