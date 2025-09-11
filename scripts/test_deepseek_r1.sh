#!/bin/bash

echo "=== DeepSeek-R1 模型调试测试 ==="
echo "启动SSH AI服务器..."

# 启动服务器（后台运行）
./sshai &
SERVER_PID=$!

# 等待服务器启动
sleep 2

echo "服务器已启动 (PID: $SERVER_PID)"
echo ""
echo "测试说明："
echo "1. 连接到deepseek-r1模型: ssh deepseek-r1@localhost -p 2212"
echo "2. 输入一个需要思考的问题，例如："
echo "   '1+1等于几？请详细解释计算过程'"
echo "3. 观察调试信息，查看API返回的字段"
echo "4. 记录除了content和role之外的其他字段名"
echo ""
echo "预期行为："
echo "- 如果有思考内容，会显示 [调试] 发现字段 xxx: ..."
echo "- 根据调试信息确定正确的思考字段名"
echo ""
echo "按任意键开始测试，或按 Ctrl+C 停止..."
read -n 1

# 清理函数
cleanup() {
    echo ""
    echo "正在停止服务器..."
    kill $SERVER_PID 2>/dev/null
    wait $SERVER_PID 2>/dev/null
    echo "服务器已停止"
    echo ""
    echo "请根据调试信息确定思考字段名，然后更新代码"
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