# MCP NPX 故障排除指南

## 问题描述

使用npx启动MCP服务时可能会遇到卡住的情况，而使用uvx启动相同的服务却没有问题。

## 常见原因和解决方案

### 1. NPX 首次下载包导致的延迟

**问题**: npx首次运行时需要下载包，可能需要较长时间
**解决方案**:
```bash
# 预先安装包到本地缓存
npx bing-cn-mcp --help

# 或者全局安装包
npm install -g bing-cn-mcp
```

### 2. NPX 交互式提示

**问题**: npx可能会显示交互式提示，导致进程挂起
**解决方案**: SSHAI已自动设置以下环境变量来避免交互：
- `NPM_CONFIG_YES=true`
- `NPM_CONFIG_AUDIT=false`
- `NPM_CONFIG_FUND=false`
- `NPM_CONFIG_UPDATE_NOTIFIER=false`
- `CI=true`

### 3. 网络连接问题

**问题**: npm registry连接缓慢或失败
**解决方案**:
```bash
# 使用国内镜像
npm config set registry https://registry.npmmirror.com/

# 或者临时使用镜像
npx --registry https://registry.npmmirror.com/ bing-cn-mcp
```

### 4. 权限问题

**问题**: npm缓存目录权限不足
**解决方案**:
```bash
# 修复npm权限
sudo chown -R $(whoami) ~/.npm
sudo chown -R $(whoami) ~/.config
```

## 改进措施

SSHAI v1.1+ 包含以下改进来解决npx问题：

### 1. 连接超时和重试机制
- 首次连接超时：15秒
- 重试连接超时：20秒、25秒
- 最多重试3次

### 2. 预检查机制
- 启动前检查命令是否存在
- 对npx包进行可用性预检查

### 3. 环境变量优化
- 自动设置非交互式环境变量
- 继承系统PATH和相关环境变量

### 4. 详细的错误日志
- 显示连接进度和重试信息
- 提供具体的故障排除建议

## 手动测试方法

### 测试NPX命令
```bash
# 测试命令是否可用
npx --yes --quiet bing-cn-mcp --help

# 测试环境变量设置
NPM_CONFIG_YES=true NPM_CONFIG_AUDIT=false npx bing-cn-mcp --help
```

### 测试UVX命令
```bash
# 测试uvx命令
uvx mcp-server-time --help
uvx mcp-server-fetch --help
```

## 配置建议

### 推荐配置（混合使用）
```yaml
mcp:
  enabled: true
  servers:
    # 使用uvx的稳定服务
    - name: "time"
      transport: "stdio"
      command: ["uvx", "mcp-server-time"]
      enabled: true
      
    - name: "fetch"
      transport: "stdio"
      command: ["uvx", "mcp-server-fetch"]
      enabled: true
    
    # 使用npx的特殊服务（如果必需）
    - name: "bing"
      transport: "stdio"
      command: ["npx", "bing-cn-mcp"]
      enabled: true
```

### 备用配置（预安装包）
```yaml
mcp:
  enabled: true
  servers:
    # 预先全局安装包，然后直接调用
    - name: "bing"
      transport: "stdio"
      command: ["bing-cn-mcp"]  # 不使用npx
      enabled: true
```

## 性能对比

| 包管理器 | 首次启动 | 后续启动 | 稳定性 | 推荐度 |
|---------|---------|---------|--------|--------|
| uvx     | 快      | 很快    | 高     | ⭐⭐⭐⭐⭐ |
| npx     | 慢      | 中等    | 中     | ⭐⭐⭐ |
| 直接命令 | 很快    | 很快    | 很高   | ⭐⭐⭐⭐⭐ |

## 如果问题仍然存在

1. **查看SSHAI日志**: 启动SSHAI时观察MCP连接日志
2. **手动测试命令**: 在终端中直接运行MCP命令
3. **检查系统环境**: 确保Node.js、npm、Python等环境正常
4. **使用uvx替代**: 如果可能，优先使用uvx而不是npx
5. **预安装包**: 将常用的MCP包全局安装，避免运行时下载

## 联系支持

如果以上方法都无法解决问题，请提供以下信息：
- 操作系统版本
- Node.js和npm版本
- 具体的错误日志
- 手动运行命令的结果