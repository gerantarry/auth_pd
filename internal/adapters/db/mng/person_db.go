package mng

import (
	"auth_pd/internal/domain/entity"
	"auth_pd/pkg/logging"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var cntx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

type Storage interface {
	GetUser(login string) (*entity.User, error)
	Insert(user entity.User) error
	Delete(id int, login string) error
	Update(id int, login string) error
}

type userStorage struct {
	client mongo.Client
	logger *logging.Logger
}

func NewUserStorage(client *mongo.Client, logger *logging.Logger) *userStorage {
	return &userStorage{client: *client, logger: logger}
}

func (s *userStorage) GetUser(login string) (*entity.User, error) {
	collection := s.client.Database("pd").Collection("persons")
	defer cancel()
	var result entity.User
	err := collection.FindOne(cntx, bson.M{"login": login}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		s.logger.Errorf("Клиент %s не найден", login)
	}
	return &result, err
}
