package check

import (
	"xiaoyin/app/service/check"
	"xiaoyin/lib/exception"
	"xiaoyin/lib/util"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func List(c *gin.Context) {
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 161010, err)
		return
	}
	r, err := check.List(uid)
	if err != nil {
		exception.Common(c, 161011, errors.Wrap(err, "获取核账记录失败"))
		return
	}
	if r == nil {
		c.JSON(200, []interface{}{})
	} else {
		c.JSON(200, r)
	}
}
