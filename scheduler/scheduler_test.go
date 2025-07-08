package scheduler

import (
	"sync"
	"testing"
)

func TestSchedule(t *testing.T) {
	Initialize()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
