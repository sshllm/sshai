package main

import (
	"fmt"
	"sshai/pkg/version"
	"sshai/pkg/ui"
)

func main() {
	fmt.Println("=== 版本信息测试 ===")
	
	// 测试版本信息
	buildInfo := version.GetBuildInfo()
	fmt.Printf("Version: %s\n", buildInfo.Version)
	fmt.Printf("Git Commit: %s\n", buildInfo.GitCommit)
	fmt.Printf("Build Time: %s\n", buildInfo.BuildTime)
	fmt.Printf("Go Version: %s\n", buildInfo.GoVersion)
	fmt.Printf("Platform: %s\n", buildInfo.Platform)
	
	fmt.Println("\n=== Banner 显示测试 ===")
	
	// 测试Banner生成
	banner := ui.GenerateBanner()
	fmt.Print(banner)
	
	fmt.Println("\n=== 彩色提示符测试 ===")
	
	// 测试彩色提示符
	prompt := ui.GeneratePrompt("testuser", "sshai.top", "gpt-oss:20b")
	fmt.Print(prompt)
	fmt.Println("这是一个测试命令")
	
	fmt.Println("\n=== 其他UI元素测试 ===")
	
	// 测试其他UI元素
	fmt.Println(ui.FormatStatus("连接成功", true))
	fmt.Println(ui.FormatStatus("连接失败", false))
	fmt.Println(ui.FormatInfo("这是一条信息"))
	fmt.Println(ui.FormatWarning("这是一条警告"))
	fmt.Println(ui.FormatError("这是一条错误"))
}
