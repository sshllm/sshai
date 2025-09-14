#!/bin/bash

# 测试思考内容输出功能
echo "=== 测试思考内容输出功能 ==="
echo ""

# 检查是否有配置文件
if [ ! -f "config.yaml" ]; then
    echo "❌ 未找到 config.yaml 文件"
    echo "请先创建配置文件"
    exit 1
fi

echo "=== 修复说明 ==="
echo "✅ 已修复思考内容输出问题："
echo "   - 使用 go-openai 库的 ReasoningContent 字段"
echo "   - 正确处理 DeepSeek 等模型的思考过程"
echo "   - 添加调试日志显示收到的思考内容"
echo "   - 保持原有的思考过程显示逻辑"
echo ""

echo "=== 支持思考的模型 ==="
echo "🧠 DeepSeek Reasoner 系列："
echo "   - deepseek-reasoner"
echo "   - deepseek-chat (部分版本)"
echo "🧠 其他支持推理的模型："
echo "   - 任何实现了 reasoning_content 字段的模型"
echo ""

echo "=== 测试步骤 ==="
echo "1. 启动服务器（在另一个终端）: ./sshai"
echo "2. SSH 连接: ssh test@localhost -p 2213"
echo "3. 选择支持思考的模型（如 deepseek-reasoner）"
echo "4. 输入需要推理的问题，例如："
echo "   - '请解释量子纠缠的原理'"
echo "   - '分析一下这个数学问题的解法'"
echo "   - '推理一下这个逻辑谜题'"
echo "5. 观察是否显示思考过程"
echo ""

echo "=== 预期行为 ==="
echo "✅ **思考阶段**："
echo "   - 显示 '🤔 思考过程:'"
echo "   - 实时输出模型的推理内容"
echo "   - 思考内容应该是模型的内部推理过程"
echo ""
echo "✅ **回答阶段**："
echo "   - 显示 '✨ 思考完成 (X.X秒)'"
echo "   - 显示 '💬 回答:'"
echo "   - 输出最终的回答内容"
echo ""

echo "=== 调试信息 ==="
echo "服务器端会显示调试日志："
echo "[DEBUG] 用户 xxx 正在使用模型: deepseek-reasoner"
echo "[DEBUG] 收到思考内容: 让我分析一下这个问题..."
echo ""

echo "=== 故障排除 ==="
echo "如果没有显示思考内容："
echo "1. 确认使用的是支持推理的模型"
echo "2. 检查服务器端调试日志"
echo "3. 确认 API 密钥有权限访问推理功能"
echo "4. 尝试不同的问题类型"
echo ""

echo "按 Enter 键启动服务器进行测试..."
read

# 启动服务器
echo "启动 SSHAI 服务器..."
./sshai