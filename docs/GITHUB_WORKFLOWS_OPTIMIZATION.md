# GitHub Workflows 优化文档

## 🎯 优化概述

本次优化确保GitHub Actions workflows与Makefile使用完全相同的版本信息注入机制，实现构建一致性和版本信息的准确性。

## 🔧 主要改进

### 1. 版本信息统一化

所有workflows现在使用与Makefile相同的版本信息获取方式：

```bash
# 版本信息变量（与Makefile完全一致）
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "v0.9.19")
GIT_COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_VERSION=$(go version | cut -d' ' -f3)
```

### 2. LDFLAGS 标准化

所有构建步骤现在使用统一的LDFLAGS格式：

```bash
LDFLAGS="-X 'sshai/pkg/version.Version=${VERSION}' \
         -X 'sshai/pkg/version.GitCommit=${GIT_COMMIT}' \
         -X 'sshai/pkg/version.BuildTime=${BUILD_TIME}' \
         -X 'sshai/pkg/version.GoVersion=${GO_VERSION}' \
         -s -w"
```

## 📋 优化的Workflows

### 1. CI Workflow (`.github/workflows/ci.yml`)

**优化内容：**
- ✅ 添加版本信息设置步骤
- ✅ 统一LDFLAGS格式
- ✅ 版本信息输出到构建日志
- ✅ 多平台构建版本信息一致性

**主要变更：**
```yaml
- name: Set version variables
  id: version
  run: |
    VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "v0.9.19-ci")
    GIT_COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
    BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    GO_VERSION=$(go version | cut -d' ' -f3)
    
    echo "VERSION=${VERSION}" >> $GITHUB_OUTPUT
    echo "GIT_COMMIT=${GIT_COMMIT}" >> $GITHUB_OUTPUT
    echo "BUILD_TIME=${BUILD_TIME}" >> $GITHUB_OUTPUT
    echo "GO_VERSION=${GO_VERSION}" >> $GITHUB_OUTPUT
```

### 2. Release Workflow (`.github/workflows/release.yml`)

**优化内容：**
- ✅ 完整版本信息注入
- ✅ 构建信息包含在README中
- ✅ 多平台构建版本一致性
- ✅ 发布包包含详细构建信息

**主要变更：**
- 版本信息获取逻辑优化
- 构建LDFLAGS与Makefile完全一致
- README文件包含完整构建信息

### 3. Manual Release Workflow (`.github/workflows/manual-release.yml`)

**优化内容：**
- ✅ 手动发布版本信息注入
- ✅ 构建信息验证
- ✅ 版本格式验证
- ✅ 完整的构建元数据

### 4. Version Test Workflow (`.github/workflows/version-test.yml`) - 新增

**功能特性：**
- 🧪 版本信息注入测试
- 🔍 与Makefile构建对比
- 🎨 UI Banner版本信息测试
- ✅ 构建一致性验证

## 🚀 使用方法

### 本地开发构建
```bash
# 开发版本（包含调试信息）
make build

# 生产版本（优化大小，去除调试信息）
make build-release

# 查看版本信息
make version
```

### GitHub Actions构建

#### 1. 持续集成
- 推送到 `main` 或 `develop` 分支自动触发
- 拉取请求到 `main` 分支自动触发
- 构建多平台二进制文件
- 版本信息自动注入

#### 2. 自动发布
- 推送标签（如 `v1.0.0`）自动触发
- 构建所有平台的发布包
- 自动创建GitHub Release
- 包含完整的版本和构建信息

#### 3. 手动发布
- 通过GitHub Actions界面手动触发
- 可指定版本号和发布选项
- 支持预发布和草稿模式

#### 4. 版本测试
- 手动触发或代码变更时自动运行
- 验证版本信息注入正确性
- 测试UI Banner显示效果

## 📊 版本信息展示

构建后的二进制文件将包含以下版本信息：

```
Version: v0.9.19-1-g1234567
Git Commit: 1234567890abcdef...
Build Time: 2025-09-18T08:00:00Z
Go Version: go1.22.0
Platform: linux/amd64
```

### UI Banner显示效果

```
  .-')     .-')    ('-. .-.         ('-.              
 ( OO ).  ( OO ). ( OO )  /        ( OO ).-.          
(_)---\_)(_)---\_),--. ,--.        / . --. /  ,-.-')  
/    _ | /    _ | |  | |  |        | \-.  \   |  |OO) 
\  :` `. \  :` `. |   .|  |      .-'-'  |  |  |  |  \ 
 '..`''.) '..`''.)|       |       \| |_.'  |  |  |(_/ 
.-._)   \.-._)   \|  .-.  |        |  .-.  | ,|  |_.' 
\       /\       /|  | |  |        |  | |  |(_|  |    
 `-----'  `-----' `--' `--'        `--' `--'  `--'    

🚀 SSH AI Assistant v0.9.19-1-g1234567

⚡ AI-Powered   ⚡ Real-time   🔒 Secure

🌍 Website: https://sshllm.top
📦 GitHub:  https://github.com/sshllm/sshai
🔨 Built:   2025-09-18 08:00:00 UTC

👨‍💻 Built for modern developers & teams
```

## 🔍 验证方法

### 1. 本地验证
```bash
# 构建并检查版本
make build
./sshai --version  # 如果支持版本参数

# 或者通过代码测试
go run test_version.go
```

### 2. CI/CD验证
- 查看GitHub Actions构建日志
- 检查构建产物的版本信息
- 运行版本测试workflow

### 3. 发布验证
- 检查Release页面的构建信息
- 下载并验证二进制文件版本
- 确认README文件包含正确信息

## 🎯 最佳实践

### 1. 版本标签规范
```bash
# 正式版本
git tag v1.0.0
git push origin v1.0.0

# 预发布版本
git tag v1.0.0-beta.1
git push origin v1.0.0-beta.1

# 开发版本
git tag v1.0.0-dev.20250918
git push origin v1.0.0-dev.20250918
```

### 2. 构建环境一致性
- 本地开发使用 `make build`
- CI/CD使用相同的LDFLAGS
- 发布使用 `make build-release`

### 3. 版本信息检查
- 每次发布前运行版本测试
- 确认UI Banner显示正确
- 验证构建信息完整性

## 🐛 故障排除

### 常见问题

1. **版本信息为空或默认值**
   - 检查git仓库状态
   - 确认LDFLAGS格式正确
   - 验证构建命令

2. **构建时间格式错误**
   - 确认使用UTC时间格式
   - 检查date命令兼容性

3. **Git信息获取失败**
   - 确认git仓库完整性
   - 检查CI环境git配置

### 调试命令
```bash
# 检查版本变量
make version

# 测试构建
make build-release

# 验证版本注入
go run test_version.go
```

## 📚 相关文档

- [Makefile使用指南](../Makefile)
- [版本管理模块](../pkg/version/version.go)
- [UI系统文档](UI_BANNER_OPTIMIZATION.md)
- [SSH Keys功能](SSH_KEYS_GUIDE.md)

---

**更新时间：** 2025-09-18  
**版本：** v1.0.0  
**维护者：** SSHAI开发团队