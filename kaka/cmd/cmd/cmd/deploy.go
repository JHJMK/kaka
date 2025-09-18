package cmd

import (
	"encoding/json"
	"fmt"
	"kaka/cmd/cmd/config"
	"kaka/cmd/cmd/handler"
	"os"

	"github.com/spf13/cobra"
)

// TODO：防呆设计，检测目标文件夹是否为空，不为空给提示是否清空目标目录
var deployCmd = &cobra.Command{
	Use:   "deploy",      //命令
	Short: "deploy kaka", //简短描述
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
		// 各个节点创建目录
		handler.CreateDir(c)
		// 复制安装包
		err = handler.CopyFile(c)
		if err != nil {
			fmt.Println("copy binary error:", err)
			os.Exit(1)
		}
		// 解压安装包
		err = handler.DecompressionAndRename(c)
		if err != nil {
			fmt.Println("decompression error:", err)
			os.Exit(1)
		}
		// 配置kafka配置文件
		err = handler.ConfigKafka(c)
		if err != nil {
			fmt.Println("config kafka error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
