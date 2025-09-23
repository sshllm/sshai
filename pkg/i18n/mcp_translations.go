package i18n

// MCP相关的翻译键值对
var mcpTranslations = map[Language]map[string]string{
	LanguageZhCN: {
		"mcp.calling_tool":         "正在调用工具",
		"mcp.tool_error":           "工具调用错误",
		"mcp.tool_execution_error": "工具执行失败",
		"mcp.tool_success":         "工具执行成功:",
		"mcp.manager_starting":     "启动MCP管理器...",
		"mcp.manager_started":      "MCP管理器启动成功",
		"mcp.manager_stopped":      "MCP管理器已停止",
		"mcp.server_connected":     "已连接到MCP服务器:",
		"mcp.server_disconnected":  "与MCP服务器断开连接:",
		"mcp.tools_refreshed":      "工具列表已刷新",
		"mcp.no_tools":             "没有可用的MCP工具",
	},
	LanguageEnUS: {
		"mcp.calling_tool":         "Calling tool",
		"mcp.tool_error":           "Tool call error",
		"mcp.tool_execution_error": "Tool execution failed",
		"mcp.tool_success":         "Tool executed successfully:",
		"mcp.manager_starting":     "Starting MCP manager...",
		"mcp.manager_started":      "MCP manager started successfully",
		"mcp.manager_stopped":      "MCP manager stopped",
		"mcp.server_connected":     "Connected to MCP server:",
		"mcp.server_disconnected":  "Disconnected from MCP server:",
		"mcp.tools_refreshed":      "Tool list refreshed",
		"mcp.no_tools":             "No MCP tools available",
	},
}

// AddMCPTranslations 添加MCP翻译到i18n系统
func AddMCPTranslations() {
	if globalI18n == nil {
		return
	}

	globalI18n.mutex.Lock()
	defer globalI18n.mutex.Unlock()

	// 将MCP翻译合并到扁平化消息映射中
	for lang, translations := range mcpTranslations {
		if _, exists := globalI18n.flatMessages[lang]; !exists {
			globalI18n.flatMessages[lang] = make(map[string]string)
		}
		for key, value := range translations {
			globalI18n.flatMessages[lang][key] = value
		}
	}
}