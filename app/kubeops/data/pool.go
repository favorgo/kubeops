package data

import (
	"context"
	"github.com/pipperman/kubeops/app/pkg/config"

	"github.com/go-kratos/kratos/v2/log"
)

type TaskFunc func()

type Pool struct {
	taskQueue  chan TaskFunc
	workerSize int
	log        *log.Helper
}

type PoolRepo interface {
	Commit(ctx context.Context, taskFunc TaskFunc)
	Len() int
	IsEmpty() bool
}

func NewPoolRepo(server *config.Server, logger log.Logger) PoolRepo {
	ctx := context.Background()
	p := &Pool{
		taskQueue:  make(chan TaskFunc, server.Task.QueueSize),
		workerSize: int(server.Task.WorkerSize),
		log:        log.NewHelper(log.With(logger, "module", "usecase/pool")),
	}
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
