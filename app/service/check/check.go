package check

import (
	"xiaoyin/app/model/check"
	"xiaoyin/lib/db"
)

type Info = check.Check

func List(uid uint) (list []Info, err error) {
	err = db.DB.Where("user_id = ?", uid).Order("create_time desc").Find(&list).Error
	return
}
