package log

import (
	"fmt"
	"os"
	"xiaoyin/lib/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Log       *zap.Logger
	zapConfig map[string]interface{}
	levelMap  = map[string]zapcore.Level{
		"debug":  zapcore.DebugLevel,
		"info":   zapcore.InfoLevel,
		"warn":   zapcore.WarnLevel,
		"error":  zapcore.ErrorLevel,
		"dpanic": zapcore.DPanicLevel,
		"panic":  zapcore.PanicLevel,
		"fatal":  zapcore.FatalLevel,
	}
)

func init() {
	zapConfig = config.Config.GetStringMap("log")
	Log = zap.New(getNewTee(), zap.AddStacktrace(zapcore.ErrorLevel))
}
func getNewTee() zapcore.Core {
	allCore := getNewCore()
	return zapcore.NewTee(allCore...)
}

func getNewCore() (allCore []zapcore.Core) {
	encoder := getEncoder()
	level := fmt.Sprintf("%s", zapConfig["level"])
	//如果没有匹配到，zapLevel=0，默认为info级别
	zapLevel := levelMap[level]
	//遍历所有level，增加多个日志等级文件
	for k, v := range levelMap {
		if zapLevel <= v {
			//获取日志分割
			hook := getLumberjackConfig(k)
			//每隔24小时主动分割
			//go func() {
			//	for {
			//		<-time.After(time.Hour * 24)
			//		_ = hook.Rotate()
			//	}
			//}()
			allCore = append(allCore, zapcore.NewCore(
				encoder,
				zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), //打印控制台和日志
				getLevel(v),
			))
		}
	}
	return
}

func getLevel(level zapcore.Level) zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		return lvl == level
	}
}

func getEncoder() zapcore.Encoder {
	switch zapConfig["format"] {
	case "console":
		return zapcore.NewConsoleEncoder(getEncoderConfig())
	case "json":
		return zapcore.NewJSONEncoder(getEncoderConfig())
	default:
		return zapcore.NewConsoleEncoder(getEncoderConfig())
	}
}

func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "file",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 将级别转换成大写
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //格式化Duration
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器,zapcore.ShortCallerEncoder
		EncodeName:     zapcore.FullNameEncoder,
	}
	switch zapConfig["encodelevel"] {
	case "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	return
}

//日志文件切割
func getLumberjackConfig(name string) lumberjack.Logger {
	path := zapConfig["path"]
	return lumberjack.Logger{
		//Filename:   fmt.Sprintf("%s/%s/%s.log", path, time.Now().Format("2006-01-02"), name), // 日志文件路径
		Filename:   fmt.Sprintf("%s/%s.log", path, name), // 日志文件路径
		MaxSize:    128,                                  // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,                                   // 被分割日志最大留存个数
		MaxAge:     15,                                   // 被分割日志最大留存天数
		Compress:   false,                                // 是否压缩
		LocalTime:  true,                                 //被分割的日志是否使用本地时间
	}
}
