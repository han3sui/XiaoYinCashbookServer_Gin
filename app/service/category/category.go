package category

import (
	"xiaoyin/app/model/category"
	"xiaoyin/app/model/detail"
	"xiaoyin/lib/db"

	"github.com/pkg/errors"
)

type Tree struct {
	Id       uint    `json:"id"`
	ParentId uint    `json:"parent_id"`
	Name     string  `json:"name"`
	Sort     uint    `json:"sort"`
	Icon     string  `json:"icon"`
	Nodes    []*Tree `json:"nodes"`
}

type ParentDetail struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type Info = category.Category

//根据分类ID查找关联明细条数
func GetDetailsCount(uid uint, id uint) (count int64, err error) {
	err = db.DB.Model(&detail.Detail{}).Where("user_id = ? AND category_id = ?", uid, id).Count(&count).Error
	return
}

func ListTreeByUid(uid uint) (list map[string][]*Tree, err error) {
	r, err := category.ListByUid(uid)
	if err != nil {
		return
	}
	incomeList := make([]Info, 0, len(r))
	outList := make([]Info, 0, len(r))
	for _, v := range r {
		if v.Direction == 1 {
			incomeList = append(incomeList, v)
		}
		if v.Direction == 2 {
			outList = append(outList, v)
		}
	}
	list = make(map[string][]*Tree)
	list["income"] = getTree(incomeList, 0)
	list["out"] = getTree(outList, 0)
	return
}

func ListByUid(uid uint) (list []Info, err error) {
	list, err = category.ListByUid(uid)
	if err != nil {
		return
	}
	return
}

func GetParent(id uint, list []Info) (data ParentDetail) {
	for _, v := range list {
		if v.ID == id {
			if v.ParentId == 0 {
				data = ParentDetail{
					Id:   v.ID,
					Name: v.Name,
					Icon: v.Icon,
				}
				break
			} else {
				data = GetParent(v.ParentId, list)
			}
		}
	}
	return
}

//根据分类ID查找分类详情
func GetDetail(id uint, list []Info) (data Info) {
	for _, v := range list {
		if v.ID == id {
			data = v
			break
		}
	}
	return
}

func Save(data *Info) (id uint, err error) {
	id, err = category.CheckExist(data.UserId, data.Name)
	if err != nil {
		err = errors.Wrap(err, "分类重复检查失败")
		return
	}
	if id != 0 {
		err = errors.New("该分类已存在")
		return
	}
	id, err = category.Save(data)
	return
}

func Update(data *Info) (err error) {
	id, err := category.CheckExist(data.UserId, data.Name)
	if err != nil {
		err = errors.Wrap(err, "分类重复检查失败")
		return
	}
	if id != 0 && id != data.ID {
		err = errors.New("该分类已存在")
		return
	}
	err = category.Update(data)
	return
}

func DelWithDetails(id uint, uid uint) (err error) {
	boolean, err := category.CheckChildren(id, uid)
	if err != nil {
		err = errors.Wrap(err, "检查子分类失败")
		return
	}
	if boolean {
		err = errors.New("请先删除子分类")
		return
	}
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = errors.New("分类删除失败，数据已回滚")
		}
	}()
	err = tx.Error
	if err != nil {
		return
	}
	err = tx.Where("category_id = ? AND user_id = ?", id, uid).Delete(detail.Detail{}).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Where("id = ? AND user_id = ?", id, uid).Delete(category.Category{}).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit().Error
	if err != nil {
		return
	}
	return
}

func getTree(list []Info, pid uint) (tree []*Tree) {
	for _, v := range list {
		if v.ParentId == pid {
			nodes := getTree(list, v.ID)
			obj := Tree{
				Id:       v.ID,
				ParentId: v.ParentId,
				Name:     v.Name,
				Sort:     v.Sort,
				Icon:     v.Icon,
				Nodes:    nodes,
			}
			tree = append(tree, &obj)
		}
	}
	return
}
