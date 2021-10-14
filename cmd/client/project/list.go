package project

import (
	"fmt"
	"log"

	"github.com/pipperman/kubeops/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var projectListCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("server.host")
		port := viper.GetInt("server.port")
		c := client.NewKubeOpsClient(host, port)
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
