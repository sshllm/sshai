#!/bin/bash

echo "=== SSH AI 历史命令功能测试 ==="
echo "功能：上下方向键浏览命令历史"
echo ""
echo "测试步骤："
echo "1. 输入几个命令"
echo "2. 按上方向键浏览历史"
echo "3. 按下方向键向前浏览"
echo ""

cd /home/dev/sshai
go build -o sshai ./cmd/main.go
./sshai