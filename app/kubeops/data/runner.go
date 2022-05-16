package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/patrickmn/go-cache"
	"github.com/pipperman/kubeops/app/kubeops/biz"
	"github.com/pipperman/kubeops/app/kubeops/biz/types"
	"github.com/pipperman/kubeops/app/pkg/ansible"
	uuid "github.com/satori/go.uuid"
	"time"
)

type runnerRepo struct {
	pool        PoolRepo
	cacheRepo   biz.CacheRepo
	factoryRepo biz.RunnerFactoryRepo
	log         *log.Helper
}

func NewRunnerRepo(poolRepo PoolRepo, repo biz.CacheRepo, factoryRepo biz.RunnerFactoryRepo, logger log.Logger) biz.RunnerRepo {
	return &runnerRepo{
		pool:        poolRepo,
		cacheRepo:   repo,
		factoryRepo: factoryRepo,
		log:         log.NewHelper(log.With(logger, "module", "repo/runner")),
	}
}

func (r *runnerRepo) RunAdhoc(ctx context.Context, ivn *types.Inventory, adhoc *types.Adhoc) (*types.RunAdhocResult, error) {
	// create adhoc runner
	adhocRunner, err := r.factoryRepo.CreateAdhocRunner(ctx, adhoc.Pattern, adhoc.Module, adhoc.Param)
	if err != nil {
		return nil, err
	}
	// create task
	dataCh := make(chan []byte)
	id := uuid.NewV4().String()
	result := r.initialResult(id)
	r.setCache(ctx, result, dataCh, ivn)
	task := func() {
		adhocRunner.Run(dataCh, result)
		result.Finished = true
		result.EndTime = time.Now().Format("2006-01-02 15:04:05")
		// update task status
		r.cacheRepo.Cache(ctx).Set(types.TaskID(result.Id), &result, cache.DefaultExpiration)
	}
	// commit task
	r.log.WithContext(ctx).Infof("commit a adhoc task: %s", result.Id)
	r.pool.Commit(ctx, task)
	return &types.RunAdhocResult{
		Result: result,
	}, nil
}

func (r *runnerRepo) RunPlaybook(ctx context.Context, playbook *types.Playbook) (*types.RunPlaybookResult, error) {
	panic("implement me")
}

func (r *runnerRepo) initialResult(id string, param ...string) *types.Result {
	result := &types.Result{
		Id:        id,
		StartTime: time.Now().Format("2006-01-02 15:04:05"),
		EndTime:   "",
		Message:   "",
		Success:   false,
		Finished:  false,
		Content:   "",
	}

	if param != nil {
		result.Project = param[0]
	}
	return result
}

func (r *runnerRepo) setCache(ctx context.Context, result *types.Result, dataCh chan []byte, inventory *types.Inventory) {
	r.cacheRepo.Cache(ctx).Set(types.TaskID(result.Id), result, cache.DefaultExpiration)
	r.cacheRepo.Cache(ctx).Set(types.DataID(result.Id), dataCh, cache.DefaultExpiration)
	r.cacheRepo.Cache(ctx).Set(types.InventoryID(result.Id), inventory, cache.DefaultExpiration)
}

type runnerFactoryRepo struct {
	proj biz.ProjectRepo
}

func NewRunnerFactoryRepo(proj biz.ProjectRepo) biz.RunnerFactoryRepo {
	return &runnerFactoryRepo{
		proj: proj,
	}
}

func (r *runnerFactoryRepo) CreatePlaybookRunner(ctx context.Context, projectName, playbookName, tag string) (*ansible.PlaybookRunner, error) {
	err := r.PreRunPlaybook(ctx, projectName, playbookName)
	if err != nil {
		return nil, err
	}
	p, err := r.proj.GetProject(ctx, projectName)
	if err != nil {
		return nil, err
	}
	return &ansible.PlaybookRunner{
		Project:  p,
		Playbook: playbookName,
		Tag:      tag,
	}, nil
}

func (r *runnerFactoryRepo) CreateAdhocRunner(ctx context.Context, pattern, module, param string) (*ansible.AdhocRunner, error) {
	return &ansible.AdhocRunner{
		Module:  module,
		Param:   param,
		Pattern: pattern,
	}, nil
}

func (r *runnerFactoryRepo) PreRunPlaybook(ctx context.Context, projectName, playbookName string) error {
	p, err := r.proj.GetProject(ctx, projectName)
	if err != nil {
		return err
	}
	exists := false
	for _, playbook := range p.Playbooks {
		if playbook == playbookName {
			exists = true
		}
	}
	if !exists {
		return errors.New(fmt.Sprintf("can not find playbook:%s in project:%s", playbookName, projectName))
	}
	return nil
}
