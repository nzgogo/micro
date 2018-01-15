package context

import (
	"net/http"
	"sync"

	"github.com/satori/go.uuid"
)

var mutex sync.Mutex

type Context interface {
	Add(*Conversation) string
	Get(string) *Conversation
	Delete(id string)
	Wait(string)
	Done(string)
}

type context struct {
	pool map[string]*Conversation
}

type Conversation struct {
	done     chan int
	ID       string
	Request  string
	Response http.ResponseWriter
}

func (ctx *context) Add(c *Conversation) string {
	if _, err := uuid.FromString(c.ID); err != nil {
		newUUID, _ := uuid.NewV4()
		c.ID = newUUID.String()
	}
	c.done = make(chan int)

	mutex.Lock()
	ctx.pool[c.ID] = c
	mutex.Unlock()

	return c.ID
}

func (ctx *context) Get(id string) *Conversation {
	return ctx.pool[id]
}

func (ctx *context) Wait(id string) {
	select {
	case sig := <-ctx.pool[id].done:
		if sig == 1 {
			return
		}
	}
}

func (ctx *context) Delete(id string) {
	mutex.Lock()
	delete(ctx.pool, id)
	mutex.Unlock()
}

func (ctx *context) Done(id string) {
	mutex.Lock()
	ctx.pool[id].done <- 1
	mutex.Unlock()
}

func NewContext() Context {
	ctx := &context{
		pool: make(map[string]*Conversation),
	}

	return ctx
}
