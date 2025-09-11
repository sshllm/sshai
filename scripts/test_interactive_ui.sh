#!/bin/bash

# 测试交互式用户界面功能
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== 测试交互式用户界面功能 ===${NC}"
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
cat > test_config.yaml << 'EOF'
# SSH AI 服务器配置文件 - 测试版本

# 服务器配置
server:
  port: "2214"
  welcome_message: "Hello!欢迎使用SSHAI测试环境！"
  prompt_template: "%s@test> "

# 认证配置
auth:
  password: ""  # 无密码认证便于测试
  login_prompt: ""
  login_success_message: "🎉 欢迎使用 SSHAI 测试环境\n\n💡 测试功能:\n  - 上下方向键: 浏览历史命令\n  - 左右方向键: 移动光标编辑\n  - Ctrl+A: 移动到行首\n  - Ctrl+E: 移动到行尾\n  - 支持中文输入和编辑\n\n"

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
echo -e "${BLUE}服务器将在端口 2214 启动${NC}"
echo -e "${BLUE}连接命令: ssh test@localhost -p 2214${NC}"
echo
echo -e "${YELLOW}测试说明:${NC}"
echo "1. 输入一些命令，如: 你好, hello world, 测试中文"
echo "2. 使用上下方向键浏览历史命令"
echo "3. 使用左右方向键移动光标进行编辑"
echo "4. 使用 Ctrl+A 移动到行首，Ctrl+E 移动到行尾"
echo "5. 测试中文字符的输入和编辑"
echo "6. 输入 'exit' 退出测试"
echo
echo -e "${RED}按 Ctrl+C 停止服务器${NC}"
echo

# 使用测试配置启动服务器
CONFIG_FILE=test_config.yaml ./sshai

# 清理测试文件
echo -e "${YELLOW}清理测试文件...${NC}"
rm -f test_config.yaml
mv config_backup.yaml config.yaml

echo -e "${GREEN}✓ 测试完成${NC}"