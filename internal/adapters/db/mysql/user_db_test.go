package mysql

import (
	"auth_pd/internal/domain/entity"
	"auth_pd/pkg/logging"
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	login    = "user1_login"
	password = "so61ty872nst"
	name     = "Carla"
)

var logger *logging.Logger

func createDbClient() *sql.DB {
	db, _ := sql.Open(mysqlDriver, dataSourceName)
	return db
}

func createUser() *entity.User {
	return &entity.User{
		ID:         0,
		FirstName:  "Alexa",
		SecondName: "Test",
		Login:      "test_login",
		Password:   "test_pass",
		Email:      "test@mail.ru",
	}
}

func init() {
	//TODO путь должен браться из конфигов
	err := os.Setenv("PROJECT_DIR", "C:\\Users\\Anton\\GolandProjects\\auth_pd")
	if err != nil {
		panic(any(err))
	}
	logger = logging.GetLogger()
}

func TestNewUserStorage(t *testing.T) {
	db := createDbClient()
	assert.NotEmpty(t, db, "ошибка при создании db клиента")

	storage := NewUserStorage(db, logger)
	assert.NotEmpty(t, storage, "хранилище не инициализировано")
}

func TestUserStorage_Get(t *testing.T) {
	storage := NewUserStorage(createDbClient(), logger)
	user, err := storage.Get(context.Background(), login, password)
	require.Nil(t, err)
	assert.EqualValues(t, name, user.FirstName, "данные в записи не совпадают")
}

// покрывает сразу 2 теста для Insert и Delete
func TestUserStorage_Insert(t *testing.T) {
	storage := NewUserStorage(createDbClient(), logger)
	user := createUser()
	errI := storage.Insert(context.Background(), *user)
	require.Nil(t, errI)
	getUser, errG := storage.Get(context.Background(), user.Login, user.Password)
	require.Nil(t, errG)
	assert.EqualValues(t, user.Login, getUser.Login)
	assert.EqualValues(t, user.Password, getUser.Password)
	assert.EqualValues(t, user.Email, getUser.Email)
	assert.EqualValues(t, user.FirstName, getUser.FirstName)
	assert.EqualValues(t, user.SecondName, getUser.SecondName)

	errD := storage.Delete(context.Background(), getUser.ID, getUser.Login)
	require.Nil(t, errD)
}
