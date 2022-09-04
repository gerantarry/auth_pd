package user

import (
	"auth_pd/internal/adapters/db/mysql_"
	"auth_pd/internal/domain/dto"
	"auth_pd/internal/domain/entity"
	"auth_pd/internal/usecase/password"
	"auth_pd/pkg/logging"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-sql-driver/mysql"
	"net/http"
	"time"
)

type Controller struct {
	storage mysql_.Storage
	logger  *logging.Logger
}

func NewController(stg mysql_.Storage, l *logging.Logger) *Controller {
	return &Controller{
		storage: stg,
		logger:  l,
	}
}

//TODO нарушено логирование при двойном вызове c.Request.Body
func (ctrl *Controller) Register(c *gin.Context) {
	ctrl.logger.Debug("Получен запрос. Начинаем биндить тело")
	var regForm dto.RegisterForm
	if err := c.MustBindWith(&regForm, binding.JSON); err != nil {
		ctrl.logger.Errorf("Не удалось разобрать запрос. Причина - %v", err.Error())
		return
	}
	fmt.Println(regForm) //удалить, когда научусь логировать тело

	hashedPsw := password.HashPassword(regForm.Password)
	user := entity.User{
		FirstName: regForm.FirstName,
		Login:     regForm.Username,
		Password:  hashedPsw,
		Email:     regForm.Email,
	}

	var resp dto.StatusResponse
	//TODO defer cancel() - прописать кейс, когда отменяется по таймауту. Что делать тогда?
	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	err := ctrl.storage.Insert(ctx, user)
	defer cancel()
	if err != nil {
		tErr, ok := err.(*mysql.MySQLError)

		if ok {
			ctrl.logger.Error(tErr.Message)
			if tErr.Number == 1062 {
				resp.Description = "Это имя пользователя уже зарегистрировано."
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
