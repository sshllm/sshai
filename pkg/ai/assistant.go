package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"sshai/pkg/config"
	"sshai/pkg/i18n"
	"sshai/pkg/models"
)

// Assistant AI助手结构体
type Assistant struct {
	messages     []models.ChatMessage
	currentModel string
	username     string
}

// NewAssistant 创建新的AI助手
func NewAssistant(username string) *Assistant {
	cfg := config.Get()

	// 初始化消息列表
	messages := make([]models.ChatMessage, 0)

	// 如果配置了系统提示词，添加到消息列表开头
	if cfg.Prompt.SystemPrompt != "" {
		messages = append(messages, models.ChatMessage{
			Role:    "system",
			Content: cfg.Prompt.SystemPrompt,
		})
	}

	return &Assistant{
		messages:     messages,
		currentModel: cfg.API.DefaultModel,
		username:     username,
	}
}

// SetModel 设置当前使用的模型
func (ai *Assistant) SetModel(model string) {
	ai.currentModel = model
}

// ClearContext 清空对话上下文
func (ai *Assistant) ClearContext() {
	ai.messages = make([]models.ChatMessage, 0)
}

// ProcessMessage 处理用户消息
func (ai *Assistant) ProcessMessage(input string, channel ssh.Channel, interrupt chan bool) {
	ai.ProcessMessageWithOptions(input, channel, interrupt, true)
}

func (ai *Assistant) ProcessMessageWithOptions(input string, channel ssh.Channel, interrupt chan bool, showAnimation bool) {
	// 添加用户消息到上下文
	ai.messages = append(ai.messages, models.ChatMessage{
		Role:    "user",
		Content: input,
	})

	// 调用AI API，根据参数决定是否显示动画
	if showAnimation {
		ai.callAIAPIWithLoading(channel, interrupt)
	} else {
		ai.callAIAPIWithoutLoading(channel, interrupt)
	}
}

// callAIAPIWithLoading 带加载动画的AI API调用
func (ai *Assistant) callAIAPIWithLoading(channel ssh.Channel, interrupt chan bool) {
	stopLoading := make(chan bool, 1)

	// 启动加载动画
	// go ai.showLoadingAnimation(channel, stopLoading)

	// 调用AI API
	ai.callAIAPI(channel, stopLoading, interrupt)
}

// callAIAPIWithoutLoading 不带加载动画的AI API调用
func (ai *Assistant) callAIAPIWithoutLoading(channel ssh.Channel, interrupt chan bool) {
	// 直接调用AI API，不显示动画
	ai.callAIAPIDirectly(channel, interrupt)
}

// callAIAPIDirectly 直接调用AI API，不显示加载动画
func (ai *Assistant) callAIAPIDirectly(channel ssh.Channel, interrupt chan bool) {
	cfg := config.Get()

	// 构建请求数据
	requestData := map[string]interface{}{
		"model":    ai.currentModel,
		"messages": ai.messages,
		"stream":   true,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		channel.Write([]byte(fmt.Sprintf("构建请求失败: %v\r\n", err)))
		return
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", cfg.API.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		channel.Write([]byte(fmt.Sprintf("创建请求失败: %v\r\n", err)))
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.API.APIKey)

	// 发送请求
	client := &http.Client{Timeout: time.Duration(cfg.API.Timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		channel.Write([]byte(fmt.Sprintf("请求失败: %v\r\n", err)))
		return
	}
	defer resp.Body.Close()

	// 直接处理流式响应，不显示加载动画
	ai.handleStreamResponse(resp, channel, interrupt)
}

// showLoadingAnimation 显示加载动画
func (ai *Assistant) showLoadingAnimation(channel ssh.Channel, stop chan bool) {
	cfg := config.Get()
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0

	ticker := time.NewTicker(time.Duration(cfg.Display.LoadingAnimationInterval) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			// 清除加载动画（只清除当前行，不影响其他内容）
			channel.Write([]byte("\r" + strings.Repeat(" ", 20) + "\r"))
			return
		case <-ticker.C:
			channel.Write([]byte(fmt.Sprintf("\r%s %s", spinner[i%len(spinner)], i18n.T("ai.thinking"))))
			i++
		}
	}
}

// showThinkingAnimation 显示思考动画
func (ai *Assistant) showThinkingAnimation(channel ssh.Channel, stop chan bool) {
	cfg := config.Get()
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0

	ticker := time.NewTicker(time.Duration(cfg.Display.ThinkingAnimationInterval) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			channel.Write([]byte(fmt.Sprintf("\r%s %s", spinner[i%len(spinner)], i18n.T("ai.thinking"))))
			i++
		}
	}
}

// callAIAPI 调用AI API
func (ai *Assistant) callAIAPI(channel ssh.Channel, stopLoading chan bool, interrupt chan bool) {
	cfg := config.Get()

	// 构建请求数据
	requestData := map[string]interface{}{
		"model":    ai.currentModel,
		"messages": ai.messages,
		"stream":   true,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		stopLoading <- true
		channel.Write([]byte(fmt.Sprintf("构建请求失败: %v\r\n", err)))
		return
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", cfg.API.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		stopLoading <- true
		channel.Write([]byte(fmt.Sprintf("创建请求失败: %v\r\n", err)))
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.API.APIKey)

	// 发送请求
	client := &http.Client{Timeout: time.Duration(cfg.API.Timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		stopLoading <- true
		channel.Write([]byte(fmt.Sprintf("请求失败: %v\r\n", err)))
		return
	}
	defer resp.Body.Close()

	// 停止加载动画
	stopLoading <- true

	// 处理流式响应
	ai.handleStreamResponse(resp, channel, interrupt)
}

// handleStreamResponse 处理流式响应 - 简化版本，避免截断
func (ai *Assistant) handleStreamResponse(resp *http.Response, channel ssh.Channel, interrupt chan bool) {
	var assistantMessage strings.Builder
	isThinking := false
	thinkingStartTime := time.Now()

	// 思考动画控制
	thinkingAnimationStop := make(chan bool, 1)
	var thinkingAnimationStarted bool

	defer func() {
		if thinkingAnimationStarted {
			thinkingAnimationStop <- true
		}
	}()

	// 逐行读取响应 - 使用更大的缓冲区并改进错误处理
	buffer := make([]byte, 32768) // 增加缓冲区大小到32KB
	var leftover []byte

	// 添加超时机制，防止长时间无响应
	lastDataTime := time.Now()
	maxWaitTime := 10 * 60 * time.Second // 最大等待时间

	for {
		select {
		case <-interrupt:
			channel.Write([]byte("\r\n[已中断]\r\n"))
			return
		default:
		}

		// 检查是否超时
		if time.Since(lastDataTime) > maxWaitTime {
			channel.Write([]byte("\r\n[响应超时，连接已断开]\r\n"))
			break
		}

		n, err := resp.Body.Read(buffer)
		if n == 0 {
			if err != nil {
				// 只有在没有数据且有错误时才退出
				break
			}
			// 如果没有数据但没有错误，继续等待
			time.Sleep(50 * time.Millisecond) // 增加等待时间
			continue
		}

		// 更新最后接收数据的时间
		lastDataTime = time.Now()

		data := append(leftover, buffer[:n]...)
		lines := bytes.Split(data, []byte("\n"))
		leftover = lines[len(lines)-1]
		lines = lines[:len(lines)-1]

		for _, line := range lines {
			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}

			if bytes.HasPrefix(line, []byte("data: ")) {
				jsonStr := string(line[6:])
				if jsonStr == "[DONE]" {
					// 正常结束流式响应
					goto streamEnd
				}

				var response models.ChatResponse
				if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
					// JSON解析错误，记录但继续处理
					continue
				}

				if len(response.Choices) > 0 {
					delta := response.Choices[0].Delta

					// 检查是否有思考内容
					thinkingText := delta.Reasoning
					if thinkingText == "" {
						thinkingText = delta.ReasoningContent
					}

					if thinkingText != "" {
						if !isThinking {
							// 第一次接收到思考内容
							isThinking = true
							thinkingStartTime = time.Now()

							// 显示思考开始信息（不清除屏幕，避免内容丢失）
							channel.Write([]byte(i18n.T("ai.thinking_process") + "\r\n"))
							// channel.Write([]byte("\r\n" + i18n.T("ai.thinking_process") + "\r\n"))
						}

						// 实时显示思考内容 - 直接输出，不处理
						if thinkingText != "" {
							// 转换换行符
							thinkingText = strings.ReplaceAll(thinkingText, "\n", "\r\n")
							channel.Write([]byte(thinkingText))
						}
					}

					// 处理正常回答内容 - 修复：先处理思考结束，再输出内容
					if delta.Content != "" {
						if isThinking {
							// 思考阶段结束，先显示完成信息
							thinkingDuration := time.Since(thinkingStartTime)
							// 显示思考完成信息
							// channel.Write([]byte(fmt.Sprintf("\r\n\n%s\r\n\n", i18n.T("ai.thinking_complete", thinkingDuration.Seconds()))))
							channel.Write([]byte(fmt.Sprintf("\r\n%s\r\n\n", i18n.T("ai.thinking_complete", thinkingDuration.Seconds()))))
							channel.Write([]byte(i18n.T("ai.response") + "\r\n"))
							isThinking = false
						}

						// 然后输出当前的内容（包括第一个delta.Content）
						content := delta.Content
						// 只转换换行符
						content = strings.ReplaceAll(content, "\n", "\r\n")
						channel.Write([]byte(content))

						// 保存到消息历史
						assistantMessage.WriteString(delta.Content)
					}
				}
			}
		}
	}

streamEnd:
	// 处理剩余的leftover数据，确保不丢失最后的内容
	if len(leftover) > 0 {
		leftover = bytes.TrimSpace(leftover)
		if len(leftover) > 0 && bytes.HasPrefix(leftover, []byte("data: ")) {
			jsonStr := string(leftover[6:])
			if jsonStr != "[DONE]" {
				var response models.ChatResponse
				if err := json.Unmarshal([]byte(jsonStr), &response); err == nil {
					if len(response.Choices) > 0 {
						delta := response.Choices[0].Delta
						if delta.Content != "" {
							content := delta.Content
							content = strings.ReplaceAll(content, "\n", "\r\n")
							channel.Write([]byte(content))
							assistantMessage.WriteString(delta.Content)
						}
					}
				}
			}
		}
	}

	// 如果还在思考状态，结束思考动画
	if isThinking && thinkingAnimationStarted {
		thinkingAnimationStop <- true
		thinkingDuration := time.Since(thinkingStartTime)
		channel.Write([]byte(fmt.Sprintf("\r%s\r\n\n", i18n.T("ai.thinking_complete", thinkingDuration.Seconds()))))

		if assistantMessage.Len() > 0 {
			channel.Write([]byte(i18n.T("ai.response") + "\r\n"))
			// 直接输出，不使用WrapText
			content := assistantMessage.String()
			content = strings.ReplaceAll(content, "\n", "\r\n")
			channel.Write([]byte(content))
		}
	}

	// 添加助手回复到上下文
	if assistantMessage.Len() > 0 {
		ai.messages = append(ai.messages, models.ChatMessage{
			Role:    "assistant",
			Content: assistantMessage.String(),
		})
	}

	channel.Write([]byte("\r\n"))
}
