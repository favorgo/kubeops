package project

import (
	"fmt"
	"github.com/pipperman/kubeops/app/opsctl"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var projectListCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("server.host")
		port := viper.GetInt("server.port")
		c := opsctl.NewKubeOpsClient(host, port)
		ps, err := c.ListProject()
		if err != nil {
			log.Fatal(err)
		}
		for _, p := range ps {
			fmt.Println(p.Name)
		}
	},
}

func init() {
}
