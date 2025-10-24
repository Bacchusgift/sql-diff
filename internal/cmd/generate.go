package cmd

import (
	"fmt"
	"os"

	"github.com/Bacchusgift/sql-diff/internal/ai"
	"github.com/Bacchusgift/sql-diff/internal/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	generateDesc   string
	generateOutput string
)

// generateCmd 生成 CREATE TABLE 命令
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "根据自然语言描述生成 CREATE TABLE 语句",
	Long: `使用 AI 根据自然语言描述生成标准的 MySQL CREATE TABLE 语句。

功能特点：
  ✓ 自动推断字段类型（VARCHAR、INT、DECIMAL、DATETIME 等）
  ✓ 自动添加主键、索引、唯一约束
  ✓ 应用 MySQL 最佳实践（InnoDB、UTF8MB4、注释等）
  ✓ 使用标准命名规范（snake_case）
  ✓ 支持输出到文件

注意：此功能需要启用 AI（配置 .sql-diff-config.yaml 或使用 --ai 参数）`,
	Example: `  # 基础用法
  sql-diff generate -d "创建用户表，包含 ID、用户名、邮箱、密码、创建时间"
  
  # 复杂示例
  sql-diff generate -d "创建订单表：订单号（唯一）、用户ID（外键）、订单金额（精确到分）、订单状态、下单时间"
  
  # 输出到文件
  sql-diff generate -d "创建商品表：商品ID、名称、价格、库存" -o product.sql
  
  # 启用 AI（如果配置文件中未启用）
  sql-diff generate --ai -d "创建博客文章表"`,
	RunE: runGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&generateDesc, "description", "d", "", "表结构的自然语言描述（必需）")
	generateCmd.Flags().StringVarP(&generateOutput, "output", "o", "", "输出文件路径（可选，默认输出到控制台）")
	generateCmd.MarkFlagRequired("description")
}

func runGenerate(cmd *cobra.Command, args []string) error {
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
		fmt.Println("  1. 使用 --ai 参数: sql-diff generate --ai -d \"...\"")
		fmt.Println("  2. 在配置文件中启用: .sql-diff-config.yaml")
		return fmt.Errorf("AI 功能未启用")
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		errorColor.Printf("✗ 配置验证失败: %v\n", err)
		return err
	}

	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	infoColor.Println("       AI 生成 CREATE TABLE 语句")
	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	color.New(color.FgCyan).Printf("📝 需求描述: %s\n", generateDesc)
	fmt.Println()

	infoColor.Println("🤖 正在使用 AI 生成 SQL...")

	// 创建 AI Provider
	provider, err := ai.NewProvider(&cfg.AI)
	if err != nil {
		errorColor.Printf("✗ AI 初始化失败: %v\n", err)
		return err
	}

	// 调用 AI 生成 SQL
	sql, err := provider.GenerateCreateTable(generateDesc)
	if err != nil {
		errorColor.Printf("✗ 生成失败: %v\n", err)
		return err
	}

	// 显示结果
	fmt.Println()
	successColor.Println("✓ 生成成功！")
	successColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	color.New(color.FgWhite, color.Bold).Println("📋 生成的 CREATE TABLE 语句:")
	fmt.Println()
	fmt.Println(sql + ";")
	fmt.Println()

	// 输出到文件
	if generateOutput != "" {
		content := sql + ";\n"
		if err := os.WriteFile(generateOutput, []byte(content), 0644); err != nil {
			errorColor.Printf("✗ 写入文件失败: %v\n", err)
			return err
		}
		successColor.Printf("✓ SQL 已保存到: %s\n", generateOutput)
	}

	successColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	return nil
}
