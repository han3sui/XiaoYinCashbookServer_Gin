package icon

import (
	"xiaoyin/app/service/icon"
	"xiaoyin/lib/exception"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func List(c *gin.Context) {
	list, err := icon.List()
	if err != nil {
		exception.Common(c, 131010, errors.Wrap(err, "图标获取失败"))
		return
	}
	c.JSON(200, list)
}
