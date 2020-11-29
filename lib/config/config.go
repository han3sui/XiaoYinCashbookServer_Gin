package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

var Config *viper.Viper

func init() {
	fileInfo, err := os.Stat("config.toml")
	if err != nil || fileInfo.IsDir() {
		createConfigToml()
	}
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("加载配置文件失败: %s", err))
	}
	Config = viper.GetViper()
}

func createConfigToml() {
	configTmp := []byte(`[server]
#监听地址
Addr = ":8199"
#数据路类型
DbType = "mysql"
#运行模式
Debug = true

[mysql]
#Link = "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
Username = "root"
Password = "root"
Addr = "127.0.0.1:3306"
DbName = "accountbook"
Config = "charset=utf8mb4&parseTime=True&loc=Local"

[log]
#日志保存路径
Path = "log"
#日志等级，可使用 "debug", "info", "warn", "error", "dpanic", "panic", "fatal",
Level = "debug"
#日志格式，console: 控制台, json: json格式输出
Format = "console"
#前缀
Prefix = "小音记账"
#LowercaseLevelEncoder:小写, LowercaseColorLevelEncoder:小写带颜色,CapitalLevelEncoder: 大写, CapitalColorLevelEncoder: 大写带颜色,
EncodeLevel = "CapitalLevelEncoder"
#是否打印台显示
LogInConsole = true

[redis]
#redis地址
Addr = "127.0.0.1:6379"
#redis密码
Password = ""
#redis使用数据库
Db = 0

[wx]
#小程序appid
AppId = ""
#小程序secret
AppSecret = ""

[jwt]
#jwt密钥
Secret = "4f7Z3UrBaDQSjWEo"
#jwt过期时间，单位秒
Exp = 86400
`)
	err := ioutil.WriteFile("config.toml", configTmp, 0644)
	if err != nil {
		panic("创建配置文件失败，请重试")
	}
}
