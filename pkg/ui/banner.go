package ui

import (
	"fmt"
	"strings"

	"sshai/pkg/version"
)

// GenerateBanner 生成彩色Banner
func GenerateBanner() string {
	buildInfo := version.GetBuildInfo()

	// ASCII艺术画 - 使用字符串拼接避免转义问题
	line1 := "  .-')     .-')    ('-. .-.         ('-.              "
	line2 := " ( OO ).  ( OO ). ( OO )  /        ( OO ).-.          "
	line3 := "(_)---\\_)(_)---\\_),--. ,--.        / . --. /  ,-.-')  "
	line4 := "/    _ | /    _ | |  | |  |        | \\-.  \\   |  |OO) "
	line5 := "\\  :" + "`" + " " + "`" + ". \\  :" + "`" + " " + "`" + ". |   .|  |      .-'-'  |  |  |  |  \\ "
	line6 := " '.." + "`" + "''.)" + " '.." + "`" + "''.)|       |       \\| |_.'  |  |  |(_/ "
	line7 := ".-._)   \\.-._)   \\|  .-.  |        |  .-.  | ,|  |_.' "
	line8 := "\\       /\\       /|  | |  |        |  | |  |(_|  |    "
	line9 := " " + "`" + "-----'  " + "`" + "-----' " + "`" + "--' " + "`" + "--'        " + "`" + "--' " + "`" + "--'  " + "`" + "--'    "

	asciiArt := BrightCyanText(line1 + "\n" + line2 + "\n" + line3 + "\n" + line4 + "\n" + line5 + "\n" + line6 + "\n" + line7 + "\n" + line8 + "\n" + line9)

	// 版本信息
	versionLine := BrightGreenText("🚀 SSH AI Assistant ") + BrightWhiteText(buildInfo.Version)

	// 功能特性行
	features := []string{
		BrightCyanText("⚡ AI-Powered"),
		BrightYellowText("⚡ Real-time"),
		BrightMagentaText("🔒 Secure"),
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

%s

%s
%s
%s

%s

`, asciiArt, versionLine, featureLine, websiteLine, githubLine, buildLine, bottomLine)

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

// GenerateConnectionInfo 生成连接信息
func GenerateConnectionInfo(port string) string {
	return fmt.Sprintf(`%s
%s %s
%s %s`,
		BrightGreenText("🚀 SSHAI Server Started Successfully!"),
		BrightCyanText("📡 Listening on port:"),
		BrightYellowText(port),
		BrightCyanText("🔗 Connect with:"),
		BrightWhiteText(fmt.Sprintf("ssh localhost -p %s", port)),
	)
}

// GenerateStartupInfo 生成启动信息
func GenerateStartupInfo(port string) string {
	connectCmd := fmt.Sprintf("ssh localhost -p %s", port)

	return fmt.Sprintf(`%s
%s %s`,
		GenerateConnectionInfo(port),
		BrightCyanText("Connect with:"),
		BrightWhiteText(connectCmd),
	)
}
