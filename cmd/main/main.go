package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"rest/internal/config"
	"rest/internal/logging"
	"rest/internal/user"
	"rest/pkg/proto"
	"runtime"
	"sync"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Create router")
	logger.Infof("MaxProcs: %d", runtime.GOMAXPROCS(-1))
	router := httprouter.New()
	cfg := config.GetConfig()
	handler := user.NewHandler()
	handler.Register(router)
	start(router, logger, cfg)
}

func start(router *httprouter.Router, logger logging.Logger, cfg *config.Config) {
	addressHTTP := fmt.Sprintf("localhost:%s", cfg.Listen.HttpPort)
	addressGRPC := fmt.Sprintf("localhost:%s", cfg.Listen.GrpcPort)

	logger.Info("Start app")
	var listenerGrpc net.Listener
	var listenerHTTP net.Listener

	var ListenErr error
	var wg sync.WaitGroup
	wg.Add(2)
	startHTTP := func() {
		defer wg.Done()
		listenerHTTP, ListenErr = net.Listen("tcp", addressHTTP)
		logger.Info("Start HTTP")

		server := &http.Server{
			Handler:      router,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		server.Serve(listenerHTTP)
	}
	go func() {
		startHTTP()
	}()
	startGRPC := func() {
		defer wg.Done()
		listenerGrpc, ListenErr = net.Listen("tcp", addressGRPC)
		logger.Info("Start GRPC")

		if ListenErr != nil {
			panic(ListenErr)
		}
		if listenerGrpc == nil {
			logger.Fatal("Listener is nil, cannot start the server")
		}
		server := grpc.NewServer()
		proto.RegisterUserRPCServer(server, &user.UserServiceServer{})
		err := server.Serve(listenerGrpc)
		if err != nil {
			logger.Fatalf("Failed to start the server: %v", err)
		}
		logger.Infof("Server is now listening on %s", listenerGrpc.Addr())

	}
	go func() {
		startGRPC()
	}()
	wg.Wait()

}
