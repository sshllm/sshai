# MCP JSON解析错误修复说明

## 问题分析

用户遇到的"解析工具参数失败: unexpected end of JSON input"错误通常由以下原因引起：

1. **流式响应中的工具调用信息不完整**：在流式API响应中，工具调用信息可能分多个chunk传输，导致JSON不完整
2. **AI模型返回空参数**：某些工具调用可能不需要参数，但AI模型返回空字符串
3. **JSON格式错误**：AI模型生成的参数不符合JSON格式规范

## 已实施的修复

### 1. 增强错误处理
```go
// 检查参数是否为空
if toolCall.Function.Arguments == "" {
    arguments = make(map[string]interface{})
    log.Printf("工具参数为空，使用空参数对象")
} else {
    if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments); err != nil {
        channel.Write([]byte(fmt.Sprintf("\r\n❌ 解析工具参数失败: %v\r\n原始参数: %s\r\n", err, toolCall.Function.Arguments)))
        // 添加错误到对话上下文，让AI知道调用失败
        return
    }
}
```

### 2. 详细调试信息
添加了完整的调试日志，包括：
- 工具ID、名称、参数内容
- 参数长度和JSON有效性检查
- 详细的错误信息输出

### 3. 上下文错误反馈
将工具调用错误添加到对话上下文中，让AI模型知道调用失败的原因，可以尝试重新调用或采用其他方式。

## 使用调试版本

1. **编译调试版本**：
```bash
go build -o sshai cmd/main.go
```

2. **启动并观察日志**：
```bash
./sshai -c config.yaml
```

3. **测试工具调用**：
尝试说"列出文件"或"创建文件"，观察控制台输出的调试信息。

## 预期的调试输出

正常情况下应该看到：
```
=== 工具调用调试信息 ===
工具ID: call_abc123
工具类型: function
函数名称: list_files
函数参数: {"path": "/home/user"}
参数长度: 20
参数JSON格式有效
=== 调试信息结束 ===
解析后的工具参数: map[path:/home/user]
```

异常情况下会看到：
```
=== 工具调用调试信息 ===
工具ID: call_abc123
工具类型: function
函数名称: list_files
函数参数: 
参数长度: 0
参数为空字符串
=== 调试信息结束 ===
工具参数为空，使用空参数对象
解析后的工具参数: map[]
```

## 进一步的解决方案

如果问题仍然存在，可能需要：

### 1. 优化系统提示词
```yaml
prompt:
  system_prompt: |
    你是一个专业的AI助手，具备MCP工具调用能力。
    
    重要规则：
    1. 调用工具时，参数必须是有效的JSON格式
    2. 如果工具不需要参数，请传递空对象 {}
    3. 字符串参数必须用双引号包围
    
    示例：
    - 列出文件: {"path": "/path/to/directory"}
    - 读取文件: {"path": "/path/to/file.txt"}
    - 创建文件: {"path": "/path/to/file.txt", "content": "file content"}
```

### 2. 处理流式响应中的工具调用
可能需要缓存不完整的工具调用信息，等待完整后再处理：

```go
// 在OpenAIClient中添加工具调用缓存
type OpenAIClient struct {
    // ... 其他字段
    pendingToolCalls map[string]*openai.ToolCall
}

// 在处理流式响应时缓存工具调用
if len(delta.ToolCalls) > 0 {
    for _, toolCall := range delta.ToolCalls {
        if existing, ok := c.pendingToolCalls[toolCall.ID]; ok {
            // 合并工具调用信息
            existing.Function.Arguments += toolCall.Function.Arguments
        } else {
            c.pendingToolCalls[toolCall.ID] = &toolCall
        }
        
        // 检查是否完整（简单检查JSON是否有效）
        if isValidJSON(toolCall.Function.Arguments) {
            c.handleToolCall(ctx, toolCall, channel, &assistantMessage)
            delete(c.pendingToolCalls, toolCall.ID)
        }
    }
}
```

### 3. 添加重试机制
```go
func (c *OpenAIClient) handleToolCallWithRetry(ctx context.Context, toolCall openai.ToolCall, channel ssh.Channel, assistantMessage *strings.Builder) {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        if err := c.handleToolCall(ctx, toolCall, channel, assistantMessage); err == nil {
            return
        }
        
        if i < maxRetries-1 {
            log.Printf("工具调用失败，重试 %d/%d", i+1, maxRetries)
            time.Sleep(time.Second * time.Duration(i+1))
        }
    }
}
```

## 测试建议

1. **使用简单的MCP服务器**：先用模拟服务器测试基本功能
2. **逐步增加复杂度**：从简单的无参数工具开始，逐步测试复杂参数
3. **监控日志输出**：密切关注调试信息，了解具体的失败原因
4. **测试不同的AI模型**：某些模型可能在工具调用方面表现更好

通过这些修复和调试手段，应该能够有效解决JSON解析错误问题。