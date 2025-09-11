#!/bin/bash

# SSH密码认证功能测试脚本

echo "=== SSH密码认证功能测试 ==="
echo

# 检查程序是否存在
if [ ! -f "./sshai" ]; then
    echo "❌ 错误: 找不到可执行文件 ./sshai"
    echo "请先运行: go build -o sshai cmd/main.go"
    exit 1
fi

# 备份原配置
cp config.yaml config.yaml.backup

echo "📋 测试场景:"
echo "1. 无密码认证模式测试 (会显示登录成功消息)"
echo "2. 密码认证模式测试 (需要密码认证)"
echo

# 测试1: 无密码认证模式
echo "🔓 测试1: 无密码认证模式"
echo "配置密码为空，应该可以直接连接..."

# 确保密码为空
sed -i 's/password: .*/password: ""/' config.yaml

echo "启动服务器 (无密码模式)..."
./sshai &
SERVER_PID=$!
sleep 2

echo "尝试连接 (应该成功)..."
timeout 5 ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2212 testuser@localhost "echo 'Connected successfully'" 2>/dev/null
if [ $? -eq 0 ]; then
    echo "✅ 无密码认证测试通过"
else
    echo "❌ 无密码认证测试失败"
fi

# 停止服务器
kill $SERVER_PID 2>/dev/null
sleep 1

echo

# 测试2: 密码认证模式
echo "🔐 测试2: 密码认证模式"
echo "设置密码为 'test123'..."

# 设置密码
sed -i 's/password: ""/password: "test123"/' config.yaml

echo "启动服务器 (密码模式)..."
./sshai &
SERVER_PID=$!
sleep 2

echo "尝试连接 (需要密码认证)..."
echo "注意: 实际使用时需要手动输入密码 'test123'"

# 这里只是启动服务器，实际密码测试需要手动进行
echo "服务器已启动，请在另一个终端运行以下命令测试:"
echo "ssh -p 2212 testuser@localhost"
echo "密码: test123"
echo
echo "按 Enter 键停止服务器..."
read

# 停止服务器
kill $SERVER_PID 2>/dev/null

# 恢复原配置
mv config.yaml.backup config.yaml

echo "✅ 测试完成，配置已恢复"
echo
echo "📝 测试说明:"
echo "- 无密码模式: password 配置为空字符串"
echo "- 密码模式: password 配置为具体密码"
echo "- 登录成功后会显示自定义的成功消息"
echo "- 可以通过修改 config.yaml 中的 auth 部分来自定义认证行为"