package service

import (
	"encoding/json"
	"time"
	"xiaoyin/app/model"
	"xiaoyin/lib/config"

	"github.com/valyala/fasthttp"

	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"
)

var wxConfig = config.Config.GetStringMap("wx")
var jwtConfig = config.Config.GetStringMap("jwt")

type JwtInfo struct {
	ID  int64  `json:"id"`
	Exp uint64 `json:"exp"`
	Nbf uint64 `json:"nbf"`
	jwt.MapClaims
}

type LoginInfo struct {
	Code      string `json:"code" form:"code" validate:"required" label:"微信Code"`
	NickName  string `json:"nick_name" form:"nick_name"`
	AvatarUrl string `json:"avatar_url" form:"avatar_url"`
}

//type UserInfo = Info

func GrantToken(data *LoginInfo) (token string, err error) {
	r, err := wxAuth(data.Code)
	if err != nil {
		return
	}
	if _, ok := r["errcode"]; ok {
		err = errors.New(r["errmsg"].(string))
		return
	}
	id, err := SaveOrUpdate(&model.User{
		NickName:   data.NickName,
		AvatarUrl:  data.AvatarUrl,
		Openid:     r["openid"].(string),
		SessionKey: r["session_key"].(string),
	})
	if err != nil {
		return
	}
	exp := time.Now().Unix() + jwtConfig["exp"].(int64)
	claims := jwt.MapClaims{
		"exp": exp,
		"nbf": time.Now().Unix(),
		"id":  id,
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenClaims.SignedString([]byte(jwtConfig["secret"].(string)))
	if err != nil {
		return
	}
	//取消颁发token的时候存入redis，该步骤无实际意义，如果需要管理颁发的token，在吊销的时候，存入吊销token到redis，登录的时候，去吊销里面查询是否存在，如果存在，则禁止登录
	//rData, err := json.Marshal(map[string]interface{}{
	//	"id":          id,
	//	"nick_name":   data.NickName,
	//	"avatar_url":  data.AvatarUrl,
	//	"openid":      r["openid"],
	//	"session_key": r["session_key"],
	//})
	//if err != nil {
	//	return
	//}
	//err = redis.Redis.Set(context.Background(), token, rData, time.Duration(expDuration)*time.Second).Err()
	//if err != nil {
	//	err = errors.Wrap(err, "Redis写入错误")
	//	return
	//}
	return
}

func ParseToken(token string) (tokenInfo *JwtInfo, err error) {
	//取消redis查询验证，直接解析token
	//err = redis.Redis.Get(context.Background(), token).Err()
	//if err == redis.Nil {
	//	err = errors.Wrap(err, "token不存在")
	//	return
	//} else if err != nil {
	//	panic(err)
	//}
	jwtToken, err := jwt.ParseWithClaims(token, &JwtInfo{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtConfig["secret"].(string)), nil
	})
	if err != nil || !jwtToken.Valid {
		err = errors.Wrap(err, "token解析错误")
		return
	}
	claim, ok := jwtToken.Claims.(*JwtInfo)
	if !ok {
		err = errors.Wrap(err, "token解析错误")
		return
	}
	tokenInfo = claim
	return
}

func wxAuth(code string) (strMap map[string]interface{}, err error) {
	appid := wxConfig["appid"].(string)
	secret := wxConfig["appsecret"].(string)
	authUrl := "https://api.weixin.qq.com/sns/jscode2session?appid=" + appid + "&secret=" + secret + "&js_code=" + code + "&grant_type=authorization_code"
	client := &fasthttp.Client{ReadTimeout: 10 * time.Second}
	status, r, err := client.Get(nil, authUrl)
	if err != nil {
		err = errors.Wrap(err, "请求微信服务器失败")
		return
	}
	if status != fasthttp.StatusOK {
		err = errors.New("微信服务器返回状态码错误")
		return
	}
	str := string(r)
	err = json.Unmarshal([]byte(str), &strMap)
	if err != nil {
		err = errors.Wrap(err, "微信返回Map解析失败")
		return
	}
	return
}
