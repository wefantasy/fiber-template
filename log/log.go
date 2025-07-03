package log

import (
	"app/conf"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

func Initialize() {
	//encoder := zapcore.NewJSONEncoder(getEncoder())
	encoder := zapcore.NewConsoleEncoder(getEncoder())

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

	Infof("Use config: %v", conf.Conf)
}

// 生成日志编码配置
func getEncoder() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	cst, err := time.LoadLocation(conf.Timezone)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading CST location: %v\n", err)
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.In(cst).Format("2006-01-02 15:04:05.0000"))
		}
	}
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

func Debug(args ...any) {
	zap.S().Debug(args...)
}

func Info(args ...any) {
	zap.S().Info(args...)
}

func Warn(args ...any) {
	zap.S().Warn(args...)
}

func Error(args ...any) {
	zap.S().Error(args...)
}

func DPanic(args ...any) {
	zap.S().DPanic(args...)
}

func Panic(args ...any) {
	zap.S().Panic(args...)
}

func Fatal(args ...any) {
	zap.S().Fatal(args...)
}

func Debugf(template string, args ...any) {
	zap.S().Debugf(template, args...)
}

func Infof(template string, args ...any) {
	zap.S().Infof(template, args...)
}

func Warnf(template string, args ...any) {
	zap.S().Warnf(template, args...)
}

func Errorf(template string, args ...any) {
	zap.S().Errorf(template, args...)
}

func DPanicf(template string, args ...any) {
	zap.S().DPanicf(template, args...)
}

func Panicf(template string, args ...any) {
	zap.S().Panicf(template, args...)
}

func Fatalf(template string, args ...any) {
	zap.S().Fatalf(template, args...)
}

func Debugw(msg string, keysAndValues ...any) {
	zap.S().Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...any) {
	zap.S().Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...any) {
	zap.S().Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...any) {
	zap.S().Errorw(msg, keysAndValues...)
}

func DPanicw(msg string, keysAndValues ...any) {
	zap.S().DPanicw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...any) {
	zap.S().Panicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...any) {
	zap.S().Fatalw(msg, keysAndValues...)
}

func Debugln(args ...any) {
	zap.S().Debugln(args...)
}

func Infoln(args ...any) {
	zap.S().Infoln(args...)
}

func Warnln(args ...any) {
	zap.S().Warnln(args...)
}

func Errorln(args ...any) {
	zap.S().Errorln(args...)
}

func DPanicln(args ...any) {
	zap.S().DPanicln(args...)
}

func Panicln(args ...any) {
	zap.S().Panicln(args...)
}

func Fatalln(args ...any) {
	zap.S().Fatalln(args...)
}
