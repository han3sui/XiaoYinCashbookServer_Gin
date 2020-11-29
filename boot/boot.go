package boot

import (
	"fmt"
	"net/http"
	"time"
	"xiaoyin/lib/config"
	"xiaoyin/lib/db"
	"xiaoyin/lib/redis"
	"xiaoyin/router"

	"github.com/gin-gonic/gin"
)

func Run() {
	serverConfig := config.Config.GetStringMap("server")
	if !serverConfig["debug"].(bool) {
		gin.SetMode(gin.ReleaseMode)
	}
	db.Init()
	redis.Init()
	s := &http.Server{
		Addr:         serverConfig["addr"].(string),
		Handler:      router.InitRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := s.ListenAndServe()
	if err != nil {
		panic(fmt.Errorf("服务启动失败: %s", err))
	}
}
