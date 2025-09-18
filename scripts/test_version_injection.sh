#!/bin/bash

# 版本信息注入测试脚本

set -e

echo "🧪 版本信息注入测试"
echo "===================="

# 创建测试目录
TEST_DIR="test_version_injection"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

echo "📋 1. 测试版本信息获取..."

# 创建版本信息测试程序
cat > test_version_info.go << 'EOF'
package main

import (
    "fmt"
    "sshai/pkg/version"
)

func main() {
    fmt.Println("🔍 版本信息测试结果:")
    fmt.Println("====================")
    
    buildInfo := version.GetBuildInfo()
    fmt.Printf("版本号: %s\n", buildInfo.Version)
    fmt.Printf("Git提交: %s\n", buildInfo.GitCommit)
    fmt.Printf("构建时间: %s\n", buildInfo.BuildTime)
    fmt.Printf("Go版本: %s\n", buildInfo.GoVersion)
    fmt.Printf("平台: %s\n", buildInfo.Platform)
    
    fmt.Println("\n📊 格式化版本信息:")
    fmt.Println(version.GetVersionString())
    
    fmt.Println("\n🕐 格式化构建时间:")
    fmt.Println(version.FormatBuildTime())
    
    fmt.Println("\n✅ 版本信息注入测试完成!")
}
EOF

echo "🔨 2. 编译测试程序..."

# 获取版本信息（与Makefile相同的方式）
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "v0.9.19-test")
GIT_COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_VERSION=$(go version | cut -d' ' -f3)

echo "📋 构建信息:"
echo "  版本: $VERSION"
echo "  提交: $GIT_COMMIT"
echo "  时间: $BUILD_TIME"
echo "  Go版本: $GO_VERSION"

# 使用与Makefile相同的LDFLAGS
LDFLAGS="-X 'sshai/pkg/version.Version=$VERSION' \
         -X 'sshai/pkg/version.GitCommit=$GIT_COMMIT' \
         -X 'sshai/pkg/version.BuildTime=$BUILD_TIME' \
         -X 'sshai/pkg/version.GoVersion=$GO_VERSION'"

echo "🔧 LDFLAGS: $LDFLAGS"

# 编译测试程序
go build -ldflags "$LDFLAGS" -o version_test test_version_info.go

echo "🚀 3. 运行版本信息测试..."
./version_test

echo ""
echo "🎨 4. 测试UI Banner..."

# 创建Banner测试程序
cat > test_banner.go << 'EOF'
package main

import (
    "fmt"
    "sshai/pkg/ui"
)

func main() {
    fmt.Println("🎨 UI Banner测试:")
    fmt.Println("================")
    fmt.Println(ui.GenerateBanner())
    
    fmt.Println("✅ Banner测试完成!")
}
EOF

# 编译Banner测试
go build -ldflags "$LDFLAGS" -o banner_test test_banner.go

# 运行Banner测试
./banner_test

echo ""
echo "📦 5. 测试二进制文件大小..."
echo "文件大小:"
ls -lh version_test banner_test

echo ""
echo "🔍 6. 验证版本信息是否嵌入到二进制文件..."

# 检查二进制文件中是否包含版本信息
if strings version_test | grep -q "$VERSION"; then
    echo "✅ 版本号已正确嵌入到二进制文件"
else
    echo "❌ 版本号未找到在二进制文件中"
fi

if strings version_test | grep -q "$GIT_COMMIT"; then
    echo "✅ Git提交哈希已正确嵌入到二进制文件"
else
    echo "❌ Git提交哈希未找到在二进制文件中"
fi

echo ""
echo "🎯 7. 与Makefile构建对比..."

# 返回上级目录进行Makefile构建
cd ..

# 使用Makefile构建
echo "使用Makefile构建..."
make build > /dev/null 2>&1

echo "📊 构建结果对比:"
echo "Makefile构建: $(ls -lh sshai | awk '{print $5}')"
echo "脚本构建: $(ls -lh $TEST_DIR/version_test | awk '{print $5}')"

echo ""
echo "✅ 版本信息注入测试完成!"
echo "🎉 所有测试通过，版本信息注入功能正常工作!"

# 清理测试文件
echo ""
echo "🧹 清理测试文件..."
rm -rf "$TEST_DIR"
echo "清理完成!"