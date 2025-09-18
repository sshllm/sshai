# SSH Keys 免密登录功能指南

## 功能概述

SSHAI 现在支持 SSH 公钥免密登录功能，允许用户使用 SSH 密钥对进行身份验证。

## 配置方法

### 1. 更新配置文件

在 `config.yaml` 中添加 SSH 公钥配置：

```yaml
auth:
  password: "your_password"  # 必须设置密码才能启用SSH公钥认证
  login_prompt: "请输入访问密码: "
  # 方式一：直接配置公钥列表
  authorized_keys:
    - "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC... user@hostname"
    - "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI... user2@hostname"
  # 方式二：从文件读取公钥
  authorized_keys_file: "~/.ssh/authorized_keys"
```

### 2. 生成 SSH 密钥对

```bash
# 生成 Ed25519 密钥（推荐）
ssh-keygen -t ed25519 -f ~/.ssh/sshai_key

# 生成 RSA 密钥
ssh-keygen -t rsa -b 2048 -f ~/.ssh/sshai_rsa_key
```

### 3. 连接测试

```bash
# 使用私钥连接
ssh -i ~/.ssh/sshai_key -p 2213 username@localhost
```

## 测试功能

运行测试脚本：

```bash
./scripts/test_ssh_keys.sh
```

## 安全说明

- SSH 公钥认证仅在设置密码时启用
- 支持多个公钥同时配置
- 兼容所有标准 SSH 密钥类型