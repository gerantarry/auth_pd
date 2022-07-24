package router

import (
	"github.com/gin-gonic/gin"
	"time"
)

type router struct {
	*gin.Engine
}

var r *gin.Engine

func NewRouter() *router {
	return &router{r}
}

func Init() {
	gin.Logger()

}

// SetLogger установка кастомного логгера
func (r *router) SetLogger(logger Logger) {
	r.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

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

		logger.Info(params)
	})
}
