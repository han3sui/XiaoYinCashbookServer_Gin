package account

import (
	"xiaoyin/app/model"
	"xiaoyin/lib/db"
)

type Account struct {
	model.BaseModel
	Name     string   `json:"name" validate:"required" label:"账户名称"`
	UserId   uint     `json:"user_id"`
	Sort     uint     `json:"sort"`
	Balance  *float64 `json:"balance" validate:"required" label:"账户余额"`
	Icon     string   `json:"icon" validate:"required" label:"账户图标"`
	AddTotal *uint    `json:"add_total"`
}

//TableName of GORM model
func (m *Account) TableName() string {
	return "account"
}

func List(uid uint) (list []Account, err error) {
	err = db.DB.Where("user_id = ?", uid).Find(&list).Error
	return
}

func Save(data *Account) (id uint, err error) {
	err = db.DB.Create(&data).Error
	id = data.ID
	return
}

func Update(data *Account) (err error) {
	err = db.DB.Model(&data).Omit("id", "create_time").Updates(data).Error
	return
}

func Del(id uint, uid uint) (err error) {
	err = db.DB.Where("user_id = ? AND id = ?", uid, id).Delete(Account{}).Error
	return
}

func CheckExist(uid uint, name string) (id uint, balance float64, err error) {
	var account Account
	account.Balance = new(float64)
	err = db.DB.Where("user_id = ? AND name = ?", uid, name).Find(&account).Error
	id = account.ID
	balance = *account.Balance
	return
}
