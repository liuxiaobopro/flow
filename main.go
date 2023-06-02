package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/liuxiaobopro/flow/api"

	"github.com/spf13/cobra"
)

func main() {
	// 创建根命令
	rootCmd := &cobra.Command{
		Use:   "flow",
		Short: "Quick start go - A command-line tool to assist with Go development",
		Long:  "flow is a command-line tool that provides quick and handy features for Go development.",
		Run: func(cmd *cobra.Command, args []string) {
			// 默认的命令行操作
			fmt.Println("Welcome to flow! Use --help to see available commands.")
		},
	}

	// 添加子命令：api
	apiCmd := &cobra.Command{
		Use:   "api",
		Short: "Generate API interface based on a file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			file := args[0]
			// 判断文件是否存在
			if _, err := os.Stat(file); os.IsNotExist(err) {
				fmt.Println("File does not exist.")
				os.Exit(1)
			}

			// 判断是否是.json结尾
			if len(file) < 5 || file[len(file)-5:] != ".json" {
				fmt.Println("File format is not correct.")
				os.Exit(1)
			}

			// 读取文件并映射
			var conf api.Config
			f, err := os.Open(file)
			if err != nil {
				fmt.Println("File open failed.")
				os.Exit(1)
			}
			defer f.Close()

			// 解析json
			b, err := io.ReadAll(f)
			if err != nil {
				fmt.Println("File read failed.")
				os.Exit(1)
			}

			// 转换为结构体
			if err := json.Unmarshal(b, &conf); err != nil {
				fmt.Println("File format is not correct.")
				os.Exit(1)
			}
			api.Run(&conf)
		},
	}
	rootCmd.AddCommand(apiCmd)

	// 执行命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
