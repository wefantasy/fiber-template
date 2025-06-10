package scheduler

import (
	"app/log"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"time"
)

func Initialize() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				zap.L().Error("CRON SCHEDULER PANICKED!",
					zap.Any("panic_reason", r),
					zap.Stack("stacktrace"),
				)
			}
		}()

		loc, _ := time.LoadLocation("Asia/Shanghai")
		l := &loggerAdapter{}
		c := cron.New(cron.WithLogger(l), cron.WithChain(cron.Recover(l)), cron.WithSeconds(), cron.WithLocation(loc))
		TestTask()
		_, err := c.AddFunc("0 * * * * *", TestTask)
		if err != nil {
			log.Error(err)
			return
		}
		c.Start()
	}()
}

func TestTask() {
	log.Info(time.Now().Format("2006-01-02 15:04:05"))
}

type loggerAdapter struct{}

func (l *loggerAdapter) Info(msg string, keysAndValues ...interface{}) {
	log.Infow(msg, keysAndValues...)
}

func (l *loggerAdapter) Error(err error, msg string, keysAndValues ...interface{}) {
	allFields := append([]interface{}{"error", err}, keysAndValues...)
	log.Errorw(msg, allFields...)
}
