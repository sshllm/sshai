## MCP功能支持（Model Context Protocol）

使用官方库：`github.com/modelcontextprotocol/go-sdk`


## 实现功能
基于本项目（SSHAI），创建添加MCP客户端功能支持，要求支持主流的stdio、sse、streamableHttp三种协议，并支持在配置里自定义（多个）
MCP的工具列表应该在程序启动的时候加载获取（如果用户开启并且配置了MCP服务），并在与大模型交互的时候传递
工具列表应该间隔时间进行刷新，以使用最新的工具列表；
请尽量使用官方库提供的功能和函数，尽量避免自己编写HTTP客户端等方式来自定义实现；
调用MCP的过程（如参数、执行结果），需要在SSHAI的交互模式下输出显示，在SSHAI的stdin（管道模式）和exec（命令模式）中不需要显示

## 文档参考：
[client](mcp_client.md)
[protocol](mcp_protocol.md)

其他文档说明参考：
```markdown
MCP Go SDK 开发文档
1
2
MCP（Model Context Protocol）是一种开放标准，用于标准化大型语言模型（LLM）与外部数据源和工具之间的交互。MCP Go SDK 提供了官方的 Go 语言开发工具包，用于构建 MCP 客户端和服务器。

安装 MCP Go SDK

在项目根目录下运行以下命令安装 SDK：

go get github.com/modelcontextprotocol/go-sdk
复制
构建 MCP 客户端

以下是通过 MCP Go SDK 构建客户端的主要步骤：

创建客户端

使用 mcp.NewClient 方法创建一个 MCP 客户端：

client := mcp.NewClient(&mcp.Implementation{
Name: "mcp-client",
Version: "v1.0.0",
}, nil)
复制
设置传输方式

通过 exec.Command 启动 MCP 服务器，并使用 NewCommandTransport 设置传输方式：

transport := mcp.NewCommandTransport(exec.Command("myserver"))
复制
建立连接

通过传输方式与 MCP 服务器建立连接，创建会话对象：

session, err := client.Connect(ctx, transport)
if err != nil {
log.Fatal(err)
}
defer session.Close()
复制
获取提示词列表

使用 ListPrompts 方法获取服务器定义的提示词：

listPromptsParams := &mcp.ListPromptsParams{}
listPrompts, err := session.ListPrompts(ctx, listPromptsParams)
if err != nil {
panic(err)
}
for _, prompt := range listPrompts.Prompts {
fmt.Printf("- %s: %s\n", prompt.Name, prompt.Description)
}
复制
获取资源列表

通过 ListResources 方法查看服务器上的资源信息：

listResourcesParams := &mcp.ListResourcesParams{}
resources, err := session.ListResources(ctx, listResourcesParams)
if err != nil {
panic(err)
}
for _, resource := range resources.Resources {
fmt.Printf("- URI: %s, Name: %s, Description: %s\n", resource.URI, resource.Name, resource.Description)
}
复制
获取工具列表

使用 ListTools 方法获取服务器注册的工具信息：

toolsRequest := mcp.ListToolsRequest{}
tools, err := session.ListTools(ctx, toolsRequest)
if err != nil {
panic(err)
}
for _, tool := range tools.Tools {
fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
}
复制
调用工具

通过构造 CallToolParams，调用服务器上的工具并获取结果：

params := &mcp.CallToolParams{
Name: "greet",
Arguments: map[string]any{"name": "you"},
}
res, err := session.CallTool(ctx, params)
if err != nil {
log.Fatalf("CallTool failed: %v", err)
}
if res.IsError {
log.Fatal("Tool execution failed")
}
for _, c := range res.Content {
log.Print(c.(*mcp.TextContent).Text)
}
复制
构建 MCP 服务器

以下是一个简单的 MCP 服务器示例：

package main

import (
"context"
"log"
"github.com/modelcontextprotocol/go-sdk/mcp"
)

func SayHi(ctx context.Context, req *mcp.CallToolRequest, args struct {
Name string `json:"name"`
}) (*mcp.CallToolResult, any, error) {
return &mcp.CallToolResult{
Content: []mcp.Content{&mcp.TextContent{Text: "Hi " + args.Name}},
}, nil, nil
}

func main() {
server := mcp.NewServer(&mcp.Implementation{
Name: "greeter",
Version: "v1.0.0",
}, nil)

mcp.AddTool(server, &mcp.Tool{
Name: "greet",
Description: "Say hi",
}, SayHi)

if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
log.Fatal(err)
}
}
复制
总结

MCP Go SDK 提供了强大的工具，用于构建高效的 MCP 客户端和服务器。通过该 SDK，开发者可以轻松实现以下功能：

使用标准输入/输出与 MCP 服务器通信。

获取提示词、资源和工具列表。

调用远程工具并处理返回结果。

该 SDK 是构建与 LLM 交互应用的理想选择，能够显著提升开发效率。
```