package context

import (
	"net/http"

	"github.com/satori/go.uuid"
)

type Context interface {
	Add(*Conversation) string
	Get(string) *Conversation
	Wait(string)
	Done(string)
}

type context struct {
	pool map[string]*Conversation
}

type Conversation struct {
	done     chan int
	Request  string
	Response *http.ResponseWriter
}

func (ctx context) Add(c *Conversation) string {
	id := uuid.NewV4()
	c.done = make(chan int)

	ctx.pool[id.String()] = c

	return id.String()
}

func (ctx context) Get(id string) *Conversation {
	return ctx.pool[id]
}

func (ctx context) Wait(id string) {
	select {
	case sig := <-ctx.pool[id].done:
		if sig == 1 {
			return
		}
	}
}

func (ctx context) Done(id string) {
	ctx.pool[id].done <- 1
}

func NewContext() Context {
	ctx := context{
		pool: make(map[string]*Conversation),
	}

	return ctx
}
