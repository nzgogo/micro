package main

import (
	"fmt"
	"log"

	"github.com/nzgogo/micro"
	"github.com/nzgogo/micro/codec"
	"github.com/nzgogo/micro/router"
	"github.com/nzgogo/micro/selector"
)

type server struct {
	srv gogo.Service
}

func (s *server) Moive(req *codec.Message) error {
	config := s.srv.Options()

	req.Node = "/movie_cast"
	//req.Body = ""

	//service discovery
	slt := selector.NewSelector(
		selector.Registry(config.Registry),
		selector.SetStrategy(selector.RoundRobin),
	)
	if err := slt.Init(); err != nil {
		fmt.Printf("NewSelector init failed. error: %v", err)

	}
	rpy := req.ReplyTo
	req.ReplyTo = "nats-request"

	subj, err := slt.Select("gogo-core-crew", "v1")
	if err != nil {
		fmt.Printf("Selector failed. error: %v", err)
		s.errHandler(req, err)
	}
	fmt.Println("Found service: " + subj)
	resp, err := codec.Marshal(req)
	if err != nil {
		return err
	}
	return config.Transport.Request(subj,resp , func(bytes []byte) error {
		message := &codec.Message{}
		codec.Unmarshal(bytes, message)
		message.Body = "The cast of movie 'Legend of The Demon Cat' includes " + message.Body
		resp1, err := codec.Marshal(message)
		if err != nil {
			return err
		}
		return config.Transport.Publish(rpy, resp1)
	})
}

func (s *server) Cast(req *codec.Message) error {
	config := s.srv.Options()

	req.Node = "/movie_cast"
	//req.Body = ""

	//service discovery
	slt := selector.NewSelector(
		selector.Registry(config.Registry),
		selector.SetStrategy(selector.RoundRobin),
	)
	if err := slt.Init(); err != nil {
		fmt.Printf("NewSelector init failed. error: %v", err)

	}
	//rpy := req.ReplyTo
	req.ReplyTo = config.Transport.Options().Subject

	subj, err := slt.Select("gogo-core-crew", "v1")
	if err != nil {
		fmt.Printf("Selector failed. error: %v", err)
		s.errHandler(req, err)
	}
	fmt.Println("Found service: " + subj)

	resp, err := codec.Marshal(req)
	if err != nil {
		return err
	}
	return config.Transport.Publish(subj, resp)
}

func (s *server)errHandler(req *codec.Message, err error){
	response := &codec.Message{
		Type:	    "response",
		StatusCode: 500,
		Header:     make(map[string][]string, 0),
		Body:       err.Error(),
		ContextID:    req.ContextID,
	}
	resp, _ := codec.Marshal(response)
	s.srv.Options().Transport.Publish(req.ReplyTo, resp)
}

func main() {
	server := server{}
	service := gogo.NewService(
		"gogo-core-movie",
		"v1",
	)

	server.srv = service

	if err := server.srv.Init(); err != nil {
		log.Fatal(err)
	}

	server.srv.Options().Transport.SetHandler(service.ServerHandler)

	r := server.srv.Options().Router

	r.Add(&router.Node{
		Method:  "GET",
		Path:    "/movie",
		ID:      "/movie",
		Handler: server.Moive,
	})
	r.Add(&router.Node{
		Method:  "GET",
		Path:    "/cast",
		ID:      "/cast",
		Handler: server.Cast,
	})

	// Run server
	if err := server.srv.Run(); err != nil {
		log.Fatal(err)
	}
}
