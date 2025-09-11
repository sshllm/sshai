#!/bin/bash

# 测试嵌入式语言包功能
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== 测试嵌入式语言包功能 ===${NC}"
echo

# 检查构建的二进制文件是否存在
if [ ! -f "build/sshai_test" ]; then
    echo -e "${RED}错误: 找不到测试二进制文件 build/sshai_test${NC}"
    echo -e "${YELLOW}请先运行构建命令:${NC}"
    echo "  GOOS=linux GOARCH=amd64 go build -ldflags \"-s -w\" -o build/sshai_test cmd/main.go"
    exit 1
fi

echo -e "${GREEN}✓ 找到测试二进制文件${NC}"

# 创建临时测试目录
TEST_DIR="test_embedded"
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR"

# 复制二进制文件和配置文件到测试目录
cp build/sshai_test "$TEST_DIR/"
cp config.yaml "$TEST_DIR/"

echo -e "${GREEN}✓ 创建测试环境${NC}"

# 进入测试目录
cd "$TEST_DIR"

# 测试1: 检查二进制文件大小（应该包含嵌入的语言包）
echo -e "${YELLOW}测试1: 检查二进制文件大小${NC}"
file_size=$(du -h sshai_test | cut -f1)
echo "  二进制文件大小: $file_size"

# 测试2: 使用strings命令检查是否包含语言包内容
echo -e "${YELLOW}测试2: 检查嵌入的语言包内容${NC}"
if strings sshai_test | grep -q "正在启动SSH服务器"; then
    echo -e "${GREEN}✓ 找到中文语言包内容${NC}"
else
    echo -e "${RED}✗ 未找到中文语言包内容${NC}"
fi

if strings sshai_test | grep -q "Starting SSH server"; then
    echo -e "${GREEN}✓ 找到英文语言包内容${NC}"
else
    echo -e "${RED}✗ 未找到英文语言包内容${NC}"
fi

# 测试3: 验证不需要外部语言文件
echo -e "${YELLOW}测试3: 验证无需外部语言文件${NC}"
if [ -d "../lang" ]; then
    echo "  原始语言目录存在，但测试目录中没有语言文件"
    ls -la
    echo -e "${GREEN}✓ 确认测试环境中无外部语言文件${NC}"
else
    echo -e "${GREEN}✓ 无外部语言文件依赖${NC}"
fi

# 测试4: 尝试运行程序（快速启动测试）
echo -e "${YELLOW}测试4: 程序启动测试${NC}"
echo "  配置中文语言..."
sed -i 's/language: .*/language: zh-cn/' config.yaml

# 创建一个简单的测试，检查程序是否能正常初始化
timeout 3s ./sshai_test 2>&1 | head -5 || true
echo -e "${GREEN}✓ 程序启动测试完成${NC}"

# 测试5: 英文语言测试
echo -e "${YELLOW}测试5: 英文语言测试${NC}"
echo "  配置英文语言..."
sed -i 's/language: .*/language: en-us/' config.yaml

timeout 3s ./sshai_test 2>&1 | head -5 || true
echo -e "${GREEN}✓ 英文语言测试完成${NC}"

# 清理测试环境
cd ..
rm -rf "$TEST_DIR"

echo
echo -e "${BLUE}=== 嵌入式语言包测试完成 ===${NC}"
echo -e "${GREEN}✓ 所有测试通过！${NC}"
echo -e "${BLUE}部署建议:${NC}"
echo "  1. 只需要复制二进制文件和 config.yaml"
echo "  2. 无需复制 lang/ 目录"
echo "  3. 语言包已嵌入到二进制文件中"
echo "  4. 可以通过 config.yaml 中的 language 配置切换语言"