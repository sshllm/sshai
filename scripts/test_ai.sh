#!/bin/bash

echo "=== SSH AI 服务器 - AI功能测试 ==="
echo ""

# 检查编译结果
if [ ! -f "./sshai" ]; then
    echo "❌ 可执行文件不存在，请先编译："
    echo "   go build -o sshai main.go"
    exit 1
fi

echo "✅ 可执行文件存在"
echo ""

# 检查网络连接
echo "🌐 测试API连接..."
curl -s --connect-timeout 5 https://ds.openugc.com/v1 > /dev/null
if [ $? -eq 0 ]; then
    echo "✅ API服务可访问"
else
    echo "⚠️  API服务连接测试失败，但可能是正常的（某些API不响应根路径）"
fi
echo ""

echo "🚀 启动SSH AI服务器..."
echo "服务器将在后台启动，监听端口2212"
echo ""

# 启动服务器（后台运行）
./sshai &
SERVER_PID=$!

echo "服务器PID: $SERVER_PID"
echo ""

# 等待服务器启动
sleep 2

echo "📋 测试说明："
echo "1. 服务器已在后台启动"
echo "2. 在新终端中运行以下命令连接："
echo "   ssh gpt-5@localhost -p 2212"
echo ""
echo "3. 测试功能："
echo "   - 输入普通文本进行AI对话"
echo "   - 输入 '/new' 创建新会话"
echo "   - 输入 'exit' 退出连接"
echo ""
echo "4. 停止服务器："
echo "   kill $SERVER_PID"
echo ""

# 等待用户输入
read -p "按回车键停止服务器..." 

# 停止服务器
kill $SERVER_PID 2>/dev/null
echo "服务器已停止"