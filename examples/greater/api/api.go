package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/nzgogo/micro/selector"
	"github.com/nzgogo/micro/context"
	"github.com/nzgogo/micro/api"
	micro "github.com/nzgogo/micro"
	"strings"

)

type MyHandler struct {
	srv micro.Service
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("http ResponseWriter; %v\n",w)
	// map the HTTP request to internal transport request message struct.
	request, err := gogoapi.HTTPReqToNatsSReq(r)
	if err != nil {
		http.Error(w, "Cannot process request", http.StatusInternalServerError)
		return
	}
	contxt := h.srv.Options().Context
	ctxId := contxt.Add(&context.Conversation{
		Response:	w,
	})
	request.Context = ctxId

	//look up registered service in kv store
	err = h.srv.Options().Router.HttpMatch(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	srvName := micro.URLToIntnlTrans(request.Host, request.Path)
	fmt.Println("Dispatch to server: " + srvName)

	//service discovery
	slt := selector.NewSelector(selector.Registry(h.srv.Options().Registry), selector.SetStrategy(selector.RoundRobin))
	if err := slt.Init(); err != nil {
		fmt.Printf("NewSelector init failed. error: %v", err)
		http.Error(w, "Cannot process request", http.StatusInternalServerError)
		return
	}

	subj, err := slt.Select(srvName, "v1")
	if err != nil {
		fmt.Printf("Selector failed. error: %v", err)
		http.Error(w, "Cannot process request", http.StatusInternalServerError)
		return
	}
	fmt.Println("Found service: " + subj)
	subj = strings.Replace(subj, "-",".",-1)
	//transport
	natsClient := h.srv.Options().Transport
	request.ReplyTo = natsClient.Options().Subject
	c := h.srv.Options().Codec
	bytes, _ := c.Marshal(request)


	respErr := natsClient.Publish(subj,bytes)

	if respErr != nil {
		fmt.Printf("failed to send message . error: %v", err)
		http.Error(w, "No response", http.StatusInternalServerError)
		return
	}
	contxt.Wait(ctxId)
}

func main() {
	service := micro.NewService(
		"gogo-core-api",
		"v1",
	)

	service.Options().Transport.SetHandler(service.ApiHandler)

	if err := service.Init(); err != nil {

		log.Fatal(err)
	}

	go func() {
		if err := service.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	handler := MyHandler{service}
	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: &handler,
	}
	server.ListenAndServe()
}
