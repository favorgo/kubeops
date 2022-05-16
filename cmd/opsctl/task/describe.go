package task

import (
	"fmt"
	"github.com/pipperman/kubeops/app/opsctl"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var taskDescribeCmd = &cobra.Command{
	Use: "describe",
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("server.host")
		port := viper.GetInt("server.port")
		c := opsctl.NewKubeOpsClient(host, port)
		if len(args) < 1 {
			log.Fatal("task_id missing")
		}
		taskId := args[0]
		result, err := c.GetResult(taskId)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(fmt.Sprintf("result_id: %s", result.Id))
		fmt.Println(fmt.Sprintf("star_time: %s", result.StartTime))
		fmt.Println(fmt.Sprintf("end_time: %s", result.EndTime))
		fmt.Println(fmt.Sprintf("finished: %t", result.Finished))
		fmt.Println(fmt.Sprintf("success: %t", result.Success))
		fmt.Println(fmt.Sprintf("message:%s", result.Message))
		fmt.Println(fmt.Sprintf("content:"))
		fmt.Println(result.Content)
	},
}

func init() {

}
