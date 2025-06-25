package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"ai-chat-cli/internal/config"

	"github.com/spf13/cobra"
)

// Message 表示对话消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest OpenAI API请求结构
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

// ChatResponse OpenAI API响应结构
type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

var (
	chatProvider string
)

// chatCmd represents the chat command
var simpleChatCmd = &cobra.Command{
	Use:   "chat [问题]",
	Short: "与AI进行对话",
	Long: `与AI模型进行对话交流。可以：

• 直接指定问题：ai-chat-cli chat "你好，介绍一下自己"
• 进入交互模式：ai-chat-cli chat （然后输入问题）
• 指定提供商：ai-chat-cli chat --provider free-oai "问题"

支持的提供商：
• openai (官方API)
• free-oai (第三方兼容API)
• 或任何您在配置文件中定义的提供商`,
	Args: cobra.MaximumNArgs(1),
	Run:  runSimpleChat,
}

func runSimpleChat(cmd *cobra.Command, args []string) {
	// 检查配置
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("❌ 配置加载失败: %v\n", err)
		fmt.Println("💡 请先运行 'ai-chat-cli config init' 初始化配置")
		return
	}

	// 如果没有指定提供商，尝试找到第一个可用的
	if chatProvider == "" {
		for name, providerCfg := range cfg.Providers {
			if providerCfg.APIKey != "" {
				chatProvider = name
				fmt.Printf("💡 自动选择提供商: %s\n", name)
				break
			}
		}
	}

	// 获取指定的提供商配置
	providerCfg, exists := cfg.Providers[chatProvider]
	if !exists {
		fmt.Printf("❌ 提供商 '%s' 未找到\n", chatProvider)
		fmt.Println("📋 可用的提供商:")
		for name := range cfg.Providers {
			fmt.Printf("  • %s\n", name)
		}
		return
	}

	if providerCfg.APIKey == "" {
		fmt.Printf("❌ 提供商 '%s' 的API密钥未设置\n", chatProvider)
		fmt.Printf("💡 请运行以下命令设置API密钥：\n")
		fmt.Printf("   ai-chat-cli config set providers.%s.api_key YOUR_API_KEY\n", chatProvider)
		return
	}

	fmt.Printf("🚀 使用提供商: %s\n", chatProvider)
	if providerCfg.BaseURL != "" && providerCfg.BaseURL != "https://api.openai.com/v1" {
		fmt.Printf("🌐 API地址: %s\n", providerCfg.BaseURL)
	}
	if providerCfg.Model != "" {
		fmt.Printf("🤖 使用模型: %s\n", providerCfg.Model)
	}

	// 初始化对话历史
	var conversationHistory []Message

	if len(args) > 0 {
		// 单次对话模式
		question := args[0]
		err = askQuestionWithHistory(providerCfg, question, &conversationHistory)
		if err != nil {
			fmt.Printf("❌ 对话失败: %v\n", err)
		}
	} else {
		// 交互模式
		runInteractiveChatWithHistory(providerCfg, &conversationHistory)
	}
}

func askQuestionWithHistory(providerCfg config.ProviderConfig, question string, history *[]Message) error {
	fmt.Print("🤖 AI: ")

	// 设置默认值
	model := providerCfg.Model
	if model == "" {
		model = "gpt-3.5-turbo"
	}

	maxTokens := providerCfg.MaxTokens
	if maxTokens == 0 {
		maxTokens = 2000
	}

	baseURL := providerCfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	// 添加用户问题到历史
	*history = append(*history, Message{Role: "user", Content: question})

	// 构建请求（包含完整历史）
	reqBody := ChatRequest{
		Model:       model,
		Messages:    *history, // 发送完整的对话历史
		MaxTokens:   maxTokens,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("构建请求失败: %w", err)
	}

	// 创建HTTP请求
	apiURL := baseURL + "/chat/completions"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+providerCfg.APIKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求发送失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API返回错误 %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return fmt.Errorf("API返回空响应")
	}

	response := chatResp.Choices[0].Message.Content
	fmt.Println(response)

	// 添加AI回复到历史
	*history = append(*history, Message{Role: "assistant", Content: response})

	// 显示使用统计
	usage := chatResp.Usage
	fmt.Printf("\n📊 Token使用: %d (输入: %d, 输出: %d) | 对话轮次: %d\n",
		usage.TotalTokens, usage.PromptTokens, usage.CompletionTokens, len(*history)/2)

	return nil
}

func runInteractiveChatWithHistory(providerCfg config.ProviderConfig, history *[]Message) {
	fmt.Println("🤖 AI Chat CLI - 交互模式 (支持上下文记忆)")
	fmt.Println("💡 输入问题开始对话")
	fmt.Println("💡 特殊命令:")
	fmt.Println("   • quit/exit - 退出程序")
	fmt.Println("   • clear - 清屏")
	fmt.Println("   • reset - 重置对话历史")
	fmt.Println("   • history - 显示对话历史")
	fmt.Println("   • help - 显示帮助")
	fmt.Println("💡 如果输入出现问题，直接按回车重新输入")
	fmt.Println("---")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("👤 你: ")

		if !scanner.Scan() {
			// 处理EOF或其他错误
			if err := scanner.Err(); err != nil {
				fmt.Printf("\n❌ 输入错误: %v\n", err)
			}
			break
		}

		input := strings.TrimSpace(scanner.Text())

		// 检查是否是空输入
		if input == "" {
			continue
		}

		// 检查是否包含不可见字符或控制字符
		cleanInput := cleanInput(input)
		if cleanInput == "" {
			fmt.Println("⚠️  输入包含无效字符，请重新输入")
			continue
		}

		// 检查特殊命令
		lowerInput := strings.ToLower(cleanInput)
		switch lowerInput {
		case "quit", "exit":
			fmt.Println("👋 再见！")
			return
		case "clear":
			fmt.Print("\033[2J\033[H") // ANSI清屏命令
			fmt.Println("🤖 AI Chat CLI - 交互模式 (支持上下文记忆)")
			fmt.Printf("💡 当前对话历史: %d 轮次\n", len(*history)/2)
			fmt.Println("💡 输入问题开始对话，输入 'help' 查看命令")
			fmt.Println("---")
			continue
		case "reset":
			*history = []Message{} // 清空对话历史
			fmt.Println("🔄 对话历史已重置")
			continue
		case "history":
			showHistory(*history)
			continue
		case "help":
			fmt.Println("🆘 可用命令:")
			fmt.Println("   • quit/exit - 退出程序")
			fmt.Println("   • clear - 清屏")
			fmt.Println("   • reset - 重置对话历史")
			fmt.Println("   • history - 显示对话历史")
			fmt.Println("   • help - 显示此帮助")
			fmt.Println("   • 直接输入问题开始对话")
			continue
		}

		// 显示清理后的输入（仅在有差异时）
		if cleanInput != input {
			fmt.Printf("📝 已清理输入: %s\n", cleanInput)
		}

		err := askQuestionWithHistory(providerCfg, cleanInput, history)
		if err != nil {
			fmt.Printf("❌ 对话失败: %v\n", err)
			fmt.Println("💡 请检查网络连接或重试，输入 'help' 查看可用命令")
		}
		fmt.Println()
	}
}

// showHistory 显示对话历史
func showHistory(history []Message) {
	if len(history) == 0 {
		fmt.Println("📝 暂无对话历史")
		return
	}

	fmt.Println("📝 对话历史:")
	for i, msg := range history {
		if msg.Role == "user" {
			fmt.Printf("  %d. 👤 你: %s\n", i/2+1, msg.Content)
		} else if msg.Role == "assistant" {
			fmt.Printf("     🤖 AI: %s\n", truncateString(msg.Content, 100))
		}
	}
	fmt.Printf("📊 总计 %d 轮对话\n", len(history)/2)
}

// truncateString 截断长字符串用于显示
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// cleanInput 清理输入字符串，移除控制字符和不可见字符
func cleanInput(input string) string {
	// 移除首尾空白字符
	cleaned := strings.TrimSpace(input)

	// 检查是否包含有效的可打印字符
	hasValidChar := false
	var result strings.Builder

	for _, r := range cleaned {
		// 保留可打印字符（包括中文）和空格
		if r == ' ' || r == '\t' || (r >= 32 && r <= 126) || r > 127 {
			result.WriteRune(r)
			if r != ' ' && r != '\t' {
				hasValidChar = true
			}
		}
		// 跳过控制字符和其他不可见字符
	}

	if !hasValidChar {
		return ""
	}

	// 压缩连续的空格
	finalResult := strings.Join(strings.Fields(result.String()), " ")
	return finalResult
}

func init() {
	rootCmd.AddCommand(simpleChatCmd)

	// 添加提供商选择参数
	simpleChatCmd.Flags().StringVarP(&chatProvider, "provider", "p", "", "指定AI提供商 (如: openai, free-oai)")
}
