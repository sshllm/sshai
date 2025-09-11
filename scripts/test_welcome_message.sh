#!/bin/bash

# 测试新的欢迎信息显示
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== 测试新的欢迎信息 ===${NC}"
echo

# 创建测试用的Go程序来读取和显示配置
cat > test_welcome.go << 'EOF'
package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Auth struct {
		LoginSuccessMessage string `yaml:"login_success_message"`
	} `yaml:"auth"`
}

func main() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("Failed to read config file:", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatal("Failed to parse config file:", err)
	}

	fmt.Println("=== 欢迎信息预览 ===")
	fmt.Println(config.Auth.LoginSuccessMessage)
	fmt.Println("=== 预览结束 ===")
}
EOF

echo -e "${YELLOW}编译测试程序...${NC}"
go build -o test_welcome test_welcome.go

echo -e "${YELLOW}测试配置文件解析...${NC}"
if ./test_welcome; then
    echo -e "${GREEN}✓ 配置文件解析成功，欢迎信息显示正常${NC}"
else
    echo -e "${RED}✗ 配置文件解析失败${NC}"
    exit 1
fi

# 清理测试文件
rm -f test_welcome test_welcome.go

echo
echo -e "${GREEN}✓ 欢迎信息测试完成！${NC}"
echo -e "${BLUE}新的ASCII艺术欢迎信息已成功配置并可以正常显示。${NC}"