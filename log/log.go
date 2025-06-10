package log

import (
	"app/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

func Initialize() {
	encoder := zapcore.NewJSONEncoder(getEncoder())
	//encoder := zapcore.NewConsoleEncoder(getEncoder())

	logHook1 := getWriter()
	logHook2 := os.Stdout

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(logHook1), conf.Logger.Level),
		zapcore.NewCore(encoder, zapcore.AddSync(logHook2), conf.Logger.Level),
	)
	// 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数
	logger := zap.New(core, zap.AddCallerSkip(1), zap.AddCaller(), zap.AddStacktrace(conf.Logger.StackTraceLevel))
	zap.ReplaceGlobals(logger)
	defer func() {
		if logHook1, ok := logHook1.(*lumberjack.Logger); ok {
			if err := logHook1.Close(); err != nil {
				panic(err)
			}
		}
	}()
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
		MaxSize:    conf.Logger.MaxSize,
		MaxBackups: conf.Logger.MaxBackups,
		MaxAge:     conf.Logger.MaxAge,
		Compress:   conf.Logger.EnableCompress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func Debug(args ...interface{}) {
	zap.S().Debug(args...)
}

func Info(args ...interface{}) {
	zap.S().Info(args...)
}

func Warn(args ...interface{}) {
	zap.S().Warn(args...)
}

func Error(args ...interface{}) {
	zap.S().Error(args...)
}

func DPanic(args ...interface{}) {
	zap.S().DPanic(args...)
}

func Panic(args ...interface{}) {
	zap.S().Panic(args...)
}

func Fatal(args ...interface{}) {
	zap.S().Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	zap.S().Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	zap.S().Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	zap.S().Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	zap.S().Errorf(template, args...)
}

func DPanicf(template string, args ...interface{}) {
	zap.S().DPanicf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	zap.S().Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	zap.S().Fatalf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	zap.S().Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	zap.S().Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	zap.S().Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	zap.S().Errorw(msg, keysAndValues...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	zap.S().DPanicw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	zap.S().Panicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	zap.S().Fatalw(msg, keysAndValues...)
}

func Debugln(args ...interface{}) {
	zap.S().Debugln(args...)
}

func Infoln(args ...interface{}) {
	zap.S().Infoln(args...)
}

func Warnln(args ...interface{}) {
	zap.S().Warnln(args...)
}

func Errorln(args ...interface{}) {
	zap.S().Errorln(args...)
}

func DPanicln(args ...interface{}) {
	zap.S().DPanicln(args...)
}

func Panicln(args ...interface{}) {
	zap.S().Panicln(args...)
}

func Fatalln(args ...interface{}) {
	zap.S().Fatalln(args...)
}
