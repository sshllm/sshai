# MCP 使用指南

## 什么是 MCP？

Model Context Protocol (MCP) 是一个开放标准，允许AI模型与外部工具和数据源进行安全、受控的交互。通过MCP，SSHAI可以访问文件系统、网络资源、时间信息等各种外部服务。

## 快速开始

### 1. 启用MCP功能

在 `config.yaml` 中启用MCP：

```yaml
mcp:
  enabled: true
  refresh_interval: 300  # 工具列表刷新间隔（秒）
  servers: []  # 服务器列表
```

### 2. 配置MCP服务器

#### 使用uvx（推荐）

uvx是Python的包管理器，启动速度快，稳定性好：

```yaml
mcp:
  enabled: true
  servers:
    # 时间服务器
    - name: "time"
      transport: "stdio"
      command: ["uvx", "mcp-server-time"]
      enabled: true
    
    # 网络请求服务器
    - name: "fetch"
      transport: "stdio"
      command: ["uvx", "mcp-server-fetch"]
      enabled: true
    
    # 文件系统服务器
    - name: "filesystem"
      transport: "stdio"
      command: ["uvx", "mcp-server-filesystem", "/tmp"]
      enabled: true
```

#### 使用npx

npx是Node.js的包管理器，支持更多的MCP服务器：

```yaml
mcp:
  enabled: true
  servers:
    # 必应搜索服务器
    - name: "bing"
      transport: "stdio"
      command: ["npx", "bing-cn-mcp"]
      enabled: true
    
    # GitHub服务器
    - name: "github"
      transport: "stdio"
      command: ["npx", "@modelcontextprotocol/server-github"]
      enabled: true
      env:
        GITHUB_PERSONAL_ACCESS_TOKEN: "your-github-token"
```

### 3. 使用预安装的服务器

如果你已经预安装了MCP服务器，可以直接使用：

```yaml
mcp:
  enabled: true
  servers:
    - name: "local-server"
      transport: "stdio"
      command: ["/usr/local/bin/my-mcp-server"]
      enabled: true
```

## 常用MCP服务器

### 官方服务器

| 服务器名称 | 包管理器 | 命令 | 功能描述 |
|-----------|---------|------|----------|
| mcp-server-time | uvx | `["uvx", "mcp-server-time"]` | 获取当前时间和日期 |
| mcp-server-fetch | uvx | `["uvx", "mcp-server-fetch"]` | 发送HTTP请求 |
| mcp-server-filesystem | uvx | `["uvx", "mcp-server-filesystem", "/path"]` | 文件系统操作 |
| server-github | npx | `["npx", "@modelcontextprotocol/server-github"]` | GitHub API访问 |
| server-sqlite | npx | `["npx", "@modelcontextprotocol/server-sqlite"]` | SQLite数据库操作 |

### 第三方服务器

| 服务器名称 | 包管理器 | 命令 | 功能描述 |
|-----------|---------|------|----------|
| bing-cn-mcp | npx | `["npx", "bing-cn-mcp"]` | 必应搜索（中文优化） |
| mcp-server-docker | npx | `["npx", "mcp-server-docker"]` | Docker容器管理 |
| mcp-server-postgres | npx | `["npx", "mcp-server-postgres"]` | PostgreSQL数据库操作 |

## 高级配置

### 环境变量配置

某些MCP服务器需要环境变量：

```yaml
mcp:
  enabled: true
  servers:
    - name: "github"
      transport: "stdio"
      command: ["npx", "@modelcontextprotocol/server-github"]
      enabled: true
      env:
        GITHUB_PERSONAL_ACCESS_TOKEN: "ghp_xxxxxxxxxxxx"
        GITHUB_API_URL: "https://api.github.com"
    
    - name: "postgres"
      transport: "stdio"
      command: ["npx", "mcp-server-postgres"]
      enabled: true
      env:
        DATABASE_URL: "postgresql://user:pass@localhost:5432/dbname"
```

### 工作目录配置

为MCP服务器指定工作目录：

```yaml
mcp:
  enabled: true
  servers:
    - name: "filesystem"
      transport: "stdio"
      command: ["uvx", "mcp-server-filesystem", "."]
      enabled: true
      working_dir: "/home/user/projects"
```

### 超时和重试配置

针对不稳定的服务器配置超时和重试：

```yaml
mcp:
  enabled: true
  connection_timeout: 15  # 连接超时（秒）
  max_retries: 3         # 最大重试次数
  servers:
    - name: "slow-server"
      transport: "stdio"
      command: ["npx", "some-slow-mcp-server"]
      enabled: true
```

## 使用示例

### 1. 时间查询

配置时间服务器后，你可以询问：

```bash
ssh user@localhost -p 2213 "现在几点了？"
ssh user@localhost -p 2213 "今天是星期几？"
```

### 2. 网络请求

配置fetch服务器后，你可以：

```bash
ssh user@localhost -p 2213 "帮我获取 https://api.github.com/users/octocat 的信息"
ssh user@localhost -p 2213 "检查 https://www.google.com 是否可访问"
```

### 3. 文件操作

配置文件系统服务器后，你可以：

```bash
ssh user@localhost -p 2213 "列出 /tmp 目录下的文件"
ssh user@localhost -p 2213 "读取 /etc/hosts 文件的内容"
```

### 4. 搜索功能

配置必应搜索服务器后，你可以：

```bash
ssh user@localhost -p 2213 "搜索最新的AI技术发展"
ssh user@localhost -p 2213 "查找Python异步编程的教程"
```

## 故障排除

### 常见问题

1. **服务器启动失败**
   - 检查命令是否正确
   - 确认包管理器已安装
   - 查看错误日志

2. **连接超时**
   - 增加 `connection_timeout` 值
   - 检查网络连接
   - 尝试手动运行命令

3. **权限问题**
   - 确认命令有执行权限
   - 检查工作目录权限
   - 验证环境变量设置

### 调试模式

启用详细日志来调试MCP连接：

```yaml
mcp:
  enabled: true
  debug: true  # 启用调试日志
  servers:
    # ... 你的服务器配置
```

### NPX特殊问题

如果使用npx遇到问题，请参考 [MCP NPX故障排除指南](MCP_NPX_TROUBLESHOOTING.md)。

## 安全考虑

### 权限控制

- 文件系统服务器应限制访问路径
- 网络服务器应配置允许的域名
- 数据库服务器应使用只读用户

### 环境变量安全

- 不要在配置文件中硬编码敏感信息
- 使用环境变量或密钥文件
- 定期轮换API密钥

```yaml
# 推荐做法
mcp:
  servers:
    - name: "github"
      command: ["npx", "@modelcontextprotocol/server-github"]
      env:
        GITHUB_PERSONAL_ACCESS_TOKEN: "${GITHUB_TOKEN}"  # 从环境变量读取
```

## 开发自定义MCP服务器

如果现有的MCP服务器不满足需求，你可以开发自定义服务器：

### Python示例

```python
from mcp.server import Server
from mcp.types import Tool

server = Server("my-custom-server")

@server.tool("hello")
async def hello_tool(name: str) -> str:
    return f"Hello, {name}!"

if __name__ == "__main__":
    server.run()
```

### 配置自定义服务器

```yaml
mcp:
  servers:
    - name: "custom"
      transport: "stdio"
      command: ["python", "/path/to/my-server.py"]
      enabled: true
```

## 最佳实践

1. **优先使用uvx** - 更稳定，启动更快
2. **合理配置超时** - 避免长时间等待
3. **定期更新服务器** - 获取最新功能和修复
4. **监控资源使用** - 避免服务器占用过多资源
5. **备份重要配置** - 防止配置丢失

## 更多资源

- [MCP官方文档](https://modelcontextprotocol.io/)
- [MCP服务器列表](https://github.com/modelcontextprotocol/servers)
- [MCP开发指南](https://modelcontextprotocol.io/docs/building-servers)