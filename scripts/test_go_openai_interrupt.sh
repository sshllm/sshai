#!/bin/bash

# 测试基于 go-openai 库的中断功能
echo "=== 测试基于 go-openai 库的 Ctrl+C 中断功能 ==="
echo ""

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
echo "=== 重构说明 ==="
echo "本次重构使用了成熟的 go-openai 库："
echo "- 原生支持流式响应和上下文取消"
echo "- 更好的错误处理和连接管理"
echo "- 专门为 OpenAI API 优化的实现"
echo ""
echo "=== 测试说明 ==="
echo "1. 在另一个终端执行: ssh test@localhost -p 2213"
echo "2. 输入长问题: '请详细介绍一下深度学习的发展历史，包括各个重要节点和突破'"
echo "3. 在回答过程中按 Ctrl+C"
echo "4. 观察是否立即中断（应该在几毫秒内响应）"
echo "5. 输入下一个问题，检查是否正常工作"
echo ""
echo "=== 预期改进 ==="
echo "✅ 立即中断: go-openai 库的原生 context 支持"
echo "✅ 更好的错误处理: 专门的 OpenAI 错误类型"
echo "✅ 稳定的连接: 成熟的连接池和重试机制"
echo "✅ 标准化: 遵循 OpenAI 官方 API 规范"
echo ""
echo "按 Ctrl+C 停止服务器..."

# 启动服务器
./sshai