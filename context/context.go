package context

import (
	"net/http"
	"sync"

	"log"

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
	pool sync.Map
}

type Conversation struct {
	done     chan int
	ID       string
	Request  string
	Response http.ResponseWriter
}

func (ctx *context) Add(c *Conversation) string {
	if c == nil {
		return ""
	}
	if _, err := uuid.FromString(c.ID); err != nil {
		newUUID, _ := uuid.NewV4()
		c.ID = newUUID.String()
	}
	c.done = make(chan int)
	ctx.pool.Store(c.ID, c)

	return c.ID
}

func (ctx *context) Get(id string) *Conversation {
	conv, ok := ctx.pool.Load(id)

	if conv == nil || !ok {
		return nil
	}
	return conv.(*Conversation)
}

func (ctx *context) Wait(id string) {
	conv, ok := ctx.pool.Load(id)
	if conv == nil || !ok {
		log.Println("context wait error")
		return
	}
	select {
	case sig := <-conv.(*Conversation).done:
		if sig == 1 {
			return
		}
	}
}

func (ctx *context) Delete(id string) {
	ctx.pool.Delete(id)
}

func (ctx *context) Done(id string) {
	conv, ok := ctx.pool.Load(id)
	if conv == nil || !ok {
		log.Println("context done error")
		return
	}
	conv.(*Conversation).done <- 1
}

func NewContext() Context {
	return &context{}
}
