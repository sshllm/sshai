#!/bin/bash

# 测试大模型回复截断问题修复

echo "=== 大模型回复截断问题修复测试 ==="
echo

# 检查程序是否存在
if [ ! -f "./sshai" ]; then
    echo "❌ 错误: 找不到可执行文件 ./sshai"
    echo "请先运行: go build -o sshai cmd/main.go"
    exit 1
fi

echo "🔧 修复说明:"
echo "- 问题: 大模型回复的前几个字符被截断"
echo "- 原因: WrapText函数在处理流式响应时误删除了开头字符"
echo "- 修复: 移除WrapText调用，直接输出内容"
echo

echo "📋 测试步骤:"
echo "1. 启动SSHAI服务器"
echo "2. 连接并测试简单问答"
echo "3. 检查回复是否完整（特别是开头字符）"
echo

echo "🚀 启动服务器..."
./sshai &
SERVER_PID=$!
sleep 3

echo "✅ 服务器已启动 (PID: $SERVER_PID)"
echo
echo "🔗 请在另一个终端运行以下命令进行测试:"
echo "ssh -p 2212 testuser@localhost"
echo
echo "💡 测试建议:"
echo "- 问: 'who are you' (检查回复是否以 'I am' 开头)"
echo "- 问: 'hello' (检查回复是否完整)"
echo "- 问: '你好' (检查中文回复是否完整)"
echo
echo "⚠️  修复前: 'who are you' 的回复可能是 ' am Gemma...' (缺少开头的 'I')"
echo "✅ 修复后: 'who are you' 的回复应该是 'I am Gemma...' (完整回复)"
echo
echo "按 Enter 键停止服务器..."
read

# 停止服务器
kill $SERVER_PID 2>/dev/null
echo "🛑 服务器已停止"
echo
echo "📝 如果测试通过，说明截断问题已修复！"