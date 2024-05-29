package protocol

import (
	"context"
	"fmt"
	"time"
)

// withTimeout 函数接受一个函数作为参数，执行该函数，并设置超时时间
func WithTimeout(fn func() error, timeout time.Duration) error {
	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 使用匿名函数执行传入的函数，以便在超时时取消执行
	done := make(chan error, 1)
	go func() {
		done <- fn()
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return fmt.Errorf("超时: %v", ctx.Err())
	}
}
