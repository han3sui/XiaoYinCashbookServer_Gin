package detail

import (
	"strconv"
	"xiaoyin/app/model"
	"xiaoyin/app/model/account"
	"xiaoyin/app/model/category"
	"xiaoyin/lib/db"
	"xiaoyin/lib/util"
)

type Detail struct {
	model.BaseModel
	Money           *float64          `json:"money"`
	UserId          uint              `json:"user_id"`
	AccountId       uint              `json:"account_id"`
	IncomeAccountId uint              `json:"income_account_id"`
	CategoryId      uint              `json:"category_id"`
	Time            string            `json:"time"`
	Remark          *string           `json:"remark"`
	Direction       uint              `json:"direction"`
	Claim           *int              `json:"claim"`
	Category        category.Category `json:"-" gorm:"foreignKey:ID;references:CategoryId"`
	Account         account.Account   `json:"-" gorm:"foreignKey:ID;references:AccountId"`
	IncomeAccount   account.Account   `json:"-" gorm:"foreignKey:ID;references:IncomeAccountId"`
}

//TableName of GORM model
func (m *Detail) TableName() string {
	return "detail"
}

type SearchParams struct {
	AccountId  int    `json:"account_id" form:"account_id"`
	CategoryId int    `json:"category_id" form:"category_id"`
	Remark     string `json:"remark" form:"remark"`
	Year       int    `json:"year" form:"year"`
	Month      int    `json:"month" form:"month"`
	CheckTime  int    `json:"check_time" form:"check_time"`
	Claim      int    `json:"claim" form:"claim"`
	Id         uint   `json:"id"`
	PageNo     int    `json:"page_no" form:"page_no"`
	PageSize   int    `json:"page_size" form:"page_size"`
}

type ListMoney struct {
	Direction int     `json:"direction"`
	Total     float64 `json:"total"`
}

func ListByUid(uid uint) (list []Detail, err error) {
	err = db.DB.Where("user_id = ?", uid).Find(&list).Error
	return
}

func ListClaim(uid uint, claim int) (list []Detail, err error) {
	err = db.DB.Where("user_id = ? AND claim = ?", uid, claim).Preload("Category").Preload("Account").Find(&list).Error
	return
}

func ListMoneyByParams(uid uint, params SearchParams) (data []ListMoney, err error) {
	var sdb = db.DB
	sdb = sdb.Where("user_id = ?", uid)
	if params.AccountId != 0 {
		sdb = sdb.Where(db.DB.Where("account_id = ?", params.AccountId).Or("income_account_id = ?", params.AccountId))
	}
	if params.Year != 0 {
		year := strconv.Itoa(params.Year)
		var timeStart string
		var timeEnd string
		if params.Month != 0 {
			date := util.GetMonthStartAndEnd(year, strconv.Itoa(params.Month))
			timeStart = date["start"]
			timeEnd = date["end"]
		} else {
			timeStart = util.GetMonthStartAndEnd(year, "1")["start"]
			timeEnd = util.GetMonthStartAndEnd(year, "12")["end"]
		}
		sdb.Where("time >= ? AND time <= ?", timeStart, timeEnd)
	}
	err = sdb.Model(&Detail{}).Select("direction, sum(money) as total").Group("direction").Find(&data).Error
	return
}

func ListByParams(uid uint, params SearchParams, paginate bool) (list []Detail, err error) {
	var sdb = db.DB
	sdb = sdb.Where("user_id = ?", uid)
	if params.Remark != "" {
		sdb = sdb.Where("remark LIKE ?", "%"+params.Remark+"%")
	}
	if params.AccountId != 0 {
		sdb = sdb.Where(db.DB.Where("account_id = ?", params.AccountId).Or("income_account_id = ?", params.AccountId))
	}
	if params.CategoryId != 0 {
		sdb = sdb.Where("category_id = ?", params.CategoryId)
	}
	if params.Year != 0 {
		year := strconv.Itoa(params.Year)
		var timeStart string
		var timeEnd string
		if params.Month != 0 {
			date := util.GetMonthStartAndEnd(year, strconv.Itoa(params.Month))
			timeStart = date["start"]
			timeEnd = date["end"]
		} else {
			timeStart = util.GetMonthStartAndEnd(year, "1")["start"]
			timeEnd = util.GetMonthStartAndEnd(year, "12")["end"]
		}
		sdb.Where("time >= ? AND time <= ?", timeStart, timeEnd)
	}
	if params.CheckTime != 0 {
		sdb = sdb.Where("update_time >= ?", params.CheckTime)
	}
	if params.Claim != 0 {
		sdb = sdb.Where("claim = ?", params.Claim)
	}
	if params.Id != 0 {
		sdb = sdb.Where("id = ?", params.Id)
	}
	if paginate {
		err = sdb.Scopes(model.Paginate(params.PageNo, params.PageSize)).Preload("Category").Preload("Account").Preload("IncomeAccount").Order("time desc,create_time desc").Find(&list).Error
	} else {
		err = sdb.Preload("Category").Preload("Account").Preload("IncomeAccount").Order("time desc,create_time desc").Find(&list).Error
	}
	return
}

func Save(data *Detail) (id uint, err error) {
	err = db.DB.Omit("Category", "Account", "IncomeAccount").Create(&data).Error
	id = data.ID
	return
}

func Update(data *Detail) (err error) {
	err = db.DB.Omit("Category", "Account", "IncomeAccount", "id", "create_time").Updates(&data).Error
	return
}

func Del(id uint, uid uint) (err error) {
	err = db.DB.Where("id = ? AND user_id = ?", id, uid).Delete(Detail{}).Error
	return
}
