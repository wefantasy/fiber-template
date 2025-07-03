package scheduler

import (
	"app/conf"
	"app/log"
	"app/util/collect"
	"github.com/robfig/cron/v3"
	"time"
)

type ExampleTask struct {
}

func NewExampleTask() *ExampleTask {
	return &ExampleTask{}
}

func (o *ExampleTask) Name() string {
	return "ExampleTask"
}

func (o *ExampleTask) Register(c *cron.Cron) {
	isRunAtStartup := collect.Contains(conf.Scheduler.RunAtStartupTasks, func(task string) bool {
		return task == o.Name()
	})
	if isRunAtStartup {
		go o.Run()
	}

	isEnable := collect.Contains(conf.Scheduler.EnableTasks, func(task string) bool {
		return task == o.Name()
	})
	if isEnable {
		_, err := c.AddFunc("*/5 * * * * *", func() {
			o.Run()
		})
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func (o *ExampleTask) Run() {
	log.Info(time.Now().Format("2006-01-02 15:04:05"))
}
