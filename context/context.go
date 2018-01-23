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
	// sync.RWMutex
	// pool map[string]*Conversation
	pool sync.Map
}

type Conversation struct {
	done     chan int
	ID       string
	Request  string
	Response http.ResponseWriter
}

func (ctx *context) Add(c *Conversation) string {
	// defer ctx.Unlock()
	// ctx.Lock()

	if _, err := uuid.FromString(c.ID); err != nil {
		newUUID, _ := uuid.NewV4()
		c.ID = newUUID.String()
	}
	c.done = make(chan int)

	// ctx.pool[c.ID] = c
	ctx.pool.Store(c.ID, c)

	return c.ID
}

func (ctx *context) Get(id string) *Conversation {
	// defer ctx.Unlock()
	// ctx.Lock()
	//
	// return ctx.pool[id]
	conv, _ := ctx.pool.Load(id)
	return conv.(*Conversation)
}

func (ctx *context) Wait(id string) {
	// defer ctx.RUnlock()
	// ctx.RLock()
	//
	// select {
	// case sig := <-ctx.pool[id].done:
	// 	if sig == 1 {
	// 		return
	// 	}
	// }

	conv, _ := ctx.pool.Load(id)
	select {
	case sig := <-conv.(*Conversation).done:
		if sig == 1 {
			return
		}
	}
}

func (ctx *context) Delete(id string) {
	// defer ctx.Unlock()
	// ctx.Lock()
	//
	// delete(ctx.pool, id)

	ctx.pool.Delete(id)
}

func (ctx *context) Done(id string) {
	// defer ctx.Unlock()
	// ctx.Lock()
	//
	// ctx.pool[id].done <- 1

	conv, _ := ctx.pool.Load(id)
	conv.(*Conversation).done <- 1
}

func NewContext() Context {
	return &context{}
}
