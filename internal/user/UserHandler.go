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
	//var Id string
	_, err = Repository.Create(context.Background(), *CreatedUser)
	logger.Infof("ID=: %d", CreatedUser.Id)
	if err != nil {
		logger.Errorf("Ошибка при создании пользователя: %v", err)
		// Обработка ошибки
		return
	}
	logger.Infof("Пользователь успешно создан с ID: %d", CreatedUser.Id)

}
func (h *Handler) GetList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := logging.GetLogger()
	cfg := config.GetConfig()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	pgsClient, err := pgclient.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	poolClient := pgsClient.(*pgclient.PgxPoolClient)
	pool := poolClient.Pool                       // Получить подлежащий пул
	Repository := db.NewRepository(pool, &logger) // Передача пула и логгера

	// Использование Repository
	ListOfUsers, err := Repository.GetList(context.Background())
	if err != nil {
		logger.Errorf("Ошибка при получении списка пользователей: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	logger.Infof("Список пользователей получен: %v", ListOfUsers)

	response, err := json.Marshal(ListOfUsers)
	if err != nil {
		logger.Errorf("Ошибка при маршалинге списка пользователей в JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
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
