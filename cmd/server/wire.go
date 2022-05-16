// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/pipperman/kubeops/app/kubeops/biz"
	"github.com/pipperman/kubeops/app/kubeops/data"
	"github.com/pipperman/kubeops/app/kubeops/server"
	"github.com/pipperman/kubeops/app/kubeops/service"
	"github.com/pipperman/kubeops/app/pkg/config"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// initApp init kratos application.
func initApp(*config.Server, *config.Cache, log.Logger, *tracesdk.TracerProvider) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
