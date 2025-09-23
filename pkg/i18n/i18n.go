package i18n

import (
	"embed"
	"fmt"
	"sync"

	"gopkg.in/yaml.v2"
)

//go:embed lang/*.yaml
var langFS embed.FS

// Language 语言类型
type Language string

const (
	// LanguageZhCN 简体中文
	LanguageZhCN Language = "zh-cn"
	// LanguageEnUS 英文
	LanguageEnUS Language = "en-us"
)

// LanguageMessages 语言消息结构体
type LanguageMessages struct {
	Server struct {
		Starting string `yaml:"starting"`
		Started  string `yaml:"started"`
		Error    string `yaml:"error"`
		Shutdown string `yaml:"shutdown"`
	} `yaml:"server"`

	Connection struct {
		New         string `yaml:"new"`
		AuthFailed  string `yaml:"auth_failed"`
		Established string `yaml:"established"`
		Closed      string `yaml:"closed"`
	} `yaml:"connection"`

	Auth struct {
		PasswordPrompt    string `yaml:"password_prompt"`
		PasswordIncorrect string `yaml:"password_incorrect"`
		LoginSuccess      string `yaml:"login_success"`
		LoginFailed       string `yaml:"login_failed"`
	} `yaml:"auth"`

	Model struct {
		Loading         string `yaml:"loading"`
		Loaded          string `yaml:"loaded"`
		Error           string `yaml:"error"`
		AutoSelected    string `yaml:"auto_selected"`
		MultipleMatches string `yaml:"multiple_matches"`
		NoMatches       string `yaml:"no_matches"`
		NoAvailable     string `yaml:"no_available"`
		AutoOnly        string `yaml:"auto_only"`
		SelectPrompt    string `yaml:"select_prompt"`
		Selected        string `yaml:"selected"`
		InvalidChoice   string `yaml:"invalid_choice"`
		CacheHit        string `yaml:"cache_hit"`
		CacheMiss       string `yaml:"cache_miss"`
	} `yaml:"model"`

	AI struct {
		Thinking         string `yaml:"thinking"`
		ThinkingProcess  string `yaml:"thinking_process"`
		ThinkingComplete string `yaml:"thinking_complete"`
		Response         string `yaml:"response"`
		Interrupted      string `yaml:"interrupted"`
		Error            string `yaml:"error"`
		RequestFailed    string `yaml:"request_failed"`
		NetworkError     string `yaml:"network_error"`
	} `yaml:"ai"`

	User struct {
		Welcome      string `yaml:"welcome"`
		Prompt       string `yaml:"prompt"`
		InputEmpty   string `yaml:"input_empty"`
		Exit         string `yaml:"exit"`
		Help         string `yaml:"help"`
		ClearContext string `yaml:"clear_context"`
	} `yaml:"user"`

	Cmd struct {
		Exit    string `yaml:"exit"`
		Help    string `yaml:"help"`
		Clear   string `yaml:"clear"`
		Unknown string `yaml:"unknown"`
	} `yaml:"cmd"`

	Error struct {
		ConfigLoad   string `yaml:"config_load"`
		KeyLoad      string `yaml:"key_load"`
		ServerStart  string `yaml:"server_start"`
		Connection   string `yaml:"connection"`
		Internal     string `yaml:"internal"`
		LangLoad     string `yaml:"lang_load"`
		LangNotFound string `yaml:"lang_not_found"`
	} `yaml:"error"`

	System struct {
		Version    string `yaml:"version"`
		ProjectURL string `yaml:"project_url"`
		ConfigFile string `yaml:"config_file"`
	} `yaml:"system"`
}

// I18n 国际化管理器
type I18n struct {
	currentLang  Language
	messages     map[Language]*LanguageMessages
	flatMessages map[Language]map[string]string
	mutex        sync.RWMutex
}

// 全局i18n实例
var globalI18n *I18n
var once sync.Once

// Init 初始化i18n系统
func Init(lang Language, langDir ...string) error {
	var err error
	once.Do(func() {
		globalI18n = &I18n{
			currentLang:  lang,
			messages:     make(map[Language]*LanguageMessages),
			flatMessages: make(map[Language]map[string]string),
		}
		err = globalI18n.loadLanguages()
		// 加载MCP翻译
		AddMCPTranslations()
	})
	return err
}

// SetLanguage 设置当前语言
func SetLanguage(lang Language) error {
	if globalI18n == nil {
		return Init(lang)
	}

	globalI18n.mutex.Lock()
	defer globalI18n.mutex.Unlock()

	// 检查语言是否已加载
	if _, exists := globalI18n.messages[lang]; !exists {
		if err := globalI18n.loadLanguage(lang); err != nil {
			return err
		}
	}

	globalI18n.currentLang = lang
	return nil
}

// GetLanguage 获取当前语言
func GetLanguage() Language {
	if globalI18n == nil {
		return LanguageZhCN
	}

	globalI18n.mutex.RLock()
	defer globalI18n.mutex.RUnlock()
	return globalI18n.currentLang
}

// T 翻译函数
func T(key string, args ...interface{}) string {
	if globalI18n == nil {
		return key
	}

	globalI18n.mutex.RLock()
	defer globalI18n.mutex.RUnlock()

	// 尝试从当前语言获取翻译
	if flatMessages, exists := globalI18n.flatMessages[globalI18n.currentLang]; exists {
		if message, exists := flatMessages[key]; exists {
			if len(args) > 0 {
				return fmt.Sprintf(message, args...)
			}
			return message
		}
	}

	// 回退到中文
	if globalI18n.currentLang != LanguageZhCN {
		if flatMessages, exists := globalI18n.flatMessages[LanguageZhCN]; exists {
			if message, exists := flatMessages[key]; exists {
				if len(args) > 0 {
					return fmt.Sprintf(message, args...)
				}
				return message
			}
		}
	}

	// 如果都没有找到，返回key
	return key
}

// loadLanguages 加载所有语言
func (i *I18n) loadLanguages() error {
	// 加载默认语言
	if err := i.loadLanguage(i.currentLang); err != nil {
		return err
	}

	// 如果当前语言不是中文，也加载中文作为回退
	if i.currentLang != LanguageZhCN {
		if err := i.loadLanguage(LanguageZhCN); err != nil {
			// 中文加载失败不是致命错误，记录但继续
			fmt.Printf("Warning: Failed to load Chinese language pack: %v\n", err)
		}
	}

	return nil
}

// loadLanguage 从嵌入文件系统加载指定语言
func (i *I18n) loadLanguage(lang Language) error {
	filename := fmt.Sprintf("lang/lang-%s.yaml", string(lang))

	// 从嵌入文件系统读取文件
	data, err := langFS.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read embedded language pack file %s: %v", filename, err)
	}

	// 解析YAML
	var messages LanguageMessages
	if err := yaml.Unmarshal(data, &messages); err != nil {
		return fmt.Errorf("failed to parse language pack file %s: %v", filename, err)
	}

	// 存储结构化消息
	i.messages[lang] = &messages

	// 创建扁平化的消息映射
	i.flatMessages[lang] = i.flattenMessages(&messages)

	return nil
}

// flattenMessages 将结构化消息转换为扁平化的key-value映射
func (i *I18n) flattenMessages(messages *LanguageMessages) map[string]string {
	flat := make(map[string]string)

	// Server
	flat["server.starting"] = messages.Server.Starting
	flat["server.started"] = messages.Server.Started
	flat["server.error"] = messages.Server.Error
	flat["server.shutdown"] = messages.Server.Shutdown

	// Connection
	flat["connection.new"] = messages.Connection.New
	flat["connection.auth_failed"] = messages.Connection.AuthFailed
	flat["connection.established"] = messages.Connection.Established
	flat["connection.closed"] = messages.Connection.Closed

	// Auth
	flat["auth.password_prompt"] = messages.Auth.PasswordPrompt
	flat["auth.password_incorrect"] = messages.Auth.PasswordIncorrect
	flat["auth.login_success"] = messages.Auth.LoginSuccess
	flat["auth.login_failed"] = messages.Auth.LoginFailed

	// Model
	flat["model.loading"] = messages.Model.Loading
	flat["model.loaded"] = messages.Model.Loaded
	flat["model.error"] = messages.Model.Error
	flat["model.auto_selected"] = messages.Model.AutoSelected
	flat["model.multiple_matches"] = messages.Model.MultipleMatches
	flat["model.no_matches"] = messages.Model.NoMatches
	flat["model.no_available"] = messages.Model.NoAvailable
	flat["model.auto_only"] = messages.Model.AutoOnly
	flat["model.select_prompt"] = messages.Model.SelectPrompt
	flat["model.selected"] = messages.Model.Selected
	flat["model.invalid_choice"] = messages.Model.InvalidChoice
	flat["model.cache_hit"] = messages.Model.CacheHit
	flat["model.cache_miss"] = messages.Model.CacheMiss

	// AI
	flat["ai.thinking"] = messages.AI.Thinking
	flat["ai.thinking_process"] = messages.AI.ThinkingProcess
	flat["ai.thinking_complete"] = messages.AI.ThinkingComplete
	flat["ai.response"] = messages.AI.Response
	flat["ai.interrupted"] = messages.AI.Interrupted
	flat["ai.error"] = messages.AI.Error
	flat["ai.request_failed"] = messages.AI.RequestFailed
	flat["ai.network_error"] = messages.AI.NetworkError

	// User
	flat["user.welcome"] = messages.User.Welcome
	flat["user.prompt"] = messages.User.Prompt
	flat["user.input_empty"] = messages.User.InputEmpty
	flat["user.exit"] = messages.User.Exit
	flat["user.help"] = messages.User.Help
	flat["user.clear_context"] = messages.User.ClearContext

	// Cmd
	flat["cmd.exit"] = messages.Cmd.Exit
	flat["cmd.help"] = messages.Cmd.Help
	flat["cmd.clear"] = messages.Cmd.Clear
	flat["cmd.unknown"] = messages.Cmd.Unknown

	// Error
	flat["error.config_load"] = messages.Error.ConfigLoad
	flat["error.key_load"] = messages.Error.KeyLoad
	flat["error.server_start"] = messages.Error.ServerStart
	flat["error.connection"] = messages.Error.Connection
	flat["error.internal"] = messages.Error.Internal
	flat["error.lang_load"] = messages.Error.LangLoad
	flat["error.lang_not_found"] = messages.Error.LangNotFound

	// System
	flat["system.version"] = messages.System.Version
	flat["system.project_url"] = messages.System.ProjectURL
	flat["system.config_file"] = messages.System.ConfigFile

	return flat
}

// GetAvailableLanguages 获取可用的语言列表（从嵌入文件系统）
func GetAvailableLanguages() []Language {
	// 返回硬编码的可用语言列表，因为我们知道嵌入了哪些语言包
	return []Language{LanguageZhCN, LanguageEnUS}
}

// ReloadLanguage 重新加载指定语言（用于热更新）
func ReloadLanguage(lang Language) error {
	if globalI18n == nil {
		return fmt.Errorf("i18n system not initialized")
	}

	globalI18n.mutex.Lock()
	defer globalI18n.mutex.Unlock()

	return globalI18n.loadLanguage(lang)
}

// GetLoadedLanguages 获取已加载的语言列表
func GetLoadedLanguages() []Language {
	if globalI18n == nil {
		return nil
	}

	globalI18n.mutex.RLock()
	defer globalI18n.mutex.RUnlock()

	var languages []Language
	for lang := range globalI18n.messages {
		languages = append(languages, lang)
	}

	return languages
}
