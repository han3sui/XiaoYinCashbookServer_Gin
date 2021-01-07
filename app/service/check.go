package service

import (
	"xiaoyin/app/model"
	"xiaoyin/lib/db"
)

type CheckSearchParams struct {
	PageNo   int `form:"page_no"`
	PageSize int `form:"page_size"`
}

func ListCheck(uid uint, params CheckSearchParams) (list []model.Check, err error) {
	err = db.DB.Scopes(model.Paginate(params.PageNo, params.PageSize)).Where("user_id = ?", uid).Order("create_time desc").Find(&list).Error
	return
}
