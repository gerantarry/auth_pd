package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"strings"
	"unicode"
)

type router struct {
	*gin.Engine
}

var r *gin.Engine

func NewRouter() *router {
	r = gin.New()
	r.Use(gin.Recovery(), gin.Logger())
	if err := r.SetTrustedProxies(nil); err != nil {
		fmt.Println("Ошибка при настройке разрешённых proxy")
		panic(any(err))
	}
	return &router{r}
}

//убрать символы переноса строки и т.д. из входящих данных
func formatReaderData(r io.Reader) (string, error) {
	buf := make([]byte, 2048)
	n, err := r.Read(buf)
	var str, resStr string

	for {
		for _, b := range buf[:n] {
			str = str + string(b)
		}
		if err == io.EOF {
			resStr = cleanString(str)
			return resStr, nil
		}
		if err != nil {
			return "", err
		}

		resStr = cleanString(str)
		return resStr, nil
	}

}

//убирает сначала управляющие символы затем '\'
func cleanString(str string) string {
	clean := strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, str)
	return clean
}
