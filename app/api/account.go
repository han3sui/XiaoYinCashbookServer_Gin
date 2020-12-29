package api

import (
	"strconv"
	"xiaoyin/app/service"
	"xiaoyin/lib/exception"
	"xiaoyin/lib/util"
	"xiaoyin/lib/validate"

	"github.com/gin-gonic/gin/binding"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

type BalanceInfo struct {
	TmpBalance *float64 `json:"tmp_balance" validate:"required" label:"账户临时余额"`
}

func ListAccountsByUid(c *gin.Context) {
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 101010, err)
		return
	}
	list, err := service.ListAccountsByUid(uid)
	if err != nil {
		exception.Common(c, 101011, errors.Wrap(err, "获取账户列表失败"))
		return
	}
	c.JSON(200, list)
}

func AccountManageList(c *gin.Context) {
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 101110, err)
		return
	}
	list, err := service.GetManageList(uid)
	if err != nil {
		exception.Common(c, 101111, errors.Wrap(err, "获取账户列表失败"))
		return
	}
	c.JSON(200, list)
}

func SaveAccount(c *gin.Context) {
	var req service.AccountInfo
	err := c.BindJSON(&req)
	if err != nil {
		exception.Common(c, 101210, errors.Wrap(err, "参数绑定失败"))
		return
	}
	err = validate.Validate.Struct(req)
	if err != nil {
		exception.Common(c, 101211, err)
		return
	}
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 101212, err)
		return
	}
	req.UserId = uid
	id, err := service.SaveAccount(&req)
	if err != nil {
		exception.Common(c, 101213, err)
		return
	}
	c.JSON(200, id)
}

func UpdateAccount(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		exception.Common(c, 101310, errors.New("更新参数错误"))
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		exception.Common(c, 101311, err)
		return
	}
	var balanceReq BalanceInfo
	err = c.ShouldBindBodyWith(&balanceReq, binding.JSON)
	if err != nil {
		exception.Common(c, 101312, errors.Wrap(err, "参数绑定失败"))
		return
	}
	err = validate.Validate.Struct(balanceReq)
	if err != nil {
		exception.Common(c, 101313, err)
		return
	}
	var req service.AccountInfo
	err = c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		exception.Common(c, 101314, errors.Wrap(err, "参数绑定失败"))
		return
	}
	err = validate.Validate.Struct(req)
	if err != nil {
		exception.Common(c, 101315, err)
		return
	}

	req.UserId, err = util.GetUid(c)
	if err != nil {
		exception.Common(c, 101316, err)
		return
	}
	req.ID = uint(id)
	err = service.UpdateAccount(*balanceReq.TmpBalance, &req)
	if err != nil {
		exception.Common(c, 101317, err)
		return
	}
	c.Status(200)
}

func DelAccountWithDetails(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		exception.Common(c, 101410, errors.New("删除参数错误"))
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		exception.Common(c, 101411, err)
		return
	}
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 101412, err)
		return
	}
	err = service.DelAccountWithDetails(uint(id), uid)
	if err != nil {
		exception.Common(c, 101413, errors.Wrap(err, "删除失败"))
	} else {
		c.Status(200)
	}
}

func GetDetailsCountByAid(c *gin.Context) {
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 101510, err)
		return
	}
	idStr := c.Param("id")
	if idStr == "" {
		exception.Common(c, 101511, errors.New("获取账户参数失败"))
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		exception.Common(c, 101512, err)
		return
	}
	count, err := service.GetDetailsCountByAid(uid, uint(id))
	if err != nil {
		exception.Common(c, 101513, errors.Wrap(err, "查询关联账单失败"))
		return
	}
	c.JSON(200, count)
}
