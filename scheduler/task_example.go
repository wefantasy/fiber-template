package scheduler

import (
	"app/log"
	"time"
)

type ExampleTask struct {
}

func (o *ExampleTask) Run() {
	log.Info(time.Now().Format("2006-01-02 15:04:05"))
}
