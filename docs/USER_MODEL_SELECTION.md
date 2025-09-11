# SSH AI 用户名和模型选择功能

## 功能概述

SSH AI 现在支持基于用户名的智能模型选择功能，提供更个性化的AI体验。

## 连接方式

### 1. 带用户名连接
```bash
ssh gpt@localhost -p 2212
ssh claude@localhost -p 2212
ssh deepseek@localhost -p 2212
```

### 2. 无用户名连接
```bash
ssh localhost -p 2212
```

## 模型选择逻辑

### 智能匹配
当用户提供用户名时，系统会：
1. 获取所有可用模型列表
2. 根据用户名匹配相关模型
3. 自动选择或提供选择菜单

### 匹配规则
- **精确匹配**：用户名完全包含在模型名中
- **部分匹配**：用户名的部分字符串匹配模型名
- **大小写不敏感**：匹配时忽略大小写

### 选择流程

#### 情况1：找到唯一匹配
```
欢迎, gpt!
正在获取可用模型...
已为用户 'gpt' 自动选择模型: gpt-4
```

#### 情况2：找到多个匹配
```
欢迎, claude!
正在获取可用模型...
为用户 'claude' 找到以下匹配的模型:
1. claude-3-sonnet
2. claude-3-opus
3. claude-3-haiku

请选择模型 (输入数字): 2
已选择模型: claude-3-opus
```

#### 情况3：无用户名或无匹配
```
Hello!
正在获取可用模型...
可用模型列表:
1. gpt-4
2. gpt-3.5-turbo
3. claude-3-sonnet
4. deepseek-v3

请选择模型 (输入数字): 4
已选择模型: deepseek-v3
```

## 用户体验改进

### 个性化欢迎
- 有用户名：`欢迎, {username}!`
- 无用户名：`Hello!`

### 动态提示符
连接后的提示符会显示当前使用的模型：
```
deepseek-v3@sshai> 你好
gpt-4@sshai> Hello
claude-3-sonnet@sshai> 请介绍一下你自己
```

### 智能默认
- 如果API调用失败，自动使用默认模型 `deepseek-v3`
- 保持向后兼容性

## 技术实现

### API集成
```go
// 获取模型列表
func getAvailableModels() ([]ModelInfo, error) {
    req, err := http.NewRequest("GET", APIBaseURL+"/models", nil)
    req.Header.Set("Authorization", "Bearer "+APIKey)
    // ...
}
```

### 用户名获取
```go
// 从SSH连接获取用户名
username := sshConn.User()
```

### 模型匹配
```go
// 智能匹配算法
func matchModelsByUsername(models []ModelInfo, username string) []ModelInfo {
    // 大小写不敏感的字符串匹配
    // ...
}
```

## 使用示例

### 示例1：GPT用户
```bash
$ ssh gpt@localhost -p 2212
欢迎, gpt!
正在获取可用模型...
已为用户 'gpt' 自动选择模型: gpt-4

gpt-4@sshai> 你好，请介绍一下你自己
我是GPT-4，一个大型语言模型...
```

### 示例2：多模型选择
```bash
$ ssh ai@localhost -p 2212
欢迎, ai!
正在获取可用模型...
为用户 'ai' 找到以下匹配的模型:
1. gpt-4
2. claude-3-sonnet
3. deepseek-v3

请选择模型 (输入数字): 2
已选择模型: claude-3-sonnet

claude-3-sonnet@sshai> 请解释量子计算的基本原理
```

### 示例3：无用户名连接
```bash
$ ssh localhost -p 2212
Hello!
正在获取可用模型...
可用模型列表:
1. gpt-4
2. gpt-3.5-turbo
3. claude-3-sonnet
4. deepseek-v3

请选择模型 (输入数字): 1
已选择模型: gpt-4

gpt-4@sshai> 开始对话
```

## 错误处理

### API失败处理
```
正在获取可用模型...
获取模型列表失败: network timeout
使用默认模型: deepseek-v3

deepseek-v3@sshai> 
```

### 无效选择处理
```
请选择模型 (输入数字): 99
请输入 1-4 之间的数字: 2
已选择模型: gpt-3.5-turbo
```

## 配置说明

### 默认模型设置
在 `main.go` 中修改：
```go
const (
    DefaultModel = "deepseek-v3"  // 修改默认模型
)
```

### API配置
```go
const (
    APIBaseURL = "https://ds.openugc.com/v1"
    APIKey     = ""
)
```

## 兼容性

- ✅ 向后兼容原有连接方式
- ✅ 支持所有SSH客户端
- ✅ 保持原有命令和功能
- ✅ 错误时自动降级到默认模型

## 性能优化

- **缓存机制**：可考虑缓存模型列表
- **超时控制**：API调用设置10秒超时
- **异步处理**：模型选择不阻塞其他功能
- **错误恢复**：网络失败时优雅降级

这个功能大大提升了SSH AI的用户体验，让每个用户都能快速找到最适合的AI模型进行对话。