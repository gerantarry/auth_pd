package main

import (
	"auth_pd/internal/adapters/router"
	"auth_pd/pkg/logging"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
)

/*
1. необходимо сделать многопоточность при обработки запросов (возможно все запросы и так работают в многопоточности)
2. пробрасывать ctx context.Context; Определить контекст ContextWithTimeout()
3. сделать конфиг файл с которого должны подтягиваться значения { логин и пароль к БД, название схемы БД }
4. добавить логирование



Trace — вывод всего подряд. На тот случай, если Debug не позволяет локализовать ошибку. В нем полезно отмечать вызовы разнообразных блокирующих и асинхронных операций.
Debug — журналирование моментов вызова «крупных» операций. Старт/остановка потока, запрос пользователя и т.п.
Info — разовые операции, которые повторяются крайне редко, но не регулярно. (загрузка конфига, плагина, запуск бэкапа)
Warning — неожиданные параметры вызова, странный формат запроса, использование дефолтных значений в замен не корректных. Вообще все, что может свидетельствовать о не штатном использовании.
Error — повод для внимания разработчиков. Тут интересно окружение конкретного места ошибки.
Fatal — тут и так понятно. Выводим все до чего дотянуться можем, так как дальше приложение работать не будет.
*/
func init() {
	projectPath := filepath.Dir(os.Args[0])
	err := os.Setenv("PROJECT_DIR", projectPath)
	if err != nil {
		panic(any(err))
	}
	logging.Init()
}
func main() {
	l := logging.GetLogger()
	startServer(l)
}

func startServer(l *logging.Logger) {
	r := router.NewRouter()
	r.SetLogger(l)

	r.GET("/test", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.POST("/test", func(c *gin.Context) {
		c.JSON(200, "{answer: 6}")
	})

	err := r.Run(":8080")
	if err != nil {
		panic(any(err))
	}
}
