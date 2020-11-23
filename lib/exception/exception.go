package exception

import (
	"fmt"
	"strings"
	"time"
	"xiaoyin/lib/log"
	"xiaoyin/lib/validate"

	"go.uber.org/zap"

	"github.com/go-sql-driver/mysql"

	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

type ErrorInfo struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Request   string `json:"request"`
	Timestamp int64  `json:"timestamp"`
}

func Common(c *gin.Context, code int, err error) {
	message := ""
	switch v := errors.Cause(err).(type) {
	case validator.ValidationErrors:
		message = validate.Translate(v)
	case *mysql.MySQLError:
		//message = v.Message
		message = strings.Split(fmt.Sprintf("%v", err), ":")[0]
	default:
		message = strings.Split(fmt.Sprintf("%v", err), ":")[0]
		//message = v.Error()
	}
	_ = c.Error(err)
	log.Log.Error(fmt.Sprintf("%v", err),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("query", c.Request.URL.RawQuery),
		zap.String("ip", c.ClientIP()),
		zap.String("user-agent", c.Request.UserAgent()),
		zap.String("token", c.GetHeader("Authorization")),
		//zap.String("stack", string(debug.Stack())),
		//zap.String("statck", fmt.Sprintf("%+v", err)),
	)
	c.AbortWithStatusJSON(400, ErrorInfo{
		Code:      code,
		Message:   message,
		Request:   c.Request.Method + " " + c.Request.URL.Path,
		Timestamp: time.Now().Unix(),
	})
}
