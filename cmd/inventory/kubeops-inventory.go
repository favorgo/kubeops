package main

import (
	"flag"
	"fmt"
	"github.com/pipperman/kubeops/pkg/config"
	"os"

	"github.com/pipperman/kubeops/pkg/inventory"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	rootCmd.Flags().Bool("list", false, "")
}

var rootCmd = &cobra.Command{
	Use:   "inventory",
	Short: "A inventory provider for kubeops",
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("server.host")
		port := viper.GetInt("server.port")
		provider := inventory.NewKubeOpsInventoryProvider(host, port)
		list, err := cmd.Flags().GetBool("list")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if list {
			result, err := provider.ListHandler()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(result)
			os.Exit(0)
		}
		fmt.Println(host)
		fmt.Println(port)
	},
}

func main() {
	config.Load(configPath, baseDir, ansibleConfDir, ansibleTemplateFilePath)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
