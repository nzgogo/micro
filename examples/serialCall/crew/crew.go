package main

import (
	"fmt"
	"log"

	"github.com/nzgogo/micro"
	"github.com/nzgogo/micro/codec"
	"github.com/nzgogo/micro/router"
)

type server struct {
	srv gogo.Service
}

func (s *server) Cast(req *codec.Message) error {
	response := &codec.Message{
		Type:	    "response",
		StatusCode: 200,
		Header:     make(map[string][]string, 0),
		Body:       "Hunagxuan",
		ContextID:    req.ContextID,
	}
	resp, err := codec.Marshal(response)
	if err != nil {
		return err
	}
	fmt.Println("Message received: " + req.ReplyTo)
	return s.srv.Options().Transport.Publish(req.ReplyTo, resp)
}

func main() {
	server := server{}
	service := gogo.NewService(
		"gogo-core-crew",
		"v1",
	)

	server.srv = service

	if err := server.srv.Init(); err != nil {
		log.Fatal(err)
	}

	server.srv.Options().Transport.SetHandler(service.ServerHandler)

	r := server.srv.Options().Router

	r.Add(&router.Node{
		ID:      "/movie_cast",
		Handler: server.Cast,
	})


	// Run server
	if err := server.srv.Run(); err != nil {
		log.Fatal(err)
	}
}
