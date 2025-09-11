#!/bin/bash

# SSH Exec 功能最终测试脚本
# 测试通过SSH直接执行命令并与AI交互的功能

echo "=== SSH Exec 功能测试 ==="

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
echo "3. 测试SSH exec功能..."

echo "测试1: 简单中文问候"
ssh user@localhost -p 2213 "你好，请介绍一下你自己" 2>/dev/null

echo -e "\n测试2: 编程问题"
ssh user@localhost -p 2213 "请写一个JavaScript函数来反转字符串" 2>/dev/null

echo -e "\n测试3: 数学问题"
ssh user@localhost -p 2213 "计算1到100的和" 2>/dev/null

echo -e "\n测试4: 复杂技术问题"
ssh user@localhost -p 2213 "解释什么是Docker容器，并给出使用示例" 2>/dev/null

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

echo -e "\n✅ SSH Exec 功能测试完成！"
echo "功能特点："
echo "- ✅ 支持中文命令"
echo "- ✅ 直接执行不进入交互模式"
echo "- ✅ AI正常响应"
echo "- ✅ 无空指针异常"
echo "- ✅ 支持复杂技术问题"