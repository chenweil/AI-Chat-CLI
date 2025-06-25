package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version 应用程序版本
var Version = "1.0.0"

// BuildDate 构建日期
var BuildDate = "2024-01-XX"

// GitCommit Git提交哈希
var GitCommit = "dev"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  `显示 ai-chat-cli 的版本信息，包括版本号、构建日期和Git提交信息。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ai-chat-cli version %s\n", Version)
		fmt.Printf("构建日期: %s\n", BuildDate)
		fmt.Printf("Git提交: %s\n", GitCommit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
