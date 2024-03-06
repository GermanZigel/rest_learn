package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"rest/internal/config"
	"rest/internal/logging"
	"rest/internal/user"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Create router")
	router := httprouter.New()
	cfg := config.GetConfig()
	handler := user.NewHandler()
	handler.Register(router)
	start(router, logger, cfg)
}

func start(router *httprouter.Router, logger logging.Logger, cfg *config.Config) {
	logger.Info("Start app")
	var listener net.Listener
	var ListenErr error
	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketpath := path.Join(appDir, "app.sock")
		logger.Debugf("socket path: %s", socketpath)
		logger.Info("create unix socket")
		listener, ListenErr = net.Listen("unix", socketpath)
		if ListenErr != nil {
			logger.Fatal(ListenErr)
		}
	} else {
		logger.Info("listen tcp")
		listener, ListenErr = net.Listen("tcp", fmt.Sprintf(":%s", cfg.Listen.Port))
	}

	if ListenErr != nil {
		panic(ListenErr)
	}
	if listener == nil {
		logger.Fatal("Listener is nil, cannot start the server")
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	server.Serve(listener)
	log.Println("server is listening port:", listener.Addr())
	log.Fatal(server.Serve(listener))
}
