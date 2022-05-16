package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pipperman/kubeops/app/kubeops/biz/types"
	"github.com/pipperman/kubeops/app/pkg/ansible"
)

type RunnerRepo interface {
	RunAdhoc(ctx context.Context, ivn *types.Inventory, adhoc *types.Adhoc) (*types.RunAdhocResult, error)
	RunPlaybook(ctx context.Context, playbook *types.Playbook) (*types.RunPlaybookResult, error)
}

type RunnerFactoryRepo interface {
	CreateAdhocRunner(ctx context.Context, pattern, module, param string) (*ansible.AdhocRunner, error)
	CreatePlaybookRunner(ctx context.Context, projectName, playbookName, tag string) (*ansible.PlaybookRunner, error)
	PreRunPlaybook(ctx context.Context, projectName, playbookName string) error
}

type RunnerUseCase struct {
	repo RunnerRepo
	log  *log.Helper
}

func NewRunnerUseCase(repo RunnerRepo, logger log.Logger) *RunnerUseCase {
	return &RunnerUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "usecasse/runner")),
	}
}

func (r *RunnerUseCase) RunAdhoc(ctx context.Context, ivn *types.Inventory, adhoc *types.Adhoc) (*types.RunAdhocResult, error) {
	return r.repo.RunAdhoc(ctx, ivn, adhoc)
}

func (r *RunnerUseCase) RunPlaybook(ctx context.Context, playbook *types.Playbook) (*types.RunPlaybookResult, error) {
	return r.repo.RunPlaybook(ctx, playbook)
}
