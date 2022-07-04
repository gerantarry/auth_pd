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
	Insert(ctx context.Context, user entity.User) error
	Delete(ctx context.Context, id int, login string)
	Update(ctx context.Context, id int, login string)
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

func (s *userStorage) Insert(ctx context.Context, user entity.User) error {
	checkPing(s)
	_, err := s.storage.Exec(
		"insert pd.person(first_name, second_name, login, password_hash, email) values(%s,%s,%s,%s,%s)",
		user.FirstName, user.SecondName, user.Login, user.Password, user.Email)
	return err
}

func checkPing(s *userStorage) {
	if err := s.storage.Ping(); err != nil {
		fmt.Println(err.Error())
	}
}

func checkUserPassword(user entity.User, password string) bool {
	return user.Password == password
}
