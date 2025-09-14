#!/bin/bash

echo "=== SSHAI 中断功能测试 ==="
echo "1. 启动 SSHAI"
echo "2. 输入一个会产生长回复的问题"
echo "3. 在回复过程中按 Ctrl+C"
echo "4. 观察是否立即中断"
echo ""

echo "启动 SSHAI..."
cd /home/dev/sshai
./sshai

echo "测试完成"