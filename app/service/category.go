package service

import (
	"xiaoyin/app/model"
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

type Category = model.Category

//根据分类ID查找关联明细条数
func GetDetailsCountByCid(uid uint, id uint) (count int64, err error) {
	err = db.DB.Model(&model.Detail{}).Where("user_id = ? AND category_id = ?", uid, id).Count(&count).Error
	return
}

func ListTreeByUid(uid uint) (list map[string][]*Tree, err error) {
	r, err := model.ListCategorysByUid(uid)
	if err != nil {
		return
	}
	incomeList := make([]Category, 0, len(r))
	outList := make([]Category, 0, len(r))
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

func ListCategorysByUid(uid uint) (list []Category, err error) {
	list, err = model.ListCategorysByUid(uid)
	if err != nil {
		return
	}
	return
}

func GetParent(id uint, list []Category) (data ParentDetail) {
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
func GetDetail(id uint, list []Category) (data Category) {
	for _, v := range list {
		if v.ID == id {
			data = v
			break
		}
	}
	return
}

func Save(data *Category) (id uint, err error) {
	id, err = model.CheckCategoryExist(data.UserId, data.Name)
	if err != nil {
		err = errors.Wrap(err, "分类重复检查失败")
		return
	}
	if id != 0 {
		err = errors.New("该分类已存在")
		return
	}
	id, err = data.Save()
	return
}

func Update(data *Category) (err error) {
	id, err := model.CheckCategoryExist(data.UserId, data.Name)
	if err != nil {
		err = errors.Wrap(err, "分类重复检查失败")
		return
	}
	if id != 0 && id != data.ID {
		err = errors.New("该分类已存在")
		return
	}
	err = data.Update()
	return
}

func DelCategoryWithDetails(id uint, uid uint) (err error) {
	boolean, err := model.CheckCategoryChildren(id, uid)
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
	err = tx.Where("category_id = ? AND user_id = ?", id, uid).Delete(model.Detail{}).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Where("id = ? AND user_id = ?", id, uid).Delete(model.Category{}).Error
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

func getTree(list []Category, pid uint) (tree []*Tree) {
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
