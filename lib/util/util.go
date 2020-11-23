package util

import (
	"strconv"
	"time"
	"xiaoyin/lib/exception"

	"github.com/shopspring/decimal"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

func HasError(c *gin.Context, code int, err error) {
	if err != nil {
		exception.Common(c, code, err)
	}
}

func GetUid(c *gin.Context) (uid uint, err error) {
	v, ok := c.Get("uid")
	if ok {
		v1, ok := v.(int64)
		if ok {
			uid = uint(v1)
			return
		} else {
			err = errors.New("用户信息错误")
			return
		}
	} else {
		err = errors.New("用户信息错误")
		return
	}
}

//GetMonthStartAndEnd 获取月份的第一天和最后一天
func GetMonthStartAndEnd(myYear string, myMonth string) (r map[string]string) {
	// 数字月份必须前置补零
	if len(myMonth) == 1 {
		myMonth = "0" + myMonth
	}
	yInt, _ := strconv.Atoi(myYear)

	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, myYear+"-"+myMonth+"-01 00:00:00", loc)
	newMonth := theTime.Month()

	t1 := time.Date(yInt, newMonth, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	t2 := time.Date(yInt, newMonth+1, 0, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	r = map[string]string{"start": t1, "end": t2}
	return
}

//加法
func FloatAdd(num1 float64, num2 float64, point int) (res float64) {
	decimal.DivisionPrecision = point
	res, _ = decimal.NewFromFloat(num1).Add(decimal.NewFromFloat(num2)).Float64()
	return
}

//减法
func FloatSub(num1 float64, num2 float64, point int) (res float64) {
	decimal.DivisionPrecision = point
	res, _ = decimal.NewFromFloat(num1).Sub(decimal.NewFromFloat(num2)).Float64()
	return
}

//除法
func FloatDiv(num1 float64, num2 float64, point int) (res float64) {
	decimal.DivisionPrecision = point
	res, _ = decimal.NewFromFloat(num1).Div(decimal.NewFromFloat(num2)).Float64()
	return
}

//乘法
func FloatMul(num1 float64, num2 float64, point int) (res float64) {
	decimal.DivisionPrecision = point
	res, _ = decimal.NewFromFloat(num1).Mul(decimal.NewFromFloat(num2)).Float64()
	return
}
