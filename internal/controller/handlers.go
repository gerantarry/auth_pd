package controller

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
	userNotFoundMsg         = "Неправильно набрал логин и/или пароль."
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
	h.logger.Debug("Получен запрос. Начинаем биндить тело")
	var regDto dto.RegisterRequestDto
	if err := c.MustBindWith(&regDto, binding.JSON); err != nil {
		h.logger.Errorf("Не удалось разобрать запрос. Причина - %v", err.Error())
		return
	}
	fmt.Println(regDto) //удалить, когда научусь логировать тело

	hashedPsw := password.HashPassword(regDto.Password)
	user := entity.User{
		FirstName: regDto.FirstName,
		Login:     regDto.Username,
		Password:  hashedPsw,
		Email:     regDto.Email,
	}

	var resp dto.StatusResponse
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err := h.storage.Insert(ctx, user)
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

func (h *Handler) Login(c *gin.Context) {
	gin.BasicAuth()

}
