package scheduler

import "github.com/robfig/cron/v3"

type Task interface {
	Name() string
	Register(*cron.Cron)
	Run()
}
