package icon

import "xiaoyin/app/model/icon"

type Info struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

func List() (list []icon.Icon, err error) {
	list, err = icon.List()
	return
}
