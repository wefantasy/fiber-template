package pool

import (
	"app/util/copier"
	"strconv"
	"testing"
	"time"
)

func TestPoolExecuteBatch(t *testing.T) {
	tasks := make([]Task, 0, 5)
	for i := 0; i < 10; i++ {
		tasks = append(tasks, func() any {
			time.Sleep(1000 * time.Millisecond)
			t.Logf("协程%d：执行任务", i)
			return strconv.Itoa(i)
		})
	}
	results := ExecuteBatch(tasks, 5)
	var target []string
	copier.TransferListType(results, &target)
	t.Logf("执行结果：%v", target)
}

func TestPoolExecuteBatchWithTimeout(t *testing.T) {
	tasks := make([]Task, 0, 5)
	for i := 0; i < 10; i++ {
		tasks = append(tasks, func() any {
			time.Sleep(1000 * time.Millisecond)
			t.Logf("协程%d：执行任务", i)
			return i
		})
	}
	results := ExecuteBatchWithTimeout(tasks, 1, 10000*time.Millisecond)
	t.Logf("执行结果：%v", results)
}

func TestPoolSubmit(t *testing.T) {
	pool := NewPool(5)
	go func() {
		time.Sleep(1000 * time.Millisecond)
		pool.Submit(func() any {
			t.Log("协程A：执行任务")
			return nil
		})
	}()
	time.Sleep(50 * time.Millisecond)
	pool.Close() // 主协程：关闭 pool
	time.Sleep(5000 * time.Millisecond)
}
