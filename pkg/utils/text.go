package utils

import (
	"strings"
	"unicode/utf8"
)

// WrapText 文本换行处理
func WrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	lines := strings.Split(text, "\n")
	var result []string

	for _, line := range lines {
		if GetDisplayWidth(line) <= width {
			result = append(result, line)
			continue
		}

		// 需要换行的长行
		for len(line) > 0 {
			breakPos := FindBreakPosition(line, width)
			if breakPos == 0 {
				// 如果找不到合适的断点，强制断行
				runes := []rune(line)
				if len(runes) > width {
					breakPos = len(string(runes[:width]))
				} else {
					breakPos = len(line)
				}
			}

			result = append(result, line[:breakPos])
			line = strings.TrimLeft(line[breakPos:], " ")
		}
	}

	return strings.Join(result, "\n")
}

// GetDisplayWidth 计算文本显示宽度（中文字符占2个宽度）
func GetDisplayWidth(text string) int {
	width := 0
	for _, r := range text {
		if utf8.RuneLen(r) > 1 {
			width += 2 // 中文字符占2个宽度
		} else {
			width += 1 // 英文字符占1个宽度
		}
	}
	return width
}

// FindBreakPosition 寻找合适的断行位置
func FindBreakPosition(text string, maxWidth int) int {
	if len(text) == 0 {
		return 0
	}

	// 寻找空格、标点符号等合适的断行位置
	breakChars := []rune{' ', ',', '.', '!', '?', ';', ':', '，', '。', '！', '？', '；', '：'}

	currentWidth := 0
	lastBreakBytePos := 0
	currentBytePos := 0

	for _, r := range text {
		runeWidth := 1
		if utf8.RuneLen(r) > 1 {
			runeWidth = 2
		}

		currentWidth += runeWidth

		// 检查是否是断行字符
		for _, bc := range breakChars {
			if r == bc {
				lastBreakBytePos = currentBytePos + utf8.RuneLen(r)
				break
			}
		}

		if currentWidth >= maxWidth {
			if lastBreakBytePos > 0 {
				return lastBreakBytePos
			}
			// 如果没有找到合适的断行位置，返回当前字符的起始位置
			return currentBytePos
		}

		currentBytePos += utf8.RuneLen(r)
	}

	return 0
}
