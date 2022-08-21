package main

import (
	"auth_pd/internal/adapters/db/mysql"
	"auth_pd/internal/adapters/router"
	"auth_pd/internal/adapters/router/handlers"
	"auth_pd/internal/config"
	"auth_pd/pkg/logging"
	"database/sql"
)

var cfg *config.Config

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
	dataSourceName := formatDBSourceString()
	logger := logging.GetLogger()
	db, err := sql.Open(mysql.DriverMySQL, dataSourceName)
	if err != nil {
		logger.Panicf("Ошибка при открытии соединения с БД: %v", err)
	}
	var storage mysql.Storage = mysql.NewUserStorage(db, logger)
	handler := handlers.NewHandler(&storage)
	startServer(handler, logger)
}

func startServer(h *handlers.Handler, l *logging.Logger) {
	r := router.NewRouter()
	r.SetLogger(l)

	r.POST("/register", h.Register)

	err := r.Run(":8080")
	if err != nil {
		panic(any(err))
	}
}

//приводим данные для бд к форматудля открытия соединения
func formatDBSourceString() string {
	return cfg.Database.Login +
		":" + cfg.Database.Password +
		"@tcp(" + cfg.Database.BindIp + ":" + cfg.Database.Port + ")/" + cfg.Database.Scheme
}
