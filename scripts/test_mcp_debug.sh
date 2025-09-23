#!/bin/bash

# MCP调试测试脚本

set -e

echo "=== SSHAI MCP调试测试 ==="

# 编译项目
echo "编译项目..."
go build -o sshai cmd/main.go

# 创建调试配置文件
echo "创建调试配置文件..."
cat > config_mcp_debug.yaml << 'EOF'
# SSH AI 服务器配置文件 - MCP调试版本

server:
  port: "2215"
  welcome_message: "Hello! 欢迎使用SSHAI MCP调试版本！"
  prompt_template: "%s@sshai-debug> "

auth:
  password: ""

api:
  base_url: "http://localhost:11434/v1"
  api_key: ""
  default_model: "gpt-oss:20b"
  timeout: 600
  temperature: 0.7

display:
  line_width: 80
  thinking_animation_interval: 150
  loading_animation_interval: 100

security:
  host_key_file: "host_key.pem"

i18n:
  language: "zh-cn"

prompt:
  system_prompt: "你是一个专业的AI助手，现在具备了MCP工具调用能力。当用户请求文件操作时，你应该调用相应的工具。可用的工具包括：list_files（列出文件）、read_file（读取文件）、write_file（写入文件）等。请根据用户需求选择合适的工具。"
  user_prompt: "用户问题："
  assistant_prompt: "AI助手："
  stdin_prompt: "请分析以下内容并提供相关的帮助、建议或解释："
  exec_prompt: ""

# MCP配置 - 启用调试
mcp:
  enabled: true
  refresh_interval: 60
  servers:
    # 模拟文件系统工具（用于调试）
    - name: "filesystem"
      transport: "stdio"
      command: ["echo", '{"tools":[{"name":"list_files","description":"List files in directory","inputSchema":{"type":"object","properties":{"path":{"type":"string","description":"Directory path"}}}}]}']
      enabled: false  # 先禁用，避免实际调用
EOF

echo "调试配置文件已创建: config_mcp_debug.yaml"

echo ""
echo "=== 调试说明 ==="
echo "1. 现在启动SSHAI时会输出详细的工具调用调试信息"
echo "2. 当AI尝试调用工具时，会在日志中显示："
echo "   - 工具ID和名称"
echo "   - 原始参数内容"
echo "   - JSON解析结果"
echo "3. 这将帮助我们诊断'unexpected end of JSON input'错误"
echo ""
echo "启动命令: ./sshai -c config_mcp_debug.yaml"
echo "然后尝试说: '列出当前目录的文件'"
echo ""

# 清理
rm -f config_mcp_debug.yaml

echo "调试版本准备完成！"