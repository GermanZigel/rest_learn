package Handlers

import (
	"github.com/julienschmidt/httprouter"
	"rest/internal/user/storage"
)

type Handler interface {
	Register(router *httprouter.Router, repo storage.Repository)
}
