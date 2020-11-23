package db

import (
	"fmt"
	"time"
	"xiaoyin/lib/config"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm/logger"

	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	dbErr error
)

func init() {
	dbConfig := config.Config.GetStringMap("mysql")
	link := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", dbConfig["username"], dbConfig["password"], dbConfig["addr"], dbConfig["dbname"], dbConfig["config"])
	gormConfig := &gorm.Config{}
	if gin.Mode() == "debug" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}
	DB, dbErr = gorm.Open(mysql.Open(link), gormConfig)
	if dbErr != nil {
		panic(fmt.Errorf("mysql connect faield: %s \n", dbErr))
	}
	if DB == nil {
		panic(fmt.Errorf("db init is nil"))
	}
	sqlDB, _ := DB.DB()
	//设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(50)
	//设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(50)
	//设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Minute)
}
