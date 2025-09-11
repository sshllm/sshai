#!/bin/bash

echo "=== DeepSeek-R1 思考功能最终测试 ==="
echo "启动SSH AI服务器..."

# 启动服务器（后台运行）
./sshai &
SERVER_PID=$!

# 等待服务器启动
sleep 2

echo "服务器已启动 (PID: $SERVER_PID)"
echo ""
echo "修复说明："
echo "✅ 已添加对 reasoning_content 字段的支持"
echo "✅ 同时支持 reasoning 和 reasoning_content 字段"
echo "✅ 移除了调试代码"
echo ""
echo "测试步骤："
echo "1. 连接: ssh deepseek-r1@localhost -p 2212"
echo "2. 或者: ssh r1@localhost -p 2212"
echo "3. 输入需要思考的问题，例如："
echo "   '请解释量子纠缠的原理'"
echo "   '1+1为什么等于2？'"
echo "   '请分析人工智能的发展趋势'"
echo ""
echo "预期行为："
echo "✅ 显示 '--- 思考中 ---' 和思考动画"
echo "✅ 实时显示思考过程内容"
echo "✅ 显示 '--- 开始回答 ---' 分隔符"
echo "✅ 显示最终回答内容"
echo "✅ 支持 Ctrl+C 中断"
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
    echo "测试完成！如果思考内容正常显示，说明修复成功。"
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