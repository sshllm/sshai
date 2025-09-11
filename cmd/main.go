package main

import (
	"log"

	"sshai/pkg/config"
	"sshai/pkg/i18n"
	"sshai/pkg/ssh"
)

func main() {
	// 加载配置文件
	err := config.Load("config.yaml")
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
