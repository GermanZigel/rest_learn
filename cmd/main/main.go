package main

import (
	"rest/internal/config"
	"rest/internal/logging"
	"rest/internal/user"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Create router")
	router := httprouter.New()
	cfg := config.GetConfig()
	handler := user.NewHandler()
	handler.Register(router)
	user.Start(router, logger, cfg)
}
