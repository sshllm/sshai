# Ctrl+C 中断信号修复文档

## 问题描述

用户在大模型回答时按下 Ctrl+C 无法立即中断响应，而是在大模型回复完成后才触发中断信号，导致下次输入时显示"已中断"提示。

### 问题现象
```log
deepseek-v3@sshai.top> 介绍一下自己
你好！我是 **SSHAI**（开源项目：[GitHub链接](https://github.com/sshllm/sshai)），一个专注于高效、简洁回答的AI助手。我的特点是：

1. **精准**——直接回答核心问题，避免冗余。
2. **实用**——提供可操作的解决方案或信息。
3. **开源透明**——代码公开，欢迎开发者参与改进。

无论是技术问题、生活建议，还是学习指导，都可以问我！需要什么帮助？ 😊
deepseek-v3@sshai.top>
^C
deepseek-v3@sshai.top> 哈哈

[已中断]
deepseek-v3@sshai.top> 
```

## 问题分析

### 根本原因

1. **阻塞读取问题**：在 `handleStreamResponse` 方法中，`resp.Body.Read(buffer)` 是阻塞调用，在大模型流式输出过程中会持续阻塞，导致无法及时检查中断信号。

2. **中断检查位置不当**：中断信号检查只在循环开始处进行，但大部分时间都在阻塞读取上。

3. **中断信号残留**：中断信号处理后没有正确清理，导致下次输入时仍然检测到中断状态。

### 技术细节

**原始代码问题**：
```go
for {
    select {
    case <-interrupt:
        channel.Write([]byte("\r\n[已中断]\r\n"))
        return
    default:
    }
    
    // 这里会长时间阻塞，无法及时响应中断
    n, err := resp.Body.Read(buffer)
    // ...
}
```

## 解决方案

### 核心改进：HTTP 请求级别的中断

**关键发现**：原始问题不仅仅是流式响应处理的问题，更重要的是 HTTP 请求本身的阻塞性质。即使优化了流式响应处理，如果 HTTP 请求没有开始或者在建立连接阶段，中断信号仍然无法生效。

**解决方案**：使用 Go 的 `context.Context` 机制在 HTTP 请求级别实现真正的中断：

```go
// 创建可取消的上下文
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// 启动中断监听 goroutine
go func() {
    select {
    case <-interrupt:
        cancel() // 立即取消HTTP请求
    case <-ctx.Done():
        return
    }
}()

// 创建带上下文的HTTP请求
req, err := http.NewRequestWithContext(ctx, "POST", cfg.API.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
```

### 1. 异步读取机制

使用 goroutine 将阻塞的读取操作与中断检查分离：

```go
// 创建数据读取通道
dataChan := make(chan []byte, 100)
errorChan := make(chan error, 1)
doneChan := make(chan bool, 1)

// 启动读取 goroutine
go func() {
    defer close(dataChan)
    defer close(errorChan)
    
    buffer := make([]byte, 32768)
    for {
        n, err := resp.Body.Read(buffer)
        if n > 0 {
            data := make([]byte, n)
            copy(data, buffer[:n])
            select {
            case dataChan <- data:
            case <-doneChan:
                return
            }
        }
        if err != nil {
            select {
            case errorChan <- err:
            case <-doneChan:
            }
            return
        }
    }
}()
```

### 2. 实时中断检查

在主循环中使用 `select` 语句同时监听多个通道：

```go
for !interrupted {
    select {
    case <-interrupt:
        // 立即响应中断信号
        interrupted = true
        doneChan <- true // 通知读取 goroutine 停止
        channel.Write([]byte("\r\n[已中断]\r\n"))
        return

    case data := <-dataChan:
        // 处理接收到的数据
        // ...

    case err := <-errorChan:
        // 处理读取错误
        // ...

    case <-time.After(100 * time.Millisecond):
        // 定期检查超时
        // ...
    }
}
```

### 3. 中断信号清理

在 `handleUserInput` 中添加中断信号清理逻辑：

```go
case 3: // Ctrl+C
    // 发送中断信号
    select {
    case interrupt <- true:
    default:
    }
    
    // 清空中断通道中的旧信号，避免影响下次输入
    go func() {
        time.Sleep(100 * time.Millisecond) // 等待中断处理完成
        // 清空中断通道
        for {
            select {
            case <-interrupt:
                // 清空通道中的信号
            default:
                return
            }
        }
    }()
    
    channel.Write([]byte("\r\n^C\r\n"))
    inputState.Clear()
    channel.Write([]byte(dynamicPrompt))
```

## 修复效果

### 修复前
- Ctrl+C 无法立即中断大模型回复
- 中断信号在回复完成后才被处理
- 下次输入时显示"已中断"

### 修复后
- Ctrl+C 能够立即中断大模型回复
- 中断响应时间 < 100ms
- 中断后正常返回提示符，不影响后续输入

## 测试方法

### 自动化测试
```bash
# 运行测试脚本
./scripts/test_interrupt_fix.sh
```

### 手动测试步骤

1. **启动服务器**：
   ```bash
   ./sshai
   ```

2. **连接测试**：
   ```bash
   ssh test@localhost -p 2213
   ```

3. **中断测试**：
   - 输入会产生长回复的问题：`请详细介绍一下人工智能的发展历史`
   - 在大模型回答过程中按下 `Ctrl+C`
   - 观察是否立即中断并返回提示符
   - 输入下一个问题，检查是否正常工作

### 预期结果

✅ **正确行为**：
```log
deepseek-v3@sshai.top> 请详细介绍一下人工智能的发展历史
人工智能（Artificial Intelligence，AI）的发展历史可以追溯到20世纪中叶，经历了多个重要阶段：

## 早期发展（1940s-1950s）
人工智能的概念最早可以追溯到古希腊神话，但现代AI的发展始于20世纪40年代：
- 1943年：麦卡洛克和皮茨提出了第一个神经网络模型
- 1950年：艾伦·图灵发表了著名的论文《计算机器与智能》，提出了"图灵测试"
^C
[已中断]
deepseek-v3@sshai.top> 你好
你好！有什么可以帮助你的吗？
deepseek-v3@sshai.top> 
```

## 技术改进点

### 1. 并发安全
- 使用通道进行 goroutine 间通信
- 避免数据竞争和内存泄漏
- 正确处理 goroutine 生命周期

### 2. 响应性能
- 中断响应时间从 "回复完成后" 优化到 "< 100ms"
- 保持流式输出的实时性
- 不影响正常的数据处理流程

### 3. 用户体验
- 立即响应用户的中断操作
- 清晰的中断提示信息
- 不影响后续交互

## 相关文件

- `pkg/ai/assistant.go` - 主要修复文件，包含流式响应处理逻辑
- `pkg/ssh/session.go` - 中断信号处理和清理逻辑
- `scripts/test_interrupt_fix.sh` - 测试脚本

## 注意事项

1. **兼容性**：修复保持了与现有功能的完全兼容
2. **性能**：异步处理不会影响正常的响应性能
3. **稳定性**：添加了错误处理和超时机制
4. **资源管理**：正确关闭 goroutine 和通道，避免资源泄漏

---

**修复日期**：2025年9月12日  
**状态**：✅ 已完成并测试通过  
**影响范围**：SSH 交互模式下的 Ctrl+C 中断功能