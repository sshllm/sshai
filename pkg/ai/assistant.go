package ai

import (
	"golang.org/x/crypto/ssh"
)

// Assistant AI助手结构体 - 重构为使用 go-openai 库
type Assistant struct {
	client   *OpenAIClient
	username string
}

// NewAssistant 创建新的AI助手
func NewAssistant(username string) *Assistant {
	return &Assistant{
		client:   NewOpenAIClient(username),
		username: username,
	}
}

// SetModel 设置当前使用的模型
func (ai *Assistant) SetModel(model string) {
	ai.client.SetModel(model)
}

// ClearContext 清空对话上下文
func (ai *Assistant) ClearContext() {
	ai.client.ClearContext()
}

// ProcessMessage 处理用户消息
func (ai *Assistant) ProcessMessage(input string, channel ssh.Channel, interrupt chan bool) {
	ai.client.ProcessMessage(input, channel, interrupt)
}

// ProcessMessageWithOptions 处理用户消息（可选动画）
func (ai *Assistant) ProcessMessageWithOptions(input string, channel ssh.Channel, interrupt chan bool, showAnimation bool) {
	ai.client.ProcessMessageWithOptions(input, channel, interrupt, showAnimation)
}

// GetCurrentModel 获取当前使用的模型
func (ai *Assistant) GetCurrentModel() string {
	return ai.client.GetCurrentModel()
}
