package handlers

import (
	"auth_pd/internal/adapters/db/mysql"
	"github.com/gin-gonic/gin"
)

type handler struct {
	storage *mysql.Storage
}

func NewHandler(stg *mysql.Storage) *handler {
	return &handler{storage: stg}
}

func Register(c *gin.Context) {

}
