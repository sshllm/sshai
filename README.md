# SSHAI - SSH AI Assistant

[English](./README_EN.md) | 简体中文

一个通过SSH连接提供AI模型服务的智能助手程序，让你可以在任何支持SSH的环境中使用AI助手。 

支持[三种调用模式](https://mp.weixin.qq.com/s/_sSEC15WOfeF0t8AaQ6Qbg)：
- **交互模式** - 通过SSH连接后，直接在终端中输入命令即可调用AI助手（`ssh your-bot@sshllm.top`）    
- **命令行模式** - 通过SSH连接后，直接在终端中执行命令即可调用AI助手（`ssh bash@sshllm.top 查看进程占用`）    
- **管道模式** - 通过SSH连接后，通过管道将内容输入到AI助手即可调用AI助手（`cat doc.txt | ssh fanyi@sshllm.top`）

此项目采用`CodeBuddy`进行开发，完全不写一行代码。    
关于开发的经验心得，请参考：[不写一行代码！我用 AI 打造了一款 AI 客户端！（开源）](https://mp.weixin.qq.com/s/-5GC3TDAP_CXAcAkGO7tMQ)    


## 🚀 精选案例
`SSHLLM`，基于当前开源版深度定制的多用户多配置版，支持用户注册、配置助手，并分享公开或者私有使用。随时随地通过SSH即可调用AI助手完成如自动生成bash脚本、代码、识别图片验证码等功能。

官网：[https://sshllm.top](https://sshllm.top)


## 🚀 体验（开源版本）
打开你的终端，输入如下命令即可立即体验在线AI服务！
```bash
ssh test.sshai.top -p 9527
```

![](docs/screenshot.png)

## ✨ 主要特性

- 🔐 **SSH安全连接** - 通过SSH协议提供加密的AI服务访问
- 🔑 **灵活认证** - 支持密码认证、SSH公钥免密登录和无密码模式
- 🗝️ **SSH Keys支持** - 支持多个SSH公钥免密登录，兼容RSA、Ed25519等密钥类型
- 🤖 **多模型支持** - 支持DeepSeek、Hunyuan等多种AI模型
- 💭 **实时思考显示** - 支持DeepSeek R1等模型的思考过程实时展示
- 🎨 **美观界面** - 彩色输出、动画效果和ASCII艺术
- ⚙️ **灵活配置** - 支持动态指定配置文件（-c参数）和完整的YAML配置
- 🌐 **多语言支持** - 支持中文和英文界面
- 📝 **自定义提示词** - 可配置的AI提示词系统
- 🚀 **启动欢迎页** - 程序启动时显示美观的欢迎banner
- 🏗️ **模块化设计** - 清晰的代码架构，易于扩展

## 🚀 快速开始

### 1. 下载和编译

```bash
# 克隆项目
git clone https://github.com/sshllm/sshai.git
cd sshai

# 编译程序
make build
# 或者
go build -o sshai cmd/main.go
```

### 2. 配置设置

编辑 `config.yaml` 文件，设置你的API密钥：

```yaml
# API配置
api:
  base_url: "https://api.deepseek.com/v1"
  api_key: "your-api-key-here"
  default_model: "deepseek-v3"

# 服务器配置
server:
  port: 2213
  welcome_message: "欢迎使用SSHAI！"

# 认证配置（可选）
auth:
  password: ""  # 留空=无密码认证
  login_prompt: "请输入访问密码: "
  # SSH公钥免密登录配置（仅在设置password时生效）
  authorized_keys:
    - "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC... user@hostname"
    - "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI... user2@hostname"
  authorized_keys_file: "~/.ssh/authorized_keys"  # 可选：从文件读取公钥

# 自定义提示词配置
prompt:
  system_prompt: "你是一个专业的AI助手，请用中文回答问题。"
  stdin_prompt: "请分析以下内容并提供相关的帮助或建议："
  exec_prompt: "请回答以下问题或执行以下任务："
```

### 3. 运行服务器

```bash
# 直接运行（使用默认配置文件 config.yaml）
./sshai

# 指定配置文件运行
./sshai -c config.yaml
./sshai -c /path/to/your/config.yaml

# 后台运行
./sshai > server.log 2>&1 &

# 使用脚本运行
./scripts/run.sh
```

#### 命令行参数

- `-c <config_file>` - 指定配置文件路径
  - 如果不指定，默认使用当前目录下的 `config.yaml`
  - 如果配置文件不存在，程序会显示错误提示并退出

```bash
# 使用示例
./sshai -c config.yaml          # 使用当前目录的配置文件
./sshai -c /etc/sshai/config.yaml  # 使用绝对路径的配置文件
./sshai                         # 默认使用 config.yaml
```

### 4. 连接使用

```bash
# 交互模式
ssh user@localhost -p 2213

# 直接执行命令
ssh user@localhost -p 2213 "你好，请介绍一下你自己"

# 管道输入分析
cat file.txt | ssh user@localhost -p 2213
echo "分析这段代码" | ssh user@localhost -p 2213
```

## 📁 项目结构

```
sshai/
├── README.md              # 中文说明文档
├── README_EN.md           # 英文说明文档
├── LICENSE                # 开源协议
├── config.yaml           # 主配置文件
├── config-en.yaml        # 英文配置文件
├── go.mod                # Go模块依赖
├── Makefile              # 构建脚本
├── cmd/                  # 程序入口
│   └── main.go           # 主程序文件
├── pkg/                  # 核心模块
│   ├── config/           # 配置管理
│   ├── models/           # 数据模型
│   ├── ai/               # AI助手功能
│   ├── ssh/              # SSH服务器
│   └── utils/            # 工具函数
├── docs/                 # 项目文档
├── scripts/              # 测试和运行脚本
└── keys/                 # SSH密钥文件
```

## 🔧 配置指南

### API配置

支持多个API端点配置：

```yaml
api:
  base_url: "https://api.deepseek.com/v1"
  api_key: "your-deepseek-key"
  default_model: "deepseek-v3"
  timeout: 600

# 可配置多个API
api_endpoints:
  - name: "deepseek"
    base_url: "https://api.deepseek.com/v1"
    api_key: "your-key"
    default_model: "deepseek-v3"
  - name: "local"
    base_url: "http://localhost:11434/v1"
    api_key: "ollama"
    default_model: "gemma2:27b"
```

### 认证配置

#### 密码认证
```yaml
auth:
  password: "your-secure-password"  # 设置访问密码
  login_prompt: "请输入访问密码: "
```

#### SSH公钥免密登录
```yaml
auth:
  password: "your-secure-password"  # 必须设置密码才能启用SSH公钥认证
  login_prompt: "请输入访问密码: "
  # 方式一：直接配置公钥列表
  authorized_keys:
    - "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC... user@hostname"
    - "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI... user2@hostname"
  # 方式二：从文件读取公钥
  authorized_keys_file: "~/.ssh/authorized_keys"
```

**SSH公钥使用方法**：
```bash
# 生成SSH密钥对
ssh-keygen -t ed25519 -f ~/.ssh/sshai_key

# 使用私钥连接（免密登录）
ssh -i ~/.ssh/sshai_key -p 2213 user@localhost

# 查看公钥内容（用于配置）
cat ~/.ssh/sshai_key.pub
```

**注意**: 
- SSH公钥认证仅在设置密码时启用，提供额外的安全保障
- 支持多个公钥同时配置，兼容RSA、Ed25519、ECDSA等密钥类型
- 登录成功后会自动显示程序内置的欢迎信息，无需在配置文件中设置

### 提示词配置

```yaml
prompt:
  system_prompt: "你是一个专业的AI助手..."
  stdin_prompt: "请分析以下内容："
  exec_prompt: "请回答以下问题："
```

## 🧪 测试

项目包含完整的测试脚本：

```bash
# 基础功能测试
./scripts/test.sh

# SSH执行功能测试
./scripts/test_ssh_exec_final.sh

# 标准输入功能测试
./scripts/test_stdin_feature.sh

# 认证功能测试
./scripts/test_auth.sh

# DeepSeek R1思考模式测试
./scripts/test_deepseek_r1.sh

# SSH Keys免密登录功能测试
./scripts/test_ssh_keys.sh
```

## 📚 文档

- [配置指南](docs/CONFIG_GUIDE.md) - 详细的配置说明
- [使用指南](docs/USAGE.md) - 功能介绍和使用方法
- [架构说明](docs/MODULAR_ARCHITECTURE.md) - 模块化架构设计
- [认证配置](docs/AUTH_CONFIG_EXAMPLE.md) - SSH认证配置示例
- [SSH Keys指南](docs/SSH_KEYS_GUIDE.md) - SSH公钥免密登录配置指南

## 🤝 贡献

欢迎提交Issue和Pull Request！

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开Pull Request

## 📄 许可证

本项目采用 Apache License 2.0 开源许可证。详情请查看 [LICENSE](LICENSE) 文件。

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者和用户！

---

**注意**: 本项目遵循 Apache 2.0 开源协议，欢迎个人和商业使用。