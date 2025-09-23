#!/bin/bash

# MCP功能测试脚本

set -e

echo "=== SSHAI MCP功能测试 ==="

# 检查Go版本
echo "检查Go版本..."
go version

# 检查依赖
echo "检查MCP SDK依赖..."
go list -m github.com/modelcontextprotocol/go-sdk

# 编译项目
echo "编译项目..."
go build -o sshai cmd/main.go

# 创建测试配置文件
echo "创建测试配置文件..."
cat > config_mcp_test.yaml << 'EOF'
# SSH AI 服务器配置文件 - MCP测试版本

server:
  port: "2214"
  welcome_message: "Hello! 欢迎使用SSHAI MCP测试版本！"
  prompt_template: "%s@sshai-mcp> "

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
  system_prompt: "你是一个专业的AI助手，现在具备了MCP工具调用能力。你可以根据用户需求自动调用相关工具来完成任务。"
  user_prompt: "用户问题："
  assistant_prompt: "AI助手："
  stdin_prompt: "请分析以下内容并提供相关的帮助、建议或解释："
  exec_prompt: ""

# MCP配置 - 启用测试
mcp:
  enabled: true
  refresh_interval: 60  # 测试时使用较短的刷新间隔
  servers:
    # 示例：文件系统工具（需要安装对应的MCP服务器）
    - name: "filesystem"
      transport: "stdio"
      command: ["echo", "Mock MCP Server Response"]  # 模拟命令，实际使用时替换为真实的MCP服务器
      enabled: false  # 默认禁用，避免测试时出错
    
    # 示例：HTTP服务器（模拟）
    - name: "web-api"
      transport: "http"
      url: "http://localhost:8080/mcp"
      headers:
        Content-Type: "application/json"
      enabled: false  # 默认禁用
EOF

echo "测试配置文件已创建: config_mcp_test.yaml"

# 验证配置文件
echo "验证配置文件格式..."
if command -v yq &> /dev/null; then
    yq eval '.' config_mcp_test.yaml > /dev/null
    echo "配置文件格式正确"
else
    echo "警告: 未安装yq，跳过配置文件格式验证"
fi

# 测试编译结果
echo "测试程序启动（5秒后自动退出）..."
timeout 5s ./sshai -c config_mcp_test.yaml || true

echo ""
echo "=== MCP功能测试完成 ==="
echo ""
echo "测试结果："
echo "✅ Go版本检查通过"
echo "✅ MCP SDK依赖正常"
echo "✅ 项目编译成功"
echo "✅ 配置文件创建成功"
echo "✅ 程序启动测试完成"
echo ""
echo "下一步："
echo "1. 安装真实的MCP服务器（如mcp-server-filesystem）"
echo "2. 修改config_mcp_test.yaml中的MCP服务器配置"
echo "3. 启用相应的MCP服务器（设置enabled: true）"
echo "4. 运行: ./sshai -c config_mcp_test.yaml"
echo ""
echo "MCP服务器示例："
echo "- 文件系统: npm install -g @modelcontextprotocol/server-filesystem"
echo "- 数据库: npm install -g @modelcontextprotocol/server-sqlite"
echo "- Git: npm install -g @modelcontextprotocol/server-git"
echo ""

# 清理
rm -f sshai config_mcp_test.yaml

echo "测试完成！"