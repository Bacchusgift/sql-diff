package cmd

import (
	"fmt"
	"os"

	"github.com/Bacchusgift/sql-diff/internal/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// 配置命令参数
	aiEnabled  bool
	aiProvider string
	aiAPIKey   string
	aiEndpoint string
	aiModel    string
	aiTimeout  int
	showEnv    bool
	quietMode  bool // 静默模式，只输出 export 命令
)

// configCmd 配置管理命令
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "配置 AI 参数",
	Long: `配置 SQL-Diff 的 AI 功能参数，保存到环境变量中。

配置优先级: 环境变量 > 配置文件 > 默认值

支持的环境变量:
  SQL_DIFF_AI_ENABLED   - 是否启用 AI (true/false)
  SQL_DIFF_AI_PROVIDER  - AI 提供商 (deepseek/openai)
  SQL_DIFF_AI_API_KEY   - API 密钥
  SQL_DIFF_AI_ENDPOINT  - API 端点
  SQL_DIFF_AI_MODEL     - 使用的模型
  SQL_DIFF_AI_TIMEOUT   - 超时时间（秒）`,
	Example: `  # 配置 DeepSeek API
  sql-diff config --ai-enabled --provider deepseek --api-key sk-xxx

  # 查看当前环境变量配置
  sql-diff config --show

  # 生成 export 命令
  sql-diff config --ai-enabled --provider deepseek --api-key sk-xxx > ~/.sql-diff-env
  source ~/.sql-diff-env`,
	RunE: runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().BoolVar(&aiEnabled, "ai-enabled", false, "启用 AI 功能")
	configCmd.Flags().StringVar(&aiProvider, "provider", "", "AI 提供商 (deepseek/openai)")
	configCmd.Flags().StringVar(&aiAPIKey, "api-key", "", "API 密钥")
	configCmd.Flags().StringVar(&aiEndpoint, "endpoint", "", "API 端点")
	configCmd.Flags().StringVar(&aiModel, "model", "", "使用的模型")
	configCmd.Flags().IntVar(&aiTimeout, "timeout", 0, "超时时间（秒）")
	configCmd.Flags().BoolVar(&showEnv, "show", false, "显示当前环境变量配置")
	configCmd.Flags().BoolVarP(&quietMode, "quiet", "q", false, "静默模式，只输出 export 命令")
}

func runConfig(cmd *cobra.Command, args []string) error {
	successColor := color.New(color.FgGreen, color.Bold)
	infoColor := color.New(color.FgCyan)
	warnColor := color.New(color.FgYellow)

	// 如果是显示当前配置
	if showEnv {
		return showCurrentConfig()
	}

	// 检查是否有任何配置参数
	hasAnyFlag := cmd.Flags().Changed("ai-enabled") ||
		cmd.Flags().Changed("provider") ||
		cmd.Flags().Changed("api-key") ||
		cmd.Flags().Changed("endpoint") ||
		cmd.Flags().Changed("model") ||
		cmd.Flags().Changed("timeout")

	if !hasAnyFlag {
		// 没有任何参数，显示帮助
		return cmd.Help()
	}

	// 构建配置
	cfg := config.DefaultConfig()

	// 应用用户指定的配置
	if cmd.Flags().Changed("ai-enabled") {
		cfg.AI.Enabled = aiEnabled
	}
	if aiProvider != "" {
		cfg.AI.Provider = aiProvider
	}
	if aiAPIKey != "" {
		cfg.AI.APIKey = aiAPIKey
	}
	if aiEndpoint != "" {
		cfg.AI.APIEndpoint = aiEndpoint
	}
	if aiModel != "" {
		cfg.AI.Model = aiModel
	}
	if aiTimeout > 0 {
		cfg.AI.Timeout = aiTimeout
	}

	// 生成 export 命令
	exports := cfg.SaveToEnv()

	// 静默模式：只输出 export 命令，方便重定向
	if quietMode {
		for _, exp := range exports {
			fmt.Println(exp)
		}
		return nil
	}

	// 正常模式：显示完整的提示信息
	fmt.Println()
	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	infoColor.Println("       配置环境变量")
	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	successColor.Println("✓ 生成的环境变量配置:")
	fmt.Println()

	for _, exp := range exports {
		fmt.Println(exp)
	}

	fmt.Println()
	warnColor.Println("💡 使用方法:")
	fmt.Println()
	fmt.Println("  方法 1 - 自动保存到 ~/.bashrc (推荐):")
	fmt.Println("    sql-diff config --ai-enabled --provider deepseek --api-key YOUR_KEY -q >> ~/.bashrc")
	fmt.Println("    source ~/.bashrc")
	fmt.Println()
	fmt.Println("  方法 2 - 保存到独立文件:")
	fmt.Println("    sql-diff config --ai-enabled --provider deepseek --api-key YOUR_KEY -q > ~/.sql-diff-env")
	fmt.Println("    echo 'source ~/.sql-diff-env' >> ~/.bashrc")
	fmt.Println("    source ~/.bashrc")
	fmt.Println()
	fmt.Println("  方法 3 - 临时使用(当前会话):")
	fmt.Println("    eval \"$(sql-diff config --ai-enabled --provider deepseek --api-key YOUR_KEY -q)\"")
	fmt.Println()
	fmt.Println("  验证配置:")
	fmt.Println("    sql-diff config --show")
	fmt.Println()

	return nil
}

func showCurrentConfig() error {
	infoColor := color.New(color.FgCyan)
	successColor := color.New(color.FgGreen)
	warnColor := color.New(color.FgYellow)

	fmt.Println()
	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	infoColor.Println("       当前环境变量配置")
	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	envVars := config.GetEnvVars()
	hasConfig := false

	for _, envVar := range envVars {
		value := os.Getenv(envVar)
		if value != "" {
			hasConfig = true
			// 隐藏 API Key 的部分内容
			if envVar == "SQL_DIFF_AI_API_KEY" && len(value) > 10 {
				value = value[:6] + "..." + value[len(value)-4:]
			}
			successColor.Printf("✓ %s = %s\n", envVar, value)
		} else {
			warnColor.Printf("  %s = (未设置)\n", envVar)
		}
	}

	fmt.Println()
	if hasConfig {
		successColor.Println("✓ 已检测到环境变量配置")
	} else {
		warnColor.Println("⚠ 未检测到环境变量配置")
		fmt.Println()
		fmt.Println("使用以下命令配置:")
		fmt.Println("  sql-diff config --ai-enabled --provider deepseek --api-key YOUR_KEY")
	}
	fmt.Println()

	// 加载并显示完整配置（环境变量 + 配置文件）
	cfg, err := config.LoadConfig(".sql-diff-config.yaml")
	if err != nil {
		return err
	}

	fmt.Println()
	infoColor.Println("📋 最终生效的配置:")
	fmt.Println()
	fmt.Printf("  AI 启用状态: %v\n", cfg.AI.Enabled)
	fmt.Printf("  AI 提供商:   %s\n", cfg.AI.Provider)

	if cfg.AI.APIKey != "" {
		maskedKey := cfg.AI.APIKey[:6] + "..." + cfg.AI.APIKey[len(cfg.AI.APIKey)-4:]
		fmt.Printf("  API Key:     %s\n", maskedKey)
	} else {
		fmt.Printf("  API Key:     (未设置)\n")
	}

	fmt.Printf("  API 端点:    %s\n", cfg.AI.APIEndpoint)
	fmt.Printf("  模型:        %s\n", cfg.AI.Model)
	fmt.Printf("  超时时间:    %d 秒\n", cfg.AI.Timeout)
	fmt.Println()

	return nil
}
