package main

import (
	"auth_pd/internal/adapters/router"
	"auth_pd/internal/config"
	"auth_pd/internal/controller"
	"auth_pd/pkg/logging"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var cfg *config.Config
var ctx, cancel = context.WithTimeout(context.Background(), time.Second*15)

/*
1. необходимо сделать многопоточность при обработки запросов (возможно все запросы и так работают в многопоточности)
2. пробрасывать ctx context.Context; Определить контекст ContextWithTimeout()


Trace — вывод всего подряд. На тот случай, если Debug не позволяет локализовать ошибку. В нем полезно отмечать вызовы разнообразных блокирующих и асинхронных операций.
Debug — журналирование моментов вызова «крупных» операций. Старт/остановка потока, запрос пользователя и т.п.
Info — разовые операции, которые повторяются крайне редко, но не регулярно. (загрузка конфига, плагина, запуск бэкапа)
Warning — неожиданные параметры вызова, странный формат запроса, использование дефолтных значений в замен не корректных. Вообще все, что может свидетельствовать о не штатном использовании.
Error — повод для внимания разработчиков. Тут интересно окружение конкретного места ошибки.
Fatal — тут и так понятно. Выводим все до чего дотянуться можем, так как дальше приложение работать не будет.
*/
func main() {
	cfg = config.GetConfig()
	logger := logging.GetLogger()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+cfg.Database.BindIp+":"+cfg.Database.Port))
	if err != nil {
		logger.Panicf("Ошибка при открытии соединения с БД: %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(any(err))
		}
	}()
	/*var storage mng.Storage = mng.NewUserStorage(client, logger)
	handler := controller.NewHandler(storage, logger)
	startServer(handler)*/
}

func startServer(h *controller.Handler) {
	r := router.NewRouter()
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	err := r.Run(":8080")
	if err != nil {
		panic(any(err))
	}
}
