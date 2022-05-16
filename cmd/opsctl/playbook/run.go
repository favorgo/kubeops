package playbook

import (
	"errors"
	"fmt"
	"github.com/pipperman/kubeops/app/opsctl"
	"io/ioutil"
	"log"
	"os"
	"time"

	api "github.com/pipperman/kubeops/api/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var playbookRunCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("server.host")
		port := viper.GetInt("server.port")
		c := opsctl.NewKubeOpsClient(host, port)
		project, err := cmd.Flags().GetString("project")
		if err != nil {
			log.Fatal(err)
		}
		if project == "" {
			log.Fatal(errors.New("you must specify project name"))
		}
		inventoryPath, err := cmd.Flags().GetString("inventory")
		if err != nil {
			log.Fatal(err)
		}
		content, err := ioutil.ReadFile(inventoryPath)
		if err != nil {
			log.Fatal(err)
		}
		var inventory api.Inventory
		err = yaml.Unmarshal(content, &inventory)
		if err != nil {
			log.Fatal(err)
		}
		if len(args) < 1 {
			log.Fatal("invalid playbook name")
		}
		playbook, tag := args[0], args[1]
		result, err := c.RunPlaybook(project, playbook, tag, &inventory)
		if err != nil {
			log.Fatal(err)
		}
		detach, err := cmd.Flags().GetBool("detach")
		if err != nil {
			log.Fatal(err)
		}
		if detach {
			fmt.Println(result.Id)
		} else {
			sign := make(chan int)
			go func() {
				for {
					result, err = c.GetResult(result.Id)
					if err != nil {
						log.Fatal(err)
					}
					if result.Finished {
						sign <- 1
					}
					time.Sleep(1 * time.Second)
				}
			}()
			err := c.WatchRun(result.Id, os.Stdout)
			if err != nil {
				log.Fatal(err)
			}
			select {
			case signal := <-sign:
				if signal == 1 {
					if !result.Success {
						log.Fatal(result.Message)
					}
				}
			}
		}

	},
}

func init() {
	playbookRunCmd.Flags().StringP("project", "prj", "", "specify project name")
	playbookRunCmd.Flags().BoolP("detach", "d", false, "run in background")
	playbookRunCmd.Flags().StringP("inventory", "inv", "", "specify inventory file path")
}
