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
	dataSourceName = "root:!QAZ2wsx#EDC@tcp(127.0.0.1:6603)/pd" //TODO заменить хардкод на значение из конфиг файла
)

type Storage interface {
	Get(ctx context.Context, login, password string) (*entity.User, error)
	Insert(ctx context.Context, user entity.User) error
	Delete(ctx context.Context, id int, login string) error
	Update(ctx context.Context, id int, login string) error
}

type userStorage struct {
	storage sql.DB
}

//NewUserStorage - конструктор
func NewUserStorage(db *sql.DB) *userStorage {
	return &userStorage{storage: *db}
}

func (s *userStorage) Get(ctx context.Context, login, password string) (*entity.User, error) {
	checkPing(ctx, s)
	row := s.storage.QueryRowContext(ctx, "select * from pd.person where login = ?", login) //TODO заменить хардкод названия схемы на значение из конфиг файла
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
	checkPing(ctx, s)
	_, err := s.storage.ExecContext(
		ctx,
		"insert pd.person(first_name, second_name, login, password_hash, email) values(?, ?, ?, ?, ?)", //TODO заменить хардкод названия схемы на значение из конфиг файла
		user.FirstName, user.SecondName, user.Login, user.Password, user.Email)
	return err
}

func (s *userStorage) Delete(ctx context.Context, id int, login string) error {
	checkPing(ctx, s)
	_, err := s.storage.ExecContext(ctx,
		"delete from pd.person where id = ? and login = ?",
		id,
		login) //TODO заменить хардкод названия схемы на значение из конфиг файла
	return err
}

//Update не доделана - необходим способ передачи параметра который собираемся апдейтить
func (s *userStorage) Update(ctx context.Context, id int, login string) error {
	checkPing(ctx, s)
	_, err := s.storage.ExecContext(ctx,
		"update pd.person",
		id,
		login)
	return err
}

func checkPing(ctx context.Context, s *userStorage) {
	if err := s.storage.PingContext(ctx); err != nil {
		fmt.Println(err.Error())
	}
}

func checkUserPassword(user entity.User, password string) bool {
	return user.Password == password
}
