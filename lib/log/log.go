package log

import (
	"fmt"
	"os"
	"time"
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
	if zapLevel <= zapcore.FatalLevel {
		fatalHook := getLumberjackConfig("fatal")
		fataLevel := getLevel(zapcore.FatalLevel)
		allCore = append(allCore, zapcore.NewCore(
			encoder,
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&fatalHook)),
			fataLevel,
		))
	}
	if zapLevel <= zapcore.PanicLevel {
		panicHook := getLumberjackConfig("panic")
		panicLevel := getLevel(zapcore.PanicLevel)
		allCore = append(allCore, zapcore.NewCore(
			encoder,
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&panicHook)),
			panicLevel,
		))
	}
	if zapLevel <= zapcore.DPanicLevel {
		dpanicHook := getLumberjackConfig("dpanic")
		dpanicLevel := getLevel(zapcore.DPanicLevel)
		allCore = append(allCore, zapcore.NewCore(
			encoder,
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&dpanicHook)),
			dpanicLevel,
		))
	}
	if zapLevel <= zapcore.ErrorLevel {
		errorHook := getLumberjackConfig("error")
		errorLevel := getLevel(zapcore.ErrorLevel)
		allCore = append(allCore, zapcore.NewCore(
			encoder,
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&errorHook)),
			errorLevel,
		))
	}
	if zapLevel <= zapcore.WarnLevel {
		warnHook := getLumberjackConfig("warn")
		warnLevel := getLevel(zapcore.WarnLevel)
		allCore = append(allCore, zapcore.NewCore(
			encoder,
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&warnHook)),
			warnLevel,
		))
	}
	if zapLevel <= zapcore.InfoLevel {
		infoHook := getLumberjackConfig("info")
		infoLevel := getLevel(zapcore.InfoLevel)
		allCore = append(allCore, zapcore.NewCore(
			encoder,
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&infoHook)),
			infoLevel,
		))
	}
	if zapLevel <= zapcore.DebugLevel {
		debugHook := getLumberjackConfig("debug")
		debugLevel := getLevel(zapcore.DebugLevel)
		allCore = append(allCore, zapcore.NewCore(
			encoder,
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&debugHook)),
			debugLevel,
		))
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

func getLumberjackConfig(name string) lumberjack.Logger {
	path := zapConfig["path"]
	return lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s/%s.log", path, time.Now().Format("2006-01-02"), name), // 日志文件路径
		MaxSize:    128,                                                                      // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,                                                                       // 日志文件最多保存多少个备份
		MaxAge:     7,                                                                        // 文件最多保存多少天
		Compress:   true,                                                                     // 是否压缩
	}
}
