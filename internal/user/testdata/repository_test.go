package testdata

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"rest/internal/logging"
	"rest/internal/user"
	"rest/internal/userProxy"
	"strconv"
	"testing"
)

type testRepository struct {
	logger *logging.Logger
}

func (r *testRepository) Create(ctx context.Context, user userProxy.User) (string, error) {
	// Генерация Id для пользователя
	userId := strconv.Itoa(user.Id)
	return userId, nil
}
func (r *testRepository) GetList(ctx context.Context) ([]userProxy.User, error) {
	var users []userProxy.User
	users = append(users, userProxy.User{})
	return users, nil
}
func (r *testRepository) GetOnce(ctx context.Context, id int) (userProxy.User, error) {
	usr := userProxy.User{}
	return usr, nil
}

func (r *testRepository) DeleteOnce(ctx context.Context, id int) (bool, error) {

	return true, nil
}
func (r *testRepository) Update(ctx context.Context, u userProxy.User) (userProxy.User, error) {
	var updatedUser userProxy.User
	updatedUser = u
	return updatedUser, nil
}
func TestDelete(t *testing.T) {
	storage := &testRepository{}
	handler := user.NewHandler(storage)
	router := httprouter.New()
	handler.Register(router, storage)
	req, err := http.NewRequest("DELETE", "/user/v3?id=7", nil)
	assert.NoError(t, err)

	// Создаем ResponseRecorder для получения ответа
	rr := httptest.NewRecorder()

	// Выполняем запрос
	router.ServeHTTP(rr, req)

	// Проверяем статус-код ответа
	assert.Equal(t, http.StatusNoContent, rr.Code)

}
