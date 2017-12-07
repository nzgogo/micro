package transport

import (
	"strings"
	"time"
	"errors"
	"github.com/nats-io/go-nats"
	"github.com/nzgogo/micro/codec"
)

type Message struct {
	Header map[string]string
	Body   []byte
}

type Client interface {
	Recv(*Message) error
	Send(*Message) error
	Close() error
}

type ntport struct {
	addrs []string
	opts  Options
}

type ntportClient struct {
	conn *nats.Conn
	addr string
	id   string
	sub  *nats.Subscription
	opts Options
}


type Option func(*Options)

type DialOption func(*DialOptions)

var (
	DefaultDialTimeout = time.Second * 5
)

// Transport is an interface which is used for communication between
// services. It uses NATS implementations
type Transport interface {
	Dial(addr string, opts ...DialOption) (Client, error)
}

var (
	DefaultTimeout = time.Minute
)

func (n *ntportClient) AsynSend(m *Message) error {
	var Codec codec.Codec
	b, err := Codec.Marshal(m)
	if err != nil {
		return err
	}

	// no deadline
	if n.opts.Timeout == time.Duration(0) {
		return n.conn.Publish(n.addr, b)
	}

	// use the deadline
	ch := make(chan error, 1)

	go func() {
		ch <- n.conn.Publish(n.addr, b)
	}()

	select {
	case err := <-ch:
		return err
	case <-time.After(n.opts.Timeout):
		return errors.New("deadline exceeded")
	}
}

func (n *ntportClient) Send(m *Message) error {
	var Codec codec.Codec
	b, err := Codec.Marshal(m)
	if err != nil {
		return err
	}

	// no deadline
	if n.opts.Timeout == time.Duration(0) {
		return n.conn.PublishRequest(n.addr, n.id, b)
	}

	// use the deadline
	ch := make(chan error, 1)

	go func() {
		ch <- n.conn.PublishRequest(n.addr, n.id, b)
	}()

	select {
	case err := <-ch:
		return err
	case <-time.After(n.opts.Timeout):
		return errors.New("deadline exceeded")
	}
}

func (n *ntportClient) Recv(m *Message) error {
	timeout := time.Second * 10
	if n.opts.Timeout > time.Duration(0) {
		timeout = n.opts.Timeout
	}

	rsp, err := n.sub.NextMsg(timeout)
	if err != nil {
		return err
	}

	var mr Message
	var Codec codec.Codec
	if err := Codec.Unmarshal(rsp.Data, &mr); err != nil {
		return err
	}

	*m = mr
	return nil
}

func (n *ntportClient) Close() error {
	n.sub.Unsubscribe()
	n.conn.Close()
	return nil
}

func (n *ntport) Dial(addr string, dialOpts ...DialOption) (Client, error) {
	dopts := DialOptions{
		Timeout: DefaultDialTimeout,
	}

	for _, o := range dialOpts {
		o(&dopts)
	}

	opts := nats.GetDefaultOptions()
	opts.Servers = n.addrs
	opts.Secure = n.opts.Secure
	opts.TLSConfig = n.opts.TLSConfig
	opts.Timeout = dopts.Timeout

	// secure might not be set
	if n.opts.TLSConfig != nil {
		opts.Secure = true
	}

	c, err := opts.Connect()
	if err != nil {
		return nil, err
	}

	id := nats.NewInbox()
	sub, err := c.SubscribeSync(id)
	if err != nil {
		return nil, err
	}

	return &ntportClient{
		conn: c,
		addr: addr,
		id:   id,
		sub:  sub,
		opts: n.opts,
	}, nil
}

func NewTransport(opts ...Option) Transport {
	options := Options{
		Timeout: DefaultTimeout,
	}

	for _, o := range opts {
		o(&options)
	}

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

	return &ntport{
		addrs: cAddrs,
		opts:  options,
	}
}
