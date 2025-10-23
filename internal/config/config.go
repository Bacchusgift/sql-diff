package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config 应用配置结构
type Config struct {
	AI AIConfig `yaml:"ai"`
}

// AIConfig AI 相关配置
type AIConfig struct {
	Enabled     bool   `yaml:"enabled"`      // 是否启用 AI 功能
	Provider    string `yaml:"provider"`     // AI 提供商：deepseek, openai, custom
	APIKey      string `yaml:"api_key"`      // API 密钥
	APIEndpoint string `yaml:"api_endpoint"` // API 端点
	Model       string `yaml:"model"`        // 使用的模型
	Timeout     int    `yaml:"timeout"`      // 请求超时时间（秒）
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		AI: AIConfig{
			Enabled:     false,
			Provider:    "deepseek",
			APIEndpoint: "https://api.deepseek.com/v1",
			Model:       "deepseek-chat",
			Timeout:     30,
		},
	}
}

// LoadConfig 从环境变量或配置文件加载配置
// 优先级: 环境变量 > 配置文件 > 默认值
func LoadConfig(path string) (*Config, error) {
	// 从默认配置开始
	config := DefaultConfig()

	// 尝试从配置文件加载（如果存在）
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}

		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("解析配置文件失败: %w", err)
		}
	}

	// 从环境变量覆盖配置（优先级最高）
	loadFromEnv(config)

	return config, nil
}

// loadFromEnv 从环境变量加载配置
func loadFromEnv(config *Config) {
	// AI 功能启用状态
	if enabled := os.Getenv("SQL_DIFF_AI_ENABLED"); enabled != "" {
		config.AI.Enabled = enabled == "true" || enabled == "1"
	}

	// AI 提供商
	if provider := os.Getenv("SQL_DIFF_AI_PROVIDER"); provider != "" {
		config.AI.Provider = provider
	}

	// API Key
	if apiKey := os.Getenv("SQL_DIFF_AI_API_KEY"); apiKey != "" {
		config.AI.APIKey = apiKey
	}

	// API 端点
	if endpoint := os.Getenv("SQL_DIFF_AI_ENDPOINT"); endpoint != "" {
		config.AI.APIEndpoint = endpoint
	}

	// 模型
	if model := os.Getenv("SQL_DIFF_AI_MODEL"); model != "" {
		config.AI.Model = model
	}

	// 超时时间
	if timeout := os.Getenv("SQL_DIFF_AI_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			config.AI.Timeout = t
		}
	}
}

// SaveToEnv 保存配置到环境变量（生成 export 命令）
func (c *Config) SaveToEnv() []string {
	exports := []string{}

	if c.AI.Enabled {
		exports = append(exports, "export SQL_DIFF_AI_ENABLED=true")
	} else {
		exports = append(exports, "export SQL_DIFF_AI_ENABLED=false")
	}

	if c.AI.Provider != "" {
		exports = append(exports, fmt.Sprintf("export SQL_DIFF_AI_PROVIDER=%s", c.AI.Provider))
	}

	if c.AI.APIKey != "" {
		exports = append(exports, fmt.Sprintf("export SQL_DIFF_AI_API_KEY=%s", c.AI.APIKey))
	}

	if c.AI.APIEndpoint != "" {
		exports = append(exports, fmt.Sprintf("export SQL_DIFF_AI_ENDPOINT=%s", c.AI.APIEndpoint))
	}

	if c.AI.Model != "" {
		exports = append(exports, fmt.Sprintf("export SQL_DIFF_AI_MODEL=%s", c.AI.Model))
	}

	if c.AI.Timeout > 0 {
		exports = append(exports, fmt.Sprintf("export SQL_DIFF_AI_TIMEOUT=%d", c.AI.Timeout))
	}

	return exports
}

// GetEnvVars 获取所有环境变量名称
func GetEnvVars() []string {
	return []string{
		"SQL_DIFF_AI_ENABLED",
		"SQL_DIFF_AI_PROVIDER",
		"SQL_DIFF_AI_API_KEY",
		"SQL_DIFF_AI_ENDPOINT",
		"SQL_DIFF_AI_MODEL",
		"SQL_DIFF_AI_TIMEOUT",
	}
}

// Validate 验证配置是否有效
func (c *Config) Validate() error {
	if c.AI.Enabled {
		if c.AI.APIKey == "" {
			return fmt.Errorf("AI 功能已启用，但未配置 API Key")
		}
		if c.AI.APIEndpoint == "" {
			return fmt.Errorf("AI 功能已启用，但未配置 API Endpoint")
		}
	}
	return nil
}
