package gogo

import (
	"github.com/nats-io/go-nats"
	"github.com/nzgogo/micro/api"
	"github.com/nzgogo/micro/codec"
	"github.com/nzgogo/micro/context"
)

func (s *service) ServerHandler(nMsg *nats.Msg) {
	message := &codec.Message{}
	codec.Unmarshal(nMsg.Data, message)
	if message.ReplyTo =="nats-request" {
		message.ReplyTo = nMsg.Reply
	}
	if message.Type == "request" {
		//TODO if this is last endpoint in a serial call, we should not add this conversation
		contxt := s.Options().Context
		contxt.Add(&context.Conversation{
			ID:		 message.ContextID,
			Request: message.ReplyTo,
		})

		handler, routerErr := s.opts.Router.Dispatch(message)
		if routerErr != nil {
			resp, _ := codec.Marshal(codec.Message{
				StatusCode: 404,
				Header:     make(map[string][]string, 0),
				Body:       routerErr.Error(),
			})
			s.opts.Transport.Publish(message.ReplyTo, resp)
		}
		go handler(message)
	} else {
		rpl := s.opts.Context.Get(message.ContextID).Request
		s.opts.Transport.Publish(rpl, nMsg.Data)
		s.opts.Context.Delete(message.ContextID)
	}
}

//Example MsgHandler
func (s *service) ApiHandler(nMsg *nats.Msg) {
	message := &codec.Message{}
	codec.Unmarshal(nMsg.Data, message)
	ctx := s.opts.Context

	r := ctx.Get(message.ContextID).Response

	gogoapi.WriteResponse(r, message)

	ctx.Done(message.ContextID)
	ctx.Delete(message.ContextID)

}
