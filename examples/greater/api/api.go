package main

import (
	"log"
	"fmt"
	"net/http"
	micro "micro"
	"micro/api"
	"micro/selector"
	"micro/router"
)

type MyHandler struct{
	srv micro.Service
}
func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// map the HTTP request to internal transport request struct.
	request, err := gogoapi.HTTPReqToNatsSReq(r)
	if err != nil {
		http.Error(w, "Cannot process request", http.StatusInternalServerError)
		return
	}

	//look up registered service in kv store
	err = h.srv.Options().Router.HttpMatch(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	var response = gogoapi.NewResponse()

	srvName := micro.URLToIntnlTrans(request.Host,request.Path)
	fmt.Println("Dispatch to server: "+srvName)

	//service discovery
	slt := selector.NewSelector(selector.Registry(h.srv.Options().Registry), selector.SetStrategy(selector.RoundRobin))
	if err := slt.Init(); err != nil {
		fmt.Printf("NewSelector init failed. error: %v", err)
		http.Error(w, "Cannot process request", http.StatusInternalServerError)
		return
	}

	subj, err := slt.Select(srvName,"v1")
	if err != nil {
		fmt.Printf("Selector failed. error: %v", err)
		http.Error(w, "Cannot process request", http.StatusInternalServerError)
		return
	}
	fmt.Println("Found service: "+subj)

	//transport
	natsClient := h.srv.Options().Transport
	c := h.srv.Options().Codec
	bytes,_ := c.Marshal(request)
	respErr := natsClient.Request(subj, bytes, func(bytes []byte) error {
		return c.Unmarshal(bytes, response)
	})

	if respErr != nil {
		fmt.Printf("Get response failed. error: %v", err)
		http.Error(w, "No response", http.StatusInternalServerError)
		return
	}

	//write response to http
	gogoapi.WriteResponse(w, response)
}

func main() {
	route := router.NewRouter(router.Name("gogox/v1/api"))

	service := micro.NewService(
		"gogox.core.api",
		"v1",
	)

	if err := service.Init(micro.Router(route)); err != nil{
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