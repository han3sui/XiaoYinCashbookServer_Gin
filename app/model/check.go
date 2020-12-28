package model

type Check struct {
	BaseModel
	UserId uint   `json:"user_id"`
	Data   string `json:"data"`
}

//TableName of GORM model
func (m *Check) TableName() string {
	return "check"
}
