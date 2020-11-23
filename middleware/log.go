package middleware

import (
	"time"
	"xiaoyin/lib/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LogInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		duration := time.Since(start)
		if len(c.Errors) == 0 && c.Writer.Status() == 200 {
			log.Log.Info("Success",
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.Duration("duration", duration),
			)
		}
	}
}
