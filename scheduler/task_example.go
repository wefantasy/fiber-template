package scheduler

import (
	"app/conf"
	"app/log"
	"app/util"
	"app/util/collect"
	"app/util/pool"
	"time"

	"github.com/robfig/cron/v3"
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

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	tasks := make([]pool.Task[int], 0, len(data))
	for _, d := range data {
		ctx := util.NewChildContext(rootCtx)
		tasks = append(tasks, func() *int {
			log.T(ctx).Infof("Process Task %d", d)
			time.Sleep(time.Second * 1)
			return &d
		})
	}
	results := pool.ExecuteBatch(tasks, conf.Goroutines)
	log.T(rootCtx).Infof("Result is %v", results)
	log.T(rootCtx).Infof("End Scheduling %s At %s, Used Time: %s", o.Name(), time.Now().Format("2006-01-02 15:04:05"), time.Since(startTime))
}
