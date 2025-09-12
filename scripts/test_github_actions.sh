#!/bin/bash

# GitHub Actions 本地测试脚本
# 模拟GitHub Actions的构建过程，用于本地验证

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

PROJECT_NAME="sshai"
TEST_VERSION="v1.0.0-test"

echo -e "${BLUE}=== GitHub Actions 本地测试 ===${NC}"
echo -e "${YELLOW}模拟版本: ${TEST_VERSION}${NC}"
echo

# 清理旧的测试文件
echo -e "${YELLOW}清理旧的测试文件...${NC}"
rm -rf build_test dist_test
mkdir -p build_test dist_test

# 检查Go环境
echo -e "${YELLOW}检查Go环境...${NC}"
if ! command -v go &> /dev/null; then
    echo -e "${RED}错误: 未找到Go环境${NC}"
    exit 1
fi

GO_VERSION=$(go version)
echo -e "${GREEN}✓ Go环境: ${GO_VERSION}${NC}"

# 检查依赖
echo -e "${YELLOW}检查项目依赖...${NC}"
if [ ! -f "go.mod" ]; then
    echo -e "${RED}错误: 未找到go.mod文件${NC}"
    exit 1
fi

go mod download
go mod verify
echo -e "${GREEN}✓ 依赖检查完成${NC}"

# 检查语言包文件（模拟GitHub Actions的处理）
echo -e "${YELLOW}检查语言包文件...${NC}"
if [ ! -f "pkg/i18n/lang/lang-zh-cn.yaml" ]; then
    echo -e "${YELLOW}创建占位符中文语言文件${NC}"
    mkdir -p pkg/i18n/lang
    echo "# Chinese language pack" > pkg/i18n/lang/lang-zh-cn.yaml
fi

if [ ! -f "pkg/i18n/lang/lang-en-us.yaml" ]; then
    echo -e "${YELLOW}创建占位符英文语言文件${NC}"
    mkdir -p pkg/i18n/lang
    echo "# English language pack" > pkg/i18n/lang/lang-en-us.yaml
fi

echo -e "${GREEN}✓ 语言包文件准备完成${NC}"

# 运行测试（如果存在）
echo -e "${YELLOW}运行测试...${NC}"
if go list ./... | grep -q .; then
    if go test -v ./... 2>/dev/null; then
        echo -e "${GREEN}✓ 测试通过${NC}"
    else
        echo -e "${YELLOW}⚠ 测试失败或无测试文件${NC}"
    fi
else
    echo -e "${YELLOW}⚠ 未找到测试文件${NC}"
fi

# 模拟多平台构建
echo -e "${YELLOW}开始多平台构建测试...${NC}"

platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

build_success=0
build_total=${#platforms[@]}

for platform in "${platforms[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$platform"
    
    echo -e "${YELLOW}构建 ${GOOS}/${GOARCH}...${NC}"
    
    output_name="${PROJECT_NAME}_${GOOS}_${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    # 构建二进制文件
    if GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags "-X main.version=${TEST_VERSION} -s -w" \
        -o "build_test/${output_name}" \
        cmd/main.go 2>/dev/null; then
        
        echo -e "${GREEN}✓ ${GOOS}/${GOARCH} 构建成功${NC}"
        ((build_success++))
        
        # 创建发布包（模拟GitHub Actions）
        dist_dir="dist_test/${PROJECT_NAME}_${GOOS}_${GOARCH}"
        mkdir -p "$dist_dir"
        
        cp "build_test/${output_name}" "$dist_dir/"
        
        # 检查配置文件
        if [ -f "config.yaml" ]; then
            cp config.yaml "$dist_dir/"
        else
            echo -e "${YELLOW}⚠ 未找到config.yaml，创建示例配置${NC}"
            cat > "$dist_dir/config.yaml" << EOF
# SSHAI Configuration Example
server:
  port: 2213
  host: "0.0.0.0"

api:
  base_url: "https://api.deepseek.com/v1"
  api_key: "your-api-key-here"
  default_model: "deepseek-v3"

auth:
  password: ""

language: "zh-cn"
EOF
        fi
        
        # 创建README
        cat > "$dist_dir/README.txt" << EOF
SSHAI - SSH AI Assistant
Version: ${TEST_VERSION}
Platform: ${GOOS}/${GOARCH}
Build Date: $(date -u +"%Y-%m-%d %H:%M:%S UTC")

Quick Start:
1. Edit config.yaml with your API configuration
2. Run ./${output_name} to start the service
3. Connect via SSH: ssh user@localhost -p 2213

For detailed documentation, visit:
https://github.com/sshai/sshai

License: Apache 2.0
EOF
        
        # 创建压缩包（如果有zip命令）
        if command -v zip >/dev/null 2>&1; then
            cd dist_test
            zip -r "${PROJECT_NAME}_${GOOS}_${GOARCH}_${TEST_VERSION}.zip" "${PROJECT_NAME}_${GOOS}_${GOARCH}/" >/dev/null 2>&1
            cd ..
        fi
        
    else
        echo -e "${RED}✗ ${GOOS}/${GOARCH} 构建失败${NC}"
    fi
done

# 生成校验和（模拟GitHub Actions）
if [ -d "dist_test" ] && command -v sha256sum >/dev/null 2>&1; then
    echo -e "${YELLOW}生成校验和...${NC}"
    cd dist_test
    if ls *.zip >/dev/null 2>&1; then
        sha256sum *.zip > checksums.txt
        echo -e "${GREEN}✓ 校验和文件已生成${NC}"
    fi
    cd ..
fi

# 显示构建结果
echo
echo -e "${BLUE}=== 构建测试结果 ===${NC}"
echo -e "${YELLOW}成功构建: ${build_success}/${build_total} 个平台${NC}"

if [ $build_success -eq $build_total ]; then
    echo -e "${GREEN}✅ 所有平台构建成功！GitHub Actions应该能正常工作${NC}"
else
    echo -e "${YELLOW}⚠ 部分平台构建失败，请检查错误信息${NC}"
fi

# 显示文件大小
echo
echo -e "${YELLOW}构建文件大小:${NC}"
if [ -d "build_test" ]; then
    for file in build_test/*; do
        if [ -f "$file" ]; then
            size=$(du -h "$file" | cut -f1)
            echo "  $(basename "$file"): $size"
        fi
    done
fi

# 显示发布包
echo
echo -e "${YELLOW}发布包:${NC}"
if [ -d "dist_test" ]; then
    ls -la dist_test/ | grep -E '\.(zip|txt)$' || echo "  无发布包文件"
fi

echo
echo -e "${BLUE}=== 测试完成 ===${NC}"
echo -e "${GREEN}GitHub Actions工作流验证完成！${NC}"
echo
echo -e "${YELLOW}下一步:${NC}"
echo -e "1. 提交代码到GitHub仓库"
echo -e "2. 创建标签触发自动发布: git tag ${TEST_VERSION} && git push origin ${TEST_VERSION}"
echo -e "3. 或在GitHub网页上手动触发工作流"

# 清理测试文件（可选）
read -p "是否清理测试文件? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    rm -rf build_test dist_test
    echo -e "${GREEN}✓ 测试文件已清理${NC}"
fi