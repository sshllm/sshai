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

// readStdinInput 尝试读取stdin输入
func readStdinInput(channel ssh.Channel) (string, bool) {
	// 使用goroutine和channel来实现超时读取
	resultChan := make(chan string, 1)
	errorChan := make(chan error, 1)

	go func() {
		buffer := make([]byte, 8192) // 增大缓冲区以支持较大文件
		var content strings.Builder

		for {
			n, err := channel.Read(buffer)
			if err != nil {
				errorChan <- err
				return
			}

			if n > 0 {
				content.Write(buffer[:n])
				// 如果读取的数据小于缓冲区大小，可能已经读完
				if n < len(buffer) {
					break
				}
			}
		}

		resultChan <- content.String()
	}()

	// 等待结果或超时
	select {
	case result := <-resultChan:
		return result, len(result) > 0
	case <-errorChan:
		return "", false
	case <-time.After(300 * time.Millisecond):
		// 超时，返回空结果
		return "", false
	}
}

// handleStdinCommand 处理通过stdin传入的内容
func handleStdinCommand(channel ssh.Channel, username, content string) {
	log.Printf("处理stdin内容，用户: %s，内容长度: %d", username, len(content))

	cfg := config.Get()

	// 显示接收到的内容信息
	channel.Write([]byte(fmt.Sprintf("接收到输入内容（%d字符）\r\n", len(content))))

	// 获取并选择模型（简化版本，不显示加载过程）
	models, err := ai.GetAvailableModels()
	if err != nil {
		log.Printf("获取模型失败: %v", err)
		models = []ai.ModelInfo{{ID: cfg.API.DefaultModel}}
	}

	// 根据用户名匹配模型（exec模式下不需要交互）
	selectedModel := cfg.API.DefaultModel

	// 尝试根据用户名匹配模型
	for _, model := range models {
		if strings.Contains(strings.ToLower(username), strings.ToLower(model.ID)) {
			selectedModel = model.ID
			break
		}
	}

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

	// 直接处理内容并获取AI响应
	assistant.ProcessMessage(prompt, channel, interrupt)
}

// HandleSession 处理SSH会话
func HandleSession(channel ssh.Channel, requests <-chan *ssh.Request, username string) {
	defer channel.Close()

	var execCommand string
	isExecMode := false
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

	// 检查是否有stdin输入（无论是exec还是shell模式）
	stdinContent, hasStdin := readStdinInput(channel)

	if hasStdin && stdinContent != "" {
		// 如果有stdin输入，将其作为主要内容处理
		log.Printf("检测到stdin输入，长度: %d", len(stdinContent))
		handleStdinCommand(channel, username, stdinContent)
		return
	}

	// 如果是执行模式且有命令，处理exec命令
	if isExec && execCommand != "" {
		handleExecCommand(channel, username, execCommand)
		return
	}

	// 发送登录成功消息（无论是否需要密码认证）
	if cfg.Auth.LoginSuccessMsg != "" {
		// 处理多行消息的换行
		lines := strings.Split(cfg.Auth.LoginSuccessMsg, "\n")
		for _, line := range lines {
			channel.Write([]byte(line + "\r\n"))
		}
		channel.Write([]byte("\r\n")) // 额外的空行分隔
	}

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

	// 生成动态提示符
	dynamicPrompt := fmt.Sprintf(cfg.Server.PromptTemplate, selectedModel)
	channel.Write([]byte("\r\n" + dynamicPrompt))

	// 处理用户输入
	handleUserInput(channel, assistant, dynamicPrompt)
}

// handleExecCommand 处理执行命令模式
func handleExecCommand(channel ssh.Channel, username, command string) {
	cfg := config.Get()

	// 显示执行的命令
	channel.Write([]byte(fmt.Sprintf("执行命令: %s\r\n\r\n", command)))

	// 获取并选择模型（简化版本，不显示加载过程）
	models, err := ai.GetAvailableModels()
	if err != nil {
		log.Printf("获取模型失败: %v", err)
		models = []ai.ModelInfo{{ID: cfg.API.DefaultModel}}
	}

	// 根据用户名匹配模型（exec模式下不需要交互）
	selectedModel := cfg.API.DefaultModel

	// 尝试根据用户名匹配模型
	for _, model := range models {
		if strings.Contains(strings.ToLower(username), strings.ToLower(model.ID)) {
			selectedModel = model.ID
			break
		}
	}

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

	// 直接处理命令并获取AI响应
	assistant.ProcessMessage(fullPrompt, channel, interrupt)

	// 添加换行符结束
	channel.Write([]byte("\r\n"))
}

// handleUserInput 处理用户输入
func handleUserInput(channel ssh.Channel, assistant *ai.Assistant, dynamicPrompt string) {
	buffer := make([]byte, 1024)
	interrupt := make(chan bool, 1)
	history := NewCommandHistory()
	inputState := NewInputState()
	var escapeSequence []byte // 用于处理ANSI转义序列

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
							channel.Write([]byte("\033[C")) // 向右移动光标
						}
					case 68: // 左方向键 ESC[D
						if inputState.MoveCursorLeft() {
							channel.Write([]byte("\033[D")) // 向左移动光标
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
					// 处理AI请求
					assistant.ProcessMessage(input, channel, interrupt)
				}
				// 清空输入状态并显示提示符
				inputState.Clear()
				channel.Write([]byte(dynamicPrompt))

			case 127, 8: // Backspace/Delete键
				if inputState.DeleteRune() {
					refreshLine(channel, inputState, dynamicPrompt)
				}

			case 1: // Ctrl+A - 移动到行首
				inputState.MoveCursorToStart()
				refreshLine(channel, inputState, dynamicPrompt)

			case 5: // Ctrl+E - 移动到行尾
				inputState.MoveCursorToEnd()
				refreshLine(channel, inputState, dynamicPrompt)

			case 3: // Ctrl+C
				// 发送中断信号
				select {
				case interrupt <- true:
				default:
				}
				channel.Write([]byte("\r\n^C\r\n"))
				inputState.Clear()
				channel.Write([]byte(dynamicPrompt))

			default:
				// 处理所有可打印字符，包括中文
				if r >= 32 || (r > 127 && utf8.ValidRune(r)) {
					inputState.InsertRune(r)
					refreshLine(channel, inputState, dynamicPrompt)
				}
			}
		}
	}
}
