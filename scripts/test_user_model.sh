#!/bin/bash

echo "=== SSH AI 用户名和模型选择功能测试 ==="
echo ""

# 检查服务器是否运行
if ! pgrep -f "./sshai" > /dev/null; then
    echo "启动SSH AI服务器..."
    ./sshai &
    SERVER_PID=$!
    echo "服务器PID: $SERVER_PID"
    sleep 3
else
    echo "SSH AI服务器已在运行"
fi

echo ""
echo "=== 测试说明 ==="
echo "本脚本将演示以下功能："
echo "1. 带用户名连接 (ssh gpt@localhost -p 2212)"
echo "2. 无用户名连接 (ssh localhost -p 2212)"
echo "3. 模型自动选择和手动选择"
echo "4. 动态提示符显示"
echo ""

echo "=== 连接示例 ==="
echo ""

echo "1. 测试GPT用户连接："
echo "   ssh gpt@localhost -p 2212"
echo "   预期：自动匹配包含'gpt'的模型"
echo ""

echo "2. 测试Claude用户连接："
echo "   ssh claude@localhost -p 2212"
echo "   预期：自动匹配包含'claude'的模型"
echo ""

echo "3. 测试DeepSeek用户连接："
echo "   ssh deepseek@localhost -p 2212"
echo "   预期：自动匹配包含'deepseek'的模型"
echo ""

echo "4. 测试无用户名连接："
echo "   ssh localhost -p 2212"
echo "   预期：显示所有可用模型供选择"
echo ""

echo "5. 测试通用用户名连接："
echo "   ssh ai@localhost -p 2212"
echo "   预期：显示多个匹配模型供选择"
echo ""

echo "=== 功能特性 ==="
echo "✅ 智能用户名匹配"
echo "✅ 动态模型选择"
echo "✅ 个性化欢迎消息"
echo "✅ 动态提示符显示当前模型"
echo "✅ 错误处理和默认模型降级"
echo "✅ 向后兼容性"
echo ""

echo "=== 测试步骤 ==="
echo "1. 运行上述连接命令"
echo "2. 观察欢迎消息和模型选择过程"
echo "3. 注意提示符中的模型名称"
echo "4. 测试对话功能"
echo "5. 使用 /new 命令测试会话重置"
echo "6. 使用 Ctrl+C 测试中断功能"
echo "7. 使用 exit 命令退出"
echo ""

echo "=== 预期行为 ==="
echo "• 有用户名时显示: '欢迎, {username}!'"
echo "• 无用户名时显示: 'Hello!'"
echo "• 提示符格式: '{model}@sshai> '"
echo "• 自动匹配时直接进入对话"
echo "• 多匹配时显示选择菜单"
echo "• API失败时使用默认模型"
echo ""

echo "=== 开始测试 ==="
echo "请在新终端中运行以下命令进行测试："
echo ""
echo "# 测试1: GPT用户"
echo "ssh gpt@localhost -p 2212"
echo ""
echo "# 测试2: 无用户名"
echo "ssh localhost -p 2212"
echo ""
echo "# 测试3: 通用用户名"
echo "ssh ai@localhost -p 2212"
echo ""

echo "按 Ctrl+C 停止服务器"
if [ ! -z "$SERVER_PID" ]; then
    wait $SERVER_PID
fi