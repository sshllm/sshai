#!/bin/bash

echo "=== 构建SSHAI项目（包含语言包）==="

# 设置变量
PROJECT_NAME="sshai"
BUILD_DIR="build"
LANG_DIR="lang"

# 清理构建目录
echo "清理构建目录..."
rm -rf $BUILD_DIR
mkdir -p $BUILD_DIR

# 编译二进制文件
echo "编译二进制文件..."
cd /home/dev/sshai
go build -o $BUILD_DIR/$PROJECT_NAME cmd/main.go

if [ $? -ne 0 ]; then
    echo "编译失败！"
    exit 1
fi

echo "编译成功！"

# 复制语言包文件
echo "复制语言包文件..."
cp -r $LANG_DIR $BUILD_DIR/

# 复制配置文件
echo "复制配置文件..."
cp config.yaml $BUILD_DIR/
cp config-en.yaml $BUILD_DIR/

# 复制密钥文件夹（如果存在）
if [ -d "keys" ]; then
    echo "复制密钥文件..."
    cp -r keys $BUILD_DIR/
fi

# 复制文档
echo "复制文档..."
cp -r docs $BUILD_DIR/

# 创建启动脚本
echo "创建启动脚本..."
cat > $BUILD_DIR/start.sh << 'EOF'
#!/bin/bash
echo "启动SSHAI服务器..."
./sshai
EOF

chmod +x $BUILD_DIR/start.sh

# 显示构建结果
echo ""
echo "构建完成！文件结构："
echo "build/"
echo "├── sshai              # 主程序"
echo "├── start.sh           # 启动脚本"
echo "├── config.yaml        # 中文配置"
echo "├── config-en.yaml     # 英文配置"
echo "├── lang/              # 语言包目录"
echo "│   ├── lang-zh-cn.yaml"
echo "│   └── lang-en-us.yaml"
echo "├── keys/              # SSH密钥"
echo "└── docs/              # 文档"
echo ""
echo "使用方法："
echo "1. cd build"
echo "2. ./start.sh"
echo ""
echo "或者创建发布包："
echo "tar -czf ${PROJECT_NAME}-$(date +%Y%m%d).tar.gz -C build ."