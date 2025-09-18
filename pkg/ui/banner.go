package ui

import (
	"fmt"
	"strings"

	"sshai/pkg/version"
)

// GenerateBanner 生成彩色Banner
func GenerateBanner() string {
	buildInfo := version.GetBuildInfo()

	// 创建彩色的SSHAI LLVM标题
	title1 := Rainbow("SSH") + " " + Rainbow("LLVM")

	// 版本信息
	versionLine := BrightGreenText("🚀 Multi-User SSH AI Assistant ") + BrightWhiteText(buildInfo.Version)

	// 功能特性行
	features := []string{
		BrightCyanText("⚡ AI-Powered"),
		BrightMagentaText("👥 Multi-User"),
		BrightYellowText("⚡ Real-time"),
	}
	featureLine := strings.Join(features, "   ")

	// 链接信息
	websiteLine := BrightBlueText("🌍 Website: ") + BrightWhiteText("https://sshllm.top")
	githubLine := BrightYellowText("📦 GitHub:  ") + BrightWhiteText("https://github.com/sshllm/sshai")
	buildLine := BrightRedText("🔨 Built:   ") + BrightWhiteText(version.FormatBuildTime())

	// 底部说明
	bottomLine := BrightCyanText("👨‍💻 Built for modern developers & teams")

	// 组装Banner
	banner := fmt.Sprintf(`
%s

%s
─────────────────────────────────────────────

%s

%s
%s
%s

%s

`, title1, versionLine, featureLine, websiteLine, githubLine, buildLine, bottomLine)

	return banner
}

// GenerateWelcomeMessage 生成欢迎消息
func GenerateWelcomeMessage(username string) string {
	greeting := fmt.Sprintf("Hello, %s!", BrightGreenText(username))
	separator := strings.Repeat("=", len("Hello, "+username+"!"))

	helpText := BrightCyanText("Type 'help' for available commands")

	return fmt.Sprintf(`
%s
%s

%s

`, greeting, separator, helpText)
}

// GeneratePrompt 生成彩色提示符
func GeneratePrompt(username, hostname, model string) string {
	return FormatPrompt(username, hostname, model)
}

// GenerateModelInfo 生成模型信息显示
func GenerateModelInfo(model string) string {
	return fmt.Sprintf("Current model: %s", FormatModelName(model))
}

// GenerateConnectionInfo 生成连接信息
func GenerateConnectionInfo(port string) string {
	return fmt.Sprintf("%s %s",
		BrightGreenText("✓ SSH AI Server listening on port"),
		BrightYellowText(port))
}

// GenerateStartupInfo 生成启动信息
func GenerateStartupInfo(port string) string {
	connectCmd := fmt.Sprintf("ssh localhost -p %s", port)

	return fmt.Sprintf(`%s
%s %s`,
		GenerateConnectionInfo(port),
		BrightCyanText("Connect with:"),
		BrightWhiteText(connectCmd))
}
