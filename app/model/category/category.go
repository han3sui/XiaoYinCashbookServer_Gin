package category

import (
	"xiaoyin/app/model"
	"xiaoyin/lib/db"
)

type Category struct {
	model.BaseModel
	Name      string `json:"name" validate:"required" label:"分类名称"`
	Sort      uint   `json:"sort"`
	Direction uint   `json:"direction" validate:"required" label:"分类类型"`
	UserId    uint   `json:"user_id"`
	ParentId  uint   `json:"parent_id"`
	Icon      string `json:"icon"`
}

//TableName of GORM model
func (m *Category) TableName() string {
	return "category"
}

func Save(data *Category) (id uint, err error) {
	err = db.DB.Create(&data).Error
	id = data.ID
	return
}

func Update(data *Category) (err error) {
	err = db.DB.Model(&data).Omit("id", "create_time", "user_id").Updates(data).Error
	return
}

func ListByUid(uid uint) (list []Category, err error) {
	err = db.DB.Where("user_id = ?", uid).Find(&list).Error
	return
}

func Del(id uint, uid uint) (err error) {
	err = db.DB.Where("user_id = ? AND id = ?", uid, id).Delete(Category{}).Error
	return
}

func CheckExist(uid uint, name string) (id uint, err error) {
	var category Category
	err = db.DB.Where("user_id = ? AND name = ?", uid, name).Find(&category).Error
	id = category.ID
	return
}

func CheckChildren(id uint, uid uint) (boolean bool, err error) {
	var category Category
	err = db.DB.Where("user_id = ? AND parent_id = ?", uid, id).Find(&category).Error
	boolean = category.ID != 0
	return
}
