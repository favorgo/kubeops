package main

import (
	"flag"
	"os"

	"github.com/pipperman/kubeops/cmd/client/root"
	"github.com/pipperman/kubeops/pkg/config"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "pipper.kubeops.inventory"
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
)

func init() {
	flag.StringVar(&configPath, "conf", "./conf", "config path, eg: -conf ./conf")
	flag.StringVar(&baseDir, "basedir", "./dist/etc/kubeops", "base director, eg: -basedir /etc/kubeops")
	flag.StringVar(&ansibleConfDir, "ansibleConfDir", "/etc/ansible", "config path, eg: -ansibleConfDir /etc/ansible")
	flag.StringVar(&ansibleTemplateFilePath, "ansibleTemplateFilePath", "./dist/etc/kubeops/", "base director, eg: -ansibleTemplateFilePath /etc/kubeops/")
}

func main() {
	config.Load(configPath, baseDir, ansibleConfDir, ansibleTemplateFilePath)
	if err := root.Cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
