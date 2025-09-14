#!/bin/bash

# 测试模型切换功能
echo "=== 测试模型切换功能 ==="
echo ""

# 检查是否有配置文件
if [ ! -f "config.yaml" ]; then
    echo "❌ 未找到 config.yaml 文件"
    echo "请先创建配置文件"
    exit 1
fi

echo "=== 修复说明 ==="
echo "✅ 已修复模型切换问题："
echo "   - 在 OpenAIClient 中添加了 currentModel 字段"
echo "   - SetModel 方法现在正确存储选择的模型"
echo "   - callStreamingAPI 使用 currentModel 而不是默认模型"
echo "   - 添加了 GetCurrentModel 方法用于调试"
echo ""

echo "=== 测试步骤 ==="
echo "1. 启动服务器（在另一个终端）: ./sshai"
echo "2. SSH 连接: ssh test@localhost -p 2213"
echo "3. 在交互式终端中选择不同的模型"
echo "4. 发送消息，观察是否使用了选择的模型"
echo "5. 可以通过 API 响应或日志验证模型切换"
echo ""

echo "=== 调试信息 ==="
echo "如果需要调试，可以在代码中添加日志："
echo "fmt.Printf(\"当前使用模型: %s\\n\", c.currentModel)"
echo ""

echo "=== 预期行为 ==="
echo "✅ 用户选择模型后，后续对话应使用选择的模型"
echo "✅ 不同用户可以使用不同的模型"
echo "✅ 模型切换立即生效，无需重启"
echo ""

echo "按 Enter 键启动服务器进行测试..."
read

# 启动服务器
echo "启动 SSHAI 服务器..."
./sshai