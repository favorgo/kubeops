package main

import (
	"flag"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pipperman/kubeops/pkg/config"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "pipperman.kubeops.server"
	// Version is the version of the compiled software.
	Version string
	// configPath is the config flag.
	configPath string
	// baseDir is the config flag.
	baseDir string
	// ansibleConfDir
	ansibleConfDir string
	// ansibleTemplateFilePath
	ansibleTemplateFilePath string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&configPath, "conf", "./conf", "config path, eg: -conf ./conf")
	flag.StringVar(&baseDir, "basedir", "./dist/etc/kubeops", "base director, eg: -basedir /etc/kubeops")
	flag.StringVar(&ansibleConfDir, "ansibleConfDir", "./dist/etc/ansible", "config path, eg: -ansibleConfDir /etc/ansible")
	flag.StringVar(&ansibleTemplateFilePath, "ansibleTemplateFilePath", "./dist/etc/kubeops/", "base director, eg: -ansibleTemplateFilePath /etc/kubeops/")
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
	)

	config.Load(configPath, baseDir, ansibleConfDir, ansibleTemplateFilePath)
	if err := NewPrepare().Handlers(handlers()).Action(); err != nil {
		panic(err)
	}

	// init application
	httpServer := newHttpServer(logger)
	grpcServer := newGrpcServer(logger)
	app := newApp(logger, httpServer, grpcServer)

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
