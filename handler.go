package gogo

import (
	"github.com/nats-io/go-nats"
	"github.com/nzgogo/micro/api"
	"github.com/nzgogo/micro/codec"
)

func (s *service) ServerHandler(nMsg *nats.Msg) {
	message := &codec.Message{}
	codec.Unmarshal(nMsg.Data, message)
	if message.Type == "request" {
		//message.ReplyTo = s.name + "." + s.version + "." + s.id

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
		rpl := s.opts.Context.Get(message.Context).Request
		s.opts.Transport.Publish(rpl, nMsg.Data)
		s.opts.Context.Delete(message.Context)
	}
}

//Example MsgHandler
func (s *service) ApiHandler(nMsg *nats.Msg) {
	message := &codec.Message{}
	codec.Unmarshal(nMsg.Data, message)
	ctx := s.opts.Context

	r := ctx.Get(message.Context).Response

	gogoapi.WriteResponse(r, message)

	ctx.Done(message.Context)
	ctx.Delete(message.Context)

}
