package service

import (
	"database/sql"
	"encoding/json"
	"time"
	"xiaoyin/app/model"
	"xiaoyin/lib/db"
)

type UserInfo struct {
	ID        uint      `json:"id"`
	NickName  string    `json:"nick_name"`
	AvatarUrl string    `json:"avatar_url"`
	CheckTime CheckTime `json:"check_time"`
}

type CheckInfo = model.Check


type CheckTime sql.NullInt64

func (v CheckTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(0)
	}
}

func UpdateCheckTime(uid uint, data *CheckInfo) (checkTime int64, err error) {
	checkTime = time.Now().Unix()
	err = db.DB.Model(&model.User{}).Where("id = ?", uid).Update("check_time", checkTime).Error
	data.UserId = uid
	_ = db.DB.Create(&data)
	return
}

func (v CheckTime) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

func SaveOrUpdate(data *model.User) (id uint, err error) {
	id, err = model.CheckUserExist(data.Openid)
	if err != nil {
		return
	}
	if id == 0 {
		id, err = data.Save()
	} else {
		data.ID=id
		err = data.Update()
	}
	return
}

func GetInfo(uid uint) (data *UserInfo, err error) {
	r, err := model.GetUserById(uid)
	if err != nil {
		return
	}
	data=new(UserInfo)
	data.ID = r.ID
	data.AvatarUrl = r.AvatarUrl
	data.NickName = r.NickName
	data.CheckTime = CheckTime(r.CheckTime)
	return
}
