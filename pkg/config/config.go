package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config 配置结构体
type Config struct {
	Server struct {
		Port           string `yaml:"port"`
		WelcomeMessage string `yaml:"welcome_message"`
		PromptTemplate string `yaml:"prompt_template"`
	} `yaml:"server"`
	Auth struct {
		Password           string   `yaml:"password"`
		LoginPrompt        string   `yaml:"login_prompt"`
		AuthorizedKeys     []string `yaml:"authorized_keys"`      // SSH公钥列表，支持多个
		AuthorizedKeysFile string   `yaml:"authorized_keys_file"` // SSH公钥文件路径（可选）
	} `yaml:"auth"`
	API struct {
		BaseURL      string `yaml:"base_url"`
		APIKey       string `yaml:"api_key"`
		DefaultModel string `yaml:"default_model"`
		Timeout      int    `yaml:"timeout"`
	} `yaml:"api"`
	Display struct {
		LineWidth                 int `yaml:"line_width"`
		ThinkingAnimationInterval int `yaml:"thinking_animation_interval"`
		LoadingAnimationInterval  int `yaml:"loading_animation_interval"`
	} `yaml:"display"`
	Security struct {
		HostKeyFile string `yaml:"host_key_file"`
	} `yaml:"security"`
	I18n struct {
		Language string `yaml:"language"` // 支持的语言: zh-cn, en-us
	} `yaml:"i18n"`
	Prompt struct {
		SystemPrompt    string `yaml:"system_prompt"`    // 系统提示词
		UserPrompt      string `yaml:"user_prompt"`      // 用户消息前缀
		AssistantPrompt string `yaml:"assistant_prompt"` // 助手回复前缀
		StdinPrompt     string `yaml:"stdin_prompt"`     // stdin输入分析提示词
		ExecPrompt      string `yaml:"exec_prompt"`      // exec命令处理提示词
	} `yaml:"prompt"`
}

// GlobalConfig 全局配置实例
var GlobalConfig Config

// Load 加载配置文件
func Load(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	err = yaml.Unmarshal(data, &GlobalConfig)
	if err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	return nil
}

// Get 获取全局配置
func Get() *Config {
	return &GlobalConfig
}
