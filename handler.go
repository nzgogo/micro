package gogo

import (
	"fmt"

	"github.com/nats-io/go-nats"
	"github.com/nzgogo/micro/api"
	"github.com/nzgogo/micro/codec"
	"github.com/nzgogo/micro/context"
)

func (s *service) ServerHandler(nMsg *nats.Msg) {
	//decode message
	message := &codec.Message{}
	codec.Unmarshal(nMsg.Data, message)

	sub := s.opts.Transport.Options().Subject

	//check message type, response or request
	if message.Type == "request" {
		//check if the message is a Request or Publish.
		fmt.Println("nMsg.reply : " + nMsg.Reply)
		if nMsg.Reply != "" {
			message.ReplyTo = nMsg.Reply
		}
		contxt := s.Options().Context
		contxt.Add(&context.Conversation{
			ID:      message.ContextID,
			Request: message.ReplyTo,
		})

		handler, routerErr := s.opts.Router.Dispatch(message)
		if routerErr != nil {
			errResp := gogoapi.NewResponse(404, message.ContextID, nil, message.Header)
			resp, _ := codec.Marshal(errResp)
			s.opts.Transport.Publish(message.ReplyTo, resp)
		}
		reply := message.ReplyTo
		message.ReplyTo = sub
		//TODO: error handle
		go handler(message, reply)
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
