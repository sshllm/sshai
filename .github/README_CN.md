# GitHub Actions 工作流说明

本目录包含了SSHAI项目的GitHub Actions自动化工作流配置。

## 📁 工作流文件

### 1. `release.yml` - 自动发布工作流
**触发条件：**
- 推送标签（格式：`v*`，如 `v1.0.0`）
- 手动触发（可指定版本号）

**功能：**
- 多平台交叉编译（Linux、macOS、Windows）
- 自动创建发布包（ZIP格式）
- 生成详细的发布说明
- 自动上传到GitHub Releases

**使用方法：**
```bash
# 方法1：推送标签触发
git tag v1.0.0
git push origin v1.0.0

# 方法2：在GitHub网页上手动触发
# Actions -> Release -> Run workflow -> 输入版本号
```

### 2. `ci.yml` - 持续集成工作流
**触发条件：**
- 推送到 `main` 或 `develop` 分支
- 向 `main` 分支提交Pull Request

**功能：**
- 代码编译测试
- 多平台构建验证
- 依赖项检查
- 构建产物上传（保留7天）

### 3. `manual-release.yml` - 手动发布工作流
**触发条件：**
- 仅手动触发

**功能：**
- 灵活的版本控制
- 可选择预发布或草稿模式
- 完整的多平台构建
- SHA256校验和生成

**使用方法：**
1. 进入GitHub仓库的Actions页面
2. 选择"Manual Release"工作流
3. 点击"Run workflow"
4. 填写参数：
   - **Version**: 版本号（如 `v1.0.0`）
   - **Pre-release**: 是否为预发布版本
   - **Draft**: 是否创建为草稿

## 🚀 发布流程

### 自动发布（推荐）
1. 确保代码已合并到main分支
2. 创建并推送版本标签：
   ```bash
   git checkout main
   git pull origin main
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
3. GitHub Actions自动构建并发布

### 手动发布
1. 进入GitHub仓库
2. 点击"Actions"标签
3. 选择"Manual Release"
4. 点击"Run workflow"
5. 输入版本信息并执行

## 📦 构建产物

每次发布会生成以下文件：
- `sshai_linux_amd64_v1.0.0.zip` - Linux AMD64版本
- `sshai_linux_arm64_v1.0.0.zip` - Linux ARM64版本
- `sshai_darwin_amd64_v1.0.0.zip` - macOS Intel版本
- `sshai_darwin_arm64_v1.0.0.zip` - macOS Apple Silicon版本
- `sshai_windows_amd64_v1.0.0.zip` - Windows AMD64版本
- `checksums.txt` - SHA256校验和文件

每个ZIP包包含：
- 对应平台的可执行文件
- `config.yaml.example` 配置文件
- `README.txt` 使用说明

## 🔧 配置要求

### 必需的GitHub Secrets
- `GITHUB_TOKEN` - 自动提供，用于创建releases

### 可选配置
如需要额外功能，可以添加以下secrets：
- `SLACK_WEBHOOK` - Slack通知（需要修改工作流）
- `DISCORD_WEBHOOK` - Discord通知（需要修改工作流）

## 📋 版本号规范

版本号必须遵循语义化版本规范：
- 格式：`v主版本.次版本.修订版本`
- 示例：`v1.0.0`, `v1.2.3`, `v2.0.0-beta`
- 预发布：`v1.0.0-alpha`, `v1.0.0-beta`, `v1.0.0-rc1`

## 🐛 故障排除

### 构建失败
1. 检查Go版本兼容性
2. 确认所有依赖项可用
3. 检查语言包文件是否存在

### 发布失败
1. 确认版本号格式正确
2. 检查标签是否已存在
3. 验证GitHub token权限

### 常见问题
- **语言包文件缺失**：工作流会自动创建占位符文件
- **版本号冲突**：使用不同的版本号或删除现有标签
- **权限问题**：确保仓库设置允许GitHub Actions创建releases

## 📚 更多信息

- [GitHub Actions文档](https://docs.github.com/en/actions)
- [Go交叉编译指南](https://golang.org/doc/install/source#environment)
- [语义化版本规范](https://semver.org/lang/zh-CN/)