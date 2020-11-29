package middleware

import (
	"time"
	"xiaoyin/lib/log"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func LogInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.RequestURI
		c.Next()
		duration := time.Since(start)
		if len(c.Errors) == 0 && c.Writer.Status() == 200 {
			log.Log.Info(path,
				zap.String("method", c.Request.Method),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.Duration("duration", duration),
			)
		}
	}
}
