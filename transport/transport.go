// Transport is an interface which is used for communication between
// services. It uses NATS implementations

package transport

import (
	"errors"
	"strings"
	"time"

	"github.com/nats-io/go-nats"
)

type Transport interface {
	Options() Options
	Init() error
	Request(string, []byte, ResponseHandler) error
	Publish(string, []byte) error
	Close() error
}

type transport struct {
	conn *nats.Conn
	sub  *nats.Subscription
	opts Options
}

type ResponseHandler func([]byte) error

var (
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

	return handler(rsp.Data)
}

func (n *transport) Publish(sub string, b []byte) error {

	// no deadline
	if n.opts.Timeout == time.Duration(0) {
		return n.conn.Publish(sub, b)
	}

	// use the deadline
	ch := make(chan error, 1)

	go func() {
		ch <- n.conn.Publish(sub, b)
	}()

	select {
	case err := <-ch:
		return err
	case <-time.After(n.opts.Timeout):
		return errors.New("deadline exceeded")
	}
}

func (n *transport) Close() error {
	n.sub.Unsubscribe()
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

	sub, err := n.conn.SubscribeSync(options.Subject)
	if err != nil {
		return err
	}

	n.conn = c
	n.opts = options
	n.sub = sub

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
