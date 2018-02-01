package gogo

import (
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
			errResp := codec.NewResponse(message.ContextID, 404, nil, message.Header)
			resp, _ := codec.Marshal(errResp)
			s.opts.Transport.Publish(message.ReplyTo, resp)
		}
		reply := message.ReplyTo
		message.ReplyTo = sub
		//TODO: error handle
		go func() {
			for i := len(s.opts.HdlrWrappers); i > 0; i-- {
				handler = s.opts.HdlrWrappers[i-1](handler)
			}
			err := handler(message, reply)
			if err != nil {
				s.Respond(
					codec.NewJsonResponse(
						message.ContextID,
						err.StatusCode,
						err.Message,
					),
					reply,
				)
			}
		}()
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

	fn := gogoapi.WriteResponse
	for i := len(s.opts.HttpRespWrappers); i > 0; i-- {
		fn = s.opts.HttpRespWrappers[i-1](fn)
	}
	fn(r, message)

	ctx.Done(message.ContextID)
	ctx.Delete(message.ContextID)

}
