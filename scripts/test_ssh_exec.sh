#!/bin/bash

# 测试SSH执行命令功能
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== 测试SSH执行命令功能 ===${NC}"
echo

# 检查构建的二进制文件是否存在
if [ ! -f "sshai" ]; then
    echo -e "${RED}错误: 找不到二进制文件 sshai${NC}"
    echo -e "${YELLOW}请先运行构建命令: go build -o sshai cmd/main.go${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 找到二进制文件${NC}"

# 创建测试配置
echo -e "${YELLOW}创建测试配置...${NC}"
cp config.yaml config_backup.yaml

# 修改配置以便测试
cat > test_exec_config.yaml << 'EOF'
# SSH AI 服务器配置文件 - SSH执行命令测试版本

# 服务器配置
server:
  port: "2215"
  welcome_message: "Hello!欢迎使用SSHAI SSH执行命令测试环境！"
  prompt_template: "%s@exec-test> "

# 认证配置
auth:
  password: ""  # 无密码认证便于测试
  login_prompt: ""
  login_success_message: "🎉 SSHAI SSH执行命令测试环境\n\n💡 测试功能:\n  - 直接执行: ssh user@localhost -p 2215 \"你好\"\n  - 中文支持: ssh user@localhost -p 2215 \"什么是人工智能？\"\n  - 复杂命令: ssh user@localhost -p 2215 \"请解释一下机器学习的基本概念\"\n\n"

# AI API配置
api:
  base_url: "https://ds.openugc.com/v1"
  api_key: ""
  default_model: "test-model"
  timeout: 30

# 显示配置
display:
  line_width: 80
  thinking_animation_interval: 150
  loading_animation_interval: 100

# 证书配置
security:
  host_key_file: "keys/host_key.pem"

# 国际化配置
i18n:
  language: "zh-cn"
EOF

echo -e "${GREEN}✓ 测试配置创建完成${NC}"

# 启动测试服务器
echo -e "${YELLOW}启动测试服务器...${NC}"
echo -e "${BLUE}服务器将在端口 2215 启动${NC}"
echo

# 在后台启动服务器
CONFIG_FILE=test_exec_config.yaml ./sshai &
SERVER_PID=$!

# 等待服务器启动
echo -e "${YELLOW}等待服务器启动...${NC}"
sleep 3

# 测试SSH执行命令功能
echo -e "${BLUE}=== 开始测试SSH执行命令功能 ===${NC}"
echo

# 测试1: 简单问候
echo -e "${YELLOW}测试1: 简单问候${NC}"
echo "命令: ssh test@localhost -p 2215 \"你好\""
echo "预期: 直接执行命令并返回AI响应"
echo -e "${GREEN}执行结果:${NC}"
timeout 10s ssh -o StrictHostKeyChecking=no test@localhost -p 2215 "你好" || echo "测试完成"
echo

# 测试2: 中文问题
echo -e "${YELLOW}测试2: 中文问题${NC}"
echo "命令: ssh test@localhost -p 2215 \"什么是人工智能？\""
echo -e "${GREEN}执行结果:${NC}"
timeout 10s ssh -o StrictHostKeyChecking=no test@localhost -p 2215 "什么是人工智能？" || echo "测试完成"
echo

# 测试3: 英文问题
echo -e "${YELLOW}测试3: 英文问题${NC}"
echo "命令: ssh test@localhost -p 2215 \"What is machine learning?\""
echo -e "${GREEN}执行结果:${NC}"
timeout 10s ssh -o StrictHostKeyChecking=no test@localhost -p 2215 "What is machine learning?" || echo "测试完成"
echo

# 测试4: 复杂命令
echo -e "${YELLOW}测试4: 复杂命令${NC}"
echo "命令: ssh test@localhost -p 2215 \"请用简单的语言解释一下深度学习的基本原理\""
echo -e "${GREEN}执行结果:${NC}"
timeout 15s ssh -o StrictHostKeyChecking=no test@localhost -p 2215 "请用简单的语言解释一下深度学习的基本原理" || echo "测试完成"
echo

# 停止服务器
echo -e "${YELLOW}停止测试服务器...${NC}"
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

# 清理测试文件
echo -e "${YELLOW}清理测试文件...${NC}"
rm -f test_exec_config.yaml
mv config_backup.yaml config.yaml

echo
echo -e "${BLUE}=== SSH执行命令功能测试完成 ===${NC}"
echo -e "${GREEN}✓ 所有测试完成！${NC}"
echo
echo -e "${YELLOW}功能说明:${NC}"
echo "1. 支持直接通过SSH命令执行AI对话"
echo "2. 格式: ssh user@server \"your question\""
echo "3. 完全支持中文和英文问题"
echo "4. 自动选择模型并返回AI响应"
echo "5. 无需进入交互式会话"