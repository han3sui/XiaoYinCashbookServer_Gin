package api

import (
	"fmt"
	"xiaoyin/app/service"
	"xiaoyin/lib/exception"
	"xiaoyin/lib/util"
	"xiaoyin/lib/validate"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func ListCheck(c *gin.Context) {
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 161010, err)
		return
	}
	var req service.CheckSearchParams
	err = c.ShouldBindQuery(&req)
	if err != nil {
		exception.Common(c, 161011, errors.New("参数绑定失败"))
		return
	}
	err = validate.Validate.Struct(req)
	if err != nil {
		exception.Common(c, 161012, err)
		return
	}
	fmt.Println(req)
	r, err := service.ListCheck(uid, req)
	if err != nil {
		exception.Common(c, 161013, errors.Wrap(err, "获取核账记录失败"))
		return
	}
	if r == nil {
		c.JSON(200, []interface{}{})
	} else {
		c.JSON(200, r)
	}
}
