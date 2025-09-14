# 模型切换功能修复文档

## 问题描述

在基于 go-openai 库的重构过程中，发现了一个关键问题：

**问题现象**：
- 用户在交互式终端中选择了不同的模型
- 但实际的 AI 客户端仍然使用配置文件中的默认模型
- 模型切换功能没有生效

## 问题分析

### 根本原因

在重构过程中，`OpenAIClient` 结构体缺少存储当前模型的字段：

```go
// 问题代码 - 缺少模型存储
type OpenAIClient struct {
    client   *openai.Client
    messages []openai.ChatCompletionMessage
    username string
    // 缺少 currentModel 字段
}
```

### 具体问题

1. **SetModel 方法无效**：
   ```go
   // 原有的 SetModel 方法是空实现
   func (c *OpenAIClient) SetModel(model string) {
       // 注释说明但没有实际实现
   }
   ```

2. **始终使用默认模型**：
   ```go
   // callStreamingAPI 中使用配置文件的默认模型
   req := openai.ChatCompletionRequest{
       Model:    cfg.API.DefaultModel, // 问题：始终使用默认模型
       Messages: c.messages,
       Stream:   true,
   }
   ```

## 解决方案

### 1. 添加模型存储字段

```go
// 修复后的结构体
type OpenAIClient struct {
    client       *openai.Client
    messages     []openai.ChatCompletionMessage
    username     string
    currentModel string // 新增：存储当前使用的模型
}
```

### 2. 正确初始化模型

```go
// 在 NewOpenAIClient 中初始化当前模型
return &OpenAIClient{
    client:       client,
    messages:     messages,
    username:     username,
    currentModel: cfg.API.DefaultModel, // 初始化为默认模型
}
```

### 3. 实现 SetModel 方法

```go
// 正确实现 SetModel 方法
func (c *OpenAIClient) SetModel(model string) {
    c.currentModel = model
    fmt.Printf("[DEBUG] 用户 %s 切换到模型: %s\n", c.username, model)
}
```

### 4. 使用当前模型发送请求

```go
// 在 callStreamingAPI 中使用当前模型
func (c *OpenAIClient) callStreamingAPI(ctx context.Context, channel ssh.Channel, showAnimation bool) {
    fmt.Printf("[DEBUG] 用户 %s 正在使用模型: %s\n", c.username, c.currentModel)
    
    req := openai.ChatCompletionRequest{
        Model:    c.currentModel, // 修复：使用当前设置的模型
        Messages: c.messages,
        Stream:   true,
    }
    // ...
}
```

### 5. 添加调试和查询方法

```go
// 获取当前模型（用于调试和验证）
func (c *OpenAIClient) GetCurrentModel() string {
    return c.currentModel
}

// 在 Assistant 中也添加对应方法
func (ai *Assistant) GetCurrentModel() string {
    return ai.client.GetCurrentModel()
}
```

## 修复验证

### 调试日志

修复后，系统会输出调试信息：

```
[DEBUG] 用户 testuser 切换到模型: gpt-4
[DEBUG] 用户 testuser 正在使用模型: gpt-4
```

### 测试步骤

1. **启动服务器**：
   ```bash
   ./scripts/test_model_switching.sh
   ```

2. **SSH 连接**：
   ```bash
   ssh test@localhost -p 2213
   ```

3. **选择模型**：
   - 在交互式终端中选择不同的模型
   - 观察服务器端的调试日志

4. **验证切换**：
   - 发送消息给 AI
   - 确认使用的是选择的模型而不是默认模型

### 预期行为

✅ **模型切换立即生效**：用户选择模型后，下次对话使用新模型  
✅ **多用户独立**：不同用户可以使用不同的模型  
✅ **持久化**：在同一会话中，模型选择保持有效  
✅ **调试可见**：通过日志可以验证模型切换  

## 技术细节

### 数据流程

1. **用户选择模型**：
   ```
   用户在终端选择 → SelectModelByUsername() → assistant.SetModel() → client.SetModel()
   ```

2. **存储模型**：
   ```
   client.SetModel() → c.currentModel = model → 调试日志输出
   ```

3. **使用模型**：
   ```
   用户发送消息 → callStreamingAPI() → 使用 c.currentModel 创建请求
   ```

### 兼容性保证

- ✅ **向后兼容**：所有现有接口保持不变
- ✅ **默认行为**：未选择模型时使用配置文件默认模型
- ✅ **错误处理**：无效模型名称由 go-openai 库处理

## 相关文件

### 修改的文件

1. **pkg/ai/client.go**：
   - 添加 `currentModel` 字段
   - 实现 `SetModel` 方法
   - 添加 `GetCurrentModel` 方法
   - 修复 `callStreamingAPI` 使用当前模型

2. **pkg/ai/assistant.go**：
   - 添加 `GetCurrentModel` 方法保持接口完整性

### 新增文件

3. **scripts/test_model_switching.sh**：
   - 模型切换功能测试脚本

4. **docs/MODEL_SWITCHING_FIX.md**：
   - 本修复文档

## 总结

这个修复解决了一个在重构过程中引入的回归问题。通过：

1. **正确的数据存储**：添加 `currentModel` 字段
2. **完整的方法实现**：实现 `SetModel` 和 `GetCurrentModel`
3. **调试支持**：添加日志输出便于验证
4. **测试工具**：提供测试脚本

现在模型切换功能应该能够正常工作，用户选择的模型会被正确使用，而不是始终使用默认模型。

---

**修复日期**：2025年9月13日  
**状态**：✅ 已完成，包含调试日志  
**影响范围**：AI 客户端模型管理，保持向后兼容