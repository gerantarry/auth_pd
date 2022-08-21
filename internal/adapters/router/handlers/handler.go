package handlers

import (
	"auth_pd/internal/adapters"
	"auth_pd/internal/adapters/db/mysql"
	"auth_pd/internal/domain/dto"
	"github.com/gin-gonic/gin"
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
	c.BindJSON(&regForm)
	c.JSON(http.StatusOK, "register test OK")
}
