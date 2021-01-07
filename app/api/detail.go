package api

import (
	"strconv"
	"xiaoyin/app/service"
	"xiaoyin/lib/exception"
	"xiaoyin/lib/util"
	"xiaoyin/lib/validate"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

func GetAllDays(c *gin.Context) {
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 121010, err)
		return
	}
	r, err := service.AllDaysCount(uid)
	if err != nil {
		exception.Common(c, 121011, errors.Wrap(err, "获取记账天数失败"))
		return
	}
	c.JSON(200, r)
}
func IsExistUncheck(c *gin.Context) {
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 121110, err)
		return
	}
	checkTimeStr := c.Param("time")
	if checkTimeStr == "" {
		exception.Common(c, 121111, errors.New("核账时间参数错误"))
		return
	}
	checkTime, err := strconv.ParseUint(checkTimeStr, 10, 64)
	if err != nil {
		exception.Common(c, 121112, errors.New("核账时间错误"))
		return
	}
	boolean, err := service.IsExistUncheck(uid, uint(checkTime))
	if err != nil {
		exception.Common(c, 121113, errors.Wrap(err, "获取未核账账单失败"))
		return
	}
	if boolean {
		c.JSON(200, 0)
	} else {
		c.JSON(200, 1)
	}
}

func Bill(c *gin.Context) {
	var params service.DetailSearchParams
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 121210, err)
		return
	}
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		exception.Common(c, 121211, errors.Wrap(err, "类型转换失败"))
		return
	}
	params.Year = year
	res, err := service.Bill(uid, params)
	if err != nil {
		exception.Common(c, 121212, errors.Wrap(err, "获取账单失败"))
		return
	}
	c.JSON(200, res)
}

func Chart(c *gin.Context) {
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 121310, err)
		return
	}
	var req service.DetailSearchParams
	err = c.ShouldBindQuery(&req)
	if err != nil {
		exception.Common(c, 121311, errors.Wrap(err, "参数绑定失败"))
		return
	}
	err = validate.Validate.Struct(req)
	if err != nil {
		exception.Common(c, 121312, err)
		return
	}
	list, err := service.Chart(uid, req)
	if err != nil {
		exception.Common(c, 121313, errors.Wrap(err, "图表获取失败"))
		return
	}
	c.JSON(200, list)
}

func ListDetailsByParams(c *gin.Context) {
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 121410, err)
		return
	}
	var req service.DetailSearchParams
	err = c.ShouldBindQuery(&req)
	if err != nil {
		exception.Common(c, 121411, errors.Wrap(err, "参数绑定失败"))
		return
	}
	err = validate.Validate.Struct(req)
	if err != nil {
		exception.Common(c, 121412, err)
		return
	}
	list, err := service.ListByParams(uid, req)
	if err != nil {
		exception.Common(c, 121413, errors.Wrap(err, "账单查询出错"))
		return
	}
	if list == nil {
		c.JSON(200, []interface{}{})
	} else {
		c.JSON(200, list)
	}
}

func SaveDetail(c *gin.Context) {
	var req service.Detail
	err := c.BindJSON(&req)
	if err != nil {
		exception.Common(c, 121510, errors.Wrap(err, "参数绑定失败"))
		return
	}
	//err = validate.Validate.Struct(req)
	//if err != nil {
	//	exception.Common(c, 5012, err)
	//	return
	//}
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 121511, err)
		return
	}
	req.UserId = uid
	list, err := service.SaveDetail(&req)
	if err != nil {
		exception.Common(c, 121512, err)
		return
	}
	c.JSON(200, list[0])
}

func UpdateDetail(c *gin.Context) {
	var req service.Detail
	err := c.BindJSON(&req)
	if err != nil {
		exception.Common(c, 121610, errors.Wrap(err, "参数绑定失败"))
		return
	}
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 121611, err)
		return
	}
	req.UserId = uid
	idStr := c.Param("id")
	if idStr == "" {
		exception.Common(c, 121612, errors.New("更新参数错误"))
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		exception.Common(c, 161213, errors.New("更新参数错误"))
		return
	}
	req.ID = uint(id)
	list, err := service.UpdateDetail(&req)
	if err != nil {
		exception.Common(c, 121614, err)
		return
	}
	c.JSON(200, list[0])
}

func DelDetail(c *gin.Context) {
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 121710, err)
	}
	idStr := c.Param("id")
	if idStr == "" {
		exception.Common(c, 121711, errors.New("删除参数错误"))
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		exception.Common(c, 121712, err)
		return
	}
	err = service.Del(uint(id), uid)
	if err != nil {
		exception.Common(c, 121713, errors.Wrap(err, "删除失败"))
		return
	}
	c.Status(200)
}

//获取报销账单
func ListClaim(c *gin.Context) {
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 121810, err)
		return
	}
	claimStr := c.Param("claim")
	if claimStr == "" {
		exception.Common(c, 121812, errors.New("报销参数错误"))
		return
	}
	claim, err := strconv.ParseInt(claimStr, 10, 64)
	if err != nil {
		exception.Common(c, 121813, err)
		return
	}
	r, err := service.ListClaim(uid, int(claim))
	if err != nil {
		exception.Common(c, 121811, errors.Wrap(err, "获取报销账单失败"))
		return
	}
	if r == nil {
		c.JSON(200, []interface{}{})
	} else {
		c.JSON(200, r)
	}
}

func ListMoneyByParams(c *gin.Context) {
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 121910, err)
		return
	}
	var req service.DetailSearchParams
	err = c.ShouldBindQuery(&req)
	if err != nil {
		exception.Common(c, 121911, errors.Wrap(err, "参数绑定失败"))
		return
	}
	err = validate.Validate.Struct(req)
	if err != nil {
		exception.Common(c, 121912, err)
		return
	}
	data, err := service.ListMoneyByParams(uid, req)
	if err != nil {
		exception.Common(c, 121913, errors.Wrap(err, "账单总额查询出错"))
		return
	}
	c.JSON(200, data)
}

func BatchUpdateClaim(c *gin.Context) {
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 122010, err)
		return
	}
	var req []service.ClaimDetail
	err = c.ShouldBindJSON(&req)
	if err != nil {
		exception.Common(c, 122011, errors.Wrap(err, "参数绑定失败"))
		return
	}
	err = service.BatchUpdateClaim(uid, req)
	if err != nil {
		exception.Common(c, 122012, errors.Wrap(err, "报销更新失败"))
		return
	}
	c.Status(200)
}
