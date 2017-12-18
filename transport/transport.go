// Transport is an interface which is used for communication between
// services. It uses NATS implementations

package transport

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/go-nats"
	//"github.com/nzgogo/micro/codec"
	"micro/codec"
)

type Transport interface {
	Options() Options
	Init(...Option) error
	Request([]byte, interface{}) error
	Publish([]byte) error
	Subscribe() error
	Close() error
}

type Nats struct {
	conn    *nats.Conn
	addr    string //nats subject
	rplAddr string //refers to nats.Msg.reply, for a subscriber use only
	sub     *nats.Subscription
	opts    Options
}

var (
	DefaultTimeout     = time.Second * 15
	DefaultDialTimeout = time.Second * 5
)

func (n *Nats) TestConnection() error {
	if n.conn == nil {
		return fmt.Errorf("natsproxy: Connection cannot be nil")
	}
	if n.conn.Status() != nats.CONNECTED {
		return fmt.Errorf("Client not connected")
	}
	return nil
}

func (n *Nats) Request(req *codec.Request, resp *codec.Response) error {
	var Codec codec.Codec
	b, err := Codec.Marshal(req)
	if err != nil {
		return err
	}

	rsp, respErr := n.conn.Request(n.addr, b, n.opts.Timeout)
	if respErr != nil {
		return respErr
	}

	if err := Codec.Unmarshal(rsp.Data, resp); err != nil {
		return err
	}

	return nil
}

func (n *Nats) Publish(m *codec.Request) error {
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

func (n *Nats) Subscribe(resp *codec.Request) error {
	sub, err := n.conn.SubscribeSync(n.addr)
	n.sub = sub
	if err != nil {
		return err
	}

	timeout := time.Second * 10
	if n.opts.Timeout > time.Duration(0) {
		timeout = n.opts.Timeout
	}

	rsp, err := n.sub.NextMsg(timeout)
	if err != nil {
		return err
	}

	var mr *codec.Request
	var Codec codec.Codec
	if err := Codec.Unmarshal(rsp.Data, &mr); err != nil {
		return err
	}
	n.rplAddr = rsp.Reply
	resp = mr
	return nil
}

func (n *Nats) Close() error {
	n.sub.Unsubscribe()
	n.conn.Close()
	return nil
}

func (n *Nats) Init(opts ...Option) error {
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

	client_opts := nats.GetDefaultOptions()
	client_opts.Servers = cAddrs
	client_opts.Timeout = options.Timeout

	// secure might not be set
	if client_opts.TLSConfig != nil {
		client_opts.Secure = true
	}

	c, err := client_opts.Connect()
	if err != nil {
		return err
	}

	options.Timeout = DefaultDialTimeout

	n.conn = c
	n.addr = options.Subject
	n.opts = options

	return nil
}

func NewTransport() *Nats {
	return &Nats{}
}
