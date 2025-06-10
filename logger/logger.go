package logger

import (
	"app/conf"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

var Handler *zap.Logger

func Initialize() {
	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	encoder := zapcore.NewJSONEncoder(getEncoder())
	//encoder := zapcore.NewConsoleEncoder(getEncoder())

	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
	logHook1 := getWriter()
	logHook2 := os.Stdout

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(logHook1), conf.Logger.Level),
		zapcore.NewCore(encoder, zapcore.AddSync(logHook2), conf.Logger.Level),
	)
	// 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(conf.Logger.StackTraceLevel))
	zap.ReplaceGlobals(logger)
	fiberzap.New(fiberzap.Config{})
	l := fiberzap.NewLogger(fiberzap.LoggerConfig{
		SetLogger: logger,
		ExtraKeys: []string{"request_id"},
	})
	defer func(l *fiberzap.LoggerConfig) {
		err := l.Sync()
		if err != nil {
			if err != nil {
				panic(err)
			}
		}
	}(l)
	log.SetLogger(l)
	Handler = logger
}

// 生成日志编码配置
func getEncoder() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	// 设置时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 返回json 格式的 日志编辑器
	return encoderConfig
}

// 获取日志切割句柄
func getWriter() io.Writer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   conf.Logger.Filename,
		MaxSize:    conf.Logger.MaxSize,        // 单个文件最大100M
		MaxBackups: conf.Logger.MaxBackups,     // 多于 60 个日志文件后，清理较旧的日志
		MaxAge:     conf.Logger.MaxAge,         // 一天一切割
		Compress:   conf.Logger.EnableCompress, //是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}
