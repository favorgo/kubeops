package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	api "github.com/pipperman/kubeops/api/v1"
	"github.com/pipperman/kubeops/app/kubeops/biz"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewKubeOpsService)

var _ api.KubeOpsApiServer = &KubeOpsService{}

type KubeOpsService struct {
	api.UnimplementedKubeOpsApiServer

	kc  *biz.KubeOpsUseCase
	log *log.Helper
}

func NewKubeOpsService(kc *biz.KubeOpsUseCase, logger log.Logger) *KubeOpsService {
	return &KubeOpsService{
		kc:  kc,
		log: log.NewHelper(log.With(logger, "module", "service/kubeops")),
	}
}
