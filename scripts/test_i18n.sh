#!/bin/bash

echo "=== 测试多语言功能（YAML语言包版本）==="
echo "此测试将验证中文和英文界面切换"
echo ""

# 检查语言包文件是否存在
if [ ! -f "lang/lang-zh-cn.yaml" ]; then
    echo "错误：中文语言包文件不存在 (lang/lang-zh-cn.yaml)"
    exit 1
fi

if [ ! -f "lang/lang-en-us.yaml" ]; then
    echo "错误：英文语言包文件不存在 (lang/lang-en-us.yaml)"
    exit 1
fi

echo "✅ 语言包文件检查通过"
echo ""

# 备份原配置
cp config.yaml config.yaml.backup

echo "1. 测试中文界面（默认）"
echo "当前配置: zh-cn"
echo "语言包: lang/lang-zh-cn.yaml"
echo ""

# 启动服务器（后台运行）
cd /home/dev/sshai
./sshai &
SERVER_PID=$!

# 等待服务器启动
sleep 3

echo "服务器已启动 (PID: $SERVER_PID)"
echo "请连接测试中文界面："
echo "ssh -p 2212 test@localhost"
echo ""
echo "测试要点："
echo "- 欢迎消息应为中文"
echo "- 模型加载提示应为中文"
echo "- AI思考过程应为中文"
echo ""
echo "按任意键切换到英文测试..."
read -n 1

# 停止服务器
kill $SERVER_PID
sleep 2

echo ""
echo "2. 测试英文界面"
echo "正在切换配置到 en-us..."

# 修改配置为英文
sed -i 's/language: "zh-cn"/language: "en-us"/' config.yaml

echo "当前配置: en-us"
echo "语言包: lang/lang-en-us.yaml"
echo ""

# 重新启动服务器
./sshai &
SERVER_PID=$!

# 等待服务器启动
sleep 3

echo "服务器已启动 (PID: $SERVER_PID)"
echo "请连接测试英文界面："
echo "ssh -p 2212 test@localhost"
echo ""
echo "测试要点："
echo "- Welcome message should be in English"
echo "- Model loading prompt should be in English"
echo "- AI thinking process should be in English"
echo ""
echo "测试完成后，按任意键恢复原配置并停止服务器..."
read -n 1

# 停止服务器
kill $SERVER_PID

# 恢复原配置
mv config.yaml.backup config.yaml

echo ""
echo "配置已恢复，测试完成"
echo ""
echo "语言包系统特性："
echo "✅ 支持外部YAML语言包文件"
echo "✅ 支持动态语言切换"
echo "✅ 支持语言回退机制"
echo "✅ 支持打包部署"