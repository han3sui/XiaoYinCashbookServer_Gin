package token

import (
	"xiaoyin/app/service/token"
	"xiaoyin/lib/exception"
	"xiaoyin/lib/validate"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

func Grant(c *gin.Context) {
	var req token.LoginInfo
	_ = c.BindJSON(&req)
	err := validate.Validate.Struct(req)
	if err != nil {
		exception.Common(c, 141010, err)
		return
	}
	tokenStr, err := token.Grant(&req)
	if err != nil {
		exception.Common(c, 141011, errors.Wrap(err, "登录失败"))
		return
	}
	c.JSON(200, tokenStr)
}
