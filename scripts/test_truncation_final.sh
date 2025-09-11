#!/bin/bash

echo "=== 测试AI回复截断修复 ==="
echo "启动SSH服务器..."

# 启动服务器（后台运行）
cd /home/dev/sshai
./sshai &
SERVER_PID=$!

# 等待服务器启动
sleep 3

echo "服务器已启动 (PID: $SERVER_PID)"
echo ""
echo "请手动测试以下场景："
echo "1. ssh -p 2222 test@localhost"
echo "2. 选择一个模型（如 gemma3:270M）"
echo "3. 输入: who are you"
echo "4. 检查回复是否以 'I am' 开头（而不是 ' am'）"
echo ""
echo "测试完成后，按任意键停止服务器..."
read -n 1

# 停止服务器
kill $SERVER_PID
echo "服务器已停止"