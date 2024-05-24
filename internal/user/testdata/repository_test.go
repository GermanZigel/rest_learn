package testdata

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"rest/internal/config"
	"rest/internal/logging"
	"rest/internal/user"
	"rest/internal/userProxy"
	"strconv"
	"testing"
	"time"

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
	// Определение количества пользователей для генерации
	numUsers := 5

	// Заполняем срез пользователей случайными данными
	users := fillUsersSlice(numUsers)

	return users, nil
}

// Функция для заполнения среза структур User случайными данными
func fillUsersSlice(size int) []userProxy.User {
	rand.Seed(time.Now().UnixNano())

	users := make([]userProxy.User, size)
	for i := range users {
		users[i] = userProxy.User{
			Id:      rand.Intn(10000), // Случайное число от 0 до 9999
			Name:    fmt.Sprintf("User%d", i),
			Job:     fmt.Sprintf("Job%d", i),
			Created: time.Now(), // Текущее время
			Comment: fmt.Sprintf("Comment%d", i),
		}
	}

	return users
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

func TestCreateUserLogic(t *testing.T) {
	storage := &testRepository{}
	handler := user.NewHandler(storage)
	logger := logging.GetLogger()
	cfg := config.GetConfig()
	// Вызываем DeleteUserLogic с каким-то id
	CreatedUser := handler.(*user.Handler).CreateUserImpl()

	logger.Infof("Созданый ID=%d, id из конфигурации = %d", CreatedUser.Id, cfg.User.MinId)
	// Проверяем, что полученный код ответа соответствует ожидаемому
	assert.GreaterOrEqual(t, CreatedUser.Id, cfg.User.MinId)

}
func TestGetListImpl(t *testing.T) {
	storage := &testRepository{}
	handler := user.NewHandler(storage)
	logger := logging.GetLogger()
	//cfg := config.GetConfig()
	// Вызываем DeleteUserLogic с каким-то id
	CreatedUser, err := handler.(*user.Handler).GetListImpl()
	if err != nil {
		panic(err)
	}
	for _, user := range CreatedUser {
		logger.Infof("получено: %+v", user)
	}
	assert.NotEmpty(t, CreatedUser)

}
