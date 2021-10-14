package server

import (
	"context"
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
)

type TaskFunc func()

type Pool struct {
	taskQueue  chan TaskFunc
	workerSize int
}

func NewPool(ctx context.Context) *Pool {
	workerSize := viper.GetInt("app.worker")
	queueSize := viper.GetInt("app.queue")
	p := &Pool{taskQueue: make(chan TaskFunc, queueSize), workerSize: workerSize}
	for i := 0; i < p.workerSize; i++ {
		go p.run(ctx)
	}
	return p
}

func (p *Pool) Commit(taskFunc TaskFunc) {
	log.Infof("receive a task")
	log.Infof("current worker size: %d", p.workerSize)
	log.Infof("task queue size: %d", p.Len())
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
