package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pipperman/kubeops/app/kubeops/biz/types"
)

type ResultRepo interface {
	GetResult(ctx context.Context, taskID string) (*types.Result, error)
	ListResult(ctx context.Context, pageNum, pageSize int64, param *types.ListResultParam) (types.ResultItems, error)
	WatchResult(ctx context.Context, taskID string) (chan []byte, error)
}

type ResultUseCase struct {
	repo ResultRepo
	log  *log.Helper
}

func NewResultUseCase(repo ResultRepo, logger log.Logger) *ResultUseCase {
	return &ResultUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "usecase/result")),
	}
}

func (r *ResultUseCase) GetResult(ctx context.Context, taskID string) (*types.Result, error) {
	return nil, nil
}

func (r *ResultUseCase) ListResult(ctx context.Context, pageNum, pageSize int64, param *types.ListResultParam) (*types.ResultItems, error) {
	return nil, nil
}

func (r *ResultUseCase) WatchResult(ctx context.Context, taskID string) (chan []byte, error) {
	return nil, nil
}
