package router

import (
	"net/http"
	"xiaoyin/app/api"
	"xiaoyin/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() http.Handler {
	r := gin.New()
	r.NoMethod(middleware.NotFound())
	r.NoRoute(middleware.NotFound())
	r.Use(middleware.LogInfo(), middleware.Recover(), middleware.Cors())
	v2 := r.Group("/api/v2")
	v2.POST("token", api.GrantToken)
	v2.Use(middleware.Jwt())
	{
		accountGroup := v2.Group("/accounts")
		{
			accountGroup.GET("", api.ListAccountByUid)
			accountGroup.GET("/details/count/:id", api.GetDetailsCountByAid)
			accountGroup.GET("/manage-list", api.AccountManageList)
			accountGroup.POST("", api.SaveAccount)
			accountGroup.PUT("/:id", api.UpdateAccount)
			accountGroup.DELETE("/:id", api.DelAccountWithDetails)
		}
		iconGroup := v2.Group("/icons")
		{
			iconGroup.GET("", api.ListIcons)
		}
		categoryGroup := v2.Group("/category")
		{
			categoryGroup.GET("", api.ListCategoryByUid)
			categoryGroup.GET("/details/count/:id", api.GetDetailsCountByCid)
			categoryGroup.DELETE("/:id", api.DelCategoryWithDetails)
			categoryGroup.POST("", api.SaveCategory)
			categoryGroup.PUT("/:id", api.UpdateCategory)
		}
		detailGroup := v2.Group("/details")
		{
			detailGroup.GET("/money", api.ListMoneyByParams)
			detailGroup.GET("", api.ListDetailsByParams)
			detailGroup.GET("/bill/:year", api.Bill)
			detailGroup.GET("/chart", api.Chart)
			detailGroup.GET("/uncheck/:time", api.IsExistUncheck)
			detailGroup.GET("/days", api.GetAllDays)
			detailGroup.GET("/claim/:claim", api.ListClaim)
			detailGroup.POST("", api.SaveDetail)
			detailGroup.DELETE("/:id", api.DelDetail)
			detailGroup.PUT("/:id", api.UpdateDetail)
		}
		userGroup := v2.Group("/user")
		{
			userGroup.GET("", api.GetUser)
			userGroup.PUT("/checktime", api.UpdateCheckTime)
		}
		checkGroup := v2.Group("/check")
		{
			checkGroup.GET("", api.ListCheck)
		}
	}
	return r
}
