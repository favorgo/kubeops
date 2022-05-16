package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	api "github.com/pipperman/kubeops/api/v1"
	"github.com/pipperman/kubeops/app/kubeops/service"
	"github.com/pipperman/kubeops/app/pkg/config"
)

func NewHttpServer(c *config.Server, svc *service.KubeOpsService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
		http.Network(c.Http.Network),
		http.Address(c.Http.Addr),
		http.Timeout(c.Http.Timeout.AsDuration()),
	}
	srv := http.NewServer(opts...)
	api.RegisterKubeOpsApiHTTPServer(srv, svc)
	return srv
}
