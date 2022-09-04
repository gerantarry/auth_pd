package password

import (
	"auth_pd/pkg/logging"
	"golang.org/x/crypto/bcrypt"
)

var logger = logging.GetLogger()

func HashPassword(psw string) string {
	logger.Debug("Хэшируем пароль")
	pswHash, err := bcrypt.GenerateFromPassword([]byte(psw), 13)
	if err != nil {
		logger.Panic(err)
	}
	return string(pswHash)
}

func VerifyPassword(psw, hashedPsw string) bool {
	res := bcrypt.CompareHashAndPassword([]byte(hashedPsw), []byte(psw))
	if res != nil {
		logger.Debugf("Пароли не совпадают")
		return false
	}
	return true
}
