package scheduler

import (
	"app/log"
	"github.com/robfig/cron/v3"
	"time"
)

func Initialize() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("CRON SCHEDULER PANICKED!")
				log.Error(r)
			}
		}()

		loc, _ := time.LoadLocation("Asia/Shanghai")
		l := &loggerAdapter{}
		c := cron.New(cron.WithLogger(l), cron.WithChain(cron.Recover(l)), cron.WithSeconds(), cron.WithLocation(loc))

		tasks := []Task{
			NewExampleTask(),
		}
		for _, task := range tasks {
			task.Register(c)
		}

		c.Start()
	}()
}

type loggerAdapter struct{}

func (l *loggerAdapter) Info(msg string, keysAndValues ...any) {
	log.Infow(msg, keysAndValues...)
}

func (l *loggerAdapter) Error(err error, msg string, keysAndValues ...any) {
	allFields := append([]any{"error", err}, keysAndValues...)
	log.Errorw(msg, allFields...)
}
