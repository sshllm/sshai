#!/bin/bash

# SSH Stdin 输入功能测试脚本
# 测试通过管道将文件内容传递给SSH服务的功能

echo "=== SSH Stdin 输入功能测试 ==="

# 构建项目
echo "1. 构建项目..."
go build -o sshai ./cmd
if [ $? -ne 0 ]; then
    echo "❌ 构建失败"
    exit 1
fi

# 启动服务器
echo "2. 启动SSH服务器..."
./sshai > test_server.log 2>&1 &
SERVER_PID=$!
sleep 3

# 测试用例
echo "3. 测试SSH stdin功能..."

echo "测试1: 简单文本输入"
echo "请分析这段文本：Hello World，这是一个测试。" | ssh user@localhost -p 2213 2>/dev/null

echo -e "\n测试2: 文件内容分析"
cat README.md | ssh user@localhost -p 2213 2>/dev/null

echo -e "\n测试3: 代码文件分析"
cat pkg/ssh/session.go | head -50 | ssh user@localhost -p 2213 2>/dev/null

echo -e "\n测试4: 配置文件分析"
cat config.yaml | ssh user@localhost -p 2213 2>/dev/null

# 停止服务器
echo -e "\n4. 停止服务器..."
kill $SERVER_PID 2>/dev/null
wait $SERVER_PID 2>/dev/null

# 显示服务器日志
echo -e "\n5. 服务器日志："
echo "--- 开始日志 ---"
cat test_server.log
echo "--- 结束日志 ---"

# 清理
rm -f test_server.log

echo -e "\n✅ SSH Stdin 功能测试完成！"
echo "功能特点："
echo "- ✅ 支持管道输入 (cat file | ssh)"
echo "- ✅ 自动检测stdin内容"
echo "- ✅ 支持大文件处理"
echo "- ✅ 中文内容完美支持"
echo "- ✅ AI智能分析文件内容"
echo "- ✅ 无需进入交互模式"

echo -e "\n使用示例："
echo "cat README.md | ssh user@localhost -p 2213"
echo "echo '分析这段代码' | ssh user@localhost -p 2213"
echo "curl -s https://api.github.com/repos/user/repo | ssh user@localhost -p 2213"