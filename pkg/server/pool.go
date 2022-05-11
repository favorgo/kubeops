package server

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type TaskFunc func()

type Pool struct {
	taskQueue  chan TaskFunc
	workerSize int
	log        *log.Helper
}

func NewPool(ctx context.Context, queueSize, workerSize int, logger log.Logger) *Pool {
	p := &Pool{taskQueue: make(chan TaskFunc, queueSize), workerSize: workerSize, log: log.NewHelper(logger)}
	for i := 0; i < p.workerSize; i++ {
		go p.run(ctx)
	}
	return p
}

func (p *Pool) Commit(ctx context.Context, taskFunc TaskFunc) {
	p.log.WithContext(ctx).Infof("current worker size: %d, current task queue size: %d", p.workerSize, p.Len())
	p.taskQueue <- taskFunc
}

func (p *Pool) Len() int {
	return len(p.taskQueue)
}

func (p *Pool) IsEmpty() bool {
	return len(p.taskQueue) == 0
}

func (p *Pool) run(ctx context.Context) {
	for {
		select {
		case task := <-p.taskQueue:
			task()
		case <-ctx.Done():
			return
		}
	}
}
