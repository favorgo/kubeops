package project

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "project",
}

func init() {
	Cmd.AddCommand(projectListCmd)
	Cmd.AddCommand(projectCreateCmd)
}
