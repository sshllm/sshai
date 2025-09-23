package ai

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/sashabaranov/go-openai"
	"golang.org/x/crypto/ssh"
)

// SimpleToolCallBuffer 简单的工具调用缓冲区
type SimpleToolCallBuffer struct {
	Name      string
	Arguments string
	ID        string
}

// processToolCallSimple 简化的工具调用处理
func (c *OpenAIClient) processToolCallSimple(ctx context.Context, toolCall openai.ToolCall, channel ssh.Channel, assistantMessage *strings.Builder) {
	// 如果这是一个完整的工具调用（有名称和参数），直接处理
	if toolCall.Function.Name != "" && toolCall.Function.Arguments != "" {
		// 检查参数是否是有效的JSON
		var temp interface{}
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &temp); err == nil {
			log.Printf("检测到完整的工具调用: %s, 参数: %s", toolCall.Function.Name, toolCall.Function.Arguments)
			c.handleToolCall(ctx, toolCall, channel, assistantMessage)
			return
		}
	}

	// 否则，累积到缓冲区
	c.accumulateToolCall(toolCall)
}

// accumulateToolCall 累积工具调用信息
func (c *OpenAIClient) accumulateToolCall(toolCall openai.ToolCall) {
	// 确保pendingToolCalls已初始化
	if c.pendingToolCalls == nil {
		c.pendingToolCalls = make(map[string]*openai.ToolCall)
	}

	// 获取或创建工具调用记录
	var pending *openai.ToolCall
	if existing, exists := c.pendingToolCalls[toolCall.ID]; exists {
		pending = existing
	} else {
		pending = &openai.ToolCall{
			ID:   toolCall.ID,
			Type: toolCall.Type,
			Function: openai.FunctionCall{
				Name:      "",
				Arguments: "",
			},
		}
		c.pendingToolCalls[toolCall.ID] = pending
	}

	// 累积信息
	if toolCall.Function.Name != "" {
		pending.Function.Name = toolCall.Function.Name
		log.Printf("累积工具名称: %s (ID: %s)", toolCall.Function.Name, toolCall.ID)
	}

	if toolCall.Function.Arguments != "" {
		pending.Function.Arguments += toolCall.Function.Arguments
		log.Printf("累积工具参数: %s (当前长度: %d)", pending.Function.Arguments, len(pending.Function.Arguments))
	}
}

// processAllPendingToolCalls 处理所有待处理的工具调用
func (c *OpenAIClient) processAllPendingToolCalls(ctx context.Context, channel ssh.Channel, assistantMessage *strings.Builder) {
	if c.pendingToolCalls == nil {
		return
	}

	// 首先尝试智能合并工具调用
	c.mergeToolCalls()

	for id, toolCall := range c.pendingToolCalls {
		log.Printf("处理待处理的工具调用: %s (ID: %s)", toolCall.Function.Name, id)
		log.Printf("最终参数: %s (长度: %d)", toolCall.Function.Arguments, len(toolCall.Function.Arguments))

		// 只处理有名称的工具调用
		if toolCall.Function.Name != "" {
			// 创建副本
			finalTool := openai.ToolCall{
				ID:   toolCall.ID,
				Type: toolCall.Type,
				Function: openai.FunctionCall{
					Name:      toolCall.Function.Name,
					Arguments: toolCall.Function.Arguments,
				},
			}

			c.handleToolCall(ctx, finalTool, channel, assistantMessage)
		}

		// 清理
		delete(c.pendingToolCalls, id)
	}
}

// mergeToolCalls 智能合并工具调用
func (c *OpenAIClient) mergeToolCalls() {
	var nameOnlyCall *openai.ToolCall
	var argsOnlyCall *openai.ToolCall
	var argsOnlyID string

	// 找到只有名称的工具调用和只有参数的工具调用
	for id, toolCall := range c.pendingToolCalls {
		if toolCall.Function.Name != "" && toolCall.Function.Arguments == "" {
			nameOnlyCall = toolCall
		} else if toolCall.Function.Name == "" && toolCall.Function.Arguments != "" {
			argsOnlyCall = toolCall
			argsOnlyID = id
		}
	}

	// 如果找到了配对，进行合并
	if nameOnlyCall != nil && argsOnlyCall != nil {
		log.Printf("合并工具调用: %s + 参数(%d字符)", nameOnlyCall.Function.Name, len(argsOnlyCall.Function.Arguments))
		
		// 将参数合并到有名称的工具调用中
		nameOnlyCall.Function.Arguments = argsOnlyCall.Function.Arguments
		
		// 删除只有参数的工具调用
		delete(c.pendingToolCalls, argsOnlyID)
		
		log.Printf("合并完成: %s, 参数: %s", nameOnlyCall.Function.Name, nameOnlyCall.Function.Arguments)
	}
}