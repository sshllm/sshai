package mcp

import (
	"log"
	"sync"
)

// GlobalManager 全局MCP管理器实例
var (
	GlobalManager *MCPManager
	once          sync.Once
)

// InitGlobalManager 初始化全局MCP管理器
func InitGlobalManager() error {
	var err error
	once.Do(func() {
		GlobalManager = NewMCPManager()
		err = GlobalManager.Start()
		if err != nil {
			log.Printf("初始化MCP管理器失败: %v", err)
		}
	})
	return err
}

// GetGlobalManager 获取全局MCP管理器
func GetGlobalManager() *MCPManager {
	if GlobalManager == nil {
		if err := InitGlobalManager(); err != nil {
			log.Printf("获取MCP管理器失败: %v", err)
			return nil
		}
	}
	return GlobalManager
}

// StopGlobalManager 停止全局MCP管理器
func StopGlobalManager() {
	if GlobalManager != nil {
		GlobalManager.Stop()
		GlobalManager = nil
		once = sync.Once{} // 重置once，允许重新初始化
	}
}