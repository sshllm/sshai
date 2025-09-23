package mcp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golang.org/x/crypto/ssh"

	"sshai/pkg/config"
	"sshai/pkg/i18n"
)

// Tool MCP工具信息
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema"`
	ServerName  string                 `json:"server_name"`
}

// MCPManager MCP管理器
type MCPManager struct {
	clients   map[string]*mcp.ClientSession // 服务器名称 -> 客户端会话
	tools     []Tool                        // 可用工具列表
	mutex     sync.RWMutex                  // 读写锁
	ctx       context.Context               // 上下文
	cancel    context.CancelFunc            // 取消函数
	refreshCh chan struct{}                 // 刷新通道
}

// NewMCPManager 创建新的MCP管理器
func NewMCPManager() *MCPManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &MCPManager{
		clients:   make(map[string]*mcp.ClientSession),
		tools:     make([]Tool, 0),
		ctx:       ctx,
		cancel:    cancel,
		refreshCh: make(chan struct{}, 1),
	}
}

// Start 启动MCP管理器
func (m *MCPManager) Start() error {
	cfg := config.Get()
	if !cfg.MCP.Enabled {
		log.Println("MCP功能未启用")
		return nil
	}

	log.Println("启动MCP管理器...")

	// 初始化连接
	if err := m.initializeConnections(); err != nil {
		return fmt.Errorf("初始化MCP连接失败: %v", err)
	}

	// 启动定期刷新
	go m.startRefreshLoop()

	log.Printf("MCP管理器启动成功，已连接 %d 个服务器", len(m.clients))
	return nil
}

// Stop 停止MCP管理器
func (m *MCPManager) Stop() {
	log.Println("停止MCP管理器...")
	m.cancel()

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 关闭所有客户端连接
	for name, client := range m.clients {
		if err := client.Close(); err != nil {
			log.Printf("关闭MCP客户端 %s 失败: %v", name, err)
		}
	}
	m.clients = make(map[string]*mcp.ClientSession)
	m.tools = make([]Tool, 0)
}

// initializeConnections 初始化所有MCP连接
func (m *MCPManager) initializeConnections() error {
	cfg := config.Get()

	for _, serverCfg := range cfg.MCP.Servers {
		if !serverCfg.Enabled {
			continue
		}

		if err := m.connectToServer(serverCfg); err != nil {
			log.Printf("连接到MCP服务器 %s 失败: %v", serverCfg.Name, err)
			continue
		}
	}

	// 加载工具列表
	return m.refreshTools()
}

// connectToServer 连接到单个MCP服务器
func (m *MCPManager) connectToServer(serverCfg config.MCPServer) error {
	log.Printf("连接到MCP服务器: %s (传输方式: %s)", serverCfg.Name, serverCfg.Transport)

	// 创建MCP客户端
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "sshai",
		Version: "1.0.0",
	}, nil)

	var transport mcp.Transport
	var err error

	switch serverCfg.Transport {
	case "stdio":
		transport, err = m.createStdioTransport(serverCfg)
	case "http", "streamable":
		transport, err = m.createHTTPTransport(serverCfg)
	case "sse":
		transport, err = m.createSSETransport(serverCfg)
	default:
		return fmt.Errorf("不支持的传输方式: %s", serverCfg.Transport)
	}

	if err != nil {
		return fmt.Errorf("创建传输失败: %v", err)
	}

	// 建立连接
	session, err := client.Connect(m.ctx, transport, nil)
	if err != nil {
		return fmt.Errorf("连接失败: %v", err)
	}

	m.mutex.Lock()
	m.clients[serverCfg.Name] = session
	m.mutex.Unlock()

	log.Printf("成功连接到MCP服务器: %s", serverCfg.Name)
	return nil
}

// createStdioTransport 创建stdio传输
func (m *MCPManager) createStdioTransport(serverCfg config.MCPServer) (mcp.Transport, error) {
	if len(serverCfg.Command) == 0 {
		return nil, fmt.Errorf("stdio传输需要指定命令")
	}

	cmd := exec.Command(serverCfg.Command[0], serverCfg.Command[1:]...)
	return &mcp.CommandTransport{Command: cmd}, nil
}

// createHTTPTransport 创建HTTP传输
func (m *MCPManager) createHTTPTransport(serverCfg config.MCPServer) (mcp.Transport, error) {
	if serverCfg.URL == "" {
		return nil, fmt.Errorf("HTTP传输需要指定URL")
	}

	transport := &mcp.StreamableClientTransport{
		Endpoint: serverCfg.URL,
	}

	// 设置自定义HTTP客户端（如果需要添加请求头）
	if len(serverCfg.Headers) > 0 {
		httpClient := &http.Client{
			Transport: &headerTransport{
				base:    http.DefaultTransport,
				headers: serverCfg.Headers,
			},
		}
		transport.HTTPClient = httpClient
	}

	return transport, nil
}

// createSSETransport 创建SSE传输
func (m *MCPManager) createSSETransport(serverCfg config.MCPServer) (mcp.Transport, error) {
	// SSE传输通常也使用HTTP，但可能有不同的配置
	return m.createHTTPTransport(serverCfg)
}

// headerTransport HTTP传输包装器，用于添加自定义请求头
type headerTransport struct {
	base    http.RoundTripper
	headers map[string]string
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 添加自定义请求头
	for key, value := range t.headers {
		req.Header.Set(key, value)
	}
	return t.base.RoundTrip(req)
}

// refreshTools 刷新工具列表
func (m *MCPManager) refreshTools() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var allTools []Tool

	for serverName, client := range m.clients {
		tools, err := m.getToolsFromServer(serverName, client)
		if err != nil {
			log.Printf("从服务器 %s 获取工具失败: %v", serverName, err)
			continue
		}
		allTools = append(allTools, tools...)
	}

	m.tools = allTools
	log.Printf("刷新工具列表完成，共 %d 个工具", len(allTools))
	return nil
}

// getToolsFromServer 从服务器获取工具列表
func (m *MCPManager) getToolsFromServer(serverName string, client *mcp.ClientSession) ([]Tool, error) {
	ctx, cancel := context.WithTimeout(m.ctx, 10*time.Second)
	defer cancel()

	// 获取工具列表
	toolsResponse, err := client.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		return nil, err
	}

	var tools []Tool
	for _, tool := range toolsResponse.Tools {
		// 将jsonschema.Schema转换为map[string]interface{}
		var schema map[string]interface{}
		if tool.InputSchema != nil {
			// 简单转换，实际使用中可能需要更复杂的转换逻辑
			schema = map[string]interface{}{
				"type":       "object",
				"properties": make(map[string]interface{}),
			}
			if tool.InputSchema.Properties != nil {
				properties := make(map[string]interface{})
				for key, prop := range tool.InputSchema.Properties {
					properties[key] = map[string]interface{}{
						"type": prop.Type,
					}
					if prop.Description != "" {
						properties[key].(map[string]interface{})["description"] = prop.Description
					}
				}
				schema["properties"] = properties
			}
		}

		mcpTool := Tool{
			Name:        tool.Name,
			Description: tool.Description,
			Schema:      schema,
			ServerName:  serverName,
		}
		tools = append(tools, mcpTool)
	}

	return tools, nil
}

// startRefreshLoop 启动定期刷新循环
func (m *MCPManager) startRefreshLoop() {
	cfg := config.Get()
	refreshInterval := time.Duration(cfg.MCP.RefreshInterval) * time.Second
	if refreshInterval <= 0 {
		refreshInterval = 300 * time.Second // 默认5分钟
	}

	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			if err := m.refreshTools(); err != nil {
				log.Printf("定期刷新工具列表失败: %v", err)
			}
		case <-m.refreshCh:
			if err := m.refreshTools(); err != nil {
				log.Printf("手动刷新工具列表失败: %v", err)
			}
		}
	}
}

// GetTools 获取可用工具列表
func (m *MCPManager) GetTools() []Tool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// 返回工具列表的副本
	tools := make([]Tool, len(m.tools))
	copy(tools, m.tools)
	return tools
}

// CallTool 调用MCP工具
func (m *MCPManager) CallTool(toolName string, arguments map[string]interface{}, channel ssh.Channel) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// 查找工具
	var tool *Tool
	for _, t := range m.tools {
		if t.Name == toolName {
			tool = &t
			break
		}
	}

	if tool == nil {
		return "", fmt.Errorf("工具 %s 不存在", toolName)
	}

	// 获取对应的客户端
	client, exists := m.clients[tool.ServerName]
	if !exists {
		return "", fmt.Errorf("服务器 %s 未连接", tool.ServerName)
	}

	// 在交互模式下显示工具调用信息
	if channel != nil {
		channel.Write([]byte(fmt.Sprintf("\r\n🔧 %s %s...\r\n", i18n.T("mcp.calling_tool"), toolName)))
	}

	// 调用工具
	ctx, cancel := context.WithTimeout(m.ctx, 30*time.Second)
	defer cancel()

	params := &mcp.CallToolParams{
		Name:      toolName,
		Arguments: arguments,
	}

	result, err := client.CallTool(ctx, params)
	if err != nil {
		if channel != nil {
			channel.Write([]byte(fmt.Sprintf("❌ %s: %v\r\n", i18n.T("mcp.tool_error"), err)))
		}
		return "", fmt.Errorf("调用工具失败: %v", err)
	}

	// 处理结果
	if result.IsError {
		log.Printf("MCP工具执行错误: IsError=%v", result.IsError)
		log.Printf("MCP工具错误内容数量: %d", len(result.Content))
		for i, content := range result.Content {
			switch c := content.(type) {
			case *mcp.TextContent:
				log.Printf("错误内容[%d]: %s", i, c.Text)
			case *mcp.ImageContent:
				log.Printf("错误内容[%d]: 图片数据", i)
			default:
				log.Printf("错误内容[%d]: 未知类型 %T", i, c)
			}
		}
		if channel != nil {
			channel.Write([]byte(fmt.Sprintf("❌ %s\r\n", i18n.T("mcp.tool_execution_error"))))
		}
		return "", fmt.Errorf("工具执行失败")
	}

	// 收集工具执行结果
	var resultText string
	for _, content := range result.Content {
		switch c := content.(type) {
		case *mcp.TextContent:
			resultText += c.Text + "\n"
		case *mcp.ImageContent:
			resultText += fmt.Sprintf("[图片: %s]\n", c.Data)
		default:
			resultText += "[未知内容类型]\n"
		}
	}

	// 在交互模式下显示工具结果
	if channel != nil {
		channel.Write([]byte(fmt.Sprintf("✅ %s %s\r\n", i18n.T("mcp.tool_success"), toolName)))
		if resultText != "" {
			// 将\n转换为\r\n以适配SSH终端
			formattedResult := strings.ReplaceAll(resultText, "\n", "\r\n")
			channel.Write([]byte(formattedResult))
		}
		channel.Write([]byte("\r\n"))
	}

	return resultText, nil
}

// RefreshTools 手动刷新工具列表
func (m *MCPManager) RefreshTools() {
	select {
	case m.refreshCh <- struct{}{}:
	default:
		// 如果通道已满，忽略这次刷新请求
	}
}

// GetServerStatus 获取服务器状态
func (m *MCPManager) GetServerStatus() map[string]bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	status := make(map[string]bool)
	cfg := config.Get()

	for _, serverCfg := range cfg.MCP.Servers {
		if !serverCfg.Enabled {
			status[serverCfg.Name] = false
			continue
		}

		_, connected := m.clients[serverCfg.Name]
		status[serverCfg.Name] = connected
	}

	return status
}