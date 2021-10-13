package task

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use: "task",
}

func init() {
	Cmd.AddCommand(taskListCmd)
	Cmd.AddCommand(taskDescribeCmd)
}
