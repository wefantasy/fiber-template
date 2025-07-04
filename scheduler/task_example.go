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
		_, err := c.AddFunc("*/2 * * * * *", func() {
			o.Run()
		})
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func (o *ExampleTask) Run() {
	startTime := time.Now()
	log.Infof("Start Scheduling %s", o.Name())

	log.Info("Start Scheduling Task")
	time.Sleep(time.Second * 5)

	log.Infof("End Scheduling %s, Used Time: %s", o.Name(), time.Since(startTime))
}
