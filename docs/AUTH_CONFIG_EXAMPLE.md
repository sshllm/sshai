# SSH密码认证配置示例

## 功能说明

SSHAI支持可选的SSH密码认证功能：
- **无密码模式**: 用户可以直接连接，无需输入密码
- **密码认证模式**: 用户需要输入正确密码才能访问AI服务

## 配置方式

在 `config.yaml` 文件中的 `auth` 部分进行配置：

### 1. 无密码认证模式 (默认)

```yaml
auth:
  password: ""  # 空字符串表示无需密码
  login_prompt: "请输入访问密码: "
  login_success_message: |  # 无密码模式也会显示此消息
    🎉 欢迎使用 SSHAI v1.0 - SSH AI Assistant
    
    📖 项目地址: https://github.com/your-repo/sshai
    💡 使用说明: 输入消息与AI对话，输入 'exit' 退出
```

### 2. 密码认证模式

```yaml
auth:
  password: "your_secure_password"  # 设置访问密码
  login_prompt: "请输入SSHAI访问密码: "  # 自定义密码提示
  login_success_message: |  # 登录成功后显示的消息
    🎉 认证成功！欢迎使用 SSHAI v1.0
    
    📖 项目地址: https://github.com/your-repo/sshai
    🚀 版本信息: v1.0.0 - SSH AI Assistant
    💡 使用说明: 输入消息与AI对话，输入 'exit' 退出
    
    ════════════════════════════════════════
```

## 配置参数说明

| 参数 | 类型 | 说明 |
|------|------|------|
| `password` | string | 访问密码，空字符串表示无需密码认证 |
| `login_prompt` | string | 密码输入提示文本（SSH客户端显示） |
| `login_success_message` | string | 认证成功后显示的欢迎信息 |

## 使用示例

### 连接到无密码模式的服务器
```bash
ssh -p 2212 username@localhost
# 直接连接，无需输入密码
```

### 连接到密码认证模式的服务器
```bash
ssh -p 2212 username@localhost
# 系统提示: 请输入SSHAI访问密码:
# 输入密码: your_secure_password
# 认证成功后显示自定义欢迎信息
```

## 安全建议

1. **密码强度**: 使用强密码，包含字母、数字和特殊字符
2. **定期更换**: 定期更换访问密码
3. **访问控制**: 结合防火墙规则限制访问来源
4. **日志监控**: 监控SSH连接日志，及时发现异常访问

## 测试方法

使用提供的测试脚本验证认证功能：

```bash
./scripts/test_auth.sh
```

该脚本会自动测试无密码和密码认证两种模式。

## 故障排除

### 常见问题

1. **密码认证失败**
   - 检查 `config.yaml` 中的密码配置
   - 确认输入的密码与配置文件中的密码完全一致
   - 注意密码区分大小写

2. **无法连接服务器**
   - 检查服务器是否正常启动
   - 确认端口配置正确
   - 检查防火墙设置

3. **登录成功消息不显示**
   - 确认 `login_success_message` 配置不为空
   - 检查YAML格式是否正确（注意缩进）

### 日志查看

服务器启动时会显示当前认证模式：
```
SSH服务器配置：无密码认证模式
# 或
SSH服务器配置：密码认证模式
```

认证过程会记录在服务器日志中：
```
密码认证尝试: user=testuser
用户 testuser 认证成功
# 或
用户 testuser 认证失败