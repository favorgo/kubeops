package project

import (
	"errors"
	"fmt"
	"github.com/pipperman/kubeops/app/opsctl"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var projectCreateCmd = &cobra.Command{
	Use: "create",
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("server.host")
		port := viper.GetInt("server.port")
		c := opsctl.NewKubeOpsClient(host, port)
		if len(args) < 0 {
			log.Fatal("invalid project source")
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		if name == "" {
			log.Fatal(errors.New("you must provide a valid project name"))
		}
		source := args[0]
		p, err := c.CreateProject(name, source)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(fmt.Sprintf("%s created", p.Name))
	},
}

func init() {
	projectCreateCmd.Flags().String("name", "", "")
}
