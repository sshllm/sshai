#!/bin/bash

echo "=== SSH AI 服务器 BUG 修复测试 ==="
echo ""

echo "修复的问题："
echo "1. ✅ 用户在选择模型后，提示了两次：已选择模型"
echo "2. ✅ 大模型回答的内容没有正确换行"
echo "3. ✅ 大模型思考的内容没有输出"
echo "4. ✅ 思考模式下的界面，如思考中，思考完成等的界面混乱"
echo "5. ✅ 用户输入空白回车，并没有正确回车。用户输入exit命令，并没有退出"
echo ""

echo "启动模块化版本服务器..."
./sshai-modular &
SERVER_PID=$!

echo "服务器已启动 (PID: $SERVER_PID)"
echo ""

echo "测试说明："
echo "1. 连接服务器: ssh deepseek@localhost -p 2212"
echo "2. 测试模型选择是否只显示一次'已选择模型'"
echo "3. 测试回答内容是否正确换行"
echo "4. 测试思考内容是否正确显示"
echo "5. 测试空白回车是否正常"
echo "6. 测试 'exit' 命令是否能退出"
echo ""

echo "按任意键停止服务器..."
read -n 1

echo "停止服务器..."
kill $SERVER_PID
echo "测试完成！"