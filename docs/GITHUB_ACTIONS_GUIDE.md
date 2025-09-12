# GitHub Actions 自动化构建和发布指南

本文档详细说明了SSHAI项目的GitHub Actions自动化构建和发布系统的配置和使用方法。

## 📋 概述

SSHAI项目配置了完整的CI/CD流水线，包括：
- 持续集成（CI）
- 自动发布（Release）
- 手动发布（Manual Release）
- 多平台交叉编译
- 自动生成发布包

## 🏗️ 工作流架构

```
GitHub Repository
├── .github/workflows/
│   ├── ci.yml              # 持续集成
│   ├── release.yml         # 自动发布
│   └── manual-release.yml  # 手动发布
├── scripts/
│   ├── build_embedded.sh   # 原始构建脚本
│   └── test_github_actions.sh  # 本地测试脚本
└── docs/
    └── GITHUB_ACTIONS_GUIDE.md  # 本文档
```

## 🔄 工作流详解

### 1. 持续集成工作流 (`ci.yml`)

**触发条件：**
- 推送到 `main` 或 `develop` 分支
- 向 `main` 分支提交 Pull Request

**执行步骤：**
1. 检出代码
2. 设置Go环境（1.21版本）
3. 缓存Go模块
4. 下载并验证依赖
5. 检查语言包文件（自动创建占位符）
6. 运行测试（如果存在）
7. 多平台构建验证
8. 上传构建产物（保留7天）

**作用：**
- 确保代码质量
- 验证多平台兼容性
- 早期发现构建问题

### 2. 自动发布工作流 (`release.yml`)

**触发条件：**
- 推送版本标签（格式：`v*`）
- 手动触发（可指定版本）

**执行步骤：**
1. 多平台并行构建
2. 创建发布包
3. 生成发布说明
4. 创建GitHub Release
5. 上传所有构建产物

**支持平台：**
- Linux AMD64/ARM64
- macOS Intel/Apple Silicon
- Windows AMD64

### 3. 手动发布工作流 (`manual-release.yml`)

**触发条件：**
- 仅手动触发

**特殊功能：**
- 版本号验证
- 预发布选项
- 草稿模式
- SHA256校验和生成

## 🚀 使用指南

### 自动发布流程

1. **准备发布**
   ```bash
   # 确保在main分支
   git checkout main
   git pull origin main
   
   # 确保所有更改已提交
   git status
   ```

2. **创建版本标签**
   ```bash
   # 创建标签
   git tag -a v1.0.0 -m "Release v1.0.0"
   
   # 推送标签触发自动发布
   git push origin v1.0.0
   ```

3. **监控构建过程**
   - 访问GitHub仓库的Actions页面
   - 查看"Build and Release"工作流状态
   - 等待所有平台构建完成

4. **验证发布**
   - 检查Releases页面
   - 下载并测试发布包
   - 验证SHA256校验和

### 手动发布流程

1. **访问GitHub Actions**
   - 进入仓库的Actions页面
   - 选择"Manual Release"工作流

2. **配置发布参数**
   - **Version**: 输入版本号（如`v1.0.0`）
   - **Pre-release**: 是否为预发布版本
   - **Draft**: 是否创建为草稿

3. **执行发布**
   - 点击"Run workflow"
   - 等待构建完成
   - 检查发布结果

### 本地测试

在推送到GitHub之前，可以使用本地测试脚本验证构建：

```bash
# 运行本地测试
./scripts/test_github_actions.sh

# 检查测试结果
# - 验证所有平台能否成功构建
# - 检查发布包结构
# - 确认配置文件正确
```

## 📦 发布包结构

每个发布包包含：
```
sshai_linux_amd64_v1.0.0.zip
├── sshai_linux_amd64           # 可执行文件
├── config.yaml                 # 配置文件
└── README.txt                  # 使用说明
```

## 🔧 配置和自定义

### 修改构建平台

编辑工作流文件中的`platforms`数组：
```yaml
platforms=(
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
  # 添加新平台...
)
```

### 自定义构建参数

修改`go build`命令的ldflags：
```yaml
go build \
  -ldflags "-X main.version=${{ steps.version.outputs.VERSION }} -s -w" \
  -o "build/${output_name}" \
  cmd/main.go
```

### 添加构建后处理

在构建步骤后添加自定义处理：
```yaml
- name: Custom post-build processing
  run: |
    # 自定义处理逻辑
    echo "Processing build artifacts..."
```

## 🐛 故障排除

### 常见问题

1. **构建失败**
   - 检查Go版本兼容性
   - 验证依赖项可用性
   - 确认代码语法正确

2. **发布失败**
   - 检查版本号格式（必须是`v1.0.0`格式）
   - 确认标签不存在冲突
   - 验证GitHub token权限

3. **平台特定问题**
   - Windows: 确保文件名包含`.exe`扩展名
   - macOS: 注意Intel和Apple Silicon的区别
   - Linux: 验证ARM64构建环境

### 调试方法

1. **查看构建日志**
   ```bash
   # 在GitHub Actions页面查看详细日志
   # 关注错误信息和警告
   ```

2. **本地复现问题**
   ```bash
   # 使用相同的构建命令
   GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o sshai cmd/main.go
   ```

3. **验证依赖**
   ```bash
   go mod verify
   go mod tidy
   ```

## 📊 监控和维护

### 定期检查

- **每月检查**：
  - Go版本更新
  - GitHub Actions版本更新
  - 依赖项安全更新

- **每次发布后**：
  - 验证所有平台包可用
  - 测试下载和安装流程
  - 收集用户反馈

### 性能优化

1. **缓存优化**
   ```yaml
   - name: Cache Go modules
     uses: actions/cache@v3
     with:
       path: ~/go/pkg/mod
       key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
   ```

2. **并行构建**
   - 利用matrix策略并行构建多平台
   - 减少总体构建时间

3. **构建产物优化**
   - 使用`-s -w`标志减小二进制大小
   - 压缩发布包

## 🔒 安全考虑

### 权限管理
- 使用最小权限原则
- 定期审查GitHub token权限
- 避免在日志中暴露敏感信息

### 代码签名（可选）
```yaml
# 为macOS和Windows添加代码签名
- name: Sign macOS binary
  if: matrix.goos == 'darwin'
  run: |
    # 代码签名逻辑
```

## 📚 参考资源

- [GitHub Actions官方文档](https://docs.github.com/en/actions)
- [Go交叉编译指南](https://golang.org/doc/install/source#environment)
- [语义化版本规范](https://semver.org/)
- [软件发布最佳实践](https://docs.github.com/en/repositories/releasing-projects-on-github/about-releases)

## 🎯 最佳实践

1. **版本管理**
   - 遵循语义化版本规范
   - 使用有意义的标签消息
   - 维护CHANGELOG文件

2. **测试策略**
   - 本地测试后再推送
   - 使用预发布版本测试
   - 收集用户反馈

3. **发布节奏**
   - 定期发布稳定版本
   - 及时修复关键问题
   - 保持向后兼容性

4. **文档维护**
   - 更新发布说明
   - 维护安装指南
   - 记录已知问题

---

通过这套完整的GitHub Actions自动化系统，SSHAI项目可以实现：
- ✅ 自动化多平台构建
- ✅ 一键发布到GitHub Releases
- ✅ 完整的CI/CD流水线
- ✅ 高质量的发布包
- ✅ 详细的发布文档

这大大提高了开发效率，确保了发布质量，为用户提供了便捷的下载和使用体验。