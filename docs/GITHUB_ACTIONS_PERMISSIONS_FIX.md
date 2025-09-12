# GitHub Actions 权限问题修复指南

## 🚨 问题描述

GitHub Actions执行发布工作流时出现403权限错误：

```
⚠️ GitHub release failed with status: 403
undefined
retrying... (2 retries remaining)
❌ Too many retries. Aborting...
Error: Too many retries.
```

## 🔍 问题分析

403错误通常表示权限不足，主要原因包括：

1. **工作流权限配置缺失** - GitHub Actions默认权限策略变更
2. **softprops/action-gh-release版本过旧** - v1版本兼容性问题
3. **token配置方式过时** - 新版本要求使用`token`而非`env.GITHUB_TOKEN`

## ✅ 修复方案

### 1. 添加工作流权限配置

为所有工作流文件添加明确的权限声明：

```yaml
permissions:
  contents: write    # 允许创建releases和读写仓库内容
  packages: write    # 允许发布包（如果需要）
  actions: read      # 允许读取actions状态
```

### 2. 升级action版本

将`softprops/action-gh-release`从v1升级到v2：

```yaml
# 旧版本
- uses: softprops/action-gh-release@v1

# 新版本
- uses: softprops/action-gh-release@v2
```

### 3. 更新token配置方式

从环境变量方式改为直接参数方式：

```yaml
# 旧方式
- uses: softprops/action-gh-release@v1
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

# 新方式
- uses: softprops/action-gh-release@v2
  with:
    token: ${{ secrets.GITHUB_TOKEN }}
```

## 📋 修复的文件

### 1. `.github/workflows/release.yml`
- ✅ 添加`permissions`配置
- ✅ 升级`softprops/action-gh-release@v1` → `@v2`
- ✅ 更新token配置方式

### 2. `.github/workflows/ci.yml`
- ✅ 添加`permissions`配置（只读权限）

### 3. `.github/workflows/manual-release.yml`
- ✅ 添加`permissions`配置
- ✅ 升级`softprops/action-gh-release@v1` → `@v2`
- ✅ 更新token配置方式

## 🔧 权限配置详解

### 发布工作流权限
```yaml
permissions:
  contents: write    # 创建releases、读写代码
  packages: write    # 发布包到GitHub Packages
  actions: read      # 读取workflow状态
```

### CI工作流权限
```yaml
permissions:
  contents: read     # 只读代码权限
  actions: read      # 读取workflow状态
```

## 🚀 验证修复

### 1. 检查仓库设置
确保仓库设置允许GitHub Actions：
- 访问 `Settings` → `Actions` → `General`
- 确认"Actions permissions"设置为"Allow all actions and reusable workflows"

### 2. 检查工作流权限
确保仓库设置中的工作流权限正确：
- 访问 `Settings` → `Actions` → `General`
- 在"Workflow permissions"部分选择"Read and write permissions"

### 3. 测试发布流程
```bash
# 创建测试标签
git tag v0.9.13-test
git push origin v0.9.13-test

# 或使用手动发布
# 访问Actions页面 → Manual Release → Run workflow
```

## 🛡️ 安全考虑

### 最小权限原则
- CI工作流只需要`read`权限
- 发布工作流需要`write`权限创建releases
- 避免给予不必要的权限

### Token安全
- 使用内置的`GITHUB_TOKEN`，自动管理权限
- 避免创建额外的Personal Access Token
- 定期检查权限使用情况

## 📚 相关文档

- [GitHub Actions Permissions](https://docs.github.com/en/actions/security-guides/automatic-token-authentication)
- [softprops/action-gh-release](https://github.com/softprops/action-gh-release)
- [Workflow Permissions](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#permissions)

## 🔄 回滚方案

如果修复后仍有问题，可以临时回滚：

```yaml
# 临时回滚到v1版本（不推荐长期使用）
- uses: softprops/action-gh-release@v1
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

但建议优先解决权限配置问题，而不是回滚版本。

## ✅ 修复确认清单

- [x] 添加工作流权限配置
- [x] 升级softprops/action-gh-release到v2
- [x] 更新token配置方式
- [x] 验证仓库权限设置
- [x] 测试发布流程
- [x] 更新文档

**修复时间**: 2025年9月12日  
**影响范围**: 所有GitHub Actions发布工作流  
**风险等级**: 低（向后兼容，仅修复权限问题）