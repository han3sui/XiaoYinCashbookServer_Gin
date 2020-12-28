package api

import (
	"xiaoyin/app/service"
	"xiaoyin/lib/exception"
	"xiaoyin/lib/util"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func GetUser(c *gin.Context) {
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 151010, err)
		return
	}
	res, err := service.GetInfo(uid)
	if err != nil {
		exception.Common(c, 151011, errors.Wrap(err, "查询用户信息失败"))
		return
	}
	c.JSON(200, res)
}

func UpdateCheckTime(c *gin.Context) {
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 151110, err)
		return
	}
	var checkInfo service.CheckInfo
	err = c.BindJSON(&checkInfo)
	if err != nil {
		exception.Common(c, 151112, errors.Wrap(err, "绑定核账内容失败"))
		return
	}
	checkTime, err := service.UpdateCheckTime(uid, &checkInfo)
	if err != nil {
		exception.Common(c, 151111, errors.Wrap(err, "核账失败"))
		return
	}
	c.JSON(200, checkTime)
}
