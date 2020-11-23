package account

import (
	"strconv"
	"xiaoyin/app/service/account"
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

func ListByUid(c *gin.Context) {
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 101010, err)
		return
	}
	list, err := account.ListByUid(uid)
	if err != nil {
		exception.Common(c, 101011, errors.Wrap(err, "获取账户列表失败"))
		return
	}
	c.JSON(200, list)
}

func ManageList(c *gin.Context) {
	uid, err := util.GetUid(c)
	if err != nil {
		exception.Common(c, 101110, err)
		return
	}
	list, err := account.GetManageList(uid)
	if err != nil {
		exception.Common(c, 101111, errors.Wrap(err, "获取账户列表失败"))
		return
	}
	c.JSON(200, list)
}

func Save(c *gin.Context) {
	var req account.Info
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
	id, err := account.Save(&req)
	if err != nil {
		exception.Common(c, 101213, err)
		return
	}
	c.JSON(200, id)
}

func Update(c *gin.Context) {
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
	var req account.Info
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
	err = account.Update(*balanceReq.TmpBalance, &req)
	if err != nil {
		exception.Common(c, 101317, err)
		return
	}
	c.Status(200)
}

func DelWithDetails(c *gin.Context) {
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
	err = account.DelWithDetails(uint(id), uid)
	if err != nil {
		exception.Common(c, 101413, errors.Wrap(err, "删除失败"))
	} else {
		c.Status(200)
	}
}
