package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/patrickmn/go-cache"
)

type CacheRepo interface {
	Cache(ctx context.Context) *cache.Cache
}

type CacheUseCase struct {
	cacheRepo CacheRepo
	log       *log.Helper
}

func NewCacheUseCase(cacheRepo CacheRepo, logger log.Logger) *CacheUseCase {
	return &CacheUseCase{
		cacheRepo: cacheRepo,
		log:       log.NewHelper(log.With(logger, "module", "usecase/cache")),
	}
}

func (c *CacheUseCase) Cache(ctx context.Context) *cache.Cache {
	return c.cacheRepo.Cache(ctx)
}
