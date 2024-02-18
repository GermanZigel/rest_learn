package user

import (
	"encoding/json"
	"log"
	"net/http"
	Handlers "rest/internal"
	"rest/internal/userProxy"

	"github.com/julienschmidt/httprouter"
)

var _ Handlers.Handler = &Handler{}

const (
	usersUrl = "/users/v2"
	userUrl  = "/user/v2/:uuid"
)

type Handler struct {
}

func NewHandler() Handlers.Handler {
	return &Handler{}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.GET(usersUrl, h.GetList)
	router.GET(userUrl, h.GetUserByUid)
	router.POST(usersUrl, h.CreateUser)
	router.PUT(userUrl, h.UpdateUser)
	router.DELETE(userUrl, h.DeleteUser)

}

func (h *Handler) GetList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	CreatedUser := userProxy.Setter()
	log.Println("Получена структура созданого User с параметрами:", *CreatedUser)
	response, _ := json.Marshal(CreatedUser)
	w.Write(response)
}

func (h *Handler) GetUserByUid(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("This is  users"))

}
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	CreatedUser := userProxy.Setter()
	log.Println("Получена структура созданого User с параметрами:", *CreatedUser)
	response, _ := json.Marshal(CreatedUser)
	w.Write(response)
}
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("This is updated of user"))

}
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(204)
	w.Write([]byte("This is delete of user"))

}
