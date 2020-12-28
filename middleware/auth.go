package middleware

import (
	"strings"
	"time"
	"xiaoyin/app/service"
	"xiaoyin/lib/exception"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", -1)
		tokenInfo, err := service.ParseToken(tokenStr)
		if err != nil {
			err = errors.Wrap(err, "用户校验失败")
			exception.Common(c, 9999, err)
			c.Abort()
		} else if int64(tokenInfo.Exp) < time.Now().Unix() {
			exception.Common(c, 9999, errors.New("登录已过期"))
			c.Abort()
		} else {
			c.Set("uid", tokenInfo.ID)
			c.Next()
		}
	}
}
