package util

import (
	"context"
	"fmt"
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
func (p *Pool) Submit(ctx context.Context, task func(ctx context.Context)) {
	childCtx := NewChildContext(ctx)
	wrappedTask := func() {
		task(childCtx)
	}
	p.taskQueue <- wrappedTask
}

// Stop 停止任务池
func (p *Pool) Stop() {
	close(p.taskQueue)
	p.wg.Wait()
}

const TraceIdKey = "TraceId"
const TraceHeaderIdKey = "X-Request-ID"
const TraceInfoKey = "TraceInfo"

type TraceInfo struct {
	TraceId    string
	childCount int
	mu         sync.Mutex
}

// NewRootContext 创建一个用于根任务的初始上下文
func NewRootContext() context.Context {
	traceId := RandTraceId()
	ti := &TraceInfo{
		TraceId: traceId,
	}
	return context.WithValue(context.Background(), TraceInfoKey, ti)
}

func NewRootContextWithTraceId(traceId string) context.Context {
	ti := &TraceInfo{
		TraceId: traceId,
	}
	return context.WithValue(context.Background(), TraceInfoKey, ti)
}

// NewChildContext 从一个父上下文中创建一个子上下文
func NewChildContext(parentCtx context.Context) context.Context {
	if parentCtx == nil {
		return NewRootContext()
	}
	parentInfo, ok := parentCtx.Value(TraceInfoKey).(*TraceInfo)
	if !ok {
		return NewRootContext()
	}

	parentInfo.mu.Lock()
	parentInfo.childCount++
	childSpanID := fmt.Sprintf("%s_%d", parentInfo.TraceId, parentInfo.childCount)
	parentInfo.mu.Unlock()

	childInfo := &TraceInfo{
		TraceId: childSpanID,
	}

	return context.WithValue(context.Background(), TraceInfoKey, childInfo)
}
