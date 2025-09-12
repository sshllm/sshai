## SSH AI项目
创建一个SSH服务，提供用户连接之后立即使用AI大模型服务。

用户：ssh gpt-5@test.com 无需密码
自动登录，进入自定义交互式终端：
```sshai
Hello!
gpt-5@sshai> 
gpt-5@sshai> Who are you?
I'm your AI assistant.
gpt-5@sshai> 
```

## 功能特性

- ✅ 标准SSH协议支持，无需密码登录
- ✅ 智能模型选择：根据用户名自动匹配或手动选择模型
- ✅ 流式AI响应，支持实时对话
- ✅ 深度思考模型支持，显示思考过程
- ✅ 中文输入输出完美支持
- ✅ 上下文管理，支持 `/new` 命令清空对话历史
- ✅ 可配置化设计，支持动态修改配置
- ✅ 加载动画和思考动画，提升用户体验
- ✅ Ctrl+C 中断支持

## 快速开始

### 1. 编译程序

#### 模块化版本（推荐）
```bash
# 使用 Makefile
make build

# 或直接使用 go build
go build -o sshai cmd/main.go
```


### 2. 配置文件
程序首次运行前，请确保 `config.yaml` 配置文件存在。默认配置：

```yaml
# SSH AI 服务器配置文件
server:
  port: "2212"
  welcome_message: "Hello!"
  prompt_template: "%s@sshai> "

api:
  base_url: "https://ds.openugc.com/v1"
  api_key: ""
  default_model: "deepseek-v3"
  timeout: 30

display:
  line_width: 80
  thinking_animation_interval: 150
  loading_animation_interval: 100

security:
  host_key_file: "host_key.pem"
```

### 3. 运行服务

#### 模块化版本
```bash
# 使用 Makefile
make run

# 或直接运行
./sshai
```


### 4. 连接使用
```bash
# 使用用户名连接（会自动匹配相应模型）
ssh deepseek@localhost -p 2212
ssh gpt@localhost -p 2212
ssh claude@localhost -p 2212

# 或直接连接（手动选择模型）
ssh localhost -p 2212
```

## 架构说明

本项目采用模块化架构设计，提供更好的可维护性和可扩展性：

- **模块化架构**: 详细说明请参考 [MODULAR_ARCHITECTURE.md](MODULAR_ARCHITECTURE.md)
- **重构总结**: 重构过程和成果请参考 [MODULAR_REFACTOR_SUMMARY.md](MODULAR_REFACTOR_SUMMARY.md)

## 配置说明

详细的配置说明请参考 [CONFIG_GUIDE.md](CONFIG_GUIDE.md)

主要配置项：
- **server.port**: SSH服务端口
- **api.base_url**: AI API地址
- **api.api_key**: API密钥
- **api.default_model**: 默认模型
- **display.line_width**: 终端显示宽度

## 使用技巧

### 模型选择
- 使用包含模型名的用户名连接，系统会自动匹配相应模型
- 如果找到多个匹配模型，会提供选择界面
- 如果没有匹配模型，会显示所有可用模型供选择

### 对话管理
- 输入 `/new` 清空对话历史，开始新的对话
- 使用 Ctrl+C 可以中断AI响应
- 支持中文输入，包括删除和编辑

### 深度思考模型
- 支持 deepseek-r1 等深度思考模型
- 会显示模型的思考过程和最终回答
- 思考阶段有专门的动画提示

## 开发技术栈
1. Golang go1.22.2
2. golang.org/x/crypto/ssh - SSH协议实现
3. gopkg.in/yaml.v2 - YAML配置文件解析

## 其他要求
1. 符合标准SSH协议
2. 用户无需输入密码即可登录成功
3. 支持配置文件动态修改