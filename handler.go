package gogo

import (
	"log"
	"micro/constant"

	"github.com/nats-io/go-nats"
	"github.com/nzgogo/micro/api"
	"github.com/nzgogo/micro/codec"
	"github.com/nzgogo/micro/context"
)

func (s *service) ServerHandler(nMsg *nats.Msg) {
	if nMsg == nil {
		log.Println("Nats body empty")
		return
	}
	//decode message
	message := &codec.Message{}
	uerr := codec.Unmarshal(nMsg.Data, message)
	if uerr != nil {
		log.Println("ServerHandler respond error: Unmarshal failed")
		return
	}
	sub := s.opts.Transport.Options().Subject

	//check message type, response or request
	if message.Type == constant.REQUEST {
		// check the message type: Request or Publish.
		// If it is a Request, the reply subject should be extracted from nats.Msg struct.
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
			err := s.Respond(errResp, message.ReplyTo)
			if err != nil {
				log.Println("ServerHandler respond error: " + err.Error())
			}
			return
		}
		reply := message.ReplyTo
		message.ReplyTo = sub

		go func() {
			for i := len(s.opts.HdlrWrappers); i > 0; i-- {
				handler = s.opts.HdlrWrappers[i-1](handler)
			}
			err := handler(message, reply)
			if err != nil {
				body := map[string]string{"message": err.Message}
				err1 := s.Respond(
					codec.NewJsonResponse(
						message.ContextID,
						err.StatusCode,
						body,
					),
					reply,
				)
				if err1 != nil {
					log.Println("ServerHandler respond error: " + err1.Error())
				}
			}
		}()
	} else if message.Type == constant.HEALTHCHECK {
		go func() {
			checkStatus, feedback := healthCheck(s.config)
			msg := codec.NewJsonResponse("", checkStatus, feedback)
			replyBody, _ := codec.Marshal(msg)
			err := s.opts.Transport.Publish(nMsg.Reply, replyBody)
			if err != nil {
				log.Println("ServerHandler respond error: " + err.Error())
			}
		}()
	} else if message.Type == constant.RESPONSE {
		conversation := s.opts.Context.Get(message.ContextID)
		if conversation == nil {
			log.Println("ServerHandler respond error: conversation lost")
			return
		}
		rpl := conversation.Request
		s.opts.Transport.Publish(rpl, nMsg.Data)
		s.opts.Context.Delete(message.ContextID)
	}
}

func (s *service) ApiHandler(nMsg *nats.Msg) {
	if nMsg == nil {
		log.Println("Nats body empty")
		return
	}
	message := &codec.Message{}
	uerr := codec.Unmarshal(nMsg.Data, message)
	if uerr != nil {
		log.Println("ApiHandler respond error: Unmarshal failed")
		return
	}
	ctx := s.opts.Context
	if message.Type == constant.HEALTHCHECK {
		go func() {
			checkStatus, feedback := healthCheck(s.config)
			msg := codec.NewJsonResponse("", checkStatus, feedback)
			replyBody, _ := codec.Marshal(msg)
			err := s.opts.Transport.Publish(nMsg.Reply, replyBody)
			if err != nil {
				log.Println("ServerHandler respond error: " + err.Error())
			}
		}()

	} else if message.Type == constant.RESPONSE {
		conversation := ctx.Get(message.ContextID)
		if conversation == nil {
			log.Println("ApiHandler respond error: conversation lost")
			return
		}
		r := conversation.Response

		fn := gogoapi.WriteResponse
		for i := len(s.opts.HttpRespWrappers); i > 0; i-- {
			fn = s.opts.HttpRespWrappers[i-1](fn)
		}
		fn(r, message)
		ctx.Done(message.ContextID)
		ctx.Delete(message.ContextID)
	}
}
