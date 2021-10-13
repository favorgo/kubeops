package main

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/spf13/viper"
	"net"
	"time"

	"github.com/pipperman/kubeops/api"
	"github.com/pipperman/kubeops/pkg/server"
)

func newTcpListener(address string) (*net.Listener, error) {
	s, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func newGrpcServer() *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
		grpc.Network("tcp"),
		grpc.Address(fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.grpcPort"))),
		grpc.Timeout(time.Duration(viper.GetInt("server.timeout"))),
	}
	gs := grpc.NewServer(opts...)
	kubeOps := server.NewKubeOps()
	api.RegisterKubeOpsApiServer(gs, kubeOps)
	return gs
}
