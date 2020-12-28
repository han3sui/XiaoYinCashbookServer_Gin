package service

import (
	"fmt"
	"xiaoyin/app/model"
	"xiaoyin/lib/db"
	"xiaoyin/lib/util"

	"github.com/pkg/errors"
)

//账户管理列表
type ManageList struct {
	Id       uint    `json:"id"`
	Name     string  `json:"name"`
	Balance  float64 `json:"balance"`
	Icon     string  `json:"icon"`
	AddTotal uint    `json:"add_total"`
}

//账户列表
type AccountList struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type AccountInfo = model.Account

func SaveAccount(data *model.Account) (id uint, err error) {
	id, _, err = model.CheckAccountExist(data.UserId, data.Name)
	if err != nil {
		err = errors.Wrap(err, "账户重复检查失败")
		return
	}
	if id != 0 {
		err = errors.New("该账户已存在")
		return
	}
	id, err = data.Save()
	if err != nil {
		err = errors.Wrap(err, "账户保存失败")
	}
	return
}

func GetDetailsCountByAid(uid uint, id uint) (count int64, err error) {
	err = db.DB.Model(&model.Detail{}).Where("user_id = ? AND account_id = ?", uid, id).Count(&count).Error
	return
}

func UpdateAccount(tmpBalance float64, data *model.Account) (err error) {
	id, balance, err := model.CheckAccountExist(data.UserId, data.Name)
	if err != nil {
		err = errors.Wrap(err, "账户重复检查失败")
		return
	}
	if id != 0 && id != data.ID {
		err = errors.New("该账户已存在")
		return
	}
	if tmpBalance != *data.Balance {
		v := util.FloatAdd(balance, util.FloatSub(*data.Balance, tmpBalance, 2), 2)
		data.Balance = &v
	} else {
		data.Balance = &balance
	}
	err = data.Update()
	if err != nil {
		err = errors.Wrap(err, "更新失败")
	}
	return
}

func DelAccountWithDetails(id uint, uid uint) (err error) {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = errors.New("账户删除失败，数据已回滚")
		}
	}()
	err = tx.Error
	if err != nil {
		return
	}
	err = tx.Where("account_id = ? AND user_id = ?", id, uid).Delete(model.Detail{}).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Where("id = ? AND user_id = ?", id, uid).Delete(model.Account{}).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit().Error
	if err != nil {
		return
	}
	return
}

func ListAccountsByUid(uid uint) (list []*AccountList, err error) {
	r, err := model.ListAccountsByUid(uid)
	if err != nil {
		return
	}
	for _, v := range r {
		list = append(list, &AccountList{
			Id:   v.ID,
			Name: v.Name,
			Icon: v.Icon,
		})
	}
	return
}

func GetManageList(uid uint) (list []*ManageList, err error) {
	accountList, err := model.ListAccountsByUid(uid)
	if err != nil {
		err = errors.Wrap(err, "查询账户列表失败")
		return
	}
	detailList, err := model.ListDetailsByUid(uid)
	if err != nil {
		err = errors.Wrap(err, "查询明细列表失败")
		return
	}
	for _, v := range accountList {
		var outMoney, incomeMoney float64
		for _, v1 := range detailList {
			if v1.AccountId == v.ID {
				switch v1.Direction {
				case 1:
					incomeMoney = util.FloatAdd(incomeMoney, *v1.Money, 2)
				case 2:
					outMoney = util.FloatAdd(outMoney, *v1.Money, 2)
				case 3:
					outMoney = util.FloatAdd(outMoney, *v1.Money, 2)
				default:
					fmt.Println("unknown")
				}
			}
			if v1.IncomeAccountId == v.ID {
				incomeMoney = util.FloatAdd(incomeMoney, *v1.Money, 2)
			}
		}
		balance := util.FloatSub(util.FloatAdd(*v.Balance, incomeMoney, 2), outMoney, 2)
		obj := ManageList{
			Id:       v.ID,
			Name:     v.Name,
			Balance:  balance,
			Icon:     v.Icon,
			AddTotal: *v.AddTotal,
		}
		list = append(list, &obj)
	}
	return
}
