package api

import (
	"xiaoyin/app/service"
	"xiaoyin/lib/exception"
	"xiaoyin/lib/validate"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

func GrantToken(c *gin.Context) {
	var req service.LoginInfo
	_ = c.BindJSON(&req)
	err := validate.Validate.Struct(req)
	if err != nil {
		exception.Common(c, 141010, err)
		return
	}
	tokenStr, err := service.GrantToken(&req)
	if err != nil {
		exception.Common(c, 141011, errors.Wrap(err, "登录失败"))
		return
	}
	c.JSON(200, tokenStr)
}
