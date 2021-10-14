package task

import (
	"fmt"
	"log"

	"github.com/pipperman/kubeops/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var taskListCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("server.host")
		port := viper.GetInt("server.port")
		c := client.NewKubeOpsClient(host, port)
		rs, err := c.ListResult()
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range rs {
			out := fmt.Sprintf("%s  %s   %s   %t  %t",
				r.Id, r.StartTime, r.EndTime, r.Finished, r.Success)
			fmt.Println(out)
		}
	},
}

func init() {
}
