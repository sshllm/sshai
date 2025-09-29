package ssh

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
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

// ConversationHistory 对话历史结构体
type ConversationHistory struct {
	messages []ConversationMessage
	mutex    sync.RWMutex
}

// ConversationMessage 对话消息结构体
type ConversationMessage struct {
	Timestamp time.Time
	Role      string // "system", "user", "assistant"
	Content   string
}

// NewConversationHistory 创建新的对话历史
func NewConversationHistory() *ConversationHistory {
	return &ConversationHistory{
		messages: make([]ConversationMessage, 0),
	}
}

// AddMessage 添加消息到对话历史
func (h *ConversationHistory) AddMessage(role, content string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	h.messages = append(h.messages, ConversationMessage{
		Timestamp: time.Now(),
		Role:      role,
		Content:   content,
	})
}

// GetMessages 获取所有消息
func (h *ConversationHistory) GetMessages() []ConversationMessage {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	// 返回副本
	messages := make([]ConversationMessage, len(h.messages))
	copy(messages, h.messages)
	return messages
}

// Clear 清空对话历史
func (h *ConversationHistory) Clear() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	h.messages = h.messages[:0]
}

// CustomCommand 自定义命令结构体
type CustomCommand struct {
	Name        string
	Description string
	Handler     func(channel ssh.Channel, assistant *ai.Assistant, args []string, conversationHistory *ConversationHistory, dynamicPrompt string)
}

// getCustomCommands 获取自定义命令列表
func getCustomCommands() map[string]CustomCommand {
	return map[string]CustomCommand{
		"/help": {
			Name:        "/help",
			Description: "显示帮助信息和可用命令",
			Handler:     handleHelpCommand,
		},
		"/new": {
			Name:        "/new",
			Description: "清空上下文，开始新对话",
			Handler:     handleNewCommand,
		},
		"/history": {
			Name:        "/history",
			Description: "查看当前对话的历史记录",
			Handler:     handleHistoryCommand,
		},
		"/clear": {
			Name:        "/clear",
			Description: "清空屏幕",
			Handler:     handleClearCommand,
		},
		"/model": {
			Name:        "/model",
			Description: "切换AI模型",
			Handler:     handleModelCommand,
		},
	}
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

// ResponseCapture 用于捕获AI响应内容的包装器
type ResponseCapture struct {
	originalChannel ssh.Channel
	content         strings.Builder
}

// Write 实现ssh.Channel接口，同时捕获内容
func (rc *ResponseCapture) Write(data []byte) (int, error) {
	// 写入原始channel
	n, err := rc.originalChannel.Write(data)
	
	// 捕获内容（去除ANSI转义序列）
	cleanData := removeANSISequences(string(data))
	rc.content.WriteString(cleanData)
	
	return n, err
}

// Read 实现ssh.Channel接口，直接转发到原始channel
func (rc *ResponseCapture) Read(data []byte) (int, error) {
	return rc.originalChannel.Read(data)
}

// Close 实现ssh.Channel接口，直接转发到原始channel
func (rc *ResponseCapture) Close() error {
	return rc.originalChannel.Close()
}

// CloseWrite 实现ssh.Channel接口，直接转发到原始channel
func (rc *ResponseCapture) CloseWrite() error {
	return rc.originalChannel.CloseWrite()
}

// SendRequest 实现ssh.Channel接口，直接转发到原始channel
func (rc *ResponseCapture) SendRequest(name string, wantReply bool, payload []byte) (bool, error) {
	return rc.originalChannel.SendRequest(name, wantReply, payload)
}

// Stderr 实现ssh.Channel接口，直接转发到原始channel
func (rc *ResponseCapture) Stderr() io.ReadWriter {
	return rc.originalChannel.Stderr()
}

// removeANSISequences 移除ANSI转义序列
func removeANSISequences(text string) string {
	// 简单的ANSI转义序列移除
	result := strings.Builder{}
	inEscape := false
	
	for i, r := range text {
		if r == '\033' && i+1 < len(text) && text[i+1] == '[' {
			inEscape = true
			continue
		}
		
		if inEscape {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEscape = false
			}
			continue
		}
		
		result.WriteRune(r)
	}
	
	return result.String()
}

// handleTabCompletion 处理Tab补全
func handleTabCompletion(channel ssh.Channel, inputState *InputState, tabState *struct {
	isActive    bool
	prefix      string
	matches     []string
	currentIdx  int
}, dynamicPrompt string) {
	currentInput := inputState.String()
	
	// 如果不是第一次按Tab，或者输入发生了变化，重新计算匹配项
	if !tabState.isActive || tabState.prefix != currentInput {
		tabState.prefix = currentInput
		tabState.matches = getCommandMatches(currentInput)
		tabState.currentIdx = 0
		tabState.isActive = true
		
		// 如果没有匹配项，直接返回
		if len(tabState.matches) == 0 {
			return
		}
		
		// 如果只有一个匹配项，直接补全
		if len(tabState.matches) == 1 {
			inputState.SetText(tabState.matches[0] + " ")
			refreshLine(channel, inputState, dynamicPrompt)
			tabState.isActive = false
			return
		}
	}
	
	// 多个匹配项，循环显示
	if len(tabState.matches) > 1 {
		// 显示当前匹配项
		currentMatch := tabState.matches[tabState.currentIdx]
		inputState.SetText(currentMatch)
		refreshLine(channel, inputState, dynamicPrompt)
		
		// 移动到下一个匹配项
		tabState.currentIdx = (tabState.currentIdx + 1) % len(tabState.matches)
	}
}

// getCommandMatches 获取命令匹配项
func getCommandMatches(prefix string) []string {
	var matches []string
	customCommands := getCustomCommands()
	
	for cmd := range customCommands {
		if strings.HasPrefix(cmd, prefix) {
			matches = append(matches, cmd)
		}
	}
	
	return matches
}

// handleCustomCommand 处理自定义命令
func handleCustomCommand(channel ssh.Channel, assistant *ai.Assistant, input string, conversationHistory *ConversationHistory, dynamicPrompt string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}
	
	command := parts[0]
	args := parts[1:]
	customCommands := getCustomCommands()
	
	// 添加命令到对话历史
	conversationHistory.AddMessage("user", input)
	
	if cmd, exists := customCommands[command]; exists {
		cmd.Handler(channel, assistant, args, conversationHistory, dynamicPrompt)
	} else {
		channel.Write([]byte(fmt.Sprintf("未知命令: %s\r\n", command)))
		channel.Write([]byte("输入 /help 查看可用命令\r\n"))
	}
	
	// 显示提示符
	channel.Write([]byte(dynamicPrompt))
}

// handleHelpCommand 处理help命令
func handleHelpCommand(channel ssh.Channel, assistant *ai.Assistant, args []string, conversationHistory *ConversationHistory, dynamicPrompt string) {
	channel.Write([]byte(ui.BrightCyanText("📋 可用的自定义命令:\r\n\r\n")))
	
	customCommands := getCustomCommands()
	// 按字母顺序显示命令
	commands := []string{"/clear", "/help", "/history", "/model", "/new"}
	
	for _, cmdName := range commands {
		if cmd, exists := customCommands[cmdName]; exists {
			channel.Write([]byte(fmt.Sprintf("  %s - %s\r\n", 
				ui.BrightYellowText(cmd.Name), 
				cmd.Description)))
		}
	}
	
	channel.Write([]byte("\r\n"))
	channel.Write([]byte(ui.BrightGreenText("💡 提示:\r\n")))
	channel.Write([]byte("  • 使用 Tab 键可以自动补全命令\r\n"))
	channel.Write([]byte("  • 输入 'exit' 或 'quit' 退出程序\r\n"))
	channel.Write([]byte("  • 直接输入消息与AI对话\r\n\r\n"))
	
	// 添加系统消息到对话历史
	conversationHistory.AddMessage("system", "显示了帮助信息")
}

// handleNewCommand 处理new命令
func handleNewCommand(channel ssh.Channel, assistant *ai.Assistant, args []string, conversationHistory *ConversationHistory, dynamicPrompt string) {
	assistant.ClearContext()
	conversationHistory.Clear()
	channel.Write([]byte(ui.BrightGreenText("✅ 对话上下文已清空，开始新对话\r\n\r\n")))
	
	// 添加系统消息到对话历史
	conversationHistory.AddMessage("system", "清空了对话上下文，开始新对话")
}

// handleHistoryCommand 处理history命令
func handleHistoryCommand(channel ssh.Channel, assistant *ai.Assistant, args []string, conversationHistory *ConversationHistory, dynamicPrompt string) {
	messages := conversationHistory.GetMessages()
	
	if len(messages) == 0 {
		channel.Write([]byte(ui.BrightYellowText("📝 当前对话历史为空\r\n\r\n")))
		return
	}
	
	channel.Write([]byte(ui.BrightCyanText("📝 对话历史记录:\r\n\r\n")))
	
	for i, msg := range messages {
		// 格式化时间
		timeStr := msg.Timestamp.Format("15:04:05")
		
		// 根据角色设置不同颜色
		var roleColor, roleIcon string
		switch msg.Role {
		case "system":
			roleColor = ui.BrightMagentaText("系统")
			roleIcon = "🔧"
		case "user":
			roleColor = ui.BrightGreenText("用户")
			roleIcon = "👤"
		case "assistant":
			roleColor = ui.BrightBlueText("助手")
			roleIcon = "🤖"
		default:
			roleColor = ui.BrightWhiteText(msg.Role)
			roleIcon = "❓"
		}
		
		// 显示消息头
		channel.Write([]byte(fmt.Sprintf("%s [%s] %s %s:\r\n", 
			roleIcon,
			ui.BrightWhiteText(timeStr),
			roleColor,
			ui.BrightWhiteText(fmt.Sprintf("#%d", i+1)))))
		
		// 显示消息内容（限制长度）
		content := msg.Content
		if len(content) > 200 {
			content = content[:200] + "..."
		}
		
		// 将换行符替换为\r\n，并添加缩进
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			channel.Write([]byte(fmt.Sprintf("  %s\r\n", line)))
		}
		channel.Write([]byte("\r\n"))
	}
	
	channel.Write([]byte(ui.BrightCyanText(fmt.Sprintf("📊 总计: %d 条消息\r\n\r\n", len(messages)))))
	
	// 添加系统消息到对话历史
	conversationHistory.AddMessage("system", fmt.Sprintf("查看了对话历史，共%d条消息", len(messages)))
}

// handleClearCommand 处理clear命令
func handleClearCommand(channel ssh.Channel, assistant *ai.Assistant, args []string, conversationHistory *ConversationHistory, dynamicPrompt string) {
	// 发送清屏命令
	channel.Write([]byte("\033[2J\033[H"))
	
	// 重新显示banner和欢迎信息
	banner := ui.GenerateBanner()
	lines := strings.Split(banner, "\n")
	for _, line := range lines {
		channel.Write([]byte(line + "\r\n"))
	}
	channel.Write([]byte("\r\n"))
	
	channel.Write([]byte(ui.BrightGreenText("🧹 屏幕已清空\r\n\r\n")))
	
	// 添加系统消息到对话历史
	conversationHistory.AddMessage("system", "清空了屏幕")
}

// handleModelCommand 处理model命令
func handleModelCommand(channel ssh.Channel, assistant *ai.Assistant, args []string, conversationHistory *ConversationHistory, dynamicPrompt string) {
	channel.Write([]byte(ui.BrightCyanText("🔄 正在加载模型列表...\r\n")))
	
	// 获取可用模型
	models, err := ai.GetAvailableModels()
	if err != nil {
		channel.Write([]byte(ui.BrightRedText(fmt.Sprintf("❌ 获取模型列表失败: %v\r\n\r\n", err))))
		conversationHistory.AddMessage("system", fmt.Sprintf("获取模型列表失败: %v", err))
		return
	}
	
	if len(models) == 0 {
		channel.Write([]byte(ui.BrightYellowText("⚠️  没有找到可用的模型\r\n\r\n")))
		conversationHistory.AddMessage("system", "没有找到可用的模型")
		return
	}
	
	// 显示当前模型
	currentModel := assistant.GetCurrentModel()
	channel.Write([]byte(fmt.Sprintf("当前模型: %s\r\n\r\n", ui.BrightYellowText(currentModel))))
	
	// 显示模型选择界面
	selectedModel := showModelSelectionForCommand(channel, models)
	if selectedModel != "" && selectedModel != currentModel {
		assistant.SetModel(selectedModel)
		channel.Write([]byte(fmt.Sprintf("✅ 已切换到模型: %s\r\n\r\n", ui.BrightGreenText(selectedModel))))
		conversationHistory.AddMessage("system", fmt.Sprintf("切换到模型: %s", selectedModel))
	} else if selectedModel == currentModel {
		channel.Write([]byte(ui.BrightYellowText("ℹ️  模型未更改\r\n\r\n")))
	} else {
		channel.Write([]byte(ui.BrightYellowText("❌ 模型切换已取消\r\n\r\n")))
	}
}

// showModelSelectionForCommand 为命令显示模型选择界面
func showModelSelectionForCommand(channel ssh.Channel, models []ai.ModelInfo) string {
	cfg := config.Get()
	
	// 显示模型列表
	channel.Write([]byte(ui.BrightCyanText("📋 可用模型:\r\n")))
	for i, model := range models {
		channel.Write([]byte(fmt.Sprintf("%s. %s\r\n", 
			ui.BrightWhiteText(fmt.Sprintf("%d", i+1)), 
			ui.BrightYellowText(model.ID))))
	}
	channel.Write([]byte(fmt.Sprintf("\r\n%s", ui.BrightCyanText("请选择模型 (输入数字，按 Ctrl+C 取消): "))))
	
	// 处理用户输入
	var inputBuffer []byte
	buffer := make([]byte, 1024)
	
	for {
		n, err := channel.Read(buffer)
		if err != nil {
			return cfg.API.DefaultModel
		}
		
		data := buffer[:n]
		
		for len(data) > 0 {
			r, size := utf8.DecodeRune(data)
			if r == utf8.RuneError && size == 1 {
				break
			}
			
			data = data[size:]
			
			switch r {
			case 13: // Enter键
				if len(inputBuffer) > 0 {
					input := strings.TrimSpace(string(inputBuffer))
					channel.Write([]byte("\r\n"))
					
					if choice, err := strconv.Atoi(input); err == nil && choice >= 1 && choice <= len(models) {
						return models[choice-1].ID
					} else {
						channel.Write([]byte(ui.BrightRedText("❌ 无效选择，请重新输入: ")))
						inputBuffer = nil
						continue
					}
				}
				
			case 127, 8: // Backspace/Delete键
				if len(inputBuffer) > 0 {
					for len(inputBuffer) > 0 {
						inputBuffer = inputBuffer[:len(inputBuffer)-1]
						if utf8.Valid(inputBuffer) {
							break
						}
					}
					channel.Write([]byte("\b \b"))
				}
				
			case 3: // Ctrl+C
				channel.Write([]byte("\r\n"))
				return ""
				
			default:
				if r >= '0' && r <= '9' {
					runeBytes := make([]byte, utf8.RuneLen(r))
					utf8.EncodeRune(runeBytes, r)
					inputBuffer = append(inputBuffer, runeBytes...)
					channel.Write(runeBytes)
				}
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

	// 直接处理内容并获取AI响应，不显示动画效果和工具调用信息
	assistant.ProcessMessageWithFullOptions(prompt, channel, interrupt, false, false)
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

	// 创建对话历史
	conversationHistory := NewConversationHistory()
	
	// 处理用户输入
	handleUserInput(channel, assistant, dynamicPrompt, conversationHistory)
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

	// 直接处理命令并获取AI响应，不显示动画效果和工具调用信息
	assistant.ProcessMessageWithFullOptions(fullPrompt, channel, interrupt, false, false)

	// 添加换行符结束
	channel.Write([]byte("\r\n"))
}

// handleUserInput 处理用户输入
func handleUserInput(channel ssh.Channel, assistant *ai.Assistant, dynamicPrompt string, conversationHistory *ConversationHistory) {
	buffer := make([]byte, 1024)
	history := NewCommandHistory()
	inputState := NewInputState()
	var escapeSequence []byte      // 用于处理ANSI转义序列
	var currentInterrupt chan bool // 当前正在使用的中断通道
	var isProcessing bool          // 标记是否正在处理AI请求
	var tabCompletionState struct {
		isActive    bool
		prefix      string
		matches     []string
		currentIdx  int
	}

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
			case 9: // Tab键 - 自动补全
				// 如果正在处理AI请求，忽略tab补全
				if isProcessing {
					continue
				}
				
				currentInput := inputState.String()
				
				// 检查是否是自定义命令的补全
				if strings.HasPrefix(currentInput, "/") {
					handleTabCompletion(channel, inputState, &tabCompletionState, dynamicPrompt)
				}
				
			case 13: // Enter键
				// 重置tab补全状态
				tabCompletionState.isActive = false
				
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

				// 检查是否是自定义命令
				if strings.HasPrefix(input, "/") {
					handleCustomCommand(channel, assistant, input, conversationHistory, dynamicPrompt)
				} else if input == "exit" || input == "quit" {
					channel.Write([]byte(i18n.T("user.exit") + "\r\n"))
					return
				} else if input != "" {
					// 添加用户消息到对话历史
					conversationHistory.AddMessage("user", input)
					
					// 设置处理状态
					isProcessing = true

					// 创建新的中断通道用于这次AI请求
					currentInterrupt = make(chan bool)

					// 异步处理AI请求，这样Ctrl+C可以在处理过程中被响应
					go func(userInput string, interruptCh chan bool) {
						// 创建一个包装的channel来捕获AI响应
						responseCapture := &ResponseCapture{
							originalChannel: channel,
							content:         strings.Builder{},
						}
						
						assistant.ProcessMessage(userInput, responseCapture, interruptCh)
						
						// 添加AI响应到对话历史
						if responseCapture.content.Len() > 0 {
							conversationHistory.AddMessage("assistant", responseCapture.content.String())
						}
						
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
				// 重置tab补全状态（当用户输入其他字符时）
				tabCompletionState.isActive = false
				
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
