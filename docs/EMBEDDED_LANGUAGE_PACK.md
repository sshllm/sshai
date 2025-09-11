# 嵌入式语言包系统

## 概述

基于用户建议，我们已经将语言包系统从外部文件加载优化为嵌入式系统，使用Go的`embed`包将语言包直接嵌入到二进制文件中。这样部署时只需要二进制文件和`config.yaml`配置文件。

## 优化前后对比

### 优化前（外部语言包）
```
部署文件:
├── sshai                 # 二进制文件
├── config.yaml          # 配置文件
└── lang/                 # 语言包目录
    ├── lang-zh-cn.yaml   # 中文语言包
    └── lang-en-us.yaml   # 英文语言包
```

### 优化后（嵌入式语言包）
```
部署文件:
├── sshai                 # 二进制文件（包含嵌入的语言包）
└── config.yaml          # 配置文件
```

## 技术实现

### 1. 使用Go embed包
```go
//go:embed lang/*.yaml
var langFS embed.FS
```

### 2. 语言包文件位置
语言包文件现在位于：
- `pkg/i18n/lang/lang-zh-cn.yaml`
- `pkg/i18n/lang/lang-en-us.yaml`

### 3. 动态加载机制
```go
// 从嵌入文件系统加载指定语言
func (i *I18n) loadLanguage(lang Language) error {
    filename := fmt.Sprintf("lang/lang-%s.yaml", string(lang))
    
    // 从嵌入文件系统读取文件
    data, err := langFS.ReadFile(filename)
    if err != nil {
        return fmt.Errorf("failed to read embedded language pack file %s: %v", filename, err)
    }
    
    // 解析YAML...
}
```

## 构建和部署

### 1. 嵌入式构建脚本
使用 `scripts/build_embedded.sh` 进行构建：

```bash
./scripts/build_embedded.sh
```

该脚本会：
- 构建多平台二进制文件
- 自动嵌入语言包
- 创建发布包
- 生成部署说明

### 2. 支持的平台
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### 3. 部署步骤
1. 选择对应平台的二进制文件
2. 复制二进制文件和 `config.yaml` 到目标服务器
3. 根据需要修改配置文件
4. 运行程序

## 功能验证

### 测试脚本
使用 `scripts/test_embedded.sh` 验证功能：

```bash
./scripts/test_embedded.sh
```

### 测试结果
✅ **二进制文件大小**: 6.0M（包含嵌入的语言包）
✅ **语言包嵌入**: 成功嵌入中英文语言包
✅ **无外部依赖**: 无需外部语言文件
✅ **多语言切换**: 通过config.yaml正常切换
✅ **程序启动**: 中英文界面正常显示

## 配置说明

在 `config.yaml` 中配置语言：

```yaml
# 界面语言配置
language: zh-cn  # 中文
# language: en-us  # 英文
```

支持的语言代码：
- `zh-cn`: 简体中文
- `en-us`: 英文

## 优势

### 1. 部署简化
- **优化前**: 需要复制3个文件/目录
- **优化后**: 只需要复制2个文件

### 2. 文件管理
- 无需担心语言文件丢失
- 版本一致性保证
- 减少部署错误

### 3. 性能优化
- 语言包直接从内存加载
- 无文件系统I/O开销
- 启动速度更快

### 4. 安全性
- 语言包无法被外部修改
- 减少文件权限问题
- 降低安全风险

## 开发说明

### 修改语言包
1. 编辑 `pkg/i18n/lang/lang-zh-cn.yaml` 或 `pkg/i18n/lang/lang-en-us.yaml`
2. 重新构建二进制文件
3. 语言包会自动嵌入到新的二进制文件中

### 添加新语言
1. 在 `pkg/i18n/lang/` 目录下创建新的语言文件
2. 在 `pkg/i18n/i18n.go` 中添加新的语言常量
3. 更新 `GetAvailableLanguages()` 函数
4. 重新构建

## 兼容性

- **Go版本**: 需要Go 1.16+（embed包支持）
- **操作系统**: 支持所有Go支持的平台
- **配置文件**: 与之前版本完全兼容

## 总结

嵌入式语言包系统成功实现了用户的需求：
- ✅ 简化部署流程
- ✅ 减少文件依赖
- ✅ 保持功能完整性
- ✅ 提升用户体验

现在用户只需要二进制文件和配置文件即可完成部署，大大简化了运维工作。