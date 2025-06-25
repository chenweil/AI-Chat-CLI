package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config 应用程序配置结构
type Config struct {
	// AI提供商配置
	Providers map[string]ProviderConfig `mapstructure:"providers" yaml:"providers" json:"providers"`

	// 默认设置
	Default DefaultConfig `mapstructure:"default" yaml:"default" json:"default"`

	// 高级设置
	Advanced AdvancedConfig `mapstructure:"advanced" yaml:"advanced" json:"advanced"`

	// 日志设置
	Logging LoggingConfig `mapstructure:"logging" yaml:"logging" json:"logging"`
}

// ProviderConfig AI提供商配置
type ProviderConfig struct {
	APIKey    string            `mapstructure:"api_key" yaml:"api_key" json:"api_key"`
	BaseURL   string            `mapstructure:"base_url" yaml:"base_url" json:"base_url"`
	Model     string            `mapstructure:"model" yaml:"model" json:"model"`
	MaxTokens int               `mapstructure:"max_tokens" yaml:"max_tokens" json:"max_tokens"`
	Extra     map[string]string `mapstructure:"extra" yaml:"extra" json:"extra"`
}

// DefaultConfig 默认配置
type DefaultConfig struct {
	Provider string `mapstructure:"provider" yaml:"provider" json:"provider"`
	Model    string `mapstructure:"model" yaml:"model" json:"model"`
	Stream   bool   `mapstructure:"stream" yaml:"stream" json:"stream"`
}

// AdvancedConfig 高级配置
type AdvancedConfig struct {
	MaxRetries    int     `mapstructure:"max_retries" yaml:"max_retries" json:"max_retries"`
	Timeout       int     `mapstructure:"timeout" yaml:"timeout" json:"timeout"`
	CostLimit     float64 `mapstructure:"cost_limit" yaml:"cost_limit" json:"cost_limit"`
	SaveHistory   bool    `mapstructure:"save_history" yaml:"save_history" json:"save_history"`
	HistoryLength int     `mapstructure:"history_length" yaml:"history_length" json:"history_length"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level    string `mapstructure:"level" yaml:"level" json:"level"`
	File     string `mapstructure:"file" yaml:"file" json:"file"`
	Requests bool   `mapstructure:"requests" yaml:"requests" json:"requests"`
}

// GlobalConfig 全局配置实例
var GlobalConfig *Config

// Load 加载配置
func Load() (*Config, error) {
	cfg := &Config{}

	// 设置默认值
	setDefaults()

	// 尝试解析配置
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	GlobalConfig = cfg
	return cfg, nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 检查默认提供商是否存在
	if c.Default.Provider != "" {
		if _, exists := c.Providers[c.Default.Provider]; !exists {
			return fmt.Errorf("默认提供商 '%s' 未配置", c.Default.Provider)
		}
	}

	// 验证API密钥
	for name, provider := range c.Providers {
		if provider.APIKey == "" {
			// 尝试从环境变量获取
			envKey := fmt.Sprintf("%s_API_KEY", name)
			if os.Getenv(envKey) == "" {
				return fmt.Errorf("提供商 '%s' 缺少API密钥，请设置 %s 环境变量或在配置文件中指定", name, envKey)
			}
		}
	}

	return nil
}

// Save 保存配置到文件
func (c *Config) Save(filename string) error {
	// 确保目录存在
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 写入配置文件
	viper.SetConfigFile(filename)
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("保存配置文件失败: %w", err)
	}

	return nil
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 提供商默认配置
	viper.SetDefault("providers.openai.base_url", "https://api.openai.com/v1")
	viper.SetDefault("providers.openai.model", "gpt-4o")
	viper.SetDefault("providers.openai.max_tokens", 4096)

	viper.SetDefault("providers.anthropic.base_url", "https://api.anthropic.com")
	viper.SetDefault("providers.anthropic.model", "claude-3-sonnet-20240229")
	viper.SetDefault("providers.anthropic.max_tokens", 4096)

	// 默认设置
	viper.SetDefault("default.provider", "openai")
	viper.SetDefault("default.stream", true)

	// 高级设置
	viper.SetDefault("advanced.max_retries", 3)
	viper.SetDefault("advanced.timeout", 30)
	viper.SetDefault("advanced.cost_limit", 10.0)
	viper.SetDefault("advanced.save_history", true)
	viper.SetDefault("advanced.history_length", 10)

	// 日志设置
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.requests", false)
}

// GetProvider 获取指定提供商配置
func (c *Config) GetProvider(name string) (ProviderConfig, error) {
	provider, exists := c.Providers[name]
	if !exists {
		return ProviderConfig{}, fmt.Errorf("提供商 '%s' 未配置", name)
	}

	// 如果配置中没有API密钥，尝试从环境变量获取
	if provider.APIKey == "" {
		envKey := fmt.Sprintf("%s_API_KEY", name)
		if apiKey := os.Getenv(envKey); apiKey != "" {
			provider.APIKey = apiKey
		}
	}

	return provider, nil
}

// GetDefaultConfigPath 获取默认配置文件路径
func GetDefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".ai-chat-cli", "config.yaml"), nil
}

// LoadConfig 加载配置文件
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置到结构体
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	return cfg, nil
}
