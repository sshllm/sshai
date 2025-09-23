package ai

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
	"golang.org/x/crypto/ssh"

	"sshai/pkg/config"
	"sshai/pkg/i18n"
	"sshai/pkg/mcp"
)

// OpenAIClient 基于 go-openai 库的客户端
type OpenAIClient struct {
	client            *openai.Client
	messages          []openai.ChatCompletionMessage
	username          string
	currentModel      string // 添加当前模型字段
	pendingToolCalls  map[string]*openai.ToolCall // 缓存不完整的工具调用
}

// NewOpenAIClient 创建新的 OpenAI 客户端
func NewOpenAIClient(username string) *OpenAIClient {
	cfg := config.Get()

	// 创建 OpenAI 客户端配置
	clientConfig := openai.DefaultConfig(cfg.API.APIKey)
	clientConfig.BaseURL = cfg.API.BaseURL

	// 创建自定义 HTTP 客户端，支持更好的取消机制
	clientConfig.HTTPClient = &http.Client{
		Timeout: time.Duration(cfg.API.Timeout) * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true, // 禁用连接复用，便于快速取消
		},
	}

	client := openai.NewClientWithConfig(clientConfig)

	// 初始化消息列表
	messages := make([]openai.ChatCompletionMessage, 0)

	// 如果配置了系统提示词，添加到消息列表开头
	if cfg.Prompt.SystemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: cfg.Prompt.SystemPrompt,
		})
	}

	return &OpenAIClient{
		client:            client,
		messages:          messages,
		username:          username,
		currentModel:      cfg.API.DefaultModel, // 初始化为默认模型
		pendingToolCalls:  make(map[string]*openai.ToolCall), // 初始化工具调用缓存
	}
}

// ProcessMessage 处理用户消息（带动画）
func (c *OpenAIClient) ProcessMessage(input string, channel ssh.Channel, interrupt chan bool) {
	c.ProcessMessageWithFullOptions(input, channel, interrupt, true, true)
}

// ProcessMessageWithOptions 处理用户消息（可选动画）
func (c *OpenAIClient) ProcessMessageWithOptions(input string, channel ssh.Channel, interrupt chan bool, showAnimation bool) {
	c.ProcessMessageWithFullOptions(input, channel, interrupt, showAnimation, true)
}

// ProcessMessageWithFullOptions 处理用户消息（完整选项）
func (c *OpenAIClient) ProcessMessageWithFullOptions(input string, channel ssh.Channel, interrupt chan bool, showAnimation bool, showToolOutput bool) {
	// 添加用户消息到上下文
	c.messages = append(c.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	})

	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动专门的中断监听 goroutine，直接监听原始中断通道
	go func() {
		for {
			select {
			case <-interrupt:
				cancel() // 立即取消请求
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	// 调用流式 API
	c.callStreamingAPI(ctx, channel, showAnimation, showToolOutput)
}

// GetAvailableTools 获取可用的MCP工具列表（用于传递给AI模型）
func (c *OpenAIClient) GetAvailableTools() []openai.Tool {
	mcpManager := mcp.GetGlobalManager()
	if mcpManager == nil {
		return nil
	}

	mcpTools := mcpManager.GetTools()
	if len(mcpTools) == 0 {
		return nil
	}

	var tools []openai.Tool
	for _, mcpTool := range mcpTools {
		tool := openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        mcpTool.Name,
				Description: mcpTool.Description,
				Parameters:  mcpTool.Schema,
			},
		}
		tools = append(tools, tool)
	}

	return tools
}

// callStreamingAPI 调用流式 API
func (c *OpenAIClient) callStreamingAPI(ctx context.Context, channel ssh.Channel, showAnimation bool, showToolOutput bool) {
	// 创建聊天完成请求
	req := openai.ChatCompletionRequest{
		Model:    c.currentModel, // 使用当前设置的模型
		Messages: c.messages,
		Stream:   true,
	}

	// 设置温度参数（如果配置了的话）
	cfg := config.Get()
	if cfg.API.Temperature > 0 {
		req.Temperature = float32(cfg.API.Temperature)
	}

	// 添加MCP工具（如果可用）
	if tools := c.GetAvailableTools(); len(tools) > 0 {
		req.Tools = tools
		req.ToolChoice = "auto" // 让AI自动决定是否使用工具
	}

	// 创建流式响应
	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		// 检查是否是因为上下文取消导致的错误
		if ctx.Err() == context.Canceled {
			channel.Write([]byte("\r\n[已中断]\r\n"))
			return
		}
		channel.Write([]byte(fmt.Sprintf("创建流式请求失败: %v\r\n", err)))
		return
	}
	defer stream.Close()

	// 处理流式响应
	c.handleStreamResponse(ctx, stream, channel, showAnimation, showToolOutput)
}

// handleStreamResponse 处理流式响应
func (c *OpenAIClient) handleStreamResponse(ctx context.Context, stream *openai.ChatCompletionStream, channel ssh.Channel, showAnimation bool, showToolOutput bool) {
	var assistantMessage strings.Builder
	isThinking := false
	thinkingStartTime := time.Now()

	// 创建响应通道
	responseChan := make(chan openai.ChatCompletionStreamResponse, 1)
	errorChan := make(chan error, 1)

	// 启动接收 goroutine
	go func() {
		defer close(responseChan)
		defer close(errorChan)

		for {
			// 在每次 Recv 前检查 context 状态
			select {
			case <-ctx.Done():
				return
			default:
			}

			response, err := stream.Recv()
			if err != nil {
				errorChan <- err
				return
			}

			select {
			case responseChan <- response:
			case <-ctx.Done():
				return
			}
		}
	}()

	// 主处理循环
	for {
		select {
		case <-ctx.Done():
			stream.Close()
			channel.Write([]byte("\r\n[已中断]\r\n"))
			return

		case err := <-errorChan:
			if err.Error() == "EOF" {
				goto finish
			}
			if ctx.Err() == context.Canceled {
				channel.Write([]byte("\r\n[已中断]\r\n"))
				return
			}
			channel.Write([]byte(fmt.Sprintf("\r\n流式响应错误: %v\r\n", err)))
			return

		case response := <-responseChan:
			// 在处理每个响应前再次检查 context
			select {
			case <-ctx.Done():
				stream.Close()
				channel.Write([]byte("\r\n[已中断]\r\n"))
				return
			default:
			}

			// 处理响应数据
			if len(response.Choices) > 0 {
				delta := response.Choices[0].Delta

				// 检查是否有思考内容（DeepSeek 等模型支持）
				if delta.ReasoningContent != "" {
					if !isThinking {
						isThinking = true
						thinkingStartTime = time.Now()
						channel.Write([]byte(i18n.T("ai.thinking_process") + "\r\n"))
					}

					// 输出思考内容
					thinkingText := strings.ReplaceAll(delta.ReasoningContent, "\n", "\r\n")
					channel.Write([]byte(thinkingText))
				}

				// 处理工具调用
				if len(delta.ToolCalls) > 0 {
					for _, toolCall := range delta.ToolCalls {
						c.processToolCallSimpleWithOptions(ctx, toolCall, channel, &assistantMessage, showToolOutput)
					}
				}

				// 处理正常回答内容
				if delta.Content != "" {
					if isThinking {
						thinkingDuration := time.Since(thinkingStartTime)
						channel.Write([]byte(fmt.Sprintf("\r\n%s\r\n\n", i18n.T("ai.thinking_complete", thinkingDuration.Seconds()))))
						channel.Write([]byte(i18n.T("ai.response") + "\r\n"))
						isThinking = false
					}

					content := delta.Content
					content = strings.ReplaceAll(content, "\n", "\r\n")
					channel.Write([]byte(content))
					assistantMessage.WriteString(delta.Content)
				}
			}
		}
	}

finish:
	// 处理任何未完成的工具调用
	c.processAllPendingToolCallsWithOptions(ctx, channel, &assistantMessage, showToolOutput)

	// 添加助手回复到上下文
	if assistantMessage.Len() > 0 {
		c.messages = append(c.messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: assistantMessage.String(),
		})
	}

	channel.Write([]byte("\r\n"))
}

// ClearContext 清空对话上下文
func (c *OpenAIClient) ClearContext() {
	cfg := config.Get()
	c.messages = make([]openai.ChatCompletionMessage, 0)
	c.pendingToolCalls = make(map[string]*openai.ToolCall) // 清空待处理的工具调用

	// 重新添加系统提示词
	if cfg.Prompt.SystemPrompt != "" {
		c.messages = append(c.messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: cfg.Prompt.SystemPrompt,
		})
	}
}

// SetModel 设置当前使用的模型
func (c *OpenAIClient) SetModel(model string) {
	c.currentModel = model
}

// GetCurrentModel 获取当前使用的模型
func (c *OpenAIClient) GetCurrentModel() string {
	return c.currentModel
}
