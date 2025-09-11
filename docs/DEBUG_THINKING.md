# DeepSeek-R1 思考功能调试指南

## 问题描述
DeepSeek-R1 模型的思考内容没有显示，需要确定正确的API响应字段名。

## 调试步骤

### 1. 运行调试测试
```bash
./test_deepseek_r1.sh
```

### 2. 连接并测试
```bash
ssh deepseek-r1@localhost -p 2212
```

### 3. 输入测试问题
输入一个简单但需要思考的问题：
```
1+1等于几？请详细解释计算过程
```

### 4. 观察调试输出
查看是否有类似以下的调试信息：
```
[调试] 发现字段 thinking: 让我来计算1+1...
[调试] 发现字段 reasoning: 这是一个基本的数学问题...
[调试] 发现字段 thought: 我需要解释加法的原理...
```

## 可能的字段名

根据不同模型的实现，思考内容可能使用以下字段名：
- `reasoning` (当前代码使用)
- `thinking`
- `thought`
- `rationale`
- `explanation`
- `process`

## 修复方法

### 方法1：如果发现了正确字段名
假设调试显示字段名是 `thinking`，修改 `main.go` 中的代码：

```go
// 将这行：
if choice.Delta.Reasoning != "" {

// 改为：
if choice.Delta.Thinking != "" {

// 并将结构体定义中的：
Reasoning string `json:"reasoning,omitempty"`

// 改为：
Thinking string `json:"thinking,omitempty"`
```

### 方法2：支持多个可能的字段名
如果不确定字段名，可以同时检查多个字段：

```go
// 检查多个可能的思考字段
thinkingContent := ""
if choice.Delta.Reasoning != "" {
    thinkingContent = choice.Delta.Reasoning
} else if choice.Delta.Thinking != "" {
    thinkingContent = choice.Delta.Thinking
} else if choice.Delta.Thought != "" {
    thinkingContent = choice.Delta.Thought
}

if thinkingContent != "" {
    // 处理思考内容...
}
```

## 测试验证

修改代码后：
1. 重新编译：`go build -o sshai main.go`
2. 重新测试：`./test_deepseek_r1.sh`
3. 验证思考内容是否正确显示

## 常见问题

### Q: 没有看到任何调试信息
A: 可能模型不支持思考功能，或者思考内容在其他地方

### Q: 看到调试信息但内容很短
A: 可能是字段名正确但内容被截断，检查完整响应

### Q: 思考内容显示但格式混乱
A: 检查换行符处理和字符编码

## 下一步

1. 运行调试测试确定字段名
2. 根据结果修改代码
3. 测试验证功能正常
4. 移除调试代码（生产环境）

## 联系支持

如果调试后仍有问题，请提供：
- 调试输出的完整信息
- 使用的模型名称
- API响应的原始数据