package middleware

import (
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
				log.Log.Error("[Recovery from panic]",
					zap.Any("error", err),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.RequestURI),
					zap.String("ip", c.ClientIP()),
					zap.String("user-agent", c.Request.UserAgent()),
					zap.String("token", c.GetHeader("Authorization")),
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
