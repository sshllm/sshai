# SSH AI 服务器配置指南

## 配置文件说明

程序现在支持通过 `config.yaml` 文件进行配置，用户可以根据需要修改各项设置。

## 配置文件结构

```yaml
# SSH AI 服务器配置文件

# 服务器配置
server:
  port: "2212"                    # SSH服务器监听端口
  welcome_message: "Hello!"       # 欢迎消息
  prompt_template: "%s@sshai> "   # 提示符模板，%s会被替换为模型名

# AI API配置
api:
  base_url: "https://ds.openugc.com/v1"  # API基础URL
  api_key: "your-api-key"                        # API密钥
  default_model: "deepseek-v3"           # 默认模型
  timeout: 30                            # 请求超时时间（秒）

# 显示配置
display:
  line_width: 80                         # 终端显示宽度
  thinking_animation_interval: 150       # 思考动画间隔（毫秒）
  loading_animation_interval: 100        # 加载动画间隔（毫秒）

# 证书配置
security:
  host_key_file: "host_key.pem"         # SSH主机密钥文件路径
```

## 配置项详细说明

### 服务器配置 (server)

- **port**: SSH服务器监听的端口号，默认为 "2212"
- **welcome_message**: 用户连接时显示的欢迎消息
- **prompt_template**: 命令提示符的模板，使用 `%s` 作为模型名的占位符

### API配置 (api)

- **base_url**: AI API的基础URL地址
- **api_key**: 访问API所需的密钥
- **default_model**: 当无法获取模型列表或用户未选择时使用的默认模型
- **timeout**: HTTP请求的超时时间，单位为秒

### 显示配置 (display)

- **line_width**: 终端显示的行宽度，用于文本换行
- **thinking_animation_interval**: 思考动画的刷新间隔，单位为毫秒
- **loading_animation_interval**: 加载动画的刷新间隔，单位为毫秒

### 安全配置 (security)

- **host_key_file**: SSH主机密钥文件的路径，程序会自动生成或加载此文件

## 使用方法

1. 确保 `config.yaml` 文件与可执行文件在同一目录
2. 根据需要修改配置文件中的各项设置
3. 运行程序：`./sshai`
4. 使用SSH客户端连接：`ssh username@localhost -p 端口号`

## 配置示例

### 更改API服务商

```yaml
api:
  base_url: "https://api.openai.com/v1"
  api_key: "your-openai-api-key"
  default_model: "gpt-4"
```

### 自定义提示符

```yaml
server:
  prompt_template: "[%s]$ "  # 显示为 [模型名]$
```

### 调整动画速度

```yaml
display:
  thinking_animation_interval: 200  # 更慢的思考动画
  loading_animation_interval: 80    # 更快的加载动画
```

## 注意事项

1. 修改配置文件后需要重启程序才能生效
2. 确保API密钥的安全性，不要将包含真实密钥的配置文件提交到版本控制系统
3. 端口号需要确保没有被其他程序占用
4. 主机密钥文件会在首次运行时自动生成，请妥善保管