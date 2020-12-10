package account

import (
	"fmt"
	"xiaoyin/app/model/account"
	"xiaoyin/app/model/detail"
	"xiaoyin/lib/db"

	"github.com/shopspring/decimal"

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
type List struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}
type Info = account.Account

func Save(data *Info) (id uint, err error) {
	id, _, err = account.CheckExist(data.UserId, data.Name)
	if err != nil {
		err = errors.Wrap(err, "账户重复检查失败")
		return
	}
	if id != 0 {
		err = errors.New("该账户已存在")
		return
	}
	id, err = account.Save(data)
	if err != nil {
		err = errors.Wrap(err, "账户保存失败")
	}
	return
}

func GetDetailsCount(uid uint, id uint) (count int64, err error) {
	err = db.DB.Model(&detail.Detail{}).Where("user_id = ? AND account_id = ?", uid, id).Count(&count).Error
	return
}

func Update(tmpBalance float64, data *Info) (err error) {
	id, balance, err := account.CheckExist(data.UserId, data.Name)
	if err != nil {
		err = errors.Wrap(err, "账户重复检查失败")
		return
	}
	if id != 0 && id != data.ID {
		err = errors.New("该账户已存在")
		return
	}
	if tmpBalance != *data.Balance {
		v, _ := decimal.NewFromFloatWithExponent(balance, -2).Add(decimal.NewFromFloatWithExponent(*data.Balance, -2).Sub(decimal.NewFromFloatWithExponent(tmpBalance, -2))).Float64()
		data.Balance = &v
	} else {
		data.Balance = &balance
	}
	err = account.Update(data)
	if err != nil {
		err = errors.Wrap(err, "更新失败")
	}
	return
}

func DelWithDetails(id uint, uid uint) (err error) {
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
	err = tx.Where("account_id = ? AND user_id = ?", id, uid).Delete(detail.Detail{}).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Where("id = ? AND user_id = ?", id, uid).Delete(account.Account{}).Error
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

func ListByUid(uid uint) (list []*List, err error) {
	r, err := account.List(uid)
	if err != nil {
		return
	}
	for _, v := range r {
		list = append(list, &List{
			Id:   v.ID,
			Name: v.Name,
			Icon: v.Icon,
		})
	}
	return
}

func GetManageList(uid uint) (list []*ManageList, err error) {
	accountList, err := account.List(uid)
	if err != nil {
		err = errors.Wrap(err, "查询账户列表失败")
		return
	}
	detailList, err := detail.ListByUid(uid)
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
					incomeMoney += v1.Money
				case 2:
					outMoney += v1.Money
				case 3:
					outMoney += v1.Money
				default:
					fmt.Println("unknown")
				}
			}
			if v1.IncomeAccountId == v.ID {
				incomeMoney += v1.Money
			}
		}
		balance, _ := (decimal.NewFromFloatWithExponent(*v.Balance, -2).Add(decimal.NewFromFloatWithExponent(incomeMoney, -2))).Sub(decimal.NewFromFloatWithExponent(outMoney, -2)).Float64()
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
