package playbook

import (
	"errors"
	"fmt"
	"github.com/pipperman/kubeops/api"
	"github.com/pipperman/kubeops/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var playbookRunCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("server.host")
		port := viper.GetInt("server.port")
		c := client.NewkubeOpsClient(host, port)
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
		backend, err := cmd.Flags().GetBool("b")
		if err != nil {
			log.Fatal(err)
		}
		if backend {
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
			case a := <-sign:
				if a == 1 {
					if !result.Success {
						log.Fatal(result.Message)
					}
				}
			}
		}

	},
}

func init() {
	playbookRunCmd.Flags().StringP("project", "p", "", "specify project name")
	playbookRunCmd.Flags().BoolP("b", "b", false, "run in background")
	playbookRunCmd.Flags().StringP("inventory", "i", "", "specify inventory file path")
}
