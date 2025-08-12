package util

import (
	"app/code"
	"context"
	"fmt"
	"sync"
	"time"
)

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
	return context.WithValue(context.Background(), code.TraceInfoKey, ti)
}

func NewRootContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(NewRootContext(), timeout)
}

func NewRootContextWithTraceId(traceId string) context.Context {
	ti := &TraceInfo{
		TraceId: traceId,
	}
	return context.WithValue(context.Background(), code.TraceInfoKey, ti)
}

// NewChildContext 从一个父上下文中创建一个子上下文
func NewChildContext(parentCtx context.Context) context.Context {
	if parentCtx == nil {
		return NewRootContext()
	}
	parentInfo, ok := parentCtx.Value(code.TraceInfoKey).(*TraceInfo)
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

	return context.WithValue(context.Background(), code.TraceInfoKey, childInfo)
}
