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

const (
	usernameAlreadyExistMsg = "Это имя пользователя уже зарегистрировано."
	registerSuccessMsg      = "Регистрация прошла успешно!"
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
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	err := ctrl.storage.Insert(ctx, user)
	defer cancel()
	if err != nil {
		switch e := err.(type) {
		case *mysql.MySQLError:
			{
				if e.Number == 1062 {
					resp.Description = usernameAlreadyExistMsg
				} else {
					resp.Description = e.Message
				}
				c.AbortWithStatusJSON(http.StatusBadRequest, resp)
				return
			}
		default:
			break
		}
		if err.Error() == context.DeadlineExceeded.Error() {
			c.AbortWithStatus(http.StatusRequestTimeout)
			return
		}

		resp.Description = err.Error()
		c.AbortWithStatusJSON(http.StatusOK, resp)
		return
	}

	resp = dto.StatusResponse{Status: true, Description: registerSuccessMsg}
	c.JSON(http.StatusOK, resp)
}
