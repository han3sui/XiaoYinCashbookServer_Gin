package model

import (
	"xiaoyin/lib/db"
)

type Icon struct {
	BaseModel
	Name string `json:"name"`
}

//TableName of GORM model
func (m *Icon) TableName() string {
	return "icon"
}

func ListIcons() (list []Icon, err error) {
	err = db.DB.Find(&list).Error
	return
}
