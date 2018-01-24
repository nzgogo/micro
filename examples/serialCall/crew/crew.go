package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/nzgogo/micro"
	"github.com/nzgogo/micro/codec"
	"github.com/nzgogo/micro/db"
	"github.com/nzgogo/micro/router"
)

type server struct {
	srv    gogo.Service
	castDB db.DB
}

type Casts struct {
	gorm.Model
	Name    string
	Role    string
	MovieId uint
}

func (s *server) Cast(req *codec.Message, reply string) error {
	fmt.Println("Message received: " + req.Body)

	db := s.castDB.DB()
	movieid, _ := strconv.ParseUint(req.Body, 10, 32)
	casts := []Casts{}
	db.Where(&Casts{MovieId: uint(movieid)}).Find(&casts)
	fmt.Println(casts)

	castlist := ""
	for _, cast := range casts {
		castlist = castlist + cast.Name + ". "
	}

	response := codec.NewResponse(200, req.ContextID, &castlist, req.Header)

	return s.srv.Respond(response, reply)
}

func main() {
	server := server{}
	server.castDB = db.NewDB("kai", "gogo1234", "mydb")
	if err := server.castDB.Connect(); err != nil {
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
