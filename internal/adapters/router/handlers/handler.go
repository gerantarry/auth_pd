package handlers

import (
	"auth_pd/internal/adapters/db/mysql"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	storage *mysql.Storage
}

func NewHandler(stg *mysql.Storage) *Handler {
	return &Handler{storage: stg}
}

func (h *Handler) Register(c *gin.Context) {
	c.JSON(http.StatusOK, "register test OK")
}
