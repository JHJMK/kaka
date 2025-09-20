package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"kaka/config"
	"kaka/handler"
	"os"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the terminal",
	Run: func(cmd *cobra.Command, args []string) {
		// 读取配置文件
		bytes, err := os.ReadFile("config/config.json")
		if err != nil {
			fmt.Println("read config.json error")
			os.Exit(1)
		}
		var c config.Config
		err = json.Unmarshal(bytes, &c)
		if err != nil {
			fmt.Println("read config.json error")
			os.Exit(1)
		}
		fmt.Println("read config ...")
		err = handler.Clear(c)
		if err != nil {
			fmt.Println("clear error")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
