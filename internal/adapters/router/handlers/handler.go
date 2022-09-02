package handlers

import (
	"auth_pd/internal/adapters/db/mysql_"
	"auth_pd/internal/domain/dto"
	"auth_pd/internal/domain/entity"
	"auth_pd/pkg/logging"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-sql-driver/mysql"
	"net/http"
)

type Handler struct {
	storage mysql_.Storage
	logger  *logging.Logger
}

func NewHandler(stg mysql_.Storage, l *logging.Logger) *Handler {
	return &Handler{
		storage: stg,
		logger:  l,
	}
}

//TODO нарушено логирование при двойном вызове c.Request.Body
func (h *Handler) Register(c *gin.Context) {
	var regForm dto.RegisterForm
	h.logger.Debug("Получен запрос. Начинаем биндить тело")
	if err := c.MustBindWith(&regForm, binding.JSON); err != nil {
		h.logger.Errorf("Не удалось разобрать запрос. Причина - %v", err.Error())
		return
	}
	fmt.Println(regForm)
	user := entity.User{
		FirstName: regForm.FirstName,
		Login:     regForm.Username,
		Password:  regForm.Password,
		Email:     regForm.Email,
	}

	err := h.storage.Insert(context.Background(), user)
	if err != nil {
		var errResp dto.StatusResponse
		tErr, ok := err.(*mysql.MySQLError)
		if ok {
			h.logger.Error(tErr.Message)
			if tErr.Number == 1062 {
				errResp = dto.StatusResponse{
					Description: "Это имя пользователя уже зарегистрировано.",
				}
			} else {
				errResp.Description = err.Error()
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errResp)
		return
	}

	c.JSON(http.StatusOK, "register test OK")
}
