package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pipperman/kubeops/app/kubeops/biz/types"
)

type KubeOpsRepo interface {
	ProjectRepo
	RunnerRepo
	InventoryRepo
	ResultRepo
}

type KubeOpsUseCase struct {
	projRepo   ProjectRepo
	runnerRepo RunnerRepo
	invRepo    InventoryRepo
	resultRepo ResultRepo
	log        *log.Helper
}

func NewKubeOpsUseCase(projRepo ProjectRepo, runnerRepo RunnerRepo, invRepo InventoryRepo, resultRepo ResultRepo, logger log.Logger) *KubeOpsUseCase {
	return &KubeOpsUseCase{
		projRepo:   projRepo,
		runnerRepo: runnerRepo,
		invRepo:    invRepo,
		resultRepo: resultRepo,
		log:        log.NewHelper(log.With(logger, "module", "usecase/kubeops")),
	}
}

func (k *KubeOpsUseCase) CreateProject(ctx context.Context, name, source string) (*types.Project, error) {
	return k.projRepo.CreateProject(ctx, name, source)
}

func (k *KubeOpsUseCase) ListProject(ctx context.Context, pageNum, pageSize int64, param *types.ListProjectParam) ([]*types.Project, error) {
	return k.projRepo.ListProject(ctx, pageNum, pageSize, param)
}

func (k *KubeOpsUseCase) RunAdhoc(ctx context.Context, ivn *types.Inventory, adhoc *types.Adhoc) (*types.RunAdhocResult, error) {
	return k.runnerRepo.RunAdhoc(ctx, ivn, adhoc)
}

func (k *KubeOpsUseCase) RunPlaybook(ctx context.Context, playbook *types.Playbook) (*types.RunPlaybookResult, error) {
	return k.runnerRepo.RunPlaybook(ctx, playbook)
}

func (k *KubeOpsUseCase) GetInventory(ctx context.Context, taskID string) (*types.Inventory, error) {
	return k.invRepo.GetInventory(ctx, taskID)
}

func (k *KubeOpsUseCase) GetResult(ctx context.Context, taskID string) (*types.Result, error) {
	return k.resultRepo.GetResult(ctx, taskID)
}

func (k *KubeOpsUseCase) ListResult(ctx context.Context, pageNum, pageSize int64, param *types.ListResultParam) (types.ResultItems, error) {
	return k.resultRepo.ListResult(ctx, pageNum, pageSize, param)
}

func (k *KubeOpsUseCase) WatchResult(ctx context.Context, taskID string) (chan []byte, error) {
	return k.resultRepo.WatchResult(ctx, taskID)
}
