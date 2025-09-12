package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"sshai/pkg/config"
	"sshai/pkg/i18n"
	"sshai/pkg/ssh"
)

// showStartupBanner 显示程序启动时的欢迎banner
func showStartupBanner() {
	fmt.Println(config.WelcomeBanner)
}

func main() {
	// 定义命令行参数
	var configFile string
	flag.StringVar(&configFile, "c", "", "指定配置文件路径")
	flag.Parse()

	// 确定配置文件路径
	if configFile == "" {
		configFile = "config.yaml" // 默认配置文件
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("错误: 配置文件 '%s' 不存在\n", configFile)
		fmt.Println("请确保配置文件存在，或使用 -c 参数指定正确的配置文件路径")
		fmt.Println("用法: sshai -c config.yaml")
		os.Exit(1)
	}

	// 加载配置文件
	err := config.Load(configFile)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化i18n系统
	cfg := config.Get()
	language := i18n.Language(cfg.I18n.Language)
	if language == "" {
		language = i18n.LanguageZhCN // 默认使用中文
	}
	if err := i18n.Init(language); err != nil {
		log.Fatal(i18n.T("error.lang_load", err))
	}

	// 显示程序启动banner
	showStartupBanner()

	// 创建SSH服务器
	server, err := ssh.NewServer()
	if err != nil {
		log.Fatal(i18n.T("error.server_start", err))
	}

	// 启动服务器
	log.Println(i18n.T("server.starting", cfg.Server.Port))
	if err := server.Start(); err != nil {
		log.Fatal(i18n.T("error.server_start", err))
	}
}
