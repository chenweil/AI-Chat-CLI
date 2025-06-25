package cmd

import (
	"fmt"
	"os"
	"strings"

	"ai-chat-cli/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "配置管理",
	Long:  `管理 ai-chat-cli 的配置文件，包括初始化、查看和设置配置项。`,
}

// configInitCmd 初始化配置
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化配置文件",
	Long:  `创建默认的配置文件，包含常用的AI提供商设置。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取配置文件路径
		configPath, err := config.GetDefaultConfigPath()
		if err != nil {
			fmt.Printf("错误：无法获取配置文件路径: %v\n", err)
			return
		}

		// 检查配置文件是否已存在
		if _, err := os.Stat(configPath); err == nil {
			fmt.Printf("配置文件已存在: %s\n", configPath)
			fmt.Println("如果要重新初始化，请先删除现有配置文件。")
			return
		}

		// 创建示例配置
		if err := createExampleConfig(configPath); err != nil {
			fmt.Printf("错误：创建配置文件失败: %v\n", err)
			return
		}

		fmt.Printf("✓ 配置文件已创建: %s\n", configPath)
		fmt.Println("\n请编辑配置文件，设置您的API密钥：")
		fmt.Printf("  - OpenAI API密钥: 设置 OPENAI_API_KEY 环境变量\n")
		fmt.Printf("  - Anthropic API密钥: 设置 ANTHROPIC_API_KEY 环境变量\n")
		fmt.Println("\n或者直接在配置文件中设置 api_key 字段。")
	},
}

// configShowCmd 显示配置
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "显示当前配置",
	Long:  `显示当前的配置内容，敏感信息（如API密钥）将被隐藏。`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("错误：加载配置失败: %v\n", err)
			return
		}

		fmt.Println("当前配置:")
		fmt.Printf("  默认提供商: %s\n", cfg.Default.Provider)
		fmt.Printf("  流式输出: %t\n", cfg.Default.Stream)
		fmt.Printf("  最大重试: %d\n", cfg.Advanced.MaxRetries)
		fmt.Printf("  超时时间: %d秒\n", cfg.Advanced.Timeout)
		fmt.Printf("  成本限制: $%.2f\n", cfg.Advanced.CostLimit)

		fmt.Println("\n已配置的提供商:")
		for name, provider := range cfg.Providers {
			apiKeyStatus := "未设置"
			if provider.APIKey != "" {
				apiKeyStatus = "已设置"
			} else if os.Getenv(strings.ToUpper(name)+"_API_KEY") != "" {
				apiKeyStatus = "环境变量"
			}
			fmt.Printf("  %s:\n", name)
			fmt.Printf("    模型: %s\n", provider.Model)
			fmt.Printf("    API密钥: %s\n", apiKeyStatus)
			fmt.Printf("    最大Token: %d\n", provider.MaxTokens)
		}
	},
}

// configSetCmd 设置配置项
var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "设置配置项",
	Long: `设置指定的配置项。

示例:
  ai-chat-cli config set default.provider openai
  ai-chat-cli config set default.stream true
  ai-chat-cli config set providers.openai.model gpt-4`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		// 设置配置值
		viper.Set(key, value)

		// 保存配置
		if err := viper.WriteConfig(); err != nil {
			fmt.Printf("错误：保存配置失败: %v\n", err)
			return
		}

		fmt.Printf("✓ 已设置 %s = %s\n", key, value)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
}

// createExampleConfig 创建示例配置文件
func createExampleConfig(filename string) error {
	configContent := `# AI Chat CLI 配置文件

# AI提供商配置
providers:
  openai:
    # API密钥（推荐使用环境变量 OPENAI_API_KEY）
    api_key: ""
    base_url: "https://api.openai.com/v1"
    model: "gpt-4o"
    max_tokens: 4096

  anthropic:
    # API密钥（推荐使用环境变量 ANTHROPIC_API_KEY）
    api_key: ""
    base_url: "https://api.anthropic.com"
    model: "claude-3-sonnet-20240229"
    max_tokens: 4096

# 默认设置
default:
  provider: "openai"   # 默认使用的AI提供商
  stream: true         # 是否启用流式输出

# 高级设置
advanced:
  max_retries: 3       # 最大重试次数
  timeout: 30          # 请求超时时间（秒）
  cost_limit: 10.0     # 每日成本限制（美元）
  save_history: true   # 是否保存对话历史
  history_length: 10   # 保存的历史对话数量

# 日志设置
logging:
  level: "info"        # 日志级别: debug, info, warn, error
  requests: false      # 是否记录API请求日志
`

	// 确保目录存在
	if err := os.MkdirAll(strings.TrimSuffix(filename, "/config.yaml"), 0755); err != nil {
		return err
	}

	// 写入配置文件
	return os.WriteFile(filename, []byte(configContent), 0644)
}
