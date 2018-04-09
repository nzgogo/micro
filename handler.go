package gogo

import (
	"github.com/nats-io/go-nats"
	"github.com/nzgogo/micro/codec"
	"github.com/nzgogo/micro/constant"
	"github.com/nzgogo/micro/context"
	recpro "github.com/nzgogo/micro/recover"
)

func (s *service) ServerHandler(nMsg *nats.Msg) {
	defer recpro.Recover(s.config[constant.SLACKCHANNELADDR], s.Options().Transport.Options().Subject, "ServerHandler", nMsg)

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
	if message.Type == constant.REQUEST {
		s.serverHandlerRequest(message, nMsg.Reply)
	} else if message.Type == constant.HEALTHCHECK {
		s.healthCheckHandler(message, nMsg.Reply)
	} else if message.Type == constant.RESPONSE {
		s.serverHandlerResponse(message, nMsg.Data)
	} else if message.Type == constant.PUBLISH {
		s.serverHandlerPublish(message, nMsg.Reply)
	}

}

func (s *service) ApiHandler(nMsg *nats.Msg) {
	defer recpro.Recover(s.config[constant.SLACKCHANNELADDR], s.Options().Transport.Options().Subject, "ApiHandler", nMsg)

	if nMsg == nil {
		panic("Nats body empty")
	}

	message := &codec.Message{}
	uerr := codec.Unmarshal(nMsg.Data, message)
	if uerr != nil {
		panic("ApiHandler respond error: Unmarshal failed")
	}

	if message.Type == constant.HEALTHCHECK {
		s.healthCheckHandler(message, nMsg.Reply)
	} else if message.Type == constant.RESPONSE {
		s.apiHandlerResponse(message)
	}
}

func (s *service) healthCheckHandler(message *codec.Message, Reply string) {
	go func() {
		defer recpro.Recover(s.config[constant.SLACKCHANNELADDR], s.Options().Transport.Options().Subject, "Micro->HealthCheckHandler", message)
		checkStatus, feedback := healthCheck(s.config)
		msg := codec.NewJsonResponse("", checkStatus, feedback)
		replyBody, _ := codec.Marshal(msg)
		err := s.opts.Transport.Publish(Reply, replyBody)
		if err != nil {
			panic("ServerHandler respond error: " + err.Error())
		}
	}()
}

func (s *service) apiHandlerResponse(message *codec.Message) {
	defer recpro.Recover(s.config[constant.SLACKCHANNELADDR], s.Options().Transport.Options().Subject, "Micro->apiHandlerResponse", message)
	ctx := s.opts.Context
	conversation := ctx.Get(message.ContextID)
	if conversation == nil {
		panic("ApiHandler respond error: conversation lost")
	}
	r := conversation.Response

	fn := message.WriteHTTPResponse
	for i := len(s.opts.HttpRespWrappers); i > 0; i-- {
		fn = s.opts.HttpRespWrappers[i-1](fn)
	}
	fn(r)
	ctx.Done(message.ContextID)
	ctx.Delete(message.ContextID)
}

func (s *service) serverHandlerRequest(message *codec.Message, Reply string) {
	defer recpro.Recover(s.config[constant.SLACKCHANNELADDR], s.Options().Transport.Options().Subject, "Micro->ServerHandlerRequest", message)
	// check the nats message type: Request or Publish.
	// If it is a nats Request, the reply subject should be extracted from nats.Msg struct.
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
				recpro.PostProc(s.config[constant.SLACKCHANNELADDR], s.Options().Transport.Options().Subject, "Micro->RoutesHandler", rMsg, message)
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
			body := map[string]interface{}{"message": err.Message}
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

func (s *service) serverHandlerPublish(message *codec.Message, Reply string) {
	defer recpro.Recover(s.config[constant.SLACKCHANNELADDR], s.Options().Transport.Options().Subject, "Micro->ServerHandlerPublish", message)
	sub := s.opts.Transport.Options().Subject
	if Reply != "" {
		message.ReplyTo = Reply
	}

	handler, routerErr := s.opts.Router.Dispatch(message)
	if routerErr != nil {
		errResp := codec.NewResponse(message.ContextID, 404, nil, message.Header)
		err := s.Respond(errResp, message.ReplyTo)
		if err != nil {
			panic("ServerHandler respond error: " + err.Error())
		}
	}
	message.ReplyTo = sub

	go func() {
		defer recpro.Recover(s.config[constant.SLACKCHANNELADDR], s.Options().Transport.Options().Subject, "Micro->RoutesHandler", message)

		for i := len(s.opts.HdlrWrappers); i > 0; i-- {
			handler = s.opts.HdlrWrappers[i-1](handler)
		}

		err := handler(message, "")
		if err != nil {
			panic("ServerHandler error: " + err.Message)
		}
	}()
}

func (s *service) serverHandlerResponse(message *codec.Message, data []byte) {
	defer recpro.Recover(s.config[constant.SLACKCHANNELADDR], s.Options().Transport.Options().Subject, "Micro->ServerHandlerResponse", message)
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
