package testdata

import (
	"context"
	"net/http"
	"rest/internal/logging"
	"rest/internal/user"
	"rest/internal/userProxy"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestDeleteUserLogic(t *testing.T) {
	storage := &testRepository{}
	handler := user.NewHandler(storage)

	// Вызываем DeleteUserLogic с каким-то id
	statusCode := handler.(*user.Handler).DeleteUserLogic(1, storage)

	// Проверяем, что полученный код ответа соответствует ожидаемому
	assert.Equal(t, http.StatusNoContent, statusCode)
}
