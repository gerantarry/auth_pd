package handlers

import (
	"auth_pd/internal/adapters"
	"auth_pd/internal/adapters/db/mysql"
	"auth_pd/internal/domain/dto"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type Handler struct {
	storage *mysql.Storage
	logger  *adapters.AppLogger
}

func NewHandler(stg *mysql.Storage, l *adapters.AppLogger) *Handler {
	return &Handler{
		storage: stg,
		logger:  l,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var regForm dto.RegisterForm
	h.logger.Debug("Получен запрос. Начинаем биндить тело")
	if err := c.MustBindWith(&regForm, binding.JSON); err != nil {
		h.logger.Errorf("Не удалось разобрать запрос. Причина - %v", err.Error())
		return
	}
	fmt.Println(regForm)
	c.JSON(http.StatusOK, "register test OK")

}
