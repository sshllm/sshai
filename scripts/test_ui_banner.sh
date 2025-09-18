#!/bin/bash

# UI Banner 和彩色输出功能测试脚本

set -e

echo "=== UI Banner 和彩色输出功能测试 ==="

# 创建测试目录
TEST_DIR="test_ui_banner"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

# 创建测试配置文件
echo "1. 创建测试配置文件..."
cat > config_ui_test.yaml << EOF
# SSH AI 服务器配置文件 - UI Banner 测试

# 服务器配置
server:
  port: "2216"
  welcome_message: "欢迎使用SSHAI UI测试！"
  prompt_template: "%s@sshai-ui.top> "

# 认证配置
auth:
  password: ""  # 无密码认证，方便测试
  login_prompt: "请输入访问密码: "

# AI API配置
api:
  base_url: "http://localhost:11434/v1"
  api_key: ""
  default_model: "gpt-oss:20b"
  timeout: 600

# 显示配置
display:
  line_width: 80
  thinking_animation_interval: 150
  loading_animation_interval: 100

# 证书配置
security:
  host_key_file: "host_key_ui_test.pem"

# 国际化配置
i18n:
  language: "zh-cn"

# AI提示词配置
prompt:
  system_prompt: "你是一个专业的AI助手，正在测试新的UI Banner功能。"
  user_prompt: "用户问题："
  assistant_prompt: "AI助手："
  stdin_prompt: "请分析以下内容："
  exec_prompt: ""
EOF

echo "   - 配置文件已创建: config_ui_test.yaml"
echo "   - 端口: 2216"
echo "   - 无密码认证模式"

# 创建测试脚本
cat > test_banner_display.sh << 'EOF'
#!/bin/bash

echo "=== UI Banner 显示测试 ==="

echo "测试说明："
echo "1. 新的彩色Banner设计"
echo "2. 版本信息显示（包含编译时间）"
echo "3. 彩色提示符输出"
echo "4. 优化的编译版本"
echo ""

echo "启动服务器命令："
echo "cd .. && ./sshai -c $PWD/config_ui_test.yaml"
echo ""

echo "连接测试命令："
echo "ssh -p 2216 testuser@localhost"
echo ""

echo "预期效果："
echo "- 启动时显示彩色Banner"
echo "- 包含版本、编译时间等信息"
echo "- SSH连接后显示彩色提示符"
echo "- 用户名、主机名、模型名使用不同颜色"
echo ""

echo "版本信息测试："
echo "cd .. && make version"
EOF

chmod +x test_banner_display.sh

# 创建版本信息测试脚本
cat > test_version_info.go << 'EOF'
package main

import (
	"fmt"
	"sshai/pkg/version"
	"sshai/pkg/ui"
)

func main() {
	fmt.Println("=== 版本信息测试 ===")
	
	// 测试版本信息
	buildInfo := version.GetBuildInfo()
	fmt.Printf("Version: %s\n", buildInfo.Version)
	fmt.Printf("Git Commit: %s\n", buildInfo.GitCommit)
	fmt.Printf("Build Time: %s\n", buildInfo.BuildTime)
	fmt.Printf("Go Version: %s\n", buildInfo.GoVersion)
	fmt.Printf("Platform: %s\n", buildInfo.Platform)
	
	fmt.Println("\n=== Banner 显示测试 ===")
	
	// 测试Banner生成
	banner := ui.GenerateBanner()
	fmt.Print(banner)
	
	fmt.Println("\n=== 彩色提示符测试 ===")
	
	// 测试彩色提示符
	prompt := ui.GeneratePrompt("testuser", "sshai.top", "gpt-oss:20b")
	fmt.Print(prompt)
	fmt.Println("这是一个测试命令")
	
	fmt.Println("\n=== 其他UI元素测试 ===")
	
	// 测试其他UI元素
	fmt.Println(ui.FormatStatus("连接成功", true))
	fmt.Println(ui.FormatStatus("连接失败", false))
	fmt.Println(ui.FormatInfo("这是一条信息"))
	fmt.Println(ui.FormatWarning("这是一条警告"))
	fmt.Println(ui.FormatError("这是一条错误"))
}
EOF

echo "2. 测试完成！"
echo ""
echo "测试文件已创建在目录: $TEST_DIR/"
echo "- config_ui_test.yaml (UI测试配置)"
echo "- test_banner_display.sh (Banner显示测试脚本)"
echo "- test_version_info.go (版本信息和UI测试程序)"
echo ""
echo "使用方法："
echo "1. 版本信息测试: cd .. && make version"
echo "2. UI元素测试: cd .. && go run $TEST_DIR/test_version_info.go"
echo "3. 启动UI测试服务器: cd .. && ./sshai -c $TEST_DIR/config_ui_test.yaml"
echo "4. 连接测试: ssh -p 2216 testuser@localhost"
echo ""
echo "UI Banner 和彩色输出功能测试准备完成！"

cd ..