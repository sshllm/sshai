#!/bin/bash

# SSH Keys 免密登录功能测试脚本

set -e

echo "=== SSH Keys 免密登录功能测试 ==="

# 创建测试目录
TEST_DIR="test_ssh_keys"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

# 生成测试用的SSH密钥对
echo "1. 生成测试SSH密钥对..."
ssh-keygen -t rsa -b 2048 -f test_key -N "" -C "test@sshai.top" > /dev/null 2>&1
ssh-keygen -t ed25519 -f test_key_ed25519 -N "" -C "test2@sshai.top" > /dev/null 2>&1

echo "   - RSA密钥已生成: test_key, test_key.pub"
echo "   - Ed25519密钥已生成: test_key_ed25519, test_key_ed25519.pub"

# 读取公钥内容
RSA_PUBLIC_KEY=$(cat test_key.pub)
ED25519_PUBLIC_KEY=$(cat test_key_ed25519.pub)

# 创建测试配置文件
echo "2. 创建测试配置文件..."
cat > config_ssh_keys.yaml << EOF
# SSH AI 服务器配置文件 - SSH Keys 测试

# 服务器配置
server:
  port: "2214"
  welcome_message: "Hello! 欢迎使用SSHAI SSH Keys测试！"
  prompt_template: "%s@sshai-keys.top> "

# 认证配置 - 启用密码认证和SSH公钥认证
auth:
  password: "test123"  # 设置密码以启用SSH公钥认证
  login_prompt: "请输入访问密码: "
  # SSH公钥免密登录配置
  authorized_keys:
    - "$RSA_PUBLIC_KEY"
    - "$ED25519_PUBLIC_KEY"
  authorized_keys_file: ""

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
  host_key_file: "host_key_ssh_keys.pem"

# 国际化配置
i18n:
  language: "zh-cn"

# AI提示词配置
prompt:
  system_prompt: "你是一个专业的AI助手，正在测试SSH Keys免密登录功能。"
  user_prompt: "用户问题："
  assistant_prompt: "AI助手："
  stdin_prompt: "请分析以下内容："
  exec_prompt: ""
EOF

echo "   - 配置文件已创建: config_ssh_keys.yaml"
echo "   - 端口: 2214"
echo "   - 密码: test123"
echo "   - 已配置 2 个SSH公钥"

# 创建authorized_keys文件测试
echo "3. 创建authorized_keys文件测试..."
mkdir -p .ssh
cat > .ssh/authorized_keys << EOF
# SSH Keys for SSHAI testing
$RSA_PUBLIC_KEY
$ED25519_PUBLIC_KEY
EOF

cat > config_ssh_keys_file.yaml << EOF
# SSH AI 服务器配置文件 - SSH Keys 文件测试

# 服务器配置
server:
  port: "2215"
  welcome_message: "Hello! 欢迎使用SSHAI SSH Keys文件测试！"
  prompt_template: "%s@sshai-keys-file.top> "

# 认证配置 - 使用authorized_keys文件
auth:
  password: "test456"
  login_prompt: "请输入访问密码: "
  authorized_keys: []
  authorized_keys_file: ".ssh/authorized_keys"

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
  host_key_file: "host_key_ssh_keys_file.pem"

# 国际化配置
i18n:
  language: "zh-cn"

# AI提示词配置
prompt:
  system_prompt: "你是一个专业的AI助手，正在测试SSH Keys文件免密登录功能。"
  user_prompt: "用户问题："
  assistant_prompt: "AI助手："
  stdin_prompt: "请分析以下内容："
  exec_prompt: ""
EOF

echo "   - authorized_keys文件已创建: .ssh/authorized_keys"
echo "   - 文件配置已创建: config_ssh_keys_file.yaml"
echo "   - 端口: 2215"

# 创建测试脚本
cat > test_connections.sh << 'EOF'
#!/bin/bash

echo "=== SSH Keys 连接测试 ==="

echo "测试说明："
echo "1. 配置列表模式 (端口2214) - 公钥直接配置在config中"
echo "2. 配置文件模式 (端口2215) - 公钥从authorized_keys文件读取"
echo ""

echo "连接测试命令："
echo "# RSA密钥连接测试 (端口2214)"
echo "ssh -i test_key -p 2214 testuser@localhost"
echo ""
echo "# Ed25519密钥连接测试 (端口2214)"
echo "ssh -i test_key_ed25519 -p 2214 testuser@localhost"
echo ""
echo "# RSA密钥连接测试 (端口2215 - 文件模式)"
echo "ssh -i test_key -p 2215 testuser@localhost"
echo ""
echo "# 密码认证测试 (端口2214)"
echo "ssh -p 2214 testuser@localhost"
echo "# 输入密码: test123"
echo ""

echo "启动服务器命令："
echo "# 启动配置列表模式服务器"
echo "cd .. && go run cmd/main.go -c $PWD/config_ssh_keys.yaml"
echo ""
echo "# 启动配置文件模式服务器"
echo "cd .. && go run cmd/main.go -c $PWD/config_ssh_keys_file.yaml"
EOF

chmod +x test_connections.sh

echo "4. 测试完成！"
echo ""
echo "测试文件已创建在目录: $TEST_DIR/"
echo "- test_key, test_key.pub (RSA密钥对)"
echo "- test_key_ed25519, test_key_ed25519.pub (Ed25519密钥对)"
echo "- config_ssh_keys.yaml (配置列表模式)"
echo "- config_ssh_keys_file.yaml (配置文件模式)"
echo "- .ssh/authorized_keys (公钥文件)"
echo "- test_connections.sh (连接测试脚本)"
echo ""
echo "使用方法："
echo "1. 启动服务器: cd .. && go run cmd/main.go -c $TEST_DIR/config_ssh_keys.yaml"
echo "2. 测试连接: cd $TEST_DIR && ./test_connections.sh"
echo ""
echo "SSH Keys 免密登录功能测试准备完成！"

cd ..