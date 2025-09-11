# SSH AI - 基于SSH协议的AI助手服务

## 项目简介

SSH AI是一个创新的SSH服务器，用户通过标准SSH客户端连接后，可以直接与DeepSeek-V3大模型进行实时对话。无需安装额外软件，只需一个SSH命令即可享受AI助手服务。

## 快速开始

### 1. 启动服务器
```bash
# 方式一：使用启动脚本
./run.sh

# 方式二：直接运行可执行文件
./sshai

# 方式三：从源码运行
go run main.go
```

### 2. 连接服务器
```bash
ssh gpt-5@localhost -p 2212
```

### 3. 开始对话
```
Hello!
gpt-5@sshai> 你好，请介绍一下自己
我是DeepSeek开发的AI助手...

gpt-5@sshai> /new
[新会话已创建]

gpt-5@sshai> exit
Goodbye!
```

## 核心特性

### 🔐 无密码登录
- 任何用户名都可以直接连接，无需密码验证
- 完全兼容标准SSH协议

### 🤖 真实AI集成
- 集成DeepSeek-V3大模型
- 支持流式响应，实时显示AI回复
- OpenAI兼容API接口

### 💬 智能会话管理
- 自动维护对话上下文
- 支持多轮对话记忆
- `/new` 命令创建新会话

### 🚀 高性能架构
- 支持多用户并发连接
- 每个连接独立的会话上下文
- 优雅的错误处理和恢复

## 命令说明

| 命令 | 功能 |
|------|------|
| 普通文本 | 与AI进行对话 |
| `/new` | 清除当前会话，开始新对话 |
| `exit` | 退出SSH连接 |

## 技术架构

### 后端技术栈
- **语言**: Golang 1.22.2
- **SSH库**: golang.org/x/crypto/ssh
- **HTTP客户端**: 标准库 net/http
- **JSON处理**: 标准库 encoding/json

### AI服务集成
- **API地址**: https://ds.openugc.com/v1
- **模型**: DeepSeek-V3
- **协议**: OpenAI兼容接口
- **响应方式**: Server-Sent Events流式响应

### 系统架构图
```
SSH客户端 ←→ SSH服务器 ←→ AI API服务
    ↓           ↓           ↓
  用户输入   会话管理    模型推理
            上下文      流式响应
```

## 安装部署

### 环境要求
- Go 1.22.2 或更高版本
- 网络连接（访问AI API）
- 开放端口2212

### 编译安装
```bash
# 克隆项目
git clone <repository-url>
cd sshai

# 安装依赖
go mod tidy

# 编译程序
go build -o sshai main.go

# 启动服务
./sshai
```

### Docker部署（可选）
```dockerfile
FROM golang:1.22.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o sshai main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/sshai .
EXPOSE 2212
CMD ["./sshai"]
```

## 使用示例

### 基础对话
```
gpt-5@sshai> 写一个Python的Hello World程序
当然！这是一个简单的Python Hello World程序：

```python
print("Hello, World!")
```

这是最基础的Python程序，它会在控制台输出"Hello, World!"。
```

### 上下文对话
```
gpt-5@sshai> 我刚才问了什么？
您刚才问我写一个Python的Hello World程序。

gpt-5@sshai> 能解释一下这个程序吗？
当然！刚才的Python程序很简单：
- `print()` 是Python的内置函数，用于输出内容
- "Hello, World!" 是要输出的字符串
- 运行这个程序会在终端显示：Hello, World!
```

### 新会话管理
```
gpt-5@sshai> /new
[新会话已创建]

gpt-5@sshai> 我刚才问了什么？
抱歉，我没有看到您之前的问题记录，因为这是一个新的会话。
```

## 测试验证

### 功能测试
```bash
# 运行完整测试
./test_ai.sh

# 基础功能测试
./test.sh
```

### 连接测试
```bash
# 测试SSH连接
ssh gpt-5@localhost -p 2212

# 测试不同用户名
ssh test@localhost -p 2212
ssh admin@localhost -p 2212
```

## 故障排除

### 常见问题

**Q: 连接被拒绝**
```bash
# 检查服务器是否运行
ps aux | grep sshai

# 检查端口是否被占用
netstat -tlnp | grep 2212
```

**Q: AI响应错误**
```bash
# 检查网络连接
curl -I https://ds.openugc.com/v1

# 查看服务器日志
./sshai 2>&1 | tee server.log
```

**Q: 首次连接主机密钥警告**
```bash
# 接受主机密钥
ssh -o StrictHostKeyChecking=no gpt-5@localhost -p 2212
```

## 安全说明

### 当前安全级别
- ⚠️ 无密码验证（适用于内网测试）
- ⚠️ 临时生成主机密钥
- ⚠️ 明文API密钥

### 生产环境建议
- 添加用户认证机制
- 使用持久化主机密钥
- 配置环境变量存储API密钥
- 启用TLS加密
- 添加访问日志和监控

## 扩展开发

### 添加新命令
```go
// 在ProcessMessage函数中添加
if strings.HasPrefix(input, "/help") {
    channel.Write([]byte("\r\n可用命令：\r\n/new - 新会话\r\nexit - 退出\r\n"))
    return
}
```

### 集成其他AI模型
```go
// 修改API配置
const (
    APIBaseURL = "https://api.openai.com/v1"
    ModelName  = "gpt-4"
)
```

### 添加用户认证
```go
// 在SSH配置中添加
PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
    if validateUser(conn.User(), string(password)) {
        return nil, nil
    }
    return nil, fmt.Errorf("password rejected")
}
```

## 贡献指南

1. Fork项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建Pull Request

## 许可证

本项目采用MIT许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 联系方式

- 项目地址: [GitHub Repository]
- 问题反馈: [Issues]
- 文档: [Wiki]

---

**享受与AI的SSH对话体验！** 🚀