package service

import (
	"sort"
	"strconv"
	"xiaoyin/app/model"
	"xiaoyin/lib/db"
	"xiaoyin/lib/util"

	"github.com/pkg/errors"
)

type Detail = model.Detail

//搜索条件
type DetailSearchParams = model.SearchParams

//图表数据实体
type ChartData struct {
	TotalOut    float64           `json:"total_out"`
	TotalIncome float64           `json:"total_income"`
	Out         []ChartDataDetail `json:"out"`
	Income      []ChartDataDetail `json:"income"`
}

type ChartDataDetail struct {
	Id      uint              `json:"id"`
	Name    string            `json:"name"`
	Money   float64           `json:"money"`
	Percent float64           `json:"percent"`
	Icon    string            `json:"icon"`
	Nodes   []ChartDataDetail `json:"nodes"`
}

//账单数据实体
type BillData struct {
	TotalOut    float64         `json:"total_out"`
	TotalIncome float64         `json:"total_income"`
	Data        []BillDataMonth `json:"data"`
}

type BillDataMonth struct {
	Month       string  `json:"month"`
	OutMoney    float64 `json:"out_money"`
	IncomeMoney float64 `json:"income_money"`
}

//明细列表实体
type TotalData struct {
	TotalOut    float64 `json:"total_out"`
	TotalIncome float64 `json:"total_income"`
	Data        []Data  `json:"data"`
}

type Data struct {
	Time   string  `json:"time"`
	Income float64 `json:"income"`
	Out    float64 `json:"out"`
	List   []List  `json:"list"`
}

type List struct {
	Id            uint         `json:"id"`
	Money         float64      `json:"money"`
	Time          string       `json:"time"`
	Remark        string       `json:"remark"`
	Direction     uint         `json:"direction"`
	Claim         uint         `json:"claim"`
	UpdateTime    int64        `json:"update_time"`
	Account       accountInfo  `json:"account"`
	IncomeAccount accountInfo  `json:"income_account"`
	Category      categoryInfo `json:"category"`
}

type categoryInfo struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type accountInfo struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

//报销更新
type ClaimDetail struct {
	Id              uint `json:"id"`
	Claim           uint `json:"claim"`
	IncomeAccountId uint `json:"income_account_id"`
}

//金额总额
type ListMoney struct {
	Total     float64 `json:"total"`
	Direction int     `json:"direction"`
}

//根据月份获取总支出/总收入
func ListMoneyByParams(uid uint, params DetailSearchParams) (data [2]ListMoney, err error) {
	r, err := model.ListAllDetailsByParams(uid, params)
	if err != nil {
		err = errors.Wrap(err, "获取总额失败")
		return
	}
	for _, v := range r {
		if *v.Claim != 2 {
			if v.Direction == 1 {
				data[0].Direction = 1
				data[0].Total = util.FloatAdd(data[0].Total, *v.Money, 2)
			}
			if v.Direction == 2 {
				data[1].Direction = 2
				data[1].Total = util.FloatAdd(data[1].Total, *v.Money, 2)
			}
		}
	}
	return
}

//条件查询
func ListByParams(uid uint, params DetailSearchParams) (list []*List, err error) {
	r, err := model.ListDetailsByParams(uid, params, true)
	if err != nil {
		return
	}
	for _, v := range r {
		list = append(list, &List{
			Id:         v.ID,
			Money:      *v.Money,
			Time:       v.Time,
			Remark:     *v.Remark,
			Direction:  v.Direction,
			Claim:      *v.Claim,
			UpdateTime: v.UpdateTime,
			Account: accountInfo{
				Id:   v.Account.ID,
				Name: v.Account.Name,
				Icon: v.Account.Icon,
			},
			IncomeAccount: accountInfo{
				Id:   v.IncomeAccount.ID,
				Name: v.IncomeAccount.Name,
				Icon: v.IncomeAccount.Icon,
			},
			Category: categoryInfo{
				Id:   v.Category.ID,
				Name: v.Category.Name,
				Icon: v.Category.Icon,
			},
		})
	}
	return
}

//查询报销账单
func ListClaim(uid uint, claim int) (list []*List, err error) {
	r, err := model.ListClaim(uid, claim)
	if err != nil {
		return
	}
	for _, v := range r {
		list = append(list, &List{
			Id:         v.ID,
			Money:      *v.Money,
			Time:       v.Time,
			Remark:     *v.Remark,
			Claim:      *v.Claim,
			Direction:  v.Direction,
			UpdateTime: v.UpdateTime,
			Account: accountInfo{
				Id:   v.Account.ID,
				Name: v.Account.Name,
				Icon: v.Account.Icon,
			},
			Category: categoryInfo{
				Id:   v.Category.ID,
				Name: v.Category.Name,
				Icon: v.Category.Icon,
			},
		})
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Time > list[j].Time
	})
	return
}

//获取 年度/月份 账单数据
func Bill(uid uint, params DetailSearchParams) (list *BillData, err error) {
	r, err := model.ListDetailsByParams(uid, params, false)
	if err != nil {
		return
	}
	list = new(BillData)
	months := [12]string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}
	data := make(map[string]*BillDataMonth)
	for _, v := range r {
		if *v.Claim != 2 {
			if v.Direction == 2 {
				list.TotalOut = util.FloatAdd(list.TotalOut, *v.Money, 2)
			}
			if v.Direction == 1 {
				list.TotalIncome = util.FloatAdd(list.TotalIncome, *v.Money, 2)
			}
			for _, v1 := range months {
				_, ok := data[v1]
				if !ok {
					data[v1] = &BillDataMonth{
						Month: v1,
					}
				}
				if v.Time[:7] == (strconv.Itoa(params.Year) + "-" + v1) {
					if v.Direction == 2 {
						data[v1].OutMoney = util.FloatAdd(data[v1].OutMoney, *v.Money, 2)
					}
					if v.Direction == 1 {
						data[v1].IncomeMoney = util.FloatAdd(data[v1].IncomeMoney, *v.Money, 2)
					}
				}
			}
		}
	}
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	for _, v := range keys {
		list.Data = append(list.Data, *data[v])
	}
	return
}

//获取明细所属父分类的切片下标
func GetDetailKey(id uint, list []ChartDataDetail) (key *int) {
	for k, v := range list {
		if id == v.Id {
			key = &k
			break
		}
	}
	return
}

//获取账单图表
func Chart(uid uint, params DetailSearchParams) (list *ChartData, err error) {
	r, err := model.ListDetailsByParams(uid, params, false)
	if err != nil {
		return
	}
	categoryList, err := ListCategorysByUid(uid)
	if err != nil {
		return
	}
	list = new(ChartData)
	var outSlice []Detail
	var incomeSlice []Detail
	for _, v := range r {
		if *v.Claim != 2 {
			if v.Direction == 2 {
				outSlice = append(outSlice, v)
				list.TotalOut = util.FloatAdd(list.TotalOut, *v.Money, 2)
			}
			if v.Direction == 1 {
				incomeSlice = append(incomeSlice, v)
				list.TotalIncome = util.FloatAdd(list.TotalIncome, *v.Money, 2)
			}
		}
	}
	list.Out = getChartTree(outSlice, categoryList, list.TotalOut)
	list.Income = getChartTree(incomeSlice, categoryList, list.TotalIncome)
	return
}

//生成图表需要的树结构
func getChartTree(data []Detail, categoryList []Category, totalMoney float64) (list []ChartDataDetail) {
	list = *new([]ChartDataDetail)
	for _, v := range data {
		//获取父分类详情
		parentDetail := GetParent(v.CategoryId, categoryList)
		//获取当前分类详情
		nodeDetail := GetDetail(v.CategoryId, categoryList)
		//获取父分类的切片下标
		key := GetDetailKey(parentDetail.Id, list)
		if key == nil {
			//如果该明细父分类没有在父分类切片中存在，则同时添加父分类和子分类
			list = append(list, ChartDataDetail{
				Id:      parentDetail.Id,
				Name:    parentDetail.Name,
				Money:   *v.Money,
				Percent: util.FloatMul(util.FloatDiv(*v.Money, totalMoney, 4), 100, 2),
				Nodes: []ChartDataDetail{
					{
						Id:      v.CategoryId,
						Name:    nodeDetail.Name,
						Icon:    nodeDetail.Icon,
						Money:   *v.Money,
						Percent: util.FloatMul(util.FloatDiv(*v.Money, totalMoney, 4), 100, 2),
					},
				},
			})
		} else {
			//如果该明细父分类已存在，则父分类的金额叠加，并重新计算百分比
			list[*key].Money = util.FloatAdd(list[*key].Money, *v.Money, 2)
			list[*key].Percent = util.FloatMul(util.FloatDiv(list[*key].Money, totalMoney, 4), 100, 2)
			//获取该明细在所属子分类中的切片下标
			key1 := GetDetailKey(v.CategoryId, list[*key].Nodes)
			if key1 == nil {
				//如果该明细不存在于所属子分类中，则添加子分类
				list[*key].Nodes = append(list[*key].Nodes, ChartDataDetail{
					Id:      v.CategoryId,
					Name:    nodeDetail.Name,
					Icon:    nodeDetail.Icon,
					Money:   *v.Money,
					Percent: util.FloatMul(util.FloatDiv(*v.Money, totalMoney, 4), 100, 2),
				})
			} else {
				//如果该明细存在于所属子分类，则子分类的金额叠加，并重新计算百分比
				list[*key].Nodes[*key1].Money = util.FloatAdd(list[*key].Nodes[*key1].Money, *v.Money, 2)
				list[*key].Nodes[*key1].Percent = util.FloatMul(util.FloatDiv(list[*key].Nodes[*key1].Money, totalMoney, 4), 100, 2)
			}
		}
	}
	//父分类切片按照金额重新排序
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Money > list[j].Money
	})
	//子分类切片按照金额重新排序
	for _, v := range list {
		sort.SliceStable(v.Nodes, func(i, j int) bool {
			return v.Nodes[i].Money > v.Nodes[j].Money
		})
	}
	return
}

func SaveDetail(data *Detail) (list []*List, err error) {
	id, err := data.Save()
	if err != nil {
		err = errors.Wrap(err, "账单保存失败")
		return
	}
	list, err = ListByParams(data.UserId, DetailSearchParams{Id: id})
	if err != nil {
		err = errors.Wrap(err, "查询账单失败")
		return
	}
	return
}

func UpdateDetail(data *Detail) (list []*List, err error) {
	err = data.Update()
	if err != nil {
		err = errors.Wrap(err, "账单更新失败")
		return
	}
	list, err = ListByParams(data.UserId, DetailSearchParams{Id: data.ID})
	if err != nil {
		err = errors.Wrap(err, "查询账单失败")
		return
	}
	return
}

func Del(id uint, uid uint) (err error) {
	err = model.DelDetailById(id, uid)
	if err != nil {
		err = errors.Wrap(err, "明细删除失败")
		return
	}
	return
}

//查询是否存在未核账明细
func IsExistUncheck(uid uint, checkTime uint) (boolean bool, err error) {
	var count int64
	err = db.DB.Model(&Detail{}).Where("user_id = ? AND update_time >= ?", uid, checkTime).Count(&count).Error
	if count == 0 {
		boolean = false
	} else {
		boolean = true
	}
	return
}

//获取总记账天数
func AllDaysCount(uid uint) (count int64, err error) {
	err = db.DB.Model(&Detail{}).Where("user_id = ?", uid).Group("time").Count(&count).Error
	return
}

//批量更新报销
func BatchUpdateClaim(uid uint, data []ClaimDetail) (err error) {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = errors.New("数据更新失败，已回滚")
		}
	}()
	err = tx.Error
	if err != nil {
		return
	}
	for _, v := range data {
		err = tx.Where("user_id = ?", uid).Where("id = ?", v.Id).Updates(&Detail{Claim: &v.Claim, IncomeAccountId: v.IncomeAccountId}).Error
		if err != nil {
			tx.Rollback()
			return
		}
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return
	}
	return
}
