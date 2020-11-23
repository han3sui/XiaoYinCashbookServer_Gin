package user

import (
	"database/sql"
	"xiaoyin/app/model"
	"xiaoyin/app/model/account"
	"xiaoyin/app/model/category"
	"xiaoyin/lib/db"

	"github.com/pkg/errors"
)

type User struct {
	model.BaseModel
	NickName   string        `json:"nick_name"`
	Openid     string        `json:"openid"`
	SessionKey string        `json:"session_key"`
	AvatarUrl  string        `json:"avatar_url"`
	CheckTime  sql.NullInt64 `json:"check_time"`
}

//TableName of GORM model
func (m *User) TableName() string {
	return "user"
}

//账户初始化数据
var InitAccountData = []account.Account{{
	Name:    "支付宝",
	Balance: new(float64),
	Icon:    "icon-account-zhifubao",
}, {
	Name:    "微信",
	Balance: new(float64),
	Icon:    "icon-account-weixin",
}, {
	Name:    "现金",
	Balance: new(float64),
	Icon:    "icon-account-xianjin",
}}

//定义初始化数据结构
type InitCategory struct {
	MainCategory category.Category
	Nodes        []category.Category
}

//分类初始化数据
var InitCategoryData = []InitCategory{{
	MainCategory: category.Category{
		Name:      "餐饮",
		Direction: 2,
	},
	Nodes: []category.Category{{
		Name:      "买菜",
		Direction: 2,
		Icon:      "icon-canyin-shucai",
	}, {
		Name:      "水果",
		Direction: 2,
		Icon:      "icon-canyin-shuiguo",
	}, {
		Name:      "外卖",
		Direction: 2,
		Icon:      "icon-canyin-waimai",
	}, {
		Name:      "餐厅",
		Direction: 2,
		Icon:      "icon-canyin-canting",
	}, {
		Name:      "零食",
		Direction: 2,
		Icon:      "icon-canyin-lingshi",
	}},
}, {
	MainCategory: category.Category{
		Name:      "生活缴费",
		Direction: 2,
	},
	Nodes: []category.Category{{
		Name:      "水费",
		Direction: 2,
		Icon:      "icon-shenghuo-shuifei",
	}, {
		Name:      "电费",
		Direction: 2,
		Icon:      "icon-shenghuo-dianfei",
	}, {
		Name:      "燃气费",
		Direction: 2,
		Icon:      "icon-shenghuo-ranqifei",
	}, {
		Name:      "宽带话费",
		Direction: 2,
		Icon:      "icon-shenghuo-huafei",
	}, {
		Name:      "物业费",
		Direction: 2,
		Icon:      "icon-account-xianjin",
	}, {
		Name:      "银行费用",
		Direction: 2,
		Icon:      "icon-account-xianjin",
	}},
}, {
	MainCategory: category.Category{
		Name:      "车辆费用",
		Direction: 2,
	},
	Nodes: []category.Category{{
		Name:      "加油",
		Direction: 2,
		Icon:      "icon-cheliang-jiayou",
	}, {
		Name:      "维修",
		Direction: 2,
		Icon:      "icon-cheliang-weixiu",
	}, {
		Name:      "保养",
		Direction: 2,
		Icon:      "icon-cheliang-baoyang",
	}, {
		Name:      "汽车用品",
		Direction: 2,
		Icon:      "icon-cheliang-baoyang",
	}, {
		Name:      "停车费",
		Direction: 2,
		Icon:      "icon-chuxing-tingche",
	}, {
		Name:      "车险",
		Direction: 2,
		Icon:      "icon-yiliao-baoxian",
	}, {
		Name:      "违章",
		Direction: 2,
		Icon:      "icon-chuxing-weizhang",
	}, {
		Name:      "高速费",
		Direction: 2,
		Icon:      "icon-chuxing-gaosu",
	}},
}, {
	MainCategory: category.Category{
		Name:      "日用品",
		Direction: 2,
	},
	Nodes: []category.Category{{
		Name:      "纸品",
		Direction: 2,
		Icon:      "icon-riyong-zhipin",
	}, {
		Name:      "调料",
		Direction: 2,
		Icon:      "icon-riyong-tiaoliao",
	}, {
		Name:      "餐具",
		Direction: 2,
		Icon:      "icon-riyong-canju",
	}, {
		Name:      "收纳物件",
		Direction: 2,
		Icon:      "icon-riyong-shouna",
	}, {
		Name:      "清洁用品",
		Direction: 2,
		Icon:      "icon-riyong-qingjie",
	}},
}, {
	MainCategory: category.Category{
		Name:      "出行",
		Direction: 2,
	},
	Nodes: []category.Category{{
		Name:      "交通费",
		Direction: 2,
		Icon:      "icon-chuxing-dache",
	}, {
		Name:      "旅游",
		Direction: 2,
		Icon:      "icon-chuxing-lvyou",
	}, {
		Name:      "酒店",
		Direction: 2,
		Icon:      "icon-chuxing-jiudian",
	}},
}, {
	MainCategory: category.Category{
		Name:      "人情",
		Direction: 2,
	},
	Nodes: []category.Category{{
		Name:      "人情红包",
		Direction: 2,
		Icon:      "icon-renqing-yasuiqian",
	}, {
		Name:      "伴手礼",
		Direction: 2,
		Icon:      "icon-renqing-banshouli",
	}, {
		Name:      "宴请招待",
		Direction: 2,
		Icon:      "icon-renqing-yanqing",
	}},
}, {
	MainCategory: category.Category{
		Name:      "医疗健康",
		Direction: 2,
	},
	Nodes: []category.Category{{
		Name:      "医院",
		Direction: 2,
		Icon:      "icon-yiliao-yiyuan",
	}, {
		Name:      "药店",
		Direction: 2,
		Icon:      "icon-yiliao-yaodian",
	}, {
		Name:      "保健品",
		Direction: 2,
		Icon:      "icon-yiliao-baojianpin",
	}, {
		Name:      "保险",
		Direction: 2,
		Icon:      "icon-yiliao-baoxian",
	}},
}, {
	MainCategory: category.Category{
		Name:      "服饰",
		Direction: 2,
	},
	Nodes: []category.Category{{
		Name:      "衣",
		Direction: 2,
		Icon:      "icon-fushi-yifu",
	}, {
		Name:      "裤",
		Direction: 2,
		Icon:      "icon-fushi-kuzi",
	}, {
		Name:      "鞋",
		Direction: 2,
		Icon:      "icon-fushi-xiezi",
	}, {
		Name:      "配件",
		Direction: 2,
		Icon:      "icon-fushi-peijian",
	}},
}, {
	MainCategory: category.Category{
		Name:      "贷款",
		Direction: 2,
	},
	Nodes: []category.Category{{
		Name:      "房贷",
		Direction: 2,
		Icon:      "icon-daikuan-fangdai",
	}, {
		Name:      "车贷",
		Direction: 2,
		Icon:      "icon-daikuan-cheweidai",
	}, {
		Name:      "借款",
		Direction: 2,
		Icon:      "icon-daikuan-jiekuan",
	}},
}, {
	MainCategory: category.Category{
		Name:      "职业收入",
		Direction: 1,
	},
	Nodes: []category.Category{{
		Name:      "工资收入",
		Direction: 1,
		Icon:      "icon-shouru-gongzi",
	}},
}, {
	MainCategory: category.Category{
		Name:      "其他收入",
		Direction: 1,
		Icon:      "icon-shouru-gongzi",
	},
	Nodes: []category.Category{{
		Name:      "红包",
		Direction: 1,
		Icon:      "icon-shouru-hongbao",
	}, {
		Name:      "收益",
		Direction: 1,
		Icon:      "icon-shouru-shouyi",
	}},
}}

func Save(data *User) (id uint, err error) {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = errors.New("用户初始化失败，数据已回滚")
		}
	}()
	err = tx.Error
	if err != nil {
		return
	}
	err = tx.Create(&data).Error
	if err != nil {
		tx.Rollback()
		return
	}
	id = data.ID
	var accountData []account.Account
	for _, v := range InitAccountData {
		accountData = append(accountData, account.Account{
			Name:    v.Name,
			Icon:    v.Icon,
			Balance: v.Balance,
			UserId:  id,
		})
	}
	err = tx.Create(&accountData).Error
	if err != nil {
		tx.Rollback()
		return
	}
	var mainCategory []category.Category
	for _, v := range InitCategoryData {
		mainCategory = append(mainCategory, category.Category{
			Name:      v.MainCategory.Name,
			Direction: v.MainCategory.Direction,
			UserId:    id,
		})
	}
	err = tx.Create(&mainCategory).Error
	if err != nil {
		tx.Rollback()
		return
	}
	for k, v := range mainCategory {
		InitCategoryData[k].MainCategory.ID = v.ID
	}
	var nodeCategory []category.Category
	for k, v := range InitCategoryData {
		for _, v1 := range InitCategoryData[k].Nodes {
			nodeCategory = append(nodeCategory, category.Category{
				Name:      v1.Name,
				Direction: v1.Direction,
				UserId:    id,
				ParentId:  v.MainCategory.ID,
				Icon:      v1.Icon,
			})
		}
	}
	err = tx.Create(&nodeCategory).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit().Error
	return
}

func Update(id uint, data *User) (err error) {
	err = db.DB.Model(&data).Where("id = ?", id).Updates(data).Error
	return
}

func Info(uid uint) (data User, err error) {
	err = db.DB.Where("id = ?", uid).Find(&data).Error
	return
}

func CheckExist(openid string) (id uint, err error) {
	var user User
	r := db.DB.Where("openid = ?", openid).Find(&user)
	if r.Error != nil {
		err = r.Error
		return
	}
	id = user.ID
	return
}
