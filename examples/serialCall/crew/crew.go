package main

import (
	"fmt"
	"log"

	"github.com/nzgogo/micro"
	"github.com/nzgogo/micro/codec"
	"github.com/nzgogo/micro/router"
	"github.com/nzgogo/micro/db"
	"github.com/jinzhu/gorm"
	"strconv"
	"micro/api"
)

type server struct {
	srv gogo.Service
	castDB db.DB
}

type Casts struct {
	gorm.Model
	Name string
	Role string
	MovieId uint
}

func (s *server) Cast(req *codec.Message) error {
	fmt.Println("Message received: " + req.ContextID)

	db := s.castDB.DB()
	movieid, _ := strconv.ParseUint(req.Body, 10, 32)
	casts := []Casts{}
	db.Where(&Casts{MovieId: uint(movieid)}).Find(&casts)
	fmt.Println(casts)

	castlist := ""
	for _,cast:=range casts {
		castlist = castlist + cast.Name + ". "
	}

	response := gogoapi.NewResponse(200, "response", req.ContextID, &castlist, req.Header)

	resp, err := codec.Marshal(response)
	if err != nil {
		return err
	}
	return s.srv.Options().Transport.Publish(req.ReplyTo, resp)
}

func main() {
	server := server{}
	server.castDB = db.NewDB("kai","qiekai1234","mydb")
	if err := server.castDB.Connect(); err!=nil {
		log.Fatal(err)
	}
	defer server.castDB.Close()

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
