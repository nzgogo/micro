package context

import (
	"net/http"
	"sync"

	"github.com/satori/go.uuid"
)

type Context interface {
	Add(*Conversation) string
	Get(string) *Conversation
	Delete(id string)
	Wait(string)
	Done(string)
}

type context struct {
	sync.Mutex
	pool map[string]*Conversation
}

type Conversation struct {
	done     chan int
	ID       string
	Request  string
	Response http.ResponseWriter
}

func (ctx *context) Add(c *Conversation) string {
	defer ctx.Unlock()
	ctx.Lock()

	if _, err := uuid.FromString(c.ID); err != nil {
		newUUID, _ := uuid.NewV4()
		c.ID = newUUID.String()
	}
	c.done = make(chan int)

	ctx.pool[c.ID] = c

	return c.ID
}

func (ctx *context) Get(id string) *Conversation {
	defer ctx.Unlock()
	ctx.Lock()

	return ctx.pool[id]
}

func (ctx *context) Wait(id string) {
	defer ctx.Unlock()
	ctx.Lock()

	select {
	case sig := <-ctx.pool[id].done:
		if sig == 1 {
			return
		}
	}
}

func (ctx *context) Delete(id string) {
	defer ctx.Unlock()
	ctx.Lock()

	delete(ctx.pool, id)
}

func (ctx *context) Done(id string) {
	defer ctx.Unlock()
	ctx.Lock()

	ctx.pool[id].done <- 1
}

func NewContext() Context {
	ctx := &context{
		pool: make(map[string]*Conversation),
	}

	return ctx
}
