#!/bin/bash

echo "=== 深度思考模型测试 ==="
echo "启动SSH AI服务器..."

# 启动服务器（后台运行）
./sshai &
SERVER_PID=$!

# 等待服务器启动
sleep 2

echo "服务器已启动 (PID: $SERVER_PID)"
echo ""
echo "测试说明："
echo "1. 连接到服务器: ssh deepseek-v3@localhost -p 2212"
echo "2. 输入一个需要深度思考的问题，例如："
echo "   '请详细分析量子计算的工作原理，并解释它与传统计算的区别'"
echo "3. 观察是否显示思考过程和思考动画"
echo "4. 按 Ctrl+C 可以中断AI回答"
echo "5. 输入 /new 可以开始新对话"
echo "6. 输入 exit 退出"
echo ""
echo "预期行为："
echo "- 如果模型支持深度思考，会先显示 '--- 思考中 ---' 和思考内容"
echo "- 然后显示 '--- 开始回答 ---' 和最终回答"
echo "- 如果模型不支持深度思考，直接显示回答内容"
echo ""
echo "按任意键继续测试，或按 Ctrl+C 停止服务器..."
read -n 1

# 清理函数
cleanup() {
    echo ""
    echo "正在停止服务器..."
    kill $SERVER_PID 2>/dev/null
    wait $SERVER_PID 2>/dev/null
    echo "服务器已停止"
    exit 0
}

# 设置信号处理
trap cleanup SIGINT SIGTERM

# 等待用户操作
wait $SERVER_PID