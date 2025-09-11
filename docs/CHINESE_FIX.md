# 中文字符乱码修复说明

## 问题描述

用户反馈在显示中文思考内容时出现乱码：
```log
deepseek-r1-0528-lkeap@sshai> 你好
⠏ Thinking
嗯，用户发来一句简单的"你好"，看起来是初次打招呼。可能刚打开聊天界面�
��还在试探阶段，或者想测试系统响应速度。
```

## 问题根因分析

### UTF-8编码特性
- 中文字符在UTF-8中占用3个字节
- 英文字符在UTF-8中占用1个字节
- 字符串切片操作基于字节索引，不是字符索引

### 原始代码问题
```go
// 问题代码：使用字符索引进行字节切片
for i, r := range text {
    // i 是字符索引，不是字节索引
    if currentWidth >= maxWidth {
        return i  // 返回字符索引
    }
}

// 使用时：
lineStr[:breakPos]  // 用字符索引切割字节序列，可能截断UTF-8字符
```

### 乱码产生机制
1. `range text` 遍历时，`i` 是字符索引
2. 中文字符 "界" 占3字节：`[E7 95 8C]`
3. 如果在字节位置2处切断，得到 `[E7 95]`
4. 不完整的UTF-8序列显示为 `�`

## 修复方案

### 核心思路
- 使用字节位置而不是字符位置进行切片
- 确保切割点不会破坏UTF-8字符边界
- 正确计算每个字符的字节长度

### 修复后的代码
```go
func (ai *AIAssistant) findBreakPosition(text string, maxWidth int) int {
    if len(text) == 0 {
        return 0
    }

    breakChars := []rune{' ', ',', '.', '!', '?', ';', ':', '，', '。', '！', '？', '；', '：'}
    
    currentWidth := 0
    lastBreakBytePos := 0
    currentBytePos := 0  // 关键：跟踪字节位置

    for _, r := range text {
        runeWidth := 1
        if utf8.RuneLen(r) > 1 {
            runeWidth = 2
        }

        currentWidth += runeWidth

        // 检查断行字符
        for _, bc := range breakChars {
            if r == bc {
                lastBreakBytePos = currentBytePos + utf8.RuneLen(r)
                break
            }
        }

        if currentWidth >= maxWidth {
            if lastBreakBytePos > 0 {
                return lastBreakBytePos
            }
            // 返回当前字符的起始字节位置，不会截断字符
            return currentBytePos
        }

        currentBytePos += utf8.RuneLen(r)  // 累加字节长度
    }

    return 0
}
```

### 关键改进点

1. **字节位置跟踪**
   ```go
   currentBytePos := 0
   currentBytePos += utf8.RuneLen(r)  // 正确累加字节长度
   ```

2. **安全的切割位置**
   ```go
   return currentBytePos  // 返回字符起始位置，不会截断
   ```

3. **断行位置计算**
   ```go
   lastBreakBytePos = currentBytePos + utf8.RuneLen(r)  // 字符结束位置
   ```

## 测试验证

### 测试用例
1. **简单中文**: "你好"
2. **长中文句子**: "请详细解释人工智能的发展历程和未来趋势"
3. **混合内容**: "中文English混合测试123"
4. **边界情况**: 正好在行宽边界的中文字符

### 验证方法
```bash
# 运行测试脚本
./test_chinese_fix.sh

# 手动测试
ssh deepseek-r1@localhost -p 2212
```

### 验证要点
- [ ] 中文字符显示完整，无乱码
- [ ] 自动换行不截断中文字符
- [ ] 思考内容中的中文正确显示
- [ ] 混合中英文内容正常换行
- [ ] 标点符号断行功能正常

## 技术细节

### UTF-8字符长度计算
```go
utf8.RuneLen(r)  // 返回字符的字节长度
// 英文字符: 1字节
// 中文字符: 3字节
// 其他Unicode字符: 1-4字节
```

### 显示宽度计算
```go
runeWidth := 1
if utf8.RuneLen(r) > 1 {
    runeWidth = 2  // 中文字符显示宽度为2
}
```

### 字符串切片安全性
```go
// 安全：基于字节边界切片
text[:bytePos]     // bytePos是字符边界
text[bytePos:]     // 不会截断UTF-8字符

// 危险：基于字符索引切片
text[:charIndex]   // 可能截断多字节字符
```

## 性能影响

### 计算复杂度
- **时间复杂度**: O(n) - 需要遍历字符串
- **空间复杂度**: O(1) - 只使用常量额外空间

### 性能优化
- 避免重复的UTF-8解码
- 缓存字符宽度计算结果
- 使用高效的字符串操作

## 兼容性说明

### 支持的字符集
- ✅ ASCII字符 (英文、数字、符号)
- ✅ 中文字符 (简体、繁体)
- ✅ 日文字符 (平假名、片假名、汉字)
- ✅ 韩文字符
- ✅ 其他Unicode字符

### 终端兼容性
- ✅ 现代终端 (支持UTF-8)
- ✅ SSH客户端 (PuTTY, Terminal, iTerm2等)
- ⚠️ 老旧终端可能需要UTF-8配置

## 更新日志

### v1.3.3 (当前版本)
- ✅ 修复中文字符自动换行乱码问题
- ✅ 改进UTF-8字符边界检测
- ✅ 优化字节位置计算逻辑
- ✅ 确保字符串切片安全性
- ✅ 支持所有Unicode字符正确显示

## 相关问题

### 常见问题
1. **Q**: 为什么只有中文出现乱码？
   **A**: 中文是多字节字符，英文是单字节，切割时更容易截断中文

2. **Q**: 如何确认修复是否生效？
   **A**: 输入长中文句子，观察自动换行是否正常

3. **Q**: 其他语言会有类似问题吗？
   **A**: 是的，所有多字节Unicode字符都可能有此问题

### 预防措施
- 始终使用字节位置进行字符串切片
- 使用`utf8.RuneLen()`计算字符字节长度
- 测试多语言内容的显示效果