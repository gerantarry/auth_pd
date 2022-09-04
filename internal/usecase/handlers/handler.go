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
	"time"
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

	var resp dto.StatusResponse
	//TODO defer cancel() - прописать кейс, когда отменяется по таймауту. Что делать тогда?
	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	err := h.storage.Insert(ctx, user)
	defer cancel()
	if err != nil {
		//TODO нужно захэшировать пароли ? изза одинаковых паролей БД кидает ошибку
		tErr, ok := err.(*mysql.MySQLError)
		if ok {
			h.logger.Error(tErr.Message)
			if tErr.Number == 1062 {
				resp = dto.StatusResponse{
					Description: "Это имя пользователя уже зарегистрировано.",
				}
			} else {
				resp.Description = err.Error()
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, resp)
		return
	}

	resp = dto.StatusResponse{Status: true, Description: "Регистрация прошла успешно!"}

	c.JSON(http.StatusOK, resp)
}
