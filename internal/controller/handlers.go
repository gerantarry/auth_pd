package controller

import (
	"auth_pd/internal/adapters/db/mysql_"
	"auth_pd/internal/domain/dto"
	"auth_pd/internal/domain/entity"
	"auth_pd/internal/usecase/authorization/password"
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

	resp = dto.StatusResponse{Success: true, Description: registerSuccessMsg}
	c.JSON(http.StatusOK, resp)
}

// Login
// запрос вида
//{"username": "$username", "password": "$password"}
func (h *Handler) Login(c *gin.Context) {
	cCp := c.Copy()
	var bodyBuff []byte
	_, err := cCp.Request.Body.Read(bodyBuff)
	h.logger.Debugf("Поступил запрос, %v", string(bodyBuff))
	var loginDto dto.LoginRequestDto
	if err := c.MustBindWith(&loginDto, binding.JSON); err != nil {
		h.logger.Errorf("Не удалось разобрать запрос. Причина - %v", err.Error())
		return
	}
	h.logger.Infof("Запрос на login пользователем %s", loginDto.Username)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var response dto.StatusResponse
	user, err := h.storage.GetUser(ctx, loginDto.Username)
	if err != nil {
		h.logger.Errorf("Ошибка при поиске пользователя в базе: %v", err)
		response.Description = userNotFoundMsg
		c.JSON(http.StatusBadRequest, response)
	}

	isSame := password.VerifyPassword(loginDto.Password, user.Password)
	if isSame {
		//TODO формируем и возвращаем токен. Через c.Writer.Header установить Authorize: Bearer token
		h.logger.Info("Формирование токена для пользователя")
		response.Success = true
		response.Description = "Авторизация успешно завершена."
		c.JSON(http.StatusOK, response)
	} else {
		response.Description = userNotFoundMsg
		c.JSON(http.StatusBadRequest, response)
	}

}
