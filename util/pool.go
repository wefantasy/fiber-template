package util

import (
	"sync"
)

type Pool struct {
	taskQueue chan func()
	wg        sync.WaitGroup
}

// NewPool 创建一个新的任务池
func NewPool(size int) *Pool {
	pool := &Pool{
		taskQueue: make(chan func(), size*2),
	}

	for i := 0; i < size; i++ {
		pool.wg.Add(1)
		go pool.worker()
	}

	return pool
}

// worker 处理任务队列中的任务
func (p *Pool) worker() {
	defer p.wg.Done()
	for task := range p.taskQueue {
		task()
	}
}

// Submit 提交任务到任务队列
func (p *Pool) Submit(task func()) {
	p.taskQueue <- task
}

// Stop 停止任务池
func (p *Pool) Stop() {
	close(p.taskQueue)
	p.wg.Wait()
}
