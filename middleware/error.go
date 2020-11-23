package middleware

import (
	"net/http/httputil"
	"time"
	"xiaoyin/lib/log"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type ErrorInfo struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Request   string `json:"request"`
	Timestamp int64  `json:"timestamp"`
}

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatusJSON(500, ErrorInfo{
					Code:      500,
					Message:   "系统未知错误",
					Request:   c.Request.Method + " " + c.Request.URL.Path,
					Timestamp: time.Now().Unix(),
				})
				//str := fmt.Sprintf("%s", err)
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				log.Log.Error("[Recovery from panic]",
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("query", c.Request.URL.RawQuery),
					zap.String("ip", c.ClientIP()),
					zap.String("user-agent", c.Request.UserAgent()),
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
					//zap.String("stack", string(debug.Stack())),
				)
			}
		}()
		c.Next()
	}
}

func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(404, ErrorInfo{
			Code:      404,
			Message:   "NOT FOUND",
			Request:   c.Request.Method + " " + c.Request.URL.Path,
			Timestamp: time.Now().Unix(),
		})
	}
}
