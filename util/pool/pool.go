package pool

import (
	"context"
	"sync"
	"time"
)

type Task[T any] func() *T

type Pool[T any] struct {
	poolSize  int
	taskQueue chan Task[T]
	Results   chan *T
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	mu        sync.Mutex
	isClosed  bool
}

// NewPool 创建一个新的任务池
func NewPool[T any](poolSize int) *Pool[T] {
	return NewPoolWithContext[T](context.Background(), poolSize)
}

// NewPoolWithContext 使用提供的上下文创建任务池
func NewPoolWithContext[T any](ctx context.Context, poolSize int) *Pool[T] {
	ctx, cancel := context.WithCancel(ctx)

	pool := &Pool[T]{
		poolSize:  poolSize,
		taskQueue: make(chan Task[T], poolSize*2),
		Results:   make(chan *T, poolSize*2),
		wg:        sync.WaitGroup{},
		ctx:       ctx,
		cancel:    cancel,
	}

	pool.run()
	return pool
}

func (p *Pool[T]) run() {
	for i := 0; i < p.poolSize; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case task, ok := <-p.taskQueue:
					if !ok {
						// 任务队列关闭，worker退出
						return
					}
					p.Results <- task()
				case <-p.ctx.Done():
					// 上下文取消，worker退出
					return
				}
			}
		}()
	}

	// 等所有 worker 完成后关闭 results channel
	go func() {
		p.wg.Wait()
		close(p.Results)
	}()
}

// Submit 提交任务到任务队列；返回 false 表示池已关闭，任务提交失败
func (p *Pool[T]) Submit(task Task[T]) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isClosed {
		return false
	}

	p.taskQueue <- task
	return true
}

// Close 关闭任务池，不再接受新任务；已在队列中的任务仍会被执行
func (p *Pool[T]) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.isClosed {
		return
	}

	p.isClosed = true
	close(p.taskQueue)
	// 如需立刻打断所有 worker，可启用下面这一行（会丢弃队列中未取出的任务）
	// p.cancel()
}

// Shutdown 会等待所有任务执行完毕并排空 Results
func (p *Pool[T]) Shutdown() {
	p.Close()
	// 排空可能未被接收的结果（若已全部取完，此处为快速返回）
	for range p.Results {
	}
}

// ExecuteBatch 批量执行任务并返回结果
func ExecuteBatch[T any](tasks []Task[T], poolSize int) []T {
	if len(tasks) == 0 {
		return nil
	}
	if len(tasks) < poolSize {
		poolSize = len(tasks)
	}

	pool := NewPool[T](poolSize)
	// 使用 Shutdown 而不是 Close，确保所有任务完成
	defer pool.Shutdown()

	go func() {
		for _, task := range tasks {
			pool.Submit(task)
		}
		pool.Close()
	}()

	results := make([]T, 0, len(tasks))
	for r := range pool.Results {
		// 与原实现一致：忽略 nil（对非可空类型永远为非 nil）
		if any(r) != nil && r != nil {
			results = append(results, *r)
		}
	}
	return results
}

// ExecuteBatchWithTimeout 批量执行任务，带有超时控制
func ExecuteBatchWithTimeout[T any](tasks []Task[T], poolSize int, timeout time.Duration) []T {
	if len(tasks) == 0 {
		return nil
	}
	if len(tasks) < poolSize {
		poolSize = len(tasks)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pool := NewPoolWithContext[T](ctx, poolSize)
	defer pool.Shutdown()

	// 在单独协程中提交任务，防止超时阻塞
	go func() {
		for _, task := range tasks {
			select {
			case <-ctx.Done():
				// 超时发生，关闭任务队列，让已提交的任务继续执行
				pool.Close()
				return
			default:
				pool.Submit(task)
			}
		}
		pool.Close()
	}()

	var results []T
	for r := range pool.Results {
		if any(r) != nil {
			results = append(results, *r)
		}
	}
	return results
}
