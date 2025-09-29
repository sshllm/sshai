package ssh

import (
	"testing"
	"time"
)

func TestConversationHistory(t *testing.T) {
	history := NewConversationHistory()
	
	// 测试添加消息
	history.AddMessage("user", "Hello")
	history.AddMessage("assistant", "Hi there!")
	history.AddMessage("system", "Test message")
	
	// 验证消息数量
	messages := history.GetMessages()
	if len(messages) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(messages))
	}
	
	// 验证消息内容
	if messages[0].Role != "user" || messages[0].Content != "Hello" {
		t.Errorf("First message incorrect: %+v", messages[0])
	}
	
	if messages[1].Role != "assistant" || messages[1].Content != "Hi there!" {
		t.Errorf("Second message incorrect: %+v", messages[1])
	}
	
	if messages[2].Role != "system" || messages[2].Content != "Test message" {
		t.Errorf("Third message incorrect: %+v", messages[2])
	}
	
	// 测试清空历史
	history.Clear()
	messages = history.GetMessages()
	if len(messages) != 0 {
		t.Errorf("Expected 0 messages after clear, got %d", len(messages))
	}
}

func TestGetCustomCommands(t *testing.T) {
	commands := getCustomCommands()
	
	// 验证所有必需的命令都存在
	expectedCommands := []string{"/help", "/new", "/history", "/clear", "/model"}
	
	for _, cmdName := range expectedCommands {
		if _, exists := commands[cmdName]; !exists {
			t.Errorf("Expected command %s not found", cmdName)
		}
	}
	
	// 验证命令数量
	if len(commands) != len(expectedCommands) {
		t.Errorf("Expected %d commands, got %d", len(expectedCommands), len(commands))
	}
}

func TestGetCommandMatches(t *testing.T) {
	// 测试完全匹配
	matches := getCommandMatches("/help")
	if len(matches) != 1 || matches[0] != "/help" {
		t.Errorf("Expected [/help], got %v", matches)
	}
	
	// 测试前缀匹配
	matches = getCommandMatches("/h")
	expectedMatches := []string{"/help", "/history"}
	if len(matches) != len(expectedMatches) {
		t.Errorf("Expected %d matches for '/h', got %d", len(expectedMatches), len(matches))
	}
	
	// 测试无匹配
	matches = getCommandMatches("/xyz")
	if len(matches) != 0 {
		t.Errorf("Expected no matches for '/xyz', got %v", matches)
	}
	
	// 测试空输入
	matches = getCommandMatches("/")
	if len(matches) != 5 { // 应该返回所有命令
		t.Errorf("Expected 5 matches for '/', got %d", len(matches))
	}
}

func TestMessageTimestamp(t *testing.T) {
	history := NewConversationHistory()
	
	before := time.Now()
	history.AddMessage("user", "Test message")
	after := time.Now()
	
	messages := history.GetMessages()
	if len(messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(messages))
	}
	
	msgTime := messages[0].Timestamp
	if msgTime.Before(before) || msgTime.After(after) {
		t.Errorf("Message timestamp %v is not between %v and %v", msgTime, before, after)
	}
}