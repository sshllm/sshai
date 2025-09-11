#!/bin/bash

echo "=== 思考界面优化测试 ==="
echo "启动SSH AI服务器..."

# 启动服务器（后台运行）
./sshai &
SERVER_PID=$!

# 等待服务器启动
sleep 2

echo "服务器已启动 (PID: $SERVER_PID)"
echo ""
echo "界面优化说明："
echo "✅ 移除了混乱的思考动画"
echo "✅ 简化思考状态显示为 '⠏ Thinking'"
echo "✅ 直接显示思考内容，无动画干扰"
echo "✅ 思考完成后显示 '⠙ Done. Xs' 时间统计"
echo "✅ 清理了界面显示逻辑"
echo ""
echo "预期界面效果："
echo "deepseek-r1@sshai> Hello"
echo "⠏ Thinking"
echo "Okay, the user said \"Hello\". That's a friendly greeting..."
echo "I should respond in a welcoming manner."
echo "⠙ Done. 2.3s"
echo "Hello! How can I assist you today?"
echo "deepseek-r1@sshai>"
echo ""
echo "测试步骤："
echo "1. 连接: ssh deepseek-r1@localhost -p 2212"
echo "2. 输入: Hello"
echo "3. 观察界面是否清晰整洁"
echo "4. 测试其他需要思考的问题"
echo ""
echo "按任意键开始测试..."
read -n 1

# 清理函数
cleanup() {
    echo ""
    echo "正在停止服务器..."
    kill $SERVER_PID 2>/dev/null
    wait $SERVER_PID 2>/dev/null
    echo "服务器已停止"
    echo ""
    echo "测试完成！界面应该更加清晰整洁。"
    exit 0
}

# 设置信号处理
trap cleanup SIGINT SIGTERM

echo ""
echo "现在可以连接测试："
echo "ssh deepseek-r1@localhost -p 2212"
echo ""
echo "按 Ctrl+C 停止服务器..."

# 等待用户操作
wait $SERVER_PID