package main

import (
	"log"
	"net"
	"net/http"
	"rest/internal/user"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	log.Println("Create router")
	router := httprouter.New()
	handler := user.NewHandler()
	handler.Register(router)
	start(router)

}
func start(router *httprouter.Router) {
	log.Println("Start app")
	listenener, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	server.Serve(listenener)
	log.Println("server is listening port:", listenener.Addr())
	log.Fatal(server.Serve(listenener))
}
