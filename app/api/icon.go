package api

import (
	"xiaoyin/app/service"
	"xiaoyin/lib/exception"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func ListIcons(c *gin.Context) {
	list, err := service.ListIcons()
	if err != nil {
		exception.Common(c, 131010, errors.Wrap(err, "图标获取失败"))
		return
	}
	c.JSON(200, list)
}
