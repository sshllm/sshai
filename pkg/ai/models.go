package ai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/ssh"

	"sshai/pkg/config"
	"sshai/pkg/i18n"
	"sshai/pkg/models"
)

// ModelInfo 模型信息别名，方便使用
type ModelInfo = models.ModelInfo

// ModelCache 模型缓存结构体
type ModelCache struct {
	models    []ModelInfo
	cacheTime time.Time
	mutex     sync.RWMutex
}

// 全局模型缓存实例
var modelCache = &ModelCache{}

// 缓存有效期（5分钟）
const cacheExpiration = 5 * time.Minute

// GetAvailableModels 获取可用模型列表（带缓存）
func GetAvailableModels() ([]ModelInfo, error) {
	// 先检查缓存
	modelCache.mutex.RLock()
	if time.Since(modelCache.cacheTime) < cacheExpiration && len(modelCache.models) > 0 {
		// 缓存有效，直接返回缓存的模型列表
		cachedModels := make([]ModelInfo, len(modelCache.models))
		copy(cachedModels, modelCache.models)
		modelCache.mutex.RUnlock()
		return cachedModels, nil
	}
	modelCache.mutex.RUnlock()

	// 缓存无效或为空，需要重新获取
	return fetchAndCacheModels()
}

// fetchAndCacheModels 从远程获取模型列表并缓存
func fetchAndCacheModels() ([]ModelInfo, error) {
	cfg := config.Get()

	req, err := http.NewRequest("GET", cfg.API.BaseURL+"/models", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.API.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var modelsResp models.ModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, err
	}

	// 更新缓存
	modelCache.mutex.Lock()
	modelCache.models = modelsResp.Data
	modelCache.cacheTime = time.Now()
	modelCache.mutex.Unlock()

	return modelsResp.Data, nil
}

// ClearModelCache 清空模型缓存（用于测试或强制刷新）
func ClearModelCache() {
	modelCache.mutex.Lock()
	modelCache.models = nil
	modelCache.cacheTime = time.Time{}
	modelCache.mutex.Unlock()
}

// GetCacheInfo 获取缓存信息（用于调试）
func GetCacheInfo() (int, time.Duration, bool) {
	modelCache.mutex.RLock()
	defer modelCache.mutex.RUnlock()

	count := len(modelCache.models)
	age := time.Since(modelCache.cacheTime)
	valid := age < cacheExpiration && count > 0

	return count, age, valid
}

// MatchModelsByUsername 根据用户名匹配模型
func MatchModelsByUsername(models []ModelInfo, username string) []ModelInfo {
	if username == "" {
		return nil
	}

	var matches []ModelInfo
	usernameLower := strings.ToLower(username)

	for _, model := range models {
		modelLower := strings.ToLower(model.ID)
		if strings.Contains(modelLower, usernameLower) || strings.Contains(usernameLower, modelLower) {
			matches = append(matches, model)
		}
	}

	return matches
}

// SelectModelByUsername 根据用户名选择模型
func SelectModelByUsername(channel ssh.Channel, models []ModelInfo, username string) string {
	cfg := config.Get()

	// 尝试根据用户名匹配模型
	matchedModels := MatchModelsByUsername(models, username)

	if len(matchedModels) == 1 {
		// 找到唯一匹配的模型
		selectedModel := matchedModels[0].ID
		channel.Write([]byte(fmt.Sprintf(i18n.T("model.auto_selected", username, selectedModel) + "\r\n")))
		return selectedModel
	} else if len(matchedModels) > 1 {
		// 找到多个匹配的模型，让用户选择
		channel.Write([]byte(fmt.Sprintf(i18n.T("model.multiple_matches", username) + "\r\n")))
		return showModelSelection(channel, matchedModels, username, models)
	} else {
		// 没有找到匹配的模型
		if len(models) > 0 {
			channel.Write([]byte(fmt.Sprintf(i18n.T("model.no_matches", username) + "\r\n")))
			return showModelSelection(channel, models, username, models)
		} else {
			channel.Write([]byte(i18n.T("model.no_available") + "\r\n"))
			return cfg.API.DefaultModel
		}
	}
}

// showModelSelection 显示模型选择界面
func showModelSelection(channel ssh.Channel, models []ModelInfo, username string, allModels []ModelInfo) string {
	cfg := config.Get()

	if len(models) == 1 {
		selectedModel := models[0].ID
		channel.Write([]byte(fmt.Sprintf(i18n.T("model.auto_only", selectedModel) + "\r\n")))
		return selectedModel
	}

	// 显示模型列表
	for i, model := range models {
		channel.Write([]byte(fmt.Sprintf("%d. %s\r\n", i+1, model.ID)))
	}
	channel.Write([]byte(i18n.T("model.select_prompt")))

	// 处理用户输入
	var inputBuffer []byte
	buffer := make([]byte, 1024)

	for {
		n, err := channel.Read(buffer)
		if err != nil {
			return cfg.API.DefaultModel
		}

		data := buffer[:n]

		// 处理可能的不完整UTF-8序列
		for len(data) > 0 {
			r, size := utf8.DecodeRune(data)
			if r == utf8.RuneError && size == 1 {
				// 可能是不完整的UTF-8序列，等待更多数据
				break
			}

			data = data[size:]

			switch r {
			case 13: // Enter键
				if len(inputBuffer) > 0 {
					input := strings.TrimSpace(string(inputBuffer))
					channel.Write([]byte("\r\n"))

					// 解析用户输入的数字
					if choice, err := strconv.Atoi(input); err == nil && choice >= 1 && choice <= len(models) {
						selectedModel := models[choice-1].ID
						channel.Write([]byte(fmt.Sprintf(i18n.T("model.selected", selectedModel) + "\r\n")))
						return selectedModel
					} else {
						channel.Write([]byte(i18n.T("model.invalid_choice")))
						inputBuffer = nil
						continue
					}
				}

			case 127, 8: // Backspace/Delete键
				if len(inputBuffer) > 0 {
					// 处理UTF-8字符的删除
					for len(inputBuffer) > 0 {
						inputBuffer = inputBuffer[:len(inputBuffer)-1]
						if utf8.Valid(inputBuffer) {
							break
						}
					}

					// 发送退格序列来清除字符
					channel.Write([]byte("\b \b"))
				}

			case 3: // Ctrl+C
				channel.Write([]byte("\r\n^C\r\n"))
				return cfg.API.DefaultModel

			default:
				// 处理数字输入
				if r >= '0' && r <= '9' {
					runeBytes := make([]byte, utf8.RuneLen(r))
					utf8.EncodeRune(runeBytes, r)
					inputBuffer = append(inputBuffer, runeBytes...)
					channel.Write(runeBytes) // 回显字符
				}
			}
		}
	}
}
