package ssh

import (
	"fmt"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/ssh"

	"sshai/pkg/ai"
	"sshai/pkg/config"
	"sshai/pkg/i18n"
	"sshai/pkg/ui"
)

// CommandHistory 命令历史结构体
type CommandHistory struct {
	commands []string
	index    int // 当前历史索引，-1表示没有在浏览历史
}

// NewCommandHistory 创建新的命令历史
func NewCommandHistory() *CommandHistory {
	return &CommandHistory{
		commands: make([]string, 0),
		index:    -1,
	}
}

// AddCommand 添加命令到历史
func (h *CommandHistory) AddCommand(cmd string) {
	if cmd != "" && (len(h.commands) == 0 || h.commands[len(h.commands)-1] != cmd) {
		h.commands = append(h.commands, cmd)
		// 限制历史命令数量
		if len(h.commands) > 100 {
			h.commands = h.commands[1:]
		}
	}
	h.index = -1 // 重置索引
}

// GetPrevious 获取上一个命令
func (h *CommandHistory) GetPrevious() string {
	if len(h.commands) == 0 {
		return ""
	}

	if h.index == -1 {
		h.index = len(h.commands) - 1
	} else if h.index > 0 {
		h.index--
	}

	return h.commands[h.index]
}

// GetNext 获取下一个命令
func (h *CommandHistory) GetNext() string {
	if len(h.commands) == 0 || h.index == -1 {
		return ""
	}

	if h.index < len(h.commands)-1 {
		h.index++
		return h.commands[h.index]
	} else {
		h.index = -1
		return ""
	}
}

// InputState 输入状态管理
type InputState struct {
	buffer     []rune // 使用rune数组以更好地处理中文
	cursorPos  int    // 光标位置（以rune为单位）
	displayPos int    // 显示位置（以字符宽度为单位）
}

// NewInputState 创建新的输入状态
func NewInputState() *InputState {
	return &InputState{
		buffer:     make([]rune, 0),
		cursorPos:  0,
		displayPos: 0,
	}
}

// String 返回当前输入的字符串
func (is *InputState) String() string {
	return string(is.buffer)
}

// Clear 清空输入状态
func (is *InputState) Clear() {
	is.buffer = is.buffer[:0]
	is.cursorPos = 0
	is.displayPos = 0
}

// SetText 设置输入文本
func (is *InputState) SetText(text string) {
	is.buffer = []rune(text)
	is.cursorPos = len(is.buffer)
	is.displayPos = calculateDisplayWidth(is.buffer)
}

// InsertRune 在光标位置插入字符
func (is *InputState) InsertRune(r rune) {
	// 在光标位置插入字符
	if is.cursorPos >= len(is.buffer) {
		is.buffer = append(is.buffer, r)
	} else {
		is.buffer = append(is.buffer[:is.cursorPos+1], is.buffer[is.cursorPos:]...)
		is.buffer[is.cursorPos] = r
	}
	is.cursorPos++
	is.updateDisplayPos()
}

// DeleteRune 删除光标前的字符
func (is *InputState) DeleteRune() bool {
	if is.cursorPos > 0 {
		is.buffer = append(is.buffer[:is.cursorPos-1], is.buffer[is.cursorPos:]...)
		is.cursorPos--
		is.updateDisplayPos()
		return true
	}
	return false
}

// MoveCursorLeft 向左移动光标
func (is *InputState) MoveCursorLeft() bool {
	if is.cursorPos > 0 {
		is.cursorPos--
		is.updateDisplayPos()
		return true
	}
	return false
}

// MoveCursorRight 向右移动光标
func (is *InputState) MoveCursorRight() bool {
	if is.cursorPos < len(is.buffer) {
		is.cursorPos++
		is.updateDisplayPos()
		return true
	}
	return false
}

// MoveCursorToStart 移动光标到开始
func (is *InputState) MoveCursorToStart() {
	is.cursorPos = 0
	is.displayPos = 0
}

// MoveCursorToEnd 移动光标到结束
func (is *InputState) MoveCursorToEnd() {
	is.cursorPos = len(is.buffer)
	is.updateDisplayPos()
}

// updateDisplayPos 更新显示位置
func (is *InputState) updateDisplayPos() {
	is.displayPos = calculateDisplayWidth(is.buffer[:is.cursorPos])
}

// calculateDisplayWidth 计算字符串的显示宽度
func calculateDisplayWidth(runes []rune) int {
	width := 0
	for _, r := range runes {
		if r < 128 {
			width++ // ASCII字符宽度为1
		} else {
			width += 2 // 中文字符宽度为2
		}
	}
	return width
}

// clearCurrentLine 清除当前行并重新显示提示符和输入内容
func clearCurrentLine(channel ssh.Channel, inputState *InputState, prompt string) {
	// 移动到行首
	channel.Write([]byte("\r"))

	// 清除整行 - 使用ANSI转义序列
	channel.Write([]byte("\033[K"))

	// 显示提示符
	channel.Write([]byte(prompt))

	// 显示当前输入内容
	if len(inputState.buffer) > 0 {
		channel.Write([]byte(inputState.String()))

		// 如果光标不在末尾，需要移动光标到正确位置
		if inputState.cursorPos < len(inputState.buffer) {
			// 计算需要向左移动的字符数
			rightPart := inputState.buffer[inputState.cursorPos:]
			rightWidth := calculateDisplayWidth(rightPart)

			// 向左移动光标
			for i := 0; i < rightWidth; i++ {
				channel.Write([]byte("\033[D")) // 向左移动一个位置
			}
		}
	}
}

// refreshLine 刷新当前行显示
func refreshLine(channel ssh.Channel, inputState *InputState, prompt string) {
	clearCurrentLine(channel, inputState, prompt)
}

// tryReadStdinInput 尝试读取stdin输入（仅在确认有管道输入时调用）
func tryReadStdinInput(channel ssh.Channel) string {
	// 使用单个goroutine和正确的退出机制
	result := make(chan string, 1)

	go func() {
		var localContent strings.Builder
		localBuffer := make([]byte, 8192)
		localHasData := false

		for {
			n, err := channel.Read(localBuffer)

			if err != nil {
				// 检查是否是EOF
				if err.Error() == "EOF" {
					log.Printf("遇到EOF，读取结束")
					break
				}
				// 其他错误
				log.Printf("读取错误: %v", err)
				break
			}

			if n > 0 {
				localHasData = true
				localContent.Write(localBuffer[:n])
				log.Printf("读取到数据，长度: %d", n)

				// 如果读取的数据小于缓冲区大小，可能已经读完
				if n < len(localBuffer) {
					log.Printf("数据读取完成")
					break
				}
			} else {
				// n == 0 且没有错误，表示没有更多数据
				log.Printf("没有更多数据")
				break
			}
		}

		if localHasData {
			result <- localContent.String()
		} else {
			result <- ""
		}
	}()

	// 等待读取完成或超时
	select {
	case content := <-result:
		log.Printf("读取完成，内容长度: %d", len(content))
		return content
	case <-time.After(2 * time.Second):
		log.Printf("读取超时，没有检测到stdin数据")
		return ""
	}
}

// readAllStdinContent 读取所有stdin内容
func readAllStdinContent(initialData []byte, channel ssh.Channel) string {
	var content strings.Builder
	content.Write(initialData)

	// 继续读取剩余数据
	buffer := make([]byte, 4096)
	for {
		// 使用短超时检查是否还有更多数据
		done := make(chan int, 1)
		errorChan := make(chan error, 1)

		go func() {
			n, err := channel.Read(buffer)
			if err != nil {
				errorChan <- err
				return
			}
			done <- n
		}()

		select {
		case n := <-done:
			if n > 0 {
				content.Write(buffer[:n])
				// 如果读取的数据小于缓冲区大小，可能已经读完
				if n < len(buffer) {
					break
				}
			} else {
				break
			}
		case <-errorChan:
			break
		case <-time.After(50 * time.Millisecond):
			// 短超时，没有更多数据
			break
		}
	}

	return content.String()
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// isTextContent 检测内容是否为纯文本
func isTextContent(content string) bool {
	// 检查是否包含大量二进制字符
	binaryCount := 0
	totalChars := len(content)

	if totalChars == 0 {
		return false
	}

	for _, b := range []byte(content) {
		// 检查是否为控制字符（除了常见的换行、制表符等）
		if b < 32 && b != 9 && b != 10 && b != 13 {
			binaryCount++
		}
		// 检查是否为高位字节（可能是二进制数据）
		if b > 127 {
			// 对于UTF-8编码的中文等，这里需要更精确的检测
			// 简单起见，我们允许一定比例的高位字节
		}
	}

	// 如果二进制字符超过5%，认为不是纯文本
	binaryRatio := float64(binaryCount) / float64(totalChars)
	if binaryRatio > 0.05 {
		return false
	}

	// 检查是否包含常见的二进制文件头
	contentLower := strings.ToLower(content)
	binaryHeaders := []string{
		"\x89png",      // PNG
		"\xff\xd8\xff", // JPEG
		"pk\x03\x04",   // ZIP
		"%pdf",         // PDF
		"\x00\x00\x00", // 可能的二进制文件
	}

	for _, header := range binaryHeaders {
		if strings.HasPrefix(contentLower, header) {
			return false
		}
	}

	return true
}

// handleStdinCommand 处理通过stdin传入的内容
func handleStdinCommand(channel ssh.Channel, username, content string) {
	log.Printf("处理stdin内容，用户: %s，内容长度: %d", username, len(content))

	cfg := config.Get()

	// 显示接收到的内容信息
	// channel.Write([]byte(fmt.Sprintf("接收到输入内容（%d字符）\r\n", len(content))))

	// 检查内容类型
	if !isTextContent(content) {
		channel.Write([]byte("错误：检测到非文本内容（如图片、PDF等二进制文件）\r\n"))
		channel.Write([]byte("本系统仅支持处理纯文本内容，请确保输入的是文本文件。\r\n"))
		channel.Write([]byte("支持的格式：.txt, .md, .log, .json, .yaml, .xml 等文本文件\r\n"))
		return
	}

	// 检查内容长度
	if len(content) == 0 {
		channel.Write([]byte("错误：输入内容为空\r\n"))
		return
	}

	// 直接使用默认模型，不加载模型列表
	selectedModel := cfg.API.DefaultModel

	// 创建AI助手
	assistant := ai.NewAssistant(username)
	assistant.SetModel(selectedModel)

	// 创建中断通道
	interrupt := make(chan bool, 1)

	// 构造提示消息，使用配置文件中的自定义提示词
	stdinPrompt := cfg.Prompt.StdinPrompt
	if stdinPrompt == "" {
		// 如果配置为空，使用默认提示词
		stdinPrompt = "请分析以下内容并提供相关的帮助或建议："
	}
	prompt := fmt.Sprintf("%s\n\n%s", stdinPrompt, content)

	// 直接处理内容并获取AI响应，不显示动画效果
	assistant.ProcessMessageWithOptions(prompt, channel, interrupt, false)
}

// HandleSession 处理SSH会话
func HandleSession(channel ssh.Channel, requests <-chan *ssh.Request, username string) {
	defer channel.Close()

	var execCommand string
	isExecMode := false
	hasPty := false // 标记是否有伪终端
	execReady := make(chan bool, 1)

	// 处理会话请求
	go func() {
		for req := range requests {
			switch req.Type {
			case "shell":
				req.Reply(true, nil)
				// 如果不是exec模式，通知可以继续
				if !isExecMode {
					select {
					case execReady <- false:
					default:
					}
				}
			case "pty-req":
				hasPty = true // 标记有伪终端
				req.Reply(true, nil)
			case "exec":
				// 处理执行命令请求
				if len(req.Payload) > 4 {
					// SSH exec请求的payload格式: [4字节长度][命令字符串]
					cmdLen := int(req.Payload[0])<<24 | int(req.Payload[1])<<16 | int(req.Payload[2])<<8 | int(req.Payload[3])
					if cmdLen > 0 && len(req.Payload) >= 4+cmdLen {
						execCommand = string(req.Payload[4 : 4+cmdLen])
						isExecMode = true
						log.Printf("接收到执行命令: %s", execCommand)
					}
				}
				req.Reply(true, nil)
				// 通知exec命令已准备好
				select {
				case execReady <- true:
				default:
				}
			default:
				req.Reply(false, nil)
			}
		}
	}()

	// 等待请求处理完成
	isExec := <-execReady

	cfg := config.Get()

	// 如果是执行模式且有命令，处理exec命令
	if isExec && execCommand != "" {
		handleExecCommand(channel, username, execCommand)
		return
	}

	// 如果没有伪终端，很可能是管道输入模式
	if !hasPty {
		log.Printf("检测到非PTY连接，尝试读取stdin输入")
		stdinContent := tryReadStdinInput(channel)
		if len(stdinContent) > 0 {
			log.Printf("读取到stdin内容，长度: %d", len(stdinContent))
			handleStdinCommand(channel, username, stdinContent)
			return
		}
	}

	// 发送登录成功消息（使用新的UI系统）
	banner := ui.GenerateBanner()
	lines := strings.Split(banner, "\n")
	for _, line := range lines {
		channel.Write([]byte(line + "\r\n"))
	}
	channel.Write([]byte("\r\n")) // 额外的空行分隔

	// 发送欢迎消息
	if username != "" {
		channel.Write([]byte(fmt.Sprintf(i18n.T("user.welcome")+", %s!\r\n", username)))
	} else {
		channel.Write([]byte(cfg.Server.WelcomeMessage + "\r\n"))
	}

	// 获取并选择模型
	channel.Write([]byte(i18n.T("model.loading") + "\r\n"))
	models, err := ai.GetAvailableModels()
	if err != nil {
		channel.Write([]byte(fmt.Sprintf(i18n.T("model.error", err) + "\r\n")))
		channel.Write([]byte(fmt.Sprintf("使用默认模型: %s\r\n", cfg.API.DefaultModel)))
		models = []ai.ModelInfo{{ID: cfg.API.DefaultModel}}
	}

	// 根据用户名匹配模型
	selectedModel := ai.SelectModelByUsername(channel, models, username)

	// 创建AI助手
	assistant := ai.NewAssistant(username)
	assistant.SetModel(selectedModel)

	// 生成彩色动态提示符
	hostname := "sshai.top" // 可以从配置或系统获取
	dynamicPrompt := ui.FormatPrompt(username, hostname, selectedModel)
	channel.Write([]byte("\r\n" + dynamicPrompt))

	// 处理用户输入
	handleUserInput(channel, assistant, dynamicPrompt)
}

// handleExecCommand 处理执行命令模式
func handleExecCommand(channel ssh.Channel, username, command string) {
	cfg := config.Get()

	// 显示执行的命令
	// channel.Write([]byte(fmt.Sprintf("执行命令: %s\r\n\r\n", command)))

	// 直接使用默认模型，不加载模型列表
	selectedModel := cfg.API.DefaultModel

	// 创建AI助手
	assistant := ai.NewAssistant(username)
	assistant.SetModel(selectedModel)

	// 创建中断通道
	interrupt := make(chan bool, 1)

	// 构造提示消息，使用配置文件中的自定义提示词
	execPrompt := cfg.Prompt.ExecPrompt
	if execPrompt == "" {
		// 如果配置为空，使用默认提示词
		execPrompt = "请回答以下问题或执行以下任务："
	}
	fullPrompt := fmt.Sprintf("%s\n\n%s", execPrompt, command)

	// 直接处理命令并获取AI响应，不显示动画效果
	assistant.ProcessMessageWithOptions(fullPrompt, channel, interrupt, false)

	// 添加换行符结束
	channel.Write([]byte("\r\n"))
}

// handleUserInput 处理用户输入
func handleUserInput(channel ssh.Channel, assistant *ai.Assistant, dynamicPrompt string) {
	buffer := make([]byte, 1024)
	history := NewCommandHistory()
	inputState := NewInputState()
	var escapeSequence []byte      // 用于处理ANSI转义序列
	var currentInterrupt chan bool // 当前正在使用的中断通道
	var isProcessing bool          // 标记是否正在处理AI请求

	for {
		n, err := channel.Read(buffer)
		if err != nil {
			log.Printf("读取输入失败: %v", err)
			return
		}

		data := buffer[:n]

		// 处理输入数据
		for i := 0; i < len(data); i++ {
			b := data[i]

			// 处理ANSI转义序列
			if len(escapeSequence) > 0 || b == 27 { // ESC键开始转义序列
				escapeSequence = append(escapeSequence, b)

				// 检查是否是完整的方向键序列
				if len(escapeSequence) == 3 && escapeSequence[0] == 27 && escapeSequence[1] == 91 {
					switch escapeSequence[2] {
					case 65: // 上方向键 ESC[A
						cmd := history.GetPrevious()
						inputState.SetText(cmd)
						refreshLine(channel, inputState, dynamicPrompt)
					case 66: // 下方向键 ESC[B
						cmd := history.GetNext()
						inputState.SetText(cmd)
						refreshLine(channel, inputState, dynamicPrompt)
					case 67: // 右方向键 ESC[C
						if inputState.MoveCursorRight() {
							// 获取刚移动过的字符，计算其显示宽度
							if inputState.cursorPos > 0 {
								movedChar := inputState.buffer[inputState.cursorPos-1]
								charWidth := 1
								if movedChar >= 128 {
									charWidth = 2 // 中文字符宽度为2
								}
								// 根据字符宽度移动终端光标
								for i := 0; i < charWidth; i++ {
									channel.Write([]byte("\033[C"))
								}
							}
						}
					case 68: // 左方向键 ESC[D
						if inputState.cursorPos > 0 {
							// 获取即将移动过的字符，计算其显示宽度
							charToMove := inputState.buffer[inputState.cursorPos-1]
							charWidth := 1
							if charToMove >= 128 {
								charWidth = 2 // 中文字符宽度为2
							}
							if inputState.MoveCursorLeft() {
								// 根据字符宽度移动终端光标
								for i := 0; i < charWidth; i++ {
									channel.Write([]byte("\033[D"))
								}
							}
						}
					}
					escapeSequence = nil
				} else if len(escapeSequence) > 3 {
					// 重置转义序列如果太长
					escapeSequence = nil
				}
				continue
			}

			// 处理普通字符
			r, size := utf8.DecodeRune(data[i:])
			if r == utf8.RuneError && size == 1 {
				// 跳过无效字符
				continue
			}

			// 跳过已处理的字节
			i += size - 1

			switch r {
			case 13: // Enter键
				// 如果正在处理AI请求，忽略新的输入
				if isProcessing {
					continue
				}

				input := strings.TrimSpace(inputState.String())
				channel.Write([]byte("\r\n"))

				// 添加非空命令到历史记录
				if input != "" {
					history.AddCommand(input)
				}

				if input == "/new" {
					assistant.ClearContext()
					channel.Write([]byte(i18n.T("user.clear_context") + "\r\n"))
				} else if input == "exit" || input == "quit" {
					channel.Write([]byte(i18n.T("user.exit") + "\r\n"))
					return
				} else if input != "" {
					// 设置处理状态
					isProcessing = true

					// 创建新的中断通道用于这次AI请求
					currentInterrupt = make(chan bool)

					// 异步处理AI请求，这样Ctrl+C可以在处理过程中被响应
					go func(userInput string, interruptCh chan bool) {
						assistant.ProcessMessage(userInput, channel, interruptCh)
						// 请求完成后清空引用和状态
						currentInterrupt = nil
						isProcessing = false
						// 显示提示符
						channel.Write([]byte(dynamicPrompt))
					}(input, currentInterrupt)
				} else {
					// 空输入，直接显示提示符
					channel.Write([]byte(dynamicPrompt))
				}
				// 清空输入状态
				inputState.Clear()

			case 127, 8: // Backspace/Delete键
				// 如果正在处理AI请求，忽略编辑操作
				if isProcessing {
					continue
				}
				if inputState.DeleteRune() {
					refreshLine(channel, inputState, dynamicPrompt)
				}

			case 1: // Ctrl+A - 移动到行首
				// 如果正在处理AI请求，忽略编辑操作
				if isProcessing {
					continue
				}
				inputState.MoveCursorToStart()
				refreshLine(channel, inputState, dynamicPrompt)

			case 5: // Ctrl+E - 移动到行尾
				// 如果正在处理AI请求，忽略编辑操作
				if isProcessing {
					continue
				}
				inputState.MoveCursorToEnd()
				refreshLine(channel, inputState, dynamicPrompt)

			case 3: // Ctrl+C
				// 如果有正在进行的AI请求，发送中断信号
				if currentInterrupt != nil {
					// 使用 goroutine 异步发送，避免阻塞
					go func(ch chan bool) {
						select {
						case ch <- true:
						case <-time.After(50 * time.Millisecond):
						}
					}(currentInterrupt)
				}

				channel.Write([]byte("\r\n^C\r\n"))
				inputState.Clear()
				channel.Write([]byte(dynamicPrompt))

			default:
				// 如果正在处理AI请求，忽略字符输入
				if isProcessing {
					continue
				}
				// 处理所有可打印字符，包括中文
				if r >= 32 || (r > 127 && utf8.ValidRune(r)) {
					inputState.InsertRune(r)
					refreshLine(channel, inputState, dynamicPrompt)
				}
			}
		}
	}
}
