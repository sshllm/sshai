#!/bin/bash

echo "=== 测试模型缓存功能 ==="
echo "此测试将验证模型列表的5分钟缓存机制"
echo ""

# 启动服务器（后台运行）
cd /home/dev/sshai
./sshai &
SERVER_PID=$!

# 等待服务器启动
sleep 3

echo "服务器已启动 (PID: $SERVER_PID)"
echo ""
echo "测试步骤："
echo "1. 第一次连接 - 应该从远程获取模型列表"
echo "2. 立即第二次连接 - 应该使用缓存的模型列表（更快）"
echo "3. 等待6分钟后连接 - 缓存过期，重新获取模型列表"
echo ""
echo "请按以下步骤手动测试："
echo ""
echo "第一次测试（获取新数据）："
echo "ssh -p 2222 test@localhost"
echo "观察模型加载时间"
echo "输入 exit 退出"
echo ""
echo "第二次测试（使用缓存）："
echo "ssh -p 2222 test@localhost"  
echo "观察模型加载时间（应该更快）"
echo "输入 exit 退出"
echo ""
echo "如需测试缓存过期，请等待6分钟后再次连接"
echo ""
echo "测试完成后，按任意键停止服务器..."
read -n 1

# 停止服务器
kill $SERVER_PID
echo ""
echo "服务器已停止"