# 项目结构说明

## 文件夹组织

### 根目录文件
- `README.md` - 项目主要说明文档
- `config.yaml` - 主配置文件
- `go.mod` / `go.sum` - Go模块依赖管理
- `Makefile` - 构建和管理脚本
- `sshai` - 编译后的可执行文件

### 核心代码 (`cmd/` 和 `pkg/`)
```
cmd/
└── main.go              # 程序入口点

pkg/
├── config/              # 配置管理
│   └── config.go
├── models/              # 数据模型定义
│   └── models.go
├── ai/                  # AI功能模块
│   ├── assistant.go     # AI助手核心功能
│   └── models.go        # 模型管理
├── ssh/                 # SSH服务器
│   ├── server.go        # SSH服务器实现
│   └── session.go       # 会话处理
└── utils/               # 工具函数
    └── text.go          # 文本处理工具
```

### 文档目录 (`docs/`)
- `readme.md` - 原始项目需求文档
- `CONFIG_GUIDE.md` - 配置文件详细说明
- `USAGE.md` - 使用指南
- `MODULAR_ARCHITECTURE.md` - 模块化架构设计
- `PROJECT_STRUCTURE.md` - 本文件，项目结构说明
- 其他开发过程文档...

### 脚本目录 (`scripts/`)
- `run.sh` - 程序运行脚本
- `test.sh` - 基础测试脚本
- `test_*.sh` - 各种功能测试脚本
- `main.go.bak` - 原始单文件版本备份

### 密钥目录 (`keys/`)
- `host_key.pem` - SSH服务器主机密钥

## 设计原则

1. **关注点分离** - 不同类型的文件分别存放
2. **模块化** - 代码按功能模块组织
3. **文档集中** - 所有文档统一管理
4. **脚本独立** - 测试和运行脚本单独存放
5. **安全文件隔离** - 密钥文件独立目录管理

## 文件移动记录

从原始的扁平结构重组为分层结构：
- 所有 `*.md` 文件 → `docs/`
- 所有 `test*.sh` 和 `run.sh` → `scripts/`
- `host_key.pem` → `keys/`
- `main.go.bak` → `scripts/`

## 配置更新

- 更新了 `config.yaml` 中的 `host_key_file` 路径为 `keys/host_key.pem`
- 保持其他配置不变，确保程序正常运行