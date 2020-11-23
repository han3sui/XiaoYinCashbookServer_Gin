package boot

import (
	"fmt"
	"net/http"
	"time"
	"xiaoyin/lib/config"
	"xiaoyin/router"
)

func Run() {
	server()
}

func server() {
	serverConfig := config.Config.GetStringMap("server")
	s := &http.Server{
		Addr:         serverConfig["addr"].(string),
		Handler:      router.InitRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := s.ListenAndServe()
	if err != nil {
		panic(fmt.Errorf("server start faield: %s \n", err))
	}
}
