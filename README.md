# AI Chat CLI

一个强大的Go语言AI对话命令行工具，支持多AI提供商和上下文记忆。

## ✨ 特性

- 🤖 **多AI提供商支持** - OpenAI、Claude及兼容API
- 💬 **上下文记忆** - 支持连续对话，AI能记住前面的内容
- 🔧 **灵活配置** - 支持自定义API地址、模型、参数
- 🚀 **简单易用** - 命令行直接对话或交互模式
- 📊 **Token统计** - 实时显示使用量和成本
- 🌐 **第三方API支持** - 支持各种OpenAI兼容的API服务

## 🚀 快速开始

### 安装

```bash
# 克隆项目
git clone https://github.com/YOUR_USERNAME/ai-chat-cli.git
cd ai-chat-cli

# 编译
go build -o ai-chat-cli

# 或者直接运行
go run main.go
```

### 配置

```bash
# 初始化配置文件
./ai-chat-cli config init

# 设置API密钥
./ai-chat-cli config set providers.openai.api_key sk-your-api-key-here

# 或设置第三方API
./ai-chat-cli config set providers.free-oai.api_key your-api-key
./ai-chat-cli config set providers.free-oai.base_url https://api.example.com/v1
./ai-chat-cli config set providers.free-oai.model gpt-4.1-nano
```

## 💡 使用方法

### 直接对话模式

```bash
# 直接提问
./ai-chat-cli chat "介绍一下Go语言的特性"

# 指定提供商
./ai-chat-cli chat --provider free-oai "写一首关于编程的诗"
```

### 交互模式

```bash
# 进入交互模式
./ai-chat-cli chat

# 交互模式支持的命令：
# - quit/exit: 退出
# - clear: 清屏
# - reset: 重置对话历史
# - history: 显示对话历史
# - help: 显示帮助
```

## 🔧 配置文件示例

配置文件位置：`~/.ai-chat-cli/config.yaml`

```yaml
providers:
  openai:
    api_key: "sk-your-openai-key"
    base_url: "https://api.openai.com/v1"
    model: "gpt-3.5-turbo"
    max_tokens: 2000
    
  free-oai:
    api_key: "your-key"
    base_url: "https://api.lianwusuoai.top/v1"
    model: "gpt-4.1-nano"
    max_tokens: 8192

default:
  provider: "openai"
  model: "gpt-3.5-turbo"
  max_tokens: 2000
  temperature: 0.7

advanced:
  timeout: 30
  retry_times: 3

logging:
  level: "info"
  file: "~/.ai-chat-cli/logs/app.log"
```

## 📋 命令参考

```bash
# 查看帮助
./ai-chat-cli --help

# 版本信息
./ai-chat-cli version

# 配置管理
./ai-chat-cli config init              # 初始化配置
./ai-chat-cli config show              # 显示当前配置
./ai-chat-cli config set key value     # 设置配置项

# 对话功能
./ai-chat-cli chat [问题]              # 直接对话
./ai-chat-cli chat --provider name     # 指定提供商
./ai-chat-cli chat                     # 交互模式
```

## 🎯 支持的AI提供商

- **OpenAI** - 官方API (GPT-3.5, GPT-4等)
- **第三方兼容API** - 支持各种OpenAI兼容的服务
- **自定义提供商** - 可在配置文件中添加任意兼容的API

## 📦 项目结构

```
ai-chat-cli/
├── cmd/                    # CLI命令定义
│   ├── root.go            # 根命令
│   ├── config.go          # 配置命令
│   ├── simple_chat.go     # 对话命令
│   └── version.go         # 版本命令
├── internal/
│   ├── config/            # 配置管理
│   └── providers/         # AI提供商接口
├── configs/               # 配置文件模板
├── main.go               # 程序入口
└── README.md
```

## 🛠️ 开发

```bash
# 安装依赖
go mod tidy

# 运行测试
go test ./...

# 编译
go build -o ai-chat-cli

# 跨平台编译
GOOS=linux GOARCH=amd64 go build -o ai-chat-cli-linux
GOOS=windows GOARCH=amd64 go build -o ai-chat-cli.exe
```

## 📝 更新日志

### v1.0.0

- ✅ 基础CLI框架
- ✅ 配置管理系统
- ✅ OpenAI API集成
- ✅ 直接对话模式
- ✅ 交互模式
- ✅ 多提供商支持
- ✅ 上下文记忆
- ✅ 输入处理优化

## 🤝 贡献

欢迎贡献代码、报告问题或提出建议！

1. Fork 这个仓库
2. 创建您的功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [Cobra](https://github.com/spf13/cobra) - 强大的CLI框架
- [Viper](https://github.com/spf13/viper) - 配置管理
- [OpenAI](https://openai.com/) - AI API服务

---

如果这个项目对您有帮助，请给个 ⭐️ Star！
