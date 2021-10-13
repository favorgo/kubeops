package config

import (
	"fmt"
	"github.com/pipperman/kubeops/pkg/constant"
	"github.com/spf13/viper"
)

const (
	defaultServerHost     = "127.0.0.1"
	defaultHttpServerPort = 8080
	defaultGrpcServerPort = 9090
	defaultBaseDir        = "/var/kubeops"
	defaultConfigPath     = "/etc/kubeops"
)

func Load(configPath, basedir, ansibleConfDir, ansibleTemplateFilePath string) {
	// add config path
	if configPath != "" {
		viper.AddConfigPath(configPath)
	} else {
		viper.AddConfigPath(defaultConfigPath)
	}

	// set basedir
	if basedir != "" {
		viper.SetDefault("base", basedir)
	} else {
		viper.SetDefault("base", defaultBaseDir)
	}

	// set ansibleConfDir
	viper.SetDefault("ansible_conf_dir", ansibleConfDir)
	// set ansibleTemplateFilePath
	viper.SetDefault("ansible_template_file_path", ansibleTemplateFilePath)

	// set config name and config type
	viper.SetConfigName("app")
	viper.SetConfigType("yml")

	// set default server config
	viper.SetDefault("server", server{
		host: defaultServerHost,
		httpPort: defaultHttpServerPort,
		grpcPort: defaultGrpcServerPort,
	})

	constant.Init()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

type server struct {
	host     string
	grpcPort int
	httpPort int
	timeout  int
}
