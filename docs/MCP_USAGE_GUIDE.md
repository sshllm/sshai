# MCP使用指南

## 快速开始

### 1. 启用MCP功能

在`config.yaml`中设置：

```yaml
mcp:
  enabled: true
  refresh_interval: 300
  servers: []  # 先留空，稍后添加服务器
```

### 2. 安装MCP服务器

#### 文件系统服务器
```bash
npm install -g @modelcontextprotocol/server-filesystem
```

#### SQLite数据库服务器
```bash
npm install -g @modelcontextprotocol/server-sqlite
```

#### Git服务器
```bash
npm install -g @modelcontextprotocol/server-git
```

### 3. 配置MCP服务器

```yaml
mcp:
  enabled: true
  refresh_interval: 300
  servers:
    # 文件系统工具
    - name: "filesystem"
      transport: "stdio"
      command: ["mcp-server-filesystem", "/home/user/documents"]
      enabled: true
    
    # SQLite数据库工具
    - name: "database"
      transport: "stdio"
      command: ["mcp-server-sqlite", "/path/to/database.db"]
      enabled: true
    
    # Git工具
    - name: "git"
      transport: "stdio"
      command: ["mcp-server-git", "/path/to/git/repo"]
      enabled: true
```

## 使用场景

### 场景1：文件操作

**用户**: "请帮我查看当前目录下的所有Python文件"

**AI响应**: 
```
🔧 正在调用工具 list_files...
✅ 工具执行成功: list_files

找到以下Python文件：
- main.py
- utils.py
- config.py
- tests/test_main.py

这些是您目录中的所有Python文件。需要我查看某个文件的内容吗？
```

### 场景2：数据库查询

**用户**: "查询用户表中的所有记录"

**AI响应**:
```
🔧 正在调用工具 execute_query...
✅ 工具执行成功: execute_query

查询结果：
| ID | Name  | Email           | Created    |
|----|-------|-----------------|------------|
| 1  | Alice | alice@email.com | 2024-01-01 |
| 2  | Bob   | bob@email.com   | 2024-01-02 |

共找到2条记录。
```

### 场景3：Git操作

**用户**: "检查当前Git仓库的状态"

**AI响应**:
```
🔧 正在调用工具 git_status...
✅ 工具执行成功: git_status

Git仓库状态：
- 当前分支: main
- 未提交的更改: 3个文件
  - modified: src/main.go
  - new file: docs/guide.md
  - deleted: old_file.txt

建议先提交这些更改再进行其他操作。
```

## 高级配置

### HTTP传输配置

```yaml
mcp:
  servers:
    - name: "web-api"
      transport: "http"
      url: "https://api.example.com/mcp"
      headers:
        Authorization: "Bearer your-api-token"
        Content-Type: "application/json"
        X-Custom-Header: "custom-value"
      enabled: true
```

### SSE传输配置

```yaml
mcp:
  servers:
    - name: "realtime-data"
      transport: "sse"
      url: "https://stream.example.com/mcp/events"
      headers:
        X-API-Key: "your-api-key"
      enabled: true
```

## 故障排除

### 问题1：MCP服务器连接失败

**症状**: 日志显示"连接到MCP服务器失败"

**解决方案**:
1. 检查MCP服务器是否已安装：
   ```bash
   which mcp-server-filesystem
   ```

2. 验证命令路径是否正确：
   ```bash
   mcp-server-filesystem --help
   ```

3. 检查权限设置：
   ```bash
   ls -la /path/to/target/directory
   ```

### 问题2：工具调用超时

**症状**: 工具调用时出现超时错误

**解决方案**:
1. 增加超时时间（在代码中默认为30秒）
2. 检查MCP服务器性能
3. 验证网络连接

### 问题3：工具列表为空

**症状**: AI提示没有可用工具

**解决方案**:
1. 检查MCP配置是否启用：
   ```yaml
   mcp:
     enabled: true
   ```

2. 验证服务器配置：
   ```yaml
   servers:
     - name: "test"
       enabled: true  # 确保启用
   ```

3. 查看启动日志确认连接状态

## 开发自定义MCP服务器

### 基本结构

```javascript
// server.js
import { Server } from '@modelcontextprotocol/sdk/server/index.js';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';

const server = new Server({
  name: "custom-server",
  version: "1.0.0"
});

// 注册工具
server.setRequestHandler('tools/list', async () => ({
  tools: [
    {
      name: "custom_tool",
      description: "A custom tool",
      inputSchema: {
        type: "object",
        properties: {
          input: { type: "string" }
        }
      }
    }
  ]
}));

// 处理工具调用
server.setRequestHandler('tools/call', async (request) => {
  const { name, arguments: args } = request.params;
  
  if (name === "custom_tool") {
    return {
      content: [
        {
          type: "text",
          text: `处理输入: ${args.input}`
        }
      ]
    };
  }
  
  throw new Error(`未知工具: ${name}`);
});

// 启动服务器
const transport = new StdioServerTransport();
await server.connect(transport);
```

### 配置使用

```yaml
mcp:
  servers:
    - name: "custom"
      transport: "stdio"
      command: ["node", "/path/to/server.js"]
      enabled: true
```

## 最佳实践

### 1. 安全考虑
- 限制文件系统访问路径
- 使用只读数据库连接
- 验证所有输入参数
- 设置适当的超时时间

### 2. 性能优化
- 合理设置刷新间隔
- 缓存常用查询结果
- 限制返回数据大小
- 使用连接池

### 3. 错误处理
- 提供清晰的错误信息
- 实现重试机制
- 记录详细日志
- 优雅降级

### 4. 用户体验
- 提供工具使用说明
- 显示执行进度
- 格式化输出结果
- 支持中断操作

## 常用MCP服务器

| 服务器 | 功能 | 安装命令 |
|--------|------|----------|
| filesystem | 文件操作 | `npm install -g @modelcontextprotocol/server-filesystem` |
| sqlite | SQLite数据库 | `npm install -g @modelcontextprotocol/server-sqlite` |
| git | Git操作 | `npm install -g @modelcontextprotocol/server-git` |
| postgres | PostgreSQL | `npm install -g @modelcontextprotocol/server-postgres` |
| puppeteer | 网页自动化 | `npm install -g @modelcontextprotocol/server-puppeteer` |

## 社区资源

- [MCP官方文档](https://modelcontextprotocol.io/)
- [Go SDK文档](https://pkg.go.dev/github.com/modelcontextprotocol/go-sdk)
- [MCP服务器列表](https://github.com/modelcontextprotocol/servers)
- [示例项目](https://github.com/modelcontextprotocol/examples)