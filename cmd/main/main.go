package main

import (
	"log"
	"net"
	"net/http"
	"rest/internal/logging"
	"rest/internal/user"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Create router")
	router := httprouter.New()
	handler := user.NewHandler()
	handler.Register(router)
	start(router, logger)

}
func start(router *httprouter.Router, logger logging.Logger) {
	logger.Info("Start app")
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
