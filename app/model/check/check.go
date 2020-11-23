package check

import "xiaoyin/app/model"

type Check struct {
	model.BaseModel
	UserId uint   `json:"user_id"`
	Data   string `json:"data"`
}

//TableName of GORM model
func (m *Check) TableName() string {
	return "check"
}
