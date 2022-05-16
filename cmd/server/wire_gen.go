// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pipperman/kubeops/app/kubeops/biz"
	"github.com/pipperman/kubeops/app/kubeops/data"
	"github.com/pipperman/kubeops/app/kubeops/server"
	"github.com/pipperman/kubeops/app/kubeops/service"
	"github.com/pipperman/kubeops/app/pkg/config"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(configServer *config.Server, cache *config.Cache, logger log.Logger, tracerProvider *trace.TracerProvider) (*kratos.App, func(), error) {
	projectRepo := data.NewProjectRepo(logger)
	poolRepo := data.NewPoolRepo(configServer, logger)
	cacheRepo := data.NewCacheRepo(cache)
	runnerFactoryRepo := data.NewRunnerFactoryRepo(projectRepo)
	runnerRepo := data.NewRunnerRepo(poolRepo, cacheRepo, runnerFactoryRepo, logger)
	inventoryRepo := biz.NewInventoryRepo(cacheRepo, logger)
	resultRepo := data.NewResultRepo(cacheRepo, logger)
	kubeOpsUseCase := biz.NewKubeOpsUseCase(projectRepo, runnerRepo, inventoryRepo, resultRepo, logger)
	kubeOpsService := service.NewKubeOpsService(kubeOpsUseCase, logger)
	httpServer := server.NewHttpServer(configServer, kubeOpsService, logger)
	grpcServer := server.NewGrpcServer(configServer, kubeOpsService, logger)
	app := newApp(logger, httpServer, grpcServer)
	return app, func() {
	}, nil
}