#!/bin/bash

# 测试中断信号修复
echo "=== 测试 Ctrl+C 中断信号修复 (深度修复版) ==="
echo ""

# 检查配置文件
if [ ! -f "config.yaml" ]; then
    echo "错误: 未找到 config.yaml 配置文件"
    echo "请先复制 config.yaml.example 并配置 API 密钥"
    exit 1
fi

# 启动服务器
echo "启动 SSHAI 服务器..."
./sshai &
SERVER_PID=$!

# 等待服务器启动
sleep 3

echo ""
echo "=== 修复说明 ==="
echo "本次修复实现了 HTTP 请求级别的中断机制："
echo "- 使用 Go context.Context 实现请求取消"
echo "- HTTP 请求可以在任何阶段被立即中断"
echo "- 解决了 SSE 流式响应中的阻塞问题"
echo ""
echo "=== 测试说明 ==="
echo "1. 服务器已启动，PID: $SERVER_PID"
echo "2. 请在另一个终端中执行以下命令进行测试："
echo "   ssh test@localhost -p 2213"
echo ""
echo "3. 测试步骤："
echo "   a) 输入一个会产生长回复的问题，如：'请详细介绍一下人工智能的发展历史'"
echo "   b) 在大模型回答过程中按下 Ctrl+C"
echo "   c) 观察是否能立即中断回复（应该在 100ms 内响应）"
echo "   d) 检查下次输入时是否还会显示 [已中断]"
echo ""
echo "4. 预期结果："
echo "   - Ctrl+C 应该能立即中断大模型的回复（包括 HTTP 请求）"
echo "   - 中断后应该显示 [已中断] 并立即返回到提示符"
echo "   - 下次输入时不应该显示 [已中断]"
echo "   - 中断响应时间应该 < 100ms"
echo ""
echo "按 Ctrl+C 停止服务器..."

# 等待用户中断
trap "echo ''; echo '停止服务器...'; kill $SERVER_PID 2>/dev/null; exit 0" INT

# 保持脚本运行
while true; do
    sleep 1
done