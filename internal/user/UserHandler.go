package user

import (
	"context"
	"encoding/json"
	"net/http"
	Handlers "rest/internal"
	"rest/internal/config"
	"rest/internal/logging"
	"rest/internal/user/db"
	"rest/internal/userProxy"
	"rest/pkg/client/pgclient"

	"github.com/julienschmidt/httprouter"
)

var _ Handlers.Handler = &Handler{}

type Handler struct {
	logger logging.Logger
}

func NewHandler() Handlers.Handler {
	return &Handler{
		logger: logging.GetLogger(),
	}
}

func (h *Handler) Register(router *httprouter.Router) {
	cfg := config.GetConfig()
	router.GET(cfg.Listen.URI_List, h.GetList)
	router.GET(cfg.Listen.URI_Once, h.GetUserByUid)
	router.POST(cfg.Listen.URI_Once, h.CreateUser)
	router.PUT(cfg.Listen.URI_Once, h.UpdateUser)
	router.DELETE(cfg.Listen.URI_Once, h.DeleteUser)

}
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := logging.GetLogger()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	CreatedUser := userProxy.Setter()
	cfg := config.GetConfig()
	pgsClient, err := pgclient.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("%v", err)
	}

	poolClient := pgsClient.(*pgclient.PgxPoolClient)
	pool := poolClient.Pool // Получить подлежащий пул

	Repository := db.NewRepository(pool, &logger) // Передача пула и логгера

	h.logger.Info("Получена структура созданного User с параметрами:", *CreatedUser)
	response, _ := json.Marshal(CreatedUser)
	w.Write(response)
	h.logger.Infof("Вернули ответ:%s", string(response))

	// Использование Repository
	id, err := Repository.Create(context.Background(), *CreatedUser)
	if err != nil {
		logger.Errorf("Ошибка при создании пользователя: %v", err)
		// Обработка ошибки
		return
	}
	logger.Infof("Пользователь успешно создан с ID: %s", id)

}
func (h *Handler) GetList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := logging.GetLogger()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	CreatedUser := userProxy.Setter()
	logger.Info("Получена структура созданного User с параметрами:", *CreatedUser)
	response, err := json.Marshal(CreatedUser)
	if err != nil {
		panic(err)
	}
	w.Write(response)
}

func (h *Handler) GetUserByUid(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("This is  users"))

}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("This is updated of user"))

}
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(204)
	w.Write([]byte("This is delete of user"))

}
