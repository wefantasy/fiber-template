package scheduler

import (
	"app/conf"
	"app/log"
	"sync"
	"testing"
)

func TestSchedule(t *testing.T) {
	conf.Initialize()
	log.Initialize()
	Initialize()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
