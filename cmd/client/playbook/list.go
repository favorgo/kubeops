package playbook

import (
	"fmt"
	"log"

	"github.com/pipperman/kubeops/api"
	"github.com/pipperman/kubeops/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var playbookListCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("server.host")
		port := viper.GetInt("server.port")
		c := client.NewkubeOpsClient(host, port)
		ps, err := c.ListProject()
		if err != nil {
			log.Fatal(err.Error())
		}
		sp, err := cmd.Flags().GetString("project")
		if err != nil {
			log.Fatal(err)
		}
		displayProjects := make([]*api.Project, 0)
		if sp != "" {
			for _, p := range ps {
				if p.Name == sp {
					displayProjects = append(displayProjects, p)
				}
			}
		} else {
			displayProjects = ps
		}
		for _, p := range displayProjects {
			fmt.Println(p.Name)
			for _, pb := range p.Playbooks {
				str := fmt.Sprintf("-- %s", pb)
				fmt.Println(str)
			}
		}
	},
}

func init() {
	playbookListCmd.Flags().StringP("project", "p", "", "specify project name")
}
