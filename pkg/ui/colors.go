package ui

import (
	"fmt"
	"strings"
)

// ANSI颜色代码
const (
	// 重置
	Reset = "\033[0m"

	// 前景色
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	// 亮色前景色
	BrightBlack   = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"

	// 背景色
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"

	// 样式
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	Blink     = "\033[5m"
	Reverse   = "\033[7m"
	Strike    = "\033[9m"
)

// ColorFunc 颜色函数类型
type ColorFunc func(string) string

// 预定义颜色函数
var (
	RedText     = ColorFunc(func(s string) string { return Red + s + Reset })
	GreenText   = ColorFunc(func(s string) string { return Green + s + Reset })
	YellowText  = ColorFunc(func(s string) string { return Yellow + s + Reset })
	BlueText    = ColorFunc(func(s string) string { return Blue + s + Reset })
	MagentaText = ColorFunc(func(s string) string { return Magenta + s + Reset })
	CyanText    = ColorFunc(func(s string) string { return Cyan + s + Reset })
	WhiteText   = ColorFunc(func(s string) string { return White + s + Reset })

	BrightRedText     = ColorFunc(func(s string) string { return BrightRed + s + Reset })
	BrightGreenText   = ColorFunc(func(s string) string { return BrightGreen + s + Reset })
	BrightYellowText  = ColorFunc(func(s string) string { return BrightYellow + s + Reset })
	BrightBlueText    = ColorFunc(func(s string) string { return BrightBlue + s + Reset })
	BrightMagentaText = ColorFunc(func(s string) string { return BrightMagenta + s + Reset })
	BrightCyanText    = ColorFunc(func(s string) string { return BrightCyan + s + Reset })
	BrightWhiteText   = ColorFunc(func(s string) string { return BrightWhite + s + Reset })

	BoldText      = ColorFunc(func(s string) string { return Bold + s + Reset })
	ItalicText    = ColorFunc(func(s string) string { return Italic + s + Reset })
	UnderlineText = ColorFunc(func(s string) string { return Underline + s + Reset })
)

// Colorize 为文本添加颜色
func Colorize(text, color string) string {
	return color + text + Reset
}

// Rainbow 彩虹文字效果
func Rainbow(text string) string {
	colors := []string{Red, Yellow, Green, Cyan, Blue, Magenta}
	var result strings.Builder

	for i, char := range text {
		color := colors[i%len(colors)]
		result.WriteString(color + string(char))
	}
	result.WriteString(Reset)

	return result.String()
}

// GradientText 渐变文字效果
func GradientText(text string, startColor, endColor string) string {
	// 简单的渐变实现，这里使用交替颜色
	var result strings.Builder
	for i, char := range text {
		if i%2 == 0 {
			result.WriteString(startColor + string(char))
		} else {
			result.WriteString(endColor + string(char))
		}
	}
	result.WriteString(Reset)
	return result.String()
}

// FormatPrompt 格式化提示符，为不同部分添加颜色
func FormatPrompt(username, hostname, model string) string {
	// 用户名：绿色
	coloredUsername := BrightGreenText(username)
	// 主机名：蓝色
	coloredHostname := BrightBlueText(hostname)
	// 模型名：黄色
	coloredModel := BrightYellowText(model)
	// 符号：白色
	symbols := BrightWhiteText("@")
	arrow := BrightWhiteText("> ")

	return fmt.Sprintf("%s%s%s:%s %s", coloredUsername, symbols, coloredHostname, coloredModel, arrow)
}

// FormatModelName 格式化模型名称
func FormatModelName(model string) string {
	return BrightYellowText(model)
}

// FormatHostname 格式化主机名
func FormatHostname(hostname string) string {
	return BrightBlueText(hostname)
}

// FormatUsername 格式化用户名
func FormatUsername(username string) string {
	return BrightGreenText(username)
}

// FormatStatus 格式化状态信息
func FormatStatus(status string, isSuccess bool) string {
	if isSuccess {
		return BrightGreenText("✓ " + status)
	}
	return BrightRedText("✗ " + status)
}

// FormatInfo 格式化信息文本
func FormatInfo(text string) string {
	return BrightCyanText(text)
}

// FormatWarning 格式化警告文本
func FormatWarning(text string) string {
	return BrightYellowText("⚠ " + text)
}

// FormatError 格式化错误文本
func FormatError(text string) string {
	return BrightRedText("✗ " + text)
}
