package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Bacchusgift/sql-diff/internal/ai"
	"github.com/Bacchusgift/sql-diff/internal/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	alterTable      string
	alterDesc       string
	alterOutput     string
	alterInteractive bool
)

// alterCmd 生成 ALTER TABLE 命令
var alterCmd = &cobra.Command{
	Use:   "alter",
	Short: "根据自然语言描述生成 ALTER TABLE 语句",
	Long: `使用 AI 根据现有表结构和自然语言描述生成 MySQL ALTER TABLE 语句。

示例：
  # 命令行模式
  sql-diff alter -t "CREATE TABLE users ..." -d "添加手机号字段、邮箱改为唯一索引"
  
  # 交互式模式
  sql-diff alter -i -d "添加商品状态字段，默认值为上架"`,
	RunE: runAlter,
}

func init() {
	rootCmd.AddCommand(alterCmd)
	alterCmd.Flags().StringVarP(&alterTable, "table", "t", "", "现有表的 CREATE TABLE 语句")
	alterCmd.Flags().StringVarP(&alterDesc, "description", "d", "", "修改需求的自然语言描述（必需）")
	alterCmd.Flags().StringVarP(&alterOutput, "output", "o", "", "输出文件路径（可选）")
	alterCmd.Flags().BoolVarP(&alterInteractive, "interactive", "i", false, "交互式输入表结构")
	alterCmd.MarkFlagRequired("description")
}

func runAlter(cmd *cobra.Command, args []string) error {
	// 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		errorColor.Printf("✗ 加载配置失败: %v\n", err)
		return err
	}

	// 如果命令行指定了 --ai，覆盖配置文件
	if enableAI {
		cfg.AI.Enabled = true
	}

	// 检查 AI 是否启用
	if !cfg.AI.Enabled {
		errorColor.Println("✗ 该功能需要启用 AI 功能")
		fmt.Println()
		fmt.Println("请通过以下方式之一启用 AI：")
		fmt.Println("  1. 使用 --ai 参数: sql-diff alter --ai -d \"...\"")
		fmt.Println("  2. 在配置文件中启用: .sql-diff-config.yaml")
		return fmt.Errorf("AI 功能未启用")
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		errorColor.Printf("✗ 配置验证失败: %v\n", err)
		return err
	}

	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	infoColor.Println("       AI 生成 ALTER TABLE 语句")
	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 获取表结构
	var currentDDL string
	if alterInteractive {
		color.New(color.FgYellow, color.Bold).Println("📋 请粘贴现有表的 CREATE TABLE 语句：")
		color.New(color.FgWhite).Println("（粘贴完成后输入 'END' 或连续按两次 Enter）")
		fmt.Println()

		ddl, err := readMultilineInput()
		if err != nil {
			return fmt.Errorf("读取表结构失败: %v", err)
		}

		if strings.TrimSpace(ddl) == "" {
			return fmt.Errorf("表结构不能为空")
		}

		currentDDL = ddl
		successColor.Printf("✓ 已读取 %d 个字符\n", len(currentDDL))
		fmt.Println()
	} else {
		if alterTable == "" {
			errorColor.Println("✗ 必须指定 -t 参数或使用 -i 交互式输入表结构")
			return fmt.Errorf("缺少表结构")
		}
		currentDDL = alterTable
	}

	color.New(color.FgCyan).Printf("📝 修改需求: %s\n", alterDesc)
	fmt.Println()

	infoColor.Println("🤖 正在使用 AI 生成 SQL...")
	
	// 创建 AI Provider
	provider, err := ai.NewProvider(&cfg.AI)
	if err != nil {
		errorColor.Printf("✗ AI 初始化失败: %v\n", err)
		return err
	}

	// 调用 AI 生成 SQL
	sql, err := provider.GenerateAlterTable(currentDDL, alterDesc)
	if err != nil {
		errorColor.Printf("✗ 生成失败: %v\n", err)
		return err
	}

	// 显示结果
	fmt.Println()
	successColor.Println("✓ 生成成功！")
	successColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	
	color.New(color.FgWhite, color.Bold).Println("📋 生成的 ALTER TABLE 语句:")
	fmt.Println()
	
	// 处理多条 SQL 语句
	sqlStatements := strings.Split(sql, "\n")
	for _, stmt := range sqlStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			fmt.Println(stmt + ";")
		}
	}
	fmt.Println()

	// 输出到文件
	if alterOutput != "" {
		var content strings.Builder
		for _, stmt := range sqlStatements {
			stmt = strings.TrimSpace(stmt)
			if stmt != "" {
				content.WriteString(stmt + ";\n")
			}
		}
		
		if err := os.WriteFile(alterOutput, []byte(content.String()), 0644); err != nil {
			errorColor.Printf("✗ 写入文件失败: %v\n", err)
			return err
		}
		successColor.Printf("✓ SQL 已保存到: %s\n", alterOutput)
	}

	successColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	return nil
}
