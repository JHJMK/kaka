package cmd

import (
	"encoding/json"
	"fmt"
	"kaka/cmd/cmd/config"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init kaka tool",
	Run: func(cmd *cobra.Command, args []string) {
		// 在当前目录下创建一个config目录
		stat, err := os.Stat("./config")
		if err != nil {
			if os.IsNotExist(err) {
				os.Mkdir("./config", 0755)
			} else {
				fmt.Println("init config path error:", err)
				os.Exit(1)
			}
		} else {
			// 只有当stat不为nil时才检查是否为目录
			if !stat.IsDir() {
				fmt.Println("init config dir error: exits the same name file")
				os.Exit(1)
			}
		}

		// 创建data目录用来存放kafka
		//stat, err = os.Stat("./data")
		//if err != nil {
		//	if os.IsNotExist(err) {
		//		os.Mkdir("./data", 0755)
		//	} else {
		//		fmt.Println("init data path error:", err)
		//		os.Exit(1)
		//	}
		//} else {
		//	// 只有当stat不为nil时才检查是否为目录
		//	if !stat.IsDir() {
		//		fmt.Println("init data dir error: exits the same name file")
		//		os.Exit(1)
		//	}
		//}

		// 创建一个config.json文件
		_, err = os.Stat("./config/config.json")
		if err != nil {
			// 如果没有就创建
			if os.IsNotExist(err) {
				f, err := os.Create("./config/config.json")
				if err != nil {
					fmt.Println("init config.json error")
					os.Exit(1)
				}
				defer f.Close()
				// 创建config的空对象，并写入文件
				var c config.Config
				c.ManageHost = []config.Host{
					{
						IP:       "input you become manage host",
						Port:     "22",
						User:     "input you become manage user",
						Password: "If you configure sshpass password-free for user interoperability, ignore it",
					},
				}
				c.InstallPath = "/opt/kaka"
				c.KafkaLogDir = "/data/kaka-logs"
				bytes, err := json.MarshalIndent(&c, "", "  ")
				if err != nil {
					fmt.Println("init config.json error:", err)
					os.Exit(1)
				}
				os.WriteFile("./config/config.json", bytes, 0644)
				fmt.Println("please fill in the config/config.json")
			} else {
				fmt.Println("init config.json error:", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
