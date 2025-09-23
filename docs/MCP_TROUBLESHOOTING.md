# MCP故障排除指南

## 常见问题及解决方案

### 1. "解析工具参数失败: unexpected end of JSON input"

**问题描述**: 在AI尝试调用MCP工具时出现JSON解析错误。

**可能原因**:
1. AI模型返回的工具参数为空字符串
2. AI模型返回的JSON格式不正确
3. 流式响应中工具调用信息不完整

**解决方案**:

#### 步骤1: 启用调试模式
在代码中已添加详细的调试日志，重新编译并运行程序：

```bash
go build -o sshai cmd/main.go
./sshai -c config.yaml
```

#### 步骤2: 查看调试信息
当出现错误时，查看控制台输出的调试信息：

```
=== 工具调用调试信息 ===
工具ID: call_xxx
工具类型: function
函数名称: list_files
函数参数: {"path": "/home/user"}
参数长度: 20
参数JSON格式有效
=== 调试信息结束 ===
```

#### 步骤3: 根据调试信息修复

**情况A: 参数为空字符串**
```
函数参数: 
参数长度: 0
参数为空字符串
```

**解决方法**: 代码已自动处理空参数情况，使用空对象`{}`。

**情况B: JSON格式无效**
```
函数参数: {path: /home/user}  # 缺少引号
参数JSON格式无效: invalid character 'p' after '{'
```

**解决方法**: 这通常是AI模型的问题，需要：
1. 调整系统提示词，明确要求返回有效JSON
2. 检查模型配置和温度设置
3. 考虑使用不同的AI模型

#### 步骤4: 优化系统提示词

在配置文件中添加更明确的工具使用指导：

```yaml
prompt:
  system_prompt: |
    你是一个专业的AI助手，具备MCP工具调用能力。
    
    重要：当调用工具时，请确保：
    1. 参数必须是有效的JSON格式
    2. 字符串值必须用双引号包围
    3. 如果工具不需要参数，传递空对象 {}
    
    可用工具示例：
    - list_files: 列出目录文件，参数 {"path": "目录路径"}
    - read_file: 读取文件，参数 {"path": "文件路径"}
    - write_file: 写入文件，参数 {"path": "文件路径", "content": "文件内容"}
```

### 2. "MCP管理器未初始化"

**问题描述**: 工具调用时提示MCP管理器未初始化。

**解决方案**:
1. 检查配置文件中`mcp.enabled`是否为`true`
2. 确认至少有一个MCP服务器配置且`enabled: true`
3. 查看启动日志确认MCP管理器是否成功启动

### 3. "工具调用失败"

**问题描述**: MCP工具调用返回错误。

**可能原因**:
1. MCP服务器未运行或连接失败
2. 工具参数不符合服务器要求
3. 权限问题

**解决方案**:

#### 检查MCP服务器状态
```bash
# 如果使用stdio传输，测试命令是否可执行
mcp-server-filesystem --help

# 如果使用HTTP传输，测试URL是否可访问
curl -X POST http://localhost:8080/mcp
```

#### 验证工具参数
查看MCP服务器文档，确认工具所需的参数格式：

```yaml
# 文件系统服务器示例
servers:
  - name: "filesystem"
    transport: "stdio"
    command: ["mcp-server-filesystem", "/allowed/path"]
    enabled: true
```

### 4. 工具列表为空

**问题描述**: 启动时显示"共 0 个工具"。

**解决方案**:
1. 确认MCP服务器配置正确
2. 检查服务器是否支持`tools/list`请求
3. 验证服务器权限和路径设置

### 5. 连接超时

**问题描述**: MCP服务器连接超时。

**解决方案**:
1. 增加连接超时时间（代码中默认10秒）
2. 检查网络连接
3. 验证服务器地址和端口

## 调试技巧

### 1. 启用详细日志
代码已包含详细的调试日志，包括：
- 工具调用详情
- 参数解析过程
- MCP服务器连接状态
- 错误详细信息

### 2. 测试MCP服务器
独立测试MCP服务器是否正常工作：

```bash
# 测试stdio服务器
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list", "params": {}}' | mcp-server-filesystem /path

# 测试HTTP服务器
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "id": 1, "method": "tools/list", "params": {}}'
```

### 3. 逐步排查
1. 先确认MCP管理器启动成功
2. 再确认工具列表获取成功
3. 最后测试工具调用功能

### 4. 使用模拟服务器
创建简单的模拟MCP服务器进行测试：

```javascript
// mock-mcp-server.js
const readline = require('readline');

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout
});

rl.on('line', (line) => {
  try {
    const request = JSON.parse(line);
    let response;
    
    if (request.method === 'tools/list') {
      response = {
        jsonrpc: "2.0",
        id: request.id,
        result: {
          tools: [
            {
              name: "test_tool",
              description: "A test tool",
              inputSchema: {
                type: "object",
                properties: {
                  message: { type: "string" }
                }
              }
            }
          ]
        }
      };
    } else if (request.method === 'tools/call') {
      response = {
        jsonrpc: "2.0",
        id: request.id,
        result: {
          content: [
            {
              type: "text",
              text: `Tool called with: ${JSON.stringify(request.params.arguments)}`
            }
          ]
        }
      };
    }
    
    console.log(JSON.stringify(response));
  } catch (e) {
    console.error('Error:', e);
  }
});
```

使用方法：
```yaml
servers:
  - name: "mock"
    transport: "stdio"
    command: ["node", "mock-mcp-server.js"]
    enabled: true
```

## 性能优化

### 1. 调整刷新间隔
```yaml
mcp:
  refresh_interval: 600  # 10分钟，减少频繁刷新
```

### 2. 禁用不需要的服务器
```yaml
servers:
  - name: "unused-server"
    enabled: false  # 禁用不需要的服务器
```

### 3. 设置合理的超时时间
在代码中调整超时设置，平衡响应速度和稳定性。

## 联系支持

如果问题仍然存在，请提供以下信息：
1. 完整的错误日志
2. 配置文件内容
3. MCP服务器类型和版本
4. 系统环境信息

这将帮助快速定位和解决问题。