package icon

import (
	"xiaoyin/app/model"
	"xiaoyin/lib/db"
)

type Icon struct {
	model.BaseModel
	Name string `json:"name"`
}

//TableName of GORM model
func (m *Icon) TableName() string {
	return "icon"
}

func List() (list []Icon, err error) {
	err = db.DB.Find(&list).Error
	return
}
