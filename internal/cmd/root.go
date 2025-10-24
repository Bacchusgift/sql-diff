package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Bacchusgift/sql-diff/internal/ai"
	"github.com/Bacchusgift/sql-diff/internal/config"
	"github.com/Bacchusgift/sql-diff/internal/differ"
	"github.com/Bacchusgift/sql-diff/internal/parser"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// 版本信息
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"

	// 命令行参数
	sourceSQL   string
	targetSQL   string
	enableAI    bool
	configPath  string
	outputFile  string
	interactive bool

	// 颜色输出
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	infoColor    = color.New(color.FgCyan)
	warnColor    = color.New(color.FgYellow)
)

// SetVersion 设置版本信息（由 main 包调用）
func SetVersion(v, bt, gc string) {
	if v != "" {
		version = v
		rootCmd.Version = v // 更新 rootCmd 的版本
	}
	if bt != "" {
		buildTime = bt
	}
	if gc != "" {
		gitCommit = gc
	}
}

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "sql-diff",
	Short: "SQL 表结构比对工具",
	Long: `sql-diff 是一个基于 AST 的 SQL 表结构比对工具。
	
可以比对两个表结构的差异，并自动生成 DDL 补全语句。
支持可选的 AI 智能分析功能，提供优化建议。`,
	Example: `  # 交互式模式（推荐，支持多行粘贴）
  sql-diff -i
  
  # 基础用法
  sql-diff -s "CREATE TABLE users (id INT)" -t "CREATE TABLE users (id INT, name VARCHAR(100))"
  
  # 启用 AI 分析
  sql-diff -s "..." -t "..." --ai
  
  # 交互式 + AI
  sql-diff -i --ai
  
  # 输出到文件
  sql-diff -i -o output.sql`,
	Version: version, // 设置版本号，支持 --version 和 -v
	RunE:    run,
}

// Execute 执行命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// 自定义版本输出模板（简洁版）
	rootCmd.SetVersionTemplate(`{{.Version}}
`)
	
	rootCmd.Flags().StringVarP(&sourceSQL, "source", "s", "", "源表的 CREATE TABLE 语句")
	rootCmd.Flags().StringVarP(&targetSQL, "target", "t", "", "目标表的 CREATE TABLE 语句")
	rootCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "交互式模式（支持多行粘贴）")
	rootCmd.Flags().BoolVar(&enableAI, "ai", false, "启用 AI 智能分析")
	rootCmd.Flags().StringVar(&configPath, "config", ".sql-diff-config.yaml", "配置文件路径")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "输出文件路径（默认输出到控制台）")
	
	// 添加 version 命令（详细版）
	rootCmd.AddCommand(versionCmd)
}

// run 执行主逻辑
func run(cmd *cobra.Command, args []string) error {
	// 交互式模式
	if interactive {
		if err := runInteractive(); err != nil {
			return err
		}
		return nil
	}

	// 命令行参数模式：验证必需参数
	if sourceSQL == "" || targetSQL == "" {
		errorColor.Println("✗ 错误: 必须指定 -s 和 -t 参数，或使用 -i 进入交互式模式")
		fmt.Println()
		fmt.Println("使用方法:")
		fmt.Println("  交互式模式: sql-diff -i")
		fmt.Println("  命令行模式: sql-diff -s \"...\" -t \"...\"")
		return fmt.Errorf("缺少必需参数")
	}

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

	// 验证配置
	if err := cfg.Validate(); err != nil {
		errorColor.Printf("✗ 配置验证失败: %v\n", err)
		return err
	}

	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	infoColor.Println("       SQL 表结构比对工具")
	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 解析源表结构
	infoColor.Println("📖 正在解析源表结构...")
	p := parser.NewParser()
	sourceSchema, err := p.Parse(sourceSQL)
	if err != nil {
		errorColor.Printf("✗ 解析源表失败: %v\n", err)
		return err
	}
	successColor.Printf("✓ 源表: %s (%d 列)\n", sourceSchema.Name, len(sourceSchema.Columns))
	fmt.Println()

	// 解析目标表结构
	infoColor.Println("📖 正在解析目标表结构...")
	targetSchema, err := p.Parse(targetSQL)
	if err != nil {
		errorColor.Printf("✗ 解析目标表失败: %v\n", err)
		return err
	}
	successColor.Printf("✓ 目标表: %s (%d 列)\n", targetSchema.Name, len(targetSchema.Columns))
	fmt.Println()

	// 比对差异
	infoColor.Println("🔍 正在比对表结构...")
	d := differ.NewDiffer(sourceSchema, targetSchema)
	diff := d.Compare()

	if !diff.HasChanges() {
		successColor.Println("✓ 两个表结构完全相同，无需修改！")
		return nil
	}

	// 显示差异摘要
	fmt.Println()
	warnColor.Println("📊 差异摘要:")
	warnColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Print(diff.Summary())
	fmt.Println()

	// 生成 DDL
	infoColor.Println("🔧 生成 DDL 语句...")
	ddls := diff.GenerateDDL(sourceSchema.Name)

	fmt.Println()
	successColor.Println("✓ 生成的 DDL 语句:")
	successColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	var output strings.Builder

	// 分类显示 DDL 语句
	addColumns := make([]string, 0)
	modifyColumns := make([]string, 0)
	dropColumns := make([]string, 0)
	addIndexes := make([]string, 0)
	dropIndexes := make([]string, 0)

	for _, ddl := range ddls {
		ddlUpper := strings.ToUpper(ddl)
		if strings.Contains(ddlUpper, "ADD COLUMN") {
			addColumns = append(addColumns, ddl)
		} else if strings.Contains(ddlUpper, "MODIFY COLUMN") {
			modifyColumns = append(modifyColumns, ddl)
		} else if strings.Contains(ddlUpper, "DROP COLUMN") {
			dropColumns = append(dropColumns, ddl)
		} else if strings.Contains(ddlUpper, "ADD INDEX") || strings.Contains(ddlUpper, "ADD UNIQUE") {
			addIndexes = append(addIndexes, ddl)
		} else if strings.Contains(ddlUpper, "DROP INDEX") {
			dropIndexes = append(dropIndexes, ddl)
		}
		output.WriteString(ddl + ";\n")
	}

	// 显示新增列
	if len(addColumns) > 0 {
		color.New(color.FgGreen, color.Bold).Printf("➕ 新增列 (%d):\n", len(addColumns))
		for i, ddl := range addColumns {
			color.New(color.FgGreen).Printf("  %d. %s;\n", i+1, ddl)
		}
		fmt.Println()
	}

	// 显示修改列
	if len(modifyColumns) > 0 {
		color.New(color.FgYellow, color.Bold).Printf("🔄 修改列 (%d):\n", len(modifyColumns))
		for i, ddl := range modifyColumns {
			color.New(color.FgYellow).Printf("  %d. %s;\n", i+1, ddl)
		}
		fmt.Println()
	}

	// 显示删除列（注释）
	if len(dropColumns) > 0 {
		color.New(color.FgRed, color.Bold).Printf("🗑️  删除列 (%d) [已注释]:\n", len(dropColumns))
		for i, ddl := range dropColumns {
			color.New(color.FgRed).Printf("  %d. %s;\n", i+1, ddl)
		}
		fmt.Println()
	}

	// 显示新增索引
	if len(addIndexes) > 0 {
		color.New(color.FgCyan, color.Bold).Printf("📇 新增索引 (%d):\n", len(addIndexes))
		for i, ddl := range addIndexes {
			color.New(color.FgCyan).Printf("  %d. %s;\n", i+1, ddl)
		}
		fmt.Println()
	}

	// 显示删除索引（注释）
	if len(dropIndexes) > 0 {
		color.New(color.FgMagenta, color.Bold).Printf("🗂️  删除索引 (%d) [已注释]:\n", len(dropIndexes))
		for i, ddl := range dropIndexes {
			color.New(color.FgMagenta).Printf("  %d. %s;\n", i+1, ddl)
		}
		fmt.Println()
	}

	// 显示完整的可执行 SQL
	if len(ddls) > 0 {
		color.New(color.FgWhite, color.Bold).Println("📋 完整执行脚本:")
		color.New(color.FgWhite, color.Bold).Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		for _, ddl := range ddls {
			fmt.Println(ddl + ";")
		}
		fmt.Println()
	}

	// AI 分析
	if cfg.AI.Enabled {
		fmt.Println()
		infoColor.Println("🤖 正在进行 AI 智能分析...")

		provider, err := ai.NewProvider(&cfg.AI)
		if err != nil {
			warnColor.Printf("⚠ AI 初始化失败: %v\n", err)
		} else {
			result, err := provider.Analyze(sourceSQL, targetSQL, diff.Summary())
			if err != nil {
				warnColor.Printf("⚠ AI 分析失败: %v\n", err)
			} else {
				fmt.Println()
				infoColor.Println("💡 AI 分析结果:")
				infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

				// 显示摘要
				if result.Summary != "" {
					fmt.Println()
					color.New(color.FgWhite, color.Bold).Println("📊 差异分析:")
					fmt.Println(result.Summary)
				}

				// 显示优化建议
				if len(result.Suggestions) > 0 {
					fmt.Println()
					color.New(color.FgGreen, color.Bold).Println("✨ 优化建议:")
					for i, suggestion := range result.Suggestions {
						fmt.Printf("  %d. %s\n", i+1, suggestion)
					}
				}

				// 显示潜在风险
				if len(result.Risks) > 0 {
					fmt.Println()
					color.New(color.FgRed, color.Bold).Println("⚠️  潜在风险:")
					for i, risk := range result.Risks {
						fmt.Printf("  %d. %s\n", i+1, risk)
					}
				}

				// 显示最佳实践
				if len(result.BestPractice) > 0 {
					fmt.Println()
					color.New(color.FgBlue, color.Bold).Println("📖 最佳实践:")
					for i, practice := range result.BestPractice {
						fmt.Printf("  %d. %s\n", i+1, practice)
					}
				}
			}
		}
	}

	// 输出到文件
	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(output.String()), 0644); err != nil {
			errorColor.Printf("✗ 写入文件失败: %v\n", err)
			return err
		}
		successColor.Printf("✓ DDL 已保存到: %s\n", outputFile)
	}

	fmt.Println()
	successColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	successColor.Println("           完成！")
	successColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	return nil
}

// runInteractive 交互式模式
func runInteractive() error {
	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	infoColor.Println("       SQL 表结构比对工具 - 交互式模式")
	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 提前加载配置，显示 AI 状态
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		errorColor.Printf("✗ 加载配置失败: %v\n", err)
		return err
	}

	// 如果命令行指定了 --ai，覆盖配置文件
	if enableAI {
		cfg.AI.Enabled = true
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		errorColor.Printf("✗ 配置验证失败: %v\n", err)
		return err
	}

	// 显示 AI 配置状态
	if cfg.AI.Enabled {
		color.New(color.FgGreen, color.Bold).Println("🤖 AI 智能分析: 已启用")
		color.New(color.FgGreen).Printf("   提供商: %s\n", cfg.AI.Provider)
		if cfg.AI.Model != "" {
			color.New(color.FgGreen).Printf("   模型: %s\n", cfg.AI.Model)
		}
	} else {
		color.New(color.FgWhite).Println("ℹ️  AI 智能分析: 未启用")
		color.New(color.FgWhite).Println("   （可通过 --ai 参数或配置文件启用）")
	}
	fmt.Println()

	// 读取源表 SQL
	color.New(color.FgYellow, color.Bold).Println("📋 请粘贴源表的 CREATE TABLE 语句：")
	color.New(color.FgWhite).Println("（直接粘贴完整 SQL，粘贴完成后输入 'END' 或连续按两次 Enter）")
	color.New(color.FgCyan).Println("（提示：建议在文本编辑器中准备好 SQL，然后直接粘贴）")
	fmt.Println()

	sourceSQL, err := readMultilineInput()
	if err != nil {
		return fmt.Errorf("读取源表 SQL 失败: %v", err)
	}

	if strings.TrimSpace(sourceSQL) == "" {
		return fmt.Errorf("源表 SQL 不能为空")
	}

	successColor.Printf("✓ 已读取 %d 个字符\n", len(sourceSQL))
	fmt.Println()

	// 读取目标表 SQL
	color.New(color.FgYellow, color.Bold).Println("📋 请粘贴目标表的 CREATE TABLE 语句：")
	color.New(color.FgWhite).Println("（直接粘贴完整 SQL，粘贴完成后输入 'END' 或连续按两次 Enter）")
	color.New(color.FgCyan).Println("（提示：建议在文本编辑器中准备好 SQL，然后直接粘贴）")
	fmt.Println()

	targetSQL, err := readMultilineInput()
	if err != nil {
		return fmt.Errorf("读取目标表 SQL 失败: %v", err)
	}

	if strings.TrimSpace(targetSQL) == "" {
		return fmt.Errorf("目标表 SQL 不能为空")
	}

	successColor.Printf("✓ 已读取 %d 个字符\n", len(targetSQL))
	fmt.Println()

	// 调用核心比对逻辑
	return processComparison(sourceSQL, targetSQL, cfg)
}

// readMultilineInput 从标准输入读取多行文本
// 支持粘贴多行 SQL，使用特殊标记结束输入
func readMultilineInput() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	emptyLineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// 方式1: 单独一行输入 END 结束（推荐）
		if trimmedLine == "END" {
			break
		}

		// 方式2: 连续两个空行结束（但保留 SQL 中的单个空行）
		if trimmedLine == "" {
			emptyLineCount++
			if emptyLineCount >= 2 {
				// 移除最后的空行
				if len(lines) > 0 {
					lines = lines[:len(lines)-1]
				}
				break
			}
		} else {
			emptyLineCount = 0
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return "", err
	}

	result := strings.Join(lines, "\n")
	return strings.TrimSpace(result), nil
}

// processComparison 执行 SQL 比对逻辑
func processComparison(sourceSQL, targetSQL string, cfg *config.Config) error {
	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	infoColor.Println("       开始比对")
	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 解析源表结构
	infoColor.Println("📖 正在解析源表结构...")
	p := parser.NewParser()
	sourceSchema, err := p.Parse(sourceSQL)
	if err != nil {
		errorColor.Printf("✗ 解析源表失败: %v\n", err)
		return err
	}
	successColor.Printf("✓ 源表: %s (%d 列)\n", sourceSchema.Name, len(sourceSchema.Columns))
	fmt.Println()

	// 解析目标表结构
	infoColor.Println("📖 正在解析目标表结构...")
	targetSchema, err := p.Parse(targetSQL)
	if err != nil {
		errorColor.Printf("✗ 解析目标表失败: %v\n", err)
		return err
	}
	successColor.Printf("✓ 目标表: %s (%d 列)\n", targetSchema.Name, len(targetSchema.Columns))
	fmt.Println()

	// 比对差异
	infoColor.Println("🔍 正在比对表结构...")
	d := differ.NewDiffer(sourceSchema, targetSchema)
	diff := d.Compare()

	if !diff.HasChanges() {
		successColor.Println("✓ 两个表结构完全相同，无需修改！")
		return nil
	}

	// 显示差异摘要
	fmt.Println()
	warnColor.Println("📊 差异摘要:")
	warnColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Print(diff.Summary())
	fmt.Println()

	// 生成 DDL
	infoColor.Println("🔧 生成 DDL 语句...")
	ddls := diff.GenerateDDL(sourceSchema.Name)

	fmt.Println()
	successColor.Println("✓ 生成的 DDL 语句:")
	successColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	var output strings.Builder

	// 分类显示 DDL 语句
	addColumns := make([]string, 0)
	modifyColumns := make([]string, 0)
	dropColumns := make([]string, 0)
	addIndexes := make([]string, 0)
	dropIndexes := make([]string, 0)

	for _, ddl := range ddls {
		ddlUpper := strings.ToUpper(ddl)
		if strings.Contains(ddlUpper, "ADD COLUMN") {
			addColumns = append(addColumns, ddl)
		} else if strings.Contains(ddlUpper, "MODIFY COLUMN") {
			modifyColumns = append(modifyColumns, ddl)
		} else if strings.Contains(ddlUpper, "DROP COLUMN") {
			dropColumns = append(dropColumns, ddl)
		} else if strings.Contains(ddlUpper, "ADD INDEX") || strings.Contains(ddlUpper, "ADD UNIQUE") {
			addIndexes = append(addIndexes, ddl)
		} else if strings.Contains(ddlUpper, "DROP INDEX") {
			dropIndexes = append(dropIndexes, ddl)
		}
		output.WriteString(ddl + ";\n")
	}

	// 显示新增列
	if len(addColumns) > 0 {
		color.New(color.FgGreen, color.Bold).Printf("➕ 新增列 (%d):\n", len(addColumns))
		for i, ddl := range addColumns {
			color.New(color.FgGreen).Printf("  %d. %s;\n", i+1, ddl)
		}
		fmt.Println()
	}

	// 显示修改列
	if len(modifyColumns) > 0 {
		color.New(color.FgYellow, color.Bold).Printf("🔄 修改列 (%d):\n", len(modifyColumns))
		for i, ddl := range modifyColumns {
			color.New(color.FgYellow).Printf("  %d. %s;\n", i+1, ddl)
		}
		fmt.Println()
	}

	// 显示删除列（注释）
	if len(dropColumns) > 0 {
		color.New(color.FgRed, color.Bold).Printf("🗑️  删除列 (%d) [已注释]:\n", len(dropColumns))
		for i, ddl := range dropColumns {
			color.New(color.FgRed).Printf("  %d. %s;\n", i+1, ddl)
		}
		fmt.Println()
	}

	// 显示新增索引
	if len(addIndexes) > 0 {
		color.New(color.FgCyan, color.Bold).Printf("📇 新增索引 (%d):\n", len(addIndexes))
		for i, ddl := range addIndexes {
			color.New(color.FgCyan).Printf("  %d. %s;\n", i+1, ddl)
		}
		fmt.Println()
	}

	// 显示删除索引（注释）
	if len(dropIndexes) > 0 {
		color.New(color.FgMagenta, color.Bold).Printf("🗂️  删除索引 (%d) [已注释]:\n", len(dropIndexes))
		for i, ddl := range dropIndexes {
			color.New(color.FgMagenta).Printf("  %d. %s;\n", i+1, ddl)
		}
		fmt.Println()
	}

	// 显示完整的可执行 SQL
	if len(ddls) > 0 {
		color.New(color.FgWhite, color.Bold).Println("📋 完整执行脚本:")
		color.New(color.FgWhite, color.Bold).Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		for _, ddl := range ddls {
			fmt.Println(ddl + ";")
		}
		fmt.Println()
	}

	// AI 分析
	if cfg.AI.Enabled {
		fmt.Println()
		infoColor.Println("🤖 正在进行 AI 智能分析...")

		provider, err := ai.NewProvider(&cfg.AI)
		if err != nil {
			warnColor.Printf("⚠ AI 初始化失败: %v\n", err)
		} else {
			result, err := provider.Analyze(sourceSQL, targetSQL, diff.Summary())
			if err != nil {
				warnColor.Printf("⚠ AI 分析失败: %v\n", err)
			} else {
				fmt.Println()
				infoColor.Println("💡 AI 分析结果:")
				infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

				// 显示摘要
				if result.Summary != "" {
					fmt.Println()
					color.New(color.FgWhite, color.Bold).Println("📊 差异分析:")
					fmt.Println(result.Summary)
				}

				// 显示优化建议
				if len(result.Suggestions) > 0 {
					fmt.Println()
					color.New(color.FgGreen, color.Bold).Println("✨ 优化建议:")
					for i, suggestion := range result.Suggestions {
						fmt.Printf("  %d. %s\n", i+1, suggestion)
					}
				}

				// 显示潜在风险
				if len(result.Risks) > 0 {
					fmt.Println()
					color.New(color.FgRed, color.Bold).Println("⚠️  潜在风险:")
					for i, risk := range result.Risks {
						fmt.Printf("  %d. %s\n", i+1, risk)
					}
				}

				// 显示最佳实践
				if len(result.BestPractice) > 0 {
					fmt.Println()
					color.New(color.FgBlue, color.Bold).Println("📖 最佳实践:")
					for i, practice := range result.BestPractice {
						fmt.Printf("  %d. %s\n", i+1, practice)
					}
				}
			}
		}
	}

	// 输出到文件
	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(output.String()), 0644); err != nil {
			errorColor.Printf("✗ 写入文件失败: %v\n", err)
			return err
		}
		successColor.Printf("✓ DDL 已保存到: %s\n", outputFile)
	}

	fmt.Println()
	successColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	successColor.Println("           完成！")
	successColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	return nil
}
