# 语言包系统 (Language Pack System)

## 概述 / Overview

SSHAI项目采用基于YAML文件的外部语言包系统，支持动态加载多语言内容。语言包文件独立存储，便于维护和扩展，同时支持与二进制文件一起打包部署。

SSHAI project uses an external language pack system based on YAML files, supporting dynamic loading of multilingual content. Language pack files are stored independently for easy maintenance and expansion, while supporting packaging and deployment with binary files.

## 系统架构 / System Architecture

### 目录结构 / Directory Structure
```
sshai/
├── lang/                    # 语言包目录 / Language pack directory
│   ├── lang-zh-cn.yaml     # 简体中文语言包 / Simplified Chinese
│   └── lang-en-us.yaml     # 英文语言包 / English
├── pkg/i18n/               # i18n核心模块 / i18n core module
│   └── i18n.go            # 语言包加载器 / Language pack loader
├── config.yaml             # 主配置文件 / Main config file
└── sshai                   # 主程序 / Main program
```

### 语言包文件格式 / Language Pack File Format

语言包使用YAML格式，采用分层结构组织翻译内容：

```yaml
# lang/lang-zh-cn.yaml
server:
  starting: "正在启动SSH服务器，端口: %s"
  started: "SSH服务器已启动，监听端口: %s"

model:
  loading: "正在加载模型列表..."
  selected: "已选择模型: %s"

ai:
  thinking: "正在思考..."
  response: "💬 回答:"
```

## 核心功能 / Core Features

### 1. 动态语言加载 / Dynamic Language Loading

系统启动时根据配置文件自动加载对应语言包：

```go
// 初始化i18n系统
language := i18n.Language(cfg.I18n.Language)
if err := i18n.Init(language); err != nil {
    log.Fatal("Failed to load language pack:", err)
}
```

### 2. 结构化消息管理 / Structured Message Management

使用Go结构体映射YAML结构，提供类型安全的访问：

```go
type LanguageMessages struct {
    Server struct {
        Starting string `yaml:"starting"`
        Started  string `yaml:"started"`
    } `yaml:"server"`
    
    Model struct {
        Loading  string `yaml:"loading"`
        Selected string `yaml:"selected"`
    } `yaml:"model"`
}
```

### 3. 扁平化键值访问 / Flattened Key-Value Access

支持点分隔符的键名访问方式：

```go
// 使用方式
message := i18n.T("server.starting", port)
message := i18n.T("model.selected", modelName)
```

### 4. 语言回退机制 / Language Fallback Mechanism

- 优先使用当前设置的语言
- 如果翻译缺失，自动回退到中文
- 如果中文也缺失，返回原始键名

### 5. 并发安全 / Concurrent Safety

使用读写锁保护语言包数据，支持多个SSH连接同时访问：

```go
type I18n struct {
    currentLang  Language
    messages     map[Language]*LanguageMessages
    flatMessages map[Language]map[string]string
    mutex        sync.RWMutex
}
```

## 支持的语言 / Supported Languages

### 当前支持 / Currently Supported
- **zh-cn**: 简体中文 / Simplified Chinese
- **en-us**: 英文 / English

### 语言代码规范 / Language Code Convention
- 使用小写字母和连字符
- 格式：`语言-地区` (language-region)
- 示例：`zh-cn`, `en-us`, `ja-jp`, `ko-kr`

## 配置方法 / Configuration

### 在config.yaml中设置语言 / Set Language in config.yaml

```yaml
# 国际化配置
i18n:
  language: "zh-cn"  # 支持: zh-cn, en-us
```

### 语言包文件命名规范 / Language Pack File Naming Convention

```
lang-{语言代码}.yaml
lang-{language-code}.yaml

例如 / Examples:
- lang-zh-cn.yaml  # 简体中文
- lang-en-us.yaml  # 英文
- lang-ja-jp.yaml  # 日文
```

## 翻译内容分类 / Translation Content Categories

### 1. 服务器相关 (server)
- 启动消息 / Startup messages
- 状态信息 / Status information
- 关闭消息 / Shutdown messages

### 2. 连接相关 (connection)
- 新连接提示 / New connection notifications
- 认证状态 / Authentication status
- 连接管理 / Connection management

### 3. 认证相关 (auth)
- 密码提示 / Password prompts
- 登录状态 / Login status
- 认证结果 / Authentication results

### 4. 模型相关 (model)
- 加载状态 / Loading status
- 选择界面 / Selection interface
- 缓存信息 / Cache information

### 5. AI对话相关 (ai)
- 思考过程 / Thinking process
- 回答标识 / Response indicators
- 错误处理 / Error handling

### 6. 用户交互 (user)
- 欢迎消息 / Welcome messages
- 帮助信息 / Help information
- 操作提示 / Operation prompts

### 7. 命令相关 (cmd)
- 命令名称 / Command names
- 帮助文本 / Help text
- 错误提示 / Error messages

### 8. 错误消息 (error)
- 系统错误 / System errors
- 网络错误 / Network errors
- 配置错误 / Configuration errors

### 9. 系统信息 (system)
- 版本信息 / Version information
- 项目链接 / Project links
- 配置说明 / Configuration descriptions

## 开发指南 / Development Guide

### 添加新的翻译键 / Adding New Translation Keys

1. **在语言包文件中添加翻译**
```yaml
# lang/lang-zh-cn.yaml
new_feature:
  welcome: "欢迎使用新功能"
  help: "这是帮助信息"

# lang/lang-en-us.yaml
new_feature:
  welcome: "Welcome to new feature"
  help: "This is help information"
```

2. **更新Go结构体定义**
```go
// pkg/i18n/i18n.go
type LanguageMessages struct {
    // ... 现有字段
    NewFeature struct {
        Welcome string `yaml:"welcome"`
        Help    string `yaml:"help"`
    } `yaml:"new_feature"`
}
```

3. **在扁平化函数中添加映射**
```go
// flattenMessages函数中添加
flat["new_feature.welcome"] = messages.NewFeature.Welcome
flat["new_feature.help"] = messages.NewFeature.Help
```

4. **在代码中使用**
```go
message := i18n.T("new_feature.welcome")
help := i18n.T("new_feature.help")
```

### 添加新语言支持 / Adding New Language Support

1. **创建语言包文件**
```bash
# 创建日文语言包
cp lang/lang-zh-cn.yaml lang/lang-ja-jp.yaml
# 然后翻译内容
```

2. **添加语言常量**
```go
// pkg/i18n/i18n.go
const (
    LanguageZhCN Language = "zh-cn"
    LanguageEnUS Language = "en-us"
    LanguageJaJP Language = "ja-jp"  // 新增
)
```

3. **更新配置文档**
```yaml
# config.yaml
i18n:
  language: "ja-jp"  # 新增支持
```

## 构建和部署 / Build and Deployment

### 使用构建脚本 / Using Build Script

```bash
# 构建包含语言包的完整项目
./scripts/build_with_lang.sh
```

构建脚本会：
- 编译二进制文件
- 复制所有语言包文件
- 复制配置文件和文档
- 创建启动脚本
- 生成完整的部署包

### 手动构建 / Manual Build

```bash
# 编译程序
go build -o sshai cmd/main.go

# 确保语言包目录存在
mkdir -p lang
cp lang-*.yaml lang/

# 运行程序
./sshai
```

### 部署包结构 / Deployment Package Structure

```
sshai-deployment/
├── sshai              # 主程序
├── start.sh           # 启动脚本
├── config.yaml        # 配置文件
├── lang/              # 语言包目录
│   ├── lang-zh-cn.yaml
│   └── lang-en-us.yaml
├── keys/              # SSH密钥
└── docs/              # 文档
```

## 测试验证 / Testing and Verification

### 自动化测试 / Automated Testing

```bash
# 运行多语言测试
./scripts/test_i18n.sh
```

### 手动测试步骤 / Manual Testing Steps

1. **验证语言包文件**
```bash
# 检查文件是否存在
ls -la lang/
# 验证YAML格式
yaml-lint lang/lang-zh-cn.yaml
```

2. **测试中文界面**
```bash
# 设置中文
sed -i 's/language: .*/language: "zh-cn"/' config.yaml
./sshai
```

3. **测试英文界面**
```bash
# 设置英文
sed -i 's/language: .*/language: "en-us"/' config.yaml
./sshai
```

### 验证要点 / Verification Points

- ✅ 语言包文件正确加载
- ✅ 翻译内容正确显示
- ✅ 参数化翻译正常工作
- ✅ 语言回退机制生效
- ✅ 并发访问安全

## 性能优化 / Performance Optimization

### 1. 启动时加载 / Load at Startup
- 所有语言包在程序启动时加载到内存
- 避免运行时的文件I/O操作
- 提供快速的翻译查找

### 2. 内存管理 / Memory Management
- 使用结构化存储减少内存占用
- 扁平化映射提供O(1)查找性能
- 读写锁最小化锁竞争

### 3. 缓存策略 / Caching Strategy
- 翻译结果缓存在内存中
- 支持热重载（开发时使用）
- 生产环境建议重启更新

## 故障排除 / Troubleshooting

### 常见问题 / Common Issues

#### 1. 语言包文件未找到
```
Error: language pack file not found: lang/lang-zh-cn.yaml
```
**解决方案**：
- 检查lang目录是否存在
- 确认语言包文件名格式正确
- 验证文件权限

#### 2. YAML格式错误
```
Error: failed to parse language pack file: yaml: line 10: mapping values are not allowed in this context
```
**解决方案**：
- 使用YAML验证工具检查格式
- 注意缩进和冒号后的空格
- 检查特殊字符是否需要引号

#### 3. 翻译键未找到
```
# 显示原始键名而不是翻译内容
server.starting
```
**解决方案**：
- 检查语言包文件中是否包含该键
- 验证键名拼写是否正确
- 确认扁平化映射是否正确

#### 4. 参数化翻译错误
```
Error: wrong number of arguments for format string
```
**解决方案**：
- 检查翻译字符串中的%s、%d等占位符数量
- 确保调用T()函数时参数数量匹配
- 验证参数类型是否正确

## 最佳实践 / Best Practices

### 1. 翻译键命名 / Translation Key Naming
- 使用层次化命名：`模块.功能.具体内容`
- 保持键名简洁明了
- 使用英文和下划线

### 2. 翻译内容编写 / Translation Content Writing
- 保持翻译准确性和一致性
- 考虑上下文和用户体验
- 使用合适的标点符号和格式

### 3. 参数化设计 / Parameterization Design
- 合理使用参数化翻译
- 避免过度复杂的格式字符串
- 考虑不同语言的语序差异

### 4. 版本管理 / Version Management
- 语言包文件纳入版本控制
- 翻译更新时同步更新所有语言
- 保持向后兼容性

## 扩展计划 / Extension Plans

### 1. 更多语言支持 / More Language Support
- 日文 (ja-jp)
- 韩文 (ko-kr)
- 法文 (fr-fr)
- 德文 (de-de)

### 2. 高级功能 / Advanced Features
- 复数形式处理
- 日期时间格式化
- 数字格式化
- 文本方向支持（RTL）

### 3. 工具支持 / Tool Support
- 翻译管理工具
- 自动化翻译验证
- 翻译覆盖率检查
- 热重载开发模式

## 贡献指南 / Contribution Guidelines

### 翻译贡献 / Translation Contributions
1. Fork项目仓库
2. 添加或更新语言包文件
3. 测试翻译效果
4. 提交Pull Request

### 代码贡献 / Code Contributions
1. 遵循现有代码风格
2. 添加适当的测试用例
3. 更新相关文档
4. 确保向后兼容性

### 质量要求 / Quality Requirements
- 翻译准确性：内容准确无误
- 一致性：术语使用统一
- 完整性：覆盖所有功能模块
- 可维护性：结构清晰易扩展