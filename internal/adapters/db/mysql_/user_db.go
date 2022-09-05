package mysql_

import (
	"auth_pd/internal/domain/entity"
	"auth_pd/pkg/logging"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

const DriverMySQL = "mysql"

type Storage interface {
	GetUser(ctx context.Context, login, password string) (*entity.User, error)
	Insert(ctx context.Context, user entity.User) error
	Delete(ctx context.Context, id int, login string) error
	Update(ctx context.Context, id int, login string) error

	GetTricks(ctx context.Context) ([]entity.Trick, error)
}

type userStorage struct {
	storage sql.DB
	logger  *logging.Logger
}

//NewUserStorage - конструктор
func NewUserStorage(db *sql.DB, logger *logging.Logger) *userStorage {
	return &userStorage{storage: *db, logger: logger}
}

func (s *userStorage) GetUser(ctx context.Context, login, password string) (*entity.User, error) {
	checkPing(ctx, s)
	s.logger.Infof("Поиск пользователя с login:%s", login)
	row := s.storage.QueryRowContext(ctx, "select * from pd.person where login = ?", login)
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
	s.logger.Debugf("Найден пользователь:\n %v", user)
	switch checkUserPassword(user, password) {
	case true:
		{
			s.logger.Debugf("Проверка пароля пользователя %s . SUCCESS", user.Login)
			return &user, nil
		}
	case false:
		{
			s.logger.Debugf("Проверка пароля пользователя %s . Неправильно набран пароль: %s", user.Login, password)
			return &user, errors.New("wrong password")
		}
	}
	return &user, nil
}

func (s *userStorage) Insert(ctx context.Context, user entity.User) error {
	s.logger.Infof("Добавление в БД пользователя: %v", user)
	checkPing(ctx, s)
	_, err := s.storage.ExecContext(
		ctx,
		"insert pd.person(first_name, second_name, login, password_hash, email) values(?, ?, ?, ?, ?)",
		user.FirstName, user.SecondName, user.Login, user.Password, user.Email)
	return err
}

func (s *userStorage) Delete(ctx context.Context, id int, login string) error {
	s.logger.Infof("Удаление в БД пользователя. id: %d, login: %s", id, login)
	checkPing(ctx, s)
	_, err := s.storage.ExecContext(ctx,
		"delete from pd.person where person_id = ? and login = ?",
		id,
		login)
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

func (s *userStorage) GetTricks(ctx context.Context) ([]*entity.Trick, error) {
	checkPing(ctx, s)
	s.logger.Debug("Ищем в базе все трюки")
	rows, err := s.storage.Query("select * from pd.tricks")
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			s.logger.Panic(err)
		}
	}()

	var tricks []*entity.Trick
	for rows.Next() {
		trick := entity.Trick{}
		err := rows.Scan(
			&trick.TrickId,
			&trick.DifficultyLevel,
			&trick.VideoId,
			&trick.Name,
			&trick.AdditionalNames,
			&trick.Description,
		)
		if err != nil {
			s.logger.Error(err.Error())
			continue
		}
		tricks = append(tricks, &trick)
	}

	return tricks, nil
}

func checkPing(ctx context.Context, s *userStorage) {
	if err := s.storage.PingContext(ctx); err != nil {
		fmt.Println(err.Error())
	}
}

func checkUserPassword(user entity.User, password string) bool {
	return user.Password == password
}
