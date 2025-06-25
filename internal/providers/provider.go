package providers

import (
	"context"
	"io"
)

// Message 表示一条对话消息
type Message struct {
	Role    string `json:"role"`    // "user", "assistant", "system"
	Content string `json:"content"` // 消息内容
}

// ChatRequest 对话请求
type ChatRequest struct {
	Messages    []Message `json:"messages"`    // 对话历史
	Model       string    `json:"model"`       // 使用的模型
	MaxTokens   int       `json:"max_tokens"`  // 最大token数
	Temperature float64   `json:"temperature"` // 温度参数
	Stream      bool      `json:"stream"`      // 是否流式响应
}

// ChatResponse 对话响应
type ChatResponse struct {
	Content      string `json:"content"`       // 响应内容
	Model        string `json:"model"`         // 使用的模型
	Usage        Usage  `json:"usage"`         // 使用统计
	FinishReason string `json:"finish_reason"` // 结束原因
}

// Usage 使用统计
type Usage struct {
	PromptTokens     int     `json:"prompt_tokens"`     // 输入token数
	CompletionTokens int     `json:"completion_tokens"` // 输出token数
	TotalTokens      int     `json:"total_tokens"`      // 总token数
	Cost             float64 `json:"cost"`              // 估算成本
}

// StreamChunk 流式响应的数据块
type StreamChunk struct {
	Content string `json:"content"` // 增量内容
	Done    bool   `json:"done"`    // 是否完成
	Error   error  `json:"-"`       // 错误信息
}

// Provider AI提供商接口
type Provider interface {
	// GetName 获取提供商名称
	GetName() string

	// Chat 发送对话请求（非流式）
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)

	// ChatStream 发送对话请求（流式）
	ChatStream(ctx context.Context, req *ChatRequest) (<-chan StreamChunk, error)

	// GetModels 获取可用模型列表
	GetModels(ctx context.Context) ([]string, error)

	// ValidateConfig 验证配置
	ValidateConfig() error
}

// StreamReader 流式响应读取器
type StreamReader interface {
	io.ReadCloser
	// ReadChunk 读取一个数据块
	ReadChunk() (*StreamChunk, error)
}

// ProviderError 提供商错误
type ProviderError struct {
	Provider string // 提供商名称
	Code     string // 错误代码
	Message  string // 错误消息
	Cause    error  // 原始错误
}

func (e *ProviderError) Error() string {
	if e.Cause != nil {
		return e.Provider + ": " + e.Message + " (" + e.Cause.Error() + ")"
	}
	return e.Provider + ": " + e.Message
}

func (e *ProviderError) Unwrap() error {
	return e.Cause
}

// NewProviderError 创建提供商错误
func NewProviderError(provider, code, message string, cause error) *ProviderError {
	return &ProviderError{
		Provider: provider,
		Code:     code,
		Message:  message,
		Cause:    cause,
	}
}
