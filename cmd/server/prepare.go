package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"os"
	"os/exec"
	"text/template"

	"github.com/pipperman/kubeops/pkg/constant"
	"github.com/spf13/viper"
)

type Handler func() error

type Prepare interface {
	Handlers([]Handler) Prepare
	Action() error
}

var _ Prepare = &prepare{}

type prepare struct {
	handlers []Handler
}

func (p *prepare) Action() error  {
	for _, handler := range p.handlers {
		if err := handler(); err != nil {
			return err
		}
	}
	return nil
}

func (p *prepare) Handlers(handlers []Handler) Prepare {
	p.handlers = handlers
	return p
}

func NewPrepare() Prepare {
	return &prepare{}
}

func handlers() []Handler {
	return []Handler{
		makeDataDir,
		makeCacheDir,
		makeKeyDir,
		makeAnsibleCfgDir,
		lookUpAnsibleBinPath,
		lookupKubeOpsInventoryBinPath,
		cleanWorkPath,
		renderAnsibleConfig,
	}
}

func makeDataDir() error {
	return os.MkdirAll(constant.ProjectDir, 0755)

}

func makeAnsibleCfgDir() error {
	return os.MkdirAll(constant.AnsibleConfDir, 0755)
}

func makeCacheDir() error {
	return os.MkdirAll(constant.CacheDir, 0755)
}

func makeKeyDir() error {
	return os.MkdirAll(constant.KeyDir, 0755)
}

func lookUpAnsibleBinPath() error {
	_, err := exec.LookPath(constant.AnsiblePlaybookBinPath)
	if err != nil {
		return err
	}
	return nil
}

func lookupKubeOpsInventoryBinPath() error {
	_, err := exec.LookPath(constant.InventoryProviderBinPath)
	if err != nil {
		return err
	}
	return nil
}

func cleanWorkPath() error {
	_ = os.RemoveAll(constant.WorkDir)
	return nil
}

func renderAnsibleConfig() error {
	tmpl := constant.AnsibleTemplateFilePath
	file, err := os.OpenFile(constant.AnsibleConfPath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		return err
	}
	data := viper.GetStringMap("ansible")
	if err := t.Execute(file, data); err != nil {
		return err
	}
	return nil
}

func newApp(logger log.Logger, hs *http.Server, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
	)
}