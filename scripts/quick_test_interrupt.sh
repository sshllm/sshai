#!/bin/bash

# 快速测试中断功能
echo "=== 快速测试 Ctrl+C 中断功能 ==="

# 检查是否有配置文件
if [ ! -f "config.yaml" ]; then
    echo "创建测试配置文件..."
    cat > config.yaml << 'EOF'
server:
  host: "0.0.0.0"
  port: 2213
  welcome_message: "欢迎使用 SSHAI！"
  prompt_template: "%s@sshai.top> "

api:
  base_url: "https://api.deepseek.com/v1"
  api_key: "your-api-key-here"
  default_model: "deepseek-chat"
  timeout: 300

prompt:
  system_prompt: "你是一个有用的AI助手。"
  exec_prompt: "请回答以下问题："
  stdin_prompt: "请分析以下内容："

display:
  loading_animation_interval: 100
  thinking_animation_interval: 150

auth:
  enabled: false
EOF
    echo "请编辑 config.yaml 文件，设置正确的 API 密钥"
    echo "然后重新运行此脚本"
    exit 1
fi

# 检查 API 密钥
if grep -q "your-api-key-here" config.yaml; then
    echo "警告: 请在 config.yaml 中设置正确的 API 密钥"
    echo "当前使用的是示例密钥，可能无法正常工作"
fi

echo ""
echo "启动服务器进行测试..."
echo "请按照以下步骤测试："
echo ""
echo "1. 在另一个终端执行: ssh test@localhost -p 2213"
echo "2. 输入长问题: '请详细介绍一下深度学习的发展历史，包括各个重要节点'"
echo "3. 在回答过程中按 Ctrl+C"
echo "4. 观察是否立即中断"
echo "5. 输入下一个问题，检查是否正常"
echo ""
echo "预期结果: Ctrl+C 应该立即中断，不再显示 [已中断]"
echo ""
echo "按 Ctrl+C 停止服务器..."

# 启动服务器
./sshai