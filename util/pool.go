package util

import "sync"

type Pool struct {
	taskQueue chan func()
	wg        sync.WaitGroup
}

func NewPool(maxWorkers int) *Pool {
	pool := &Pool{
		taskQueue: make(chan func(), 1000),
	}

	// 启动工作协程
	for i := 0; i < maxWorkers; i++ {
		pool.wg.Add(1)
		go pool.worker()
	}

	return pool
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for task := range p.taskQueue {
		task()
	}
}

func (p *Pool) Submit(task func()) {
	p.taskQueue <- task
}

func (p *Pool) Stop() {
	close(p.taskQueue)
	p.wg.Wait()
}
