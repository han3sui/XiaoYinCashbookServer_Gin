package user

import (
	"database/sql"
	"encoding/json"
	"time"
	"xiaoyin/app/model/check"
	"xiaoyin/app/model/user"
	"xiaoyin/lib/db"
)

type UserInfo struct {
	ID        uint      `json:"id"`
	NickName  string    `json:"nick_name"`
	AvatarUrl string    `json:"avatar_url"`
	CheckTime CheckTime `json:"check_time"`
}

type CheckInfo = check.Check

type Info = user.User

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
	err = db.DB.Model(&Info{}).Where("id = ?", uid).Update("check_time", checkTime).Error
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

func SaveOrUpdate(data *Info) (id uint, err error) {
	id, err = user.CheckExist(data.Openid)
	if err != nil {
		return
	}
	if id == 0 {
		id, err = user.Save(data)
	} else {
		err = user.Update(id, data)
	}
	return
}

func GetInfo(uid uint) (data UserInfo, err error) {
	r, err := user.Info(uid)
	if err != nil {
		return
	}
	data.ID = r.ID
	data.AvatarUrl = r.AvatarUrl
	data.NickName = r.NickName
	data.CheckTime = CheckTime(r.CheckTime)
	return
}
