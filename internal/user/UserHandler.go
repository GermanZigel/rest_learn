package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	Handlers "rest/internal"
	"rest/internal/config"
	"rest/internal/logging"
	"rest/internal/user/db"
	"rest/internal/user/storage"
	"rest/internal/userProxy"
	"rest/pkg/client/pgclient"
	"rest/pkg/proto"
	"strconv"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
)

var _ Handlers.Handler = &Handler{}

type Handler struct {
	logger logging.Logger
	repo   storage.Repository
}

func NewHandler(repo storage.Repository) Handlers.Handler {
	return &Handler{
		logger: logging.GetLogger(),
		repo:   repo,
	}
}

func (h *Handler) Register(router *httprouter.Router, repo storage.Repository) {
	cfg := config.GetConfig()
	router.GET(cfg.Listen.URI_List, h.GetList)
	router.GET(cfg.Listen.URI_Once, h.GetUserByUid)
	router.POST(cfg.Listen.URI_Once, h.CreateUser)
	router.PUT(cfg.Listen.URI_Once, h.UpdateUser)
	router.DELETE(cfg.Listen.URI_Once, h.DeleteUser)
	router.DELETE(cfg.Listen.URI_List, h.DeleteUsers)

}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := logging.GetLogger()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	CreatedUser := h.CreateUserImpl()
	logger.Infof("ID=: %d", CreatedUser.Id)
	logger.Infof("Пользователь успешно создан с ID: %d", CreatedUser.Id)
	response, _ := json.Marshal(CreatedUser)
	w.Write(response)
	h.logger.Infof("Вернули ответ:%s", string(response))

}
func (h *Handler) CreateUserImpl() *userProxy.User {
	logger := logging.GetLogger()
	User := userProxy.Setter()
	h.logger.Info("Получена структура созданного User с параметрами:", &User)
	logger.Infof("ID=: %d", User.Id)
	ctx := context.Background()
	_, err := h.repo.Create(ctx, *User)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	return User

}
func (h *Handler) GetList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := logging.GetLogger()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	ListOfUsers, err := h.GetListImpl()
	if err != nil {
		logger.Errorf("Ошибка при получении списка пользователей: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	response, err := json.Marshal(ListOfUsers)
	if err != nil {
		logger.Errorf("Ошибка при маршалинге списка пользователей в JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func (h *Handler) GetListImpl() ([]userProxy.User, error) {
	logger := logging.GetLogger()
	ctx := context.Background()
	ListOfUsers, err := h.repo.GetList(ctx)
	if err != nil {
		return nil, err
	}
	logger.Infof("Список пользователей получен: %v", ListOfUsers)
	return ListOfUsers, err
}
func (h *Handler) GetUserByUid(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := logging.GetLogger()
	logger.Infof("Query params http in:%s", r.URL.Query())
	SearchIdstr := r.URL.Query().Get("id")
	SearchId, err := strconv.Atoi(SearchIdstr)
	if err != nil {
		// Обработка ошибки, если строка не является числом
		fmt.Println("Ошибка:", err)
		// Вероятно, вам также нужно отправить ошибку клиенту HTTP
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	logger.Infof("Query params pars:%s", SearchId)

	FOundUsers, err := h.repo.GetOnce(context.Background(), SearchId)
	if err == pgx.ErrNoRows {
		logger.Errorf("Пользователь не найден: %v", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	} else if err != nil && err != sql.ErrNoRows {
		logger.Errorf("Ошибка при получении пользователя: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	} else {
		logger.Infof("Found user %s", FOundUsers)
		response, _ := json.Marshal(FOundUsers)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(response)
		return
	}
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := logging.GetLogger()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// обработка ошибки чтения тела запроса
		return
	}
	var user userProxy.User
	if err := json.Unmarshal(buf, &user); err != nil {
		logger.Errorf("Ошибка при разборе JSON: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	logger.WithFields(logrus.Fields{
		"Id":      user.Id,
		"Name":    user.Name,
		"Job":     user.Job,
		"Created": user.Created,
		"Comment": user.Comment,
	}).Info("Recieved User")
	ctx := context.Background()
	updatedUser, err := h.repo.Update(ctx, user)
	logger.WithFields(logrus.Fields{
		"Id": updatedUser.Id,
	}).Info("User ID= ")
	if err != nil {
		logger.Errorf("Ошибка при получении списка пользователей: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if updatedUser.Id == 0 {
		logger.WithFields(logrus.Fields{
			"Id": user.Id,
		}).Info("Не найен пользователь с")
		http.Error(w, "user not found", http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
	} else {
		logger.WithFields(logrus.Fields{
			"Id":      updatedUser.Id,
			"Name":    updatedUser.Name,
			"Job":     updatedUser.Job,
			"Created": updatedUser.Created,
			"Comment": updatedUser.Comment,
		}).Info("Updated User")
		responseUser := struct {
			Id      int    `json:"Id"`
			Name    string `json:"name"`
			Job     string `json:"job"`
			Comment string `json:"Comment,omitempty"`
		}{
			Id:      updatedUser.Id,
			Name:    updatedUser.Name,
			Job:     updatedUser.Job,
			Comment: updatedUser.Comment,
		}

		response, err := json.Marshal(responseUser)
		if err != nil {
			logger.Errorf("Ошибка при формировании json: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(response)
	}
}
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := logging.GetLogger()
	logger.Infof("Query params http in:%s", r.URL.Query())
	SearchIdstr := r.URL.Query().Get("id")
	SearchId, err := strconv.Atoi(SearchIdstr)
	if err != nil {
		// Обработка ошибки, если строка не является числом
		fmt.Println("Ошибка:", err)
		// Вероятно, вам также нужно отправить ошибку клиенту HTTP
		http.Error(w, "Invalid ID", http.StatusBadRequest)

	}
	storage := h.repo
	logger.Infof("Query params pars:%s", SearchId)
	statusCode := h.DeleteUserLogic(SearchId, storage)
	if statusCode == http.StatusNoContent {
		http.Error(w, "User deleted", http.StatusNoContent)
	} else {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
	}

}
func (h *Handler) DeleteUserLogic(SearchId int, storage storage.Repository) int {
	logger := logging.GetLogger()
	ctx := context.Background()                      // Используйте контекст, который вы хотите
	delRes, err := storage.DeleteOnce(ctx, SearchId) // Исправлено: передаем ctx и SearchId
	if err != nil {
		logger.Fatalf("%v", err)
	}
	if delRes == true {
		return http.StatusNoContent
	} else {
		return http.StatusBadRequest
	}
}

func (h *Handler) DeleteUsers(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(204)
	w.Write([]byte("This is delete of user"))
}

type UserServiceServer struct {
}

func (UserServiceServer) GetUser(ctx context.Context, input *proto.GetUserInput) (*proto.User, error) {
	// Создаем новый экземпляр структуры proto.User
	logger := logging.GetLogger()
	cfg := config.GetConfig()
	logger.WithFields(logrus.Fields{
		"Id": input.Id,
	}).Info("Input req")
	var GrpcSearchId int
	GrpcSearchId = int(input.Id)
	pgsClient, err := pgclient.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	poolClient := pgsClient.(*pgclient.PgxPoolClient)
	pool := poolClient.Pool                                             // Получить подлежащий пул
	Repository := db.NewRepository(context.Background(), pool, &logger) // Передача пула и логгера
	FOundUsers, err := Repository.GetOnce(context.Background(), GrpcSearchId)
	GRPCFOundUsers := new(proto.User)

	GRPCFOundUsers.Id = int32(FOundUsers.Id)
	GRPCFOundUsers.Name = FOundUsers.Name
	GRPCFOundUsers.Job = FOundUsers.Job
	createdTimestamp := timestamppb.New(FOundUsers.Created)
	GRPCFOundUsers.Created = createdTimestamp

	logger.WithFields(logrus.Fields{
		"Id":      GRPCFOundUsers.Id,
		"Name":    GRPCFOundUsers.Name,
		"Job":     GRPCFOundUsers.Job,
		"Created": FOundUsers.Created,
	}).Info("Found User")
	return GRPCFOundUsers, nil
}
