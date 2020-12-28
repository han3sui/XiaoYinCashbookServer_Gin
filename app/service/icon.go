package service

import (
	"xiaoyin/app/model"
)

func ListIcons() (list []model.Icon, err error) {
	list, err = model.ListIcons()
	return
}
