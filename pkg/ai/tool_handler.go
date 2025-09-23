package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/sashabaranov/go-openai"
	"golang.org/x/crypto/ssh"

	"sshai/pkg/config"
	"sshai/pkg/mcp"
)

// handleToolCall 处理工具调用
func (c *OpenAIClient) handleToolCall(ctx context.Context, toolCall openai.ToolCall, channel ssh.Channel, assistantMessage *strings.Builder) {
	if toolCall.Function.Name == "" {
		return
	}

	// 调试信息：输出原始工具调用信息
	log.Printf("=== 工具调用调试信息 ===")
	log.Printf("工具ID: %s", toolCall.ID)
	log.Printf("工具类型: %s", toolCall.Type)
	log.Printf("函数名称: %s", toolCall.Function.Name)
	log.Printf("函数参数: %s", toolCall.Function.Arguments)
	log.Printf("参数长度: %d", len(toolCall.Function.Arguments))
	log.Printf("=== 调试信息结束 ===")

	// 解析工具参数
	var arguments map[string]interface{}
	
	// 检查参数是否为空
	if toolCall.Function.Arguments == "" {
		arguments = make(map[string]interface{})
		log.Printf("工具参数为空，使用空参数对象")
	} else {
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments); err != nil {
			channel.Write([]byte(fmt.Sprintf("\r\n❌ 解析工具参数失败: %v\r\n原始参数: %s\r\n", err, toolCall.Function.Arguments)))
			
			// 将错误信息添加到对话上下文
			c.messages = append(c.messages, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    fmt.Sprintf("工具参数解析失败: %v", err),
				ToolCallID: toolCall.ID,
			})
			return
		}
	}

	log.Printf("解析后的工具参数: %+v", arguments)

	// 获取MCP管理器
	mcpManager := mcp.GetGlobalManager()
	if mcpManager == nil {
		channel.Write([]byte("\r\n❌ MCP管理器未初始化\r\n"))
		
		// 将错误信息添加到对话上下文
		c.messages = append(c.messages, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			Content:    "MCP管理器未初始化",
			ToolCallID: toolCall.ID,
		})
		return
	}

	// 调用MCP工具
	log.Printf("开始调用MCP工具: %s, 参数: %+v", toolCall.Function.Name, arguments)
	result, err := mcpManager.CallTool(toolCall.Function.Name, arguments, channel)
	if err != nil {
		log.Printf("MCP工具调用失败: %v", err)
		channel.Write([]byte(fmt.Sprintf("\r\n❌ 工具调用失败: %v\r\n", err)))
		
		// 将错误信息添加到对话上下文
		c.messages = append(c.messages, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			Content:    fmt.Sprintf("工具调用失败: %v", err),
			ToolCallID: toolCall.ID,
		})
		return
	}

	log.Printf("MCP工具调用成功: %s, 结果长度: %d", toolCall.Function.Name, len(result))
	// 工具调用成功，将实际结果添加到对话上下文
	c.messages = append(c.messages, openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleTool,
		Content:    result,
		ToolCallID: toolCall.ID,
	})

	// 记录工具调用到助手消息中
	toolCallInfo := fmt.Sprintf("[调用工具: %s]", toolCall.Function.Name)
	assistantMessage.WriteString(toolCallInfo)
	
	// 工具执行完成后，需要继续对话让AI根据工具结果生成回复
	c.continueConversationAfterTool(ctx, channel, assistantMessage)
}

// continueConversationAfterTool 工具执行后继续对话
func (c *OpenAIClient) continueConversationAfterTool(ctx context.Context, channel ssh.Channel, assistantMessage *strings.Builder) {
	log.Printf("工具执行完成，继续对话...")
	
	cfg := config.Get()
	
	// 创建新的聊天完成请求，让AI根据工具结果继续回复
	req := openai.ChatCompletionRequest{
		Model:       cfg.API.DefaultModel,
		Messages:    c.messages,
		MaxTokens:   4000, // 使用默认值
		Temperature: float32(cfg.API.Temperature),
		Stream:      true,
	}

	// 发起新的流式请求
	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		log.Printf("创建工具后续对话流失败: %v", err)
		channel.Write([]byte(fmt.Sprintf("\r\n❌ 继续对话失败: %v\r\n", err)))
		return
	}
	defer stream.Close()

	// 处理流式响应
	for {
		response, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Printf("接收工具后续对话流响应失败: %v", err)
			break
		}

		if len(response.Choices) > 0 {
			delta := response.Choices[0].Delta
			
			// 输出AI的后续回复
			if delta.Content != "" {
				assistantMessage.WriteString(delta.Content)
				// 将\n转换为\r\n以适配SSH终端
				formattedContent := strings.ReplaceAll(delta.Content, "\n", "\r\n")
				channel.Write([]byte(formattedContent))
			}
			
			// 检查是否完成
			if response.Choices[0].FinishReason == "stop" {
				break
			}
		}
	}
	
	// 将完整的助手回复添加到消息历史
	if assistantMessage.Len() > 0 {
		c.messages = append(c.messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: assistantMessage.String(),
		})
	}
	
	log.Printf("工具后续对话完成")
}