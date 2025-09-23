package ai

import (
	"encoding/json"
	"log"

	"github.com/sashabaranov/go-openai"
)

// debugToolCall 调试工具调用信息
func debugToolCall(toolCall openai.ToolCall) {
	log.Printf("=== 工具调用调试信息 ===")
	log.Printf("工具ID: %s", toolCall.ID)
	log.Printf("工具类型: %s", toolCall.Type)
	log.Printf("函数名称: %s", toolCall.Function.Name)
	log.Printf("函数参数: %s", toolCall.Function.Arguments)
	log.Printf("参数长度: %d", len(toolCall.Function.Arguments))
	
	// 检查参数是否为有效JSON
	if toolCall.Function.Arguments == "" {
		log.Printf("参数为空字符串")
	} else {
		// 尝试验证JSON格式
		var temp interface{}
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &temp); err != nil {
			log.Printf("参数JSON格式无效: %v", err)
			log.Printf("原始参数字节: %v", []byte(toolCall.Function.Arguments))
		} else {
			log.Printf("参数JSON格式有效")
		}
	}
	log.Printf("=== 调试信息结束 ===")
}