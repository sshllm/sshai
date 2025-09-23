package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/sashabaranov/go-openai"
	"golang.org/x/crypto/ssh"
)

// processToolCallDelta 处理流式响应中的工具调用增量
func (c *OpenAIClient) processToolCallDelta(ctx context.Context, toolCall openai.ToolCall, channel ssh.Channel, assistantMessage *strings.Builder) {
	// 处理可能的ID为空的情况，使用索引作为临时ID
	toolID := toolCall.ID
	if toolID == "" {
		// 如果没有ID，尝试找到第一个未完成的工具调用
		for id, pending := range c.pendingToolCalls {
			if pending.Function.Name == "" || pending.Function.Arguments == "" {
				toolID = id
				break
			}
		}
		// 如果还是没找到，创建一个临时ID
		if toolID == "" {
			toolID = fmt.Sprintf("temp_%d", len(c.pendingToolCalls))
		}
	}

	// 获取或创建待处理的工具调用
	var pendingCall *openai.ToolCall
	if existing, exists := c.pendingToolCalls[toolID]; exists {
		pendingCall = existing
	} else {
		// 创建新的工具调用记录
		pendingCall = &openai.ToolCall{
			ID:   toolID,
			Type: toolCall.Type,
			Function: openai.FunctionCall{
				Name:      "",
				Arguments: "",
			},
		}
		c.pendingToolCalls[toolID] = pendingCall
	}

	// 累积工具调用信息
	if toolCall.Function.Name != "" {
		pendingCall.Function.Name = toolCall.Function.Name
		log.Printf("工具调用名称: %s (ID: %s)", toolCall.Function.Name, toolID)
	}

	if toolCall.Function.Arguments != "" {
		pendingCall.Function.Arguments += toolCall.Function.Arguments
		log.Printf("累积工具参数: %s (当前长度: %d)", pendingCall.Function.Arguments, len(pendingCall.Function.Arguments))
	}

	// 检查工具调用是否完整
	if c.isToolCallComplete(pendingCall) {
		log.Printf("工具调用完整，开始处理: %s", pendingCall.Function.Name)
		log.Printf("完整的工具调用参数: %s", pendingCall.Function.Arguments)
		log.Printf("传递给handleToolCall的参数长度: %d", len(pendingCall.Function.Arguments))
		
		// 创建一个副本来传递，避免引用问题
		completeTool := openai.ToolCall{
			ID:   pendingCall.ID,
			Type: pendingCall.Type,
			Function: openai.FunctionCall{
				Name:      pendingCall.Function.Name,
				Arguments: pendingCall.Function.Arguments,
			},
		}
		
		log.Printf("副本工具调用参数: %s", completeTool.Function.Arguments)
		c.handleToolCall(ctx, completeTool, channel, assistantMessage)
		// 清理已处理的工具调用
		delete(c.pendingToolCalls, toolID)
	}
}

// isToolCallComplete 检查工具调用是否完整
func (c *OpenAIClient) isToolCallComplete(toolCall *openai.ToolCall) bool {
	// 必须有工具名称
	if toolCall.Function.Name == "" {
		log.Printf("工具调用不完整: 缺少工具名称")
		return false
	}

	// 如果没有参数，暂时认为不完整，等待更多数据
	// 只有在流结束时才处理无参数的工具调用
	if toolCall.Function.Arguments == "" {
		log.Printf("工具调用不完整: 参数为空，等待更多数据")
		return false
	}

	// 检查JSON是否有效且完整
	var temp interface{}
	err := json.Unmarshal([]byte(toolCall.Function.Arguments), &temp)
	if err != nil {
		log.Printf("工具调用参数JSON不完整: %v, 参数: %s", err, toolCall.Function.Arguments)
		return false
	}

	log.Printf("工具调用JSON格式有效，认为完整")
	return true
}

// finalizePendingToolCalls 在流结束时处理所有待处理的工具调用
func (c *OpenAIClient) finalizePendingToolCalls(ctx context.Context, channel ssh.Channel, assistantMessage *strings.Builder) {
	for id, toolCall := range c.pendingToolCalls {
		log.Printf("处理未完成的工具调用: %s (ID: %s)", toolCall.Function.Name, id)
		log.Printf("最终工具调用参数: %s (长度: %d)", toolCall.Function.Arguments, len(toolCall.Function.Arguments))
		
		// 即使参数不完整，也尝试处理
		if toolCall.Function.Name != "" {
			// 创建副本
			completeTool := openai.ToolCall{
				ID:   toolCall.ID,
				Type: toolCall.Type,
				Function: openai.FunctionCall{
					Name:      toolCall.Function.Name,
					Arguments: toolCall.Function.Arguments,
				},
			}
			c.handleToolCall(ctx, completeTool, channel, assistantMessage)
		}
		
		// 清理
		delete(c.pendingToolCalls, id)
	}
}