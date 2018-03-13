package gogo

import (
	"github.com/nats-io/go-nats"
	"github.com/nzgogo/micro/api"
	"github.com/nzgogo/micro/codec"
	"github.com/nzgogo/micro/context"
	recpro "github.com/nzgogo/micro/recover"
)

const (
	REQUEST     = "request"
	RESPONSE    = "response"
	HEALTHCHECK = "healthCheck"
)

func (s *service) ServerHandler(nMsg *nats.Msg) {
	defer func() {
		if rMsg := recover(); rMsg != nil {
			recpro.Recover(s.Options().Transport.Options().Subject, "ServerHandler", rMsg, nMsg)
		}
	}()
	if nMsg == nil {
		panic("Nats body empty")
	}
	//decode message
	message := &codec.Message{}
	uerr := codec.Unmarshal(nMsg.Data, message)
	if uerr != nil {
		panic("ServerHandler respond error: Unmarshal failed")
	}

	//check message type, response or request
	if message.Type == REQUEST {
		s.ServerHandlerRequest(message, nMsg.Reply)
	} else if message.Type == HEALTHCHECK {
		s.healthCheckHandler(message, nMsg.Reply)
	} else if message.Type == RESPONSE {
		s.ServerHandlerResponse(message, nMsg.Data)
	}
}

func (s *service) ApiHandler(nMsg *nats.Msg) {
	defer func() {
		if rMsg := recover(); rMsg != nil {
			recpro.Recover(s.Options().Transport.Options().Subject, "ApiHandler", rMsg, nMsg)
		}
	}()

	if nMsg == nil {
		panic("Nats body empty")
	}

	message := &codec.Message{}
	uerr := codec.Unmarshal(nMsg.Data, message)
	if uerr != nil {
		panic("ApiHandler respond error: Unmarshal failed")
	}

	if message.Type == HEALTHCHECK {
		s.healthCheckHandler(message, nMsg.Reply)
	} else if message.Type == RESPONSE {
		s.apiHandlerResponse(message)
	}
}

func (s *service) healthCheckHandler(message *codec.Message, Reply string) {
	go func() {
		defer func() {
			if rMsg := recover(); rMsg != nil {
				recpro.Recover(s.Options().Transport.Options().Subject,"Micro->HealthCheckHandler", rMsg, message)
			}
		}()
		checkStatus, feedback := healthCheck(s.config)
		msg := codec.NewJsonResponse("", checkStatus, feedback)
		replyBody, _ := codec.Marshal(msg)
		err := s.opts.Transport.Publish(Reply ,replyBody)
		if err != nil {
			panic("ServerHandler respond error: " + err.Error())
		}
	}()
}

func (s *service) apiHandlerResponse(message *codec.Message) {
	defer func() {
		if rMsg := recover(); rMsg != nil {
			recpro.Recover(s.Options().Transport.Options().Subject, "Micro->apiHandlerResponse", rMsg, message)
		}
	}()
	ctx := s.opts.Context
	conversation := ctx.Get(message.ContextID)
	if conversation == nil {
		panic("ApiHandler respond error: conversation lost")
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

func (s *service) ServerHandlerRequest(message *codec.Message, Reply string) {
	defer func() {
		if rMsg := recover(); rMsg != nil {
			recpro.Recover(s.Options().Transport.Options().Subject, "Micro->ServerHandlerRequest", rMsg, message)
		}
	}()
	// check the message type: Request or Publish.
	// If it is a Request, the reply subject should be extracted from nats.Msg struct.
	sub := s.opts.Transport.Options().Subject
	if Reply != "" {
		message.ReplyTo = Reply
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
			panic("ServerHandler respond error: " + err.Error())
		}
	}
	reply := message.ReplyTo
	message.ReplyTo = sub

	go func() {
		defer func() {
			if rMsg := recover(); rMsg != nil {
				recpro.Recover(s.Options().Transport.Options().Subject, "Micro->RoutesHandler", rMsg, message)
				s.Respond(
					codec.NewJsonResponse(
						message.ContextID,
						500,
						nil,
					),
					reply,
				)
			}
		}()

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
				panic("ServerHandler respond error: " + err1.Error())
			}
		}
	}()
}

func (s *service) ServerHandlerResponse(message *codec.Message, data []byte)  {
	defer func() {
		if rMsg := recover(); rMsg != nil {
			recpro.Recover(s.Options().Transport.Options().Subject, "Micro->ServerHandlerResponse", rMsg, message)
		}
	}()
	conversation := s.opts.Context.Get(message.ContextID)
	if conversation == nil {
		panic("ServerHandler respond error: conversation lost")
	}
	rpl := conversation.Request
	if err := s.opts.Transport.Publish(rpl, data); err != nil {
		panic("ServerHandlerResponse Publish error: " + err.Error())
	}
	s.opts.Context.Delete(message.ContextID)
}