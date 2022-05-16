package main

import (
	"flag"
	kconf "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"gopkg.in/yaml.v2"
	"os"

	"github.com/pipperman/kubeops/app/pkg/config"
	"github.com/pipperman/kubeops/cmd/opsctl/root"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "pipperman.kubeops.client"
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
	// ansibleVariablesName
	ansibleVariablesName string
)

func init() {
	flag.StringVar(&configPath, "conf", "./conf", "config path, eg: -conf ./conf")
	flag.StringVar(&baseDir, "basedir", "./dist/etc/kubeops", "base director, eg: -basedir /etc/errors")
	flag.StringVar(&ansibleConfDir, "ansibleConfDir", "/etc/ansible", "config path, eg: -ansibleConfDir /etc/ansible")
	flag.StringVar(&ansibleTemplateFilePath, "ansibleTemplateFilePath", "./dist/etc/kubeops/", "base director, eg: -ansibleTemplateFilePath /etc/errors/")
	flag.StringVar(&ansibleVariablesName, "variablesName", "variable.yml", "variable name, eg: -variablesName variable.yml")
}

func main() {
	c := kconf.New(
		kconf.WithSource(
			file.NewSource(configPath),
		),
		kconf.WithDecoder(func(kv *kconf.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)
	var bc config.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	if err := root.Cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
