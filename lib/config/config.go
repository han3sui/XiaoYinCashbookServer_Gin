package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var Config *viper.Viper

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("load config faield: %s \n", err))
	}
	Config = viper.GetViper()
}
