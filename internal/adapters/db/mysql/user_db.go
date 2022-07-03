package mysql

import (
	"auth_pd/internal/domain/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

const (
	mysqlDriver    = "mysql"
	dataSourceName = "root:!QAZ2wsx#EDC@/schema" //TODO заменить хардкод на значение из конфиг файла
)

type Storage interface {
	Get(ctx context.Context, login string, password string) (*entity.User, error)
	Insert(ctx context.Context, user entity.User)
	Delete()
	Update()
}

type userStorage struct {
	storage sql.DB
}

//NewUserStorage - конструктор
func NewUserStorage(db *sql.DB) *userStorage {
	return &userStorage{storage: *db}
}

func (s *userStorage) Get(ctx context.Context, login string, password string) (*entity.User, error) {
	checkPing(s)
	row := s.storage.QueryRow("select * from pd.person where login = %s", login) //TODO заменить хардкод названия схемы на значение из конфиг файла
	user := entity.User{}
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.SecondName,
		&user.Login,
		&user.Password,
		&user.Email)
	if err != nil {
		fmt.Println(err.Error())
		return &user, err
	}
	switch checkUserPassword(user, password) {
	case true:
		return &user, nil
	case false:
		return &user, errors.New("wrong password")
	}
	return &user, nil
}

func checkPing(s *userStorage) {
	if err := s.storage.Ping(); err != nil {
		fmt.Println(err.Error())
	}
}

func checkUserPassword(user entity.User, password string) bool {
	return user.Password == password
}
