#!/bin/bash

echo "=== SSH AI 服务器 - 改进功能测试 ==="
echo ""

# 检查编译结果
if [ ! -f "./sshai" ]; then
    echo "❌ 可执行文件不存在，请先编译："
    echo "   go build -o sshai main.go"
    exit 1
fi

echo "✅ 可执行文件存在"
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

echo "📋 新功能测试说明："
echo ""
echo "🎯 1. 中文输入测试："
echo "   连接后输入中文：你好，请介绍一下自己"
echo "   验证中文输入和显示是否正常"
echo ""
echo "⏳ 2. 加载动画测试："
echo "   输入任何问题后观察加载动画"
echo "   应该看到旋转的加载指示器"
echo ""
echo "📝 3. 自动换行测试："
echo "   询问一个需要长回答的问题"
echo "   观察AI回复是否自动换行"
echo ""
echo "🔄 4. 上下文记忆测试："
echo "   先问一个问题，然后问'刚才我问了什么？'"
echo "   验证AI是否记住了对话历史"
echo ""
echo "🆕 5. 新会话测试："
echo "   输入 '/new' 创建新会话"
echo "   再问'刚才我问了什么？'验证上下文已清除"
echo ""
echo "🔌 连接命令："
echo "   ssh gpt-5@localhost -p 2212"
echo ""
echo "💡 测试建议："
echo "   - 尝试输入长句子测试换行"
echo "   - 尝试中英文混合输入"
echo "   - 观察加载动画的流畅性"
echo "   - 测试退格键删除中文字符"
echo ""

# 等待用户输入
read -p "按回车键停止服务器..." 

# 停止服务器
kill $SERVER_PID 2>/dev/null
echo "服务器已停止"
echo ""
echo "🎉 测试完成！如有问题请检查服务器日志。"