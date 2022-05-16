package playbook

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "playbook",
}

func init() {
	Cmd.AddCommand(playbookListCmd)
	Cmd.AddCommand(playbookRunCmd)
}
