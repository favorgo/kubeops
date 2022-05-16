package adhoc

import (
	"fmt"
	"github.com/pipperman/kubeops/app/opsctl"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/ghodss/yaml"
	api "github.com/pipperman/kubeops/api/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var adhocRunCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("server.host")
		port := viper.GetInt("server.port")
		c := opsctl.NewKubeOpsClient(host, port)
		module, _ := cmd.Flags().GetString("module")
		pattern, _ := cmd.Flags().GetString("pattern")
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
		var param string
		if len(args) > 0 {
			param = args[0]
		}
		result, err := c.RunAdhoc(pattern, module, param, &inventory)
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
	adhocRunCmd.Flags().BoolP("b", "b", false, "run in background")
	adhocRunCmd.Flags().StringP("inventory", "i", "", "specify inventory file path")
	adhocRunCmd.Flags().StringP("module", "m", "shell", "specify ansible module")
	adhocRunCmd.Flags().StringP("pattern", "p", "all", "specify inventory pattern")
}
