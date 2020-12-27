package router

import (
	"net/http"
	"xiaoyin/app/api/account"
	"xiaoyin/app/api/category"
	"xiaoyin/app/api/check"
	"xiaoyin/app/api/detail"
	"xiaoyin/app/api/icon"
	"xiaoyin/app/api/token"
	"xiaoyin/app/api/user"
	"xiaoyin/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() http.Handler {
	r := gin.New()
	r.NoMethod(middleware.NotFound())
	r.NoRoute(middleware.NotFound())
	r.Use(middleware.LogInfo(), middleware.Recover(), middleware.Cors())
	v2 := r.Group("/api/v2")
	v2.POST("token", token.Grant)
	v2.Use(middleware.Jwt())
	{
		accountGroup := v2.Group("/accounts")
		{
			accountGroup.GET("", account.ListByUid)
			accountGroup.GET("/details/count/:id", account.GetDetailsCount)
			accountGroup.GET("/manage-list", account.ManageList)
			accountGroup.POST("", account.Save)
			accountGroup.PUT("/:id", account.Update)
			accountGroup.DELETE("/:id", account.DelWithDetails)
		}
		iconGroup := v2.Group("/icons")
		{
			iconGroup.GET("", icon.List)
		}
		categoryGroup := v2.Group("/category")
		{
			categoryGroup.GET("", category.ListByUid)
			categoryGroup.GET("/details/count/:id", category.GetDetailsCount)
			categoryGroup.DELETE("/:id", category.DelWithDetails)
			categoryGroup.POST("", category.Save)
			categoryGroup.PUT("/:id", category.Update)
		}
		detailGroup := v2.Group("/details")
		{
			detailGroup.GET("/money", detail.ListMoneyByParams)
			detailGroup.GET("", detail.ListByParams)
			detailGroup.GET("/bill/:year", detail.Bill)
			detailGroup.GET("/chart", detail.Chart)
			detailGroup.GET("/uncheck/:time", detail.IsExistUncheck)
			detailGroup.GET("/days", detail.GetAllDays)
			detailGroup.GET("/claim/:claim", detail.ListClaim)
			detailGroup.POST("", detail.Save)
			detailGroup.DELETE("/:id", detail.Del)
			detailGroup.PUT("/:id", detail.Update)
		}
		userGroup := v2.Group("/user")
		{
			userGroup.GET("", user.Info)
			userGroup.PUT("/checktime", user.UpdateCheckTime)
		}
		checkGroup := v2.Group("/check")
		{
			checkGroup.GET("", check.List)
		}
	}
	return r
}
