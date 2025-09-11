#!/bin/bash

# 嵌入式语言包构建脚本
# 此脚本构建包含嵌入语言包的二进制文件，部署时只需要二进制文件和config.yaml

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目信息
PROJECT_NAME="sshai"
VERSION=$(date +%Y%m%d_%H%M%S)
BUILD_DIR="build"
DIST_DIR="dist"

echo -e "${BLUE}=== SSHAI 嵌入式语言包构建脚本 ===${NC}"
echo -e "${YELLOW}版本: ${VERSION}${NC}"
echo

# 清理旧的构建文件
echo -e "${YELLOW}清理旧的构建文件...${NC}"
rm -rf ${BUILD_DIR} ${DIST_DIR}
mkdir -p ${BUILD_DIR} ${DIST_DIR}

# 检查语言包文件
echo -e "${YELLOW}检查语言包文件...${NC}"
if [ ! -f "pkg/i18n/lang/lang-zh-cn.yaml" ]; then
    echo -e "${RED}错误: 找不到中文语言包文件 pkg/i18n/lang/lang-zh-cn.yaml${NC}"
    exit 1
fi

if [ ! -f "pkg/i18n/lang/lang-en-us.yaml" ]; then
    echo -e "${RED}错误: 找不到英文语言包文件 pkg/i18n/lang/lang-en-us.yaml${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 语言包文件检查完成${NC}"

# 构建不同平台的二进制文件
platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$platform"
    
    echo -e "${YELLOW}构建 ${GOOS}/${GOARCH}...${NC}"
    
    # 设置输出文件名
    output_name="${PROJECT_NAME}_${GOOS}_${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    # 构建二进制文件
    GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags "-X main.version=${VERSION} -s -w" \
        -o "${BUILD_DIR}/${output_name}" \
        cmd/main.go
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ ${GOOS}/${GOARCH} 构建成功${NC}"
        
        # 创建发布包
        dist_dir="${DIST_DIR}/${PROJECT_NAME}_${GOOS}_${GOARCH}"
        mkdir -p "$dist_dir"
        
        # 复制二进制文件
        cp "${BUILD_DIR}/${output_name}" "$dist_dir/"
        
        # 复制配置文件
        cp config.yaml "$dist_dir/"
        
        # 创建README
        cat > "$dist_dir/README.txt" << EOF
SSHAI - SSH AI Assistant (嵌入式语言包版本)
版本: ${VERSION}
平台: ${GOOS}/${GOARCH}

部署说明:
1. 将 ${output_name} 和 config.yaml 复制到目标服务器
2. 根据需要修改 config.yaml 配置文件
3. 运行 ./${output_name} 启动服务

注意: 此版本已将语言包嵌入到二进制文件中，无需额外的语言文件。

配置文件说明:
- server.host: 服务器监听地址
- server.port: 服务器监听端口
- api.base_url: AI API 基础URL
- language: 界面语言 (zh-cn 或 en-us)

更多信息请访问: https://github.com/your-repo/sshai
EOF
        
        # 创建压缩包
        cd "$DIST_DIR"
        if command -v zip >/dev/null 2>&1; then
            zip -r "${PROJECT_NAME}_${GOOS}_${GOARCH}_${VERSION}.zip" "${PROJECT_NAME}_${GOOS}_${GOARCH}/"
            echo -e "${GREEN}✓ 创建压缩包: ${PROJECT_NAME}_${GOOS}_${GOARCH}_${VERSION}.zip${NC}"
        fi
        cd - >/dev/null
        
    else
        echo -e "${RED}✗ ${GOOS}/${GOARCH} 构建失败${NC}"
    fi
done

# 显示构建结果
echo
echo -e "${BLUE}=== 构建完成 ===${NC}"
echo -e "${YELLOW}构建文件位置:${NC}"
ls -la ${BUILD_DIR}/
echo
echo -e "${YELLOW}发布包位置:${NC}"
ls -la ${DIST_DIR}/

# 显示文件大小
echo
echo -e "${YELLOW}二进制文件大小:${NC}"
for file in ${BUILD_DIR}/*; do
    if [ -f "$file" ]; then
        size=$(du -h "$file" | cut -f1)
        echo "  $(basename "$file"): $size"
    fi
done

echo
echo -e "${GREEN}✓ 所有构建任务完成！${NC}"
echo -e "${BLUE}部署说明: 只需要将对应平台的二进制文件和 config.yaml 复制到目标服务器即可${NC}"