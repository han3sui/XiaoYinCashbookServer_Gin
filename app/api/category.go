package api

import (
	"strconv"
	"xiaoyin/app/service"
	"xiaoyin/lib/exception"
	"xiaoyin/lib/util"
	"xiaoyin/lib/validate"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func ListCategoryByUid(c *gin.Context) {
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 111010, err)
		return
	}
	list, err := service.ListTreeByUid(uid)
	if err != nil {
		exception.Common(c, 111011, errors.Wrap(err, "获取分类列表失败"))
		return
	}
	c.JSON(200, list)
}

func SaveCategory(c *gin.Context) {
	var req service.Category
	_ = c.BindJSON(&req)
	err := validate.Validate.Struct(req)
	if err != nil {
		exception.Common(c, 111110, err)
		return
	}
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 111211, err)
		return
	}
	req.UserId = uid
	id, err := service.Save(&req)
	if err != nil {
		exception.Common(c, 111212, err)
		return
	}
	c.JSON(200, id)
}

func UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		exception.Common(c, 111310, errors.New("更新参数错误"))
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		exception.Common(c, 111311, err)
		return
	}
	var req service.Category
	_ = c.BindJSON(&req)
	err = validate.Validate.Struct(req)
	if err != nil {
		exception.Common(c, 111312, err)
		return
	}
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 111313, err)
		return
	}
	req.UserId = uid
	req.ID = uint(id)
	err = service.Update(&req)
	if err != nil {
		exception.Common(c, 111314, err)
		return
	}
	c.Status(200)
}

func DelCategoryWithDetails(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		exception.Common(c, 111410, errors.New("删除参数错误"))
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		exception.Common(c, 111411, err)
		return
	}
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 111412, err)
		return
	}
	err = service.DelCategoryWithDetails(uint(id), uid)
	if err != nil {
		exception.Common(c, 111413, errors.Wrap(err, "删除失败"))
		return
	} else {
		c.Status(200)
	}
}

func GetDetailsCountByCid(c *gin.Context) {
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 111510, err)
		return
	}
	idStr := c.Param("id")
	if idStr == "" {
		exception.Common(c, 111511, errors.New("获取分类参数失败"))
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		exception.Common(c, 111512, err)
		return
	}
	count, err := service.GetDetailsCountByCid(uid, uint(id))
	if err != nil {
		exception.Common(c, 111513, errors.Wrap(err, "查询关联账单失败"))
		return
	}
	c.JSON(200, count)
}
