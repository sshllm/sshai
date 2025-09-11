#!/bin/bash

echo "=== 中文字符乱码修复测试 ==="
echo "启动SSH AI服务器..."

# 启动服务器（后台运行）
./sshai &
SERVER_PID=$!

# 等待服务器启动
sleep 2

echo "服务器已启动 (PID: $SERVER_PID)"
echo ""
echo "修复说明："
echo "✅ 修复了UTF-8字符截断问题"
echo "✅ 确保字符串切片不会在多字节字符中间断开"
echo "✅ 改进了findBreakPosition函数的字节位置计算"
echo "✅ 保持中文字符完整性"
echo ""
echo "测试用例："
echo "1. 输入: 你好"
echo "2. 输入: 请详细解释人工智能的发展历程和未来趋势"
echo "3. 输入: 中文测试：这是一个很长的中文句子，用来测试自动换行功能是否正确处理中文字符"
echo ""
echo "预期结果："
echo "- 中文字符显示正常，无乱码"
echo "- 自动换行不会截断中文字符"
echo "- 思考内容中的中文正确显示"
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
    echo "测试完成！中文字符应该正常显示，无乱码。"
    exit 0
}

# 设置信号处理
trap cleanup SIGINT SIGTERM

echo ""
echo "现在可以连接测试："
echo "ssh deepseek-r1@localhost -p 2212"
echo ""
echo "测试命令："
echo "1. 你好"
echo "2. 请详细解释人工智能的发展历程"
echo "3. 中文测试：这是一个很长的中文句子，用来测试自动换行功能"
echo ""
echo "按 Ctrl+C 停止服务器..."

# 等待用户操作
wait $SERVER_PID