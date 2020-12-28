package service

import (
	"xiaoyin/app/model"
	"xiaoyin/lib/db"
)

//type Info = model.Check

func ListCheck(uid uint) (list []model.Check, err error) {
	err = db.DB.Where("user_id = ?", uid).Order("create_time desc").Find(&list).Error
	return
}