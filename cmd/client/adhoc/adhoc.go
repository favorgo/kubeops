package adhoc

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "adhoc",
}

func init() {
	Cmd.AddCommand(adhocRunCmd)
}

