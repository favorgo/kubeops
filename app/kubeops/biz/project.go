package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pipperman/kubeops/app/kubeops/biz/types"
)

type ProjectRepo interface {
	CreateProject(ctx context.Context, name, source string) (*types.Project, error)
	ListProject(ctx context.Context, pageNum, pageSize int64, param *types.ListProjectParam) ([]*types.Project, error)
	SearchPlaybooks(ctx context.Context, projectName string) ([]string, error)
	GetProject(ctx context.Context, name string) (*types.Project, error)
	SearchProjects(ctx context.Context) ([]*types.Project, error)
	IsProjectExists(ctx context.Context, name string) (bool, error)
}

type ProjectUseCase struct {
	repo ProjectRepo
	log  *log.Helper
}

func NewProjectUseCase(repo ProjectRepo, logger log.Logger) *ProjectUseCase {
	return &ProjectUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "usecase/project")),
	}
}

func (proj *ProjectUseCase) CreateProject(ctx context.Context, name, source string) (*types.Project, error) {
	return nil, nil
}

func (proj *ProjectUseCase) ListProject(ctx context.Context, pageNum, pageSize int64, param *types.ListProjectParam) ([]*types.Project, error) {
	return nil, nil
}
