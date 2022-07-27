package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type router struct {
	*gin.Engine
}

var r *gin.Engine

func NewRouter() *router {
	r = gin.New()
	r.Use(gin.Recovery())
	return &router{r}
}

// SetLogger установка кастомного логгера
func (r *router) SetLogger(logger Logger) {
	r.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		//читаем тело запроса
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Printf("Cannot read body")
		}
		defer c.Request.Body.Close()

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
			string(body),
			params.ErrorMessage,
			params.Latency)

		logger.Infof(paramsString)
	})
}
