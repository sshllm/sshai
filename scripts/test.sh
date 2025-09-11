#!/bin/bash

echo "=== SSH AI 服务器测试脚本 ==="
echo ""

# 检查Go版本
echo "1. 检查Go版本："
go version
echo ""

# 检查编译结果
echo "2. 检查编译结果："
if [ -f "./sshai" ]; then
    echo "✓ 可执行文件 sshai 已生成"
    ls -la sshai
else
    echo "✗ 可执行文件未找到"
fi
echo ""

# 检查依赖
echo "3. 检查Go模块依赖："
go mod verify
echo ""

echo "4. 启动服务器测试："
echo "运行以下命令启动服务器："
echo "  ./run.sh"
echo "或者："
echo "  ./sshai"
echo ""
echo "然后在另一个终端中连接："
echo "  ssh gpt-5@localhost -p 2212"
echo ""
echo "注意：首次连接时选择 'yes' 接受主机密钥"