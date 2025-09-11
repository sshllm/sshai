# SSH执行命令功能

## 概述

SSHAI 现在支持通过SSH直接执行命令的功能，用户可以使用 `ssh user@server "command"` 的格式直接向AI发送问题并获取响应，无需进入交互式会话。

## 功能特性

### 1. 直接命令执行
- **一次性对话**: 发送问题并立即获取AI响应
- **无需交互**: 不需要进入交互式SSH会话
- **快速响应**: 适合脚本化和自动化场景

### 2. 完整的中文支持
- **UTF-8编码**: 完全支持中文问题和回答
- **字符处理**: 正确处理中英文混合内容
- **编码兼容**: 与各种SSH客户端兼容

### 3. 智能模型选择
- **自动匹配**: 根据用户名自动选择合适的AI模型
- **无界面选择**: 在后台完成模型选择，不显示选择过程
- **快速启动**: 跳过交互式会话的初始化步骤

## 使用方法

### 基本语法
```bash
ssh [username]@[server] -p [port] "[your question]"
```

### 实际示例

#### 1. 简单问候
```bash
ssh gpt@localhost -p 2213 "你好"
```

#### 2. 中文问题
```bash
ssh gpt@localhost -p 2213 "什么是人工智能？"
```

#### 3. 英文问题
```bash
ssh gpt@localhost -p 2213 "What is machine learning?"
```

#### 4. 复杂问题
```bash
ssh gpt@localhost -p 2213 "请用简单的语言解释一下深度学习的基本原理"
```

#### 5. 技术问题
```bash
ssh gpt@localhost -p 2213 "如何在Python中实现一个简单的神经网络？"
```

## 技术实现

### SSH协议处理

#### exec请求处理
```go
case "exec":
    // 处理执行命令请求
    if len(req.Payload) > 4 {
        // SSH exec请求的payload格式: [4字节长度][命令字符串]
        cmdLen := int(req.Payload[0])<<24 | int(req.Payload[1])<<16 | 
                 int(req.Payload[2])<<8 | int(req.Payload[3])
        if cmdLen > 0 && len(req.Payload) >= 4+cmdLen {
            execCommand = string(req.Payload[4 : 4+cmdLen])
            isExecMode = true
        }
    }
    req.Reply(true, nil)
```

#### 命令执行流程
1. **接收命令**: 从SSH exec请求中提取命令内容
2. **模型选择**: 根据用户名自动选择AI模型
3. **AI处理**: 将命令发送给AI助手处理
4. **返回响应**: 将AI响应直接返回给客户端
5. **关闭连接**: 完成后自动关闭SSH连接

### 核心函数

#### handleExecCommand
```go
func handleExecCommand(channel ssh.Channel, username, command string) {
    // 显示执行的命令
    channel.Write([]byte(fmt.Sprintf("执行命令: %s\r\n\r\n", command)))
    
    // 获取并选择模型（简化版本）
    models, err := ai.GetAvailableModels()
    selectedModel := ai.SelectModelByUsername(nil, models, username)
    
    // 创建AI助手并处理命令
    assistant := ai.NewAssistant(username)
    assistant.SetModel(selectedModel)
    assistant.ProcessMessage(command, channel, interrupt)
}
```

## 应用场景

### 1. 脚本自动化
```bash
#!/bin/bash
# 自动化AI咨询脚本
QUESTION="今天的天气如何？"
RESPONSE=$(ssh ai@server -p 2213 "$QUESTION")
echo "AI回答: $RESPONSE"
```

### 2. 批量处理
```bash
# 批量问题处理
questions=(
    "什么是机器学习？"
    "深度学习的应用领域有哪些？"
    "如何开始学习人工智能？"
)

for question in "${questions[@]}"; do
    echo "问题: $question"
    ssh ai@server -p 2213 "$question"
    echo "---"
done
```

### 3. API集成
```python
import subprocess

def ask_ai(question):
    """通过SSH向AI提问"""
    cmd = ['ssh', 'ai@localhost', '-p', '2213', question]
    result = subprocess.run(cmd, capture_output=True, text=True)
    return result.stdout

# 使用示例
answer = ask_ai("什么是Python？")
print(answer)
```

### 4. 命令行工具
```bash
# 创建别名简化使用
alias ask='ssh ai@localhost -p 2213'

# 使用别名
ask "今天适合做什么？"
ask "推荐一些学习资源"
```

## 配置说明

### 服务器配置
无需特殊配置，使用标准的SSHAI配置文件即可。exec模式会自动检测并处理。

### 客户端配置
```bash
# ~/.ssh/config
Host ai
    HostName your-server.com
    Port 2213
    User gpt
    StrictHostKeyChecking no
```

使用配置后可以简化命令：
```bash
ssh ai "你的问题"
```

## 性能优化

### 1. 快速启动
- 跳过交互式会话初始化
- 不显示欢迎信息和模型选择过程
- 直接进入AI处理流程

### 2. 资源管理
- 自动关闭连接释放资源
- 最小化内存使用
- 优化网络传输

### 3. 并发处理
- 支持多个并发exec请求
- 独立的会话处理
- 无状态设计

## 安全考虑

### 1. 命令验证
- 输入长度限制
- 字符编码验证
- 恶意命令过滤

### 2. 访问控制
- 用户身份验证
- 权限管理
- 频率限制

### 3. 日志记录
```go
log.Printf("接收到执行命令: %s", execCommand)
```

## 故障排除

### 常见问题

#### 1. 命令未执行
**症状**: SSH连接成功但命令没有被处理
**原因**: exec请求解析失败
**解决**: 检查SSH客户端是否正确发送exec请求

#### 2. 中文乱码
**症状**: 中文问题或回答显示乱码
**原因**: 字符编码问题
**解决**: 确保SSH客户端和服务器都使用UTF-8编码

#### 3. 连接超时
**症状**: SSH连接建立后长时间无响应
**原因**: AI处理时间过长
**解决**: 增加客户端超时时间或优化AI响应速度

### 调试方法

#### 1. 启用详细日志
```bash
# 服务器端
./sshai -v

# 客户端
ssh -v ai@server "test"
```

#### 2. 检查网络连接
```bash
# 测试基本连接
telnet server 2213

# 测试SSH连接
ssh -T ai@server
```

#### 3. 验证配置
```bash
# 检查配置文件
cat config.yaml

# 测试配置解析
./sshai --check-config
```

## 测试方法

### 自动化测试
```bash
# 运行测试脚本
./scripts/test_ssh_exec.sh
```

### 手动测试
```bash
# 基本功能测试
ssh test@localhost -p 2213 "hello"

# 中文测试
ssh test@localhost -p 2213 "你好"

# 复杂问题测试
ssh test@localhost -p 2213 "请解释量子计算的基本原理"
```

## 与交互式模式的对比

| 特性 | 执行命令模式 | 交互式模式 |
|------|-------------|-----------|
| 使用方式 | 一次性命令 | 持续对话 |
| 连接时间 | 短暂 | 长期 |
| 资源消耗 | 低 | 中等 |
| 适用场景 | 脚本化、自动化 | 人机交互 |
| 历史记录 | 无 | 有 |
| 上下文 | 无 | 有 |

## 未来改进

### 计划功能
- [ ] 支持多轮对话的exec模式
- [ ] 命令结果缓存
- [ ] 批量命令处理
- [ ] 异步响应模式

### 性能优化
- [ ] 连接池管理
- [ ] 响应压缩
- [ ] 智能超时控制

## 总结

SSH执行命令功能为SSHAI增加了强大的自动化能力，用户可以通过简单的SSH命令直接与AI交互，无需进入交互式会话。这个功能特别适合：

1. **脚本自动化**: 在脚本中集成AI功能
2. **批量处理**: 处理大量问题
3. **API集成**: 作为其他系统的AI接口
4. **快速查询**: 获取即时AI响应

通过完整的中文支持、智能模型选择和优化的性能，这个功能大大扩展了SSHAI的应用场景。