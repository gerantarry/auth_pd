package router

import (
	"auth_pd/internal/adapters"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"strings"
	"time"
	"unicode"
)

type router struct {
	*gin.Engine
}

var r *gin.Engine

func NewRouter() *router {
	r = gin.New()
	r.Use(gin.Recovery())
	if err := r.SetTrustedProxies(nil); err != nil {
		fmt.Println("Ошибка при настройке разрешённых proxy")
		panic(any(err))
	}
	return &router{r}
}

// SetLogger установка кастомного логгера
func (r *router) SetLogger(logger adapters.Logger) {
	r.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		//читаем тело запроса
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Printf("Cannot read body")
		}
		defer func() {
			err := c.Request.Body.Close()
			if err != nil {
				logger.Errorf("Could not close Reader: \n" + err.Error())
			}
		}()
		fmtBody, err := formatReaderData(bytes.NewReader(body))
		if err != nil {
			logger.Warnf(err.Error())
		}

		// before request

		c.Next()

		// after request
		latency := time.Since(start)

		params := gin.LogFormatterParams{
			Request:      c.Request,
			TimeStamp:    time.Now(),
			StatusCode:   c.Writer.Status(),
			Latency:      latency,
			ClientIP:     c.ClientIP(),
			Method:       c.Request.Method,
			Path:         path,
			ErrorMessage: c.Errors.String(),
			BodySize:     c.Writer.Size(),
			Keys:         c.Keys,
		}
		if raw != "" {
			path = path + "?" + raw
		}
		params.Path = path

		paramsString := fmt.Sprintf("Method: %s | Status code: %d | Path: %s | Query: %v | Body: %v | ErrorMessage: %s | Latency: %v |",
			params.Method,
			params.StatusCode,
			params.Path,
			params.Keys,
			fmtBody,
			params.ErrorMessage,
			params.Latency)

		logger.Infof(paramsString)
	})
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
