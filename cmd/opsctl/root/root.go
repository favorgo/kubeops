package root

import (
	"github.com/pipperman/kubeops/cmd/opsctl/adhoc"
	"github.com/pipperman/kubeops/cmd/opsctl/playbook"
	"github.com/pipperman/kubeops/cmd/opsctl/project"
	"github.com/pipperman/kubeops/cmd/opsctl/task"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "opsctl",
	Short: "A kubeops client cli",
}

func init() {
	Cmd.AddCommand(project.Cmd)
	Cmd.AddCommand(playbook.Cmd)
	Cmd.AddCommand(task.Cmd)
	Cmd.AddCommand(adhoc.Cmd)
}
