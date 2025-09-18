package auth

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"

	"sshai/pkg/config"
)

// AuthorizedKeysManager SSH公钥管理器
type AuthorizedKeysManager struct {
	keys []ssh.PublicKey
}

// NewAuthorizedKeysManager 创建新的SSH公钥管理器
func NewAuthorizedKeysManager() (*AuthorizedKeysManager, error) {
	cfg := config.Get()
	manager := &AuthorizedKeysManager{
		keys: make([]ssh.PublicKey, 0),
	}

	// 加载配置中的公钥列表
	for _, keyStr := range cfg.Auth.AuthorizedKeys {
		if keyStr = strings.TrimSpace(keyStr); keyStr != "" {
			if err := manager.addKeyFromString(keyStr); err != nil {
				log.Printf("警告：无法解析公钥: %v", err)
			}
		}
	}

	// 加载公钥文件（如果配置了）
	if cfg.Auth.AuthorizedKeysFile != "" {
		if err := manager.loadKeysFromFile(cfg.Auth.AuthorizedKeysFile); err != nil {
			log.Printf("警告：无法加载公钥文件 %s: %v", cfg.Auth.AuthorizedKeysFile, err)
		}
	}

	log.Printf("SSH公钥管理器初始化完成，共加载 %d 个公钥", len(manager.keys))
	return manager, nil
}

// addKeyFromString 从字符串添加公钥
func (m *AuthorizedKeysManager) addKeyFromString(keyStr string) error {
	// 解析SSH公钥字符串
	publicKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(keyStr))
	if err != nil {
		return fmt.Errorf("解析公钥失败: %v", err)
	}

	m.keys = append(m.keys, publicKey)
	return nil
}

// loadKeysFromFile 从文件加载公钥
func (m *AuthorizedKeysManager) loadKeysFromFile(filePath string) error {
	// 展开用户主目录路径
	if strings.HasPrefix(filePath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("无法获取用户主目录: %v", err)
		}
		filePath = strings.Replace(filePath, "~", homeDir, 1)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("无法打开公钥文件: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if err := m.addKeyFromString(line); err != nil {
			log.Printf("警告：公钥文件第 %d 行解析失败: %v", lineNum, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取公钥文件失败: %v", err)
	}

	return nil
}

// VerifyPublicKey 验证公钥是否在授权列表中
func (m *AuthorizedKeysManager) VerifyPublicKey(key ssh.PublicKey) bool {
	keyData := key.Marshal()
	keyType := key.Type()

	for _, authorizedKey := range m.keys {
		if authorizedKey.Type() == keyType {
			if string(authorizedKey.Marshal()) == string(keyData) {
				return true
			}
		}
	}

	return false
}

// GetKeyCount 获取已加载的公钥数量
func (m *AuthorizedKeysManager) GetKeyCount() int {
	return len(m.keys)
}

// IsEnabled 检查SSH公钥认证是否启用
// 只有在设置了密码认证时，SSH公钥认证才生效
func IsEnabled() bool {
	cfg := config.Get()
	return cfg.Auth.Password != "" && (len(cfg.Auth.AuthorizedKeys) > 0 || cfg.Auth.AuthorizedKeysFile != "")
}
