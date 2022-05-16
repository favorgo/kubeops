package data

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/pipperman/kubeops/app/kubeops/biz"
	"github.com/pipperman/kubeops/app/pkg/config"
)

type cacheRepo struct {
	conf *config.Cache
}

func (c cacheRepo) Cache(ctx context.Context) *cache.Cache {
	return cache.New(time.Hour*c.conf.Memcached.DefaultExpiration.AsDuration(), time.Minute*c.conf.Memcached.CleanupInterval.AsDuration())
}

func NewCacheRepo(conf *config.Cache) biz.CacheRepo {
	return &cacheRepo{
		conf: conf,
	}
}
