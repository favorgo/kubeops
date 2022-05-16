package biz

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pipperman/kubeops/app/kubeops/biz/types"
)

type InventoryRepo interface {
	GetInventory(ctx context.Context, taskID string) (*types.Inventory, error)
}

type Inventory struct {
	cache CacheRepo
	log   *log.Helper
}

func NewInventoryRepo(cache CacheRepo, logger log.Logger) InventoryRepo {
	return &Inventory{
		cache: cache,
		log:   log.NewHelper(log.With(logger, "module", "usecase/inventory")),
	}
}

func (ivn *Inventory) GetInventory(ctx context.Context, taskID string) (*types.Inventory, error) {
	item, _ := ivn.cache.Cache(ctx).Get(types.TaskID(taskID))
	if item == nil {
		return nil, errors.New("inventory is expire")
	}
	resp, ok := item.(*types.Inventory)
	if !ok {
		return nil, errors.New("internal error")
	}

	return resp, nil
}
