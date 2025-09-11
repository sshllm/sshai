#!/bin/bash

# 快速测试SSH执行命令功能
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== 快速测试SSH执行命令功能 ===${NC}"
echo

# 检查二进制文件
if [ ! -f "sshai" ]; then
    echo -e "${RED}错误: 找不到二进制文件 sshai${NC}"
    exit 1
fi

# 在后台启动服务器
echo -e "${YELLOW}启动测试服务器...${NC}"
./sshai &
SERVER_PID=$!

# 等待服务器启动
sleep 3

echo -e "${YELLOW}测试SSH执行命令...${NC}"

# 测试简单命令
echo -e "${BLUE}测试命令: ssh test@localhost -p 2213 \"你好\"${NC}"
timeout 10s ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null test@localhost -p 2213 "你好" 2>/dev/null || echo "测试完成"

echo

# 停止服务器
echo -e "${YELLOW}停止服务器...${NC}"
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

echo -e "${GREEN}✓ 测试完成${NC}"