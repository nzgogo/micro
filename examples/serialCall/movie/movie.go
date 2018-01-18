package main

import (
	"fmt"
	"log"
	"time"
	"errors"

	"github.com/nzgogo/micro"
	"github.com/nzgogo/micro/codec"
	"github.com/nzgogo/micro/router"
	"github.com/nzgogo/micro/db"
	"github.com/jinzhu/gorm"

)

var(
	SrvCastCastHandler = "/movie_cast"
	SrvCast = "gogo-core-crew"
	ErrQueryFailure = errors.New("Query Faileld")
)

type Movies struct {
	gorm.Model
	Name string
	Director string
	Budget   int
	Producer string
	InitRlease time.Time
}

type server struct {
	srv gogo.Service
	movieDB db.DB
}



func (s *server) GetMoiveInfo(req *codec.Message) error {
	config := s.srv.Options()
	db := s.movieDB.DB()

	if len(req.Query["movie"]) == 0 {
		fmt.Printf("Query failed. \n")
		s.errHandler(req, ErrQueryFailure)
		return ErrQueryFailure
	}
	movie := Movies{}
	//search in database
	db.Where(&Movies{Name: req.Query["movie"][0]}).Find(&movie)
	fmt.Println(movie)
	req.Body = fmt.Sprint(movie.ID)

	//service discovery
	rpy := req.ReplyTo
	req.ReplyTo = "nats-request"
	req.Node = SrvCastCastHandler
	subj, err := config.Selector.Select(SrvCast, "v1")
	if err != nil {
		fmt.Printf("Selector failed. error: %v\n", err)
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
		message.Body = "The cast of movie " + movie.Name + " includes " + message.Body
		resp1, err := codec.Marshal(message)
		if err != nil {
			return err
		}
		return config.Transport.Publish(rpy, resp1)
	})
}

func (s *server) Cast(req *codec.Message) error {
	config := s.srv.Options()
	db := s.movieDB.DB()

	if len(req.Query["movie"]) == 0 {
		fmt.Printf("Query failed. \n")
		s.errHandler(req, ErrQueryFailure)
		return ErrQueryFailure
	}
	movie := Movies{}
	//search in database
	db.Where(&Movies{Name: req.Query["movie"][0]}).Find(&movie)
	fmt.Println(movie)
	req.Body = fmt.Sprint(movie.ID)

	//service discovery
	req.ReplyTo = config.Transport.Options().Subject
	req.Node = SrvCastCastHandler
	subj, err := config.Selector.Select(SrvCast, "v1")
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
	server.movieDB = db.NewDB("kai","qiekai1234","mydb")
	if err := server.movieDB.Connect(); err!=nil {
		log.Fatal(err)
	}
	defer server.movieDB.Close()

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
		Handler: server.GetMoiveInfo,
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
