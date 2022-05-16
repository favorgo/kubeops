package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	api "github.com/pipperman/kubeops/api/v1"
	"github.com/pipperman/kubeops/app/kubeops/service"
	"github.com/pipperman/kubeops/app/pkg/config"
)

func NewGrpcServer(c *config.Server, svc *service.KubeOpsService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
		grpc.Network(c.Grpc.Network),
		grpc.Address(c.Grpc.Addr),
		grpc.Timeout(c.Grpc.Timeout.AsDuration()),
	}
	gs := grpc.NewServer(opts...)
	api.RegisterKubeOpsApiServer(gs, svc)
	return gs
}
