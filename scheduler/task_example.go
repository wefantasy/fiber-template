package scheduler

import (
	"app/conf"
	"app/log"
	"app/util"
	"app/util/collect"
	"context"
	"github.com/robfig/cron/v3"
	"sync"
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
	rootCtx := util.NewRootContext()
	startTime := time.Now()
	log.T(rootCtx).Infof("Start Scheduling %s At %s", o.Name(), startTime.Format("2006-01-02 15:04:05"))

	tasks := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	pool := util.NewPool(conf.Goroutines)
	var wg sync.WaitGroup
	wg.Add(len(tasks))

	for _, task := range tasks {
		pool.Submit(rootCtx, func(ctx context.Context) {
			defer wg.Done()
			log.T(ctx).Infof("Process Task %d", task)
			time.Sleep(time.Second * 5)
		})
	}

	wg.Wait()
	pool.Stop()

	log.T(rootCtx).Infof("End Scheduling %s At %s, Used Time: %s", o.Name(), time.Now().Format("2006-01-02 15:04:05"), time.Since(startTime))
}
