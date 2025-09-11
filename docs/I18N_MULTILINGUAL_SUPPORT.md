# 多语言国际化支持 (i18n Multilingual Support)

## 功能概述 / Overview

SSHAI项目现已支持多语言国际化，用户可以通过配置文件选择界面语言。目前支持简体中文和英文两种语言。

SSHAI project now supports multilingual internationalization. Users can select the interface language through the configuration file. Currently supports Simplified Chinese and English.

## 支持的语言 / Supported Languages

- **简体中文 (zh-cn)** - 默认语言 / Default language
- **English (en-us)** - 英文支持 / English support

## 配置方法 / Configuration

### 在config.yaml中设置语言 / Set Language in config.yaml

```yaml
# 国际化配置 / Internationalization Configuration
i18n:
  language: "zh-cn"  # 支持的语言: zh-cn (简体中文), en-us (英文)
                     # Supported languages: zh-cn (Simplified Chinese), en-us (English)
```

### 语言选项 / Language Options

- `zh-cn`: 简体中文界面 / Simplified Chinese interface
- `en-us`: 英文界面 / English interface

## 技术实现 / Technical Implementation

### 架构设计 / Architecture Design

```
pkg/i18n/
├── i18n.go          # 核心i18n管理器 / Core i18n manager
└── messages         # 消息定义 / Message definitions
```

### 核心组件 / Core Components

#### 1. I18n管理器 / I18n Manager
```go
type I18n struct {
    currentLang Language
    messages    map[Language]map[string]string
    mutex       sync.RWMutex
}
```

#### 2. 翻译函数 / Translation Function
```go
func T(key string, args ...interface{}) string
```

#### 3. 语言类型 / Language Types
```go
const (
    LanguageZhCN Language = "zh-cn"  // 简体中文
    LanguageEnUS Language = "en-us"  // 英文
)
```

### 使用方法 / Usage

#### 在代码中使用翻译 / Using Translation in Code
```go
import "sshai/pkg/i18n"

// 简单翻译 / Simple translation
message := i18n.T("user.welcome")

// 带参数翻译 / Translation with parameters
message := i18n.T("model.selected", modelName)
```

#### 初始化i18n系统 / Initialize i18n System
```go
// 在main.go中 / In main.go
cfg := config.Get()
language := i18n.Language(cfg.I18n.Language)
i18n.Init(language)
```

## 翻译覆盖范围 / Translation Coverage

### 服务器相关 / Server Related
- 服务器启动消息 / Server startup messages
- 连接状态信息 / Connection status information
- 错误消息 / Error messages

### 认证相关 / Authentication Related
- 密码提示 / Password prompts
- 登录成功/失败消息 / Login success/failure messages

### 模型相关 / Model Related
- 模型加载状态 / Model loading status
- 模型选择界面 / Model selection interface
- 缓存状态信息 / Cache status information

### AI对话相关 / AI Conversation Related
- 思考过程提示 / Thinking process prompts
- 回答标识 / Response indicators
- 中断消息 / Interruption messages

### 用户交互 / User Interaction
- 欢迎消息 / Welcome messages
- 帮助信息 / Help information
- 命令提示 / Command prompts

## 配置示例 / Configuration Examples

### 中文配置示例 / Chinese Configuration Example
```yaml
# config.yaml
server:
  welcome_message: "Hello!欢迎使用SSHAI！"
  
auth:
  login_prompt: "请输入访问密码: "
  login_success_message: |
    🎉 欢迎使用 SSHAI v1.0 - SSH AI Assistant
    📖 项目地址: https://github.com/your-repo/sshai

i18n:
  language: "zh-cn"
```

### 英文配置示例 / English Configuration Example
```yaml
# config-en.yaml
server:
  welcome_message: "Hello! Welcome to SSHAI!"
  
auth:
  login_prompt: "Enter access password: "
  login_success_message: |
    🎉 Welcome to SSHAI v1.0 - SSH AI Assistant
    📖 Project URL: https://github.com/your-repo/sshai

i18n:
  language: "en-us"
```

## 测试验证 / Testing and Verification

### 自动化测试 / Automated Testing
```bash
# 运行多语言测试脚本 / Run multilingual test script
./scripts/test_i18n.sh
```

### 手动测试步骤 / Manual Testing Steps

#### 1. 测试中文界面 / Test Chinese Interface
```bash
# 设置中文配置 / Set Chinese configuration
echo 'i18n:\n  language: "zh-cn"' >> config.yaml

# 启动服务器 / Start server
./sshai

# 连接测试 / Connect and test
ssh -p 2212 test@localhost
```

#### 2. 测试英文界面 / Test English Interface
```bash
# 设置英文配置 / Set English configuration
sed -i 's/language: "zh-cn"/language: "en-us"/' config.yaml

# 重启服务器 / Restart server
./sshai

# 连接测试 / Connect and test
ssh -p 2212 test@localhost
```

### 验证要点 / Verification Points

- ✅ 服务器启动消息语言正确 / Server startup messages in correct language
- ✅ 模型选择界面语言正确 / Model selection interface in correct language
- ✅ AI对话提示语言正确 / AI conversation prompts in correct language
- ✅ 错误消息语言正确 / Error messages in correct language
- ✅ 用户交互语言正确 / User interaction in correct language

## 扩展新语言 / Adding New Languages

### 1. 添加语言常量 / Add Language Constant
```go
// 在pkg/i18n/i18n.go中 / In pkg/i18n/i18n.go
const (
    LanguageZhCN Language = "zh-cn"
    LanguageEnUS Language = "en-us"
    LanguageJaJP Language = "ja-jp"  // 新增日文 / Add Japanese
)
```

### 2. 添加翻译消息 / Add Translation Messages
```go
// 在loadMessages()函数中添加 / Add in loadMessages() function
i.messages[LanguageJaJP] = map[string]string{
    "user.welcome": "ようこそ",
    "model.loading": "モデルを読み込んでいます...",
    // ... 更多翻译 / More translations
}
```

### 3. 更新配置文档 / Update Configuration Documentation
```yaml
i18n:
  language: "ja-jp"  # 新增日文支持 / Add Japanese support
```

## 最佳实践 / Best Practices

### 1. 翻译键命名规范 / Translation Key Naming Convention
```
模块.功能.具体内容
module.function.specific_content

例如 / Examples:
- server.starting
- model.selected
- user.welcome
- error.connection
```

### 2. 参数化翻译 / Parameterized Translation
```go
// 好的做法 / Good practice
i18n.T("model.selected", modelName)

// 避免 / Avoid
i18n.T("model.selected") + ": " + modelName
```

### 3. 回退机制 / Fallback Mechanism
- 如果当前语言缺少翻译，自动回退到中文 / Auto fallback to Chinese if translation missing
- 如果中文也缺少翻译，返回原始key / Return original key if Chinese translation also missing

## 性能考虑 / Performance Considerations

### 1. 内存使用 / Memory Usage
- 所有翻译消息在启动时加载到内存 / All translation messages loaded into memory at startup
- 使用读写锁保证并发安全 / Use RWMutex for concurrent safety

### 2. 查找效率 / Lookup Efficiency
- O(1)时间复杂度的消息查找 / O(1) time complexity for message lookup
- 无需文件I/O操作 / No file I/O operations required

### 3. 线程安全 / Thread Safety
- 支持多个SSH连接同时使用 / Support multiple SSH connections simultaneously
- 读写锁保护共享数据 / RWMutex protects shared data

## 故障排除 / Troubleshooting

### 常见问题 / Common Issues

#### 1. 翻译不生效 / Translation Not Working
```bash
# 检查配置 / Check configuration
grep -A2 "i18n:" config.yaml

# 检查语言设置 / Check language setting
# 确保语言代码正确: zh-cn 或 en-us / Ensure correct language code
```

#### 2. 部分消息未翻译 / Some Messages Not Translated
```bash
# 检查是否有遗漏的翻译键 / Check for missing translation keys
# 查看日志输出的原始key / Check logs for original keys
```

#### 3. 编译错误 / Compilation Errors
```bash
# 确保导入了i18n包 / Ensure i18n package is imported
import "sshai/pkg/i18n"

# 检查函数调用 / Check function calls
i18n.T("key.name")
```

## 未来计划 / Future Plans

### 1. 更多语言支持 / More Language Support
- 日文 (ja-jp) / Japanese
- 韩文 (ko-kr) / Korean
- 法文 (fr-fr) / French
- 德文 (de-de) / German

### 2. 动态语言切换 / Dynamic Language Switching
- 运行时切换语言 / Runtime language switching
- 用户级别语言设置 / User-level language settings

### 3. 外部翻译文件 / External Translation Files
- JSON/YAML格式的翻译文件 / JSON/YAML translation files
- 热重载翻译更新 / Hot reload translation updates

## 贡献指南 / Contribution Guidelines

### 添加新翻译 / Adding New Translations
1. 在`pkg/i18n/i18n.go`中添加翻译键值对 / Add key-value pairs in `pkg/i18n/i18n.go`
2. 确保所有支持的语言都有对应翻译 / Ensure all supported languages have corresponding translations
3. 更新文档和测试用例 / Update documentation and test cases
4. 提交PR进行代码审查 / Submit PR for code review

### 翻译质量要求 / Translation Quality Requirements
- 准确性：翻译内容准确无误 / Accuracy: Translations are accurate and error-free
- 一致性：术语使用保持一致 / Consistency: Consistent terminology usage
- 自然性：符合目标语言表达习惯 / Naturalness: Natural expression in target language