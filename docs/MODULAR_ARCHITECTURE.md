# SSH AI 服务器模块化架构文档

## 架构概述

SSH AI 服务器已重构为模块化设计，将原来的单文件程序拆分为多个独立的模块，提高了代码的可维护性、可扩展性和可测试性。

## 目录结构

```
sshai/
├── cmd/                    # 应用程序入口
│   └── main.go            # 主程序文件
├── pkg/                   # 核心包目录
│   ├── config/           # 配置管理模块
│   │   └── config.go     # 配置加载和管理
│   ├── models/           # 数据模型模块
│   │   └── models.go     # 数据结构定义
│   ├── ai/               # AI助手模块
│   │   ├── assistant.go  # AI助手核心功能
│   │   └── models.go     # 模型管理功能
│   ├── ssh/              # SSH服务器模块
│   │   ├── server.go     # SSH服务器实现
│   │   └── session.go    # SSH会话处理
│   └── utils/            # 工具函数模块
│       └── text.go       # 文本处理工具
├── config.yaml           # 配置文件
├── go.mod                # Go模块定义
├── go.sum                # 依赖版本锁定
├── Makefile              # 构建脚本
├── main.go               # 原版本（兼容性保留）
└── README.md             # 项目说明
```

## 模块详细说明

### 1. 配置管理模块 (`pkg/config`)

**职责**: 负责配置文件的加载、解析和全局配置管理

**主要功能**:
- 配置文件加载和解析
- 全局配置实例管理
- 配置结构体定义

**核心类型**:
```go
type Config struct {
    Server   ServerConfig   // 服务器配置
    API      APIConfig      // API配置
    Display  DisplayConfig  // 显示配置
    Security SecurityConfig // 安全配置
}
```

**主要函数**:
- `Load(configPath string) error` - 加载配置文件
- `Get() *Config` - 获取全局配置实例

### 2. 数据模型模块 (`pkg/models`)

**职责**: 定义所有数据结构和类型

**主要功能**:
- OpenAI API 相关数据结构
- 模型信息数据结构
- 聊天消息数据结构

**核心类型**:
```go
type ChatMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type ChatRequest struct {
    Model    string        `json:"model"`
    Messages []ChatMessage `json:"messages"`
    Stream   bool          `json:"stream"`
}

type ModelInfo struct {
    ID      string `json:"id"`
    Object  string `json:"object"`
    Created int64  `json:"created"`
}
```

### 3. AI助手模块 (`pkg/ai`)

**职责**: 处理AI相关的所有功能

**主要功能**:
- AI助手实例管理
- API调用和响应处理
- 模型获取和选择
- 流式响应处理
- 思考内容显示

**核心类型**:
```go
type Assistant struct {
    messages     []models.ChatMessage
    currentModel string
    username     string
}
```

**主要函数**:
- `NewAssistant(username string) *Assistant` - 创建AI助手
- `ProcessMessage(input string, channel ssh.Channel, interrupt chan bool)` - 处理用户消息
- `GetAvailableModels() ([]ModelInfo, error)` - 获取可用模型
- `SelectModelByUsername(channel ssh.Channel, models []ModelInfo, username string) string` - 模型选择

### 4. SSH服务器模块 (`pkg/ssh`)

**职责**: 处理SSH协议相关的所有功能

**主要功能**:
- SSH服务器创建和启动
- SSH连接处理
- 会话管理
- 用户输入处理
- 主机密钥管理

**核心类型**:
```go
type Server struct {
    config *ssh.ServerConfig
}
```

**主要函数**:
- `NewServer() (*Server, error)` - 创建SSH服务器
- `Start() error` - 启动服务器
- `HandleSession(channel ssh.Channel, requests <-chan *ssh.Request, username string)` - 处理会话

### 5. 工具函数模块 (`pkg/utils`)

**职责**: 提供通用的工具函数

**主要功能**:
- 文本换行处理
- 显示宽度计算
- 断行位置查找
- UTF-8字符处理

**主要函数**:
- `WrapText(text string, width int) string` - 文本换行
- `GetDisplayWidth(text string) int` - 计算显示宽度
- `FindBreakPosition(text string, maxWidth int) int` - 查找断行位置

### 6. 主程序 (`cmd/main.go`)

**职责**: 应用程序入口点

**主要功能**:
- 配置加载
- 服务器创建和启动
- 错误处理

## 模块间依赖关系

```
cmd/main.go
    ├── pkg/config (配置管理)
    └── pkg/ssh (SSH服务器)
        ├── pkg/ai (AI助手)
        │   ├── pkg/config
        │   ├── pkg/models
        │   └── pkg/utils
        ├── pkg/config
        └── pkg/models

pkg/ai
    ├── pkg/config
    ├── pkg/models
    └── pkg/utils

pkg/utils (独立模块，无依赖)
pkg/models (独立模块，无依赖)
pkg/config (独立模块，无依赖)
```

## 构建和运行

### 使用 Makefile

```bash
# 构建模块化版本
make build

# 运行模块化版本
make run

# 构建原版本（兼容性）
make build-legacy

# 运行原版本
make run-legacy

# 清理构建文件
make clean

# 安装依赖
make deps

# 开发环境设置
make dev-setup
```

### 直接使用 Go 命令

```bash
# 构建模块化版本
go build -o sshai cmd/main.go

# 构建原版本
go build -o sshai main.go

# 运行原版本
./sshai
```

## 模块化的优势

### 1. 可维护性
- **单一职责**: 每个模块只负责特定的功能
- **清晰边界**: 模块间接口明确，职责分离
- **易于调试**: 问题定位更加精确

### 2. 可扩展性
- **插件化**: 可以轻松添加新的功能模块
- **接口抽象**: 便于实现不同的后端服务
- **配置驱动**: 通过配置文件控制行为

### 3. 可测试性
- **单元测试**: 每个模块可以独立测试
- **模拟依赖**: 便于创建测试替身
- **集成测试**: 模块间交互测试更容易

### 4. 可重用性
- **独立模块**: 工具函数可以在其他项目中重用
- **标准接口**: 遵循Go语言的最佳实践
- **包管理**: 便于版本控制和依赖管理

## 开发指南

### 添加新功能

1. **确定模块**: 根据功能确定应该添加到哪个模块
2. **定义接口**: 如果需要，先定义接口
3. **实现功能**: 在相应模块中实现功能
4. **更新文档**: 更新相关文档和注释
5. **编写测试**: 为新功能编写测试用例

### 修改现有功能

1. **定位模块**: 找到需要修改的模块
2. **理解依赖**: 了解模块间的依赖关系
3. **保持兼容**: 尽量保持接口的向后兼容性
4. **更新测试**: 更新相关的测试用例

### 最佳实践

1. **遵循Go惯例**: 使用标准的Go项目结构
2. **错误处理**: 合理处理和传播错误
3. **日志记录**: 在关键位置添加日志
4. **文档注释**: 为公开的函数和类型添加注释
5. **性能考虑**: 注意内存使用和并发安全

## 迁移指南

### 从单文件版本迁移

1. **备份原文件**: 保留原始的 `main.go` 文件
2. **使用新版本**: 使用 `cmd/main.go` 作为入口点
3. **配置文件**: 确保 `config.yaml` 文件存在
4. **测试功能**: 验证所有功能正常工作
5. **更新脚本**: 更新部署和运行脚本

### 兼容性说明

- 配置文件格式保持不变
- 用户界面和交互方式保持一致
- SSH协议兼容性完全保持
- 所有原有功能都得到保留

## 总结

模块化架构使SSH AI服务器更加健壮、可维护和可扩展。通过清晰的模块划分和接口设计，开发者可以更容易地理解、修改和扩展代码。这种架构为未来的功能增强和性能优化奠定了坚实的基础。