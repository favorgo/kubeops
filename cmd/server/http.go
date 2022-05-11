package main

import (
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/pipperman/kubeops/api"
	"github.com/pipperman/kubeops/pkg/server"
	"github.com/spf13/viper"
)

func newHttpServer(logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
		http.Address(fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.httpPort"))),
		http.Timeout(time.Duration(viper.GetInt("server.timeout"))),
	}
	srv := http.NewServer(opts...)
	kubeOps := server.NewKubeOps(logger, server.WithPoolOption(viper.GetInt("app.queue"), viper.GetInt("worker"), logger))
	api.RegisterKubeOpsApiHTTPServer(srv, kubeOps)
	return srv
}
