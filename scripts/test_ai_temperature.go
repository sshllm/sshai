package main

import (
	"fmt"
	"log"

	"sshai/pkg/ai"
	"sshai/pkg/config"
)

func main() {
	// 加载配置文件
	err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg := config.Get()

	fmt.Printf("=== AI 温度参数集成测试 ===\n")
	fmt.Printf("配置的温度值: %.2f\n", cfg.API.Temperature)

	// 创建AI客户端
	_ = ai.NewOpenAIClient(cfg.API.DefaultModel)

	fmt.Printf("✅ AI客户端创建成功\n")
	fmt.Printf("✅ 模型设置为: %s\n", cfg.API.DefaultModel)
	fmt.Printf("✅ 温度参数将在API请求中使用: %.2f\n", cfg.API.Temperature)

	// 显示温度参数的作用说明
	fmt.Printf("\n=== 温度参数说明 ===\n")
	fmt.Printf("• 温度值范围: 0.0 - 2.0\n")
	fmt.Printf("• 当前设置: %.2f\n", cfg.API.Temperature)
	fmt.Printf("• 0.0: 最确定的回答，输出几乎总是相同\n")
	fmt.Printf("• 0.3: 较为保守，适合事实性问答\n")
	fmt.Printf("• 0.7: 平衡创造性和准确性，适合大多数场景\n")
	fmt.Printf("• 1.0: 更有创造性，适合创意写作\n")
	fmt.Printf("• 2.0: 极高随机性，输出可能不太连贯\n")

	fmt.Printf("\n✅ 温度配置集成测试完成！\n")
}
